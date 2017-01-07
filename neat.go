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
	"runtime"
	"sort"
	"sync"
)

var (
	// globalInnovNum is a global variable that keeps track of
	// the chronology of the evolution via historical marking;
	// it is initialized as 1, since 0 is reserved for innovation
	// map's zero value.
	globalInnovNum = 1

	// innovations is a global list of structural innovation of newly added
	// connections that are added during mutations; this list of innovations
	// maps IDs of nodes that are connected due to mutations to innovation
	// numbers of connections that connect them.
	innovations = make(map[[2]int]int)

	// param is a global parameter that can only be manipulated internally;
	// it needs to be initialized before creating a new NEAT struct.
	param *Param

	// toolbox is a wrapper of all functions that are utilized widely for
	// NEAT operations, such as selection, activation, or evaluation; it needs
	// to be initialized before creating a new NEAT struct.
	toolbox *Toolbox

	// initPass is an indicator of whether param and afnSet are initialized;
	// it only becomes true when Init() is called.
	initPass = false
)

// Init initializes NEAT by setting parameters and a set of activation set;
// it returns an error if the argument parameter or activation set is invalid.
func Init(p *Param, tb *Toolbox) error {
	// check validity
	if err := p.IsValid(); err != nil {
		return err
	}
	if err := tb.IsValid(); err != nil {
		return err
	}

	param = p
	toolbox = tb
	initPass = true

	return nil
}

// NEAT is an implementation of NeuroEvolution of Augmenting
// Topologies; it includes
type NEAT struct {
	population []*Genome  // population of genomes
	species    []*Species // ordered list of species
}

// New creates NEAT and initializes its environment given a set of parameters.
func New() (*NEAT, error) {
	if !initPass {
		return nil, errors.New("Initialization check failed")
	}

	// initialize population
	population := make([]*Genome, param.PopulationSize)
	for i := range population {
		population[i] = NewGenome(i)
	}

	// initialize slice of species with one species
	species := []*Species{NewSpecies(0, population[0])}

	return &NEAT{
		population: population,
		species:    species,
	}, nil
}

// evaluate executes evaluation function on each genome of the population,
// and sets their fitness values.
func (n *NEAT) evaluate() {
	for i, genome := range n.population {
		score := toolbox.Evaluation(genome)
		n.population[i].SetFitness(score)
	}
}

func (n *NEAT) evaluateParallel(procs int) {
	runtime.GOMAXPROCS(procs)

	// var wg sync.WaitGroup

	// to be implemented

}

func (n *NEAT) speciate() {
	for i, genome := range n.population {
		// species loop
		pass := false
		for j, species := range n.species {
			d := genome.Distance(species.representative)
			if d < param.DistThreshold {
				n.species[j].AddMember(n.population[i])
				pass = true
				break
			}
		}
		if pass == false {
			// create a new species
			ns := NewSpecies(len(n.species), n.population[i])
			n.species = append(n.species, ns)
		}
	}

	// remove species with no members
	for i := range n.species {
		if len(n.species[i].members) == 0 {
			n.species = append(n.species[:i], n.species[i+1:])
		}
	}
}

// Run executes NEAT algorithm.
func (n *NEAT) Run(verbose bool) {
	for i := 0; i < param.NumGeneration; i++ {
		n.evaluate()
		n.speciate()

		// mutation
		for j := range n.population {
			n.population[j].Mutate()
		}

		// crossover and fitness share
		for j, niche := range n.species {
			niche.FitnessShare()
			niche.age++
		}
	}
}

// RunParallel executes NEAT algorithm in parallel by separating the
// evaluation of individuals in a population into different processor.
func (n *NEAT) RunParallel(verbose bool, procs int) {
	runtime.GOMAXPROCS(procs)

	var wg sync.WaitGroup

	for i := 0; i < param.NumGeneration; i++ {
		// number of evaluations per processor
		numEval := param.PopulationSize / procs

		start := 0      // iterator
		next := numEval // next iteration
		for p := 0; p < procs; p++ {
			// handle leftover genomes
			if next > param.PopulationSize {
				next = param.PopulationSize
			}

			wg.Add(1)

			go func(start int) {
				defer wg.Done()
				// iterate through this group of genomes
				iter := start
				for iter < next {
					score := toolbox.Evaluation(n.population[iter])
					n.population[iter].SetFitness(score)
					iter++
				}
			}(start)

			start = next
			next += numEval
		}
		wg.Wait()

		// genome loop

		// species loop

		// mutate

		// crossover
	}
}
