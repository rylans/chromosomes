package chromosomes

type Chromosome struct {
}

type chromosomeBuilder struct {
  traitKeys []string
}

func NewBuilder() *chromosomeBuilder {
  return &chromosomeBuilder{traitKeys: make([]string, 0)}
}
