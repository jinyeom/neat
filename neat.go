package neat

import (
	"encoding/json"
	"os"
)

// Config consists of all hyperparameter settings for NEAT. It can be imported
// from a JSON file.
type Config struct {
	// evolution settings
	NumGenerations int `json:"numGenerations"` // number of generations
	PopulationSize int `json:"populationSize"` // number of genomes in population

	// mutation rates settings
	RatePerturb float64 `json:"ratePerturb"` // mutation by perturbing weights
	RateAddNode float64 `json:"rateAddNode"` // mutation by adding a node
	RateAddConn float64 `json:"rateAddConn"` // mutation by adding a connection

	// compatibility distance coefficient settings
	CoeffUnmatching float64 `json:"coeffUnmatching"` // unmatching genes
	CoeffMatching   float64 `json:"coeffMatching"`   // matching genes
}

// NewConfig creates a new instance of Config, given the name of a JSON file
// that consists of the hyperparameter settings.
func NewConfigJSON(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

}

type NEAT struct {
	nextGenomeID int // genome ID that is assigned to a newly created genome

	Config     *Config
	Population map[*Genome]float64 // genomes mapped to fitness scores
	Species    []*Species          // subpopulations grouped by species
}
