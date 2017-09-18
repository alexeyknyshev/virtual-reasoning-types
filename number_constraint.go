package virtual_types

import "fmt"

type NumberConstraint interface {
	Equal(object interface{}) (BValue, error)
	NotEqual(object interface{}) (BValue, error)
	Inverse(subject interface{}) (NumberConstraint, error)
	Name() string
	Less(object interface{}) (BValue, error)
	Greater(object interface{}) (BValue, error)
	LessEqual(object interface{}) (BValue, error)
	GreaterEqual(object interface{}) (BValue, error)
	//Unbox() []interface{}
}

func errApplyInvalidNumber(action string) error {
	return fmt.Errorf("could not apply %s constraint to non-Number value", action)
}

func not(val BValue, err error) (BValue, error) {
	switch val {
	case BTrue:
		return BFalse, err
	case BFalse:
		return BTrue, err
	}
	return val, err
}

func or(l, r BValue) BValue {
	if l == BTrue || r == BFalse {
		return l
	}
	return r
}

func unknown(object interface{}, name string) (BValue, error) {
	if _, ok := object.(*NumberPrivate); ok {
		return BUnknown, nil
	}
	return -1, errApplyInvalidNumber(name)
}

// ====== NumberOr ======

type NumberOr struct {
	variants []NumberConstraint
}

func NewNumberOr(variants ...NumberConstraint) NumberOr {
	variantsLen := len(variants)
	if variantsLen < 2 {
		panic("expected atleast 2 variants for NumberOr construction")
	}
	variantsArr := make([]NumberConstraint, variantsLen)
	for i, v := range variants {
		variantsArr[i] = v
	}
	return NumberOr{variants: variants}
}

func (_ NumberOr) Name() string {
	return "NumberOr"
}

func (c NumberOr) Equal(object interface{}) (BValue, error) {
	result := BFalse
	for _, v := range c.variants {
		res, err := v.Equal(object)
		if err != nil || res == BTrue {
			return res, err
		} else if res == BUnknown {
			result = BUnknown
		}
	}
	return result, nil
}

func (c NumberOr) NotEqual(object interface{}) (BValue, error) {
	result := BFalse
	for _, v := range c.variants {
		res, err := v.NotEqual(object)
		if err != nil || res == BTrue {
			return res, err
		} else if res == BUnknown {
			result = BUnknown
		}
	}
	return result, nil
}

func (c NumberOr) Less(object interface{}) (BValue, error) {
	result := BFalse
	for _, v := range c.variants {
		res, err := v.Less(object)
		if err != nil || res == BTrue {
			return res, err
		} else if res == BUnknown {
			result = BUnknown
		}
	}
	return result, nil
}

func (c NumberOr) Greater(object interface{}) (BValue, error) {
	result := BFalse
	for _, v := range c.variants {
		res, err := v.Greater(object)
		if err != nil || res == BTrue {
			return res, err
		} else if res == BUnknown {
			result = BUnknown
		}
	}
	return result, nil
}

func (c NumberOr) LessEqual(object interface{}) (BValue, error) {
	lt, err := c.Less(object)
	if err != nil || lt == BTrue {
		return lt, err
	}

	eq, err := c.Equal(object)
	if err != nil || eq == BTrue {
		return eq, err
	}

	return or(lt, eq), nil
}

func (c NumberOr) GreaterEqual(object interface{}) (BValue, error) {
	gt, err := c.Greater(object)
	if err != nil || gt == BTrue {
		return gt, err
	}

	eq, err := c.Equal(object)
	if err != nil || eq == BTrue {
		return eq, err
	}

	return or(gt, eq), nil
}

func (c NumberOr) Inverse(object interface{}) (NumberConstraint, error) {
	panic("NYI")
}

// ====== NumberEqual ======

type NumberEqual struct {
	subject *NumberPrivate
}

func NewNumberEqual(subject Number) NumberEqual {
	return NumberEqual{subject: subject.p}
}

func (_ NumberEqual) Name() string {
	return "NumberEqual"
}

func (c NumberEqual) Equal(object interface{}) (BValue, error) {
	if obj, ok := object.(*NumberPrivate); ok {
		if c.subject == obj {
			return BTrue, nil
		}
		return BUnknown, nil
	}
	return -1, errApplyInvalidNumber("NumberEqual")
}

func (c NumberEqual) NotEqual(object interface{}) (BValue, error)     { return not(c.Equal(object)) }
func (c NumberEqual) Less(object interface{}) (BValue, error)         { return not(c.Equal(object)) }
func (c NumberEqual) Greater(object interface{}) (BValue, error)      { return not(c.Equal(object)) }
func (c NumberEqual) LessEqual(object interface{}) (BValue, error)    { return c.Equal(object) }
func (c NumberEqual) GreaterEqual(object interface{}) (BValue, error) { return c.Equal(object) }

func (c NumberEqual) Inverse(subject interface{}) (NumberConstraint, error) {
	panic("NYI")
}

/*func (c NumberEqual) Unbox() []interface{} {
	panic("NYI")
}*/

// ====== NumberNotEqual ======

type NumberNotEqual struct {
	subject *NumberPrivate
}

func NewNumberNotEqual(subject Number) NumberNotEqual {
	return NumberNotEqual{subject: subject.p}
}

func (_ NumberNotEqual) Name() string {
	return "NumberNotEqual"
}

func (c NumberNotEqual) Equal(object interface{}) (BValue, error) { return not(c.NotEqual(object)) }

func (c NumberNotEqual) NotEqual(object interface{}) (BValue, error) {
	if obj, ok := object.(*NumberPrivate); ok {
		if c.subject == obj {
			return BTrue, nil
		}
		return BUnknown, nil
	}
	return -1, errApplyInvalidNumber("NumberNotEqual")
}

func (_ NumberNotEqual) Less(object interface{}) (BValue, error) {
	return unknown(object, "NumberNotEqual")
}

func (_ NumberNotEqual) Greater(object interface{}) (BValue, error) {
	return unknown(object, "NumberNotEqual")
}

// ====== NumberLess ======

type NumberLess struct {
	subject *NumberPrivate
}

func NewNumberLess(subject Number) NumberLess {
	return NumberLess{subject: subject.p}
}

func (_ NumberLess) Name() string {
	return "NumberLess"
}

func (c NumberLess) Equal(object interface{}) (BValue, error)    { return not(c.Less(object)) }
func (c NumberLess) NotEqual(object interface{}) (BValue, error) { return c.Less(object) }

func (c NumberLess) Less(object interface{}) (BValue, error) {
	if obj, ok := object.(*NumberPrivate); ok {
		if c.subject == obj {
			return BTrue, nil
		}
		return BUnknown, nil
	}
	return -1, errApplyInvalidNumber("NumberLess")
}

func (c NumberLess) Greater(object interface{}) (BValue, error)      { return not(c.Less(object)) }
func (c NumberLess) LessEqual(object interface{}) (BValue, error)    { return c.Less(object) }
func (c NumberLess) GreaterEqual(object interface{}) (BValue, error) { return not(c.Less(object)) }

func (c NumberLess) Inverse(subject interface{}) (NumberConstraint, error) {
	panic("NYI")
}

// ====== NumberLessEqual ======

type NumberLessEqual struct {
	subject *NumberPrivate
}

func NewNumberLessEqual(subject Number) NumberLessEqual {
	return NumberLessEqual{subject: subject.p}
}

func (_ NumberLessEqual) Name() string {
	return "NumberLessEqual"
}

func (_ NumberLessEqual) Equal(object interface{}) (BValue, error) {
	return unknown(object, "NumberLessEqual")
}

func (_ NumberLessEqual) NotEqual(object interface{}) (BValue, error) {
	return unknown(object, "NumberLessEqual")
}

func (_ NumberLessEqual) Less(object interface{}) (BValue, error) {
	return unknown(object, "NumberLessEqual")
}

func (c NumberLessEqual) Greater(object interface{}) (BValue, error) {
	return not(c.LessEqual(object))
}

func (c NumberLessEqual) LessEqual(object interface{}) (BValue, error) {
	if obj, ok := object.(*NumberPrivate); ok {
		if c.subject == obj {
			return BTrue, nil
		}
		return BUnknown, nil
	}
	return -1, errApplyInvalidNumber("NumberLessEqual")
}

func (_ NumberLessEqual) GreaterEqual(object interface{}) (BValue, error) {
	return unknown(object, "NumberLessEqual")
}

func (c NumberLessEqual) Inverse(subject interface{}) (NumberConstraint, error) {
	panic("NYI")
}

// ====== NumberGreater ======

type NumberGreater struct {
	subject *NumberPrivate
}

func NewNumberGreater(subject Number) NumberGreater {
	return NumberGreater{subject: subject.p}
}

func (_ NumberGreater) Name() string {
	return "NumberGreater"
}

func (c NumberGreater) Equal(object interface{}) (BValue, error)    { return not(c.Greater(object)) }
func (c NumberGreater) NotEqual(object interface{}) (BValue, error) { return c.Greater(object) }
func (c NumberGreater) Less(object interface{}) (BValue, error)     { return not(c.Greater(object)) }

func (c NumberGreater) Greater(object interface{}) (BValue, error) {
	if obj, ok := object.(*NumberPrivate); ok {
		if c.subject == obj {
			return BTrue, nil
		}
		return BUnknown, nil
	}
	return -1, errApplyInvalidNumber("NumberGreater")
}

func (c NumberGreater) LessEqual(object interface{}) (BValue, error)    { return not(c.Greater(object)) }
func (c NumberGreater) GreaterEqual(object interface{}) (BValue, error) { return c.Greater(object) }

func (c NumberGreater) Inverse(subject interface{}) (NumberConstraint, error) {
	panic("NYI")
}

/*func (c NumberGreater) Unbox() []interface{} {
	panic("NYI")
	return nil
}*/

// ====== NumberGreaterEqual ======

type NumberGreaterEqual struct {
	subject *NumberPrivate
}

func NewNumberGreaterEqual(subject Number) NumberGreaterEqual {
	return NumberGreaterEqual{subject: subject.p}
}

func (_ NumberGreaterEqual) Name() string {
	return "NumberGreaterEqual"
}

func (_ NumberGreaterEqual) Equal(object interface{}) (BValue, error) {
	return unknown(object, "NumberGreaterEqual")
}

func (_ NumberGreaterEqual) NotEqual(object interface{}) (BValue, error) {
	return unknown(object, "NumberGreaterEqual")
}

func (c NumberGreaterEqual) Less(object interface{}) (BValue, error) {
	return not(c.GreaterEqual(object))
}

func (_ NumberGreaterEqual) Greater(object interface{}) (BValue, error) {
	return unknown(object, "NumberGreaterEqual")
}

func (_ NumberGreaterEqual) LessEqual(object interface{}) (BValue, error) {
	return unknown(object, "NumberGreaterEqual")
}

func (c NumberGreaterEqual) GreaterEqual(object interface{}) (BValue, error) {
	if obj, ok := object.(*NumberPrivate); ok {
		if c.subject == obj {
			return BTrue, nil
		}
		return BUnknown, nil
	}
	return -1, errApplyInvalidNumber("NumberGreaterEqual")
}

func (c NumberGreaterEqual) Inverse(subject interface{}) (NumberConstraint, error) {
	panic("NYI")
}
