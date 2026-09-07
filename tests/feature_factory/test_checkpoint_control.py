"""Deterministic regression tests for the ADH-2026-077 implementation checkpoint.

These tests prove the checkpoint control, guard, and task-only amendment
preflight fail closed exactly as ADH-2026-077 requires, without inspecting or
changing any FEATURE-0018 product semantics. They operate on isolated temporary
copies of the FEATURE-0018 control manifest, state, and spec artifacts so the
real feature files are never mutated.
"""

import importlib.util
import json
import subprocess
import sys
import unittest
from copy import deepcopy
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
sys.dont_write_bytecode = True

CONTROL_PATH = ROOT / ".automation/features/FEATURE-0018.control.json"
STATE_PATH = ROOT / ".automation/state/FEATURE-0018.json"
SPEC = ROOT / ".kiro/specs/governance-iam-approval-exception-foundation"
GUARD = ROOT / "scripts/checkpoint-guard.py"
PREFLIGHT = ROOT / "scripts/checkpoint-task-amendment-preflight.py"
APPROVE = ROOT / "scripts/checkpoint-task-amendment-approve.py"

FROZEN_DESIGN_SHA = "4eb84a0002d8e58de466806b2dcd18085608a0a62d97cc78684b7e98210d92a8"
FROZEN_REQ_SHA = "aed8f347ea154c30292cdd0d8f43841fad6d419b8cf789becacc8b57d9c2a9e4"


def load_module(name: str, path: Path):
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader
    spec.loader.exec_module(module)
    return module


control = load_module("feature_control_ckpt", ROOT / "scripts/feature-control.py")


class CheckpointControlValidationTests(unittest.TestCase):
    """The optional checkpoint block validates fail-closed."""

    def setUp(self):
        self.data = json.loads(CONTROL_PATH.read_text())

    def test_feature_0018_control_with_checkpoint_validates(self):
        control.validate(self.data, expected_feature="FEATURE-0018")

    def test_checkpoint_block_is_present_and_active(self):
        self.assertIn("checkpoint", self.data)
        cp = self.data["checkpoint"]
        self.assertEqual(cp["handoff"], "ADH-2026-077")
        self.assertIs(cp["active"], True)
        self.assertEqual(cp["last_committed_task"], 3)
        self.assertEqual(cp["current_task"], 4)
        self.assertEqual(cp["frozen_design_sha256"], FROZEN_DESIGN_SHA)
        self.assertEqual(cp["frozen_requirements_sha256"], FROZEN_REQ_SHA)

    def test_current_task_must_exceed_last_committed_task(self):
        invalid = deepcopy(self.data)
        invalid["checkpoint"]["current_task"] = 3
        with self.assertRaisesRegex(control.ControlError, "greater than last_committed_task"):
            control.validate(invalid, expected_feature="FEATURE-0018")

    def test_checkpoint_cannot_block_non_ai_stages(self):
        invalid = deepcopy(self.data)
        invalid["checkpoint"]["blocked_stages"] = ["tasks"]
        with self.assertRaisesRegex(control.ControlError, "may only block requirements/design"):
            control.validate(invalid, expected_feature="FEATURE-0018")

    def test_checkpoint_active_must_be_true_when_present(self):
        invalid = deepcopy(self.data)
        invalid["checkpoint"]["active"] = False
        with self.assertRaisesRegex(control.ControlError, "active must be true"):
            control.validate(invalid, expected_feature="FEATURE-0018")

    def test_frozen_hashes_must_be_hex_sha256(self):
        invalid = deepcopy(self.data)
        invalid["checkpoint"]["frozen_design_sha256"] = "not-a-hash"
        with self.assertRaisesRegex(control.ControlError, "frozen_design_sha256"):
            control.validate(invalid, expected_feature="FEATURE-0018")

    def test_reentry_requires_handoff_must_be_true(self):
        invalid = deepcopy(self.data)
        invalid["checkpoint"]["reentry_requires_handoff"] = False
        with self.assertRaisesRegex(control.ControlError, "reentry_requires_handoff"):
            control.validate(invalid, expected_feature="FEATURE-0018")


class CheckpointGuardTests(unittest.TestCase):
    """The deterministic guard fails closed on the prohibited actions."""

    def run_guard(self, *args):
        return subprocess.run(
            [sys.executable, str(GUARD), "--feature", "FEATURE-0018", *args],
            cwd=ROOT,
            text=True,
            capture_output=True,
        )

    def test_status_reports_active_checkpoint(self):
        result = self.run_guard("--mode", "status")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("ACTIVE", result.stdout)
        self.assertIn("ADH-2026-077", result.stdout)

    def test_design_stage_is_blocked(self):
        result = self.run_guard("--mode", "assert-stage-allowed", "--stage", "design")
        self.assertEqual(result.returncode, 3)
        self.assertIn("CHECKPOINT_BLOCK", result.stderr)

    def test_requirements_stage_is_blocked(self):
        result = self.run_guard("--mode", "assert-stage-allowed", "--stage", "requirements")
        self.assertEqual(result.returncode, 3)

    def test_tasks_stage_is_permitted(self):
        result = self.run_guard("--mode", "assert-stage-allowed", "--stage", "tasks")
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_design_review_routing_is_blocked(self):
        result = self.run_guard("--mode", "assert-review-allowed", "--stage", "design")
        self.assertEqual(result.returncode, 3)

    def test_frozen_baseline_is_intact(self):
        result = self.run_guard("--mode", "assert-frozen")
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_cursor_resume_requires_explicit_amendment_approval(self):
        result = self.run_guard("--mode", "assert-cursor-allowed")
        self.assertEqual(result.returncode, 3)
        self.assertIn("checkpoint task-amendment approval", result.stderr)

    def test_reentry_must_name_a_real_and_different_handoff(self):
        result = self.run_guard(
            "--mode", "assert-stage-allowed", "--stage", "design",
            "--reentry-handoff", "ADH-2026-077",
        )
        self.assertEqual(result.returncode, 3)
        self.assertIn("must differ", result.stderr)

    def test_reentry_unknown_handoff_is_rejected(self):
        result = self.run_guard(
            "--mode", "assert-stage-allowed", "--stage", "design",
            "--reentry-handoff", "ADH-2099-999",
        )
        self.assertEqual(result.returncode, 3)
        self.assertIn("not found", result.stderr)


class CheckpointAmendmentPreflightTests(unittest.TestCase):
    """The task-only amendment preflight proves the ADH-2026-077 §3 invariants."""

    def run_preflight(self):
        return subprocess.run(
            [sys.executable, str(PREFLIGHT), "--feature", "FEATURE-0018"],
            cwd=ROOT,
            text=True,
            capture_output=True,
        )

    def test_unchanged_tasks_pass_preflight(self):
        result = self.run_preflight()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("task-only amendment preflight", result.stdout)
        self.assertIn("amended_tasks_sha256=", result.stdout)

    def _mutate_and_restore(self, path: Path, mutate):
        original = path.read_bytes()
        try:
            mutate(path)
            return self.run_preflight()
        finally:
            path.write_bytes(original)

    def test_modifying_committed_task_is_rejected(self):
        def mutate(path: Path):
            text = path.read_text()
            text = text.replace(
                "### Task 3: RoleAssignment domain and workflow triggers",
                "### Task 3: RoleAssignment domain and workflow triggers (TAMPERED)",
            )
            path.write_text(text)

        result = self._mutate_and_restore(SPEC / "tasks.md", mutate)
        self.assertEqual(result.returncode, 1)
        self.assertIn("committed Task 3", result.stderr)

    def test_forward_dependency_is_rejected(self):
        def mutate(path: Path):
            text = path.read_text()
            old = "TASK-F18-07. Internal order: the Task-4-owned"
            new = "TASK-F18-07, TASK-F18-08. Internal order: the Task-4-owned"
            self.assertIn(old, text)
            path.write_text(text.replace(old, new))

        result = self._mutate_and_restore(SPEC / "tasks.md", mutate)
        self.assertEqual(result.returncode, 1)
        self.assertIn("forward dependency", result.stderr)

    def test_missing_writable_path_is_rejected(self):
        import re

        def mutate(path: Path):
            text = path.read_text()
            start = text.index("### Task 4:")
            end = text.index("### Task 5:")
            block = text[start:end]
            block = re.sub(r"`internal/[^`]+`", "(path removed)", block)
            path.write_text(text[:start] + block + text[end:])

        result = self._mutate_and_restore(SPEC / "tasks.md", mutate)
        self.assertEqual(result.returncode, 1)
        self.assertIn("Task 4", result.stderr)

    def test_frozen_design_change_blocks_amendment(self):
        def mutate(path: Path):
            path.write_text(path.read_text() + "\n<!-- tamper -->\n")

        result = self._mutate_and_restore(SPEC / "design.md", mutate)
        self.assertEqual(result.returncode, 1)
        self.assertIn("frozen design/requirements changed", result.stderr)


class CheckpointStateTests(unittest.TestCase):
    """The FEATURE-0018 state records the checkpoint without drift status."""

    def setUp(self):
        self.state = json.loads(STATE_PATH.read_text())

    def test_state_records_active_checkpoint(self):
        # Status and cursor position advance as implementation tasks run; the
        # checkpoint invariant is that the feature remains in cursor mode with
        # a next task beyond its last committed task.
        self.assertEqual(self.state["current_stage"], "cursor")
        self.assertGreater(int(self.state["current_task"]), int(self.state["last_committed_task"]))
        self.assertEqual(self.state["checkpoint_handoff"], "ADH-2026-077")
        self.assertEqual(self.state["frozen_design_sha256"], FROZEN_DESIGN_SHA)
        self.assertEqual(self.state["frozen_requirements_sha256"], FROZEN_REQ_SHA)

    def test_frozen_design_matches_state_hash(self):
        import hashlib

        actual = hashlib.sha256((SPEC / "design.md").read_bytes()).hexdigest()
        self.assertEqual(actual, self.state["frozen_design_sha256"])

    def test_frozen_requirements_matches_state_hash(self):
        import hashlib

        actual = hashlib.sha256((SPEC / "requirements.md").read_bytes()).hexdigest()
        self.assertEqual(actual, self.state["frozen_requirements_sha256"])


class CheckpointAmendmentApprovalTests(unittest.TestCase):
    """A human approval is mandatory and bound to the amended task digest."""

    def run_approval(self, *args):
        return subprocess.run(
            [sys.executable, str(APPROVE), "--feature", "FEATURE-0018", *args],
            cwd=ROOT,
            text=True,
            capture_output=True,
        )

    def _with_state_restored(self, action):
        original = STATE_PATH.read_bytes()
        try:
            return action()
        finally:
            STATE_PATH.write_bytes(original)

    def test_explicit_approve_is_required(self):
        result = self.run_approval("--approved-by", "Test Approver")
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("explicit --approve", result.stderr)

    def test_successful_approval_binds_current_tasks_hash(self):
        def action():
            result = self.run_approval("--approved-by", "Test Approver", "--approve")
            self.assertEqual(result.returncode, 0, result.stderr)
            state = json.loads(STATE_PATH.read_text())
            digest = __import__("hashlib").sha256((SPEC / "tasks.md").read_bytes()).hexdigest()
            self.assertEqual(state["checkpoint_task_amendment_approval_token"], "APPROVED_FOR_CURSOR")
            self.assertEqual(state["checkpoint_task_amendment_tasks_sha256"], digest)
            self.assertEqual(state["tasks_approved_sha256"], digest)
            self.assertEqual(state["current_stage"], "cursor")
        self._with_state_restored(action)

    def test_approval_preserves_later_cursor_position(self):
        def action():
            state = json.loads(STATE_PATH.read_text())
            state.update({"last_committed_task": "9", "current_task": "10"})
            STATE_PATH.write_text(json.dumps(state, indent=2, sort_keys=True) + "\n")
            result = self.run_approval("--approved-by", "Test Approver", "--approve")
            self.assertEqual(result.returncode, 0, result.stderr)
            approved = json.loads(STATE_PATH.read_text())
            self.assertEqual(approved["current_task"], "10")
            self.assertIn("resume_task=10", result.stdout)
        self._with_state_restored(action)

    def test_stale_amendment_approval_blocks_cursor(self):
        def action():
            state = json.loads(STATE_PATH.read_text())
            state.update(
                {
                    "checkpoint_task_amendment_approval_token": "APPROVED_FOR_CURSOR",
                    "checkpoint_task_amendment_tasks_sha256": "0" * 64,
                    "tasks_approval_token": "APPROVED_FOR_CURSOR",
                    "tasks_approved_sha256": "0" * 64,
                }
            )
            STATE_PATH.write_text(json.dumps(state, indent=2, sort_keys=True) + "\n")
            result = subprocess.run(
                [sys.executable, str(GUARD), "--feature", "FEATURE-0018", "--mode", "assert-cursor-allowed"],
                cwd=ROOT,
                text=True,
                capture_output=True,
            )
            self.assertEqual(result.returncode, 3)
            self.assertIn("does not match current tasks.md", result.stderr)
        self._with_state_restored(action)


if __name__ == "__main__":
    unittest.main()
