#!/usr/bin/env python3
"""Safe, resumable FEATURE-0012 implementation orchestrator.

The flow executes approved leaf tasks in dependency-wave order, enforces a
clean-tree boundary around every task, performs deterministic machine gates,
commits one task at a time, runs checkpoint barriers without Cursor, and stops
with a final human review gate. It never writes a human approval token.
"""

from __future__ import annotations

import argparse
import datetime as dt
import fcntl
import fnmatch
import json
import os
import re
import shlex
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, Sequence

FEATURE = "FEATURE-0012"
EXPECTED_BRANCH = "feature-0012-api-resource-naming-status-and-validation-standard"
EXPECTED_SLUG = "api-resource-naming-status-and-validation-standard"
TASKS_REL = Path(".kiro/specs") / EXPECTED_SLUG / "tasks.md"
STATE_REL = Path(".automation/state") / f"{FEATURE}.json"
ARCH_REL = Path("docs/architecture/api-resource-standard.md")
REQUIREMENTS_REL = Path(".kiro/specs") / EXPECTED_SLUG / "requirements.md"
DESIGN_REL = Path(".kiro/specs") / EXPECTED_SLUG / "design.md"
CHECKPOINTS = {"3", "7", "13", "18"}
DEFAULT_IMAGE = "golang:1.22"
# Task 17.2 is a verification-only leaf and may legitimately commit
# only its checkbox and durable workflow state when all checks pass.
CONTROL_ONLY_TASKS = {"17.2"}

TASK_LINE_RE = re.compile(
    r"^(?P<indent>\s*)- \[(?P<mark>[ xX])\] "
    r"(?P<id>[0-9]+[a-z]?(?:\.[0-9]+[a-z]?)*)\b(?P<rest>.*)$"
)
CHECKBOX_NORMALIZE_RE = re.compile(
    r"^(\s*- \[)[ xX](\] [0-9]+[a-z]?(?:\.[0-9]+[a-z]?)?\b.*)$",
    re.MULTILINE,
)


class FlowError(RuntimeError):
    pass


@dataclass(frozen=True)
class TaskEntry:
    task_id: str
    title: str
    line_index: int
    checked: bool
    indent: int


@dataclass(frozen=True)
class Plan:
    waves: list[list[str]]
    barrier_waves: set[int]
    tasks: dict[str, TaskEntry]

    @property
    def ordered_ids(self) -> list[str]:
        return [task for wave in self.waves for task in wave]


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
            print("+", shlex.join([str(x) for x in args]), flush=True)
        merged_env = os.environ.copy()
        if env:
            merged_env.update(env)
        result = subprocess.run(
            [str(x) for x in args],
            cwd=self.root,
            env=merged_env,
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
            raise FlowError(
                f"command failed ({result.returncode}): {shlex.join(args)}"
            )
        return result

    def output(self, args: Sequence[str]) -> str:
        return self.run(args, capture=True).stdout.strip()

    def docker_shell(self, script: str) -> None:
        self.run(
            [
                "docker",
                "run",
                "--rm",
                "-v",
                f"{self.root}:/src",
                "-w",
                "/src",
                self.image,
                "sh",
                "-ceu",
                script,
            ]
        )


def utc_now() -> str:
    return dt.datetime.now(dt.timezone.utc).isoformat()


def repo_root() -> Path:
    result = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        text=True,
        capture_output=True,
        check=False,
    )
    if result.returncode != 0:
        raise FlowError("not inside a git repository")
    return Path(result.stdout.strip()).resolve()


def read_json(path: Path) -> dict:
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except FileNotFoundError as exc:
        raise FlowError(f"missing required file: {path}") from exc
    except json.JSONDecodeError as exc:
        raise FlowError(f"invalid JSON in {path}: {exc}") from exc


def write_json(path: Path, value: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(
        json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )


def update_state(root: Path, **values: object) -> dict:
    path = root / STATE_REL
    state = read_json(path)
    state.update(values)
    state["updated_at"] = utc_now()
    write_json(path, state)
    return state


def parse_task_entries(text: str) -> dict[str, TaskEntry]:
    entries: dict[str, TaskEntry] = {}
    for index, line in enumerate(text.splitlines()):
        match = TASK_LINE_RE.match(line)
        if not match:
            continue
        task_id = match.group("id")
        if task_id in entries:
            raise FlowError(f"duplicate task checkbox id: {task_id}")
        entries[task_id] = TaskEntry(
            task_id=task_id,
            title=match.group("rest").strip(),
            line_index=index,
            checked=match.group("mark").lower() == "x",
            indent=len(match.group("indent")),
        )
    return entries


def extract_dependency_graph(text: str) -> tuple[list[list[str]], set[int]]:
    marker = "## Task Dependency Graph"
    marker_index = text.find(marker)
    if marker_index < 0:
        raise FlowError("tasks.md has no Task Dependency Graph section")
    tail = text[marker_index + len(marker) :]
    match = re.search(r"```json\s*(\{.*?\})\s*```", tail, re.DOTALL)
    if not match:
        raise FlowError("tasks.md has no fenced JSON dependency graph")
    try:
        graph = json.loads(match.group(1))
    except json.JSONDecodeError as exc:
        raise FlowError(f"invalid dependency graph JSON: {exc}") from exc
    raw_waves = graph.get("waves")
    if not isinstance(raw_waves, list) or not raw_waves:
        raise FlowError("dependency graph has no waves")
    waves: list[list[str]] = []
    barriers: set[int] = set()
    seen: set[str] = set()
    for expected_id, wave in enumerate(raw_waves):
        if wave.get("id") != expected_id:
            raise FlowError(
                f"dependency wave id mismatch: expected {expected_id}, got {wave.get('id')}"
            )
        tasks = wave.get("tasks")
        if not isinstance(tasks, list) or not tasks:
            raise FlowError(f"wave {expected_id} has no tasks")
        ids = [str(task) for task in tasks]
        for task_id in ids:
            if task_id in seen:
                raise FlowError(f"task appears in multiple waves: {task_id}")
            seen.add(task_id)
        waves.append(ids)
        if bool(wave.get("barrier")):
            barriers.add(expected_id)
    return waves, barriers


def load_plan(tasks_path: Path) -> Plan:
    text = tasks_path.read_text(encoding="utf-8")
    tasks = parse_task_entries(text)
    waves, barrier_waves = extract_dependency_graph(text)
    missing = [task_id for wave in waves for task_id in wave if task_id not in tasks]
    if missing:
        raise FlowError(
            "dependency graph references missing task checkbox(es): "
            + ", ".join(missing)
        )
    graph_checkpoints = {
        task_id
        for index in barrier_waves
        for task_id in waves[index]
    }
    if graph_checkpoints != CHECKPOINTS:
        raise FlowError(
            f"checkpoint mismatch: expected {sorted(CHECKPOINTS)}, "
            f"found {sorted(graph_checkpoints)}"
        )
    return Plan(waves=waves, barrier_waves=barrier_waves, tasks=tasks)


def normalize_task_checkboxes(text: str) -> str:
    return CHECKBOX_NORMALIZE_RE.sub(r"\1 \2", text)


def checkbox_states(text: str) -> dict[str, bool]:
    return {task_id: entry.checked for task_id, entry in parse_task_entries(text).items()}


def mark_task_checked(tasks_path: Path, task_id: str) -> bool:
    lines = tasks_path.read_text(encoding="utf-8").splitlines()
    entries = parse_task_entries("\n".join(lines))
    entry = entries.get(task_id)
    if entry is None:
        raise FlowError(f"cannot mark missing task: {task_id}")
    if entry.checked:
        return False
    line = lines[entry.line_index]
    lines[entry.line_index] = line.replace("- [ ]", "- [x]", 1)
    tasks_path.write_text("\n".join(lines) + "\n", encoding="utf-8")
    return True


def task_group(task_id: str) -> str:
    return task_id.split(".", 1)[0]


def roll_up_parent(tasks_path: Path, task_id: str) -> list[str]:
    if task_id in CHECKPOINTS or "." not in task_id:
        return []
    text = tasks_path.read_text(encoding="utf-8")
    entries = parse_task_entries(text)
    parent = task_group(task_id)
    parent_entry = entries.get(parent)
    if parent_entry is None or parent_entry.checked:
        return []
    children = [
        entry
        for child_id, entry in entries.items()
        if child_id.startswith(parent + ".")
    ]
    if children and all(child.checked for child in children):
        mark_task_checked(tasks_path, parent)
        return [parent]
    return []


def git_changed_paths(runner: Runner) -> set[str]:
    tracked = runner.output(["git", "diff", "--name-only", "HEAD"])
    untracked = runner.output(
        ["git", "ls-files", "--others", "--exclude-standard"]
    )
    return {
        line.strip()
        for line in (tracked + "\n" + untracked).splitlines()
        if line.strip()
    }


def git_staged_paths(runner: Runner) -> set[str]:
    output = runner.output(["git", "diff", "--cached", "--name-only"])
    return {line for line in output.splitlines() if line}


def require_clean_tree(runner: Runner) -> None:
    status = runner.output(
        ["git", "status", "--porcelain=v1", "--untracked-files=all"]
    )
    if status:
        raise FlowError(
            "working tree must be clean before an automated step:\n" + status
        )


def matches_rule(path: str, rule: str) -> bool:
    if rule.endswith("/"):
        return path.startswith(rule)
    if any(char in rule for char in "*?["):
        return fnmatch.fnmatch(path, rule)
    return path == rule


def feature_owned_rules() -> list[str]:
    return [
        "internal/apimeta/",
        "internal/apiref/",
        "internal/apicond/",
        "internal/apiproblem/",
        "internal/apivalid/",
        "internal/apischema/",
        "internal/apiconform/",
        "api/schemas/",
        "tests/conformance/",
        "docs/api/",
        "docs/reviews/",
        "scripts/api-conformance-check.sh",
        "scripts/feature-gate.sh",
        ".github/CODEOWNERS",
        "CODEOWNERS",
    ]


def allowed_rules(task_id: str) -> list[str]:
    common = [str(TASKS_REL), str(STATE_REL)]
    group = task_group(task_id)
    if task_id == "1.1":
        return common + ["go.mod", "go.sum"]
    if task_id == "1.2":
        return common + [
            "internal/apimeta/doc.go",
            "internal/apiref/doc.go",
            "internal/apicond/doc.go",
            "internal/apiproblem/doc.go",
            "internal/apivalid/doc.go",
            "internal/apischema/doc.go",
            "internal/apiconform/doc.go",
        ]
    if task_id == "1.3":
        return common + ["internal/apiconform/imports_test.go"]
    if group == "2":
        return common + [
            "internal/apimeta/",
            "internal/apiref/",
            "internal/apicond/",
            "internal/apiproblem/",
            "internal/apivalid/",
        ]
    if group == "4":
        return common + ["internal/apivalid/"]
    if group == "5":
        return common + ["internal/apischema/", "internal/apivalid/"]
    if group == "6":
        return common + ["internal/apivalid/"]
    if group == "7a":
        return common + ["internal/apiconform/"]
    if group == "8":
        return common + ["internal/apiconform/", "internal/apivalid/"]
    if group == "9":
        return common + ["internal/apischema/"]
    if task_id == "10.1":
        # Task 10.1 creates shared schemas and may add the focused
        # executable conformance test that validates schema support and
        # field-policy completeness for those schemas.
        return common + [
            "api/schemas/_common/",
            "internal/apiconform/common_schemas_test.go",
        ]
    if task_id == "10.2":
        # Task 10.2 creates the eight canonical resource schemas and may add
        # the focused executable conformance test that verifies their
        # annotations, field policies, registry loading, and structure.
        return common + [
            "api/schemas/project.json",
            "api/schemas/resource-pool.json",
            "api/schemas/discovered-database.json",
            "api/schemas/plugin-definition.json",
            "api/schemas/adapter-configuration.json",
            "api/schemas/placement-evaluation-request.json",
            "api/schemas/operation.json",
            "api/schemas/audit-event.json",
            "internal/apiconform/canonical_schemas_test.go",
        ]
    if group == "10":
        return common + ["api/schemas/"]
    if group == "11":
        return common + ["internal/apischema/"]
    if group == "12":
        return common + [
            "internal/apiconform/",
            "docs/api/",
            "scripts/*boundary*ledger*",
            "scripts/*ledger*",
        ]
    if group == "14":
        return common + ["tests/conformance/", "internal/apiconform/"]
    if group == "15":
        return common + ["docs/api/", "internal/apiconform/"]
    if group == "16":
        return common + [
            "internal/apiconform/",
            "scripts/api-conformance-check.sh",
            "scripts/feature-gate.sh",
        ]
    if task_id == "17.1":
        return common + [
            ".github/CODEOWNERS",
            "CODEOWNERS",
            "docs/api/",
            "docs/governance/",
            "docs/reviews/",
        ]
    if task_id == "17.2":
        return common + feature_owned_rules()
    if task_id == "17.3":
        return common + [
            "docs/reviews/",
            "docs/api/",
            ".automation/evidence/FEATURE-0012/",
        ]
    if task_id in CHECKPOINTS:
        return common
    raise FlowError(f"no path policy defined for task {task_id}")


def validate_paths(task_id: str, paths: set[str]) -> None:
    rules = allowed_rules(task_id)
    disallowed = sorted(
        path for path in paths if not any(matches_rule(path, rule) for rule in rules)
    )
    if disallowed:
        raise FlowError(
            f"task {task_id} changed path(s) outside its safe allowlist:\n  "
            + "\n  ".join(disallowed)
        )

    immutable = {
        str(REQUIREMENTS_REL),
        str(DESIGN_REL),
        str(ARCH_REL),
    }
    immutable_changes = sorted(paths & immutable)
    if immutable_changes:
        raise FlowError(
            "approved immutable specification/architecture changed:\n  "
            + "\n  ".join(immutable_changes)
        )

    if task_id != "1.1" and ({"go.mod", "go.sum"} & paths):
        raise FlowError(
            f"task {task_id} changed go.mod/go.sum; FEATURE-0012 authorizes no new dependencies"
        )


def validate_tasks_file_change(
    before: str,
    after: str,
    current_task: str,
) -> None:
    if normalize_task_checkboxes(before) != normalize_task_checkboxes(after):
        raise FlowError(
            "Cursor changed tasks.md content beyond checkbox state; review required"
        )

    before_states = checkbox_states(before)
    after_states = checkbox_states(after)

    changed = {
        task_id
        for task_id in before_states.keys() | after_states.keys()
        if before_states.get(task_id) != after_states.get(task_id)
    }

    allowed = {current_task}

    # Cursor may roll up the current task's parent itself. Accept that only
    # when the parent changes from unchecked to checked and every child of
    # that parent is checked in the resulting tasks file.
    if current_task not in CHECKPOINTS and "." in current_task:
        parent = task_group(current_task)

        if before_states.get(parent) != after_states.get(parent):
            children = sorted(
                task_id
                for task_id in after_states
                if task_id.startswith(parent + ".")
            )

            valid_parent_rollup = (
                before_states.get(parent) is False
                and after_states.get(parent) is True
                and bool(children)
                and all(after_states[child] for child in children)
            )

            if not valid_parent_rollup:
                raise FlowError(
                    f"Cursor changed parent checkbox {parent} without a valid "
                    "all-children-complete roll-up"
                )

            allowed.add(parent)

    if not changed.issubset(allowed):
        raise FlowError(
            "Cursor changed checkbox(es) other than the current task or its "
            "valid completed parent: "
            + ", ".join(sorted(changed))
        )

    if before_states.get(current_task):
        raise FlowError(f"task {current_task} was already checked before execution")

def validate_precompleted_waves(plan: Plan) -> None:
    previous_complete = True
    for wave_index, wave in enumerate(plan.waves):
        checked = [plan.tasks[task_id].checked for task_id in wave]
        if any(checked) and not previous_complete:
            raise FlowError(
                f"wave {wave_index} contains completed tasks while an earlier wave is incomplete"
            )
        previous_complete = previous_complete and all(checked)


def ensure_prerequisites(
    runner: Runner, plan: Plan, selected: list[str]
) -> dict:
    state = read_json(runner.root / STATE_REL)
    if state.get("feature_id") != FEATURE:
        raise FlowError("feature state is not FEATURE-0012")
    if state.get("slug") != EXPECTED_SLUG:
        raise FlowError("feature state slug does not match FEATURE-0012")
    if state.get("feature_branch") != EXPECTED_BRANCH:
        raise FlowError("feature state branch does not match FEATURE-0012")
    branch = runner.output(["git", "branch", "--show-current"])
    if branch != EXPECTED_BRANCH:
        raise FlowError(f"wrong branch: expected {EXPECTED_BRANCH}, found {branch}")
    if state.get("current_stage") != "cursor":
        raise FlowError("feature is not at the cursor stage")
    if state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        raise FlowError("missing APPROVED_FOR_CURSOR task approval")
    human_gate = str(state.get("human_gate_required", "false")).lower() not in {
        "false",
        "0",
    }
    completed_automation = (
        not selected
        and state.get("status")
        in {"PENDING_FINAL_CHECKPOINT_REVIEW", "PENDING_HUMAN_REVIEW"}
    )
    if human_gate and not completed_automation:
        raise FlowError(
            "a human gate is currently required; use the explicit final "
            "checkpoint target after review"
        )
    validate_precompleted_waves(plan)
    runner.run(["docker", "info"], capture=True)
    if not (
        shutil_which("cursor-agent")
        or shutil_which("cursor")
        or os.environ.get("CURSOR_AGENT_BIN")
    ):
        raise FlowError("Cursor CLI is not available")
    return state


def shutil_which(command: str) -> str | None:
    from shutil import which

    return which(command)


def run_generic_gate(runner: Runner) -> None:
    runner.run(["git", "diff", "--check"])
    runner.docker_shell(
        'test -z "$(gofmt -l .)"\n'
        "go vet ./...\n"
        "go test ./..."
    )
    runner.run(["./scripts/guardrails.sh", "--feature", FEATURE])


def run_make_fmt_without_drift(runner: Runner) -> None:
    before = git_changed_paths(runner)
    if before:
        raise FlowError("checkpoint started with a dirty tree")
    runner.run(["make", "fmt"])
    after = git_changed_paths(runner)
    if after:
        raise FlowError(
            "make fmt changed files at a checkpoint; fix and commit them in the responsible task:\n  "
            + "\n  ".join(sorted(after))
        )


def checkpoint_shell(checkpoint: str) -> str:
    primitives = (
        "./internal/apimeta/... ./internal/apiref/... ./internal/apicond/... "
        "./internal/apiproblem/..."
    )
    all_grammar = (
        "./internal/apimeta/... ./internal/apiref/... ./internal/apicond/... "
        "./internal/apiproblem/... ./internal/apivalid/... "
        "./internal/apischema/... ./internal/apiconform/..."
    )
    if checkpoint == "3":
        packages = primitives + " ./internal/apivalid/... ./internal/apiconform/..."
    elif checkpoint == "7":
        packages = all_grammar
    elif checkpoint == "13":
        packages = "./internal/apivalid/... ./internal/apischema/... ./internal/apiconform/..."
    elif checkpoint == "18":
        return (
            'test -z "$(gofmt -l .)"\n'
            "go test ./...\n"
            "go test -race ./...\n"
            "go vet ./..."
        )
    else:
        raise FlowError(f"unsupported checkpoint: {checkpoint}")
    return (
        'test -z "$(gofmt -l .)"\n'
        f"go test {packages}\n"
        f"go test -race {packages}\n"
        f"go vet {packages}"
    )


def run_checkpoint_gate(runner: Runner, checkpoint: str) -> None:
    require_clean_tree(runner)
    run_make_fmt_without_drift(runner)
    runner.run(["git", "diff", "--check"])
    runner.docker_shell(checkpoint_shell(checkpoint))
    runner.run(["./scripts/guardrails.sh", "--feature", FEATURE])
    if checkpoint == "18":
        runner.run(["make", "ff-feature-gate", f"FEATURE={FEATURE}"])


def safe_commit(
    runner: Runner,
    task_id: str,
    message: str,
    *,
    pending_final_checkpoint: bool = False,
    final_pending_review: bool = False,
) -> str:
    values: dict[str, object] = {
        "current_task": task_id,
        "last_verified_task": task_id,
        "last_committed_task": task_id,
        "automation_flow_status": "running",
        "automation_last_machine_gate": task_id,
        "automation_last_machine_gate_result": "PASS",
    }
    if task_id in CHECKPOINTS:
        values["status"] = f"checkpoint_{task_id}_passed"
    else:
        values["status"] = f"cursor_task_{task_id}_committed"
    if pending_final_checkpoint:
        values.update(
            {
                "status": "PENDING_HUMAN_REVIEW",
                "human_gate_required": "true",
                "automation_flow_status": "awaiting_final_checkpoint_review",
                "automation_tasks_completed_at": utc_now(),
            }
        )
    if final_pending_review:
        values.update(
            {
                "status": "PENDING_HUMAN_REVIEW",
                "human_gate_required": "true",
                "automation_flow_status": "final_checkpoint_passed_pending_approval",
                "automation_completed_at": utc_now(),
            }
        )
    update_state(runner.root, **values)

    paths = git_changed_paths(runner)
    if not paths:
        raise FlowError(f"nothing to commit for task {task_id}")
    validate_paths(task_id, paths)
    runner.run(["git", "add", "-A", "--", *sorted(paths)])
    runner.run(["git", "diff", "--cached", "--check"])
    staged = git_staged_paths(runner)
    if staged != paths:
        raise FlowError(
            "staged path set does not match validated change set:\n"
            f"validated={sorted(paths)}\n"
            f"staged={sorted(staged)}"
        )
    runner.run(["git", "commit", "-m", message])
    require_clean_tree(runner)
    return runner.output(["git", "rev-parse", "HEAD"])


def write_gate_log(
    runner: Runner,
    task_id: str,
    kind: str,
    paths: Iterable[str],
    result: str,
    details: dict | None = None,
) -> None:
    log_dir = runner.root / ".automation/logs" / FEATURE / "flow"
    log_dir.mkdir(parents=True, exist_ok=True)
    payload = {
        "feature": FEATURE,
        "task": task_id,
        "kind": kind,
        "result": result,
        "timestamp": utc_now(),
        "paths": sorted(paths),
        "human_approval": False,
    }
    if details:
        payload.update(details)
    write_json(log_dir / f"{task_id}.json", payload)


def execute_cursor_task(runner: Runner, task_id: str, plan: Plan) -> None:
    require_clean_tree(runner)
    tasks_path = runner.root / TASKS_REL
    before_tasks = tasks_path.read_text(encoding="utf-8")
    before_states = checkbox_states(before_tasks)
    if before_states.get(task_id):
        raise FlowError(f"refusing to rerun completed task {task_id}")

    env = {
        "FEATURE_FACTORY_CURSOR_MODE": "auto",
        "FEATURE_FACTORY_CURSOR_VERIFY": "0",
        "GO_DOCKER_IMAGE": runner.image,
    }
    update_state(
        runner.root,
        automation_flow_status="running",
        automation_current_task=task_id,
        automation_started_at=utc_now(),
    )
    try:
        runner.run(
            [
                "make",
                "ff-cursor-task-auto",
                f"FEATURE={FEATURE}",
                f"TASK={task_id}",
            ],
            env=env,
        )
        after_tasks = tasks_path.read_text(encoding="utf-8")
        validate_tasks_file_change(before_tasks, after_tasks, task_id)
        paths = git_changed_paths(runner)
        validate_paths(task_id, paths)
        non_control_paths = paths - {str(TASKS_REL), str(STATE_REL)}
        if not non_control_paths and task_id not in CONTROL_ONLY_TASKS:
            raise FlowError(
                f"task {task_id} produced no implementation/evidence change outside control files"
            )
        run_generic_gate(runner)
        mark_task_checked(tasks_path, task_id)
        rolled_up = roll_up_parent(tasks_path, task_id)
        final_paths = git_changed_paths(runner)
        validate_paths(task_id, final_paths)
        write_gate_log(
            runner,
            task_id,
            "cursor-task",
            final_paths,
            "PASS",
            {"rolled_up_parents": rolled_up},
        )
        title = plan.tasks[task_id].title
        message = f"feat(feature-0012): complete task {task_id} {title}"[:180]
        commit = safe_commit(
            runner,
            task_id,
            message,
            pending_final_checkpoint=(task_id == "17.3"),
        )
        print(f"==> Task {task_id} committed: {commit[:12]}")
    except Exception as exc:
        update_state(
            runner.root,
            status=f"automation_task_{task_id}_failed",
            automation_flow_status="failed",
            automation_failed_task=task_id,
            automation_failure=str(exc),
        )
        write_gate_log(
            runner,
            task_id,
            "cursor-task",
            git_changed_paths(runner),
            "FAIL",
            {"error": str(exc)},
        )
        raise


def execute_checkpoint(runner: Runner, checkpoint: str) -> None:
    require_clean_tree(runner)
    tasks_path = runner.root / TASKS_REL
    plan = load_plan(tasks_path)
    if plan.tasks[checkpoint].checked:
        print(f"==> Checkpoint {checkpoint} already complete; skipping")
        return
    try:
        run_checkpoint_gate(runner, checkpoint)
        mark_task_checked(tasks_path, checkpoint)
        final_pending = checkpoint == "18"
        paths = git_changed_paths(runner)
        validate_paths(checkpoint, paths)
        write_gate_log(
            runner,
            checkpoint,
            "checkpoint",
            paths,
            "PASS",
        )
        commit = safe_commit(
            runner,
            checkpoint,
            f"chore(feature-0012): pass checkpoint {checkpoint}",
            final_pending_review=final_pending,
        )
        print(f"==> Checkpoint {checkpoint} committed: {commit[:12]}")
    except Exception as exc:
        update_state(
            runner.root,
            status=f"checkpoint_{checkpoint}_failed",
            automation_flow_status="failed",
            automation_failed_task=checkpoint,
            automation_failure=str(exc),
        )
        write_gate_log(
            runner,
            checkpoint,
            "checkpoint",
            git_changed_paths(runner),
            "FAIL",
            {"error": str(exc)},
        )
        raise


def plan_rows(plan: Plan) -> list[tuple[int, str, bool, bool]]:
    rows: list[tuple[int, str, bool, bool]] = []
    for wave_index, wave in enumerate(plan.waves):
        for task_id in wave:
            rows.append(
                (
                    wave_index,
                    task_id,
                    plan.tasks[task_id].checked,
                    task_id in CHECKPOINTS,
                )
            )
    return rows


def print_plan(plan: Plan) -> None:
    print(f"FEATURE-0012 execution plan: {len(plan.waves)} waves")
    for wave_index, task_id, checked, checkpoint in plan_rows(plan):
        status = "done" if checked else "pending"
        kind = "checkpoint" if checkpoint else "cursor"
        print(f"  wave {wave_index:02d}  {task_id:>5}  {kind:10}  {status}")


def select_steps(
    plan: Plan,
    start_task: str | None,
    stop_after: str | None,
    max_steps: int,
) -> list[str]:
    ordered = plan.ordered_ids
    for label, value in (("start-task", start_task), ("stop-after", stop_after)):
        if value and value not in ordered:
            raise FlowError(f"unknown {label}: {value}")

    start_index = ordered.index(start_task) if start_task else 0
    selected: list[str] = []
    for task_id in ordered[start_index:]:
        if plan.tasks[task_id].checked:
            if stop_after == task_id:
                break
            continue
        selected.append(task_id)
        if stop_after == task_id:
            break
        if max_steps and len(selected) >= max_steps:
            break
    return selected


def verify_start_dependencies(plan: Plan, selected: list[str]) -> None:
    if not selected:
        return
    target = selected[0]
    target_wave = next(
        index for index, wave in enumerate(plan.waves) if target in wave
    )
    for index in range(target_wave):
        incomplete = [
            task_id
            for task_id in plan.waves[index]
            if not plan.tasks[task_id].checked
        ]
        if incomplete:
            raise FlowError(
                f"cannot start at {target}; wave {index} is incomplete: "
                + ", ".join(incomplete)
            )


def ensure_final_checkpoint_prerequisites(runner: Runner, plan: Plan) -> None:
    state = read_json(runner.root / STATE_REL)
    if state.get("feature_id") != FEATURE:
        raise FlowError("feature state is not FEATURE-0012")
    if runner.output(["git", "branch", "--show-current"]) != EXPECTED_BRANCH:
        raise FlowError(f"final checkpoint requires branch {EXPECTED_BRANCH}")
    if state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        raise FlowError("missing APPROVED_FOR_CURSOR task approval")
    incomplete_before_final = [
        task_id
        for task_id in plan.ordered_ids
        if task_id != "18" and not plan.tasks[task_id].checked
    ]
    if incomplete_before_final:
        raise FlowError(
            "final checkpoint cannot run; prior task(s) remain incomplete: "
            + ", ".join(incomplete_before_final)
        )
    if plan.tasks["18"].checked:
        raise FlowError("final checkpoint 18 is already complete")
    allowed_statuses = {
        "PENDING_HUMAN_REVIEW",
        # Retained only for recovery from any earlier experimental run.
        "PENDING_FINAL_CHECKPOINT_REVIEW",
        "checkpoint_18_failed",
    }
    if state.get("status") not in allowed_statuses:
        raise FlowError(
            "final checkpoint requires completed automated tasks and a pending "
            f"review state; found {state.get('status')!r}"
        )
    runner.run(["docker", "info"], capture=True)


def final_summary(runner: Runner) -> None:
    plan = load_plan(runner.root / TASKS_REL)
    incomplete = [task for task in plan.ordered_ids if not plan.tasks[task].checked]
    state = read_json(runner.root / STATE_REL)
    summary = {
        "feature": FEATURE,
        "completed_at": utc_now(),
        "incomplete_tasks": incomplete,
        "status": state.get("status"),
        "human_gate_required": state.get("human_gate_required"),
        "head": runner.output(["git", "rev-parse", "HEAD"]),
    }
    write_json(
        runner.root / ".automation/logs" / FEATURE / "flow" / "summary.json",
        summary,
    )
    print("\n==> FEATURE-0012 automation finished")
    print(f"    status: {state.get('status')}")
    print(f"    human_gate_required: {state.get('human_gate_required')}")
    print(f"    incomplete tasks: {len(incomplete)}")
    if incomplete:
        print("    remaining: " + ", ".join(incomplete))
    else:
        print("    All automated tasks and checkpoints passed.")
        print("    Final human semantic review is required; no approval token was written.")


def self_test(tasks_path: Path) -> None:
    plan = load_plan(tasks_path)
    ids = plan.ordered_ids
    required = {"1.1", "1.2", "6.5a", "7a.1", "17.3", "18"}
    missing = required - set(ids)
    if missing:
        raise FlowError(f"self-test missing ids: {sorted(missing)}")
    if len(ids) != len(set(ids)):
        raise FlowError("self-test found duplicate dependency ids")
    if len(plan.waves) != 33:
        raise FlowError(f"expected 33 waves, found {len(plan.waves)}")
    if set(CHECKPOINTS) != {
        task for task in ids if task in CHECKPOINTS
    }:
        raise FlowError("checkpoint self-test failed")
    print(
        f"PASS: parsed {len(ids)} executable steps across {len(plan.waves)} waves; "
        f"checkpoints={sorted(CHECKPOINTS)}"
    )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Safe FEATURE-0012 leaf-task and checkpoint orchestrator"
    )
    mode = parser.add_mutually_exclusive_group(required=False)
    mode.add_argument("--plan", action="store_true", help="print the resume plan")
    mode.add_argument("--run", action="store_true", help="execute the plan")
    mode.add_argument("--self-test", action="store_true", help="validate parser/invariants")
    mode.add_argument(
        "--final-checkpoint",
        action="store_true",
        help="run manually confirmed final checkpoint 18",
    )
    parser.add_argument("--feature", default=FEATURE)
    parser.add_argument("--start-task")
    parser.add_argument("--stop-after")
    parser.add_argument("--max-steps", type=int, default=0)
    parser.add_argument("--image", default=os.environ.get("GO_DOCKER_IMAGE", DEFAULT_IMAGE))
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    if args.feature != FEATURE:
        raise FlowError(f"this orchestrator supports only {FEATURE}")
    root = repo_root()
    runner = Runner(root=root, image=args.image)
    tasks_path = root / TASKS_REL
    plan = load_plan(tasks_path)

    if args.self_test:
        self_test(tasks_path)
        return 0
    if args.plan or not (args.run or args.final_checkpoint):
        print_plan(plan)
        return 0

    if args.run and os.environ.get("CONFIRM_FEATURE_0012_AUTORUN") != "YES":
        raise FlowError(
            "refusing unattended execution; set CONFIRM_FEATURE_0012_AUTORUN=YES"
        )
    if (
        args.final_checkpoint
        and os.environ.get("CONFIRM_FEATURE_0012_FINAL_CHECKPOINT") != "YES"
    ):
        raise FlowError(
            "final checkpoint requires explicit human confirmation; set "
            "CONFIRM_FEATURE_0012_FINAL_CHECKPOINT=YES"
        )

    lock_path = root / ".git" / "feature-0012-flow.lock"
    lock_path.parent.mkdir(parents=True, exist_ok=True)
    with lock_path.open("w", encoding="utf-8") as lock:
        try:
            fcntl.flock(lock.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)
        except BlockingIOError as exc:
            raise FlowError("another FEATURE-0012 flow is already running") from exc

        require_clean_tree(runner)
        if args.final_checkpoint:
            ensure_final_checkpoint_prerequisites(runner, plan)
            execute_checkpoint(runner, "18")
            final_summary(runner)
            return 0

        selected = select_steps(
            plan, args.start_task, args.stop_after, args.max_steps
        )
        # Checkpoint 18 is intentionally excluded from unattended execution.
        # It is the single final manual review/checkpoint requested for the
        # feature and runs only through --final-checkpoint.
        selected = [task_id for task_id in selected if task_id != "18"]
        ensure_prerequisites(runner, plan, selected)
        verify_start_dependencies(plan, selected)
        if not selected:
            print("==> No pending automated steps selected")
            final_summary(runner)
            return 0

        print("==> Selected automated steps: " + ", ".join(selected))
        for task_id in selected:
            current_plan = load_plan(tasks_path)
            if current_plan.tasks[task_id].checked:
                print(f"==> {task_id} became complete; skipping")
                continue
            if task_id in CHECKPOINTS:
                execute_checkpoint(runner, task_id)
            else:
                execute_cursor_task(runner, task_id, current_plan)
        final_summary(runner)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except FlowError as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1)
