package apivalid

// DecodeMode makes decoding operation-aware (F12-VALIDATION-002, F12-META-002,
// F12-OWNER-002, D-15). It selects the FieldPolicy that governs which ownership
// classes of fields are accepted or rejected.
type DecodeMode int

const (
	ModeCreateRequest      DecodeMode = iota // customer/operator create: reject system-owned + status
	ModeReplaceRequest                       // full replacement: reject status + immutable system fields
	ModeStatusUpdate                         // authorized controller: accept status, reject spec mutation
	ModeInternalObject                       // internal/system producer: accept system-owned fields
	ModeReadRepresentation                   // decode a stored/response object: accept all fields
)

// FieldPolicy resolves, for a DecodeMode, which field ownership classes
// (per Matrix C2: creator, system, spec-owner, status-owner) are permitted.
//
// Customer mutation modes (create, replace) reject unauthorized system/status
// fields; status-update, internal-object, read-representation, and fixture
// decoding accept them under the correct ownership rules. There is no
// unconditional rejection of status or system-owned fields.
type FieldPolicy struct {
	Mode              DecodeMode
	AllowStatus       bool
	AllowSystemOwned  bool // uid, generation, resourceVersion, timestamps
	AllowSpecMutation bool
}

// PolicyFor returns the FieldPolicy for mode per Matrix C2 ownership rules
// (D-15, F12-VALIDATION-002, F12-META-002, F12-OWNER-002).
//
//	ModeCreateRequest / ModeReplaceRequest:
//	  reject status and system-owned metadata; allow spec mutation
//	ModeStatusUpdate:
//	  allow status and system-owned metadata; reject spec mutation
//	ModeInternalObject / ModeReadRepresentation:
//	  allow status, system-owned metadata, and spec
//
// Unknown modes fail closed: all Allow* flags false.
func PolicyFor(mode DecodeMode) FieldPolicy {
	switch mode {
	case ModeCreateRequest:
		return FieldPolicy{
			Mode:              ModeCreateRequest,
			AllowStatus:       false,
			AllowSystemOwned:  false,
			AllowSpecMutation: true,
		}
	case ModeReplaceRequest:
		return FieldPolicy{
			Mode:              ModeReplaceRequest,
			AllowStatus:       false,
			AllowSystemOwned:  false,
			AllowSpecMutation: true,
		}
	case ModeStatusUpdate:
		return FieldPolicy{
			Mode:              ModeStatusUpdate,
			AllowStatus:       true,
			AllowSystemOwned:  true,
			AllowSpecMutation: false,
		}
	case ModeInternalObject:
		return FieldPolicy{
			Mode:              ModeInternalObject,
			AllowStatus:       true,
			AllowSystemOwned:  true,
			AllowSpecMutation: true,
		}
	case ModeReadRepresentation:
		return FieldPolicy{
			Mode:              ModeReadRepresentation,
			AllowStatus:       true,
			AllowSystemOwned:  true,
			AllowSpecMutation: true,
		}
	default:
		return FieldPolicy{Mode: mode}
	}
}

// String returns a stable label for the decode mode.
func (m DecodeMode) String() string {
	switch m {
	case ModeCreateRequest:
		return "ModeCreateRequest"
	case ModeReplaceRequest:
		return "ModeReplaceRequest"
	case ModeStatusUpdate:
		return "ModeStatusUpdate"
	case ModeInternalObject:
		return "ModeInternalObject"
	case ModeReadRepresentation:
		return "ModeReadRepresentation"
	default:
		return "DecodeModeUnknown"
	}
}
