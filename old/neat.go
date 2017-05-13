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
	"fmt"
	"log"
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
// and sets their fitness values. Evaluation is done with maximum of 4 parallel
// processors for fast performance.
func (n *NEAT) evaluate() {
	runtime.GOMAXPROCS(4)

	gap := len(n.population) / 4

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < gap; i++ {
			score := toolbox.Evaluation(n.population[i])
			n.population[i].SetFitness(score)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := gap; i < gap*2; i++ {
			score := toolbox.Evaluation(n.population[i])
			n.population[i].SetFitness(score)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := gap * 2; i < gap*3; i++ {
			score := toolbox.Evaluation(n.population[i])
			n.population[i].SetFitness(score)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := gap * 3; i < len(n.population); i++ {
			score := toolbox.Evaluation(n.population[i])
			n.population[i].SetFitness(score)
		}
	}()

	wg.Wait()
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
		if !pass {
			// create a new species
			ns := NewSpecies(len(n.species), n.population[i])
			n.species = append(n.species, ns)
		}
	}

	// remove species with no members
	for i, niche := range n.species {
		if len(niche.members) == 0 {
			if len(n.species) < 1 {
				n.species = []*Species{}
			} else {
				n.species = append(n.species[:i], n.species[i+1:]...)
			}
		}
	}
}

// Run executes NEAT algorithm.
func (n *NEAT) Run(verbose bool) {
	for i := 0; i < param.NumGeneration; i++ {
		n.evaluate()
		n.speciate()

		// crossover and fitness share
		n.population = []*Genome{}
		for _, niche := range n.species {
			niche.FitnessShare()
			niche.VarMembers()
			niche.age++

			// add members that survived
			n.population = append(n.population, niche.members...)
		}
		sort.Sort(byFitness(n.population))
		fmt.Printf("Best: %f\n", n.population[0].fitness)
	}
	best := n.population[0]
	nnet := NewNetwork(best)

	inputs := make([]float64, 3)
	inputs[0] = 1.0
	// 0 xor 0
	inputs[1] = 0.0
	inputs[2] = 0.0
	output, err := nnet.Activate(inputs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("0 xor 0 = %f\n", output)

	// 0 xor 1
	inputs[1] = 0.0
	inputs[2] = 1.0
	output, err = nnet.Activate(inputs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("0 xor 1 = %f\n", output)

	// 1 xor 0
	inputs[1] = 1.0
	inputs[2] = 0.0
	output, err = nnet.Activate(inputs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1 xor 0 = %f\n", output)

	// 1 xor 1
	inputs[1] = 1.0
	inputs[2] = 1.0
	output, err = nnet.Activate(inputs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1 xor 1 = %f\n", output)
}
