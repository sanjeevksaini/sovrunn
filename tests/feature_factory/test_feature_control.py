import importlib.util
import json
import tempfile
import unittest
import sys
import subprocess
from copy import deepcopy
from pathlib import Path
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]
FIXTURE = ROOT / "tests/feature_factory/fixtures/FEATURE-0014.control.json"
sys.dont_write_bytecode = True


def load_module(name: str, path: Path):
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader
    spec.loader.exec_module(module)
    return module


control = load_module("feature_control", ROOT / "scripts/feature-control.py")
delta = load_module("semantic_delta", ROOT / "scripts/semantic-delta.py")
closeout = load_module("feature_closeout", ROOT / "scripts/feature-closeout.py")
orchestrator = load_module("feature_orchestrator", ROOT / "scripts/feature-orchestrator.py")
cursor_prompt = load_module("generic_cursor_prompt", ROOT / "scripts/generic-cursor-prompt.py")
receipt = load_module("receipt_check", ROOT / "scripts/receipt-check.py")
task_boundary = load_module(
    "generic_feature_boundary_check", ROOT / "scripts/generic-feature-boundary-check.py"
)
kiro_semantic = load_module(
    "kiro_semantic_check", ROOT / "scripts/kiro-semantic-check.py"
)


class FeatureControlTests(unittest.TestCase):
    def setUp(self):
        self.data = json.loads(FIXTURE.read_text())

    def test_closed_feature_fixture_validates(self):
        control.validate(self.data, expected_feature="FEATURE-0014")

    def test_previous_feature_must_precede_current_feature(self):
        invalid = deepcopy(self.data)
        invalid["ownership"]["previous_features"][0]["id"] = "FEATURE-0015"
        with self.assertRaisesRegex(control.ControlError, "must precede"):
            control.validate(invalid, expected_feature="FEATURE-0014")

    def test_unknown_manifest_keys_are_rejected(self):
        invalid = deepcopy(self.data)
        invalid["execution"]["implicit_approval"] = True
        with self.assertRaisesRegex(control.ControlError, "extra=implicit_approval"):
            control.validate(invalid, expected_feature="FEATURE-0014")

    def test_verification_commands_are_fail_closed(self):
        invalid = deepcopy(self.data)
        invalid["guardrails"]["cursor"]["verify_commands"] = ["sh -c make test"]
        with self.assertRaisesRegex(control.ControlError, "unsafe Cursor verification command"):
            control.validate(invalid, expected_feature="FEATURE-0014")

    def test_context_is_bounded_and_contains_go_guardrails(self):
        manifest = control.context_manifest(self.data, "implementation")
        paths = {item["path"] for item in manifest["files"]}
        self.assertIn("docs/engineering/go-coding-guardrails.md", paths)
        self.assertIn("docs/engineering/go-version-standard.md", paths)
        self.assertIn("docs/architecture/api-resource-standard.md", paths)
        self.assertLessEqual(manifest["total_bytes"], manifest["budget_bytes"])

    def test_prompt_names_dependency_ownership_and_exclusions(self):
        manifest = control.context_manifest(self.data, "requirements")
        prompt = control.prompt_fragment(self.data, manifest, "requirements")
        self.assertIn("FEATURE-0012", prompt)
        self.assertIn("ResourcePool", prompt)
        self.assertIn("ARCHITECTURE_DECISION_REQUIRED", prompt)
        self.assertIn("modify only", prompt)

    def test_review_context_includes_only_existing_stage_predecessors(self):
        spec = self.data["feature"]["spec_path"]
        requirements = f"{spec}/requirements.md"
        design = f"{spec}/design.md"
        tasks = f"{spec}/tasks.md"

        requirements_review = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(self.data, "review", review_stage="requirements")
        }
        self.assertNotIn(requirements, requirements_review)
        self.assertNotIn(design, requirements_review)
        self.assertNotIn(tasks, requirements_review)

        design_review = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(self.data, "review", review_stage="design")
        }
        self.assertIn(requirements, design_review)
        self.assertNotIn(design, design_review)
        self.assertNotIn(tasks, design_review)

        tasks_review = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(self.data, "review", review_stage="tasks")
        }
        self.assertIn(requirements, tasks_review)
        self.assertIn(design, tasks_review)
        self.assertNotIn(tasks, tasks_review)

    def test_review_context_requires_explicit_review_stage(self):
        with self.assertRaisesRegex(control.ControlError, "review context requires"):
            control.resolve_context(self.data, "review")

    def test_closeout_is_never_allowed_to_commit_or_push(self):
        self.assertEqual(self.data["closeout"]["mode"], "prepare_only")
        self.assertIs(self.data["closeout"]["commit"], False)
        self.assertIs(self.data["closeout"]["push"], False)


class SemanticDeltaTests(unittest.TestCase):
    def write(self, directory: Path, name: str, content: str) -> Path:
        path = directory / name
        path.write_text(content)
        return path

    def test_no_change(self):
        with tempfile.TemporaryDirectory() as raw:
            directory = Path(raw)
            baseline = self.write(directory, "old.md", "# A\nMUST remain.\n")
            current = self.write(directory, "new.md", "# A\nMUST remain.\n")
            self.assertEqual(delta.report(baseline, current)["classification"], "NO_CHANGE")

    def test_trailing_whitespace_is_mechanical_only(self):
        with tempfile.TemporaryDirectory() as raw:
            directory = Path(raw)
            baseline = self.write(directory, "old.md", "# A  \nMUST remain.\n")
            current = self.write(directory, "new.md", "# A\nMUST remain.\n")
            self.assertEqual(delta.report(baseline, current)["classification"], "MECHANICAL_ONLY")

    def test_normative_change_requires_semantic_review(self):
        with tempfile.TemporaryDirectory() as raw:
            directory = Path(raw)
            baseline = self.write(directory, "old.md", "# A\nMUST remain.\n")
            current = self.write(directory, "new.md", "# A\nMUST NOT remain.\n")
            report = delta.report(baseline, current)
            self.assertEqual(report["classification"], "SEMANTIC_REVIEW_REQUIRED")
            self.assertTrue(report["semantic_review_required"])


class ScriptSafetyTests(unittest.TestCase):
    task_plan_fixture = """\
### Task 1: Parent implementation task

**Writable paths:**
- `internal/example/service.go`

**Tests:**
- `internal/example/service_test.go`

**Commit message:**
```
feat(example): add service

Long-form commit detail is allowed in the plan.
```

#### Work package 2 (within Task 1): Internal work

**Included writable paths:**
- `internal/example/internal.go`

**Included tests:**
- `internal/example/internal_test.go`

### Task 18: Verification checkpoint

**Writable paths:**
- None (verification only).

**Tests:**
- None (runs existing checks).

**Commit message:** None — verification only.

## 3. No-task ledger
"""

    def test_approved_task_heading_and_sections_are_parsed(self):
        plan = cursor_prompt.task_blocks(self.task_plan_fixture)
        self.assertEqual(sorted(plan), [1, 18])
        self.assertEqual(
            cursor_prompt.writable_paths(plan[1]),
            [
                "internal/example/service.go",
                "internal/example/internal.go",
                "internal/example/service_test.go",
                "internal/example/internal_test.go",
            ],
        )
        self.assertEqual(cursor_prompt.commit_message(plan[1]), "feat(example): add service")

    def test_orchestrator_excludes_non_commit_checkpoint(self):
        plan = orchestrator.blocks(self.task_plan_fixture)
        self.assertEqual(orchestrator.commit_task_ids(plan), [1])
        self.assertEqual(orchestrator.verification_checkpoint_ids(plan), [18])
        self.assertEqual(
            orchestrator.task_writable_paths(plan[1]),
            [
                "internal/example/service.go",
                "internal/example/internal.go",
                "internal/example/service_test.go",
                "internal/example/internal_test.go",
            ],
        )
        self.assertEqual(orchestrator.commit_message(plan[1]), "feat(example): add service")
        self.assertIsNone(orchestrator.commit_message(plan[18]))

    def test_verification_checkpoint_is_not_a_cursor_or_commit_task(self):
        source = (ROOT / "scripts/feature-orchestrator.py").read_text()
        checkpoint = source.index("def run_verification_checkpoint")
        self.assertIn('"make", "ff-feature-gate"', source[checkpoint:])
        self.assertIn("if completed_final_commit:", source)

    def test_interrupted_task_can_resume_only_within_declared_scope(self):
        plan = orchestrator.blocks(self.task_plan_fixture)
        allowed = orchestrator.allowed_resume_paths(
            ROOT / ".automation/state/FEATURE-0099.json", plan, [[1]]
        )
        self.assertIn("internal/example/service.go", allowed)
        self.assertIn("internal/example/internal_test.go", allowed)
        self.assertNotIn("internal/other/escape.go", allowed)

    def test_generic_scripts_compile(self):
        scripts = [
            "feature-control.py",
            "semantic-delta.py",
            "generic-cursor-prompt.py",
            "generic-feature-boundary-check.py",
            "generic-kiro-boundary-check.py",
            "generic-review-prompt.py",
            "feature-orchestrator.py",
            "feature-closeout.py",
            "human-gate.py",
            "executable-plan-report.py",
            "spec-approval-check.py",
            "receipt-check.py",
            "reviewer-openai.py",
            "kiro-semantic-check.py",
        ]
        for script in scripts:
            path = ROOT / "scripts" / script
            compile(path.read_text(), str(path), "exec")

    def test_generic_prompts_do_not_embed_legacy_feature_rules(self):
        prompts = [
            ROOT / "docs/prompts/kiro/generic-requirements.prompt.md",
            ROOT / "docs/prompts/kiro/generic-design.prompt.md",
            ROOT / "docs/prompts/kiro/generic-tasks.prompt.md",
            ROOT / "docs/prompts/cursor/generic-task.prompt.md",
            ROOT / "docs/prompts/reviewer/generic-approval-review.prompt.md",
        ]
        for prompt in prompts:
            text = prompt.read_text()
            self.assertNotIn("FEATURE-0013", text)
            self.assertNotIn("FEATURE-0014", text)

    def test_founder_approved_cursor_default_is_centralized(self):
        policy = json.loads((ROOT / ".automation/model-policy.json").read_text())
        self.assertEqual(policy["default_cursor_profile"], "implementation")
        self.assertEqual(
            policy["tools"]["cursor"]["profiles"]["implementation"][0]["model"],
            "cursor-grok-4.5-high-fast",
        )

    def test_closeout_source_cannot_invoke_commit_or_push(self):
        source = (ROOT / "scripts/feature-closeout.py").read_text()
        self.assertNotIn('["git", "commit"', source)
        self.assertNotIn('["git", "push"', source)

    def test_full_flow_uses_manifest_controlled_cursor_runner(self):
        source = (ROOT / "scripts/feature-flow.sh").read_text()
        self.assertIn('PHASE_BRANCH="$PHASE_BRANCH"', source)
        self.assertIn("make ff-controlled-run", source)
        self.assertNotIn("make ff-task-flow", source)
        self.assertIn("stop_for_gate architecture", source)
        self.assertIn("stop_for_gate executable_plan", source)
        self.assertIn("stop_for_gate final", source)

    def test_spec_flow_applies_pending_revision_before_review(self):
        source = (ROOT / "scripts/spec-flow.sh").read_text()
        pending = source.index('if [[ "$pending_status" == "${stage}_revision_required" ]]')
        review_loop = source.index("while true; do", pending)
        self.assertLess(pending, review_loop)
        self.assertIn("generic-kiro-boundary-check.py", source[pending - 4000 : review_loop])

    def test_spec_flow_runs_semantic_guardrail_before_model_review(self):
        source = (ROOT / "scripts/spec-flow.sh").read_text()
        loop = source.index("while true; do", source.index("run_stage()"))
        semantic = source.index('run_semantic_guardrails "$stage"', loop)
        reviewer = source.index("review-and-route-stage.sh", semantic)
        self.assertLess(semantic, reviewer)
        self.assertIn("${stage}.revision-count", source)

    def test_reviewer_discovery_supports_desktop_bundled_codex(self):
        common = (ROOT / "scripts/common.sh").read_text()
        reviewer = (ROOT / "scripts/reviewer-codex.sh").read_text()
        self.assertIn("/Applications/ChatGPT.app/Contents/Resources/codex", common)
        self.assertIn('CODEX_REVIEWER_BIN="$(resolve_codex_bin)"', common)
        self.assertIn('CODEX_BIN="$(resolve_codex_bin || true)"', reviewer)

    def test_reviewer_discovery_honors_explicit_executable(self):
        with tempfile.TemporaryDirectory() as raw:
            fake = Path(raw) / "codex"
            fake.write_text("#!/bin/sh\nexit 0\n")
            fake.chmod(0o755)
            completed = subprocess.run(
                [
                    "/opt/homebrew/bin/bash",
                    "-c",
                    'source scripts/common.sh; resolve_codex_bin',
                ],
                cwd=ROOT,
                env={"PATH": "/usr/bin:/bin", "CODEX_REVIEWER_BIN": str(fake)},
                check=True,
                capture_output=True,
                text=True,
            )
            self.assertEqual(completed.stdout.strip(), str(fake))


class KiroSemanticGuardrailTests(unittest.TestCase):
    feature_text = """\
## Scope Classification
| ID | Behavior | DEC/ADH | VS0 IDs |
|----|----------|---------|---------|
| REQ-F99-01 | Preserve the approved behavior | DEC-0099 | VS0-STATE-001 |

## Acceptance Criteria
| ID | Criterion | Conformance ID |
|----|-----------|----------------|
| AC-F99-01 | Prove the registered scenario | VS0-CF-F09 |
"""
    conformance = {
        "VS0-CF-F09": {
            "id": "VS0-CF-F09",
            "owner": "FEATURE-0023",
            "inputs": "unavailable-required-participation",
            "expectedState": "denied-placement",
            "expectedError": "VALIDATION_FAILED",
            "expectedSideEffects": "no-silent-substitution",
            "gate": "decision",
        }
    }

    def valid_requirements(self) -> str:
        return """\
## Canonical requirement ledger
| ID | Behavior | DEC/ADH | VS0 IDs |
|----|----------|---------|---------|
| REQ-F99-01 | Preserve the approved behavior | DEC-0099 | VS0-STATE-001 |

## Canonical acceptance ledger
| ID | Criterion | Conformance ID |
|----|-----------|----------------|
| AC-F99-01 | Prove the registered scenario | VS0-CF-F09 |

## Exact conformance semantics ledger
| ID | Owner | Inputs | Expected State | Expected Error | Expected Side Effects | Gate |
|----|-------|--------|----------------|----------------|-----------------------|------|
| VS0-CF-F09 | FEATURE-0023 | unavailable-required-participation | denied-placement | VALIDATION_FAILED | no-silent-substitution | decision |

### REQ-F99-01 — Approved behavior
The behavior remains unchanged. AC-F99-01 uses VS0-CF-F09 only for unavailable-required-participation denial.
"""

    def test_numbered_canonical_section_heading_is_accepted(self):
        body = kiro_semantic.section(
            "## 10. Canonical requirement ledger\n\n| ID | Behavior |\n",
            "Canonical requirement ledger",
        )
        self.assertEqual(body, "\n| ID | Behavior |")

    def test_registry_null_is_rendered_exactly_as_yaml(self):
        self.assertEqual(kiro_semantic.scalar(None), "null")

    def test_exact_ledgers_and_registry_semantics_pass(self):
        errors = []
        kiro_semantic.check_requirements(
            errors,
            "FEATURE-0099",
            self.feature_text,
            self.valid_requirements(),
            self.conformance,
        )
        self.assertEqual(errors, [])

    def test_repurposed_requirement_row_fails(self):
        target = self.valid_requirements().replace(
            "Preserve the approved behavior", "Repurpose the approved behavior", 1
        )
        errors = []
        kiro_semantic.check_requirements(
            errors, "FEATURE-0099", self.feature_text, target, self.conformance
        )
        self.assertTrue(any("approved REQ-F99-01 row" in error for error in errors))

    def test_cross_feature_conformance_cannot_prove_state_transition(self):
        target = self.valid_requirements().replace(
            "The behavior remains unchanged.",
            "An invalid participation transition is rejected by VS0-CF-F09.",
        )
        errors = []
        kiro_semantic.check_requirements(
            errors, "FEATURE-0099", self.feature_text, target, self.conformance
        )
        self.assertTrue(any("transition evidence" in error for error in errors))


class TaskBatchTests(unittest.TestCase):
    def test_run_all_preserves_manifest_sized_batches(self):
        self.assertEqual(
            orchestrator.task_batches([1, 2, 3, 4, 5], None, None, None, 2, True),
            [[1, 2], [3, 4], [5]],
        )

    def test_single_run_keeps_existing_batch_limit(self):
        self.assertEqual(
            orchestrator.task_batches([1, 2, 3, 4], None, None, None, 2, False),
            [[1, 2]],
        )

    def test_resume_starts_after_last_committed_task(self):
        self.assertEqual(
            orchestrator.task_batches([1, 2, 3, 4], "2", None, None, 2, True),
            [[3, 4]],
        )

    def test_completed_plan_is_idempotent_for_run_all(self):
        self.assertEqual(
            orchestrator.task_batches([1, 2], "2", None, None, 2, True),
            [],
        )

    def test_sequence_drift_is_rejected(self):
        with self.assertRaisesRegex(ValueError, "expected Task 2, got 3"):
            orchestrator.task_batches([1, 2, 3], "1", 3, None, 2, True)


class GenericTaskBoundaryTests(unittest.TestCase):
    def test_changed_paths_expands_untracked_directories_to_files(self):
        with patch.object(
            task_boundary.subprocess,
            "run",
            return_value=subprocess.CompletedProcess(
                args=[],
                returncode=0,
                stdout="?? internal/cloudmodel/validate/schema.go\n",
                stderr="",
            ),
        ) as run:
            self.assertEqual(task_boundary.changed_paths(), ["internal/cloudmodel/validate/schema.go"])
        self.assertEqual(
            run.call_args.args[0],
            ["git", "status", "--porcelain=v1", "--untracked-files=all"],
        )


class ReceiptCheckTests(unittest.TestCase):
    def test_ansi_wrapped_receipt_is_accepted_and_diff_preview_is_ignored(self):
        raw = (
            "\x1b[0m+ 390: \x1b[38;2;192;197;206mSTAGE_STATUS: COMPLETE\x1b[K\n"
            "STAGE_STATUS: COMPLETE\x1b[0m\x1b[0m\n"
        )
        self.assertEqual(receipt.find_receipts(raw, "stage"), ["STAGE_STATUS: COMPLETE"])

    def test_duplicate_exact_receipts_remain_detectable(self):
        raw = "TASK_STATUS: COMPLETE\n\x1b[32mTASK_STATUS: COMPLETE\x1b[0m\n"
        self.assertEqual(
            receipt.find_receipts(raw, "task"),
            ["TASK_STATUS: COMPLETE", "TASK_STATUS: COMPLETE"],
        )

    def test_markdown_wrapped_cursor_task_receipt_is_accepted(self):
        self.assertEqual(
            receipt.find_receipts("`TASK_STATUS: COMPLETE`\n", "task"),
            ["TASK_STATUS: COMPLETE"],
        )

    def test_blocked_receipt_is_parsed_but_not_completion(self):
        self.assertEqual(
            receipt.find_receipts("STAGE_STATUS: BLOCKED ARCHITECTURE_DECISION_REQUIRED\n", "stage"),
            ["STAGE_STATUS: BLOCKED ARCHITECTURE_DECISION_REQUIRED"],
        )

    def test_stage_document_is_valid_fallback_when_cli_omits_receipt(self):
        with tempfile.TemporaryDirectory() as raw:
            directory = Path(raw)
            log = directory / "kiro.log"
            document = directory / "requirements.md"
            log.write_text("Revision complete. Receipt retained in document.\n")
            document.write_text("# Requirements\n\nSTAGE_STATUS: COMPLETE\n")
            completed = subprocess.run(
                [
                    sys.executable,
                    str(ROOT / "scripts/receipt-check.py"),
                    "--log",
                    str(log),
                    "--document",
                    str(document),
                    "--kind",
                    "stage",
                ],
                check=True,
                capture_output=True,
                text=True,
            )
            self.assertIn("stage document fallback", completed.stdout)

    def test_blocked_cli_receipt_cannot_be_overridden_by_document(self):
        with tempfile.TemporaryDirectory() as raw:
            directory = Path(raw)
            log = directory / "kiro.log"
            document = directory / "requirements.md"
            log.write_text("STAGE_STATUS: BLOCKED ARCHITECTURE_DECISION_REQUIRED\n")
            document.write_text("STAGE_STATUS: COMPLETE\n")
            completed = subprocess.run(
                [
                    sys.executable,
                    str(ROOT / "scripts/receipt-check.py"),
                    "--log",
                    str(log),
                    "--document",
                    str(document),
                    "--kind",
                    "stage",
                ],
                check=False,
                capture_output=True,
                text=True,
            )
            self.assertNotEqual(completed.returncode, 0)


class CloseoutMetadataTests(unittest.TestCase):
    def test_feature_index_and_traceability_updates_are_bounded(self):
        with tempfile.TemporaryDirectory() as raw:
            directory = Path(raw)
            index = directory / "FEATURE_INDEX.md"
            index.write_text(
                "| Feature | Name | Phase | Status | Path | Notes |\n"
                "|---|---|---|---|---|---|\n"
                "| FEATURE-0015 | Pools | Phase 2 | Planned | file.md | Existing. |\n"
            )
            closeout.update_feature_index(index, "FEATURE-0015", "PR #17 merged")
            updated = index.read_text()
            self.assertIn("Implemented and Merged", updated)
            self.assertIn("PR #17 merged", updated)

            traceability = directory / "FEATURE_TRACEABILITY_MATRIX.md"
            traceability.write_text(
                "| Feature | Name | Status | A | B | C | Gate | D |\n"
                "|---|---|---|---|---|---|---|---|\n"
                "| FEATURE-0015 | Pools | Planned | a | b | c | Pending | d |\n"
            )
            closeout.update_traceability(traceability, "FEATURE-0015", "Final gate; PR #17")
            updated = traceability.read_text()
            self.assertIn("Implemented and Merged", updated)
            self.assertIn("Final gate; PR #17", updated)

    def test_frontmatter_update_preserves_document_body(self):
        with tempfile.TemporaryDirectory() as raw:
            path = Path(raw) / "feature.md"
            path.write_text("---\nstatus: planned\ntitle: Pools\n---\n# Body\nKeep me.\n")
            closeout.update_frontmatter(path, {"status": "implemented_and_merged", "merged_pr": "#17"})
            updated = path.read_text()
            self.assertIn("status: implemented_and_merged", updated)
            self.assertIn("merged_pr: #17", updated)
            self.assertTrue(updated.endswith("# Body\nKeep me.\n"))


if __name__ == "__main__":
    unittest.main()
