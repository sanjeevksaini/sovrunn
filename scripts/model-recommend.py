#!/usr/bin/env python3
import argparse, json, re, sys
from pathlib import Path

POLICY_PATH = Path('.automation/model-policy.json')


def load_policy():
    if not POLICY_PATH.exists():
        raise SystemExit(f'ERROR: model policy not found: {POLICY_PATH}')
    return json.loads(POLICY_PATH.read_text())


def classify_task(policy, task_text):
    text = task_text.lower()
    for rule in policy.get('task_classification_rules', []):
        for token in rule.get('match_any', []):
            if token.lower() in text:
                return rule['profile']
    return policy.get('default_cursor_profile', 'routine_code')


def extract_task_text(tasks_path, task_id):
    p = Path(tasks_path)
    if not p.exists():
        return ''
    text = p.read_text()
    # Capture a task section starting with '- [ ] N.N' until next '- [ ] X.Y' at same level.
    pat = re.compile(rf'(?ms)^\s*- \[ \] {re.escape(task_id)}\b.*?(?=^\s*- \[ \] \d+\.\d+\b|\Z)')
    m = pat.search(text)
    return m.group(0) if m else text[:4000]


def as_markdown(tool, profile, recs):
    lines = []
    lines.append('## Model recommendation')
    lines.append('')
    lines.append(f'Tool: {tool}')
    lines.append(f'Profile: {profile}')
    lines.append('')
    lines.append('Use the first available model from this prioritized list:')
    lines.append('')
    for r in recs:
        lines.append(f"{r['priority']}. {r['label']} (`{r['model']}`), effort: {r['effort']} — {r['why']}")
    lines.append('')
    lines.append('If the first model is unavailable, use the next available model in priority order.')
    lines.append('Do not use a fourth model unless all three are unavailable or blocked by governance.')
    lines.append('')
    lines.append('At the end of execution, include exactly this report:')
    lines.append('')
    lines.append('```text')
    lines.append('Model Execution Report:')
    lines.append(f'- Tool: {tool}')
    lines.append('- Stage or task: <stage/task id>')
    lines.append('- Recommended priority list: <copy the three model labels from this prompt>')
    lines.append('- Selected model: <actual model selected>')
    lines.append('- Effort/reasoning setting: <actual setting if visible>')
    lines.append('- Fallback used: yes/no')
    lines.append('- Fallback reason: <unavailable/cost/latency/user override/none>')
    lines.append('```')
    return '\n'.join(lines)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--tool', required=True, choices=['kiro','cursor'])
    ap.add_argument('--stage', choices=['requirements','design','tasks','revision'])
    ap.add_argument('--task')
    ap.add_argument('--tasks-path')
    ap.add_argument('--format', choices=['markdown','json'], default='markdown')
    args = ap.parse_args()

    policy = load_policy()
    if args.tool == 'kiro':
        profile = args.stage or 'requirements'
    else:
        task_text = extract_task_text(args.tasks_path, args.task) if args.tasks_path and args.task else ''
        profile = classify_task(policy, task_text)

    recs = policy['tools'][args.tool]['profiles'][profile]
    result = {'tool': args.tool, 'profile': profile, 'recommendations': recs}
    if args.format == 'json':
        print(json.dumps(result, indent=2))
    else:
        print(as_markdown(args.tool, profile, recs))

if __name__ == '__main__':
    main()
