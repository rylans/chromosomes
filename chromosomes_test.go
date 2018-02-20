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

  val := c1.Get("abc")
  assert.Equal(t, uint8(0xc6), val)
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

  valA := child.Get("a")
  valB := child.Get("b")
  valC := child.Get("c")
  assert.Equal(t, uint8(0xb6), valA)
  assert.Equal(t, uint8(0xe7), valB)
  assert.Equal(t, uint8(0x66), valC)
}

func TestFitnessEvaluation(t *testing.T){
  rand.Seed(1401)

  fn := func(x *Chromosome) float64 {
    v := x.Get("X")
    return float64(v)
  }

  b := NewBuilder()
  b.AddTrait("X")

  c1 := b.BuildRandom()
  c2 := b.BuildRandom()
  c3 := b.BuildRandom()

  winner := MostFit(fn, c1, c2,c3)

  winval := winner.Get("X")
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
  valA := child.Get("a")
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

func TestLength(t *testing.T) {
  b := NewBuilder()

  cZero := b.BuildRandom()

  b.AddTrait("asdf")

  cEight := b.BuildRandom()

  assert.Equal(t, 0, cZero.Len())
  assert.Equal(t, 8, cEight.Len())
}

func TestGeneticDiversityDecreases(t *testing.T){
  b := NewBuilder()
  b.MutationChance(0)
  b.AddTrait("t1")
  b.AddTrait("t2")

  parent1 := b.BuildRandom()
  parent2 := b.BuildRandom()
  parentDifference := parent1.Difference(parent2)

  child1 := parent1.Crossover(parent2)
  child2 := parent1.Crossover(parent2)
  childDifference := child1.Difference(child2)

  assert.Equal(t, true, childDifference < parentDifference)
}

func TestCloneNoMutation(t *testing.T){
  b := NewBuilder()
  b.MutationChance(0.0)
  b.AddTrait("A")
  
  c1 := b.BuildRandom()
  c2 := c1.Clone()

  assert.Equal(t, c2.Get("A"), c1.Get("A"))
}

func TestCloneWithMutation(t *testing.T){
  b := NewBuilder()
  b.MutationChance(1)
  b.AddTrait("A")
  
  c1 := b.BuildRandom()
  c2 := c1.Clone()

  assert.Equal(t, true, c2.Get("A") != c1.Get("A"))
}

func TestCustomCrossoverTakesFirstParent(t *testing.T){
  b := NewBuilder()
  b.MutationChance(0)
  b.AddTrait("A")

  b.setCrossoverRule( func(c1 *Chromosome, c2 *Chromosome) *Chromosome {
    return c1
  })

  c1 := b.BuildRandom()
  c2 := b.BuildRandom()

  child := c1.Crossover(c2)

  assert.Equal(t, c1.Get("A"), child.Get("A"))
}

func TestCustomCrossoverTakesSecondParent(t *testing.T){
  b := NewBuilder()
  b.MutationChance(0)
  b.AddTrait("A")

  b.setCrossoverRule( func(c1 *Chromosome, c2 *Chromosome) *Chromosome {
    return c2
  })

  c1 := b.BuildRandom()
  c2 := b.BuildRandom()

  child := c1.Crossover(c2)

  assert.Equal(t, c2.Get("A"), child.Get("A"))
}
