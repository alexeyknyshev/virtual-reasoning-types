package virtual_types

type Value interface {
	TypeName() string

	IsValid() bool
	IsUndefined() bool
	IsConstant() bool
	IsSame(o Value) bool

	Equal(o Value) Boolean
	NotEqual(o Value) Boolean

	ToBoolean() Boolean

	/*Less(o Value) (Boolean, error)
	LessEqual(o Value) (Boolean, error)
	Greater(o Value) (Boolean, error)
	GreaterEqual(o Value) (Boolean, error)*/
}
