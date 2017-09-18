package virtual_types

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

func withoutIntegerConstrints(n Number) Number {
	n.p.integer.p.constraints = nil
	return n
}

type NumberSuite struct {
	suite.Suite
	Zero                   Number
	One                    Number
	Two                    Number
	Four                   Number
	Five                   Number
	MinusOne               Number
	Unknown                Number
	Unknown_copy           Number
	Unknown_other          Number
	OneFiveSeg             Number
	OneFiveSeg_copy        Number
	OneFiveSeg_other       Number
	OneFiveSegInt          Number
	ZeroOneSeg             Number
	ZeroOneSegOpen         Number
	TwoFourSegOpen         Number
	MinusTenTwoAndHalfSeg  Number
	MinusTenTenSeg         Number
	TwoAndHalfFourInterval Number
	Inf                    Number
	Positive               Number
	Positive_copy          Number
	Positive_other         Number
	Negative               Number
}

func (s *NumberSuite) SetupTest() {
	s.Zero = NewNumberConst(0)
	s.One = NewNumberConst(1)
	s.Two = NewNumberConst(2)
	s.Four = NewNumberConst(4)
	s.Five = NewNumberConst(5)
	s.MinusOne = NewNumberConst(-1)

	s.Unknown = NewNumber()
	s.Unknown_copy = s.Unknown

	s.Unknown_other = NewNumber()

	s.OneFiveSeg = NewNumberSegment(1, 5)
	s.OneFiveSeg_copy = s.OneFiveSeg
	s.OneFiveSeg_other = NewNumberSegment(1, 5)

	s.OneFiveSegInt = NewNumberSegment(1, 5)
	s.OneFiveSegInt.p.integer = NewBooleanConst(BTrue, nil)

	s.ZeroOneSeg = NewNumberSegment(0, 1)
	s.ZeroOneSegOpen = NewNumberRange(&NRange{
		lVal: 0, lIncluding: true,
		rVal: 1, rIncluding: false,
	})

	s.TwoFourSegOpen = NewNumberRange(&NRange{
		lVal: 2, lIncluding: true,
		rVal: 4, rIncluding: false,
	})
	s.MinusTenTwoAndHalfSeg = NewNumberSegment(-10, 2.5)
	s.MinusTenTenSeg = NewNumberSegment(-10, 10)
	s.TwoAndHalfFourInterval = NewNumberRange(&NRange{
		lVal: 2.5, lIncluding: false,
		rVal: 4.0, rIncluding: false,
	})

	s.Inf = NewNumberConst(math.Inf(1))

	s.Positive = NewNumberRange(&NRange{
		lVal: 0, lIncluding: false,
		rVal: math.Inf(1), rIncluding: true,
	})
	s.Positive_copy = s.Positive
	s.Positive_other = NewNumberRange(&NRange{
		lVal: 0, lIncluding: false,
		rVal: math.Inf(1), rIncluding: true,
	})
	s.Negative = NewNumberRange(&NRange{
		lVal: math.Inf(-1), lIncluding: true,
		rVal: 0, rIncluding: false,
	})
}

func (s *NumberSuite) TestIsSame() {
	assert := assert.New(s.T())

	assert.True(s.Zero.IsSame(s.Zero))
	assert.True(s.Two.IsSame(NewNumberConst(2)))

	assert.False(s.Zero.IsSame(s.One))
	assert.False(s.One.IsSame(s.OneFiveSeg))

	assert.True(s.OneFiveSeg.IsSame(NewNumberSegment(1, 5)))
	assert.True(NewNumberSegment(1, 5).IsSame(s.OneFiveSeg))

	assert.True(s.OneFiveSeg.IsSame(s.OneFiveSeg_copy))
	assert.True(s.OneFiveSeg_copy.IsSame(s.OneFiveSeg))

	assert.True(s.OneFiveSeg.IsSame(s.OneFiveSeg_other))
	assert.True(s.OneFiveSeg_other.IsSame(s.OneFiveSeg))

	assert.True(s.Unknown.IsSame(s.Unknown))

	assert.True(s.Unknown.IsSame(s.Unknown_copy))
	assert.True(s.Unknown_copy.IsSame(s.Unknown))

	assert.True(s.Unknown.IsSame(s.Unknown_other))
	assert.True(s.Unknown_other.IsSame(s.Unknown))

	assert.False(s.Unknown.IsSame(s.Two))
	assert.False(s.Two.IsSame(s.Unknown))

	assert.False(s.ZeroOneSeg.IsSame(s.Unknown_other))
	assert.False(s.Unknown_other.IsSame(s.ZeroOneSeg))
}

func (s *NumberSuite) TestEqualSimple() {
	assert := assert.New(s.T())

	assert.True(s.Zero.Equal(s.Zero).IsTrue())
	assert.True(s.Zero.Equal(s.One).IsFalse())
	assert.True(s.One.Equal(s.MinusOne).IsFalse())

	assert.True(s.Zero.Equal(s.Unknown).IsUnknown())
	assert.True(s.MinusOne.Equal(s.Unknown).IsUnknown())

	assert.True(s.Unknown.Equal(s.Unknown).IsTrue())

	assert.True(s.Unknown.Equal(s.Unknown_copy).IsTrue())
	assert.True(s.Unknown_copy.Equal(s.Unknown).IsTrue())

	assert.True(s.Unknown.Equal(s.Unknown_other).IsUnknown())
	assert.True(s.Unknown_other.Equal(s.Unknown).IsUnknown())
}

func (s *NumberSuite) TestEqualRange() {
	assert := assert.New(s.T())

	assert.True(s.OneFiveSeg.Equal(s.Two).IsUnknown())
	assert.True(s.One.Equal(s.OneFiveSeg).IsUnknown())

	assert.True(s.Two.Equal(s.TwoFourSegOpen).IsUnknown())
	assert.True(s.Four.Equal(s.TwoFourSegOpen).IsFalse())

	assert.True(s.OneFiveSeg.Equal(s.TwoFourSegOpen).IsUnknown())
	assert.True(s.Zero.Equal(s.TwoFourSegOpen).IsFalse())

	assert.True(s.MinusTenTwoAndHalfSeg.Equal(s.TwoAndHalfFourInterval).IsFalse())

	assert.True(s.OneFiveSeg.Equal(s.OneFiveSeg).IsTrue())

	assert.True(s.OneFiveSeg.Equal(s.OneFiveSeg_copy).IsTrue())
	assert.True(s.OneFiveSeg_copy.Equal(s.OneFiveSeg).IsTrue())

	assert.True(s.OneFiveSeg.Equal(s.OneFiveSeg_other).IsUnknown())
	assert.True(s.OneFiveSeg_other.Equal(s.OneFiveSeg).IsUnknown())

	assert.True(s.OneFiveSeg.Equal(s.ZeroOneSeg).IsUnknown())
	assert.True(s.ZeroOneSeg.Equal(s.OneFiveSeg).IsUnknown())

	assert.True(s.OneFiveSeg.Equal(s.ZeroOneSegOpen).IsFalse())
	assert.True(s.ZeroOneSegOpen.Equal(s.OneFiveSeg).IsFalse())
}

func (s *NumberSuite) TestLessSimple() {
	assert := assert.New(s.T())

	assert.True(s.Zero.Less(s.Zero).IsFalse())
	assert.True(s.Zero.Less(s.One).IsTrue())
	assert.True(s.Zero.Less(s.Two).IsTrue())
	assert.True(s.Zero.Less(s.Unknown).IsUnknown())

	assert.True(s.Unknown.Less(s.Zero).IsUnknown())
	assert.True(s.Unknown.Less(s.One).IsUnknown())
	assert.True(s.Unknown.Less(s.Two).IsUnknown())

	assert.True(s.Unknown.Less(s.Unknown).IsFalse())

	assert.True(s.Unknown.Less(s.Unknown_copy).IsFalse())
	assert.True(s.Unknown_copy.Less(s.Unknown).IsFalse())

	assert.True(s.Unknown.Less(s.Unknown_other).IsUnknown())
	assert.True(s.Unknown_other.Less(s.Unknown).IsUnknown())
}

func (s *NumberSuite) TestLessRange() {
	assert := assert.New(s.T())

	assert.True(s.Zero.Less(s.TwoFourSegOpen).IsTrue())
	assert.True(s.TwoFourSegOpen.Less(s.Zero).IsFalse())

	assert.True(s.Two.Less(s.TwoFourSegOpen).IsUnknown())
	assert.True(s.TwoFourSegOpen.Less(s.Two).IsFalse())

	assert.True(s.OneFiveSeg.Less(s.One).IsFalse())
	assert.True(s.One.Less(s.OneFiveSeg).IsUnknown())

	assert.True(s.TwoFourSegOpen.Less(s.Four).IsTrue())
	assert.True(s.Four.Less(s.TwoFourSegOpen).IsFalse())

	assert.True(s.MinusTenTwoAndHalfSeg.Less(s.MinusOne).IsUnknown())
	assert.True(s.MinusOne.Less(s.MinusTenTwoAndHalfSeg).IsUnknown())

	assert.True(s.TwoAndHalfFourInterval.Less(s.Four).IsTrue())
	assert.True(s.Four.Less(s.TwoAndHalfFourInterval).IsFalse())

	assert.True(s.OneFiveSeg.Less(s.OneFiveSeg).IsFalse())

	assert.True(s.OneFiveSeg.Less(s.OneFiveSeg_copy).IsFalse())
	assert.True(s.OneFiveSeg_copy.Less(s.OneFiveSeg).IsFalse())

	assert.True(s.OneFiveSeg.Less(s.OneFiveSeg_other).IsUnknown())
	assert.True(s.OneFiveSeg_other.Less(s.OneFiveSeg).IsUnknown())

	assert.True(s.OneFiveSeg.Less(s.ZeroOneSeg).IsFalse())
	assert.True(s.ZeroOneSeg.Less(s.OneFiveSeg).IsUnknown())

	assert.True(s.OneFiveSeg.Less(s.ZeroOneSegOpen).IsFalse())
	assert.True(s.ZeroOneSegOpen.Less(s.OneFiveSeg).IsTrue())

	assert.True(s.ZeroOneSeg.Less(s.ZeroOneSegOpen).IsUnknown())
	assert.True(s.ZeroOneSegOpen.Less(s.ZeroOneSeg).IsUnknown())
}

func (s *NumberSuite) TestLessEqualSimple() {
	assert := assert.New(s.T())

	assert.True(s.One.LessEqual(s.One).IsTrue())

	assert.True(s.One.Equal(s.OneFiveSeg).IsUnknown())
	assert.True(s.One.LessEqual(s.OneFiveSeg).IsTrue())
	assert.True(s.OneFiveSeg.LessEqual(s.One).IsUnknown())

	assert.True(s.Unknown.LessEqual(s.Unknown).IsTrue())

	assert.True(s.Unknown.LessEqual(s.Unknown_copy).IsTrue())
	assert.True(s.Unknown_copy.LessEqual(s.Unknown).IsTrue())

	assert.True(s.Unknown.LessEqual(s.Unknown_other).IsUnknown())
	assert.True(s.Unknown_other.LessEqual(s.Unknown).IsUnknown())
}

func (s *NumberSuite) TestLessEqualRange() {
	assert := assert.New(s.T())

	assert.True(s.OneFiveSeg.LessEqual(s.OneFiveSeg).IsTrue())

	assert.True(s.OneFiveSeg.LessEqual(s.OneFiveSeg_copy).IsTrue())
	assert.True(s.OneFiveSeg_copy.LessEqual(s.OneFiveSeg).IsTrue())

	assert.True(s.OneFiveSeg.LessEqual(s.OneFiveSeg_other).IsUnknown())
	assert.True(s.OneFiveSeg_other.LessEqual(s.OneFiveSeg).IsUnknown())

	assert.True(s.OneFiveSeg.LessEqual(s.ZeroOneSeg).IsUnknown())
	assert.True(s.ZeroOneSeg.LessEqual(s.OneFiveSeg).IsUnknown())

	assert.True(s.OneFiveSeg.LessEqual(s.ZeroOneSegOpen).IsFalse())
	assert.True(s.ZeroOneSegOpen.LessEqual(s.OneFiveSeg).IsTrue())
}

func (s *NumberSuite) TestGreaterEqualSimple() {
	assert := assert.New(s.T())

	assert.True(s.Five.GreaterEqual(s.Five).IsTrue())

	assert.True(s.Five.Equal(s.OneFiveSeg).IsUnknown())
	assert.True(s.Five.GreaterEqual(s.OneFiveSeg).IsTrue())
	assert.True(s.OneFiveSeg.GreaterEqual(s.Five).IsUnknown())

	assert.True(s.Unknown.GreaterEqual(s.Unknown).IsTrue())

	assert.True(s.Unknown.GreaterEqual(s.Unknown_copy).IsTrue())
	assert.True(s.Unknown_copy.GreaterEqual(s.Unknown).IsTrue())

	assert.True(s.Unknown.GreaterEqual(s.Unknown_other).IsUnknown())
	assert.True(s.Unknown_other.GreaterEqual(s.Unknown).IsUnknown())
}

func (s *NumberSuite) TestGreaterEqualRange() {
	assert := assert.New(s.T())

	assert.True(s.OneFiveSeg.GreaterEqual(s.OneFiveSeg).IsTrue())

	assert.True(s.OneFiveSeg.GreaterEqual(s.OneFiveSeg_copy).IsTrue())
	assert.True(s.OneFiveSeg_copy.GreaterEqual(s.OneFiveSeg).IsTrue())

	assert.True(s.OneFiveSeg.GreaterEqual(s.OneFiveSeg_other).IsUnknown())
	assert.True(s.OneFiveSeg_other.GreaterEqual(s.OneFiveSeg).IsUnknown())

	assert.True(s.OneFiveSeg.GreaterEqual(s.ZeroOneSeg).IsUnknown())
	assert.True(s.ZeroOneSeg.GreaterEqual(s.OneFiveSeg).IsUnknown())

	assert.True(s.OneFiveSeg.GreaterEqual(s.ZeroOneSegOpen).IsTrue())
	assert.True(s.ZeroOneSegOpen.GreaterEqual(s.OneFiveSeg).IsFalse())
}

func (s *NumberSuite) TestAddSimple() {
	assert := assert.New(s.T())

	two, err := s.One.Add(s.One)

	assert.Nil(err)
	assert.True(two.Equal(s.Two).IsTrue())

	zero, err := s.One.Add(s.MinusOne)

	assert.Nil(err)
	assert.True(zero.Equal(s.Zero).IsTrue())

	one, err := s.One.Add(s.Zero)

	assert.Nil(err)
	assert.True(one.Equal(s.One).IsTrue())

	unknown, err := s.One.Add(s.Unknown)

	assert.Nil(err)
	assert.True(unknown.IsUnknown())
	assert.True(unknown.Greater(s.Unknown).IsTrue())
	assert.True(unknown.Greater(s.Unknown_copy).IsTrue())
	assert.True(unknown.Greater(s.Unknown_other).IsUnknown())
}

func (s *NumberSuite) TestAddRange() {
	assert := assert.New(s.T())

	two_six_seg, err := s.One.Add(s.OneFiveSeg)

	expected := NewNumberSegment(2, 6)
	expected.p.constraints = []NumberConstraint{
		NewNumberGreater(s.OneFiveSeg),
		NewNumberGreater(s.One), // Can be omitted G([2; 6] > 1)
	}

	assert.Nil(err)
	assert.Equal(expected, two_six_seg)
	assert.True(two_six_seg.Greater(s.OneFiveSeg).IsTrue())
	assert.True(two_six_seg.Greater(s.OneFiveSeg_copy).IsTrue())
	assert.True(two_six_seg.Greater(s.OneFiveSeg_other).IsUnknown())

	two_six_seg, err = s.OneFiveSeg.Add(s.One)

	c := expected.p.constraints
	c[0], c[1] = c[1], c[0] // Order of constraints matters for assert.Equal

	assert.Nil(err)
	assert.Equal(expected, withoutIntegerConstrints(two_six_seg))

	two_ten_seg, err := s.OneFiveSeg.Add(s.OneFiveSeg)

	expected = NewNumberSegment(2, 10)
	expected.p.constraints = []NumberConstraint{
		NewNumberGreater(s.OneFiveSeg),
	}

	assert.Nil(err)
	assert.Equal(expected, withoutIntegerConstrints(two_ten_seg))
	assert.True(two_ten_seg.Greater(s.OneFiveSeg).IsTrue())

	two_ten_seg, err = s.OneFiveSeg_copy.Add(s.OneFiveSeg_other)

	expected.p.constraints = append(expected.p.constraints, NewNumberGreater(s.OneFiveSeg_other))

	assert.Nil(err)
	assert.Equal(expected, withoutIntegerConstrints(two_ten_seg))
	assert.True(two_ten_seg.Greater(s.OneFiveSeg_copy).IsTrue())
	assert.True(two_ten_seg.Greater(s.OneFiveSeg_other).IsTrue())
	assert.True(s.OneFiveSeg_copy.Less(two_ten_seg).IsTrue())
	assert.True(s.OneFiveSeg_other.Less(two_ten_seg).IsTrue())

	five_six_seg_open, err := s.Five.Add(s.ZeroOneSegOpen)

	expected = NewNumberRange(&NRange{
		lVal: 5, lIncluding: true,
		rVal: 6, rIncluding: false,
	})
	expected.p.constraints = []NumberConstraint{
		NewNumberGreater(s.ZeroOneSegOpen),
	}

	assert.Nil(err)
	assert.Equal(expected, five_six_seg_open)

	assert.True(five_six_seg_open.Greater(s.ZeroOneSegOpen).IsTrue())
	assert.True(five_six_seg_open.Greater(s.Five).IsUnknown())
}

func (s *NumberSuite) TestAddZeroEdgeCases() {
	assert := assert.New(s.T())

	one, err := s.One.Add(s.Zero)

	assert.Nil(err)
	assert.Equal(s.One, one)
	assert.True(s.One.p == one.p)

	one, err = s.Zero.Add(s.One)

	assert.Nil(err)
	assert.Equal(s.One, one)
	assert.True(s.One.p == one.p)

	one_five_seg, err := s.OneFiveSeg.Add(s.Zero)

	assert.Nil(err)
	assert.Equal(s.OneFiveSeg, one_five_seg)
	assert.True(s.OneFiveSeg.p == one_five_seg.p)

	unknown, err := s.Unknown.Add(s.Zero)

	assert.Nil(err)
	assert.Equal(s.Unknown, unknown)
	assert.True(s.Unknown.p == unknown.p)
}

func (s *NumberSuite) TestSubSimple() {
	assert := assert.New(s.T())

	one, err := s.Two.Sub(s.One)

	assert.Nil(err)
	assert.Equal(s.One, one)

	unknown, err := s.Unknown.Sub(s.One)

	assert.Nil(err)
	assert.Equal(s.Unknown, withoutIntegerConstrints(unknown))
	assert.True(unknown.Less(s.Unknown).IsTrue(), "NYI")
}

func (s *NumberSuite) TestSubZeroEdgeCases() {
	assert := assert.New(s.T())

	one, err := s.One.Sub(s.Zero)

	assert.Nil(err)
	assert.Equal(s.One, one)
	assert.True(s.One.p == one.p)

	one_five_seg, err := s.OneFiveSeg.Sub(s.Zero)

	assert.Nil(err)
	assert.Equal(s.OneFiveSeg, one_five_seg)
	assert.True(s.OneFiveSeg.p == one_five_seg.p)

	unknown, err := s.Unknown.Sub(s.Zero)

	assert.Nil(err)
	assert.Equal(s.Unknown, unknown)
	assert.True(s.Unknown.p == unknown.p)
}

func (s *NumberSuite) TestSubSameEdgeCases() {
	assert := assert.New(s.T())

	zero, err := s.One.Sub(s.One)

	assert.Nil(err)
	assert.Equal(s.Zero, zero)

	zero, err = s.OneFiveSeg.Sub(s.OneFiveSeg_copy)

	assert.Nil(err)
	assert.Equal(s.Zero, zero)

	minus_4_to_4_seg, err := s.OneFiveSeg.Sub(s.OneFiveSeg_other)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(-4, 4), withoutIntegerConstrints(minus_4_to_4_seg))

	zero, err = s.Unknown.Sub(s.Unknown_copy)

	assert.Nil(err)
	assert.Equal(s.Zero, zero)

	unknown, err := s.Unknown.Sub(s.Unknown_other)

	assert.Nil(err)
	assert.Equal(s.Unknown, withoutIntegerConstrints(unknown))
	assert.True(unknown.Equal(s.Unknown).IsUnknown())
	assert.True(unknown.Equal(s.Unknown_other).IsUnknown())
}

func (s *NumberSuite) TestMulSimple() {
	assert := assert.New(s.T())

	one, err := s.One.Mul(s.One)

	assert.Nil(err)
	assert.Equal(s.One, one)
	assert.True(s.One.p == one.p)

	five, err := s.One.Mul(s.Five)

	assert.Nil(err)
	assert.Equal(s.Five, five)
	assert.True(s.Five.p == five.p)

	minus_four, err := s.Four.Mul(s.MinusOne)

	assert.Nil(err)
	assert.Equal(NewNumberConst(-4), minus_four)

	minus_twenty, err := minus_four.Mul(s.Five)

	assert.Nil(err)
	assert.Equal(NewNumberConst(-20), minus_twenty)
}

func (s *NumberSuite) TestMulRange() {
	assert := assert.New(s.T())

	_2_to_10_seg, err := s.OneFiveSeg.Mul(s.Two)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(2, 10), withoutIntegerConstrints(_2_to_10_seg))

	_2_to_10_seg, err = s.Two.Mul(s.OneFiveSeg_other)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(2, 10), _2_to_10_seg)

	_1_to_25_seg, err := s.OneFiveSeg.Mul(s.OneFiveSeg)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(1, 25), withoutIntegerConstrints(_1_to_25_seg))

	_0_to_1_seg_open, err := s.ZeroOneSegOpen.Mul(s.ZeroOneSegOpen)

	assert.Nil(err)
	assert.Equal(s.ZeroOneSegOpen, withoutIntegerConstrints(_0_to_1_seg_open))

	_minus_10_to_10, err := s.MinusTenTenSeg.Mul(s.One.Negate())

	assert.Nil(err)
	assert.Equal(s.MinusTenTenSeg, withoutIntegerConstrints(_minus_10_to_10))
}

func (s *NumberSuite) TestMulPositive() {
	assert := assert.New(s.T())

	positive, err := s.Positive.Mul(s.Two)

	assert.Nil(err)
	assert.Equal(s.Positive, withoutIntegerConstrints(positive))
	assert.True(s.Positive.Less(positive).IsUnknown())

	positive, err = s.Positive.Mul(s.Positive_other)

	assert.Nil(err)
	assert.Equal(s.Positive, withoutIntegerConstrints(positive))
	assert.True(s.Positive.Less(positive).IsUnknown())

	negative, err := s.Positive.Mul(s.One.Negate())

	assert.Nil(err)
	assert.Equal(s.Negative, withoutIntegerConstrints(negative))
}

func (s *NumberSuite) TestMulOneEdgeCases() {
	assert := assert.New(s.T())

	unknown, err := s.One.Mul(s.Unknown)

	assert.Nil(err)
	assert.Equal(s.Unknown, unknown)

	unknown, err = s.Unknown.Mul(s.One)

	assert.Nil(err)
	assert.Equal(s.Unknown, unknown)

	one_five_seg, err := s.One.Mul(s.OneFiveSeg)

	assert.Nil(err)
	assert.Equal(s.OneFiveSeg, one_five_seg)

	one_five_seg, err = s.OneFiveSeg.Mul(s.One)

	assert.Nil(err)
	assert.Equal(s.OneFiveSeg, one_five_seg)
}

func (s *NumberSuite) TestMulZeroEdgeCases() {
	assert := assert.New(s.T())

	zero, err := s.Zero.Mul(s.Unknown)

	assert.Nil(err)
	assert.Equal(s.Zero, zero)

	zero, err = s.Unknown.Mul(s.Zero)

	assert.Nil(err)
	assert.Equal(s.Zero, zero)

	zero, err = s.Zero.Mul(s.OneFiveSeg)

	assert.Nil(err)
	assert.Equal(s.Zero, zero)

	zero, err = s.OneFiveSeg.Mul(s.Zero)

	assert.Nil(err)
	assert.Equal(s.Zero, zero)
}

func (s *NumberSuite) TestDivSimple() {
	assert := assert.New(s.T())

	two_and_half, err := s.Five.Div(s.Two)

	assert.Nil(err)
	assert.Equal(NewNumberConst(2.5), two_and_half)

	unknown, err := s.Unknown.Div(s.Two)

	assert.Nil(err)
	assert.Equal(s.Unknown, unknown)
	assert.False(s.Unknown.p == unknown.p)
	assert.True(unknown.Less(s.Unknown).IsUnknown())

	unknown, err = s.Unknown.Div(s.Unknown_other)

	assert.Nil(err)
	assert.Equal(s.Unknown, unknown)
	assert.True(s.Unknown.Equal(unknown).IsUnknown())

	/*one, err := s.ZeroOneSegOpen.Div(s.ZeroOneSegOpen)

	assert.Nil(err)
	assert.Equal(s.One, one)*/
}

func (s *NumberSuite) TestDivPositive() {
	assert := assert.New(s.T())

	positive, err := s.Positive.Div(s.Two)

	assert.Nil(err)
	assert.Equal(s.Positive, positive)
	assert.False(s.Positive.p == positive.p)
	//assert.True(positive.Less(s.Positive).IsTrue())

	positive, err = s.Positive.Div(s.Positive_other)

	assert.Nil(err)
	assert.Equal(s.Positive, positive)
	assert.False(s.Positive.p == positive.p)
	assert.True(positive.Less(s.Positive).IsUnknown())
}

func (s *NumberSuite) TestDivRange() {
	assert := assert.New(s.T())

	_0_25_to_2_seg_open, err := s.Four.Div(NewNumberRange(&NRange{
		lVal: 2, lIncluding: false,
		rVal: 16, rIncluding: true,
	}))

	assert.Nil(err)
	assert.Equal(NewNumberRange(&NRange{
		lVal: 0.25, lIncluding: true,
		rVal: 2.0, rIncluding: false,
	}), _0_25_to_2_seg_open)

	_0_5_to_2_5_seg, err := s.OneFiveSeg.Div(s.Two)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(0.5, 2.5), _0_5_to_2_5_seg)

	_0_to_0_5_seg_open, err := s.ZeroOneSegOpen.Div(s.Two)

	assert.Nil(err)
	assert.Equal(NewNumberRange(&NRange{
		lVal: 0.0, lIncluding: true,
		rVal: 0.5, rIncluding: false,
	}), _0_to_0_5_seg_open)

	_0_2_to_5_seg, err := s.OneFiveSeg.Div(s.OneFiveSeg_other)

	assert.Nil(err)
	assert.Equal(NewNumberRange(&NRange{
		lVal: 0.2, lIncluding: true,
		rVal: 5.0, rIncluding: true,
	}), _0_2_to_5_seg)

	one_inf_open_ray, err := s.OneFiveSeg.Div(s.ZeroOneSegOpen)

	assert.Nil(err)
	assert.Equal(NewNumberRange(&NRange{
		lVal: 1, lIncluding: false,
		rVal: math.Inf(1), rIncluding: true,
	}), one_inf_open_ray)
}

func (s *NumberSuite) TestDivOneEdgeCases() {
	assert := assert.New(s.T())

	five, err := s.Five.Div(s.One)

	assert.Nil(err)
	assert.Equal(s.Five, five)
	assert.True(s.Five.p == five.p)

	one_five_seg, err := s.OneFiveSeg.Div(s.One)

	assert.Nil(err)
	assert.Equal(s.OneFiveSeg, one_five_seg)
	assert.True(s.OneFiveSeg.p == one_five_seg.p)
}

func (s *NumberSuite) TestDivZeroEdgeCases() {
	assert := assert.New(s.T())

	invalid_num, err := s.One.Div(s.Zero)

	if assert.NotNil(err) {
		assert.Equal("division by zero", err.Error())
	}
	assert.False(invalid_num.IsValid())

	invalid_num, err = s.OneFiveSeg.Div(s.Zero)

	if assert.NotNil(err) {
		assert.Equal("division by zero", err.Error())
	}
	assert.False(invalid_num.IsValid())

	invalid_num, err = s.ZeroOneSegOpen.Div(s.Zero)

	if assert.NotNil(err) {
		assert.Equal("division by zero", err.Error())
	}
	assert.False(invalid_num.IsValid())
}

func (s *NumberSuite) TestDivSameEdgeCases() {
	assert := assert.New(s.T())

	one, err := s.One.Div(s.One)

	assert.Nil(err)
	assert.Equal(s.One, one)
	assert.True(s.One.p == one.p)

	one, err = s.OneFiveSeg.Div(s.OneFiveSeg)

	assert.Nil(err)
	assert.Equal(s.One, one)

	one, err = s.Unknown.Div(s.Unknown_copy)

	assert.Nil(err)
	assert.Equal(s.One, one)
}

func (s *NumberSuite) TestPowSimple() {
	assert := assert.New(s.T())

	sixteen, err := s.Two.Pow(s.Four)

	assert.Nil(err)
	assert.Equal(NewNumberConst(16), sixteen)

	_0_25, err := s.Four.Pow(s.One.Negate())

	assert.Nil(err)
	assert.Equal(NewNumberConst(0.25), _0_25)
}

func (s *NumberSuite) TestPowRange() {
	assert := assert.New(s.T())

	_2_to_32_seg, err := s.Two.Pow(s.OneFiveSeg)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(2, 32), _2_to_32_seg)

	_1_to_25_seg, err := s.OneFiveSeg.Pow(s.Two)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(1, 25), _1_to_25_seg)

	_1_to_3125_seg, err := s.OneFiveSeg.Pow(s.OneFiveSeg)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(1, 3125), _1_to_3125_seg)

	zero_one_seg_open, err := s.ZeroOneSegOpen.Pow(s.Two)

	assert.Nil(err)
	assert.Equal(s.ZeroOneSegOpen, zero_one_seg_open)
	assert.False(s.ZeroOneSegOpen.p == zero_one_seg_open.p)
}

func (s *NumberSuite) TestPowZeroEdgeCases() {
	assert := assert.New(s.T())

	one, err := s.Two.Pow(s.Zero)

	assert.Nil(err)
	assert.Equal(s.One, one)
	// TODO: check why one.p != _one.p

	one, err = s.ZeroOneSeg.Pow(s.Zero)

	assert.Nil(err)
	assert.Equal(s.One, one)

	one, err = s.MinusTenTwoAndHalfSeg.Pow(s.Zero)
	assert.Nil(err)
	assert.Equal(s.One, one)

	one, err = s.Inf.Pow(s.Zero)
	assert.Nil(err)
	assert.Equal(s.One, one)
}

func (s *NumberSuite) TestPowOneEdgeCases() {
	assert := assert.New(s.T())

	two, err := s.Two.Pow(s.One)

	assert.Nil(err)
	assert.Equal(s.Two, two)
	assert.True(s.Two.p == two.p)

	zero_one_seg_open, err := s.ZeroOneSegOpen.Pow(s.One)

	assert.Nil(err)
	assert.Equal(s.ZeroOneSegOpen, zero_one_seg_open)
	assert.True(s.ZeroOneSegOpen.p == zero_one_seg_open.p)

	minus_10_to_2_5_seg, err := s.MinusTenTwoAndHalfSeg.Pow(s.One)
	assert.Nil(err)
	assert.Equal(s.MinusTenTwoAndHalfSeg, minus_10_to_2_5_seg)
	assert.True(s.MinusTenTwoAndHalfSeg.p == minus_10_to_2_5_seg.p)

	inf, err := s.Inf.Pow(s.One)
	assert.Nil(err)
	assert.Equal(s.Inf, inf)
	assert.True(s.Inf.p == inf.p)
}

/*
func (s *NumberSuite) TestPowMinusOneEdgeCases() {
	assert := assert.New(s.T())

	_0_5, err := s.Two.Pow(s.MinusOne)

	assert.Nil(err)
	assert.Equal(NewNumberConst(0.5), _0_5)

	_0_2_to_1_seg, err := s.OneFiveSeg.Pow(s.MinusOne)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(0.2, 1), _0_2_to_1_seg)

	_1_to_inf_seg, err := s.ZeroOneSeg.Pow(s.MinusOne)

	assert.Nil(err)
	assert.Equal(NewNumberSegment(1, math.Inf(1)), _1_to_inf_seg)

	_1_open_inf_seg, err := s.ZeroOneSegOpen.Pow(s.MinusOne)

	assert.Nil(err)
	assert.Equal(NewNumberRange(&NRange{
		lVal: 1, lIncluding: false,
		rVal: math.Inf(1), rIncluding: true,
	}), _1_open_inf_seg)

	_minus_inf_to_minus_1_seg_or_1_to_inf_seg, err := NewNumberSegment(-1, 1).Pow(s.MinusOne)

	_exp_minus_1_to_1_pow_minus_1 := NewNumberRange(&NRange{
		lVal: math.Inf(-1), lIncluding: true,
		rVal: -1, rIncluding: true,
	}, &NRange{
		lVal: 1, lIncluding: true,
		rVal: math.Inf(1), rIncluding: true,
	})

	assert.Nil(err)
	assert.Equal(_exp_minus_1_to_1_pow_minus_1, _minus_inf_to_minus_1_seg_or_1_to_inf_seg)
	// TODO: support multirange numbers

	// Suppres unused
	assert.Nil(_minus_inf_to_minus_1_seg_or_1_to_inf_seg.p)
}
*/
func (s *NumberSuite) TestIDivSimple() {
	assert := assert.New(s.T())

	two, err := s.Five.IDiv(s.Two)

	assert.Nil(err)
	assert.Equal(s.Two, two)

	minus_two, err := s.Five.Negate().IDiv(s.Two)

	assert.Nil(err)
	assert.Equal(NewNumberConst(-2), minus_two)

	unknown_int, err := s.Unknown.IDiv(s.Unknown_other)

	unknownInt := NewNumber()
	unknownInt.p.integer = NewBooleanConst(BTrue, nil)

	assert.Nil(err)
	assert.Equal(unknownInt, unknown_int)

	one, err := s.Unknown.IDiv(s.Unknown_copy)

	assert.Nil(err)
	assert.Equal(s.One, one)
}

func (s *NumberSuite) TestIDivRange() {
	assert := assert.New(s.T())

	_0_to_2_seg_int, err := s.OneFiveSeg.IDiv(s.Two)

	zeroTwoSegInt := NewNumberSegment(0, 2)
	zeroTwoSegInt.p.integer = NewBooleanConst(BTrue, nil)

	assert.Nil(err)
	assert.Equal(zeroTwoSegInt, _0_to_2_seg_int)

	_0_to_3_seg_int, err := s.TwoFourSegOpen.IDiv(s.OneFiveSeg)

	zeroThreeSegInt := NewNumberSegment(0, 3)
	zeroThreeSegInt.p.integer = NewBooleanConst(BTrue, nil)

	assert.Nil(err)
	assert.Equal(zeroThreeSegInt, _0_to_3_seg_int)

	unknown_int, err := s.Two.IDiv(s.MinusTenTenSeg)

	unknownInt := NewNumber()
	unknownInt.p.integer = NewBooleanConst(BTrue, nil)

	assert.Nil(err)
	assert.Equal(unknownInt, unknown_int)
}

func (s *NumberSuite) TestIDivInfEdgeCases() {
	assert := assert.New(s.T())

	nan, err := s.Inf.IDiv(s.Inf)

	assert.Nil(err)
	assert.True(nan.p.integer.IsUnknown())

	// we cannot assert.Equal due to NaN != NaN
	assert.True(NewNumberConst(math.NaN()).IsSame(nan))

	nan, err = s.Inf.IDiv(s.Inf.Negate())

	assert.Nil(err)
	assert.True(nan.p.integer.IsUnknown())

	assert.True(NewNumberConst(math.NaN()).IsSame(nan))
}

func (s *NumberSuite) TestSplitSimple() {
	assert := assert.New(s.T())

	l, r := s.One.Split(0.9999)

	assert.Equal(s.One, r)
	assert.True(s.One.p == r.p)

	assert.False(l.IsValid())

	l, r = s.One.Split(2)

	assert.Equal(s.One, l)
	assert.True(s.One.p == l.p)

	assert.False(r.IsValid())

	l, r = s.One.Split(1)

	assert.Equal(s.One, r)
	assert.True(s.One.p == r.p)

	assert.False(l.IsValid())
}

func (s *NumberSuite) TestSplitRangeSimple() {
	assert := assert.New(s.T())

	l, r := s.OneFiveSeg.Split(2)

	assert.Equal(NewNumberRange(&NRange{
		lVal: 1, lIncluding: true,
		rVal: 2, rIncluding: false,
	}), l)
	assert.Equal(NewNumberSegment(2, 5), r)

	l, r = s.OneFiveSeg.Split(1)

	assert.Equal(s.OneFiveSeg, r)
	assert.True(s.OneFiveSeg.p == r.p)

	assert.False(l.IsValid())

	l, r = s.OneFiveSeg.Split(5)

	assert.Equal(NewNumberRange(&NRange{
		lVal: 1, lIncluding: true,
		rVal: 5, rIncluding: false,
	}), l)
	assert.Equal(s.Five, r)

	l, r = s.OneFiveSeg.Split(0)

	assert.Equal(s.OneFiveSeg, r)
	assert.True(s.OneFiveSeg.p == r.p)

	assert.False(l.IsValid())

	l, r = s.OneFiveSeg.Split(100)

	assert.Equal(s.OneFiveSeg, l)
	assert.True(s.OneFiveSeg.p == l.p)

	assert.False(r.IsValid())

	l, r = s.ZeroOneSegOpen.Split(0.5)

	assert.Equal(NewNumberRange(&NRange{
		lVal: 0.0, lIncluding: true,
		rVal: 0.5, rIncluding: false,
	}), l)
	assert.Equal(NewNumberRange(&NRange{
		lVal: 0.5, lIncluding: true,
		rVal: 1.0, rIncluding: false,
	}), r)

	l, r = s.ZeroOneSegOpen.Split(0)

	assert.Equal(s.ZeroOneSegOpen, r)
	assert.True(s.ZeroOneSegOpen.p == r.p)

	assert.False(l.IsValid())

	l, r = s.ZeroOneSegOpen.Split(1)

	assert.Equal(s.ZeroOneSegOpen, l)
	assert.True(s.ZeroOneSegOpen.p == l.p)

	assert.False(r.IsValid())
}

func (s *NumberSuite) TestSplitRangeInt() {
	assert := assert.New(s.T())

	l, r := s.OneFiveSegInt.Split(2)

	_2_to_5_seg_int := NewNumberSegment(2, 5)
	_2_to_5_seg_int.p.integer = NewBooleanConst(BTrue, nil)

	assert.Equal(NewNumberConst(1), l)
	assert.True(s.OneFiveSegInt.p.integer.p == l.p.integer.p)

	assert.Equal(_2_to_5_seg_int, r)
	assert.True(s.OneFiveSegInt.p.integer.p == r.p.integer.p)

	l, r = s.OneFiveSegInt.Split(1)

	assert.Equal(s.OneFiveSegInt, r)
	assert.True(s.OneFiveSegInt.p == r.p)

	assert.False(l.IsValid())

	l, r = s.OneFiveSegInt.Split(5)

	_1_to_4_seg_int := NewNumberSegment(1, 4)
	_1_to_4_seg_int.p.integer = NewBooleanConst(BTrue, nil)

	assert.Equal(_1_to_4_seg_int, l)
	assert.True(s.OneFiveSegInt.p.integer.p == l.p.integer.p)

	assert.Equal(NewNumberConst(5), r)
	assert.True(s.OneFiveSegInt.p.integer.p == r.p.integer.p)
}

func (s *NumberSuite) TestSplitRangeNotInt() {
	assert := assert.New(s.T())

	oneFiveOpenSegNotInt := Number{p: s.OneFiveSeg.p.Clone()}
	oneFiveOpenSegNotInt.p.integer = NewBooleanConst(BFalse, nil)
	oneFiveOpenSegNotInt, err := oneFiveOpenSegNotInt.RangeAdjust()

	assert.Nil(err)

	l, r := s.OneFiveSegInt.Split(2)

	_1_to_2_seg_not_int := NewNumberSegment(1, 2)
	_1_to_2_seg_not_int.p.integer = NewBooleanConst(BFalse, nil)
	_1_to_2_seg_not_int, err = _1_to_2_seg_not_int.RangeAdjust()

	assert.Nil(err)

	_2_to_5_seg_not_int := NewNumberSegment(2, 5)
	_2_to_5_seg_not_int.p.integer = NewBooleanConst(BFalse, nil)
	_2_to_5_seg_not_int, err = _2_to_5_seg_not_int.RangeAdjust()

	assert.Equal(_1_to_2_seg_not_int, l)
	assert.True(s.OneFiveSegInt.p.integer.p == l.p.integer.p)

	assert.Equal(_2_to_5_seg_not_int, r)
	assert.True(s.OneFiveSegInt.p.integer.p == r.p.integer.p)

	/*l, r = s.OneFiveSegInt.Split(1)

	assert.Equal(s.OneFiveSegInt, r)
	assert.True(s.OneFiveSegInt.p == r.p)
	assert.False(l.IsValid())

	l, r = s.OneFiveSegInt.Split(5)

	_1_to_4_int := NewNumberSegment(1, 4)
	_1_to_4_int.p.integer = NewBooleanConst(BTrue, nil)

	assert.Equal(_1_to_4_int, l)
	assert.True(s.OneFiveSegInt.p.integer.p == l.p.integer.p)

	assert.Equal(NewNumberConst(5), r)
	assert.True(s.OneFiveSegInt.p.integer.p == r.p.integer.p)*/
}

func (s *NumberSuite) TestRangeAdjust() {
	assert := assert.New(s.T())

	one, err := NewNumberSegment(1, 1).RangeAdjust()

	assert.Nil(err)
	assert.Equal(s.One, one)

	one_two_seg_not_int := NewNumberSegment(1, 2)
	one_two_seg_not_int.p.integer = NewBooleanConst(BFalse, nil)
	one_two_seg_not_int, err = one_two_seg_not_int.RangeAdjust()

	expected := NewNumberRange(&NRange{
		lVal: 1, lIncluding: false,
		rVal: 2, rIncluding: false,
	})
	expected.p.integer = NewBooleanConst(BFalse, nil)

	assert.Nil(err)
	assert.Equal(expected, one_two_seg_not_int)

	one_two_seg_int := NewNumberSegment(1, 2)
	one_two_seg_int.p.integer = NewBooleanConst(BTrue, nil)
	one_two_seg_int, err = one_two_seg_int.RangeAdjust()

	expected = NewNumberRange(&NRange{
		lVal: 1, lIncluding: true,
		rVal: 2, rIncluding: true,
	})
	expected.p.integer = NewBooleanConst(BTrue, nil)

	assert.Nil(err)
	assert.Equal(expected, one_two_seg_int)

	invalid_seg_int := NewNumberSegment(1.5, 1.9)
	invalid_seg_int.p.integer = NewBooleanConst(BTrue, nil)
	invalid_seg_int, err = invalid_seg_int.RangeAdjust()

	assert.NotNil(err)
	assert.Equal("no integer representation for NRange", err.Error())
	assert.False(invalid_seg_int.IsValid())
}

func (s *NumberSuite) TestAbsSimple() {
	assert := assert.New(s.T())

	two := s.Two.Abs()

	assert.Equal(s.Two, two)
	assert.True(s.Two.p == two.p)

	one := s.MinusOne.Abs()

	assert.Equal(s.One, one)
}

func (s *NumberSuite) TestAbsRange() {
	assert := assert.New(s.T())

	one_five_seg := s.OneFiveSeg.Abs()

	assert.Equal(s.OneFiveSeg, one_five_seg)
	assert.True(s.OneFiveSeg.p == one_five_seg.p)

	zero_ten_seg := s.MinusTenTwoAndHalfSeg.Abs()

	expected := NewNumberSegment(0, 10)
	expected.p.constraints = []NumberConstraint{
		NewNumberGreaterEqual(s.MinusTenTwoAndHalfSeg),
	}

	assert.Equal(expected, zero_ten_seg)

	assert.True(zero_ten_seg.Equal(s.MinusTenTwoAndHalfSeg).IsUnknown())
	assert.True(zero_ten_seg.Greater(s.MinusTenTwoAndHalfSeg).IsUnknown())
	assert.True(zero_ten_seg.GreaterEqual(s.MinusTenTwoAndHalfSeg).IsTrue())
	assert.True(zero_ten_seg.Less(s.MinusTenTwoAndHalfSeg).IsFalse())

	non_negative := s.Unknown.Abs()

	expected = NewNumberSegment(0, math.Inf(1))
	expected.p.constraints = []NumberConstraint{
		NewNumberGreaterEqual(s.Unknown),
	}

	assert.Equal(expected, non_negative)

	assert.True(non_negative.Equal(s.Unknown).IsUnknown())
	assert.True(non_negative.Greater(s.Unknown).IsUnknown())
	assert.True(non_negative.GreaterEqual(s.Unknown).IsTrue())
	assert.True(non_negative.Less(s.Unknown).IsFalse())

	positive := s.Negative.Abs()

	expected = NewNumberRange(&NRange{
		lVal: 0, lIncluding: false,
		rVal: math.Inf(1), rIncluding: true,
	})
	expected.p.constraints = []NumberConstraint{
		NewNumberGreater(s.Negative),
	}

	assert.Equal(expected, positive)

	assert.True(positive.Equal(s.Negative).IsFalse())
	assert.True(positive.Greater(s.Negative).IsTrue())
	assert.True(positive.GreaterEqual(s.Negative).IsTrue())
	assert.True(positive.Less(s.Negative).IsFalse())
}

func TestNumber(t *testing.T) {
	suite.Run(t, new(NumberSuite))
}
