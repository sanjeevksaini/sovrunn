import importlib.util
import json
import tempfile
import unittest
import sys
from copy import deepcopy
from pathlib import Path

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
