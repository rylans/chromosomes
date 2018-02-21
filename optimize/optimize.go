package optimize

import (
  "github.com/rylans/chromosomes"
  "math"
)

type FitnessFn func(c *chromosomes.Chromosome) float64

// Attempt to maximize the fitness function
func Optimize(fitness FitnessFn, cb *chromosomes.ChromosomeBuilder) *chromosomes.Chromosome {
  return stepOptimize(fitness, cb, 24)
}

func stepOptimize(fitness FitnessFn, cb *chromosomes.ChromosomeBuilder, steps int) *chromosomes.Chromosome {
  poolsize := 16
  maxPoolSize := 80
  pool := make([]*chromosomes.Chromosome, 0, poolsize)

  for i := 0; i < poolsize; i++ {
    pool = append(pool, cb.BuildRandom())
  }

  for i := 0; i < steps; i++ {
    elite := mostFit(fitness, pool...)
    if len(pool) > maxPoolSize {
      pool = pool[:maxPoolSize]
    }

    poolLength := len(pool)
    pool = append(pool, elite.Crossover(pool[poolLength-1]))
    pool = append(pool, elite.Crossover(pool[poolLength-3]))
    pool = append(pool, elite.Crossover(pool[poolLength-5]))
    pool = append(pool, cb.BuildRandom())
    pool[4] = elite.Crossover(cb.BuildRandom())
    pool[1] = elite.Crossover(cb.BuildRandom())
    pool[0] = elite
  }
  
  return mostFit(fitness, pool...)
}

// mostFit returns the chromosome in the list of candidates that yields the max value when the fitness function fn is applied to it
//
// This function panics if the list of candidates is empty
func mostFit(fn func(x *chromosomes.Chromosome) float64, candidates ...*chromosomes.Chromosome) *chromosomes.Chromosome {
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
