package chromosomes

import (
  "math/rand"
  "time"
)

type Chromosome struct {
  traitKeys []string
  traits map[string]uint8
}

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


type chromosomeBuilder struct {
  traitKeys []string
  traits map[string]uint8
}

func NewBuilder() *chromosomeBuilder {
  return &chromosomeBuilder{traitKeys: make([]string, 0),
    traits: make(map[string]uint8, 0)}
}

func (builder *chromosomeBuilder) AddTrait(trait string) {
  if _, duplicate := builder.traits[trait]; duplicate {
    panic("Duplicate trait: " + trait)
  }

  builder.traitKeys = append(builder.traitKeys, trait)
  builder.traits[trait] = 1
}

func (builder *chromosomeBuilder) BuildRandom() *Chromosome {
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

func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}
