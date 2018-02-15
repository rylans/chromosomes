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
