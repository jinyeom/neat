/*


neat.go implementation of NEAT.

@licstart   The following is the entire license notice for
the Go code in this page.

Copyright (C) 2016 jin yeom, whitewolf.studio

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

As additional permission under GNU GPL version 3 section 7, you
may distribute non-source (e.g., minimized or compacted) forms of
that code without the copy of the GNU GPL normally required by
section 4, provided you include this license notice and a URL
through which recipients can access the Corresponding Source.

@licend    The above is the entire license notice
for the Go code in this page.


*/

package neat

var (
	// globalInnovNum is a global variable that keeps track of
	// the chronology of the evolution as a global innovation
	// number; it is initialized as 0. Users cannot directly
	// access globalInnovNum.
	globalInnovNum = 0
)

// Config is a wrapper of all configurations of NEAT.
type Config struct {
	numSensors     int // number of sensors
	numOutputs     int // number of outputs
	populationSize int // population size

	evalFunc *EvaluationFunc // evalutation function

	crossoverRate  float64 // crossover rate
	mutAddNodeRate float64 // mutation rate for adding a node
	mutAddConnRate float64 // mutation rate for adding a connection
	mutWeightRate  float64 // mutation rate of weights of connections
}

// NEAT is an implementation of NeuroEvolution of Augmenting
// Topologies; it includes
type NEAT struct {
	config     *Config   // NEAT configuration
	population []*Genome // population of genomes
}

// New creates NEAT and initializes its environment given a configuration.
func New(config *Config) (*NEAT, error) {
	// initialize global innovation number
	globalInnovNum = (config.numSensors + 1) * config.numOutputs

	// initialize population
	population := make([]*Genome, config.populationSize)
	for i := range population {
		genome, err := NewGenome(i, config.numSensors, config.numOutputs)
		if err != nil {
			return nil, err
		}
		population[i] = genome
	}
	return &NEAT{
		config:     config,
		population: population,
	}, nil
}

// Run starts the evolution process of NEAT.
func (n *NEAT) Run() {

}
