package virtual_types

import (
	"errors"
	"fmt"
	"math"
)

const (
	DIV_BY_ZERO_STR = "division by zero"
	NAN_STR         = "NaN"
)

var (
	ERR_DIV_BY_ZERO = errors.New(DIV_BY_ZERO_STR)
	ERR_NAN         = errors.New(NAN_STR)
)

var _zero Number = NewNumberConst(0)
var _one Number = NewNumberConst(1)

var _inf Number = NewNumberConst(math.Inf(1))
var _minus_inf Number = NewNumberConst(math.Inf(-1))

type NEdge int

const (
	NEdgeNo    NEdge = 0
	NEdgeLeft  NEdge = -1
	NEdgeRight NEdge = 1
)

type NRange struct {
	lVal float64
	rVal float64

	lIncluding bool
	rIncluding bool
}

func minmax(s []float64) (min float64, max float64) {
	sLen := len(s)
	if sLen == 0 {
		return 0, 0
	}

	min, max = s[0], s[0]
	for i := 2; i < sLen; i++ {
		if max < s[i] {
			max = s[i]
		} else if min > s[i] {
			min = s[i]
		}
	}
	return min, max
}

func newRangeSegment(l, r float64) *NRange {
	return &NRange{
		lVal: l,
		rVal: r,

		lIncluding: true,
		rIncluding: true,
	}
}

func newRange() *NRange {
	return newRangeSegment(math.Inf(-1), math.Inf(1))
}

func floor(v float64, arithmeticallyCorrect bool) float64 {
	if arithmeticallyCorrect && v < 0 {
		return -math.Floor(-v)
	}
	return math.Floor(v)
}

func (r *NRange) IsSame(o *NRange) bool {
	return r.lVal == o.lVal && r.rVal == o.rVal &&
		r.lIncluding == o.lIncluding && r.rIncluding == o.rIncluding
}

func (r *NRange) Clone() *NRange {
	return &NRange{
		lVal: r.lVal,
		rVal: r.rVal,

		lIncluding: r.lIncluding,
		rIncluding: r.rIncluding,
	}
}

func (r *NRange) Invert() *NRange {
	return &NRange{
		lVal: r.rVal,
		rVal: r.lVal,

		lIncluding: r.rIncluding,
		rIncluding: r.lIncluding,
	}
}

func (r *NRange) Negate() *NRange {
	return &NRange{
		lVal: -r.rVal,
		rVal: -r.lVal,

		lIncluding: r.rIncluding,
		rIncluding: r.lIncluding,
	}
}

func (r *NRange) Split(splitPoint float64) (*NRange, *NRange) {
	if splitPoint <= r.lVal {
		return nil, r
	}
	if splitPoint > r.rVal {
		return r, nil
	}
	if splitPoint == r.rVal {
		if !r.rIncluding {
			return r, nil
		}

		clone := r.Clone()
		clone.rIncluding = false
		return clone, newRangeSegment(splitPoint, splitPoint)
	}

	left := &NRange{
		lVal: r.lVal,
		rVal: splitPoint,

		lIncluding: r.lIncluding,
		rIncluding: false,
	}
	right := &NRange{
		lVal: splitPoint,
		rVal: r.rVal,

		lIncluding: true,
		rIncluding: r.rIncluding,
	}
	return left, right
}

func (r *NRange) Merge(o *NRange) *NRange {
	r, o = r.Order(o)
	if r.rVal < o.lVal {
		return nil
	}
	if r.rVal == o.lVal && !r.rIncluding && !o.lIncluding {
		return nil
	}

	rVal := o.rVal
	rIncluding := o.rIncluding

	if rVal < r.rVal {
		rVal = r.rVal
		rIncluding = r.rIncluding
	} else if rVal == r.rVal {
		rIncluding = rIncluding || r.rIncluding
	}

	return &NRange{
		lVal: r.lVal,
		rVal: rVal,

		lIncluding: r.lIncluding,
		rIncluding: rIncluding,
	}
}

func (r *NRange) Abs() *NRange {
	if r.lVal <= 0 && r.rVal <= 0 {
		return r.Negate()
	} else if r.lVal < 0 {
		if -r.lVal > r.rVal {
			return &NRange{
				lVal: 0,
				rVal: -r.lVal,

				lIncluding: true,
				rIncluding: r.lIncluding,
			}
		} else if -r.lVal == r.rVal {
			return &NRange{
				lVal: 0,
				rVal: r.rVal,

				lIncluding: true,
				rIncluding: r.lIncluding || r.rIncluding,
			}
		}
	}
	return r
}

func (r *NRange) Floor(arithmeticallyCorrect, inverted bool) *NRange {
	lVal := floor(r.lVal, arithmeticallyCorrect)
	if inverted {
		if lVal < r.lVal || !r.lIncluding {
			lVal++
		}
	}

	rVal := floor(r.rVal, arithmeticallyCorrect)
	if rVal == r.rVal && !r.rIncluding {
		rVal--
	}
	if inverted && rVal < r.rVal {
		rVal++
	}

	if lVal != r.lVal || rVal != r.rVal || !r.lIncluding || !r.rIncluding {
		return &NRange{
			lVal: lVal,
			rVal: rVal,

			lIncluding: true,
			rIncluding: true,
		}
	}

	return r
}

func (r *NRange) Contains(val float64) (bool, NEdge) {
	if math.IsNaN(val) {
		return false, NEdgeNo
	}

	on_edge := NEdgeNo

	if val < r.lVal {
		return false, NEdgeNo
	} else if val == r.lVal {
		if r.lIncluding {
			on_edge = NEdgeLeft
		} else {
			return false, NEdgeLeft
		}
	}

	if val > r.rVal {
		return false, NEdgeNo
	} else if val == r.rVal {
		if r.rIncluding {
			on_edge = NEdgeRight
		} else {
			return false, NEdgeRight
		}
	}

	return true, on_edge
}

func (r *NRange) IsInf(sign int) Boolean {
	if sign < 0 {
		if math.IsInf(r.lVal, sign) {
			return NewBoolean()
		}
		return NewBooleanConst(BFalse, nil)
	} else if sign > 0 {
		if math.IsInf(r.rVal, sign) {
			return NewBoolean()
		}
		return NewBooleanConst(BFalse, nil)
	}

	if math.IsInf(r.lVal, -1) || math.IsInf(r.rVal, 1) {
		return NewBoolean()
	}
	return NewBooleanConst(BFalse, nil)
}

func (r *NRange) ToIntegerRange() *NRange {
	changed := false

	lVal := math.Floor(r.lVal)
	if !math.IsInf(lVal, 0) {
		if lVal < r.lVal || !r.lIncluding {
			changed = true
			lVal++
		}
	}

	rVal := math.Floor(r.rVal)
	if !math.IsInf(rVal, 0) {
		if rVal < r.rVal {
			changed = true
		} else if !r.rIncluding {
			changed = true
			rVal--
		}
	}

	if lVal > rVal {
		return nil
	}

	if !changed {
		return r
	}

	return &NRange{
		lVal:       lVal,
		rVal:       rVal,
		lIncluding: true,
		rIncluding: true,
	}
}

func (r *NRange) Extend(o *NRange, leftExtend bool) *NRange {
	if leftExtend {
		if r.lVal > o.lVal {
			return &NRange{
				lVal: o.lVal,
				rVal: r.rVal,

				lIncluding: o.lIncluding,
				rIncluding: r.rIncluding,
			}
		} else if r.lVal == o.lVal && !r.lIncluding && o.lIncluding {
			return &NRange{
				lVal: r.lVal,
				rVal: r.rVal,

				lIncluding: o.lIncluding,
				rIncluding: r.rIncluding,
			}
		}
	} else {
		if r.rVal < o.rVal {
			return &NRange{
				lVal: r.lVal,
				rVal: o.rVal,

				lIncluding: r.lIncluding,
				rIncluding: o.rIncluding,
			}
		} else if r.rVal == o.rVal && !r.rIncluding && o.rIncluding {
			return &NRange{
				lVal: r.lVal,
				rVal: r.rVal,

				lIncluding: r.lIncluding,
				rIncluding: o.rIncluding,
			}
		}
	}
	return r
}

func edge_pick(r *NRange, left bool) (float64, bool) {
	if left {
		return r.lVal, r.lIncluding
	}
	return r.rVal, r.rIncluding
}

func edge_cmp(a, b *NRange, aLeft, bLeft bool) (int, bool) {
	aVal, aIncluding := edge_pick(a, aLeft)
	bVal, bIncluding := edge_pick(b, bLeft)

	if aVal < bVal {
		return -1, false
	} else if aVal > bVal {
		return 1, false
	}

	if aIncluding == bIncluding {
		if aIncluding || aLeft == bLeft {
			return 0, true
		}
		if bLeft {
			return -1, true
		}
		return 1, true
	}

	if aLeft == bLeft {
		if aLeft {
			if aIncluding {
				return -1, true
			}
			return 1, true
		} else {
			if aIncluding {
				return 1, true
			}
			return -1, true
		}
	} else if aLeft {
		return 1, true
	}

	return -1, true
}

func (r *NRange) Order(o *NRange) (*NRange, *NRange) {
	aLeft, bLeft := true, true

	leftEdgeCmp, _ := edge_cmp(r, o, aLeft, bLeft)
	if leftEdgeCmp == 1 {
		return o, r
	}
	return r, o
}

func (r *NRange) Less(o *NRange) Boolean {
	if r.rVal < o.lVal {
		return NewBooleanConst(BTrue, nil)
	} else if r.rVal == o.lVal {
		if r.rIncluding && o.lIncluding {
			return NewBoolean()
		}
		return NewBooleanConst(BTrue, nil)
	} else if o.rVal <= r.lVal {
		return NewBooleanConst(BFalse, nil)
	}

	return NewBoolean()
}

func (r *NRange) ArithmeticOperation(o *NRange, op ArithmeticOperationBinary) ([]*NRange, error) {
	if !op.IsClosedField() {
		lSplit, rSplit := o.Split(0)
		if lSplit != nil && rSplit != nil {
			lRes, err := r.ArithmeticOperation(lSplit, op)
			if err != nil {
				return nil, err
			}
			fmt.Printf("%+v / %+v = %+v\n", r, lSplit, lRes[0])
			fmt.Printf("left done\n")
			rRes, err := r.ArithmeticOperation(rSplit, op)
			if err != nil {
				return nil, err
			}

			if len(lRes) != 1 || len(rRes) != 1 {
				panic("got unexpected ranges cnt from subsequent ArithmeticOperation")
			}

			rUnion := lRes[0].Merge(rRes[0])
			if rUnion == nil {
				return append(lRes, rRes[0]), nil
			}
			return []*NRange{rUnion}, nil
		}
	}

	r = op.PreprocessRangeLeft(r)
	o = op.PreprocessRangeRight(o)

	changed := false

	lIncludingForce := false
	rIncludingForce := false

	lVal, err := op.Compute(r.lVal, o.lVal)
	if err != nil {
		switch err.Error() {
		case DIV_BY_ZERO_STR:
			if o.rVal == 0 { // range is 0 constant
				return nil, err
			}

			lVal = math.Copysign(math.Inf(0), r.lVal*o.rVal)
			//fmt.Printf("lVal(fix) = %f (r.lVal = %f, o.rVal = %f)\n", lVal, r.lVal, o.rVal)
			lIncludingForce = true
			changed = true
			break
		case NAN_STR:
			break
		default:
			return nil, err
		}
	} else if lVal != r.lVal {
		changed = true
	}

	//fmt.Printf("o.lVal = %f\n", o.lVal)

	rVal, err := op.Compute(r.rVal, o.rVal)
	//fmt.Printf("rVal = %f = %f / %f\n", rVal, r.rVal, o.rVal)
	if err != nil {
		switch err.Error() {
		case DIV_BY_ZERO_STR:
			rVal = math.Copysign(math.Inf(0), r.rVal*o.lVal)
			//fmt.Printf("rVal(fix) = %f (r.rVal = %f, o.lVal = %f)\n", rVal, r.rVal, o.lVal)
			rIncludingForce = true
			changed = true
			break
		case NAN_STR:
			break
		default:
			return nil, err
		}
	} else if rVal != r.rVal {
		changed = true
	}

	//fmt.Printf("l: %f / %f = %f\n", r.lVal, o.lVal, lVal)
	//fmt.Printf("r: %f / %f = %f\n", r.rVal, o.rVal, rVal)

	lIncluding := r.lIncluding && o.lIncluding || lIncludingForce
	rIncluding := r.rIncluding && o.rIncluding || rIncludingForce

	if lIncluding != r.lIncluding || rIncluding != r.rIncluding {
		changed = true
	}

	if !changed {
		return []*NRange{r}, nil
	}

	return []*NRange{&NRange{
		lVal: lVal,
		rVal: rVal,

		lIncluding: lIncluding,
		rIncluding: rIncluding,
	}}, nil
}

func (r *NRange) IsConstant() bool {
	return r.lVal == r.rVal && r.lIncluding && r.lIncluding
}

func (r *NRange) IsNaN() bool {
	if r.lVal > r.rVal || (r.lVal == r.rVal && !(r.lIncluding && r.rIncluding)) {
		return true
	}
	return math.IsNaN(r.lVal) || math.IsNaN(r.rVal)
}

func (r *NRange) Sign() []float64 {
	if r.IsNaN() {
		return nil
	}
	if r.lVal > 0 {
		return []float64{1}
	} else if r.lVal == 0 {
		if r.lIncluding {
			return []float64{0, 1}
		}
		return []float64{1}
	} else if r.rVal < 0 {
		return []float64{-1}
	} else if r.rVal == 0 {
		if r.rIncluding {
			return []float64{-1, 0}
		}
		return []float64{-1}
	}
	return []float64{-1, 0, 1}
}

func (r *NRange) Overlaps(o *NRange) (*NRange, NEdge) {
	if r == o {
		return r, NEdgeNo
	}

	lCmp, _ := edge_cmp(r, o, true, true)
	rCmp, _ := edge_cmp(r, o, false, false)

	swapped := false
	if lCmp == 0 {
		if rCmp <= 0 {
			return r, NEdgeNo
		}
	} else if lCmp < 0 {
		if rCmp >= 0 {
			return o, NEdgeNo
		}
	} else {
		if rCmp <= 0 {
			return r, NEdgeNo
		}
		r, o = o, r
		swapped = true
	}

	rRight_oLeftCmp, _ := edge_cmp(r, o, false, true)
	if rRight_oLeftCmp == -1 {
		return nil, NEdgeNo
	}

	edge := NEdgeNo
	if rRight_oLeftCmp == 0 {
		if swapped {
			edge = NEdgeLeft
		} else {
			edge = NEdgeRight
		}
	}

	return &NRange{
		lVal: o.lVal,
		rVal: r.rVal,

		lIncluding: o.lIncluding,
		rIncluding: r.rIncluding,
	}, edge
}

type Number struct {
	p *NumberPrivate
}

type NumberPrivate struct {
	val         float64
	integer     Boolean
	valRange    *NRange
	next        *NumberPrivate
	constraints []NumberConstraint
}

func newNumberPrivate(r *NRange) *NumberPrivate {
	if r == nil {
		r = &NRange{
			lVal: math.Inf(-1),
			rVal: math.Inf(1),

			lIncluding: true,
			rIncluding: true,
		}
	}

	return &NumberPrivate{
		val:      0,
		integer:  NewBoolean(),
		valRange: r,
		next:     nil,
	}
}

func (p *NumberPrivate) Clone() *NumberPrivate {
	res := &NumberPrivate{
		val:     p.val,
		integer: p.integer,
	}

	if p.valRange != nil {
		res.valRange = p.valRange.Clone()
	}
	if p.next != nil {
		res.next = p.next.Clone()
	}
	return res
}

func (p *NumberPrivate) IsInteger() Boolean {
	if p.next == nil {
		return p.integer
	}
	return p.integer.And(p.next.IsInteger())
}

func (p *NumberPrivate) IsConstant() bool {
	if p.valRange != nil {
		if p.valRange.IsConstant() {
			p.val = p.valRange.lVal
			p.valRange = nil
			return true
		}
		return false
	}

	if p.next != nil {
		return p.next.IsConstant()
	}
	return true
}

func (p *NumberPrivate) Sign() []float64 {
	if p.IsConstant() {
		return []float64{math.Copysign(1, p.val)}
	}
	return p.valRange.Sign()
}

func (n Number) IsValid() bool {
	return n.p != nil
}

func (n Number) IsSame(o Number) bool {
	if n.p == o.p {
		return true
	}

	nConst := n.IsConstant()
	oConst := o.IsConstant()

	if nConst && oConst {
		return n.p.val == o.p.val || (math.IsNaN(n.p.val) && math.IsNaN(o.p.val))
	} else if !nConst && !oConst {
		if !n.p.valRange.IsSame(o.p.valRange) {
			return false
		}
		return n.p.integer.IsSame(o.p.integer)
	}

	return false
}

func newNumberConstWithIntegerHint(val float64, hint Boolean) Number {
	intVal := BFalse
	if math.IsNaN(val) {
		intVal = BUnknown
	} else if val == math.Floor(val) {
		intVal = BTrue
	}

	if intVal == hint.p.val {
		return Number{&NumberPrivate{
			val:     val,
			integer: hint,
		}}
	}
	return NewNumberConst(val)
}

func (n Number) RangeAdjust() (Number, error) {
	if n.p.next != nil {
		panic("next is unsupported yet")
	}

	r := n.p.valRange
	if r == nil {
		return n, nil
	}

	integer := n.IsInteger()
	//fmt.Printf("RangeAdjust: p = %+v (valRange = %+v)\n", n.p, r)

	if r.lVal > r.rVal {
		r = r.Invert()
	} else if r.lVal == r.rVal {
		return newNumberConstWithIntegerHint(r.lVal, integer), nil
	} else if math.IsNaN(r.lVal) || math.IsNaN(r.rVal) {
		return Number{}, errors.New("NaN edge un NRange")
	}

	if integer.IsTrue() {
		intRange := r.ToIntegerRange()
		if intRange == nil {
			return Number{}, errors.New("no integer representation for NRange")
		} else if intRange != r {
			if intRange.IsConstant() {
				return newNumberConstWithIntegerHint(intRange.lVal, integer), nil
			}
			r = intRange
		}
	} else if integer.IsFalse() {
		if math.Floor(r.lVal) == r.lVal && r.lIncluding {
			if r == n.p.valRange {
				r = r.Clone()
			}
			r.lIncluding = false
		}
		if math.Floor(r.rVal) == r.rVal && r.rIncluding {
			if r == n.p.valRange {
				r = r.Clone()
			}
			r.rIncluding = false
		}
	}

	if r == n.p.valRange {
		return n, nil
	}

	p := n.p.Clone()
	p.valRange = r

	return Number{p: p}, nil
}

func (n Number) IsConstant() bool {
	return n.p.IsConstant()
}

func (n Number) IsUnknown() bool {
	if n.IsConstant() {
		return false
	}

	r := n.p.valRange
	if math.IsInf(r.lVal, -1) && math.IsInf(r.rVal, 1) {
		return true
	}

	return false
}

func (n Number) IsInf(sign int) Boolean {
	if n.IsConstant() {
		if math.IsInf(n.p.val, sign) {
			return NewBooleanConst(BTrue, nil)
		}
		return NewBooleanConst(BFalse, nil)
	}

	return n.p.valRange.IsInf(sign)
}

func (n Number) IsUniversum() bool {
	if !n.IsInf(-1).IsUnknown() {
		return false
	}
	if !n.IsInf(1).IsUnknown() {
		return false
	}
	return true
}

func (n Number) IsNaN() Boolean {
	r := n.p.valRange
	if r == nil {
		if math.IsNaN(n.p.val) {
			return NewBooleanConst(BTrue, nil)
		}
		return NewBooleanConst(BFalse, nil)
	}

	if n.IsUnknown() {
		return NewBoolean()
	}

	if math.IsNaN(r.lVal) || math.IsNaN(r.rVal) {
		return NewBooleanConst(BTrue, nil)
	}

	return NewBooleanConst(BFalse, nil)
}

func (n Number) IsInteger() Boolean {
	return n.p.IsInteger()
}

func extendWithConst(n Number, c Number) []Number {
	if n.p.next != nil || c.p.next != nil {
		panic("next is unsupported yet")
	}

	if n.IsUnknown() {
		return []Number{NewNumber()}
	}

	nInt := n.IsInteger()
	cVal := c.p.val
	if contains, on_edge := n.p.valRange.Contains(cVal); !contains {
		if on_edge != NEdgeNo {
			cValInt := (cVal == math.Floor(cVal))
			if on_edge == NEdgeLeft {
				if cValInt && !nInt.IsFalse() {
					p := n.p.Clone()
					p.valRange.lIncluding = true
					return []Number{Number{p: p}}
				}
			} else { // NEdgeRight
				if cValInt && !nInt.IsFalse() {
					p := n.p.Clone()
					p.valRange.rIncluding = true
					return []Number{Number{p: p}}
				}
			}
		} else {
			return []Number{n, c}
		}
	}
	return []Number{n, c}
}

func (n Number) extendTo(o Number, leftExtend bool) []Number {
	if n.p.next != nil || o.p.next != nil {
		panic("next is unsupported yet")
	}

	nConst := n.IsConstant()
	oConst := o.IsConstant()

	if nConst && oConst {
		if n.p.val == o.p.val {
			return []Number{n}
		} else if n.p.val <= o.p.val {
			return []Number{n, o}
		}
		return []Number{o, n}
	} else if !nConst && !oConst {
		// TODO: check integer
		overlappingRange, _ := n.p.valRange.Overlaps(o.p.valRange)
		if overlappingRange == nil {
			return []Number{n, o}
		}

		r := n.p.valRange.Extend(o.p.valRange, leftExtend)
		if r == n.p.valRange {
			return []Number{n}
		}
		return []Number{NewNumberRange(r)}
	} else if nConst {
		return extendWithConst(o, n)
	} else { // oConst
		return extendWithConst(n, o)
	}
	return []Number{n, o}
}

func (p *NumberPrivate) less(o *NumberPrivate) (*BooleanPrivate, NEdge) {
	if p == o {
		return &BooleanPrivate{val: BFalse}, NEdgeNo
	}

	if p.next != nil || o.next != nil {
		panic("next is unsupported yet")
	}

	pConst := p.IsConstant()
	oConst := o.IsConstant()

	pRange := p.valRange
	oRange := o.valRange

	on_edge := NEdgeNo

	bVal := BUnknown
	if pConst && oConst {
		if p.val < o.val {
			bVal = BTrue
		} else {
			bVal = BFalse
		}
	} else {
		var err error
		bVal, err = checkConstraintsNumber(p, o, constraintNumberLessAllVisitor)
		if err != nil {
			panic(err.Error())
		}
		if isBValConst(bVal) {
			return &BooleanPrivate{val: bVal}, on_edge
		}

		bVal, err = checkConstraintsNumber(o, p, constraintNumberGreaterAllVisitor)
		if err != nil {
			panic(err.Error())
		}
		if isBValConst(bVal) {
			return &BooleanPrivate{val: bVal}, on_edge
		}

		if pConst {
			if p.val < oRange.lVal {
				bVal = BTrue
			} else if p.val == oRange.lVal {
				if !oRange.lIncluding {
					bVal = BTrue
				}
				on_edge = NEdgeLeft
			} else if p.val >= oRange.rVal {
				bVal = BFalse
			}
		} else if oConst {
			if pRange.rVal < o.val {
				bVal = BTrue
			} else if pRange.rVal == o.val {
				if !pRange.rIncluding {
					bVal = BTrue
				}
				on_edge = NEdgeRight
			} else if pRange.lVal >= o.val {
				bVal = BFalse
			}
		} else {
			less := pRange.Less(oRange)
			bVal = less.p.val
		}
	}

	return &BooleanPrivate{val: bVal}, on_edge
}

func (n Number) Less(o Number) Boolean {
	lt, _ := n.p.less(o.p)
	return Boolean{p: lt}
}

func (n Number) LessEqual(o Number) Boolean {
	if n.p == o.p {
		return NewBooleanConst(BTrue, nil)
	}

	bVal, err := checkConstraintsNumber(n.p, o.p, constraintNumberLessEqualAllVisitor)
	if err != nil {
		panic(err.Error())
	}
	if isBValConst(bVal) {
		return NewBooleanConst(bVal, nil)
	}

	bVal, err = checkConstraintsNumber(o.p, n.p, constraintNumberGreaterEqualAllVisitor)
	if err != nil {
		panic(err.Error())
	}
	if isBValConst(bVal) {
		return NewBooleanConst(bVal, nil)
	}

	lt, on_edge := n.p.less(o.p)
	if lt.val == BTrue {
		return Boolean{p: lt}
	}
	if on_edge == NEdgeLeft {
		return NewBooleanConst(BTrue, nil)
	}

	eq, _ := n.p.equal(o.p)
	return Boolean{p: lt}.Or(Boolean{p: eq})
}

func (n Number) Greater(o Number) Boolean {
	return o.Less(n)
}

func (n Number) GreaterEqual(o Number) Boolean {
	if n.p == o.p {
		return NewBooleanConst(BTrue, nil)
	}

	bVal, err := checkConstraintsNumber(n.p, o.p, constraintNumberGreaterEqualAllVisitor)
	if err != nil {
		panic(err.Error())
	}
	if isBValConst(bVal) {
		return NewBooleanConst(bVal, nil)
	}

	bVal, err = checkConstraintsNumber(o.p, n.p, constraintNumberLessEqualAllVisitor)
	if err != nil {
		panic(err.Error())
	}
	if isBValConst(bVal) {
		return NewBooleanConst(bVal, nil)
	}

	gt, on_edge := o.p.less(n.p)
	if gt.val == BTrue {
		return Boolean{p: gt}
	}
	if on_edge == NEdgeRight {
		return NewBooleanConst(BTrue, nil)
	}

	eq, _ := n.p.equal(o.p)
	return Boolean{p: gt}.Or(Boolean{p: eq})
}

func rangeOverlaps(l, r *NRange) (BValue, NEdge) {
	overlappingRange, edge := l.Overlaps(r)
	if overlappingRange == nil {
		return BFalse, NEdgeNo
	}
	return BUnknown, edge
}

func rangeContains(r *NRange, val float64) (BValue, NEdge) {
	contains, edge := r.Contains(val)
	if !contains {
		return BFalse, NEdgeNo
	}
	return BUnknown, edge
}

type constraintFuncNumber func(c NumberConstraint, o *NumberPrivate) (BValue, error)

func checkConstraintsNumber(p, o *NumberPrivate, f constraintFuncNumber) (BValue, error) {
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

func constraintNumberEqualAllVisitor(c NumberConstraint, o *NumberPrivate) (BValue, error) {
	return c.Equal(o)
}

func constraintNumberLessAllVisitor(c NumberConstraint, o *NumberPrivate) (BValue, error) {
	return c.Less(o)
}

func constraintNumberGreaterAllVisitor(c NumberConstraint, o *NumberPrivate) (BValue, error) {
	return c.Greater(o)
}

func constraintNumberLessEqualAllVisitor(c NumberConstraint, o *NumberPrivate) (BValue, error) {
	return c.LessEqual(o)
}

func constraintNumberGreaterEqualAllVisitor(c NumberConstraint, o *NumberPrivate) (BValue, error) {
	return c.GreaterEqual(o)
}

func isBValConst(val BValue) bool {
	return val == BTrue || val == BFalse
}

func (p *NumberPrivate) equal(o *NumberPrivate) (*BooleanPrivate, NEdge) {
	if p == o {
		if p.IsConstant() && math.IsNaN(p.val) {
			return &BooleanPrivate{val: BFalse}, NEdgeNo
		}
		return &BooleanPrivate{val: BTrue}, NEdgeNo
	}

	if p.next != nil || o.next != nil {
		panic("next is unsupported yet")
	}

	if p.integer.Equal(o.integer).IsFalse() {
		return &BooleanPrivate{val: BFalse}, NEdgeNo
	}

	bVal := BUnknown
	on_edge := NEdgeNo

	pConst := p.IsConstant()
	oConst := o.IsConstant()

	pRange := p.valRange
	oRange := o.valRange

	if pConst && oConst {
		if p.val == o.val {
			bVal = BTrue
		} else {
			bVal = BFalse
		}
	} else {
		var err error
		bVal, err = checkConstraintsNumber(p, o, constraintNumberEqualAllVisitor)
		if err != nil {
			panic(err.Error())
		}
		if isBValConst(bVal) {
			return &BooleanPrivate{val: bVal}, on_edge
		}
		bVal, err = checkConstraintsNumber(o, p, constraintNumberEqualAllVisitor)
		if err != nil {
			panic(err.Error())
		}
		if isBValConst(bVal) {
			return &BooleanPrivate{val: bVal}, on_edge
		}

		if !pConst && !oConst {
			bVal, on_edge = rangeOverlaps(pRange, oRange)
		} else if !pConst {
			bVal, on_edge = rangeContains(pRange, o.val)
		} else if !oConst {
			bVal, on_edge = rangeContains(oRange, p.val)
		}
	}

	return &BooleanPrivate{val: bVal}, on_edge
}

func (n Number) Equal(o Number) Boolean {
	eq, _ := n.p.equal(o.p)
	return Boolean{p: eq}
}

func (n Number) Negate() Number {
	res := n.p.Clone()

	if n.IsConstant() {
		res.val = -n.p.val
	} else if n.p.valRange != nil {
		res.valRange = n.p.valRange.Negate()
	}

	return Number{p: res}
}

type ArithmeticOperationBinary interface {
	Compute(x, y float64) (float64, error)
	IsClosedField() bool
	IsStrictClosedField() bool

	DetectEdgeCaseLeft(val float64, n Number) Number
	DetectEdgeCaseRight(n Number, val float64) Number

	DetectEdgeCaseSame(n Number) Number

	PreprocessRangeLeft(r *NRange) *NRange
	PreprocessRangeRight(r *NRange) *NRange

	IsResultInt() Boolean

	ResultConstraints(x, y, result Number) Number
}

type OpAdd struct{}

func (_ OpAdd) Compute(x, y float64) (float64, error) { return x + y, nil }
func (_ OpAdd) IsClosedField() bool                   { return true }
func (_ OpAdd) IsStrictClosedField() bool             { return true }

func (_ OpAdd) DetectEdgeCaseLeft(val float64, n Number) Number {
	if val == 0 {
		return n
	}
	if math.IsInf(val, 0) {
		if val > 0 {
			return _inf
		}
		return _minus_inf
	}
	return Number{}
}

func (o OpAdd) DetectEdgeCaseRight(n Number, val float64) Number {
	return o.DetectEdgeCaseLeft(val, n)
}

func (_ OpAdd) DetectEdgeCaseSame(n Number) Number     { return Number{} }
func (_ OpAdd) PreprocessRangeLeft(r *NRange) *NRange  { return r }
func (_ OpAdd) PreprocessRangeRight(r *NRange) *NRange { return r }
func (_ OpAdd) IsResultInt() Boolean                   { return Boolean{} }

func (_ OpAdd) ResultConstraints(x, y, result Number) Number {
	if result.IsConstant() || x.IsNaN().IsTrue() || y.IsNaN().IsTrue() {
		return result
	}

	var constraints []NumberConstraint

	if result.p != y.p {
		xPossibleSigns := x.p.Sign()
		xSignMin, xSignMax := minmax(xPossibleSigns)
		if xSignMin != 0 || xSignMax != 0 {
			if xSignMin > 0 {
				constraints = append(constraints, NewNumberGreater(y))
			} else if xSignMin == 0 {
				constraints = append(constraints, NewNumberGreaterEqual(y))
			}

			if xSignMax < 0 {
				constraints = append(constraints, NewNumberLess(y))
			} else if xSignMax == 0 {
				constraints = append(constraints, NewNumberLessEqual(y))
			}
		}
	}

	if result.p != x.p && y.p != x.p {
		yPossibleSigns := y.p.Sign()
		ySignMin, ySignMax := minmax(yPossibleSigns)

		if ySignMin != 0 || ySignMax != 0 {
			if ySignMin > 0 {
				constraints = append(constraints, NewNumberGreater(x))
			} else if ySignMin == 0 {
				constraints = append(constraints, NewNumberGreaterEqual(x))
			}

			if ySignMax < 0 {
				constraints = append(constraints, NewNumberLess(x))
			} else if ySignMax == 0 {
				constraints = append(constraints, NewNumberLessEqual(x))
			}
		}
	}

	result.p.constraints = append(result.p.constraints, constraints...)
	return result
}

type OpSub struct{}

func (_ OpSub) Compute(x, y float64) (float64, error) { return x - y, nil }
func (_ OpSub) IsClosedField() bool                   { return true }
func (_ OpSub) IsStrictClosedField() bool             { return true }

func (_ OpSub) DetectEdgeCaseLeft(val float64, n Number) Number {
	if val == 0 {
		return n.Negate()
	}
	if math.IsInf(val, 0) {
		if val > 0 {
			return _inf
		}
		return _minus_inf
	}
	return Number{}
}

func (_ OpSub) DetectEdgeCaseRight(n Number, val float64) Number {
	if val == 0 {
		return n
	}
	if math.IsInf(val, 0) {
		if val > 0 {
			return _minus_inf
		}
		return _inf
	}
	return Number{}
}

func (_ OpSub) DetectEdgeCaseSame(n Number) Number {
	if !n.IsNaN().IsTrue() {
		return _zero
	}
	return Number{}
}

func (_ OpSub) PreprocessRangeLeft(r *NRange) *NRange  { return r }
func (_ OpSub) PreprocessRangeRight(r *NRange) *NRange { return r.Invert() }
func (_ OpSub) IsResultInt() Boolean                   { return Boolean{} }

func (_ OpSub) ResultConstraints(x, y, result Number) Number {
	return result
}

type OpMul struct{}

func (_ OpMul) Compute(x, y float64) (float64, error) { return x * y, nil }
func (_ OpMul) IsClosedField() bool                   { return true }
func (_ OpMul) IsStrictClosedField() bool             { return false }

func (_ OpMul) DetectEdgeCaseLeft(val float64, n Number) Number {
	if val == 0 {
		return NewNumberConst(0)
	} else if val == 1 {
		return n
	}
	return Number{}
}

func (o OpMul) DetectEdgeCaseRight(n Number, val float64) Number {
	return o.DetectEdgeCaseLeft(val, n)
}

func (_ OpMul) DetectEdgeCaseSame(n Number) Number     { return Number{} }
func (_ OpMul) PreprocessRangeLeft(r *NRange) *NRange  { return r }
func (_ OpMul) PreprocessRangeRight(r *NRange) *NRange { return r }
func (_ OpMul) IsResultInt() Boolean                   { return Boolean{} }

func (_ OpMul) ResultConstraints(x, y, result Number) Number {
	return result
}

type OpDiv struct{}

func (_ OpDiv) Compute(x, y float64) (float64, error) {
	if y == 0 {
		return 0, ERR_DIV_BY_ZERO
	}
	res := x / y
	print(x, "/", y, "=", res, "\n")
	return res, nil
}

func (_ OpDiv) IsClosedField() bool       { return false }
func (_ OpDiv) IsStrictClosedField() bool { return false }

func (_ OpDiv) DetectEdgeCaseLeft(val float64, n Number) Number {
	if val == 0 {
		return _zero
	}
	return Number{}
}

func (_ OpDiv) DetectEdgeCaseRight(n Number, val float64) Number {
	if val == 1 {
		return n
	}
	return Number{}
}

func divDetectEdgeCaseSame(n Number) Number {
	if n.IsNaN().IsTrue() {
		return NewNumberConst(math.NaN())
	}
	if n.IsConstant() && n.p.val == 0 {
		return Number{}
	}
	return _one
}

func (_ OpDiv) DetectEdgeCaseSame(n Number) Number     { return divDetectEdgeCaseSame(n) }
func (_ OpDiv) PreprocessRangeLeft(r *NRange) *NRange  { return r }
func (_ OpDiv) PreprocessRangeRight(r *NRange) *NRange { return r.Invert() }
func (_ OpDiv) IsResultInt() Boolean                   { return Boolean{} }

func (_ OpDiv) ResultConstraints(x, y, result Number) Number {
	return result
}

type OpIDiv struct{}

func (_ OpIDiv) Compute(x, y float64) (float64, error) {
	if y == 0 {
		return 0, ERR_DIV_BY_ZERO
	}
	if math.IsInf(x, 0) && math.IsInf(y, 0) {
		return math.Copysign(x, math.Copysign(1, x)*math.Copysign(1, y)), ERR_NAN
	}
	return float64(int64(x) / int64(y)), nil
}

func (_ OpIDiv) IsClosedField() bool       { return false }
func (_ OpIDiv) IsStrictClosedField() bool { return false }

func (_ OpIDiv) DetectEdgeCaseLeft(val float64, n Number) Number {
	if val == 0 {
		return _zero
	}
	return Number{}
}

func (_ OpIDiv) DetectEdgeCaseRight(n Number, val float64) Number {
	if val == 1 {
		return n.Floor()
	}
	return Number{}
}

func (_ OpIDiv) DetectEdgeCaseSame(n Number) Number     { return divDetectEdgeCaseSame(n) }
func (_ OpIDiv) PreprocessRangeLeft(r *NRange) *NRange  { return r }
func (_ OpIDiv) PreprocessRangeRight(r *NRange) *NRange { return r.Invert() }

var numberIsInteger = NewBooleanConst(BTrue, nil)

func (_ OpIDiv) IsResultInt() Boolean { return numberIsInteger }

func (_ OpIDiv) ResultConstraints(x, y, result Number) Number {
	return result
}

type OpPow struct{}

func (_ OpPow) Compute(x, y float64) (float64, error) { return math.Pow(x, y), nil }
func (_ OpPow) IsClosedField() bool                   { return false }
func (_ OpPow) IsStrictClosedField() bool             { return false }

func (_ OpPow) DetectEdgeCaseLeft(val float64, n Number) Number {
	if val == 0 {
		lt_zero := n.Less(_zero)
		if lt_zero.IsTrue() {
			return _inf
		} else if lt_zero.IsFalse() {
			return _zero
		}
	} else if val == 1 {
		return _one
	}
	return Number{}
}

func (_ OpPow) DetectEdgeCaseRight(n Number, val float64) Number {
	if val == 0 {
		return _one
	} else if val == 1 {
		return n
	}
	return Number{}
}

func (_ OpPow) DetectEdgeCaseSame(n Number) Number     { return Number{} }
func (_ OpPow) PreprocessRangeLeft(r *NRange) *NRange  { return r }
func (_ OpPow) PreprocessRangeRight(r *NRange) *NRange { return r }
func (_ OpPow) IsResultInt() Boolean                   { return Boolean{} }

func (_ OpPow) ResultConstraints(x, y, result Number) Number {
	return result
}

func operator(x, y Number, op ArithmeticOperationBinary) (Number, error) {
	if x.p.next != nil || y.p.next != nil {
		panic("next is unsupported yet")
	}

	xConstant := x.IsConstant()
	yConstant := y.IsConstant()

	if xConstant && yConstant {
		xpVal, ypVal := x.p.val, y.p.val
		res, err := op.Compute(xpVal, ypVal)
		if err != nil {
			if err.Error() == NAN_STR {
				return NewNumberConst(math.NaN()), nil
			}
			return Number{}, err
		}
		if res == xpVal {
			return x, nil
		} else if res == ypVal {
			return y, nil
		}
		return NewNumberConst(res), nil
	}

	xp, yp := x.p, y.p

	if xp == yp {
		edgeRes := op.DetectEdgeCaseSame(x)
		if edgeRes.IsValid() {
			return op.ResultConstraints(x, y, edgeRes), nil
		}
	}

	xRange := xp.valRange
	yRange := yp.valRange

	var r []*NRange
	if xRange == nil && yRange == nil {
		r = []*NRange{newRange()}
	} else {
		if xRange == nil {
			if xConstant {
				edgeRes := op.DetectEdgeCaseLeft(xp.val, y)
				if edgeRes.IsValid() {
					return edgeRes, nil
				}

				xRange = newRangeSegment(xp.val, xp.val)
			} else {
				xRange = newRange()
			}
		} else if yRange == nil {
			if yConstant {
				edgeRes := op.DetectEdgeCaseRight(x, yp.val)
				if edgeRes.IsValid() {
					return edgeRes, nil
				}

				yRange = newRangeSegment(yp.val, yp.val)
			} else {
				yRange = newRange()
			}
		}

		var err error
		r, err = xRange.ArithmeticOperation(yRange, op)
		//fmt.Printf("r: %+v\n", r)
		if err != nil {
			return Number{}, err
		}
		if len(r) == 1 && r[0].IsConstant() {
			return NewNumberConst(r[0].lVal), nil
		}
	}

	resInt := op.IsResultInt()

	needAdjust := !op.IsStrictClosedField()

	res := NewNumberRange(r...)
	if op.IsClosedField() {
		resInt = x.IsInteger().And(y.IsInteger())
		if !resInt.IsFalse() || op.IsStrictClosedField() {
			needAdjust = true
		}
	}

	if resInt.IsValid() {
		res.p.integer = resInt
		needAdjust = true
	}

	var err error
	if needAdjust {
		res, err = res.RangeAdjust()
	}
	return op.ResultConstraints(x, y, res), err
}

func (n Number) Add(o Number) (Number, error) { return operator(n, o, OpAdd{}) }
func (n Number) Sub(o Number) (Number, error) { return operator(n, o, OpSub{}) }
func (n Number) Mul(o Number) (Number, error) { return operator(n, o, OpMul{}) }
func (n Number) Div(o Number) (Number, error) { return operator(n, o, OpDiv{}) }
func (n Number) Pow(o Number) (Number, error) { return operator(n, o, OpPow{}) }

func (n Number) IDiv(o Number) (Number, error) {
	res, err := operator(n, o, OpIDiv{})
	if err == nil && !res.IsNaN().IsTrue() {
		res.p.integer = NewBooleanConst(BTrue, nil)
	}
	return res, err
}

func (n Number) Max(numbers []Number) Number {
	if n.p.next != nil {
		panic("next is unsupported yet")
	}
	if n.IsNaN().IsTrue() {
		return n
	}

	for _, num := range numbers {
		if num.IsNaN().IsTrue() {
			return num
		}

		if gt := num.Greater(n); gt.IsTrue() {
			n = num
		} else if gt.IsUnknown() {

			panic("NYI")
		}
	}
	return n
}

func (n Number) Min(numbers []Number) Number {
	if n.p.next != nil {
		panic("next is unsupported yet")
	}
	if n.IsNaN().IsTrue() {
		return n
	}

	for _, num := range numbers {
		if num.IsNaN().IsTrue() {
			return num
		}

		if lt := num.Less(n); lt.IsTrue() {
			n = num
		} else { // TODO: False & Unknown
			panic("NYI")
		}
	}
	return n
}

func (n Number) Floor() Number {
	arithmeticallyCorrect := false
	inverted := false

	return n.FloorWithOpt(arithmeticallyCorrect, inverted)
}

func (n Number) Ceil() Number {
	arithmeticallyCorrect := false
	inverted := true

	return n.FloorWithOpt(arithmeticallyCorrect, inverted)
}

func (n Number) FloorWithOpt(arithmeticallyCorrect, inverted bool) Number {
	var p *NumberPrivate
	if n.p.next != nil {
		panic("next is unsupported yet")
	}
	if n.IsConstant() {
		newVal := floor(n.p.val, arithmeticallyCorrect)
		if newVal == n.p.val {
			return n
		}
		if inverted {
			newVal++
		}
		p = &NumberPrivate{
			val:     newVal,
			integer: NewBooleanConst(BTrue, nil),
		}
	} else if n.p.valRange != nil {
		newRange := n.p.valRange.Floor(arithmeticallyCorrect, inverted)
		p = &NumberPrivate{
			val:      0,
			integer:  NewBooleanConst(BTrue, nil),
			valRange: newRange,
		}
	}

	res, _ := Number{p: p}.RangeAdjust()
	return res
}

func (n Number) Abs() Number {
	if n.p.next != nil {
		panic("next is unsupported yet")
	}
	if n.IsConstant() {
		if n.p.val < 0 {
			return NewNumberConst(-n.p.val)
		}
	} else {
		newRange := n.p.valRange.Abs()
		if newRange != n.p.valRange {
			res := NewNumberRange(newRange)
			constraints := make([]NumberConstraint, 1)
			if res.Greater(_zero).IsTrue() {
				constraints[0] = NewNumberGreater(n)
			} else {
				constraints[0] = NewNumberGreaterEqual(n)
			}
			res.p.constraints = constraints
			return res
		}
	}
	return n
}

func removeOnce(s []float64, val float64) []float64 {
	for i, v := range s {
		if v == val {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (n Number) Sign() Number {
	possibleSigns := n.p.Sign()
	canBeZero := !n.IsInteger().IsFalse()

	if !canBeZero {
		possibleSigns = removeOnce(possibleSigns, 0)
	}

	var p, curr *NumberPrivate
	for _, sign := range possibleSigns {
		next := &NumberPrivate{
			val:     sign,
			integer: NewBooleanConst(BTrue, nil),
		}

		if curr == nil {
			curr = next
			p = curr
		} else {
			curr.next = next
			curr = next
		}
	}

	return Number{p: p}
}

func (n Number) Split(splitPoint float64) (Number, Number) {
	if n.p.next != nil {
		panic("next is unsupported yet")
	}
	if n.IsConstant() {
		if n.p.val < splitPoint {
			return n, Number{}
		}
		return Number{}, n
	}

	l, r := n.p.valRange.Split(splitPoint)
	if l != nil && r != nil {
		lNum, _ := Number{p: &NumberPrivate{
			val:      0,
			integer:  n.p.integer,
			valRange: l,
		}}.RangeAdjust()

		rNum, _ := Number{p: &NumberPrivate{
			val:      0,
			integer:  n.p.integer,
			valRange: r,
		}}.RangeAdjust()

		return lNum, rNum
	} else if l == nil {
		return Number{}, n
	} else if r == nil {
		return n, Number{}
	}

	panic("failed to split number (should not reach there: something went wrong)")
	return Number{}, Number{}
}

func NewNumber() Number {
	return NewNumberRange(newRange())
}

func NewNumberConst(v float64) Number {
	is_integer := BUnknown
	if !math.IsNaN(v) {
		if v == math.Floor(v) {
			is_integer = BTrue
		} else {
			is_integer = BFalse
		}
	}
	return Number{p: &NumberPrivate{
		val:     v,
		integer: NewBooleanConst(is_integer, nil),
	}}
}

func NewNumberRange(rVec ...*NRange) Number {
	rVecLen := len(rVec)
	if rVecLen == 0 {
		panic("no ranges passed to NewNumberRange constructor")
	}

	p := newNumberPrivate(rVec[0])
	pPrev := p

	for i := 1; i < rVecLen; i++ {
		pCurr := newNumberPrivate(rVec[i])
		pPrev.next = pCurr
		pPrev = pCurr
	}

	res, _ := Number{p: p}.RangeAdjust()
	return res
}

func NewNumberSegment(l, r float64) Number {
	return NewNumberRange(newRangeSegment(l, r))
}
