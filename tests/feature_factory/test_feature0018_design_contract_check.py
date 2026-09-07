"""Deterministic tests for the FEATURE-0018 design-contract checker.

Authorized by ADH-2026-075. These tests exercise the standard-library-only
checker with positive (the real, hash-bound contract) and isolated negative
fixtures. Each negative fixture violates exactly one contract-consistency rule.
"""

import copy
import hashlib
import importlib.util
import json
import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
CHECKER_PATH = ROOT / "scripts/feature0018-design-contract-check.py"
CONTRACT_PATH = (
    ROOT
    / ".kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json"
)
DESIGN_PATH = (
    ROOT / ".kiro/specs/governance-iam-approval-exception-foundation/design.md"
)
sys.dont_write_bytecode = True


def load_module(name: str, path: Path):
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader
    spec.loader.exec_module(module)
    return module


checker = load_module("feature0018_design_contract_check", CHECKER_PATH)


class ContractCheckerPositiveTests(unittest.TestCase):
    def test_checker_source_compiles(self):
        compile(CHECKER_PATH.read_text(), str(CHECKER_PATH), "exec")

    def test_self_test_regressions_pass(self):
        self.assertEqual(checker._self_test(), 0)

    def test_checker_is_inert_under_active_checkpoint(self):
        # ADH-2026-077 §4/§5: while the FEATURE-0018 implementation checkpoint is
        # active, the Design Mechanics Contract check is an implementation check
        # only and is not a design-review authority. Its CLI must SKIP (rc=0)
        # against the frozen design.md baseline rather than reopen it.
        import subprocess

        self.assertTrue(checker.checkpoint_active())
        completed = subprocess.run(
            [sys.executable, str(CHECKER_PATH)],
            cwd=str(ROOT),
            capture_output=True,
            text=True,
        )
        self.assertEqual(completed.returncode, 0, completed.stderr)
        self.assertIn("SKIP:", completed.stdout)
        self.assertIn("ADH-2026-077", completed.stdout)

    def test_contract_binding_is_historical_draft_not_frozen_baseline(self):
        # The mechanics contract was hash-bound (ADH-2026-075) to the later
        # working design draft, which ADH-2026-077 discards. It intentionally no
        # longer matches the frozen baseline design.md; it is historical evidence.
        contract = json.loads(CONTRACT_PATH.read_text())
        recorded = contract["authority"]["binding"]["design_sha256"]
        frozen = hashlib.sha256(DESIGN_PATH.read_bytes()).hexdigest()
        self.assertNotEqual(
            recorded,
            frozen,
            "the mechanics contract must remain bound to the discarded draft, not the frozen baseline",
        )


class ContractCheckerNegativeTests(unittest.TestCase):
    def setUp(self):
        self.contract = json.loads(CONTRACT_PATH.read_text())
        self.design = DESIGN_PATH.read_bytes()

    def run_errors(self, contract):
        errors: list[str] = []
        checker.check_binding(errors, contract, self.design)
        checker.check_package_edges(errors, contract)
        checker.check_component_ledger_matches_edges(errors, contract)
        checker.check_api_surface(errors, contract)
        checker.check_state_editor_surface(errors, contract)
        checker.check_owner_finalization_table(errors, contract)
        checker.check_transaction_modes(errors, contract)
        checker.check_identity_rules(errors, contract)
        checker.check_conformance_and_proofs(errors, contract)
        checker.check_static_checker_contract(errors, contract)
        checker.check_design_citations(errors, contract, self.design)
        return errors

    def test_unstamped_binding_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["authority"]["binding"]["design_sha256"] = "PENDING"
        errors: list[str] = []
        checker.check_binding(errors, bad, self.design)
        self.assertTrue(any("unstamped" in e for e in errors))

    def test_wrong_binding_hash_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["authority"]["binding"]["design_sha256"] = "0" * 64
        errors: list[str] = []
        checker.check_binding(errors, bad, self.design)
        self.assertTrue(any("mismatch" in e for e in errors))

    def test_pseudo_union_factory_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["goApiSurface"]["factories"].append(
            {
                "name": "operation.NewPseudo",
                "signature": "operation.NewPseudo(x any) -> A | none",
                "permittedCallers": ["state"],
            }
        )
        errors: list[str] = []
        checker.check_api_surface(errors, bad)
        self.assertTrue(any("pseudo-Go" in e or "bare union" in e for e in errors))

    def test_factory_without_permitted_caller_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["goApiSurface"]["factories"].append(
            {
                "name": "state.NewOrphan",
                "signature": "state.NewOrphan() (state.Orphan, operation.MechanicalFailure)",
                "permittedCallers": [],
            }
        )
        errors: list[str] = []
        checker.check_api_surface(errors, bad)
        self.assertTrue(any("no permitted callers" in e for e in errors))

    def test_unknown_package_edge_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["packageEdgeLedger"]["edges"]["model"] = ["ghost"]
        errors: list[str] = []
        checker.check_package_edges(errors, bad)
        self.assertTrue(any("unknown node" in e for e in errors))

    def test_stdlib_import_as_ledger_edge_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["packageEdgeLedger"]["edges"]["model"] = ["time"]
        errors: list[str] = []
        checker.check_package_edges(errors, bad)
        self.assertTrue(any("standard-library" in e for e in errors))

    def test_package_edge_cycle_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["packageEdgeLedger"]["edges"]["operation"] = ["model"]
        bad["packageEdgeLedger"]["edges"]["model"] = ["operation"]
        errors: list[str] = []
        checker.check_package_edges(errors, bad)
        self.assertTrue(any("acyclic" in e for e in errors))

    def test_component_ledger_mismatch_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["componentDependencyLedger"]["componentToEdgeKey"]["ghost"] = "ghostedge"
        errors: list[str] = []
        checker.check_component_ledger_matches_edges(errors, bad)
        self.assertTrue(any("absent from the exhaustive package-edge ledger" in e for e in errors))

    def test_missing_registered_owner_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["ownerFinalizationTable"] = self.contract["ownerFinalizationTable"][:-1]
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(any("nine registered owners" in e for e in errors))

    def test_duplicate_sealed_change_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["ownerFinalizationTable"][1]["sealedFinalizedChange"] = bad[
            "ownerFinalizationTable"
        ][0]["sealedFinalizedChange"]
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(any("duplicate sealed finalized change" in e for e in errors))

    def test_missing_transaction_transition_fails(self):
        bad = copy.deepcopy(self.contract)
        del bad["transactionModes"]["stateMachine"]["CallerMutationTransaction"]["seal"]
        errors: list[str] = []
        checker.check_transaction_modes(errors, bad)
        self.assertTrue(any("missing state-machine transition" in e for e in errors))

    def test_missing_transaction_mode_fails(self):
        bad = copy.deepcopy(self.contract)
        del bad["transactionModes"]["stateMachine"]["EvidenceTransaction"]
        errors: list[str] = []
        checker.check_transaction_modes(errors, bad)
        self.assertTrue(any("exactly the three modes" in e for e in errors))

    def test_transaction_identity_not_independent_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["identityAndPublicationRules"]["transactionIdentity"]["independent_of"] = []
        errors: list[str] = []
        checker.check_identity_rules(errors, bad)
        self.assertTrue(any("independent of" in e for e in errors))

    def test_undefined_non_runtime_gate_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["conformanceProofArtifacts"]["nonRuntimeGates"] = [
            g
            for g in bad["conformanceProofArtifacts"]["nonRuntimeGates"]
            if g["id"] != "VS0-CF-F18-54"
        ]
        errors: list[str] = []
        checker.check_conformance_and_proofs(errors, bad)
        self.assertTrue(any("VS0-CF-F18-54" in e and "undefined" in e for e in errors))

    def test_gate_without_artifact_fails(self):
        bad = copy.deepcopy(self.contract)
        for gate in bad["conformanceProofArtifacts"]["nonRuntimeGates"]:
            if gate["id"] == "VS0-CF-F18-50":
                gate["artifacts"] = []
        errors: list[str] = []
        checker.check_conformance_and_proofs(errors, bad)
        self.assertTrue(any("no proof artifact" in e for e in errors))

    def test_reused_tombstone_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["conformanceProofArtifacts"]["codeAgnosticSafeErrorOracleCases"].append(
            "VS0-CF-F18-38"
        )
        errors: list[str] = []
        checker.check_conformance_and_proofs(errors, bad)
        self.assertTrue(any("tombstone" in e for e in errors))

    def test_state_editor_wrong_declarer_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["stateEditorSurface"]["declared_by"] = "uow"
        errors: list[str] = []
        checker.check_state_editor_surface(errors, bad)
        self.assertTrue(any("declared by state alone" in e for e in errors))

    def test_uncovered_factory_symbol_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["goApiSurface"]["factories"].append(
            {
                "name": "evidence.NewUncovered",
                "signature": "evidence.NewUncovered() (evidence.X, operation.MechanicalFailure)",
                "permittedCallers": ["uow"],
            }
        )
        errors: list[str] = []
        checker.check_static_checker_contract(errors, bad)
        self.assertTrue(any("protected-symbol ledger" in e for e in errors))

    def test_design_without_citation_fails(self):
        errors: list[str] = []
        checker.check_design_citations(errors, self.contract, b"# design with no citation\n")
        self.assertTrue(any("must cite" in e for e in errors))

    def test_placeholder_result_type_in_owner_row_fails(self):
        bad = copy.deepcopy(self.contract)
        row = bad["ownerFinalizationTable"][0]
        row["concreteDomainOutcomeType"] = "OwnerDomainOutcome"
        row["finalizeAt"] = (
            "func (i roleassign.PreparedRoleAssignmentIntent) "
            "FinalizeAt(state.MutationFinalizationPermit) "
            "operation.FinalizationOutcome[FinalizedPreparedChange, OwnerDomainOutcome]"
        )
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(any("prohibited placeholder type" in e for e in errors))

    def test_categorical_caller_selector_in_owner_row_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["ownerFinalizationTable"][0]["permittedCallerSelectors"] = [
            "registered owner packages with own owner constant"
        ]
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(
            any("categorical" in e or "literal selector" in e for e in errors)
        )

    def test_categorical_caller_label_in_factory_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["goApiSurface"]["factories"].append(
            {
                "name": "evidence.NewCategorical",
                "signature": "evidence.NewCategorical() (evidence.X, operation.MechanicalFailure)",
                "permittedCallers": ["state transaction"],
            }
        )
        errors: list[str] = []
        checker.check_api_surface(errors, bad)
        self.assertTrue(any("categorical caller label" in e for e in errors))

    def test_omitted_participant_claim_accessor_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["ownerFinalizationMethodSurface"]["finalizedChangeMethods"] = [
            m
            for m in bad["ownerFinalizationMethodSurface"]["finalizedChangeMethods"]
            if not m.startswith("ParticipantClaim(")
        ]
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(any("ParticipantClaim" in e for e in errors))

    def test_omitted_context_binding_accessor_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["ownerFinalizationMethodSurface"]["finalizedChangeMethods"] = [
            m
            for m in bad["ownerFinalizationMethodSurface"]["finalizedChangeMethods"]
            if not m.startswith("ContextBinding(")
        ]
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(any("ContextBinding" in e for e in errors))

    def test_owner_surface_interface_mismatch_fails(self):
        bad = copy.deepcopy(self.contract)
        bad["ownerFinalizationMethodSurface"]["finalizedChangeMethods"].append(
            "ExtraAccessor() bool"
        )
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(any("must exactly match" in e for e in errors))

    def test_owner_row_wrong_finalizeat_return_fails(self):
        bad = copy.deepcopy(self.contract)
        # Keep concrete types but declare a FinalizeAt whose return instantiation
        # does not match the row's finalized-change/domain-outcome types.
        bad["ownerFinalizationTable"][0]["finalizeAt"] = (
            "func (i roleassign.PreparedRoleAssignmentIntent) "
            "FinalizeAt(state.MutationFinalizationPermit) "
            "operation.FinalizationOutcome[membership.FinalizedMembershipChange, "
            "roleassign.RoleAssignmentNonPublicationOutcome]"
        )
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(any("finalizeAt must return" in e for e in errors))

    def test_owner_row_missing_concrete_domain_outcome_fails(self):
        bad = copy.deepcopy(self.contract)
        del bad["ownerFinalizationTable"][0]["concreteDomainOutcomeType"]
        errors: list[str] = []
        checker.check_owner_finalization_table(errors, bad)
        self.assertTrue(
            any("missing 'concreteDomainOutcomeType'" in e for e in errors)
        )


if __name__ == "__main__":
    unittest.main()
