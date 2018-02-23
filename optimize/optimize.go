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

  result := Optimize(boundedFitnessFunc1d(min, max, f), b)
  x = rescale(result.Get("X"), min, max)
  return
}

// Maximize a two-dimensional real-valued function
func BoundedMaximize2D(f RealFunction2D, min, max int) (x,y float64) {
  b := chromosomes.NewBuilder()
  b.AddTrait("X")
  b.AddTrait("Y")

  result := Optimize(boundedFitnessFunc2d(min, max, f), b)
  x = rescale(result.Get("X"), min, max)
  y = rescale(result.Get("Y"), min, max)
  return
}

// Maximize a three-dimensional real-valued function
func BoundedMaximize3D(f RealFunction3D, min, max int) (x,y,z float64) {
  b := chromosomes.NewBuilder()
  b.AddTrait("X")
  b.AddTrait("Y")
  b.AddTrait("Z")

  result := Optimize(boundedFitnessFunc3d(min, max, f), b)
  x = rescale(result.Get("X"), min, max)
  y = rescale(result.Get("Y"), min, max)
  z = rescale(result.Get("Z"), min, max)
  return
}

// Attempt to maximize the fitness function
func Optimize(fitness FitnessFn, cb *chromosomes.ChromosomeBuilder) *chromosomes.Chromosome {
  return stepOptimize(fitness, cb, 40)
}

func avgFitness(fitness FitnessFn, cs []*chromosomes.Chromosome) float64 {
  total := 0.0
  for _, c := range cs {
    total += fitness(c)
  }
  return total/float64(len(cs))
}

func aboveAverage(fitness FitnessFn, cs []*chromosomes.Chromosome) []*chromosomes.Chromosome {
  avg := avgFitness(fitness, cs)
  better := make([]*chromosomes.Chromosome, 0)
  for _, c := range cs {
    if fitness(c) > avg {
      better = append(better, c)
    }
  }
  return better
}

func stepOptimize(fitness FitnessFn, cb *chromosomes.ChromosomeBuilder, steps int) *chromosomes.Chromosome {
  poolsize := 8
  maxPoolSize := 233
  pool := make([]*chromosomes.Chromosome, 0, poolsize)

  for i := 0; i < poolsize; i++ {
    pool = append(pool, cb.BuildRandom())
  }

  for i := 0; i < steps; i++ {
    pool[0] = mostFit(fitness, pool...)
    if len(pool) > maxPoolSize {
      pool = pool[:maxPoolSize/2]
    }

    pool = append(pool, cb.BuildRandom())
    pool = append(pool, cb.BuildRandom())

    topQuarter := aboveAverage(fitness, aboveAverage(fitness, pool))
    if len(topQuarter) > poolsize {
      topQuarter = topQuarter[:poolsize]
    }

    for i := 1; i < len(topQuarter); i++ {
      for j := i; j < len(topQuarter); j++ {
	if topQuarter[i].Difference(topQuarter[j]) > 1 {
	  pool = append(pool, topQuarter[i].Crossover(topQuarter[j]))
	}
      }
      pool[i] = topQuarter[i]
      pool = append(pool, topQuarter[i].Crossover(cb.BuildRandom()))
    }
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
