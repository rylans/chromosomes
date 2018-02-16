package chromosomes

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "math/rand"
)

func TestChromosomeBuilding(t *testing.T){
  rand.Seed(1401)

  b := NewBuilder()
  b.AddTrait("abc")
  
  c1 := b.BuildRandom()

  val, exists := c1.Get("abc")
  assert.Equal(t, true, exists)
  assert.Equal(t, uint8(0xc6), val)

  noval, nope := c1.Get("anything-else")
  assert.Equal(t, false, nope)
  assert.Equal(t, uint8(0x0), noval)
}

func TestCrossover(t *testing.T){
  rand.Seed(1401)

  b := NewBuilder()
  b.AddTrait("a")
  b.AddTrait("b")
  b.AddTrait("c")

  parent1 := b.BuildRandom()
  parent2 := b.BuildRandom()

  child := parent1.Crossover(parent2)

  valA, _ := child.Get("a")
  valB, _ := child.Get("b")
  valC, _ := child.Get("c")
  assert.Equal(t, uint8(0xb6), valA)
  assert.Equal(t, uint8(0xe7), valB)
  assert.Equal(t, uint8(0x66), valC)
}

func TestFitnessEvaluation(t *testing.T){
  rand.Seed(1401)

  fn := func(x *Chromosome) float64 {
    v, _ := x.Get("X")
    return float64(v)
  }

  b := NewBuilder()
  b.AddTrait("X")

  c1 := b.BuildRandom()
  c2 := b.BuildRandom()
  c3 := b.BuildRandom()

  winner := MostFit(fn, c1, c2,c3)

  winval, _ := winner.Get("X")
  assert.Equal(t, uint8(0xff), winval)
}

func TestCrossoverWithMutation(t *testing.T){
  rand.Seed(1401)

  b := NewBuilder()
  b.MutationChance(1.0)
  b.AddTrait("a")

  parent1 := b.BuildRandom()
  parent2 := b.BuildRandom()

  child := parent1.Crossover(parent2)
  valA, _ := child.Get("a")
  assert.Equal(t, uint8(0x28), valA)
}

func TestSetMutationPanics(t *testing.T) {
  b := NewBuilder()

  assert.Panics(t,
    func() { b.MutationChance(1.01) }, "Expected panic")
}

func TestChromosomeDifferenceSymmetric(t *testing.T){
  b := NewBuilder()
  b.AddTrait("t1")
  b.AddTrait("t2")

  parent1 := b.BuildRandom()
  parent2 := b.BuildRandom()

  assert.Equal(t, parent1.Difference(parent2), parent2.Difference(parent1))
}

func TestChromosomeDifferencePanics(t *testing.T){
  build1 := NewBuilder()
  build1.AddTrait("foo")

  build2 := NewBuilder()
  build2.AddTrait("bar")

  assert.Panics(t,
    func() { build1.BuildRandom().Difference(build2.BuildRandom()) },
    "Did not panic")
}

func TestEmptyChromosomeDifferenceZero(t *testing.T){
    b := NewBuilder()

    c1 := b.BuildRandom()
    c2 := b.BuildRandom()

    assert.Equal(t, 0, c2.Difference(c1))
}
