package virtual_types

type Constraint interface {
	Equal(object interface{}) (BValue, error)
	NotEqual(object interface{}) (BValue, error)
	Inverse(subject interface{}) (Constraint, error)
	Name() string
	//Unbox() []interface{}
}
