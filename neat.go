package neat

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"sync"
)

// NEAT is the implementation of NeuroEvolution of Augmenting Topology (NEAT).
type NEAT struct {
	Config     *Config        // configuration
	Population []*Genome      // population of genome
	Species    []*Species     // subpopulations of genomes grouped by species
	Evaluation EvaluationFunc // evaluation function
	Comparison ComparisonFunc // comparison function
	Best       *Genome        // best genome
	Statistics *Statistics    // statistics

	nextGenomeID  int // genome ID that is assigned to a newly created genome
	nextSpeciesID int // species ID that is assigned to a newly created species
}

// New creates a new instance of NEAT with provided argument configuration and
// an evaluation function.
func New(config *Config, evaluation EvaluationFunc) *NEAT {
	nextGenomeID := 0
	nextSpeciesID := 0

	population := make([]*Genome, config.PopulationSize)
	for i := 0; i < config.PopulationSize; i++ {
		population[i] = NewGenome(nextGenomeID, config.NumInputs,
			config.NumOutputs, config.InitFitness)
		nextGenomeID++
	}

	// initialize the first species with a randomly selected genome
	s := NewSpecies(nextSpeciesID, population[rand.Intn(len(population))])
	species := []*Species{s}
	nextSpeciesID++

	return &NEAT{
		Config:        config,
		Population:    population,
		Species:       species,
		Evaluation:    evaluation,
		Comparison:    NewComparisonFunc(config.MinimizeFitness),
		Best:          population[rand.Intn(config.PopulationSize)],
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

	// average fitness of this generation
	avgFitness := 0.0
	for _, genome := range n.Population {
		avgFitness += genome.Fitness
	}
	avgFitness /= float64(n.Config.PopulationSize)

	// compose each line of summary and the spacing of separating line
	str := fmt.Sprintf(tmpl, gen, len(n.Species), n.Best.Fitness, avgFitness)
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
	// for parallel processing; potentially, if there are the same number of CPUs
	// as the population size, all genomes in the population can be evaluated at
	// the same time.
	runtime.GOMAXPROCS(n.Config.PopulationSize)

	var wg sync.WaitGroup

	for i := range n.Population {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			n.Population[j].Evaluate(n.Evaluation)
		}(i)
	}

	wg.Wait()

	// explicit fitness sharing
	/*

		for i, genome0 := range n.Population {
			adjustment := 0.0
			for _, genome1 := range append(n.Population[:i], n.Population[i+1:]...) {
				if Compatibility(genome0, genome1, n.Config.CoeffUnmatching,
					n.Config.CoeffMatching) <= n.Config.DistanceThreshold {
					adjustment += 1.0
				}
			}

			if adjustment != 0.0 {
				genome0.Fitness /= adjustment
			}
		}

	*/
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
					Mutate(child, n.Config.RatePerturb,
						n.Config.RateAddNode, n.Config.RateAddConn)
				} else {
					// if the two parents are identical, definitely mutate the child.
					if p0.ID == p1.ID {
						Mutate(child, n.Config.RatePerturb,
							n.Config.RateAddNode, n.Config.RateAddConn)
					}
				}
				n.nextGenomeID++

				nextGeneration = append(nextGeneration, child)
			}

			// mutate all the genomes that survived.
			for _, genome := range s.Members {
				Mutate(genome, n.Config.RatePerturb,
					n.Config.RateAddNode, n.Config.RateAddConn)
				nextGeneration = append(nextGeneration, genome)
			}
		} else {
			// otherwise, they all survive, and mutate.
			for _, genome := range s.Members {
				Mutate(genome, n.Config.RatePerturb,
					n.Config.RateAddNode, n.Config.RateAddConn)
				nextGeneration = append(nextGeneration, genome)
			}
		}

		s.Flush()
	}

	// update the population with the new generation
	n.Population = nextGeneration
}

// Run executes evolution.
func (n *NEAT) Run() {
	if n.Config.Verbose {
		n.Config.Summarize()
	}

	// for each generation
	for i := 0; i < n.Config.NumGenerations; i++ {
		n.Evaluate()

		n.Statistics.Update(i, n)
		if n.Config.Verbose {
			n.Summarize(i)
		}

		// speciate genomes and reproduce children genomes
		n.Speciate()
		n.Reproduce()

		// eliminate stagnant species
		if len(n.Species) > 1 {
			survived := make([]*Species, 0)
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
	}
}
