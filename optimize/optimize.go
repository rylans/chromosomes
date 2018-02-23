package optimize

import (
  "github.com/rylans/chromosomes"
  "math"
)

// A real-valued function of one variable
type RealFunction func(x float64) float64

// A real-valued function of two variables
type RealFunction2D func(x, y float64) float64

// A real-valued function of three variables
type RealFunction3D func(x, y, z float64) float64

type FitnessFn func(c *chromosomes.Chromosome) float64

// Maximize a one-dimensional real-valued function
func BoundedMaximize(f RealFunction, min, max int) (x float64) {
  b := chromosomes.NewBuilder()
  b.AddTrait("X")

  fitnessf := boundedFitnessFunc1d(min, max, f)
  result := Optimize(fitnessf, b)

  x = rescale(result.Get("X"), min, max)
  return
}

// Maximize a two-dimensional real-valued function
func BoundedMaximize2D(f RealFunction2D, min, max int) (x,y float64) {
  b := chromosomes.NewBuilder()
  b.AddTrait("X")
  b.AddTrait("Y")

  fitnessf := boundedFitnessFunc2d(min, max, f)
  result := Optimize(fitnessf, b)

  x = rescale(result.Get("X"), min, max)
  y = rescale(result.Get("Y"), min, max)
  return
}

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

func boundedFitnessFunc1d(min int, max int, optimizef RealFunction) FitnessFn {
	f := func(c *chromosomes.Chromosome) float64 {
		realVal := rescale(c.Get("X"), min, max)
		return optimizef(realVal)
	}
	return f
}

func boundedFitnessFunc2d(min int, max int, optimizef RealFunction2D) FitnessFn {
	f := func(c *chromosomes.Chromosome) float64 {
		realValX := rescale(c.Get("X"), min, max)
		realValY := rescale(c.Get("Y"), min, max)
		return optimizef(realValX, realValY)
	}
	return f
}

func boundedFitnessFunc3d(min int, max int, optimizef RealFunction3D) FitnessFn {
	f := func(c *chromosomes.Chromosome) float64 {
		realValX := rescale(c.Get("X"), min, max)
		realValY := rescale(c.Get("Y"), min, max)
		realValZ := rescale(c.Get("Z"), min, max)
		return optimizef(realValX, realValY, realValZ)
	}
	return f
}

func rescale(val uint8, min int, max int) float64 {
	newRange := float64(max) - float64(min)
	return (newRange/255.0)*(float64(val)-255.0) + float64(max)
}
