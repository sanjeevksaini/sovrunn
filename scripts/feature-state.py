#!/usr/bin/env python3
import argparse, json
from pathlib import Path
from datetime import datetime, timezone
STATE_DIR = Path('.automation/state')
FEATURE_DIR = Path('.automation/features')
def now(): return datetime.now(timezone.utc).isoformat()
def path_for(feature): return STATE_DIR / f'{feature}.json'
def load(feature):
    p=path_for(feature)
    if not p.exists(): raise SystemExit(f'ERROR: missing state file: {p}')
    return json.loads(p.read_text())
def save(feature, data):
    STATE_DIR.mkdir(parents=True, exist_ok=True)
    data['updated_at']=now()
    path_for(feature).write_text(json.dumps(data, indent=2, sort_keys=True)+'\n')
def init(args):
    STATE_DIR.mkdir(parents=True, exist_ok=True); FEATURE_DIR.mkdir(parents=True, exist_ok=True)
    data={'feature_id':args.feature,'slug':args.slug,'title':args.title,'phase_branch':args.phase_branch,'feature_branch':args.branch,'spec_path':f'.kiro/specs/{args.slug}','generated_prompt_path':f'docs/generated-prompts/{args.feature}','status':'branch_created','current_stage':'requirements','current_task':None,'human_gate_required':True,'created_at':now(),'updated_at':now()}
    save(args.feature,data)
    (FEATURE_DIR/f'{args.feature}.yaml').write_text('\n'.join([f'feature_id: {args.feature}',f'slug: {args.slug}',f'title: {args.title}',f'phase_branch: {args.phase_branch}',f'feature_branch: {args.branch}',f'spec_path: .kiro/specs/{args.slug}','status: branch_created','']))
    print(path_for(args.feature))
def get(args): print(json.dumps(load(args.feature), indent=2, sort_keys=True))
def get_value(args):
    data=load(args.feature); value=data.get(args.key)
    if value is None: raise SystemExit(f'ERROR: key not found or null: {args.key}')
    print(value)
def set_value(args):
    data=load(args.feature); data[args.key]=args.value; save(args.feature,data)
def main():
    parser=argparse.ArgumentParser(description='Sovrunn feature state helper'); sub=parser.add_subparsers(required=True)
    p=sub.add_parser('init'); p.add_argument('--feature',required=True); p.add_argument('--slug',required=True); p.add_argument('--title',required=True); p.add_argument('--phase-branch',required=True); p.add_argument('--branch',required=True); p.set_defaults(func=init)
    p=sub.add_parser('get'); p.add_argument('--feature',required=True); p.set_defaults(func=get)
    p=sub.add_parser('get-value'); p.add_argument('--feature',required=True); p.add_argument('--key',required=True); p.set_defaults(func=get_value)
    p=sub.add_parser('set'); p.add_argument('--feature',required=True); p.add_argument('--key',required=True); p.add_argument('--value',required=True); p.set_defaults(func=set_value)
    args=parser.parse_args(); args.func(args)
if __name__=='__main__': main()
