#!/usr/bin/env python3
import argparse, json, subprocess
from pathlib import Path
TEMPLATES={'requirements':'docs/prompts/kiro/requirements.prompt.md','design':'docs/prompts/kiro/design.prompt.md','tasks':'docs/prompts/kiro/tasks.prompt.md'}
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

def main():
    parser=argparse.ArgumentParser(); parser.add_argument('--feature',required=True); parser.add_argument('--stage',required=True,choices=TEMPLATES.keys()); args=parser.parse_args()
    state=load_state(args.feature); template_path=Path(TEMPLATES[args.stage])
    values={'FEATURE_ID':state['feature_id'],'FEATURE_SLUG':state['slug'],'FEATURE_TITLE':state['title'],'PHASE_BRANCH':state['phase_branch'],'FEATURE_BRANCH':state['feature_branch'],'SPEC_PATH':state['spec_path'],'REQUIREMENTS_PATH':f"{state['spec_path']}/requirements.md",'DESIGN_PATH':f"{state['spec_path']}/design.md",'TASKS_PATH':f"{state['spec_path']}/tasks.md",'MODEL_RECOMMENDATIONS':model_recommendation(args.stage)}
    out_dir=Path(state['generated_prompt_path']); out_dir.mkdir(parents=True,exist_ok=True)
    out_file=out_dir/f'{args.stage}.prompt.md'; out_file.write_text(render(template_path.read_text(), values)); print(out_file)
if __name__=='__main__': main()
