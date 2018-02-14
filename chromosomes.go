// Package chromosomes provides functionality for emulating an arbitrary genome
// and producing offspring from two parents.
package chromosomes

import (
  "math/rand"
  "math"
  "time"
)

// type Chromosome is a mapping between trait keys and their integer values
type Chromosome struct {
  traitKeys []string
  traits map[string]uint8
}

// Get returns the uint8 value of the given trait
func (c *Chromosome) Get(trait string) (uint8, bool) {
  v, exists := c.traits[trait]
  return v, exists
}

// Crossover takes two chromosomes and produces a random child chromosome based on the parents' genome
//
// Clients can control randomness by calling rand.Seed(x)
func (c *Chromosome) Crossover(other *Chromosome) *Chromosome {
  tmap := make(map[string]uint8, 0)

  for _, k := range c.traitKeys {
    theseBits := c.traits[k]
    thisMask := uint8(rand.Intn(256))

    thoseBits := other.traits[k]
    thatMask := uint8(255 - thisMask)

    tmap[k] = (theseBits & thisMask) | (thoseBits & thatMask)
  }

  return &Chromosome{traitKeys: c.traitKeys, traits: tmap}
}


type ChromosomeBuilder struct {
  traitKeys []string
  traits map[string]uint8
}

// NewBuilder creates a new empty ChromosomeBuilder
func NewBuilder() *ChromosomeBuilder {
  return &ChromosomeBuilder{traitKeys: make([]string, 0),
    traits: make(map[string]uint8, 0)}
}

// Add a certain trait to this ChromosomeBuilder
//
// A trait represents a particular gene.
// This function will panic if two of the same traits are added
func (builder *ChromosomeBuilder) AddTrait(trait string) {
  if _, duplicate := builder.traits[trait]; duplicate {
    panic("Duplicate trait: " + trait)
  }

  builder.traitKeys = append(builder.traitKeys, trait)
  builder.traits[trait] = 1
}

// BuildRandom creates a random Chromosome from this builder
//
// Clients can control randomness by calling rand.Seed(x)
func (builder *ChromosomeBuilder) BuildRandom() *Chromosome {
  builderTraitKeys := builder.traitKeys

  ckeys := make([]string, 0)
  for _, k := range builderTraitKeys {
    ckeys = append(ckeys, k)
  }

  traitmap := make(map[string]uint8, 0)
  for _, k := range ckeys {
    traitmap[k] = uint8(rand.Intn(256))
  }

  return &Chromosome{traitKeys: ckeys, traits: traitmap}
}

// MostFit returns the chromosome in the list of candidates that yields the max value when the fitness function fn is applied to it
//
// This function panics if the list of candidates is empty
func MostFit(fn func(x *Chromosome) float64, candidates ...*Chromosome) *Chromosome {
  topfitness := -math.MaxFloat64
  topcandidate := candidates[0]

  for _, c := range candidates {
    fitness := fn(c)
    if fitness > topfitness {
      topfitness = fitness
      topcandidate = c
    }
  }
  return topcandidate
}

func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}
