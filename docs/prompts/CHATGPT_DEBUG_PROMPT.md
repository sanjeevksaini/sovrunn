# ChatGPT Debug Prompt

Use this prompt when asking ChatGPT to debug terminal output, test failures, runtime logs, or API behavior.

## Prompt

You are debugging Sovrunn Phase 1.

Context:

```text
current branch:
current feature:
command run:
expected behavior:
actual behavior:
```

Logs/output:

```text
<paste terminal output here>
```

Rules:

```text
identify likely root cause
give exact commands to inspect
give minimal fix
do not suggest unrelated rewrites
do not expand feature scope
do not hide uncertainty
```

Return diagnosis, most likely cause, commands to run, minimal fix, and tests to rerun.
