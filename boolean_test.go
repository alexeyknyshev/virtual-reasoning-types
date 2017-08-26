package virtual_types

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BooleanSuite struct {
	suite.Suite
	True           Boolean
	False          Boolean
	Unknown_1      Boolean
	Unknown_2      Boolean
	Unknown_1_same Boolean
	Unknown_1_copy Boolean
}

func (s *BooleanSuite) SetupTest() {
	s.True = NewBooleanConst(BTrue, nil)
	s.False = NewBooleanConst(BFalse, nil)
	s.Unknown_1 = NewBooleanConst(BUnknown, nil)
	s.Unknown_2 = NewBoolean()
	s.Unknown_1_same = NewBooleanConst(BUnknown, []Constraint{
		NewBooleanEqual(s.Unknown_1),
	})
	s.Unknown_1_copy = Boolean{p: s.Unknown_1.p}
}

func (s *BooleanSuite) TestIsTrue() {
	assert := assert.New(s.T())

	assert.Equal(BTrue, s.True.p.val, "BTrue val expected")
	assert.Equal(BFalse, s.False.p.val, "BFalse val expected")
	assert.Equal(BUnknown, s.Unknown_1.p.val, "BUnknown val expected")
	assert.Equal(BUnknown, s.Unknown_2.p.val, "BUnknown val expected")
	assert.Equal(BUnknown, s.Unknown_1_same.p.val, "BUnknown val expected")

	assert.True(s.True.IsTrue(), "True.IsTrue() should be true")
	assert.False(s.False.IsTrue(), "False.IsTrue() should be false")
	assert.False(s.Unknown_1.IsTrue(), "Unknown_1.IsTrue() should be false")
	assert.False(s.Unknown_2.IsTrue(), "Unknown_2.IsTrue() should be false")
	assert.False(s.Unknown_1_same.IsTrue(), "Unknown_2.IsTrue() should be false")
}

func (s *BooleanSuite) TestIsConstant() {
	assert := assert.New(s.T())

	assert.True(s.True.IsConstant(), "True.IsConstant() should be true")
	assert.True(s.False.IsConstant(), "False.IsConstant() should be true")
	assert.False(s.Unknown_1.IsConstant(), "Unknown_1.IsConstant() should be false")
	assert.False(s.Unknown_2.IsConstant(), "Unknown_2.IsConstant() should be false")
	assert.False(s.Unknown_1_same.IsConstant(), "Unknown_2.IsConstant() should be false")
}

func (s *BooleanSuite) TestIsUnknown() {
	assert := assert.New(s.T())

	assert.False(s.True.IsUnknown(), "")
	assert.False(s.False.IsUnknown(), "")
	assert.True(s.Unknown_1.IsUnknown(), "")
	assert.True(s.Unknown_2.IsUnknown(), "")
	assert.True(s.Unknown_1_same.IsUnknown(), "")
}

func (s *BooleanSuite) TestNot() {
	assert := assert.New(s.T())

	assert.False(s.True.Not().IsTrue(), "")
	assert.True(s.False.Not().IsTrue(), "")
	assert.True(s.Unknown_1.Not().IsUnknown(), "")
	assert.True(s.Unknown_2.Not().IsUnknown(), "")
	assert.True(s.Unknown_1_same.IsUnknown(), "")

	assert.True(s.True.Not().Not().IsTrue(), "")
	assert.False(s.False.Not().Not().IsTrue(), "")
}

func (s *BooleanSuite) TestAnd() {
	assert := assert.New(s.T())

	/*	assert.True(s.True.And(s.True).IsTrue(), "")
		assert.True(s.True.And(s.False).IsFalse(), "")
		assert.True(s.True.And(s.Unknown_1).IsUnknown(), "")
		assert.True(s.True.And(s.Unknown_1_same).IsUnknown(), "")

		assert.True(s.False.And(s.True).IsFalse(), "")
		assert.True(s.False.And(s.False).IsFalse(), "")
		assert.True(s.False.And(s.Unknown_1).IsFalse(), "")
		assert.True(s.False.And(s.Unknown_1_copy).IsFalse(), "")*/

	//	assert.True(s.Unknown_1.And(s.True).IsUnknown(), "")
	assert.True(s.Unknown_1.And(s.False).IsFalse(), "")
	assert.True(s.Unknown_1.And(s.Unknown_1).IsUnknown(), "")
	assert.True(s.Unknown_1.And(s.Unknown_1_copy).IsUnknown(), "")
}

func (s *BooleanSuite) TestOr() {
	assert := assert.New(s.T())

	assert.True(s.True.Or(s.True).IsTrue(), "")
	assert.True(s.True.Or(s.False).IsTrue(), "")
	assert.True(s.True.Or(s.Unknown_1).IsTrue(), "")
	assert.True(s.True.Or(s.Unknown_2).IsTrue(), "")
	assert.True(s.True.Or(s.Unknown_1_same).IsTrue(), "")

	assert.True(s.False.Or(s.True).IsTrue(), "")
	assert.True(s.False.Or(s.False).IsFalse(), "")
	assert.True(s.False.Or(s.Unknown_1).IsUnknown(), "")
	assert.True(s.False.Or(s.Unknown_2).IsUnknown(), "")
	assert.True(s.False.Or(s.Unknown_1_same).IsUnknown(), "")

	assert.True(s.Unknown_1.Or(s.True).IsTrue(), "")
	assert.True(s.Unknown_1.Or(s.False).IsUnknown(), "")
	assert.True(s.Unknown_1.Or(s.Unknown_1).IsUnknown(), "")
	assert.True(s.Unknown_1.Or(s.Unknown_2).IsUnknown(), "")
	assert.True(s.Unknown_1.Or(s.Unknown_1_same).IsUnknown(), "")
}

func (s *BooleanSuite) TestEqual() {
	assert := assert.New(s.T())

	assert.True(s.True.Equal(s.True).IsTrue(), "")
	assert.True(s.True.Equal(s.False).IsFalse(), "")
	assert.True(s.True.Equal(s.Unknown_1).IsUnknown(), "")
	assert.True(s.True.Equal(s.Unknown_2).IsUnknown(), "")
	assert.True(s.True.Equal(s.Unknown_1_same).IsUnknown(), "")
	assert.True(s.True.Equal(s.Unknown_1_copy).IsUnknown(), "")

	assert.True(s.False.Equal(s.True).IsFalse(), "")
	assert.True(s.False.Equal(s.False).IsTrue(), "")
	assert.True(s.False.Equal(s.Unknown_1).IsUnknown(), "")
	assert.True(s.False.Equal(s.Unknown_2).IsUnknown(), "")
	assert.True(s.False.Equal(s.Unknown_1_same).IsUnknown(), "")
	assert.True(s.False.Equal(s.Unknown_1_copy).IsUnknown(), "")

	assert.True(s.Unknown_1.Equal(s.True).IsUnknown(), "")
	assert.True(s.Unknown_1.Equal(s.False).IsUnknown(), "")
	assert.True(s.Unknown_1.Equal(s.Unknown_1).IsTrue(), "")
	assert.True(s.Unknown_1.Equal(s.Unknown_2).IsUnknown(), "")
	assert.True(s.Unknown_1.Equal(s.Unknown_1_same).IsTrue(), "")
	assert.True(s.Unknown_1.Equal(s.Unknown_1_copy).IsTrue(), "")

	assert.True(s.Unknown_2.Equal(s.True).IsUnknown(), "")
	assert.True(s.Unknown_2.Equal(s.False).IsUnknown(), "")
	assert.True(s.Unknown_2.Equal(s.Unknown_1).IsUnknown(), "")
	assert.True(s.Unknown_2.Equal(s.Unknown_2).IsTrue(), "")
	assert.True(s.Unknown_2.Equal(s.Unknown_1_same).IsUnknown(), "")
	assert.True(s.Unknown_2.Equal(s.Unknown_1_copy).IsUnknown(), "")

	assert.True(s.Unknown_1_same.Equal(s.True).IsUnknown(), "")
	assert.True(s.Unknown_1_same.Equal(s.False).IsUnknown(), "")
	assert.True(s.Unknown_1_same.Equal(s.Unknown_1).IsTrue(), "")
	assert.True(s.Unknown_1_same.Equal(s.Unknown_2).IsUnknown(), "")
	assert.True(s.Unknown_1_same.Equal(s.Unknown_1_same).IsTrue(), "")
	assert.True(s.Unknown_1_same.Equal(s.Unknown_1_copy).IsTrue(), "")
}

func (s *BooleanSuite) TestEqualTransitive() {
	assert := assert.New(s.T())

	{
		// b ==> a <== c
		a := NewBoolean()
		b := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(a)})
		c := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(a)})

		assert.True(a.Equal(b).IsTrue(), "")
		assert.True(a.Equal(c).IsTrue(), "")

		assert.True(b.Equal(a).IsTrue(), "")
		assert.True(c.Equal(a).IsTrue(), "")

		assert.True(b.Equal(c).IsTrue(), "")
		assert.True(c.Equal(b).IsTrue(), "")
	}

	{
		// b ==> a <!= c
		a := NewBoolean()
		b := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(a)})
		c := NewBooleanConst(BUnknown, []Constraint{NewBooleanNotEqual(a)})

		assert.True(a.Equal(b).IsTrue(), "")
		assert.True(b.Equal(a).IsTrue(), "")

		assert.True(a.Equal(c).IsFalse(), "")
		assert.True(c.Equal(a).IsFalse(), "")

	}

	{
		// [a <== c], [b]
		a := NewBoolean()
		b := NewBoolean()
		c := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(a)})

		assert.True(b.Equal(c).IsUnknown(), "")
		assert.True(c.Equal(b).IsUnknown(), "")
	}

	{
		// a <== b <== c ... <== g
		eqChain := []Boolean{NewBooleanConst(BUnknown, nil)}

		for i := 1; i < 10; i++ {
			eqChain = append(eqChain, NewBooleanConst(BUnknown, []Constraint{
				NewBooleanEqual(eqChain[i-1]),
			}))
		}

		for i := 0; i < len(eqChain)-1; i++ {
			for j := i + 1; j < len(eqChain); j++ {
				assert.True(eqChain[i].Equal(eqChain[j]).IsTrue(), "")
				assert.True(eqChain[j].Equal(eqChain[i]).IsTrue(), "")
			}
		}
	}

	{
		// d !=> b ==> a <!= c
		a := NewBoolean()
		b := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(a)})
		c := NewBooleanConst(BUnknown, []Constraint{NewBooleanNotEqual(a)})
		d := NewBooleanConst(BUnknown, []Constraint{NewBooleanNotEqual(b)})

		assert.True(a.Equal(b).IsTrue(), "")
		assert.True(b.Equal(a).IsTrue(), "")

		assert.True(c.Equal(a).IsFalse(), "")
		assert.True(a.Equal(c).IsFalse(), "")

		assert.True(b.Equal(c).IsFalse(), "")
		assert.True(c.Equal(b).IsFalse(), "")

		assert.True(d.Equal(b).IsFalse(), "")
		assert.True(b.Equal(d).IsFalse(), "")

		assert.True(d.Equal(a).IsFalse(), "")
		assert.True(a.Equal(d).IsFalse(), "")
	}

	{
		// a ==> b ==> c <== d <== e <!= f
		c := NewBoolean()
		b := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(c)})
		a := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(b)})
		d := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(c)})
		e := NewBooleanConst(BUnknown, []Constraint{NewBooleanEqual(d)})

		f := NewBooleanConst(BUnknown, []Constraint{NewBooleanNotEqual(e)})

		assert.True(a.Equal(e).IsTrue(), "")
		assert.True(e.Equal(a).IsTrue(), "")

		assert.True(a.Equal(f).IsFalse(), "")
		assert.True(f.Equal(a).IsFalse(), "")
	}
}

func TestBoolean(t *testing.T) {
	suite.Run(t, new(BooleanSuite))
}
