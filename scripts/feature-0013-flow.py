#!/usr/bin/env python3
"""Safe, resumable FEATURE-0013 task orchestrator.

Runs approved T-xxx tasks in order through Cursor Agent, enforces the
FEATURE-0013 architecture-boundary checker before and after task execution,
validates changed paths, runs repository verification, commits one task at a
time, and stops immediately on drift.

This script does not use the OpenAI reviewer API and does not write an approval
token. Use --human-approved-for-cursor to record that a human has authorized
starting implementation from the already-approved tasks.md.
"""

from __future__ import annotations

import argparse
import datetime as dt
import json
import os
import re
import shlex
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, Sequence

FEATURE = "FEATURE-0013"
EXPECTED_BRANCH = "feature-0013-decision-record-and-auditevent-standard"
EXPECTED_SLUG = "decision-object-and-auditevent-standard"
TASKS_REL = Path(".kiro/specs") / EXPECTED_SLUG / "tasks.md"
REQ_REL = Path(".kiro/specs") / EXPECTED_SLUG / "requirements.md"
DESIGN_REL = Path(".kiro/specs") / EXPECTED_SLUG / "design.md"
ARCH_REL = Path("docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md")
STATE_REL = Path(".automation/state") / f"{FEATURE}.json"
DEFAULT_IMAGE = "golang:1.22"
TASK_RE = re.compile(r"^### (?P<id>T-\d{3})\s+[—-]\s+(?P<title>.+?)\s*$", re.MULTILINE)

GENERATED_PREFIXES = (
    "docs/generated-prompts/FEATURE-0013/",
    ".automation/generated-prompts/FEATURE-0013/",
    ".automation/reviews/FEATURE-0013/",
    ".automation/logs/FEATURE-0013/",
    ".automation/model-usage/FEATURE-0013/",
    ".automation/reports/FEATURE-0013/",
    "scripts/__pycache__/",
)

PROTECTED_EXACT = {
    str(ARCH_REL),
    str(REQ_REL),
    str(DESIGN_REL),
    str(TASKS_REL),
    "docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-consolidated-architecture.md",
    "go.mod",
    "go.sum",
}

PROTECTED_PREFIXES = (
    ".kiro/specs/decision-object-and-auditevent-standard/design.superseded-patched-draft.md",
    ".automation/reviews/",
    ".automation/generated-prompts/",
    ".automation/reports/",
    "docs/generated-prompts/",
)

ALLOWED_PREFIXES = (
    "internal/decision/",
    "internal/apiconform/",
    "api/schemas/",
    "tests/conformance/",
    "docs/traceability/",
    "docs/reviews/feature-gates/FEATURE-0013",
    "scripts/",
)

ALLOWED_EXACT = {
    str(STATE_REL),
}

FORBIDDEN_TEXT_PATTERNS = (
    (re.compile(r"\bAuditScope\b"), "AuditScope must not be introduced"),
    (re.compile(r"\b(PendingDecision|DeferredDecision|DecisionPending)\b"), "pending decision envelope must not be introduced"),
    (re.compile(r"api/schemas/SCHEMA_BASELINE_MANIFEST\.json"), "parallel schema baseline manifest must not be introduced"),
    (re.compile(r"api/schemas/diffs/"), "parallel schema diffs directory must not be introduced"),
    (re.compile(r"api/schemas/approvals/"), "parallel schema approvals directory must not be introduced"),
)

CURSOR_ALLOWED_PATHS = """internal/decision/**
internal/apiconform/**
api/schemas/**
tests/conformance/**
docs/traceability/**
docs/reviews/feature-gates/FEATURE-0013*.md
scripts/stage-matrix-e-0013.sh
.automation/state/FEATURE-0013.json
""".strip()


class FlowError(RuntimeError):
    pass


@dataclass(frozen=True)
class TaskEntry:
    task_id: str
    title: str
    start: int
    end: int


@dataclass(frozen=True)
class Plan:
    tasks: list[TaskEntry]

    @property
    def ids(self) -> list[str]:
        return [task.task_id for task in self.tasks]

    def task(self, task_id: str) -> TaskEntry:
        for entry in self.tasks:
            if entry.task_id == task_id:
                return entry
        raise FlowError(f"unknown task id: {task_id}")


@dataclass
class Runner:
    root: Path
    image: str
    verbose: bool = True

    def run(
        self,
        args: Sequence[str],
        *,
        env: dict[str, str] | None = None,
        capture: bool = False,
        check: bool = True,
    ) -> subprocess.CompletedProcess[str]:
        if self.verbose:
            print("+", shlex.join([str(arg) for arg in args]), flush=True)
        merged = os.environ.copy()
        if env:
            merged.update(env)
        result = subprocess.run(
            [str(arg) for arg in args],
            cwd=self.root,
            env=merged,
            text=True,
            capture_output=capture,
            check=False,
        )
        if check and result.returncode != 0:
            if capture:
                if result.stdout:
                    print(result.stdout, end="", file=sys.stderr)
                if result.stderr:
                    print(result.stderr, end="", file=sys.stderr)
            raise FlowError(f"command failed ({result.returncode}): {shlex.join(args)}")
        return result

    def output(self, args: Sequence[str]) -> str:
        return self.run(args, capture=True).stdout.strip()

    def docker_shell(self, script: str) -> None:
        self.run([
            "docker", "run", "--rm", "-v", f"{self.root}:/src", "-w", "/src",
            self.image, "sh", "-ceu", script,
        ])


def utc_now() -> str:
    return dt.datetime.now(dt.timezone.utc).isoformat()


def repo_root() -> Path:
    result = subprocess.run(["git", "rev-parse", "--show-toplevel"], text=True, capture_output=True, check=False)
    if result.returncode != 0:
        raise SystemExit("ERROR: not inside a git repository")
    return Path(result.stdout.strip())


def read_json(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def write_json(path: Path, payload: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def load_plan(root: Path) -> Plan:
    tasks_path = root / TASKS_REL
    text = tasks_path.read_text(encoding="utf-8")
    matches = list(TASK_RE.finditer(text))
    if not matches:
        raise FlowError(f"no T-xxx tasks found in {TASKS_REL}")
    entries: list[TaskEntry] = []
    for index, match in enumerate(matches):
        end = matches[index + 1].start() if index + 1 < len(matches) else len(text)
        entries.append(TaskEntry(match.group("id"), match.group("title").strip(), match.start(), end))
    expected = [f"T-{index:03d}" for index in range(1, len(entries) + 1)]
    actual = [entry.task_id for entry in entries]
    if actual != expected:
        raise FlowError(f"task ids are not contiguous T-001..T-{len(entries):03d}: {actual}")
    return Plan(entries)


def status_paths(root: Path) -> set[str]:
    result = subprocess.run(["git", "status", "--porcelain"], cwd=root, text=True, capture_output=True, check=True)
    paths: set[str] = set()
    for raw in result.stdout.splitlines():
        if not raw:
            continue
        path = raw[3:]
        if " -> " in path:
            path = path.split(" -> ", 1)[1]
        paths.add(path)
    return paths


def is_generated(path: str) -> bool:
    return any(path == prefix.rstrip("/") or path.startswith(prefix) for prefix in GENERATED_PREFIXES)


def non_generated_changes(root: Path) -> set[str]:
    return {path for path in status_paths(root) if not is_generated(path)}


def require_clean_for_task(runner: Runner) -> None:
    paths = non_generated_changes(runner.root)
    if paths:
        raise FlowError("working tree has non-generated changes:\n  " + "\n  ".join(sorted(paths)))


def is_allowed_change(path: str) -> bool:
    if path in ALLOWED_EXACT:
        return True
    return any(path.startswith(prefix) for prefix in ALLOWED_PREFIXES)


def validate_paths(paths: Iterable[str]) -> None:
    material = sorted(path for path in paths if not is_generated(path))
    if not material:
        raise FlowError("task produced no material changed paths")
    bad: list[str] = []
    for path in material:
        if path in PROTECTED_EXACT:
            bad.append(f"protected exact path: {path}")
        elif any(path.startswith(prefix) for prefix in PROTECTED_PREFIXES):
            bad.append(f"protected generated/spec path: {path}")
        elif not is_allowed_change(path):
            bad.append(f"not in FEATURE-0013 allowed implementation/evidence scope: {path}")
    if bad:
        raise FlowError("changed path boundary violation:\n  " + "\n  ".join(bad))


def negated_context(text: str, start: int, end: int) -> bool:
    line_start = text.rfind("\n", 0, start) + 1
    line_end = text.find("\n", end)
    if line_end == -1:
        line_end = len(text)
    line = text[line_start:line_end].lower()
    markers = (
        "no ",
        "not ",
        "never ",
        "without ",
        "do not ",
        "must not ",
        "does not ",
        "is not ",
        "are not ",
        "is never ",
        "must never ",
        "prohibit",
        "forbid",
        "reject",
        "out-of-scope",
        "non-task",
    )
    return any(marker in line for marker in markers)


def validate_text_drift(root: Path, paths: Iterable[str]) -> None:
    offenders: list[str] = []
    for rel in sorted(paths):
        if is_generated(rel):
            continue
        path = root / rel
        if not path.exists() or path.is_dir():
            continue
        try:
            text = path.read_text(encoding="utf-8")
        except UnicodeDecodeError:
            continue
        for pattern, message in FORBIDDEN_TEXT_PATTERNS:
            for match in pattern.finditer(text):
                if negated_context(text, match.start(), match.end()):
                    continue
                line = text.count("\n", 0, match.start()) + 1
                offenders.append(f"{rel}:{line}: {message}: {match.group(0)}")
    if offenders:
        raise FlowError("forbidden FEATURE-0013 drift text introduced:\n  " + "\n  ".join(offenders[:50]))


def update_state(root: Path, **values: object) -> None:
    state_path = root / STATE_REL
    state = read_json(state_path)
    state.update(values)
    state["updated_at"] = utc_now()
    write_json(state_path, state)


def ensure_prerequisites(runner: Runner, *, human_approved: bool) -> None:
    state = read_json(runner.root / STATE_REL)
    if state.get("feature_id") != FEATURE:
        raise FlowError("feature state is not FEATURE-0013")
    if state.get("slug") != EXPECTED_SLUG:
        raise FlowError("feature state slug is not FEATURE-0013 slug")
    branch = runner.output(["git", "branch", "--show-current"])
    if branch != EXPECTED_BRANCH:
        raise FlowError(f"wrong branch: expected {EXPECTED_BRANCH}, found {branch}")
    token = state.get("tasks_approval_token")
    if token != "APPROVED_FOR_CURSOR" and not human_approved:
        raise FlowError(
            "missing APPROVED_FOR_CURSOR. Because OpenAI reviewer quota is unavailable, "
            "rerun with --human-approved-for-cursor only after you have explicitly approved tasks.md."
        )
    runner.run(["git", "diff", "--check"])
    runner.run(["python3", "scripts/feature-0013-architecture-boundary-check.py", "--feature", FEATURE, "--mode", "readiness"])
    for stage in ("requirements", "design", "tasks"):
        runner.run(["python3", "scripts/feature-0013-architecture-boundary-check.py", "--feature", FEATURE, "--stage", stage, "--mode", "post"])
    runner.run(["docker", "info"], capture=True)
    if not (runner.output(["sh", "-c", "command -v cursor-agent || command -v agent || true"])):
        raise FlowError("Cursor Agent CLI is not available: install cursor-agent or set CURSOR_AGENT_BIN")


def select_tasks(plan: Plan, state: dict, start_task: str | None, stop_after: str | None, max_tasks: int) -> list[TaskEntry]:
    ids = plan.ids
    for label, value in (("start-task", start_task), ("stop-after", stop_after)):
        if value and value not in ids:
            raise FlowError(f"unknown {label}: {value}")
    if start_task:
        start_index = ids.index(start_task)
    else:
        last = state.get("last_committed_task") or state.get("automation_last_committed_task")
        start_index = ids.index(last) + 1 if last in ids else 0
    selected: list[TaskEntry] = []
    for task in plan.tasks[start_index:]:
        selected.append(task)
        if stop_after == task.task_id:
            break
        if max_tasks and len(selected) >= max_tasks:
            break
    return selected


def run_verification(runner: Runner, *, final: bool = False) -> None:
    runner.run(["git", "diff", "--check"])
    runner.run(["python3", "scripts/feature-0013-architecture-boundary-check.py", "--feature", FEATURE, "--stage", "tasks", "--mode", "post"])
    runner.run(["./scripts/guardrails.sh", "--feature", FEATURE])
    runner.run(["./scripts/verify.sh"])
    if final:
        runner.docker_shell('test -z "$(gofmt -l .)"\ngo test -race ./...\ngo vet ./...')


def commit_task(runner: Runner, task: TaskEntry, *, push: bool) -> str:
    paths = non_generated_changes(runner.root)
    validate_paths(paths)
    validate_text_drift(runner.root, paths)
    runner.run(["git", "add", "--", *sorted(paths)])
    runner.run(["git", "diff", "--cached", "--check"])
    message = f"feat(feature-0013): complete {task.task_id} {task.title}"[:180]
    runner.run(["git", "commit", "-m", message])
    commit = runner.output(["git", "rev-parse", "HEAD"])
    if push:
        runner.run(["git", "push"])
    return commit


def execute_task(runner: Runner, task: TaskEntry, *, push: bool) -> None:
    require_clean_for_task(runner)
    print("\n" + "=" * 72)
    print(f"==> FEATURE-0013 {task.task_id}: {task.title}")
    print("=" * 72)
    update_state(
        runner.root,
        current_stage="cursor",
        current_task=task.task_id,
        automation_flow_status="running",
        automation_current_task=task.task_id,
        automation_started_at=utc_now(),
    )
    env = {
        "FEATURE_FACTORY_CURSOR_MODE": "auto",
        "FEATURE_FACTORY_CURSOR_VERIFY": "0",
        "FEATURE_FACTORY_ALLOWED_PATHS": CURSOR_ALLOWED_PATHS,
        "GO_DOCKER_IMAGE": runner.image,
    }
    try:
        runner.run(["make", "ff-cursor-task-auto", f"FEATURE={FEATURE}", f"TASK={task.task_id}"], env=env)
        paths = non_generated_changes(runner.root)
        validate_paths(paths)
        validate_text_drift(runner.root, paths)
        run_verification(runner)
        update_state(
            runner.root,
            last_committed_task=task.task_id,
            automation_last_committed_task=task.task_id,
            automation_last_machine_gate=task.task_id,
            automation_last_machine_gate_result="PASS",
            status=f"cursor_task_{task.task_id}_verified",
            automation_flow_status="verified",
        )
        commit = commit_task(runner, task, push=push)
        print(f"==> committed {task.task_id}: {commit[:12]}")
    except Exception as exc:
        update_state(
            runner.root,
            status=f"automation_task_{task.task_id}_failed",
            automation_flow_status="failed",
            automation_failed_task=task.task_id,
            automation_failure=str(exc),
        )
        print(f"ERROR: FEATURE-0013 automation stopped at {task.task_id}: {exc}", file=sys.stderr)
        raise


def print_plan(plan: Plan, selected: list[TaskEntry]) -> None:
    selected_ids = {task.task_id for task in selected}
    print(f"FEATURE-0013 task plan: {len(plan.tasks)} tasks")
    for task in plan.tasks:
        marker = "run" if task.task_id in selected_ids else "skip"
        print(f"  {task.task_id}  {marker:4}  {task.title}")


def final_summary(runner: Runner, plan: Plan) -> None:
    update_state(
        runner.root,
        status="automation_tasks_completed_pending_human_review",
        human_gate_required="true",
        automation_flow_status="implementation_tasks_completed",
        automation_tasks_completed_at=utc_now(),
        current_task="",
    )
    run_verification(runner, final=True)
    paths = non_generated_changes(runner.root)
    if paths:
        runner.run(["git", "add", "--", *sorted(paths)])
        runner.run(["git", "diff", "--cached", "--check"])
        runner.run(["git", "commit", "-m", "chore(feature-0013): record implementation automation completion"])
    summary = {
        "feature": FEATURE,
        "completed_at": utc_now(),
        "head": runner.output(["git", "rev-parse", "HEAD"]),
        "tasks": plan.ids,
        "human_gate_required": True,
        "note": "Final Codex/human code review is required; no final approval token was written.",
    }
    out = runner.root / ".automation/logs" / FEATURE / "implementation-summary.json"
    write_json(out, summary)
    print("\n==> FEATURE-0013 task automation completed")
    print("    Final Codex/human code review is required before merge approval.")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run FEATURE-0013 tasks sequentially with guardrails")
    parser.add_argument("--start-task")
    parser.add_argument("--stop-after")
    parser.add_argument("--max-tasks", type=int, default=0, help="0 means no limit")
    parser.add_argument("--plan", action="store_true", help="print selected tasks and exit")
    parser.add_argument("--dry-run", action="store_true", help="validate prerequisites and print selected tasks without running Cursor")
    parser.add_argument("--push", action="store_true", help="push after each task commit")
    parser.add_argument("--human-approved-for-cursor", action="store_true", help="explicit human authorization to begin implementation without reviewer API token")
    parser.add_argument("--image", default=os.environ.get("GO_DOCKER_IMAGE", DEFAULT_IMAGE))
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    root = repo_root()
    runner = Runner(root=root, image=args.image)
    plan = load_plan(root)
    state = read_json(root / STATE_REL)
    selected = select_tasks(plan, state, args.start_task, args.stop_after, args.max_tasks)
    if args.plan or args.dry_run:
        print_plan(plan, selected)
        if args.plan:
            return
    ensure_prerequisites(runner, human_approved=args.human_approved_for_cursor)
    if args.dry_run:
        print("PASS: dry-run prerequisites and boundary checks passed")
        return
    if not selected:
        print("==> No tasks selected; nothing to run")
        return
    for task in selected:
        execute_task(runner, task, push=args.push)
    if selected[-1].task_id == plan.ids[-1]:
        final_summary(runner, plan)


if __name__ == "__main__":
    try:
        main()
    except FlowError as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1)
