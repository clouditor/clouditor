package assessment

// Operator represents a comparison operator used in metric configurations and comparison results.
//
// Note: The protobuf-generated structs currently use string fields for operators.
// To keep compatibility, we define Operator as a string-based enum. When interacting with
// protobuf-generated code, convert with string(OperatorEqual), etc.
//
// Valid operators are aligned with validation regex in assessment.proto:
// ^(<|>|<=|>=|==|isIn|allIn)$ (and by extension any operators we support in metrics)
// If additional operators are introduced in protobuf, extend this list accordingly.

type Operator string

const (
    OperatorEqual        Operator = "=="
    OperatorNotEqual     Operator = "!="
    OperatorGreaterThan  Operator = ">"
    OperatorLessThan     Operator = "<"
    OperatorGreaterEqual Operator = ">="
    OperatorLessEqual    Operator = "<="
    OperatorIsIn         Operator = "isIn"
    OperatorAllIn        Operator = "allIn"
)

// IsValidOperator reports whether the provided string is one of the known operators.
func IsValidOperator(op string) bool {
    switch Operator(op) {
    case OperatorEqual, OperatorNotEqual, OperatorGreaterThan, OperatorLessThan, OperatorGreaterEqual, OperatorLessEqual, OperatorIsIn, OperatorAllIn:
        return true
    default:
        return false
    }
}

// AsOperator converts a string to an Operator type (no validation performed).
func AsOperator(op string) Operator { return Operator(op) }

// String returns the string representation of the operator.
func (o Operator) String() string { return string(o) }
