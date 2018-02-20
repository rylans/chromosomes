// Package chromosomes provides functionality for simulating a bitstring chromosome
// and producing offspring from two parents.
package chromosomes

import (
  "math/rand"
  "math/bits"
  "math"
  "time"
)

// Default chance of a mutation occuring on a per-trait basis
const DefaultMutationChance = 1e-5

type mutate func(input uint8) uint8

type crossoverRule func(c1 *Chromosome, c2 *Chromosome) *Chromosome

func perTraitMutator(chance float64) mutate {
  f := func(input uint8) uint8 {
    if rand.Float64() < chance {
      return input ^ uint8(rand.Intn(256))
    }
    return input
  }
  return f
}

// type Chromosome is a mapping between trait keys and their integer values
type Chromosome struct {
  traitKeys []string
  traits map[string]uint8
  mutator mutate
  crossoverRule
}

// Len returns the number of bits in this chromosome
func (c *Chromosome) Len() int {
  return 8 * len(c.traitKeys)
}

// Difference returns the number of bits that differ between c and other
//
// This function panics if the two chromosomes are incompatible with each other
func (c *Chromosome) Difference(other *Chromosome) int {
  diff := 0
  for _, k := range c.traitKeys {
    diff += bits.OnesCount8(c.Get(k) ^ other.Get(k))
  }
  return diff
}

// Get returns the uint8 value of the given trait
//
// This function panics if such trait does not exist
func (c *Chromosome) Get(trait string) uint8 {
  v, exists := c.traits[trait]
  if !exists {
    panic("Chromosome has no such trait: " + trait)
  }
  return v
}

// Clone creates a child chromosome from the given chromosome crossing over itself
//
// Mutations can occur during cloning
func (c *Chromosome) Clone() *Chromosome {
  return c.Crossover(c)
}

// Crossover takes two chromosomes and produces a random child chromosome based on the parents' genome
//
// Clients can control randomness by calling rand.Seed(x)
func (c *Chromosome) Crossover(other *Chromosome) *Chromosome {
  return c.crossoverRule(c, other)
}

func standardCrossover() crossoverRule {
  f := func (c *Chromosome, other *Chromosome) *Chromosome {
    tmap := make(map[string]uint8, len(c.traitKeys))

    for _, k := range c.traitKeys {
      thisMask := uint8(rand.Intn(256))
      thatMask := uint8(255 - thisMask)

      tmap[k] = (c.Get(k) & thisMask) | (other.Get(k) & thatMask)
      tmap[k] = c.mutator(tmap[k])
    }

    return &Chromosome{
	traitKeys: c.traitKeys, 
	traits: tmap, 
	mutator: c.mutator, 
	crossoverRule: c.crossoverRule,
    }
  }
  return f
}


type ChromosomeBuilder struct {
  traitKeys []string
  traits map[string]uint8
  mutator mutate
  crossoverRule
}

// NewBuilder creates a new empty ChromosomeBuilder
func NewBuilder() *ChromosomeBuilder {
  mutationf := perTraitMutator(DefaultMutationChance)
  return &ChromosomeBuilder{
    traitKeys: make([]string, 0),
    traits: make(map[string]uint8, 0),
    mutator: mutationf,
    crossoverRule: standardCrossover(),
  }
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

// Set the chance of a mutation occuring on a trait
//
// chance must be in the interval [0.0, 1.0] otherwise this function panics
func (builder *ChromosomeBuilder) MutationChance(chance float64) {
  if chance < 0.0 || chance > 1.0 {
    panic("Chance out of range (zero to one)")
  }
  mutationf := perTraitMutator(chance)
  builder.mutator = mutationf
}

// Set the function which produces a child chromosome from two parents during the Crossover call
func (builder *ChromosomeBuilder) setCrossoverRule(f crossoverRule) {
  builder.crossoverRule = f
}

// BuildRandom creates a random Chromosome from this builder
//
// Clients can control randomness by calling rand.Seed(x)
func (builder *ChromosomeBuilder) BuildRandom() *Chromosome {
  traits := len(builder.traitKeys)

  ckeys := make([]string, 0, traits)
  for _, k := range builder.traitKeys {
    ckeys = append(ckeys, k)
  }

  traitmap := make(map[string]uint8, traits)
  for _, k := range ckeys {
    traitmap[k] = uint8(rand.Intn(256))
  }

  return &Chromosome{
    traitKeys: ckeys, 
    traits: traitmap, 
    mutator: builder.mutator, 
    crossoverRule: builder.crossoverRule,
  }
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
