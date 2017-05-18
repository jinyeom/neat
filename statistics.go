package neat

import (
	"math"
)

// Statistics is a data structure that records statistical information of each
// generation during the evolutionary process.
type Statitstics struct {
	NumSpecies []int     // number of species in each generation
	MinFitness []float64 // minimum fitness in each generation
	MaxFitness []float64 // maximum fitness in each generation
	AvgFitness []float64 // average fitness in each generation
}

// NewStatistics returns a new instance of Statistics.
func NewStatistics(numGenerations int) *Statistics {
	return &Statistics{
		NumSpecies: make([]int, numGenerations),
		MinFitness: make([]float64, numGenerations),
		MaxFitness: make([]float64, numGenerations),
		AvgFitness: make([]float64, numGenerations),
	}
}

// Update the statistics of current generation
func (s *Statistics) Update(currGen int, n *NEAT) {
	s.NumSpecies[currGen] = len(n.Species)

	// mininum and maximum
	MinFitness[currGen] = n.Population[0].Fitness
	MaxFitness[currGen] = n.Population[0].Fitness
	for _, genome := range n.Population {
		MinFitness[currGen] = math.Min(genome.Fitness, MinFitness[currGen])
		MaxFitness[currGen] = math.Max(genome.Fitness, MinFitness[currGen])
	}

	// average fitness
	s.AvgFitness[currGen] = func() float64 {
		avg := 0.0
		for _, genome := range n.Population {
			avg += genome.Fitness
		}
		return avg / float64(n.Config.PopulationSize)
	}()
}

func (s *Statistics) ExportChart() {

}
