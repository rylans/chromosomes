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
  assert.Equal(t, uint8(0xf7), valB)
  assert.Equal(t, uint8(0x2), valC)
}
