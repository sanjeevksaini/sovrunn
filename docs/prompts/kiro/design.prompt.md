Feature:
{{FEATURE_ID}} {{FEATURE_TITLE}}

Generate design.md only.

Do not generate tasks.md yet.
Do not implement code.
Do not modify source files.
Do not change requirements.md unless you find a blocking contradiction.

Input requirements:
{{REQUIREMENTS_PATH}}

Tool-output safety constraints:
- Use fs_write only for chunks of 50 lines or fewer.
- For files longer than 50 lines, create the file with fs_write using the first chunk, then use fs_append in chunks of 50 lines or fewer.
- Do not write the entire design.md in one fs_write call.
- Do not use one very large str_replace edit.
- Split content into logical sections.
- Write one section at a time.
- After writing, read the file back and verify it is complete.

Use these repo context files:
- AGENTS.md
- README.md
- docs/engineering/go-coding-guardrails.md
- docs/engineering/ai-context-loading-standard.md
- existing implementations for similar resources

Design must include overview, resolved decisions, architecture, files, data models, interfaces, validation, API/handler design where applicable, registry/storage design where applicable, operation/audit behavior, error mapping, security/privacy, testing, verification, non-goals, and resolved design questions.

Hard constraints:
- Keep Go 1.21 compatibility.
- Do not introduce external dependencies unless requirements explicitly demand them.
- internal/api must not import internal/server.
- Do not add unrelated future scope.

Generate design.md only.
