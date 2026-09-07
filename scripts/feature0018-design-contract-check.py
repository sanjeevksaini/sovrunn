#!/usr/bin/env python3
"""Deterministic, standard-library-only FEATURE-0018 design-contract checker.

Authorized by ADH-2026-075 (FEATURE-0018 executable-design normalization and
deterministic contract check). This checker runs BEFORE the LLM design review and
fails closed on the exact conditions ADH-2026-075 section 3 enumerates:

  - a mismatch between the component dependency ledger and the exhaustive
    package-edge ledger;
  - undeclared or pseudo-Go API types and result variants;
  - a missing factory, permitted caller, owner-method rule, failure path or test
    seam for a declared mechanics-contract surface;
  - an undefined or duplicate transaction-state transition, lock order,
    notification point or context-identity component;
  - a repeated, repurposed or unmapped active FEATURE-0018 conformance identifier;
    and
  - an undefined proof artifact or command for an active non-runtime conformance
    gate.

Per the ADH-2026-075/076 mechanics-contract correction it additionally fails closed on:

  - a placeholder finalized-change / domain-outcome result type
    (``FinalizedPreparedChange``, ``OwnerDomainOutcome``) in any owner-finalization
    row, method surface or named result type;
  - a categorical caller/owner label (e.g. ``registered owner packages with own
    owner constant``, ``exact semantic-owner FinalizeAt``, ``state transaction``)
    used instead of a literal Go selector in any factory or owner row;
  - an owner-finalization row that omits its concrete domain-outcome type, exact
    FinalizeAt receiver/signature or literal caller selector, or whose FinalizeAt
    return instantiation does not match that row's concrete types; and
  - an owner finalized-change method surface that does not exactly equal
    ``state.ContextBoundChange.interface_methods`` (including ``ParticipantClaim()``
    and ``ContextBinding()``).

It validates repository-owned package edges only and never prohibits approved Go
standard-library imports. It also verifies the hash-bound design-contract binding:
the contract records the design.md SHA-256 it was derived from, and a mismatch
fails closed.

Standard library only. No network access. Deterministic (sorted traversal, stable
ordering).
"""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SPEC_DIR = ROOT / ".kiro/specs/governance-iam-approval-exception-foundation"
CONTRACT_PATH = SPEC_DIR / "design-mechanics-contract.json"
DESIGN_PATH = SPEC_DIR / "design.md"


def checkpoint_active() -> bool:
    """True when the FEATURE-0018 implementation checkpoint is active.

    ADH-2026-077 §4/§5: after the checkpoint, the Design Mechanics Contract and
    its checker remain available as an implementation check only; they are no
    longer a design-review authority and must not reopen the frozen design.md.
    While the checkpoint is active this checker is inert against the frozen
    baseline (which was not the draft the contract was bound to).
    """
    control_path = ROOT / ".automation/features/FEATURE-0018.control.json"
    if not control_path.is_file():
        return False
    try:
        control = json.loads(control_path.read_text())
    except (OSError, json.JSONDecodeError):
        return False
    checkpoint = control.get("checkpoint")
    return isinstance(checkpoint, dict) and checkpoint.get("active") is True


# Approved Go 1.22 standard-library import roots and inherited FEATURE-0012/0013/
# 0017 primitive packages. These are never treated as repository-owned ledger
# edges and must never be prohibited by the package-edge closed set.
STDLIB_ROOTS = frozenset(
    {
        "time",
        "sync",
        "context",
        "errors",
        "fmt",
        "sort",
        "math",
        "strings",
        "bytes",
        "encoding/json",
        "hash",
        "io",
    }
)
INHERITED_PRIMITIVES = frozenset(
    {
        "apimeta",
        "apivalid",
        "apiref",
        "apiproblem",
        "apicond",
        "apischema",
        "apiconform",
        "decision",
        "decision/bundle",
        "policyeval",
    }
)

# The nine registered semantic owners.
REGISTERED_OWNERS = (
    "roleassign",
    "membership",
    "roledefinition",
    "accessgroup",
    "approval",
    "privileged",
    "review",
    "exception",
    "governanceprofile",
)

# The exact non-runtime conformance gates that must have proof artifacts + commands.
REQUIRED_NON_RUNTIME_GATES = ("VS0-CF-F18-50", "VS0-CF-F18-54")

# Placeholder result/outcome type names and categorical caller/owner labels that
# the normalized mechanics contract must never use (ADH-2026-075 §2: no generic
# package placeholder, open-ended caller/owner category, or unnamed algebraic
# value). ADH-2026-076 keeps deterministic mechanics the sole authority for these.
PROHIBITED_PLACEHOLDER_TYPES = (
    "FinalizedPreparedChange",
    "OwnerDomainOutcome",
)
CATEGORICAL_CALLER_LABELS = (
    "registered owner packages with own owner constant",
    "exact semantic-owner finalizeat",
    "state transaction",
    "owner finalizeat",
    "originating owner finalizeat",
    "approval owner finalizeat",
    "exception owner finalizeat",
    "semantic-owner",
)

# The nine registered owners and their concrete owner-private intent/change types.
OWNER_INTENT_TYPE = {
    "roleassign": "roleassign.PreparedRoleAssignmentIntent",
    "membership": "membership.PreparedMembershipIntent",
    "roledefinition": "roledefinition.PreparedRoleDefinitionIntent",
    "accessgroup": "accessgroup.PreparedAccessGroupIntent",
    "approval": "approval.PreparedApprovalIntent",
    "privileged": "privileged.PreparedPrivilegedIntent",
    "review": "review.PreparedReviewIntent",
    "exception": "exception.PreparedExceptionIntent",
    "governanceprofile": "governanceprofile.PreparedGovernanceProfileIntent",
}
OWNER_CHANGE_TYPE = {
    "roleassign": "roleassign.FinalizedRoleAssignmentChange",
    "membership": "membership.FinalizedMembershipChange",
    "roledefinition": "roledefinition.FinalizedRoleDefinitionChange",
    "accessgroup": "accessgroup.FinalizedAccessGroupChange",
    "approval": "approval.FinalizedApprovalChange",
    "privileged": "privileged.FinalizedPrivilegedChange",
    "review": "review.FinalizedReviewChange",
    "exception": "exception.FinalizedExceptionChange",
    "governanceprofile": "governanceprofile.FinalizedGovernanceProfileChange",
}

# The exact FinalizeAt selector for every registered owner (literal caller form).
OWNER_FINALIZE_AT_SELECTORS = tuple(
    f"{intent}.FinalizeAt" for intent in OWNER_INTENT_TYPE.values()
)


class CheckError(Exception):
    """Raised for a fail-closed contract-consistency violation."""


def digest(payload: bytes) -> str:
    return hashlib.sha256(payload).hexdigest()


def load_contract(path: Path) -> dict:
    if not path.exists():
        raise CheckError(f"design-mechanics contract not found: {path}")
    try:
        return json.loads(path.read_text())
    except json.JSONDecodeError as exc:  # pragma: no cover - defensive
        raise CheckError(f"design-mechanics contract is not valid JSON: {exc}") from exc


# --- Pseudo-Go detection ---------------------------------------------------

_PSEUDO_TOKENS = (" | none", "| none", " none", "...", "…", ": any", " any)")
_BARE_UNION_RE = re.compile(r"\b(?:[A-Za-z_][\w.]*)\s+\|\s+(?:[A-Za-z_][\w.]*)")


def _is_categorical_label(value: str) -> bool:
    """True when a caller/owner value is a categorical label, not a literal selector.

    A literal selector is a concrete Go symbol path (e.g. an owner intent's
    ``roleassign.PreparedRoleAssignmentIntent.FinalizeAt`` selector, or a concrete
    transaction ``state.CallerMutationTransaction.ApplyFinalized`` selector, or a
    single lowercase package name such as ``uow``/``state``/``evidence``). A
    categorical label (e.g. ``registered owner packages with own owner constant``,
    ``exact semantic-owner FinalizeAt``, ``state transaction``, ``owner FinalizeAt``)
    is a prose category and is rejected.
    """
    if not value:
        return False
    lowered = value.lower()
    for label in CATEGORICAL_CALLER_LABELS:
        # Match the whole categorical phrase, not an incidental substring of a
        # legitimate dotted selector.
        if label in ("owner finalizeat", "state transaction", "semantic-owner"):
            # These bare-phrase labels are only categorical when they are NOT part
            # of a concrete dotted selector path.
            if label == "owner finalizeat" and lowered == "owner finalizeat":
                return True
            if label == "state transaction" and "." not in value and label in lowered:
                return True
            if label == "semantic-owner" and "." not in value and label in lowered:
                return True
            continue
        if label in lowered:
            return True
    # A multi-word phrase with spaces and no dotted symbol path is categorical.
    if " " in value and "." not in value:
        return True
    return False


def _signature_is_closed(signature: str) -> bool:
    """A signature is closed when its input/field positions carry no pseudo-Go.

    The only permitted `|` uses are in return position for the sealed failure form
    ``(A, operation.MechanicalFailure)`` or a named closed result type, which the
    contract expresses through explicit Go two-return or named-type syntax rather
    than a bare ``A | B``. Any bare ``A | B``, ``none``, ellipsis, or open ``any``
    in the signature is a pseudo-Go violation.
    """
    lowered = signature.lower()
    if "| none" in lowered or "|none" in lowered:
        return False
    if "..." in signature or "…" in signature:
        return False
    if "interface{}" in signature or " any)" in signature or " any," in signature:
        return False
    # A bare "A | B" union anywhere in a contract signature is disallowed: the
    # contract uses Go two-return `(A, operation.MechanicalFailure)` and named
    # closed result types instead.
    if _BARE_UNION_RE.search(signature):
        return False
    return True


# --- Checks ----------------------------------------------------------------


def check_binding(errors: list[str], contract: dict, design_text: bytes) -> None:
    binding = contract.get("authority", {}).get("binding", {})
    recorded = binding.get("design_sha256")
    actual = digest(design_text)
    if recorded in (None, "", "PENDING"):
        errors.append(
            "authority.binding.design_sha256 is unstamped (PENDING); the contract "
            "must record the normalized design.md SHA-256"
        )
        return
    if recorded != actual:
        errors.append(
            "hash-bound design binding mismatch: contract records "
            f"{recorded} but design.md is {actual}"
        )


def check_package_edges(errors: list[str], contract: dict) -> None:
    ledger = contract.get("packageEdgeLedger", {})
    edges = ledger.get("edges")
    if not isinstance(edges, dict) or not edges:
        errors.append("packageEdgeLedger.edges is missing or empty")
        return

    known_nodes = set(edges) | INHERITED_PRIMITIVES | STDLIB_ROOTS
    # Every edge target must be a known repository node or inherited primitive;
    # standard-library roots are permitted but should not appear as ledger edges.
    for src in sorted(edges):
        targets = edges[src]
        if not isinstance(targets, list):
            errors.append(f"package-edge ledger entry for '{src}' is not a list")
            continue
        if len(set(targets)) != len(targets):
            errors.append(f"package-edge ledger for '{src}' has duplicate targets")
        for tgt in targets:
            if tgt in STDLIB_ROOTS:
                errors.append(
                    f"standard-library import '{tgt}' must not appear as a ledger "
                    f"edge (source '{src}')"
                )
            elif tgt not in known_nodes:
                errors.append(
                    f"package-edge ledger references unknown node '{tgt}' "
                    f"(source '{src}')"
                )

    # Acyclicity assertion: verify the declared no-reverse-edges hold and that no
    # cycle exists among repository-owned edges (ignoring inherited primitives).
    repo_edges = {
        src: [t for t in targets if src in edges and t in edges]
        for src, targets in edges.items()
    }
    _assert_acyclic(errors, repo_edges)


def _assert_acyclic(errors: list[str], edges: dict[str, list[str]]) -> None:
    WHITE, GRAY, BLACK = 0, 1, 2
    color = {node: WHITE for node in edges}

    def visit(node: str, stack: list[str]) -> bool:
        color[node] = GRAY
        for nxt in sorted(edges.get(node, [])):
            if nxt not in color:
                continue
            if color[nxt] == GRAY:
                errors.append(
                    "package-edge ledger is not acyclic: cycle through "
                    + " -> ".join(stack + [node, nxt])
                )
                return False
            if color[nxt] == WHITE and not visit(nxt, stack + [node]):
                return False
        color[node] = BLACK
        return True

    for node in sorted(edges):
        if color[node] == WHITE:
            visit(node, [])


def check_component_ledger_matches_edges(errors: list[str], contract: dict) -> None:
    edges = contract.get("packageEdgeLedger", {}).get("edges", {})
    mapping = contract.get("componentDependencyLedger", {}).get("componentToEdgeKey", {})
    if not mapping:
        errors.append("componentDependencyLedger.componentToEdgeKey is missing")
        return
    for component, edge_key in sorted(mapping.items()):
        if edge_key not in edges:
            errors.append(
                "component dependency ledger references edge key "
                f"'{edge_key}' (component '{component}') absent from the "
                "exhaustive package-edge ledger"
            )
    # Every edge-ledger node must be represented in the component mapping.
    for node in sorted(edges):
        if node not in mapping.values():
            errors.append(
                f"package-edge node '{node}' has no component dependency ledger entry"
            )


def check_api_surface(errors: list[str], contract: dict) -> None:
    api = contract.get("goApiSurface", {})
    if not api:
        errors.append("goApiSurface is missing")
        return

    # Sealed failure type present and closed.
    failure = api.get("sealedFailureType", {})
    if failure.get("name") != "operation.MechanicalFailure":
        errors.append("goApiSurface.sealedFailureType must be operation.MechanicalFailure")
    if not failure.get("zero_value_is_no_failure"):
        errors.append("operation.MechanicalFailure must document zero-value == no failure")
    if failure.get("carries_public_problem_code") is not False:
        errors.append("operation.MechanicalFailure must not carry a public Problem code")

    # Named closed result types: each must declare discriminants + a returnedBy.
    result_types = {rt.get("type") for rt in api.get("namedClosedResultTypes", [])}
    for rt in api.get("namedClosedResultTypes", []):
        name = rt.get("type", "<unnamed>")
        if not rt.get("discriminants"):
            errors.append(f"named closed result type '{name}' declares no discriminants")
        returned_by = rt.get("returnedBy", "")
        if not returned_by:
            errors.append(f"named closed result type '{name}' declares no returnedBy")
        else:
            for placeholder in PROHIBITED_PLACEHOLDER_TYPES:
                if placeholder in returned_by:
                    errors.append(
                        f"named closed result type '{name}' returnedBy uses prohibited "
                        f"placeholder type '{placeholder}'"
                    )
            # A returnedBy that references a categorical owner label instead of a
            # literal FinalizeAt selector is rejected.
            for token in returned_by.split(","):
                token = token.strip()
                if token and _is_categorical_label(token):
                    errors.append(
                        f"named closed result type '{name}' returnedBy uses a "
                        f"categorical caller label: '{token}'"
                    )

    # Factories: closed signatures + non-empty, literal permitted callers.
    factory_names = set()
    for factory in api.get("factories", []):
        name = factory.get("name", "<unnamed>")
        factory_names.add(name)
        signature = factory.get("signature", "")
        if not signature:
            errors.append(f"factory '{name}' declares no signature")
        elif not _signature_is_closed(signature):
            errors.append(
                f"factory '{name}' has a pseudo-Go signature (bare union, none, "
                f"ellipsis, or open interface): {signature}"
            )
        callers = factory.get("permittedCallers")
        if not callers:
            errors.append(f"factory '{name}' declares no permitted callers")
        else:
            for caller in callers:
                if _is_categorical_label(caller):
                    errors.append(
                        f"factory '{name}' permittedCallers uses a categorical caller "
                        f"label (not a literal selector): '{caller}'"
                    )

    # Every named result type that is returned by a factory in the surface must be
    # referenced; and every factory returning a multi-arm type must map to a named
    # closed result type or the sealed failure form.
    for factory in api.get("factories", []):
        signature = factory.get("signature", "")
        # The sealed failure two-return form is always acceptable.
        if "operation.MechanicalFailure" in signature:
            continue
        # Otherwise the return must be a single concrete type or a named closed
        # result type.
        tail = signature.split(")")[-1].strip() if ")" in signature else ""
        if "|" in tail:
            errors.append(
                f"factory '{factory.get('name')}' return uses a bare union; it must "
                "return a single concrete type, a named closed result type, or the "
                "sealed failure form"
            )

    # Optional carriers must be named (no bare "| none").
    for carrier in api.get("optionalCarriers", []):
        if not carrier.get("someFactory") or not carrier.get("noneFactory"):
            errors.append(
                f"optional carrier '{carrier.get('type')}' must declare both a "
                "present and absent named factory"
            )


def check_state_editor_surface(errors: list[str], contract: dict) -> None:
    surface = contract.get("stateEditorSurface", {})
    if surface.get("interface") != "state.ContextBoundChange":
        errors.append("stateEditorSurface.interface must be state.ContextBoundChange")
    if surface.get("declared_by") != "state":
        errors.append("state.ContextBoundChange must be declared by state alone")
    if surface.get("checker_rule") != "F18-ARCH-012_CONTEXT_BOUND_CHANGE":
        errors.append("stateEditorSurface must cite F18-ARCH-012_CONTEXT_BOUND_CHANGE")
    if not surface.get("interface_methods"):
        errors.append("stateEditorSurface.interface_methods is missing")
    allow = surface.get("ownerToMethodAllowlist", {})
    if allow.get("permit_binding_extractor") != "matching owner FinalizeAt only":
        errors.append(
            "state.MutationFinalizationPermit.Binding extraction must be confined to "
            "the matching owner FinalizeAt"
        )


def check_owner_finalization_table(errors: list[str], contract: dict) -> None:
    table = contract.get("ownerFinalizationTable", [])
    owners = [row.get("ownerPackage") for row in table]
    if sorted(owners) != sorted(REGISTERED_OWNERS):
        errors.append(
            "ownerFinalizationTable must list exactly the nine registered owners; "
            f"found: {sorted(o for o in owners if o)}"
        )
    seen_intents = set()
    seen_changes = set()
    seen_outcomes = set()
    for row in table:
        owner = row.get("ownerPackage", "<unknown>")
        # Every row must declare a concrete claim owner constant, fully-qualified
        # concrete intent, finalized-change, and domain-outcome type, an exact
        # FinalizeAt receiver/signature, and literal permitted caller selector(s).
        for field in (
            "claimOwnerConstant",
            "sealedIntent",
            "sealedFinalizedChange",
            "concreteDomainOutcomeType",
            "finalizeAt",
            "permittedCallerSelectors",
        ):
            if not row.get(field):
                errors.append(f"owner '{owner}' finalization row missing '{field}'")

        intent = row.get("sealedIntent", "")
        change = row.get("sealedFinalizedChange", "")
        outcome = row.get("concreteDomainOutcomeType", "")
        finalize_at = row.get("finalizeAt", "")
        selectors = row.get("permittedCallerSelectors", []) or []

        # Reject placeholder / category types in any owner row field.
        for value in (intent, change, outcome, finalize_at):
            for placeholder in PROHIBITED_PLACEHOLDER_TYPES:
                if placeholder in value:
                    errors.append(
                        f"owner '{owner}' finalization row uses prohibited placeholder "
                        f"type '{placeholder}'"
                    )

        # Concrete types must be the registered fully-qualified owner types.
        if owner in OWNER_INTENT_TYPE and intent != OWNER_INTENT_TYPE[owner]:
            errors.append(
                f"owner '{owner}' sealedIntent must be the concrete "
                f"'{OWNER_INTENT_TYPE[owner]}'; found '{intent}'"
            )
        if owner in OWNER_CHANGE_TYPE and change != OWNER_CHANGE_TYPE[owner]:
            errors.append(
                f"owner '{owner}' sealedFinalizedChange must be the concrete "
                f"'{OWNER_CHANGE_TYPE[owner]}'; found '{change}'"
            )
        # The concrete domain-outcome type must be fully-qualified in the owner
        # package and must not be a category label.
        if outcome and not outcome.startswith(f"{owner}."):
            errors.append(
                f"owner '{owner}' concreteDomainOutcomeType must be a fully-qualified "
                f"'{owner}.'-prefixed concrete type; found '{outcome}'"
            )
        if _is_categorical_label(outcome):
            errors.append(
                f"owner '{owner}' concreteDomainOutcomeType uses a categorical label: "
                f"'{outcome}'"
            )

        # The FinalizeAt signature must name the exact intent receiver and
        # instantiate operation.FinalizationOutcome with this row's concrete
        # finalized-change and domain-outcome types (no placeholder).
        if finalize_at:
            expected_return = (
                f"operation.FinalizationOutcome[{change}, {outcome}]"
            )
            if "FinalizeAt(state.MutationFinalizationPermit)" not in finalize_at:
                errors.append(
                    f"owner '{owner}' finalizeAt must declare "
                    "FinalizeAt(state.MutationFinalizationPermit)"
                )
            if intent and intent not in finalize_at:
                errors.append(
                    f"owner '{owner}' finalizeAt receiver must name its concrete "
                    f"intent '{intent}'"
                )
            if change and outcome and expected_return not in finalize_at:
                errors.append(
                    f"owner '{owner}' finalizeAt must return "
                    f"'{expected_return}'; found '{finalize_at}'"
                )

        # Permitted caller selectors must be literal and must be this owner's
        # own FinalizeAt selector — never a categorical caller label.
        for selector in selectors:
            if _is_categorical_label(selector):
                errors.append(
                    f"owner '{owner}' permittedCallerSelectors uses a categorical "
                    f"caller label: '{selector}'"
                )
        expected_selector = f"{intent}.FinalizeAt" if intent else ""
        if expected_selector and expected_selector not in selectors:
            errors.append(
                f"owner '{owner}' permittedCallerSelectors must include its own "
                f"literal selector '{expected_selector}'"
            )

        if intent in seen_intents:
            errors.append(f"duplicate sealed intent '{intent}' in owner-finalization table")
        if change in seen_changes:
            errors.append(
                f"duplicate sealed finalized change '{change}' in owner-finalization table"
            )
        if outcome in seen_outcomes:
            errors.append(
                f"duplicate concrete domain-outcome type '{outcome}' in "
                "owner-finalization table"
            )
        seen_intents.add(intent)
        seen_changes.add(change)
        seen_outcomes.add(outcome)

    surface = contract.get("ownerFinalizationMethodSurface", {})
    # The intent method template must not reintroduce a placeholder outcome type.
    template = surface.get("intentMethodTemplate", "")
    if "FinalizeAt(state.MutationFinalizationPermit)" not in template:
        errors.append(
            "ownerFinalizationMethodSurface.intentMethodTemplate must declare "
            "FinalizeAt(state.MutationFinalizationPermit)"
        )
    for placeholder in PROHIBITED_PLACEHOLDER_TYPES:
        if placeholder in template:
            errors.append(
                "ownerFinalizationMethodSurface.intentMethodTemplate uses prohibited "
                f"placeholder type '{placeholder}'"
            )
    # legacy placeholder field must be gone.
    if "intentMethod" in surface and any(
        p in surface.get("intentMethod", "") for p in PROHIBITED_PLACEHOLDER_TYPES
    ):
        errors.append(
            "ownerFinalizationMethodSurface.intentMethod must not use a placeholder "
            "outcome type"
        )

    # The finalized-change method surface must EXACTLY match the state-declared
    # state.ContextBoundChange.interface_methods (rule 009), including
    # ParticipantClaim() and ContextBinding().
    editor_surface = contract.get("stateEditorSurface", {})
    interface_methods = list(editor_surface.get("interface_methods", []))
    surface_methods = list(surface.get("finalizedChangeMethods", []))
    if surface_methods != interface_methods:
        errors.append(
            "ownerFinalizationMethodSurface.finalizedChangeMethods must exactly match "
            "state.ContextBoundChange.interface_methods (rule 009); "
            f"surface={surface_methods} interface={interface_methods}"
        )
    # Explicitly require the two accessors the correction mandates.
    joined = " ".join(surface_methods)
    if "ParticipantClaim() state.ParticipantClaim" not in surface_methods:
        errors.append(
            "ownerFinalizationMethodSurface.finalizedChangeMethods must include "
            "ParticipantClaim() state.ParticipantClaim"
        )
    if "ContextBinding() evidence.PublicationContextBinding" not in surface_methods:
        errors.append(
            "ownerFinalizationMethodSurface.finalizedChangeMethods must include "
            "ContextBinding() evidence.PublicationContextBinding"
        )
    if surface.get("finalizedChangeMustEqualInterfaceMethods") != (
        "state.ContextBoundChange.interface_methods"
    ):
        errors.append(
            "ownerFinalizationMethodSurface must bind finalizedChangeMethods to "
            "state.ContextBoundChange.interface_methods"
        )


def check_transaction_modes(errors: list[str], contract: dict) -> None:
    modes = contract.get("transactionModes", {})
    machine = modes.get("stateMachine", {})
    expected_modes = {"CallerMutationTransaction", "ControllerMutationTransaction", "EvidenceTransaction"}
    if set(machine) != expected_modes:
        errors.append(
            "transactionModes.stateMachine must define exactly the three modes "
            f"{sorted(expected_modes)}; found {sorted(machine)}"
        )
    required_transitions = {
        "begin",
        "admission",
        "seal",
        "commit",
        "abort",
        "reservationCleanup",
        "waiterNotification",
        "resultRelease",
    }
    for mode, transitions in sorted(machine.items()):
        declared = set(transitions)
        missing = required_transitions - declared
        if missing:
            errors.append(
                f"transaction mode '{mode}' missing state-machine transition(s): "
                + ", ".join(sorted(missing))
            )
        # No duplicate/empty transition definitions.
        for name, body in transitions.items():
            if not isinstance(body, str) or not body.strip():
                errors.append(
                    f"transaction mode '{mode}' transition '{name}' has no definition"
                )

    # Lock order + notification point must be documented on the caller mode.
    caller = machine.get("CallerMutationTransaction", {})
    if "stateMu" not in caller.get("abort", "") or "reservationMu" not in caller.get("abort", ""):
        errors.append(
            "CallerMutationTransaction.abort must document the stateMu -> reservationMu "
            "lock order"
        )
    if "after" not in caller.get("waiterNotification", "").lower():
        errors.append(
            "CallerMutationTransaction.waiterNotification must document broadcast "
            "only after unlock"
        )


def check_identity_rules(errors: list[str], contract: dict) -> None:
    rules = contract.get("identityAndPublicationRules", {})
    required_components = (
        "publicationInstant",
        "publicationSequence",
        "transactionIdentity",
        "storeNonce",
        "finalizationPermitToken",
        "completedRecordRetention",
        "failureVsError",
    )
    for component in required_components:
        if component not in rules:
            errors.append(f"identityAndPublicationRules missing '{component}'")
    ident = rules.get("transactionIdentity", {})
    independent = set(ident.get("independent_of", []))
    if not {"publicationSequence", "operation.PublicationInstant"} <= independent:
        errors.append(
            "transactionIdentity must be independent of both publicationSequence and "
            "operation.PublicationInstant"
        )
    fve = rules.get("failureVsError", {})
    if not fve.get("distinguished_from_go_error"):
        errors.append(
            "identityAndPublicationRules.failureVsError must distinguish "
            "operation.MechanicalFailure from a Go error"
        )
    if fve.get("no_public_problem_semantics_change") is not True:
        errors.append("mechanical-failure result must not change public Problem semantics")


def check_conformance_and_proofs(errors: list[str], contract: dict) -> None:
    proofs = contract.get("conformanceProofArtifacts", {})
    gates = {g.get("id"): g for g in proofs.get("nonRuntimeGates", [])}
    for gate_id in REQUIRED_NON_RUNTIME_GATES:
        gate = gates.get(gate_id)
        if gate is None:
            errors.append(f"non-runtime conformance gate '{gate_id}' is undefined")
            continue
        if not gate.get("artifacts"):
            errors.append(f"non-runtime gate '{gate_id}' declares no proof artifact")
        if not gate.get("commands"):
            errors.append(f"non-runtime gate '{gate_id}' declares no proof command")
        if not gate.get("owningTask"):
            errors.append(f"non-runtime gate '{gate_id}' declares no owning task")

    # Conformance identifier hygiene: no repeated, repurposed, or unmapped IDs.
    all_ids: list[str] = []
    all_ids.extend(gates)
    tomb = proofs.get("tombstone", {})
    if tomb.get("id"):
        all_ids.append(tomb["id"])
    all_ids.extend(proofs.get("codeAgnosticSafeErrorOracleCases", []))
    duplicates = {cid for cid in all_ids if all_ids.count(cid) > 1}
    if duplicates:
        errors.append(
            "repeated FEATURE-0018 conformance identifier(s): " + ", ".join(sorted(duplicates))
        )
    # The tombstone must not be reused as an active gate or oracle case.
    if tomb.get("id") in gates or tomb.get("id") in proofs.get(
        "codeAgnosticSafeErrorOracleCases", []
    ):
        errors.append(
            f"tombstone conformance identifier '{tomb.get('id')}' is reused as an "
            "active case"
        )
    # Every referenced VS0-CF id must match the canonical pattern.
    id_pattern = re.compile(r"^VS0-CF-F18-\d{2}$")
    for cid in all_ids:
        if not id_pattern.match(cid):
            errors.append(f"malformed FEATURE-0018 conformance identifier: '{cid}'")


def check_static_checker_contract(errors: list[str], contract: dict) -> None:
    checker = contract.get("staticCheckerContract", {})
    impl = checker.get("designOwnedImplementationChecker", {})
    if not impl.get("validatesRepositoryOwnedEdgesOnly"):
        errors.append(
            "static checker must validate repository-owned package edges only"
        )
    if impl.get("mustNotProhibitApprovedStdlibImports") is not True:
        errors.append("static checker must not prohibit approved standard-library imports")
    seams = contract.get("testOnlySeams", {})
    bound = set(seams.get("boundTo", []))
    if not {"*_test.go", "testdata"} <= bound:
        errors.append("test-only seams must be bound to *_test.go or testdata")
    if seams.get("mustNotWeakenProductionCallerRestrictions") is not True:
        errors.append("test-only seams must not weaken production caller restrictions")

    # Every declared checker rule must have an id; each declared surface factory
    # must have a permitted caller (already checked) — cross-check protected
    # symbols cover the factories.
    protected = set(contract.get("staticCheckerContract", {}).get("protectedSymbols", []))
    for factory in contract.get("goApiSurface", {}).get("factories", []):
        name = factory.get("name")
        # Generic release/state helpers are protected; require the factory name to
        # be present in the protected-symbol ledger.
        if name and name not in protected:
            errors.append(
                f"mechanics-contract factory '{name}' is not covered by the "
                "static-checker protected-symbol ledger"
            )


def check_design_citations(errors: list[str], contract: dict, design_text: bytes) -> None:
    """Every authoritative design must cite the contract, not restate mechanics."""
    text = design_text.decode("utf-8", errors="replace")
    if "design-mechanics-contract.json" not in text:
        errors.append(
            "design.md must cite the authoritative Design Mechanics Contract "
            "(design-mechanics-contract.json)"
        )


def run(contract_path: Path, design_path: Path) -> list[str]:
    errors: list[str] = []
    contract = load_contract(contract_path)
    if not design_path.exists():
        raise CheckError(f"design.md not found: {design_path}")
    design_text = design_path.read_bytes()

    check_binding(errors, contract, design_text)
    check_package_edges(errors, contract)
    check_component_ledger_matches_edges(errors, contract)
    check_api_surface(errors, contract)
    check_state_editor_surface(errors, contract)
    check_owner_finalization_table(errors, contract)
    check_transaction_modes(errors, contract)
    check_identity_rules(errors, contract)
    check_conformance_and_proofs(errors, contract)
    check_static_checker_contract(errors, contract)
    check_design_citations(errors, contract, design_text)
    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--contract", default=str(CONTRACT_PATH))
    parser.add_argument("--design", default=str(DESIGN_PATH))
    parser.add_argument(
        "--self-test",
        action="store_true",
        help="run deterministic offline regression checks and exit",
    )
    args = parser.parse_args()

    if args.self_test:
        return _self_test()

    if checkpoint_active():
        print(
            "SKIP: FEATURE-0018 implementation checkpoint active (ADH-2026-077); "
            "the Design Mechanics Contract check is an implementation check only and is "
            "not a design-review authority. The frozen design.md baseline governs."
        )
        return 0

    try:
        errors = run(Path(args.contract), Path(args.design))
    except CheckError as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        return 1

    if errors:
        print(
            f"FAIL: {len(errors)} FEATURE-0018 design-contract consistency error(s):",
            file=sys.stderr,
        )
        for error in errors:
            print(f"  - {error}", file=sys.stderr)
        return 1

    print(
        "PASS: FEATURE-0018 design-mechanics contract is internally consistent, "
        "hash-bound to design.md, and cited by design.md"
    )
    return 0


def _self_test() -> int:
    """Deterministic offline regression checks on synthetic contracts."""
    base = load_contract(CONTRACT_PATH)
    failures: list[str] = []

    # 1. A pseudo-Go signature must fail.
    import copy

    bad = copy.deepcopy(base)
    bad["goApiSurface"]["factories"].append(
        {"name": "operation.NewBad", "signature": "operation.NewBad(any) -> A | none", "permittedCallers": ["state"]}
    )
    errs: list[str] = []
    check_api_surface(errs, bad)
    if not any("pseudo-Go" in e or "bare union" in e for e in errs):
        failures.append("self-test 1: pseudo-Go signature was not rejected")

    # 2. A package edge to an unknown node must fail.
    bad2 = copy.deepcopy(base)
    bad2["packageEdgeLedger"]["edges"]["model"] = ["nonexistent"]
    errs = []
    check_package_edges(errs, bad2)
    if not any("unknown node" in e for e in errs):
        failures.append("self-test 2: unknown package edge was not rejected")

    # 3. A cycle must fail.
    bad3 = copy.deepcopy(base)
    bad3["packageEdgeLedger"]["edges"]["model"] = ["operation"]
    bad3["packageEdgeLedger"]["edges"]["operation"] = ["model"]
    errs = []
    check_package_edges(errs, bad3)
    if not any("acyclic" in e for e in errs):
        failures.append("self-test 3: package-edge cycle was not rejected")

    # 4. Missing a registered owner must fail.
    bad4 = copy.deepcopy(base)
    bad4["ownerFinalizationTable"] = base["ownerFinalizationTable"][:-1]
    errs = []
    check_owner_finalization_table(errs, bad4)
    if not any("nine registered owners" in e for e in errs):
        failures.append("self-test 4: missing registered owner was not rejected")

    # 5. A reused tombstone id must fail.
    bad5 = copy.deepcopy(base)
    bad5["conformanceProofArtifacts"]["codeAgnosticSafeErrorOracleCases"].append("VS0-CF-F18-38")
    errs = []
    check_conformance_and_proofs(errs, bad5)
    if not any("tombstone" in e for e in errs):
        failures.append("self-test 5: reused tombstone identifier was not rejected")

    # 6. A missing transaction-state transition must fail.
    bad6 = copy.deepcopy(base)
    del bad6["transactionModes"]["stateMachine"]["EvidenceTransaction"]["commit"]
    errs = []
    check_transaction_modes(errs, bad6)
    if not any("missing state-machine transition" in e for e in errs):
        failures.append("self-test 6: missing transaction transition was not rejected")

    # 7. A component ledger key absent from the edge ledger must fail.
    bad7 = copy.deepcopy(base)
    bad7["componentDependencyLedger"]["componentToEdgeKey"]["ghost"] = "ghostedge"
    errs = []
    check_component_ledger_matches_edges(errs, bad7)
    if not any("absent from the exhaustive package-edge ledger" in e for e in errs):
        failures.append("self-test 7: component/edge ledger mismatch was not rejected")

    # 8. A placeholder result/outcome type in an owner FinalizeAt must fail.
    bad8 = copy.deepcopy(base)
    row8 = bad8["ownerFinalizationTable"][0]
    row8["concreteDomainOutcomeType"] = "OwnerDomainOutcome"
    row8["finalizeAt"] = (
        "func (i roleassign.PreparedRoleAssignmentIntent) "
        "FinalizeAt(state.MutationFinalizationPermit) "
        "operation.FinalizationOutcome[FinalizedPreparedChange, OwnerDomainOutcome]"
    )
    errs = []
    check_owner_finalization_table(errs, bad8)
    if not any("prohibited placeholder type" in e for e in errs):
        failures.append("self-test 8: placeholder result type was not rejected")

    # 9. A categorical caller label on an owner row must fail.
    bad9 = copy.deepcopy(base)
    bad9["ownerFinalizationTable"][0]["permittedCallerSelectors"] = [
        "registered owner packages with own owner constant"
    ]
    errs = []
    check_owner_finalization_table(errs, bad9)
    if not any("categorical" in e or "literal selector" in e for e in errs):
        failures.append("self-test 9: categorical caller label was not rejected")

    # 10. Omitting ParticipantClaim from the finalized-change surface must fail.
    bad10 = copy.deepcopy(base)
    bad10["ownerFinalizationMethodSurface"]["finalizedChangeMethods"] = [
        m
        for m in bad10["ownerFinalizationMethodSurface"]["finalizedChangeMethods"]
        if not m.startswith("ParticipantClaim(")
    ]
    errs = []
    check_owner_finalization_table(errs, bad10)
    if not any("ParticipantClaim" in e for e in errs):
        failures.append("self-test 10: omitted ParticipantClaim was not rejected")

    # 11. Omitting ContextBinding from the finalized-change surface must fail.
    bad11 = copy.deepcopy(base)
    bad11["ownerFinalizationMethodSurface"]["finalizedChangeMethods"] = [
        m
        for m in bad11["ownerFinalizationMethodSurface"]["finalizedChangeMethods"]
        if not m.startswith("ContextBinding(")
    ]
    errs = []
    check_owner_finalization_table(errs, bad11)
    if not any("ContextBinding" in e for e in errs):
        failures.append("self-test 11: omitted ContextBinding was not rejected")

    # 12. A finalized-change surface that does not equal the interface method set
    #     (owner-surface mismatch) must fail.
    bad12 = copy.deepcopy(base)
    bad12["ownerFinalizationMethodSurface"]["finalizedChangeMethods"].append(
        "ExtraAccessor() bool"
    )
    errs = []
    check_owner_finalization_table(errs, bad12)
    if not any("must exactly match" in e for e in errs):
        failures.append("self-test 12: owner-surface/interface mismatch was not rejected")

    # 13. A categorical caller label on a factory must fail.
    bad13 = copy.deepcopy(base)
    bad13["goApiSurface"]["factories"].append(
        {
            "name": "evidence.NewCategorical",
            "signature": "evidence.NewCategorical() (evidence.X, operation.MechanicalFailure)",
            "permittedCallers": ["exact semantic-owner FinalizeAt"],
        }
    )
    errs = []
    check_api_surface(errs, bad13)
    if not any("categorical caller label" in e for e in errs):
        failures.append("self-test 13: categorical factory caller was not rejected")

    if failures:
        print("SELF-TEST FAIL:", file=sys.stderr)
        for failure in failures:
            print(f"  - {failure}", file=sys.stderr)
        return 1
    print("PASS: design-contract checker self-test (13 regression checks)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
