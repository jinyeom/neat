package neat

import (
	"encoding/json"
	"fmt"
	"os"
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
	nextGenomeID int // genome ID that is assigned to a newly created genome

	Config     *Config             // configuration
	Population map[*Genome]float64 // population of genome
	Species    []*Species          // subpopulations of genomes grouped by species
	Evaluation EvaluationFunc      // evaluation function
	Best       *Genome             // best performing genome
}

// New creates a new instance of NEAT with provided argument configuration and
// an evaluation function.
func New(config *Config, evaluation EvaluationFunc) *NEAT {
	nextGenomeID := 0
	population := make(map[*Genome]float64)
	for i := 0; i < config.PopulationSize; i++ {
		g := NewGenome(nextGenomeID, config.NumInputs, config.NumOutputs)
		population[g] = config.InitFitness
		nextGenomeID++
	}
	return &NEAT{
		nextGenomeID: nextGenomeID,
		Config:       config,
		Population:   population,
		Species:      make([]*Species, 0),
		Evaluation:   evaluation,
	}
}

// Run
func (n *NEAT) Run() {

}
