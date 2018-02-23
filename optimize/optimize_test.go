package optimize 

import (
  "math/bits"
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/rylans/chromosomes"
)

const GENE1 string = "gene1"
const GENE2 string = "gene2"
const GENE3 string = "gene3"

func mostOnesFitness(c *chromosomes.Chromosome) float64 {
  return float64(c.Get(GENE1)) + float64(c.Get(GENE2)) + float64(c.Get(GENE3))
}

func sumLeadingZeros(c *chromosomes.Chromosome) float64 {
  z1 := bits.LeadingZeros8(c.Get(GENE1))
  z2 := bits.LeadingZeros8(c.Get(GENE2))
  z3 := bits.LeadingZeros8(c.Get(GENE3))
  return float64(z1) + float64(z2) + float64(z3)
}

func TestOptimizeImprovesChromosomeFromRandom(t *testing.T){
  assertOptimizeImproves(t, mostOnesFitness, threeGeneBuilder())
}

func TestOptimizeImprovesChromosomeFromRandomLeadingZeros(t *testing.T){
  assertOptimizeImproves(t, sumLeadingZeros, threeGeneBuilder())
}

func assertOptimizeImproves(t *testing.T, f func(c *chromosomes.Chromosome) float64, b *chromosomes.ChromosomeBuilder) {
  random := b.BuildRandom()
  result := Optimize(f, b)

  assert.Equal(t, true, f(result) > f(random))
}

func threeGeneBuilder() *chromosomes.ChromosomeBuilder {
  b := chromosomes.NewBuilder()
  b.MutationChance(0.15)
  b.AddTrait(GENE1)
  b.AddTrait(GENE2)
  b.AddTrait(GENE3)
  return b
}

func TestMaximize1DStraightLine(t *testing.T){
  f := func(x float64) float64 {return x}

  x := BoundedMaximize(f, -10, 10)

  assert.Equal(t, true, x > 9.98)
}

func TestMaximize1DParabola(t *testing.T){
  f := func(x float64) float64 {
    return -1.0*(x*x)
  }

  x := BoundedMaximize(f, -10, 10)

  assert.Equal(t, true, x > -0.1)
  assert.Equal(t, true, x < 0.1)
}

func TestMaximize2DPlane(t *testing.T){
  f := func(x, y float64) float64 {return x*y}

  x, y := BoundedMaximize2D(f, -10, 10)

  assert.Equal(t, true, (x*y) > 80.0)
}
