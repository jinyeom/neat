/*


param.go parameters of NEAT.

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
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// paramFileExt is the extension of a parameter file, np, which
	// stands for NEAT parameter.
	paramFileExt = ".np"
)

// Param is a wrapper for all parameters of NEAT.
type Param struct {
	// Topology parameter
	NumSensors int // number of sensors
	NumOutputs int // number of outputs

	PopulationSize int // population size
	NumGeneration  int // number of generations

	// Crossover and Mutation rates
	CrossoverRate  float64 // crossover rate
	MutAddNodeRate float64 // mutation rate for adding a node
	MutAddConnRate float64 // mutation rate for adding a connection
	MutWeightRate  float64 // mutation rate of weights of connections

	// Coefficients for compatibility distance
	CoeffExcess   float64 // coefficient for excess
	CoeffDisjoint float64 // coefficient for disjoint
	CoeffWeight   float64 // coefficient for average weight
}

// NewParam creates a new NEAT parameter wrapper, given a name of a parameter
// file that contains its presets. A parameter file must have a ".np"
// extension and must be structured as the following example:
//
// example_param.np:
//
//  NumSensors 5
//  NumOutputs 3
//  PopulationSize 50
//  NumGeneration 50
//  CrossoverRate 0.1
//  MutAddNodeRate 0.1
//  MutAddConnRate 0.1
//  MutWeightRate 0.1
//  CoeffExcess 0.5
//  CoeffDisjoint 0.5
//  CoeffWeight 0.5
//
func NewParam(filename string) (*Param, error) {
	// parse parameter file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// check file extension
	if ext := filepath.Ext(filename); ext != paramFileExt {
		return nil, fmt.Errorf("Invalid file type: %s", ext)
	}

	// initialize presets as empty
	param := &Param{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parsed := strings.Split(scanner.Text(), " ")
		switch parsed[0] {
		case "NumSensors":
			numSensors, err := strconv.Atoi(parsed[1])
			if err != nil {
				return nil, err
			}
			param.NumSensors = numSensors
		case "NumOutputs":
			numOutputs, err := strconv.Atoi(parsed[1])
			if err != nil {
				return nil, err
			}
			param.NumOutputs = numOutputs
		case "PopulationSize":
			populationSize, err := strconv.Atoi(parsed[1])
			if err != nil {
				return nil, err
			}
			param.PopulationSize = populationSize
		case "MutAddNodeRate":
			mutAddNodeRate, err := strconv.ParseFloat(parsed[1], 64)
			if err != nil {
				return nil, err
			}
			param.MutAddNodeRate = mutAddNodeRate
		case "MutAddConnRate":
			mutAddConnRate, err := strconv.ParseFloat(parsed[1], 64)
			if err != nil {
				return nil, err
			}
			param.MutAddConnRate = mutAddConnRate
		case "MutWeightRate":
			mutWeightRate, err := strconv.ParseFloat(parsed[1], 64)
			if err != nil {
				return nil, err
			}
			param.MutWeightRate = mutWeightRate
		case "CoeffExcess":
			coeffExcess, err := strconv.ParseFloat(parsed[1], 64)
			if err != nil {
				return nil, err
			}
			param.CoeffExcess = coeffExcess
		case "CoeffDisjoint":
			coeffDisjoint, err := strconv.ParseFloat(parsed[1], 64)
			if err != nil {
				return nil, err
			}
			param.CoeffDisjoint = coeffDisjoint
		case "CoeffWeight":
			coeffWeight, err := strconv.ParseFloat(parsed[1], 64)
			if err != nil {
				return nil, err
			}
			param.CoeffWeight = coeffWeight
		}
	}

	return param, nil
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
	// coefficient for excess and disjoint genes
	if p.CoeffExcess < 0.0 || p.CoeffDisjoint < 0.0 {
		return errors.New("Invalid coefficient for excess/disjoint genes")
	}
	return nil
}
