// Package validate provides pure, deterministic FEATURE-0015 validation for
// the seven owned canonical cloud model kinds (design §5, DD-02, DD-10).
//
// Validators cover closed create-contract classification, scope-kind subset,
// immutable-field checks, reference/containment and topology graph integrity,
// assigned ISO-3166-1 alpha-2 codes, header/body preconditions, and complete
// VS0-SCHEMA-008..014 constraint validation. StageSet-compatible adapters
// plug these validators into the inherited FEATURE-0012 apivalid pipeline
// without route invocation, fork, or bypass.
//
// Functions perform no I/O and no wall-clock reads. Outputs do not depend on
// map iteration order. Validation does not import the in-memory store.
package validate
