package address

// Type defines the Filecoin address protocol (f1, f4, Ox).
type Type int32

const (
	TypeUnknown Type = 0
	TypeF1      Type = 1
	TypeF3      Type = 2
	TypeF4      Type = 3
	Type0X      Type = 4
)

// Address represents a concrete Filecoin address instance.
type Address struct {
	Type  Type
	Value string
}

// String provides a human-readable representation of the address.
func (a Address) String() string {
	return a.Value
}
