package neat

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"sync"
	"text/tabwriter"
)

// Config consists of all hyperparameter settings for NEAT. It can be imported
// from a JSON file.
type Config struct {
	// neural network settings
	NumInputs  int `json:"numInputs"`  // number of inputs (including bias)
	NumOutputs int `json:"numOutputs"` // number of outputs

	// evolution settings
	NumGenerations  int     `json:"numGenerations"`  // number of generations
	PopulationSize  int     `json:"populationSize"`  // size of population
	InitFitness     float64 `json:"initFitness"`     // initial fitness score
	MinimizeFitness bool    `json:"minimizeFitness"` // true if minimizing fitness
	SurvivalRate    float64 `json:"survivalRate"`    // survival rate
	StagnationLimit int     `json:"stagnationLimit"` // limit of stagnation
	HOFSize         int     `json:"hallOfFameSize"`  // hall of fame size

	// mutation rates settings
	RatePerturb float64 `json:"ratePerturb"` // mutation by perturbing weights
	RateAddNode float64 `json:"rateAddNode"` // mutation by adding a node
	RateAddConn float64 `json:"rateAddConn"` // mutation by adding a connection

	// compatibility distance coefficient settings
	DistanceThreshold float64 `json:"distanceThreshold"` // distance threshold
	CoeffUnmatching   float64 `json:"coeffUnmatching"`   // unmatching genes
	CoeffMatching     float64 `json:"coeffMatching"`     // matching genes
}

// NewConfig creates a new instance of Config, given the name of a JSON file
// that consists of the hyperparameter settings.
func NewConfigJSON(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := &Config{}
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

// Summarize prints the summarized configuration on terminal.
func (c *Config) Summarize() {
	w := tabwriter.NewWriter(os.Stdout, 40, 1, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "==================================================\n")
	fmt.Fprintf(w, "Summary of NEAT hyperparameter configuration\t\n")
	fmt.Fprintf(w, "==================================================\n")
	fmt.Fprintf(w, "--------------------------------------------------\n")
	fmt.Fprintf(w, "Neural network settings\t\n")
	fmt.Fprintf(w, "--------------------------------------------------\n")
	fmt.Fprintf(w, "+ Number of inputs\t%d\t\n", c.NumInputs)
	fmt.Fprintf(w, "+ Number of outputs \t%d\t\n", c.NumOutputs)
	fmt.Fprintf(w, "--------------------------------------------------\n")
	fmt.Fprintf(w, "Evolution settings\t\n")
	fmt.Fprintf(w, "--------------------------------------------------\n")
	fmt.Fprintf(w, "+ Number of generations\t%d\t\n", c.NumGenerations)
	fmt.Fprintf(w, "+ Population size\t%d\t\n", c.PopulationSize)
	fmt.Fprintf(w, "+ Initial fitness score\t%.3f\t\n", c.InitFitness)
	fmt.Fprintf(w, "+ Fitness is being minimized\t%t\t\n", c.MinimizeFitness)
	fmt.Fprintf(w, "+ Rate of survival each generation\t%.3f\t\n", c.SurvivalRate)
	fmt.Fprintf(w, "+ Limit of species' stagnation\t%d\t\n", c.StagnationLimit)
	fmt.Fprintf(w, "+ Size of Hall of Fame\t%d\t\n", c.HOFSize)
	fmt.Fprintf(w, "--------------------------------------------------\n")
	fmt.Fprintf(w, "Mutation settings\t\n")
	fmt.Fprintf(w, "--------------------------------------------------\n")
	fmt.Fprintf(w, "+ Rate of perturbation of weights\t%.3f\t\n", c.RatePerturb)
	fmt.Fprintf(w, "+ Rate of adding a node\t%.3f\t\n", c.RateAddNode)
	fmt.Fprintf(w, "+ Rate of adding a connection\t%.3f\t\n", c.RateAddConn)
	fmt.Fprintf(w, "--------------------------------------------------\n")
	fmt.Fprintf(w, "Compatibility distance settings\t\n")
	fmt.Fprintf(w, "--------------------------------------------------\n")
	fmt.Fprintf(w, "+ Distance threshold\t%.3f\t\n", c.DistanceThreshold)
	fmt.Fprintf(w, "+ Unmatching connection genes\t%.3f\t\n", c.CoeffUnmatching)
	fmt.Fprintf(w, "+ Matching connection genes\t%.3f\t\n", c.CoeffMatching)
	fmt.Fprintf(w, "--------------------------------------------------\n")
	w.Flush()
}

// NEAT is the implementation of NeuroEvolution of Augmenting Topology (NEAT).
type NEAT struct {
	sync.Mutex

	Config     *Config        // configuration
	Population []*Genome      // population of genome
	Species    []*Species     // subpopulations of genomes grouped by species
	Evaluation EvaluationFunc // evaluation function
	Comparison ComparisonFunc // comparison function
	HallOfFame *HallOfFame    // best performing genomes
	Statistics *Statistics    // statistics

	nextGenomeID  int // genome ID that is assigned to a newly created genome
	nextSpeciesID int // species ID that is assigned to a newly created species
}

// New creates a new instance of NEAT with provided argument configuration and
// an evaluation function.
func New(config *Config, evaluation EvaluationFunc) *NEAT {
	nextGenomeID := 0
	nextSpeciesID := 0

	comparison := func(g0, g1 *Genome) bool {
		return g0.Fitness < g1.Fitness
	}
	if !config.MinimizeFitness {
		comparison = func(g0, g1 *Genome) bool {
			return g0.Fitness > g1.Fitness
		}
	}

	population := make([]*Genome, config.PopulationSize)
	for i := 0; i < config.PopulationSize; i++ {
		population[i] = NewGenome(nextGenomeID, config.NumInputs, config.NumOutputs)
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
		Comparison:    comparison,
		HallOfFame:    NewHallOfFame(config.HOFSize, config.InitFitness),
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
	str := fmt.Sprintf(tmpl, gen, len(n.Species),
		n.HallOfFame.BestScore(), avgFitness)
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

		if numSurvived > 2 && numEliminated > 0 {
			sort.Slice(s.Members, func(i, j int) bool {
				return n.Comparison(s.Members[i], s.Members[j])
			})
			s.Members = s.Members[:numSurvived]

			for i := 0; i < numEliminated; i++ {
				perm := rand.Perm(numSurvived)
				p0 := s.Members[perm[0]] // parent 0
				p1 := s.Members[perm[1]] // parent 1

				child := Crossover(n.nextGenomeID, p0, p1, n.Config.InitFitness)
				Mutate(child, n.Config.RatePerturb,
					n.Config.RateAddNode, n.Config.RateAddConn)
				n.nextGenomeID++

				nextGeneration = append(nextGeneration, child)
			}

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
func (n *NEAT) Run(verbose bool) {
	if verbose {
		n.Config.Summarize()
	}

	for i := 0; i < n.Config.NumGenerations; i++ {
		n.Evaluate()
		n.Statistics.Update(i, n)
		if verbose {
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
			n.HallOfFame.Update(genome, n.Comparison)
		}
	}
}

// inheritParallel performs crossover and mutation within all species in
// parallel.
func (n *NEAT) inheritParallel() {
	runtime.GOMAXPROCS(len(n.Species))

	var wg sync.WaitGroup

	nextGeneration := struct {
		sync.Mutex
		population []*Genome // children genome for the next generation
	}{population: make([]*Genome, 0, n.Config.PopulationSize)}

	for _, species := range n.Species {
		wg.Add(1)
		go func(s *Species) {
			defer wg.Done()

			// genomes in this species can inherit to the next generation, if two or
			// more genomes survive in this species, and there is room for more
			// children, i.e., at least one genome must be eliminated.
			numSurvived := int(math.Ceil(float64(len(s.Members)) *
				n.Config.SurvivalRate))
			numEliminated := len(s.Members) - numSurvived

			if numSurvived > 2 && numEliminated > 0 {
				sort.Slice(s.Members, func(i, j int) bool {
					return n.Comparison(s.Members[i], s.Members[j])
				})
				s.Members = s.Members[:numSurvived]

				for i := 0; i < numEliminated; i++ {
					perm := rand.Perm(numSurvived)
					p0 := s.Members[perm[0]]
					p1 := s.Members[perm[1]]

					n.Lock()
					child := Crossover(n.nextGenomeID, p0, p1, n.Config.InitFitness)
					n.nextGenomeID++
					n.Unlock()

					nextGeneration.Lock()
					nextGeneration.population = append(nextGeneration.population, child)
					nextGeneration.Unlock()
				}
			}
		}(species)
	}

	wg.Wait()

	// update the population with the new generation
	n.Population = nextGeneration.population
}

// RunParallel executes evolution with parallel processing.
func (n *NEAT) RunParallel(summarize bool) {
	if summarize {
		n.Config.Summarize()
	}
	for i := 0; i < n.Config.NumGenerations; i++ {
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

		n.inheritParallel()
	}
}
