// neat.go implementation of NeuroEvolution of Augmenting Topologies (NEAT).
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
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// NEAT is the implementation of NeuroEvolution of Augmenting Topology (NEAT).
type NEAT struct {
	Config      *Config           // configuration
	Population  []*Genome         // population of genome
	Species     []*Species        // species of subpopulation of genomes
	Activations []*ActivationFunc // set of activation functions
	Evaluation  EvaluationFunc    // evaluation function
	Comparison  ComparisonFunc    // comparison function
	Best        *Genome           // best genome
	Statistics  *Statistics       // statistics

	nextGenomeID  int // genome ID that is assigned to a newly created genome
	nextSpeciesID int // species ID that is assigned to a newly created species
}

// New creates a new instance of NEAT with provided argument configuration and
// an evaluation function.
func New(config *Config, evaluation EvaluationFunc) *NEAT {
	nextGenomeID := 0
	nextSpeciesID := 0

	// in order to prevent containing multiple of the same activation function
	// in the set of activation functions, they will temporarily be added to a
	// map first, which contains Sigmoid function as a default, then be
	// transferred to a slice of ActivationFunc.
	temp := map[string]*ActivationFunc{
		"sigmoid": Sigmoid(),
	}

	// if more additional activation functions are needed,
	for _, name := range config.CPPNActivations {
		temp[name] = ActivationSet[name]
	}

	activations := make([]*ActivationFunc, 0, len(temp))
	for _, afunc := range temp {
		activations = append(activations, afunc)
	}

	population := make([]*Genome, config.PopulationSize)
	if config.FullyConnected {
		for i := 0; i < config.PopulationSize; i++ {
			population[i] = NewFCGenome(nextGenomeID, config.NumInputs,
				config.NumOutputs, config.InitFitness)
			nextGenomeID++
		}
	} else {
		for i := 0; i < config.PopulationSize; i++ {
			population[i] = NewGenome(nextGenomeID, config.NumInputs,
				config.NumOutputs, config.InitFitness)
			nextGenomeID++
		}
	}

	// initialize the first species with a randomly selected genome
	s := NewSpecies(nextSpeciesID, population[rand.Intn(len(population))])
	species := []*Species{s}
	nextSpeciesID++

	return &NEAT{
		Config:        config,
		Population:    population,
		Species:       species,
		Activations:   activations,
		Evaluation:    evaluation,
		Comparison:    NewComparisonFunc(config.MinimizeFitness),
		Best:          population[rand.Intn(config.PopulationSize)].Copy(),
		Statistics:    NewStatistics(config.NumGenerations),
		nextGenomeID:  nextGenomeID,
		nextSpeciesID: nextSpeciesID,
	}
}

// Summarize summarizes current state of evolution process.
func (n *NEAT) Summarize(gen int) {
	// summary template
	tmpl := "Gen. %4d | Num. Species: %4d | Best Fitness: %.4f | " +
		"Avg. Fitness: %.4f"

	// compose each line of summary and the spacing of separating line
	str := fmt.Sprintf(tmpl, gen, len(n.Species),
		n.Best.Fitness, n.Statistics.AvgFitness[gen])
	spacing := int(math.Max(float64(len(str)), 80.0))

	for i := 0; i < spacing; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("\n%s\n", str)
	for i := 0; i < spacing; i++ {
		fmt.Printf("-")
	}
	fmt.Println()
}

// Evaluate evaluates fitness of every genome in the population. After the
// evaluation, their fitness scores are recored in each genome.
func (n *NEAT) Evaluate() {
	for _, genome := range n.Population {
		genome.Evaluate(n.Evaluation)
	}
}

// Speciate performs speciation of each genome. The speciation mechanism is as
// follows (from http://nn.cs.utexas.edu/downloads/papers/stanley.phd04.pdf):
//
//	The Genome Loop:
//		Take next genome g from P
//		The Species Loop:
//			If all species in S have been checked:
//				create new species snew and place g in it
//			Else:
//				get next species s from S
//				If g is compatible with s:
//					add g to s
//			If g has not been placed:
//				Species Loop
//		If not all genomes in G have been placed:
//			Genome Loop
//		Else STOP
//
func (n *NEAT) Speciate() {
	for _, genome := range n.Population {
		registered := false
		for i := 0; i < len(n.Species) && !registered; i++ {
			dist := Compatibility(n.Species[i].Representative, genome,
				n.Config.CoeffUnmatching, n.Config.CoeffMatching)

			if dist <= n.Config.DistanceThreshold {
				n.Species[i].Register(genome, n.Config.MinimizeFitness)
				registered = true
			}
		}

		if !registered {
			n.Species = append(n.Species, NewSpecies(n.nextSpeciesID, genome))
			n.nextSpeciesID++
		}
	}
}

// Reproduce performs reproduction of genomes in each species. Reproduction is
// performed under the assumption of speciation being already executed. The
// number of eliminated genomes in each species is determined by rate of
// elimination specified in n.Config; after some number of genomes are
// eliminated, the empty space is filled with resulting genomes of crossover
// among surviving genomes. If the number of eliminated genomes is 0 or less
// then 2 genomes survive, every member survives and mutates.
func (n *NEAT) Reproduce() {
	nextGeneration := make([]*Genome, 0, n.Config.PopulationSize)
	for _, s := range n.Species {
		// genomes in this species can inherit to the next generation, if two or
		// more genomes survive in this species, and there is room for more
		// children, i.e., at least one genome must be eliminated.
		numSurvived := int(math.Ceil(float64(len(s.Members)) *
			n.Config.SurvivalRate))
		numEliminated := len(s.Members) - numSurvived

		// reproduction of this species is only executed, if there is enough room.
		if numSurvived > 2 && numEliminated > 0 {
			// adjust the fitness of each member genome of this species.
			s.ExplicitFitnessSharing()

			sort.Slice(s.Members, func(i, j int) bool {
				return n.Comparison(s.Members[i], s.Members[j])
			})
			s.Members = s.Members[:numSurvived]

			// fill the spaces that are made by eliminated genomes, by creating
			// children.
			for i := 0; i < numEliminated; i++ {
				perm := rand.Perm(numSurvived)
				p0 := s.Members[perm[0]] // parent 0
				p1 := s.Members[perm[1]] // parent 1

				// create a child from two chosen parents as a result of crossover;
				// mutate the child given the rate of mutation of children.
				child := Crossover(n.nextGenomeID, p0, p1, n.Config.InitFitness)
				if rand.Float64() < n.Config.RateMutateChild {
					child.MutatePerturb(n.Config.RatePerturb)
					child.MutateAddNode(n.Config.RateAddNode, n.randActivationFunc())
					child.MutateAddConn(n.Config.RateAddConn)
				} else {
					// if the two parents are identical, definitely mutate the child.
					if p0.ID == p1.ID {
						child.MutatePerturb(n.Config.RatePerturb)
						child.MutateAddNode(n.Config.RateAddNode, n.randActivationFunc())
						child.MutateAddConn(n.Config.RateAddConn)
					}
				}
				n.nextGenomeID++

				nextGeneration = append(nextGeneration, child)
			}

			// mutate all the genomes that survived.
			for _, genome := range s.Members {
				genome.MutatePerturb(n.Config.RatePerturb)
				genome.MutateAddNode(n.Config.RateAddNode, n.randActivationFunc())
				genome.MutateAddConn(n.Config.RateAddConn)
				nextGeneration = append(nextGeneration, genome)
			}
		} else {
			// otherwise, they all survive, and mutate.
			for _, genome := range s.Members {
				genome.MutatePerturb(n.Config.RatePerturb)
				genome.MutateAddNode(n.Config.RateAddNode, n.randActivationFunc())
				genome.MutateAddConn(n.Config.RateAddConn)
				nextGeneration = append(nextGeneration, genome)
			}
		}

		s.Flush()
	}

	// update the population with the new generation
	n.Population = nextGeneration
}

// randActivationFunc is a helper function that returns a random activation
// function.
func (n *NEAT) randActivationFunc() *ActivationFunc {
	return n.Activations[rand.Intn(len(n.Activations))]
}

// Run executes evolution.
func (n *NEAT) Run() {
	if n.Config.Verbose {
		n.Config.Summarize()
	}

	// for each generation
	for i := 0; i < n.Config.NumGenerations; i++ {
		n.Evaluate()

		// speciate genomes and reproduce children genomes
		n.Speciate()
		n.Reproduce()

		// eliminate stagnant species
		if len(n.Species) > 1 {
			var survived []*Species
			for j := range n.Species {
				if n.Species[j].Stagnation <= n.Config.StagnationLimit {
					n.Species[j].Stagnation++
					survived = append(survived, n.Species[j])
				}
			}
			n.Species = survived
		}

		// update the best genome
		for _, genome := range n.Population {
			if n.Comparison(genome, n.Best) {
				n.Best = genome.Copy()
			}
		}

		n.Statistics.Update(i, n)
		if n.Config.Verbose {
			n.Summarize(i)
		}
	}
}
