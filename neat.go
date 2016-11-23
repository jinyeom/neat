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

import (
	"errors"
)

var (
	// globalInnovNum is a global variable that keeps track of
	// the chronology of the evolution; it is initialized as 0.
	// Users cannot directly access globalInnovNum.
	globalInnovNum = 0
)

// Param is a wrapper for all parameters of NEAT.
type Param struct {
	NumSensors     int // number of sensors
	NumOutputs     int // number of outputs
	PopulationSize int // population size

	CrossoverRate  float64 // crossover rate
	MutAddNodeRate float64 // mutation rate for adding a node
	MutAddConnRate float64 // mutation rate for adding a connection
	MutWeightRate  float64 // mutation rate of weights of connections
}

// IsValid checks the parameter's validity. It returns an error that
// indicates the invalid portion of the parameter; return nil otherwise.
func (p *Param) IsValid() error {
	// number of sensors and outputs
	if p.NumSensors < 1 || p.NumOutputs < 1 {
		return errors.New("Invalid number of sensors and/or outputs")
	}
	// population size
	if p.PopulationSize < 1 {
		return errors.New("Invalid size of population")
	}
	// crossover rate
	if p.CrossoverRate < 0.0 {
		return errors.New("Invalid crossover rate")
	}
	// mutation rate for adding a node
	if p.MutAddNodeRate < 0.0 || p.MutAddConnRate < 0.0 || p.MutWeightRate < 0.0 {
		return errors.New("Invalid mutation rate")
	}
	return nil
}

// NEAT is an implementation of NeuroEvolution of Augmenting
// Topologies; it includes
type NEAT struct {
	param      *Param          // NEAT parameters
	evalFunc   *EvaluationFunc // evaluation function
	population []*Genome       // population of genomes
	species    []*Species      // ordered list of species
}

// New creates NEAT and initializes its environment given a set of parameters.
func New(param *Param, evalFunc *EvaluationFunc) (*NEAT, error) {
	// check if parameter is valid
	err := param.IsValid()
	if err != nil {
		return nil, err
	}
	// initialize population
	population := make([]*Genome, param.PopulationSize)
	for i := range population {
		population[i] = NewGenome(i, param)
	}
	return &NEAT{
		param:      param,
		evalFunc:   evalFunc,
		population: population,
		species:    make([]*Species, 0), // to be fixed
	}, nil
}

// Run starts the evolution process of NEAT.
func (n *NEAT) Run() {

}
