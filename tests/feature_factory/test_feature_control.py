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
review_prompt = load_module("generic_review_prompt", ROOT / "scripts/generic-review-prompt.py")
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

    def test_design_review_may_exclude_only_a_declared_handoff(self):
        excluded = self.data["feature"]["handoffs"][0]
        configured = deepcopy(self.data)
        configured["context"]["review_exclusions"] = {"design": [excluded]}
        control.validate(configured, expected_feature="FEATURE-0014")
        paths = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(configured, "review", review_stage="design")
        }
        self.assertNotIn(excluded, paths)

        invalid = deepcopy(configured)
        invalid["context"]["review_exclusions"]["design"] = ["docs/glossary.md"]
        with self.assertRaisesRegex(control.ControlError, "must name a feature handoff"):
            control.validate(invalid, expected_feature="FEATURE-0014")

    def test_review_stage_context_is_scoped_to_the_declared_review_stage(self):
        configured = deepcopy(self.data)
        context_path = "docs/engineering/go-observability-standard.md"
        configured["context"]["review_stage_context"] = {"design": [context_path]}
        control.validate(configured, expected_feature="FEATURE-0014")
        design_paths = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(configured, "review", review_stage="design")
        }
        task_paths = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(configured, "review", review_stage="tasks")
        }
        self.assertIn(context_path, design_paths)
        self.assertNotIn(context_path, task_paths)

    def test_review_authority_block_validates_and_enforces_finding_rules(self):
        feature_control = json.loads(
            (ROOT / ".automation/features/FEATURE-0018.control.json").read_text()
        )
        control.validate(feature_control, expected_feature="FEATURE-0018")
        authority = feature_control["context"]["review_authority"]["design"]
        self.assertEqual(authority["handoff"], "ADH-2026-076")
        self.assertEqual(
            authority["finding_classifications"],
            [
                "STALE_TRANSCRIPTION",
                "DESIGN_EXECUTABILITY",
                "REQUIREMENT_GAP",
                "ARCHITECTURE_CLARIFICATION_REQUIRED",
                "OUT_OF_SCOPE",
            ],
        )
        manifest = control.context_manifest(
            feature_control, "review", review_stage="design"
        )
        prompt = control.prompt_fragment(
            feature_control, manifest, "review", review_stage="design"
        )
        self.assertIn("Controlled review authority (ADH-2026-076)", prompt)
        self.assertIn(
            "You may raise `DESIGN_EXECUTABILITY` only when", prompt
        )
        self.assertIn("design-mechanics-contract.json", prompt)
        self.assertIn(
            ".kiro/specs/governance-iam-approval-exception-foundation/design.md",
            prompt,
        )

    def test_review_authority_rejects_unknown_stage(self):
        invalid = deepcopy(self.data)
        invalid.setdefault("context", {})
        invalid["context"]["review_authority"] = {
            "implementation": {
                "handoff": "ADH-2026-076",
                "mechanics_authority": "x",
                "semantic_authority": "y",
                "finding_classifications": ["DESIGN_EXECUTABILITY"],
                "design_executability_allowed_only_when": ["a"],
                "reviewer_must_not": ["b"],
                "finding_requirement": "c",
                "hash_bound_evidence": ["docs/glossary.md"],
            }
        }
        with self.assertRaisesRegex(
            control.ControlError, "review_authority has unknown review stages"
        ):
            control.validate(invalid, expected_feature="FEATURE-0014")

    def test_review_authority_requires_complete_finding_rule_keys(self):
        invalid = deepcopy(self.data)
        invalid["context"]["review_authority"] = {
            "design": {
                "handoff": "ADH-2026-076",
                "mechanics_authority": "x",
                "semantic_authority": "y",
                "finding_classifications": ["DESIGN_EXECUTABILITY"],
                "design_executability_allowed_only_when": ["a"],
                "reviewer_must_not": ["b"],
                "finding_requirement": "c"
            }
        }
        with self.assertRaisesRegex(control.ControlError, "hash_bound_evidence"):
            control.validate(invalid, expected_feature="FEATURE-0014")

    def test_review_authority_handoff_must_be_adh_identifier(self):
        invalid = deepcopy(self.data)
        invalid["context"]["review_authority"] = {
            "design": {
                "handoff": "not-an-adh",
                "mechanics_authority": "x",
                "semantic_authority": "y",
                "finding_classifications": ["DESIGN_EXECUTABILITY"],
                "design_executability_allowed_only_when": ["a"],
                "reviewer_must_not": ["b"],
                "finding_requirement": "c",
                "hash_bound_evidence": ["docs/glossary.md"],
            }
        }
        with self.assertRaisesRegex(control.ControlError, "must be an ADH-YYYY-NNN"):
            control.validate(invalid, expected_feature="FEATURE-0014")

    def test_review_prompt_without_review_authority_omits_the_block(self):
        # The closed FEATURE-0014 fixture declares no review_authority; the design
        # review prompt must not synthesize the ADH-076 finding-rule block.
        manifest = control.context_manifest(self.data, "review", review_stage="design")
        prompt = control.prompt_fragment(
            self.data, manifest, "review", review_stage="design"
        )
        self.assertNotIn("Controlled review authority", prompt)

    def test_feature_0018_tasks_review_uses_governed_requirements_projection(self):
        feature_control = json.loads(
            (ROOT / ".automation/features/FEATURE-0018.control.json").read_text()
        )
        control.validate(feature_control, expected_feature="FEATURE-0018")
        paths = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(
                feature_control, "review", review_stage="tasks"
            )
        }
        spec = feature_control["feature"]["spec_path"]
        self.assertNotIn(f"{spec}/requirements.md", paths)
        self.assertIn(
            ".automation/context-projections/FEATURE-0018.tasks-review-requirements.md",
            paths,
        )
        self.assertIn(f"{spec}/design.md", paths)

    def test_task_context_may_exclude_only_a_declared_handoff(self):
        excluded = self.data["feature"]["handoffs"][0]
        configured = deepcopy(self.data)
        configured["context"]["stage_exclusions"] = {"tasks": [excluded]}
        control.validate(configured, expected_feature="FEATURE-0014")
        paths = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(configured, "tasks")
        }
        self.assertNotIn(excluded, paths)

        invalid = deepcopy(configured)
        invalid["context"]["stage_exclusions"]["tasks"] = ["docs/glossary.md"]
        with self.assertRaisesRegex(control.ControlError, "must name a feature handoff"):
            control.validate(invalid, expected_feature="FEATURE-0014")

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

    def test_feature_task_heading_is_parsed_by_cursor_and_orchestrator(self):
        plan = self.task_plan_fixture.replace("### Task 1:", "### TASK-F16-01 —").replace(
            "### Task 18:", "### TASK-F16-18 —"
        )
        self.assertEqual(sorted(orchestrator.blocks(plan)), [1, 18])
        self.assertEqual(sorted(cursor_prompt.task_blocks(plan)), [1, 18])

    def test_grouped_path_bullets_expand_to_every_exact_path(self):
        block = """### Task 13: Grouped paths

**Writable paths:**
- `internal/api/one.go`, `internal/api/two.go`

**Tests:**
- `internal/api/one_test.go`, `internal/api/two_test.go`: the inherited `internal/apivalid` pipeline returns `VALIDATION_FAILED` as expected.

**Commit message:**
```
feat(example): grouped paths
```
"""
        expected = [
            "internal/api/one.go",
            "internal/api/two.go",
            "internal/api/one_test.go",
            "internal/api/two_test.go",
        ]
        self.assertEqual(cursor_prompt.writable_paths(block), expected)
        self.assertEqual(orchestrator.task_writable_paths(block), expected)

    def test_explicit_directory_rules_reach_cursor_and_orchestrator_allowlists(self):
        block = """### Task 14: Directory-scoped foundation

**Writable paths:**
- `internal/govaccess/` (root composition package)
- `internal/govaccess/clock/` (deterministic clock implementation)
- `internal/govaccess/model/principalref.go`

**Tests:**
- `internal/govaccess/clock/clock_test.go`

**Commit message:**
```
feat(govaccess): add deterministic foundation
```
"""
        expected = [
            "internal/govaccess/",
            "internal/govaccess/clock/",
            "internal/govaccess/model/principalref.go",
            "internal/govaccess/clock/clock_test.go",
        ]
        self.assertEqual(cursor_prompt.writable_paths(block), expected)
        self.assertEqual(orchestrator.task_writable_paths(block), expected)

    def test_tests_narrative_may_reuse_a_writable_test_path(self):
        block = """### Task 17: Conformance

**Writable paths:**
- `internal/apiconform/feature_0015_test.go`

**Tests:**
- Every local conformance scenario is asserted by the writable test file.

**Commit message:**
```
test(conformance): add local cases
```
"""
        self.assertEqual(
            cursor_prompt.writable_paths(block),
            ["internal/apiconform/feature_0015_test.go"],
        )

    def test_verification_checkpoint_is_not_a_cursor_or_commit_task(self):
        source = (ROOT / "scripts/feature-orchestrator.py").read_text()
        checkpoint = source.index("def run_verification_checkpoint")
        checkpoint_source = source[checkpoint : source.index("\ndef main()", checkpoint)]
        self.assertIn('"make", "ff-feature-gate"', checkpoint_source)
        self.assertNotIn("generic-feature-boundary-check.py", checkpoint_source)
        self.assertIn("if completed_final_commit:", source)
        self.assertIn("running verification-only Task", source)

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

    def test_task_review_delta_receipt_binds_raw_audit_without_full_lines(self):
        raw = b'{"full":"semantic delta"}\n'
        receipt_data = review_prompt.compact_delta_receipt(
            {
                "classification": "SEMANTIC_REVIEW_REQUIRED",
                "semantic_review_required": True,
                "summary": "Review the current plan.",
                "baseline": "baseline.md",
                "current": "tasks.md",
                "baseline_sha256": "baseline-sha",
                "current_sha256": "current-sha",
                "changed_lines": {"added": 3, "removed": 2},
                "identifiers": {"added": ["TASK-F18-01"], "removed": []},
                "headings": {"added": [], "removed": ["Old task"]},
                "normative_lines": {"added": ["MUST do X"], "removed": ["ONLY do Y"], "truncated": False},
            },
            raw,
        )
        self.assertEqual(receipt_data["raw_delta_sha256"], __import__("hashlib").sha256(raw).hexdigest())
        self.assertEqual(receipt_data["normative_line_delta"], {"listed_added": 1, "listed_removed": 1, "truncated": False})
        rendered = review_prompt.render_delta_receipt(receipt_data)
        self.assertIn("Hash-bound semantic-delta receipt", rendered)
        self.assertNotIn("MUST do X", rendered)

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

    def test_closeout_updates_phase_2r_context(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "CURRENT_PHASE_CONTEXT.md"
            path.write_text(
                "# Current Phase Context\n\n"
                "## Phase 2R Goal\n\nGoal.\n\n"
                "## Phase 2R Completed Features\n\n"
                "| Feature | Status | Approval |\n|---|---|---|\n\n"
                "## Phase 2R Next Planned Feature\n\n"
                "| Feature | Status | Controlling decisions |\n|---|---|---|\n"
            )
            closeout.update_phase_context(
                path,
                {"feature": {"id": "FEATURE-0018", "title": "Governance Foundation"}},
                "PR #20 merged 2026-09-07",
                "2026-09-07",
            )
            result = path.read_text()
            self.assertIn("## Phase 2R Goal", result)
            self.assertIn("FEATURE-0018 status: implemented and merged", result)
            self.assertIn("| FEATURE-0018 Governance Foundation | Implemented and merged", result)

    def test_closeout_updates_phase_2r_baseline(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "CURRENT_ARCHITECTURE_BASELINE.md"
            path.write_text("# Baseline\n\n## Phase 2R Next Planned Feature\n\nFEATURE-0015\n")
            closeout.update_architecture_baseline(
                path,
                {
                    "feature": {
                        "id": "FEATURE-0018",
                        "title": "Governance Foundation",
                        "architecture": "docs/architecture/feature-0018.md",
                        "handoffs": [],
                    },
                    "ownership": {"owned_resources": [], "excluded_features": []},
                },
                "PR #20 merged 2026-09-07",
            )
            result = path.read_text()
            self.assertIn("## Implemented Phase 2R Features", result)
            self.assertIn("Completed and merged: `FEATURE-0018: Governance Foundation`", result)

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
| REQ-F99-01 | Preserve the approved behavior | DEC-0099 | VS0-CF-F99 |

## Acceptance Criteria
| ID | Criterion | Conformance ID |
|----|-----------|----------------|
| AC-F99-01 | Prove the registered scenario | VS0-CF-F99 |
"""
    conformance = {
        "VS0-CF-F99": {
            "id": "VS0-CF-F99",
            "owner": "FEATURE-0099",
            "inputs": "approved-feature-behavior",
            "expectedState": "preserved",
            "expectedError": "none",
            "expectedSideEffects": "no side effects",
            "gate": "decision",
        },
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
| REQ-F99-01 | Preserve the approved behavior | DEC-0099 | VS0-CF-F99 |

## Canonical acceptance ledger
| ID | Criterion | Conformance ID |
|----|-----------|----------------|
| AC-F99-01 | Prove the registered scenario | VS0-CF-F99 |

## Exact conformance semantics ledger
| ID | Owner | Inputs | Expected State | Expected Error | Expected Side Effects | Gate |
|----|-------|--------|----------------|----------------|-----------------------|------|
| VS0-CF-F99 | FEATURE-0099 | approved-feature-behavior | preserved | none | no side effects | decision |

### REQ-F99-01 — Approved behavior
The behavior remains unchanged. AC-F99-01 uses VS0-CF-F99 only for approved-feature-behavior proof.
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

    def test_cross_feature_conformance_cannot_appear_in_exact_ledger(self):
        target = self.valid_requirements().replace(
            "| VS0-CF-F99 | FEATURE-0099 | approved-feature-behavior | preserved | none | no side effects | decision |",
            "| VS0-CF-F99 | FEATURE-0099 | approved-feature-behavior | preserved | none | no side effects | decision |\n"
            "| VS0-CF-F09 | FEATURE-0023 | unavailable-required-participation | denied-placement | VALIDATION_FAILED | no-silent-substitution | decision |",
        )
        errors = []
        kiro_semantic.check_requirements(
            errors, "FEATURE-0099", self.feature_text, target, self.conformance
        )
        self.assertTrue(any("not authorized" in error for error in errors))


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

    def test_cursor_transport_suffix_after_wrapped_receipt_is_accepted(self):
        raw = "`TASK_STATUS: COMPLETE`Verification finished successfully: make test passed.\n"
        self.assertEqual(
            receipt.find_receipts(raw, "task"), ["TASK_STATUS: COMPLETE"]
        )

    def test_feature_task_verification_suffix_after_wrapped_receipt_is_accepted(self):
        raw = (
            "`TASK_STATUS: COMPLETE`TASK-F16-12 verification finished successfully: "
            "conformance passed.\n"
        )
        self.assertEqual(
            receipt.find_receipts(raw, "task"), ["TASK_STATUS: COMPLETE"]
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


class Feature0018DesignReviewTargetTests(unittest.TestCase):
    """ADH-2026-076 clause 3: hash-bound semantic design-review target projection
    is rendered as the reviewer target while semantic-delta.py hashes raw design.md."""

    builder = load_module(
        "feature0018_design_review_target",
        ROOT / "scripts/feature0018-design-review-target.py",
    )
    projection = ROOT / ".automation/context-projections/FEATURE-0018.design-review-target.md"
    manifest = ROOT / ".automation/context-projections/FEATURE-0018.design-review-target.manifest.json"
    design = ROOT / ".kiro/specs/governance-iam-approval-exception-foundation/design.md"
    contract = ROOT / ".kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json"

    def test_builder_compiles(self):
        path = ROOT / "scripts/feature0018-design-review-target.py"
        compile(path.read_text(), str(path), "exec")

    def test_builder_is_inert_under_active_checkpoint(self):
        # ADH-2026-077 §4/§5: while the FEATURE-0018 implementation checkpoint is
        # active, the design-review target projection is not a design-review
        # authority and is not built. The CLI must SKIP (rc=0) rather than build
        # or compare against the frozen baseline design.md.
        self.assertTrue(self.builder.checkpoint_active())
        completed = subprocess.run(
            [sys.executable, str(ROOT / "scripts/feature0018-design-review-target.py"), "--check-only"],
            cwd=str(ROOT),
            capture_output=True,
            text=True,
        )
        self.assertEqual(completed.returncode, 0, completed.stderr)
        self.assertIn("SKIP:", completed.stdout)
        self.assertIn("ADH-2026-077", completed.stdout)

    def test_projection_manifest_binds_the_discarded_draft(self):
        # The retained ranges and manifest were bound (ADH-2026-076) to the later
        # working design draft, which ADH-2026-077 discards. The manifest binding
        # intentionally no longer matches the frozen baseline design.md; it is
        # historical review-context evidence only.
        manifest = json.loads(self.manifest.read_text())
        self.assertEqual(manifest["controlling_handoff"], "ADH-2026-076")
        frozen = __import__("hashlib").sha256(self.design.read_bytes()).hexdigest()
        self.assertNotEqual(
            manifest["binding"]["raw_design_sha256"],
            frozen,
            "the projection manifest must remain bound to the discarded draft, not the frozen baseline",
        )

    def test_review_prompt_renders_render_target_but_hashes_raw_target(self):
        source = (ROOT / "scripts/generic-review-prompt.py").read_text()
        # semantic-delta is always run against --target (raw design.md).
        self.assertIn('"--current",\n        str(target),', source)
        # The rendered document content and target path use the render target.
        self.assertIn('.replace("{{TARGET_PATH}}", render_target_rel)', source)
        self.assertIn('.replace("{{DOCUMENT_CONTENT}}", render_target.read_text())', source)
        # render_target defaults to target when --render-target is absent.
        self.assertIn(
            "render_target_rel = args.render_target if args.render_target else args.target",
            source,
        )

    def test_reviewer_stage_routes_feature_0018_design_through_projection(self):
        source = (ROOT / "scripts/reviewer-stage.sh").read_text()
        self.assertIn('if [[ "$FEATURE" == "FEATURE-0018" && "$STAGE" == "design" ]]; then', source)
        self.assertIn("scripts/feature0018-design-review-target.py", source)
        self.assertIn("--render-target", source)

    def test_control_declares_target_projection_as_design_review_context_and_evidence(self):
        feature_control = json.loads(
            (ROOT / ".automation/features/FEATURE-0018.control.json").read_text()
        )
        control.validate(feature_control, expected_feature="FEATURE-0018")
        design_review = {
            str(path.relative_to(ROOT))
            for path in control.resolve_context(
                feature_control, "review", review_stage="design"
            )
        }
        self.assertIn(
            ".automation/context-projections/FEATURE-0018.design-review-target.md",
            design_review,
        )
        evidence = feature_control["context"]["review_authority"]["design"]["hash_bound_evidence"]
        self.assertIn(
            ".automation/context-projections/FEATURE-0018.design-review-target.md",
            evidence,
        )
        # The manifest-controlled review context stays within the 650 KB budget.
        manifest = control.context_manifest(
            feature_control, "review", review_stage="design"
        )
        self.assertLessEqual(manifest["total_bytes"], manifest["budget_bytes"])


if __name__ == "__main__":
    unittest.main()
