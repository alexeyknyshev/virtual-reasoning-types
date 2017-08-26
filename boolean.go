package virtual_types

type BValue int

const (
	BFalse BValue = iota
	BTrue
	BUnknown
	BInvalid
)

type BooleanPrivate struct {
	val         BValue
	constraints []Constraint
}

func (p *BooleanPrivate) equal(o *BooleanPrivate) *BooleanPrivate {
	if p == o {
		return &BooleanPrivate{val: BTrue}
	}

	if (p.val == BFalse || p.val == BTrue) && (o.val == BFalse || o.val == BTrue) {
		if p.val == o.val {
			return &BooleanPrivate{val: BTrue}
		}
		return &BooleanPrivate{val: BFalse}
	}

	res := BUnknown
	if p.val == BUnknown {
		var err error
		res, err = checkConstraints(p, o, constraintEqualAllVisitor)
		if err != nil {
			panic(err.Error())
		}
	}

	if res == BFalse {
		return &BooleanPrivate{val: BFalse}
	}

	if o.val == BUnknown {
		new_res, err := checkConstraints(o, p, constraintEqualAllVisitor)
		if err != nil {
			panic(err.Error())
		}

		if new_res == BTrue || new_res == BFalse {
			res = new_res
		}
	}

	return &BooleanPrivate{val: res}
}

func (p *BooleanPrivate) not() *BooleanPrivate {
	res := BUnknown

	switch p.val {
	case BFalse:
		res = BTrue
		break
	case BTrue:
		res = BFalse
		break
	case BUnknown:
		break
	default:
		panic("Invalid BValue")
	}

	return &BooleanPrivate{
		val:         res,
		constraints: []Constraint{BooleanNotEqual{subject: p}},
	}
}

type constraintFunc func(c Constraint, o *BooleanPrivate) (BValue, error)

func checkConstraints(p, o *BooleanPrivate, f constraintFunc) (BValue, error) {
	res := BUnknown
	for _, c := range p.constraints {
		r, err := f(c, o)
		if err != nil || r == BFalse {
			return r, err
		}
		if r == BTrue {
			res = BTrue
		}
	}
	return res, nil
}

type Boolean struct {
	p *BooleanPrivate
}

func (b Boolean) TypeName() string {
	return "Boolean"
}

func (b Boolean) IsValid() bool {
	return b.p != nil
}

func (b Boolean) IsUnknown() bool {
	return b.p.val == BUnknown
}

func (b Boolean) IsConstant() bool {
	return b.IsTrue() || b.IsFalse()
}

func (b Boolean) IsTrue() bool {
	return b.p.val == BTrue
}

func (b Boolean) IsFalse() bool {
	return b.p.val == BFalse
}

func (b Boolean) IsSame(o Boolean) bool {
	return b.p.val == o.p.val
}

func constraintEqualAllVisitor(c Constraint, o *BooleanPrivate) (BValue, error) {
	return c.Equal(o)
}

func (b Boolean) Equal(o Boolean) Boolean {
	return Boolean{p: b.p.equal(o.p)}
}

func (b Boolean) Not() Boolean {
	return Boolean{p: b.p.not()}
}

func (b Boolean) And(o Boolean) Boolean {
	if b.p.val == BFalse {
		return b
	} else if b.p.val == BUnknown {
		res, err := checkConstraints(b.p, o.p, constraintEqualAllVisitor)
		if err != nil {
			panic(err.Error())
		}
		if res == BFalse {
			return NewBooleanConst(BFalse, []Constraint{
				NewConstraintOr([]Constraint{
					NewBooleanEqual(b), NewBooleanEqual(o),
				}),
			})
		}
	}

	if o.p.val == BFalse {
		if b.p.val == BUnknown {
			return NewBooleanConst(BFalse, []Constraint{
				NewConstraintOr([]Constraint{
					NewBooleanEqual(b), NewBooleanEqual(o),
				}),
			})
		}
	} else if o.p.val == BUnknown {
		res, err := checkConstraints(o.p, b.p, constraintEqualAllVisitor)
		if err != nil {
			panic(err.Error())
		}
		if res == BFalse {
			if b.p.val == BUnknown {
				return NewBooleanConst(BFalse, []Constraint{
					NewConstraintOr([]Constraint{
						NewBooleanEqual(b), NewBooleanEqual(o),
					}),
				})
			}
			return NewBooleanConst(BFalse, []Constraint{
				NewBooleanEqual(o),
			})
		}
	}

	if b.p.val == BUnknown {
		return NewBooleanConst(BUnknown, []Constraint{
			NewConstraintOr([]Constraint{
				NewBooleanEqual(b), NewBooleanEqual(o),
			}),
		})
	}

	return o
}

func (b Boolean) Or(o Boolean) Boolean {
	if b.p.val == BTrue {
		return b
	} else if b.p.val == BFalse {
		return o
	}

	res := BUnknown
	if o.p.val == BTrue {
		res = BTrue
	}

	return NewBooleanConst(res, []Constraint{
		NewConstraintOr([]Constraint{
			NewBooleanEqual(b), NewBooleanEqual(o),
		}),
	})
}

func NewBoolean() Boolean {
	return NewBooleanConst(BUnknown, nil)
}

func NewBooleanConst(v BValue, constraints []Constraint) Boolean {
	switch v {
	case BFalse:
	case BTrue:
	case BUnknown:
		break
	default:
		panic("Invalid BValue")
	}

	return Boolean{
		p: &BooleanPrivate{
			val:         v,
			constraints: constraints,
		},
	}
}
