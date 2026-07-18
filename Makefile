.PHONY: fmt vet test test-race build run demo clean verify-feature

APP_NAME=sovrunn-api
CONFIG=configs/sovrunn-api.local.yaml

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
	go build -o bin/$(APP_NAME) ./cmd/sovrunn-api

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

.PHONY: ff-start ff-kiro-stage ff-kiro-stage-auto ff-prompt-requirements ff-prompt-design ff-prompt-tasks ff-review ff-review-auto ff-review-route ff-approve-requirements ff-approve-design ff-approve-tasks ff-spec-flow ff-kiro-decision ff-model-recommend ff-model-record ff-commit-spec ff-cursor-task ff-verify ff-guardrails ff-commit-task ff-final ff-pr ff-state

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

