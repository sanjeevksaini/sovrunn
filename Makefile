.PHONY: fmt vet test test-race build run demo clean verify-feature

APP_NAME=sovrunn-api
CONFIG=configs/sovrunn-api.local.yaml
VERSION ?= dev
GO_VERSION ?= 1.22
GO_DOCKER_IMAGE ?= golang:$(GO_VERSION)
MODULE  = github.com/sanjeevksaini/sovrunn
LDFLAGS = -X '$(MODULE)/internal/api.buildVersion=$(VERSION)'

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

test-race:
	go test -race ./...

build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/sovrunn-api

run:
	go run ./cmd/sovrunn-api --config $(CONFIG)

demo:
	chmod +x scripts/demo_phase1.sh
	./scripts/demo_phase1.sh

verify-feature:
	./scripts/verify-feature.sh $(FEATURE)

clean:
	rm -rf bin

# Sovrunn Feature Factory Makefile targets
# Append this file to your repo Makefile, or include it from your Makefile.

PHASE_BRANCH ?= phase1-foundation

.PHONY: ff-start ff-kiro-stage ff-kiro-stage-auto ff-prompt-requirements ff-prompt-design ff-prompt-tasks ff-review ff-review-auto ff-review-route ff-approve-requirements ff-approve-design ff-approve-tasks ff-spec-flow ff-spec-report ff-kiro-decision ff-model-recommend ff-model-record ff-commit-spec ff-cursor-task ff-verify ff-guardrails ff-commit-task ff-final ff-pr ff-state ff-feature-control ff-approve-human-gate ff-controlled-run ff-semantic-delta ff-closeout ff-feature-factory-test

ff-start:
	./scripts/feature-start.sh --feature "$(FEATURE)" --slug "$(SLUG)" --title "$(TITLE)" --phase-branch "$(PHASE_BRANCH)"


ff-kiro-stage:
	./scripts/kiro-stage.sh --feature "$(FEATURE)" --stage "$(STAGE)" --mode "$${FEATURE_FACTORY_KIRO_MODE:-auto}"

ff-kiro-stage-auto:
	FEATURE_FACTORY_KIRO_MODE=auto ./scripts/kiro-stage.sh --feature "$(FEATURE)" --stage "$(STAGE)" --mode auto

ff-prompt-requirements:
	./scripts/render-prompt.py --feature "$(FEATURE)" --stage requirements

ff-prompt-design:
	./scripts/render-prompt.py --feature "$(FEATURE)" --stage design

ff-prompt-tasks:
	./scripts/render-prompt.py --feature "$(FEATURE)" --stage tasks

ff-review:
	./scripts/reviewer-stage.sh --feature "$(FEATURE)" --stage "$(STAGE)"

ff-review-auto:
	FEATURE_FACTORY_REVIEW_MODE=auto ./scripts/reviewer-stage.sh --feature "$(FEATURE)" --stage "$(STAGE)" --mode auto

ff-review-route:
	./scripts/review-and-route-stage.sh --feature "$(FEATURE)" --stage "$(STAGE)" --mode "$${FEATURE_FACTORY_REVIEW_MODE:-auto}"

ff-approve-requirements:
	./scripts/approve-stage.sh --feature "$(FEATURE)" --stage requirements

ff-approve-design:
	./scripts/approve-stage.sh --feature "$(FEATURE)" --stage design

ff-approve-tasks:
	./scripts/approve-stage.sh --feature "$(FEATURE)" --stage tasks

ff-spec-flow:
	./scripts/spec-flow.sh --feature "$(FEATURE)" --mode "$${FEATURE_FACTORY_REVIEW_MODE:-auto}" --kiro-mode "$${FEATURE_FACTORY_KIRO_MODE:-auto}"

ff-spec-report:
	@cat ".automation/reports/$(FEATURE)/spec-flow-latest.md"

ff-kiro-decision:
	./scripts/kiro-decision.sh --feature "$(FEATURE)" --stage "$(STAGE)" --question "$(QUESTION)"

ff-model-recommend:
	./scripts/model-recommend.py --tool "$(TOOL)" --stage "$(STAGE)" --task "$(TASK)" --tasks-path "$(TASKS_PATH)"

ff-model-record:
	./scripts/model-record.sh --feature "$(FEATURE)" --tool "$(TOOL)" --stage "$(STAGE)" --task "$(TASK)" --selected-model "$(SELECTED_MODEL)" --effort "$(EFFORT)" --fallback-used "$${FALLBACK_USED:-no}" --fallback-reason "$${FALLBACK_REASON:-none}"


ff-commit-spec:
	./scripts/commit-spec.sh --feature "$(FEATURE)"

ff-cursor-task:
	./scripts/cursor-task.sh --feature "$(FEATURE)" --task "$(TASK)" --mode "$${FEATURE_FACTORY_CURSOR_MODE:-auto}"

ff-cursor-task-auto:
	FEATURE_FACTORY_CURSOR_MODE=auto ./scripts/cursor-task.sh --feature "$(FEATURE)" --task "$(TASK)" --mode auto

.PHONY: ff-feature-0014-prompt ff-feature-0014-run feature-0014-cursor-boundary-check
ff-feature-0014-prompt:
	FEATURE_FACTORY_CURSOR_MODE=prompt ./scripts/cursor-task.sh --feature "FEATURE-0014" --task "$(TASK)" --mode prompt

ff-feature-0014-run:
	CURSOR_SELECTED_MODEL="$${CURSOR_SELECTED_MODEL:-cursor-grok-4.5-high-fast}" CURSOR_REQUIRE_SELECTED_MODEL=1 ./scripts/feature-0014-flow.py --start-task "$${START_TASK:-1}" $${STOP_AFTER:+--stop-after "$${STOP_AFTER}"}

feature-0014-cursor-boundary-check:
	PYTHONDONTWRITEBYTECODE=1 python3 ./scripts/feature-0014-cursor-boundary-check.py --task "$(TASK)"

ff-verify:
	./scripts/verify.sh

ff-guardrails:
	./scripts/guardrails.sh --feature "$(FEATURE)"

ff-commit-task:
	./scripts/commit-task.sh --feature "$(FEATURE)" --task "$(TASK)" --message "$(MESSAGE)"

ff-final:
	./scripts/final-verify.sh --feature "$(FEATURE)"

ff-pr:
	./scripts/create-pr.sh --feature "$(FEATURE)"

ff-state:
	./scripts/feature-state.py get --feature "$(FEATURE)"

ff-feature-control:
	./scripts/feature-control.py validate --feature "$(FEATURE)"

ff-approve-human-gate:
	./scripts/human-gate.py --feature "$(FEATURE)" --gate "$(GATE)" --approved-by "$(APPROVED_BY)" --approve

ff-controlled-run:
	./scripts/feature-orchestrator.py --feature "$(FEATURE)" $${START_TASK:+--start-task "$${START_TASK}"} $${STOP_AFTER:+--stop-after "$${STOP_AFTER}"} $${RUN_ALL:+--all-tasks}

ff-semantic-delta:
	./scripts/semantic-delta.py --current "$(CURRENT)" $${BASELINE:+--baseline "$${BASELINE}"} --json-out "$(JSON_OUT)" --markdown-out "$(MARKDOWN_OUT)"

ff-closeout:
	./scripts/feature-closeout.py --feature "$(FEATURE)" --pr "$(PR)" $${APPLY:+--apply}

ff-feature-factory-test:
	PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s tests/feature_factory -p 'test_*.py'


.PHONY: ff-task-flow
ff-task-flow:
	./scripts/task-flow.sh --feature "$(FEATURE)" --start-task "$${START_TASK:-1}"

.PHONY: ff-feature-flow
ff-feature-flow:
	./scripts/feature-flow.sh --feature "$(FEATURE)" --slug "$(SLUG)" --title "$(TITLE)" $(if $(filter file,$(origin PHASE_BRANCH)),,--phase-branch "$(PHASE_BRANCH)") $${START_TASK:+--start-task "$${START_TASK}"}

.PHONY: phase1-consistency
phase1-consistency:
	./scripts/phase1-consistency-check.sh

.PHONY: phase1-integration
phase1-integration:
	./scripts/phase1-integration-test.sh

.PHONY: ff-feature-gate
ff-feature-gate:
	@test -n "$(FEATURE)" || (echo "FEATURE is required"; exit 1)
	./scripts/feature-gate.sh $(FEATURE)

.PHONY: context-pack
context-pack:
	./scripts/context-pack.sh

.PHONY: phase2-scope-check
phase2-scope-check:
	@test -n "$(FEATURE)" || (echo "FEATURE is required"; exit 1)
	./scripts/phase2-scope-check.sh $(FEATURE)

.PHONY: arch-handoff-check
arch-handoff-check:
	@test -n "$(HANDOFF)" || (echo "HANDOFF is required"; exit 1)
	./scripts/architecture-handoff-check.sh "$(HANDOFF)"

.PHONY: feature-0013-architecture-boundary-check
feature-0013-architecture-boundary-check:
	@test -n "$(STAGE)" || (echo "STAGE is required: requirements, design, or tasks"; exit 1)
	PYTHONDONTWRITEBYTECODE=1 python3 ./scripts/feature-0013-architecture-boundary-check.py \
		--feature FEATURE-0013 --stage "$(STAGE)" --mode "$${MODE:-post}"

.PHONY: feature-0013-architecture-readiness
feature-0013-architecture-readiness:
	PYTHONDONTWRITEBYTECODE=1 python3 ./scripts/feature-0013-architecture-boundary-check.py \
		--feature FEATURE-0013 --stage requirements --mode readiness

.PHONY: feature-0014-architecture-boundary-check feature-0014-architecture-readiness
feature-0014-architecture-boundary-check:
	@test -n "$(STAGE)" || (echo "STAGE is required: requirements, design, or tasks"; exit 1)
	PYTHONDONTWRITEBYTECODE=1 python3 ./scripts/feature-0014-boundary-check.py \
		--feature FEATURE-0014 --stage "$(STAGE)" --mode "$${MODE:-post}"

feature-0014-architecture-readiness:
	PYTHONDONTWRITEBYTECODE=1 python3 ./scripts/feature-0014-boundary-check.py \
		--feature FEATURE-0014 --stage requirements --mode readiness

.PHONY: feature-0015-architecture-readiness
feature-0015-architecture-readiness:
	PYTHONDONTWRITEBYTECODE=1 python3 ./scripts/feature-0015-architecture-readiness-check.py

.PHONY: structurizr-lite
structurizr-lite:
	./scripts/structurizr-lite.sh

.PHONY: structurizr-check
structurizr-check:
	./scripts/structurizr-check.sh

.PHONY: structurizr-push
structurizr-push:
	./scripts/structurizr-push.sh

# FEATURE-0012 exact leaf-task orchestrator.
# Plan is safe/read-only. Run requires CONFIRM_FEATURE_0012_AUTORUN=YES.
.PHONY: ff-feature-0012-plan ff-feature-0012-run ff-feature-0012-final-checkpoint ff-feature-0012-flow-self-test
ff-feature-0012-plan:
	./scripts/feature-0012-flow.py --feature FEATURE-0012 --plan

ff-feature-0012-run:
	./scripts/feature-0012-flow.py \
		--feature FEATURE-0012 \
		--run \
		--start-task "$${START_TASK:-}" \
		--stop-after "$${STOP_AFTER:-}" \
		--max-steps "$${MAX_STEPS:-0}" \
		--image "$${GO_DOCKER_IMAGE:-golang:1.22}"

ff-feature-0012-final-checkpoint:
	./scripts/feature-0012-flow.py \
		--feature FEATURE-0012 \
		--final-checkpoint \
		--image "$${GO_DOCKER_IMAGE:-golang:1.22}"

ff-feature-0012-flow-self-test:
	./scripts/feature-0012-flow.py --feature FEATURE-0012 --self-test


.PHONY: ff-feature-0013-plan ff-feature-0013-flow vs000-contract-check phase2r-drift-check
ff-feature-0013-plan:
	./scripts/feature-0013-flow.py --plan

ff-feature-0013-flow:
	./scripts/feature-0013-flow.py --human-approved-for-cursor --start-task "$${START_TASK:-}" --stop-after "$${STOP_AFTER:-}" --max-tasks "$${MAX_TASKS:-0}" $${FEATURE_FACTORY_PUSH:+--push}

vs000-contract-check:
	PYTHONDONTWRITEBYTECODE=1 python3 ./scripts/vs000-contract-check.py

phase2r-drift-check:
	./scripts/phase2r-drift-check.sh
