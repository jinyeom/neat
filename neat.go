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
	"bufio"
	"strings"
)

const (
	// paramFileExt is the extension of a parameter file, np, which
	// stands for NEAT parameter.
	paramFileExt = ".np"
)

var (
	// globalInnovNum is a global variable that keeps track of
	// the chronology of the evolution; it is initialized as 0.
	globalInnovNum = 0

	// nodeInnovations is a global list of structural innovations of newly added
	// nodes that are added during mutations; this list of innovations
	// maps innovation numbers of connections that are split due to a mutation
	// to innovation numbers of nodes that split them.
	nodeInnovations = make(map[int]int)

	// connInnovations is a global list of structural innovation of newly added
	// connections that are added during mutations; this list of innovations
	// maps IDs of nodes that are connected due to mutations to innovation
	// numbers of connections that connect them.
	connInnovations = make(map[[]int]int)
)

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
