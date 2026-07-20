# Architecture Decision Handoffs

This folder stores approved or proposed handoffs from architecture discussion into Kiro repo/spec updates.

A handoff is produced after architecture discussion, usually in the Sovrunn Architecture Governor ChatGPT Project.

A handoff is not authoritative until human-approved.

Required flow:

1. Discuss architecture tradeoff.
2. Produce handoff using `docs/prompts/chatgpt/architecture-decision-handoff.prompt.md`.
3. Save handoff here as `ADH-YYYY-NNN-short-title.md`.
4. Validate format with `make arch-handoff-check HANDOFF=<file>`.
5. Kiro validates the approved handoff against the Architecture Operating System.
6. Kiro updates docs/specs/traceability.
7. Cursor implements from approved Kiro tasks only.

Handoffs should use `docs/templates/ARCHITECTURE_DECISION_HANDOFF.md`.
