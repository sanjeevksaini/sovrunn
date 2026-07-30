#!/usr/bin/env python3
import argparse, hashlib, json, subprocess
from pathlib import Path
TEMPLATES={'requirements':'docs/prompts/kiro/requirements.prompt.md','design':'docs/prompts/kiro/design.prompt.md','tasks':'docs/prompts/kiro/tasks.prompt.md'}

FEATURE_0014_DESIGN_CONTEXT = [
    Path('AGENTS.md'),
    Path('README.md'),
    Path('docs/engineering/ai-context-loading-standard.md'),
    Path('docs/foundation/constitution.md'),
    Path('docs/decisions/DECISION_INDEX.md'),
    Path('docs/glossary.md'),
    Path('docs/features/FEATURE_SEQUENCE.md'),
    Path('docs/resource-specs/RESOURCE_MODEL_PHASE1.md'),
    Path('docs/api/API_CONTRACT_PHASE1.md'),
    Path('.kiro/steering/product.md'),
    Path('.kiro/steering/architecture.md'),
    Path('.kiro/steering/engineering.md'),
    Path('docs/engineering/go-coding-guardrails.md'),
    Path('docs/engineering/go-version-standard.md'),
    Path('docs/phase2/PHASE2_SCOPE.md'),
    Path('docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md'),
    Path('docs/architecture/api-resource-standard.md'),
    Path('docs/architecture/provider-neutral-resource-model.md'),
    Path('docs/reviews/architecture-decision-handoffs/ADH-2026-018-feature-0014-provider-neutral-resource-model.md'),
    Path('docs/reviews/architecture-decision-handoffs/ADH-2026-019-feature-0014-geographic-descriptor-clarification.md'),
    Path('docs/features/FEATURE-0014-provider-neutral-resource-model.md'),
    Path('.kiro/specs/provider-neutral-resource-model/requirements.md'),
]
FEATURE_0014_TASKS_CONTEXT = [
    Path('AGENTS.md'),
    Path('docs/engineering/go-coding-guardrails.md'),
    Path('docs/engineering/go-version-standard.md'),
    Path('docs/architecture/api-resource-standard.md'),
    Path('docs/architecture/provider-neutral-resource-model.md'),
    Path('docs/reviews/architecture-decision-handoffs/ADH-2026-018-feature-0014-provider-neutral-resource-model.md'),
    Path('docs/reviews/architecture-decision-handoffs/ADH-2026-019-feature-0014-geographic-descriptor-clarification.md'),
    Path('.kiro/specs/provider-neutral-resource-model/requirements.md'),
    Path('.kiro/specs/provider-neutral-resource-model/design.md'),
]
def load_state(feature):
    p=Path(f'.automation/state/{feature}.json')
    if not p.exists(): raise SystemExit(f'ERROR: state file not found: {p}')
    return json.loads(p.read_text())
def render(template, values):
    out=template
    for k,v in values.items():
        out=out.replace('{{'+k+'}}', str(v))
    return out

def model_recommendation(stage):
    try:
        return subprocess.check_output(['./scripts/model-recommend.py','--tool','kiro','--stage',stage], text=True)
    except Exception as e:
        return f'Model recommendation unavailable: {e}'

def write_feature_0014_context_manifest(out_dir, stage):
    if stage == 'design':
        paths = FEATURE_0014_DESIGN_CONTEXT
    elif stage == 'tasks':
        paths = FEATURE_0014_TASKS_CONTEXT
    else:
        paths = [
        Path('AGENTS.md'),
        Path('docs/engineering/ai-context-loading-standard.md'),
        Path('docs/foundation/constitution.md'),
        Path('docs/decisions/DECISION_INDEX.md'),
        Path('docs/phase2/PHASE2_EXECUTION_STRATEGY.md'),
        Path('docs/phase2/PHASE2_ARCHITECTURE_SPINE.md'),
        Path('docs/phase2/PHASE2_SCOPE.md'),
        Path('docs/phase2/PHASE2_FEATURE_SEQUENCE.md'),
        Path('docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md'),
        Path('docs/architecture/api-resource-standard.md'),
        Path('docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md'),
        Path('docs/reviews/architecture-decision-handoffs/ADH-2026-013-operation-allowed-scopes.md'),
        Path('docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md'),
        Path('docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-consolidated-architecture.md'),
        Path('docs/architecture/provider-neutral-resource-model.md'),
        Path('docs/reviews/architecture-decision-handoffs/ADH-2026-018-feature-0014-provider-neutral-resource-model.md'),
        Path('docs/reviews/architecture-decision-handoffs/ADH-2026-019-feature-0014-geographic-descriptor-clarification.md'),
        Path('docs/features/FEATURE-0014-provider-neutral-resource-model.md'),
        Path('docs/features/FEATURE_INDEX.md'),
        Path('docs/context/CURRENT_ARCHITECTURE_BASELINE.md'),
        Path('docs/context/CURRENT_DECISION_SUMMARY.md'),
        Path('docs/glossary.md'),
        ]
    missing = [str(path) for path in paths if not path.is_file()]
    if missing:
        raise SystemExit('ERROR: FEATURE-0014 context manifest missing: ' + ', '.join(missing))
    files = []
    for path in paths:
        data = path.read_bytes()
        files.append({
            'path': str(path),
            'bytes': len(data),
            'lines': len(data.splitlines()),
            'sha256': hashlib.sha256(data).hexdigest(),
        })
    manifest = {'feature': 'FEATURE-0014', 'stage': stage, 'files': files}
    (out_dir / f'{stage}.context.json').write_text(json.dumps(manifest, indent=2, sort_keys=True) + '\n')

def main():
    parser=argparse.ArgumentParser(); parser.add_argument('--feature',required=True); parser.add_argument('--stage',required=True,choices=TEMPLATES.keys()); args=parser.parse_args()
    state=load_state(args.feature); template_path=Path(TEMPLATES[args.stage])
    if args.feature == 'FEATURE-0014' and args.stage == 'design':
        template_path = Path('docs/prompts/kiro/feature-0014-design.prompt.md')
    elif args.feature == 'FEATURE-0014' and args.stage == 'tasks':
        template_path = Path('docs/prompts/kiro/feature-0014-tasks.prompt.md')
    context_files = FEATURE_0014_TASKS_CONTEXT if args.feature == 'FEATURE-0014' and args.stage == 'tasks' else FEATURE_0014_DESIGN_CONTEXT
    values={'FEATURE_ID':state['feature_id'],'FEATURE_SLUG':state['slug'],'FEATURE_TITLE':state['title'],'PHASE_BRANCH':state['phase_branch'],'FEATURE_BRANCH':state['feature_branch'],'SPEC_PATH':state['spec_path'],'REQUIREMENTS_PATH':f"{state['spec_path']}/requirements.md",'DESIGN_PATH':f"{state['spec_path']}/design.md",'TASKS_PATH':f"{state['spec_path']}/tasks.md",'MODEL_RECOMMENDATIONS':model_recommendation(args.stage),'CONTEXT_FILES':'\n'.join(f'- `{path}`' for path in context_files)}
    out_dir=Path(state['generated_prompt_path']); out_dir.mkdir(parents=True,exist_ok=True)
    out_file=out_dir/f'{args.stage}.prompt.md'; out_file.write_text(render(template_path.read_text(), values))
    if args.feature == 'FEATURE-0014':
        write_feature_0014_context_manifest(out_dir, args.stage)
    print(out_file)
if __name__=='__main__': main()
