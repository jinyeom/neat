/*


genome.go implementation of the genome in NEAT.

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
	"math/rand"
)

// Genome is an implementation of genotype of an evolving network;
// it includes NodeGenes and ConnGenes.
type Genome struct {
	id int // genome ID

	numSensors int // number of sensor nodes
	numOutputs int // number of output nodes
	numHidden  int // number of hidden nodes

	nodes []*NodeGene // collection of node genes
	conns []*ConnGene // collection of connection genes

	fitness float64 // fitness value of the genome
}

// NewGenome creates a new genome in its initial state, it is
// only consist of fully connected sensor nodes and output nodes.
func NewGenome(id, numSensors, numOutputs int) (*Genome, error) {
	if numSensors < 1 || numOutputs < 1 {
		return nil, errors.New("Invalid number of sensors and/or outputs")
	}

	// initial innovation number
	initInnovNum := 0

	// number of nodes and connections including bias
	numNodes := numSensors + 1 + numOutputs
	numConns := (numSensors + 1) * numOutputs

	nodes := make([]*NodeGene, 0, numNodes)
	conns := make([]*ConnGene, 0, numConns)
	// sensor nodes
	for i := 0; i < numSensors; i++ {
		nodes = append(nodes, NewNodeGene(i, "sensor", Identity()))
	}
	// bias node as one of the sensors
	nodes = append(nodes, NewNodeGene(numSensors, "bias", Identity()))
	// output nodes and connections
	for i := numSensors + 1; i < numNodes; i++ {
		nodes = append(nodes, NewNodeGene(i, "output", Sigmoid()))
		// connect from input nodes to this node
		for j := 0; j <= numSensors; j++ {
			conns = append(conns, NewConnGene(initInnovNum, j, i))
			initInnovNum++
		}
	}

	return &Genome{
		id:         id,
		numSensors: numSensors,
		numOutputs: numOutputs,
		numHidden:  0,
		nodes:      nodes,
		conns:      conns,
	}, nil
}

// ID returns the genome's ID.
func (g *Genome) ID() int {
	return g.id
}

// NumSensors returns the number of sensor nodes in the genome.
func (g *Genome) NumSensors() int {
	return g.numSensors
}

// NumOutputs returns the number of output nodes in the genome.
func (g *Genome) NumOutputs() int {
	return g.numOutputs
}

// NumHidden returns the number of hidden nodes in the genome.
func (g *Genome) NumHidden() int {
	return g.numHidden
}

// Nodes returns all nodes in the genome.
func (g *Genome) Nodes() []*NodeGene {
	return g.nodes
}

// Conns returns all connections in the genome.
func (g *Genome) Conns() []*ConnGene {
	return g.conns
}

// Copy returns a deep copy of this genome.
func (g *Genome) Copy() *Genome {
	return &Genome{
		id:         g.id,
		numSensors: g.numSensors,
		numOutputs: g.numOutputs,
		numHidden:  g.numHidden,
		nodes: func() []*NodeGene {
			nodes := make([]*NodeGene, 0, len(g.nodes))
			for _, node := range g.nodes {
				nodes = append(nodes, node.Copy())
			}
			return nodes
		}(),
		conns: func() []*ConnGene {
			conns := make([]*ConnGene, 0, len(g.conns))
			for _, conn := range g.conns {
				conns = append(conns, conn.Copy())
			}
			return conns
		}(),
	}
}

// Crossover returns children genome created by crossover operation
// between this genome and other genome provided as an argument.
func Crossover(g0, g1 *Genome) (*Genome, *Genome) {
	child1 := g0.Copy()
	child2 := g1.Copy()

	return child1, child2
}

// Mutate mutates the genome by adding a node, adding a connection,
// and by mutating connections' weights.
func (g *Genome) Mutate(conf *Config) {
	if rand.Float64() < conf.MutAddNodeRate {
		g.mutateAddNode()
	}
	if rand.Float64() < conf.MutAddConnRate {
		g.mutateAddConn()
	}
	// mutate connections
	for i := range g.conns {
		g.conns[i].mutate(conf.MutWeightRate)
	}
}

// mutateAddNode mutates the genome by adding a node between a
// connection of two nodes.
func (g *Genome) mutateAddNode() {
	ci := rand.Intn(len(g.conns))
	oldIn := g.conns[ci].In()
	oldOut := g.conns[ci].Out()

	newNode := NewNodeGene(len(g.nodes), "hidden", Sigmoid())
	g.nodes = append(g.nodes, newNode)

	newConn1 := NewConnGene(globalInnovNum, oldIn, newNode.id)
	newConn2 := NewConnGene(globalInnovNum+1, newNode.id, oldOut)
	g.conns = append(g.conns, newConn1)
	g.conns = append(g.conns, newConn2)
	globalInnovNum += 2

	g.conns[ci].switchConn()
	g.numHidden++
}

// mutateAddConn mutates the genome by adding a connection between
// two nodes. A new connection can be added from a node to itself.
func (g *Genome) mutateAddConn() {
	in := rand.Intn(len(g.nodes))
	out := rand.Intn(len(g.nodes[g.numSensors+1:])) + g.numSensors + 1

	newConn := NewConnGene(globalInnovNum, in, out)
	g.conns = append(g.conns, newConn)
	globalInnovNum++
}
