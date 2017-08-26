package virtual_types

import "fmt"

func errApplyInvalidBoolean(action string) error {
	return fmt.Errorf("could not apply %s constraint to non-Boolean value", action)
}

func errInverseInvalidBoolean(action string) error {
	return fmt.Errorf("could not inserse %s constraint on non-Boolean value", action)
}

// ====== BooleanEqual ======

type BooleanEqual struct {
	subject *BooleanPrivate
}

func NewBooleanEqual(subject Boolean) BooleanEqual {
	return BooleanEqual{subject: subject.p}
}

func (c BooleanEqual) Name() string {
	return "BooleanEqual"
}

func (c BooleanEqual) Equal(object interface{}) (BValue, error) {
	if obj, ok := object.(*BooleanPrivate); ok {
		if c.subject == obj {
			return BTrue, nil
		} else if c.subject.val == BUnknown || obj.val == BUnknown {
			p := c.subject.equal(obj)
			return p.val, nil
		} else if c.subject.val == obj.val {
			return BTrue, nil
		}
		return BFalse, nil
	}
	return BInvalid, errApplyInvalidBoolean("BooleanEqual")
}

func (c BooleanEqual) NotEqual(object interface{}) (BValue, error) {
	if obj, ok := object.(*BooleanPrivate); ok {
		if c.subject == obj {
			return BFalse, nil
		} else if c.subject.val == BUnknown || obj.val == BUnknown {
			p := c.subject.equal(obj).not()
			return p.val, nil
		} else if c.subject.val == obj.val {
			return BFalse, nil
		}
		return BTrue, nil
	}
	return BInvalid, errApplyInvalidBoolean("BooleanEqual")
}

func (c BooleanEqual) Inverse(subject interface{}) (Constraint, error) {
	if subj, ok := subject.(*BooleanPrivate); ok {
		return BooleanEqual{subject: subj}, nil
	}
	return nil, errInverseInvalidBoolean("BooleanEqual")
}

func (c BooleanEqual) Unbox() []interface{} {
	return []interface{}{c.subject}
}

// ====== BooleanNotEqual ======

type BooleanNotEqual struct {
	subject *BooleanPrivate
}

func NewBooleanNotEqual(subject Boolean) BooleanNotEqual {
	return BooleanNotEqual{subject: subject.p}
}

func (c BooleanNotEqual) Name() string {
	return "BooleanNotEqual"
}

func (c BooleanNotEqual) Equal(object interface{}) (BValue, error) {
	if obj, ok := object.(*BooleanPrivate); ok {
		if c.subject == obj {
			return BFalse, nil
		} else if c.subject.val == BUnknown || obj.val == BUnknown {
			p := c.subject.equal(obj).not()
			return p.val, nil
		} else if c.subject.val == obj.val {
			return BFalse, nil
		}
		return BTrue, nil
	}
	return BInvalid, errApplyInvalidBoolean("BooleanNotEqual")
}

func (c BooleanNotEqual) NotEqual(object interface{}) (BValue, error) {
	if obj, ok := object.(*BooleanPrivate); ok {
		if c.subject == obj {
			return BTrue, nil
		} else if c.subject.val == BUnknown || obj.val == BUnknown {
			p := c.subject.equal(obj)
			return p.val, nil
		} else if c.subject.val == obj.val {
			return BTrue, nil
		}
		return BFalse, nil
	}
	return BInvalid, errApplyInvalidBoolean("BooleanNotEqual")
}

func (c BooleanNotEqual) Inverse(subject interface{}) (Constraint, error) {
	if subj, ok := subject.(*BooleanPrivate); ok {
		return BooleanNotEqual{subject: subj}, nil
	}
	return nil, errInverseInvalidBoolean("BooleanNotEqual")
}

func (c BooleanNotEqual) Unbox() []interface{} {
	return nil
}

// ====== BooleanDummyConstraint ======

type BooleanDummyConstraint struct{}

func (c BooleanDummyConstraint) Name() string {
	return "BooleanDummyConstraint"
}

func (c BooleanDummyConstraint) Equal(object interface{}) (BValue, error) {
	if _, ok := object.(*BooleanPrivate); ok {
		return BUnknown, nil
	}
	return BInvalid, errApplyInvalidBoolean("BooleanDummyConstraint")
}

func (c BooleanDummyConstraint) NotEqual(object interface{}) (BValue, error) {
	if _, ok := object.(*BooleanPrivate); ok {
		return BUnknown, nil
	}
	return BInvalid, errApplyInvalidBoolean("BooleanDummyConstrain")
}

func (c BooleanDummyConstraint) Inverse(subject interface{}) (Constraint, error) {
	if _, ok := subject.(*BooleanPrivate); ok {
		return c, nil
	}
	return nil, errInverseInvalidBoolean("BooleanDummyConstraint")
}

func (c BooleanDummyConstraint) Unbox() []interface{} {
	return nil
}
