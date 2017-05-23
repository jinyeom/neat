// statistics.go implementation of statistical information of the evolution.
//
// Copyright (C) 2017  Jin Yeom
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package neat

import (
	"math"
)

// Statistics is a data structure that records statistical information of each
// generation during the evolutionary process.
type Statistics struct {
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
	s.MinFitness[currGen] = n.Population[0].Fitness
	s.MaxFitness[currGen] = n.Population[0].Fitness
	for _, genome := range n.Population {
		s.MinFitness[currGen] = math.Min(genome.Fitness, s.MinFitness[currGen])
		s.MaxFitness[currGen] = math.Max(genome.Fitness, s.MinFitness[currGen])
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
