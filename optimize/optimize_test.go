package optimize 

import (
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

func TestOptimizeImprovesChromosomeFromRandom(t *testing.T){
  b := chromosomes.NewBuilder()
  b.MutationChance(0.15)
  b.AddTrait(GENE1)
  b.AddTrait(GENE2)
  b.AddTrait(GENE3)

  random := b.BuildRandom()
  result := Optimize(mostOnesFitness, b)

  assert.Equal(t, true, mostOnesFitness(result) > mostOnesFitness(random))
}
