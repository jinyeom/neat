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
	// it is initialized as nil pointer and needs to be initialized before
	// creating a new NEAT struct.
	param *Param

	// afnSet is a global set of activation functions that are used within
	// a network that is being evolved via NEAT; it is initialized as a nil
	// pointer and needs to be initialized before creating a new NEAT struct.
	afnSet ActivationSet

	// initPass is an indicator of whether param and afnSet are initialized;
	// it only becomes true when Init() is called.
	initPass = false
)

// Init initializes NEAT by setting parameters and a set of activation set;
// it returns an error if the argument parameter or activation set is invalid.
func Init(p *Param, a ActivationSet) error {
	// intialize parameter
	if err := p.IsValid(); err != nil {
		return err
	}
	param = p

	// initialize activation set
	if a == nil || len(a) == 0 {
		return errors.New("invalid activation set")
	}
	afnSet = a

	// pass and return
	initPass = true
	return nil
}

// NEAT is an implementation of NeuroEvolution of Augmenting
// Topologies; it includes
type NEAT struct {
	evalFunc   EvaluationFunc // evaluation function
	population []*Genome      // population of genomes
	species    []*Species     // ordered list of species
}

// New creates NEAT and initializes its environment given a set of parameters.
func New(evalFunc EvaluationFunc) (*NEAT, error) {
	if !initPass {
		return nil, errors.New("initializing check failed")
	}

	// initialize population
	population := make([]*Genome, param.PopulationSize)
	for i := range population {
		population[i] = NewGenome(i)
	}

	// initialize slice of species with one species
	species := []*Species{NewSpecies(0, population[0])}

	return &NEAT{
		evalFunc:   evalFunc,
		population: population,
		species:    species,
	}, nil
}

// sh implements a part of the explicit fitness sharing function, sh.
// If a compatibility distance 'd' is larger than the compatibility
// threshold 'dt', return 0; otherwise, return 1.
func sh(d float64) float64 {
	if d > param.DistThreshold {
		return 0.0
	}
	return 1.0
}

// FitnessShare computes and assigns the shared fitness of genomes,
// via explicit fitness sharing.
func (n *NEAT) FitnessShare() {
	adjusted := make(map[int]float64)
	for _, g0 := range n.population {
		adjustment := 0.0
		for _, g1 := range n.population {
			adjustment += sh(g0.Distance(g1))
		}
		if adjustment != 0.0 {
			adjusted[g0.gid] = g0.fitness / adjustment
		}
	}
	for i := range n.population {
		n.population[i].fitness = adjusted[n.population[i].gid]
	}
}

// Run executes NEAT algorithm.
func (n *NEAT) Run(verbose bool) {
	for i := 0; i < param.NumGeneration; i++ {
		for i := range n.population {
			network := n.population[i].Decode()
			n.population[i].SetFitness(n.evalFunc(network))
		}

		// genome loop

		// species loop

		// mutate

		// crossover
	}
}

// RunParallel executes NEAT algorithm in parallel by separating the
// evaluation of individuals in a population into different processor.
func (n *NEAT) RunParallel(verbose bool, procs int) {
	runtime.GOMAXPROCS(procs)

	var wg sync.WaitGroup
	wg.Add(procs)

	for i := 0; i < param.NumGeneration; i++ {
		// number of evaluations per processor
		numEval := param.PopulationSize / procs

		iter := 0       // iterator
		next := numEval // next iteration
		for p := 0; p < procs; p++ {
			go func() {
				defer wg.Done()

				// handle leftover genomes
				if next > param.PopulationSize {
					next = param.PopulationSize
				}

				// iterate through this group of genomes
				for iter < next {
					network := n.population[iter].Decode()
					n.population[iter].SetFitness(n.evalFunc(network))
					iter++
				}
			}()

			iter = next
			next += numEval
		}

		// genome loop

		// species loop

		// mutate

		// crossover
	}
}
