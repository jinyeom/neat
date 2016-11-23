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
	"math/rand"
)

// Genome is an implementation of genotype of an evolving network; it is
// consist of a genome ID, number of sensors, number of outputs, number
// of hidden nodes, nodes and connections within the network, and the
// genome's fitness value.
type Genome struct {
	id int // genome ID

	param *Param // parameters of NEAT

	nodes []*NodeGene // collection of node genes
	conns []*ConnGene // collection of connection genes

	fitness float64 // fitness value of the genome
}

// NewGenome creates a new genome in its initial state; it only creates
// sensor nodes and output nodes with no connections. Connections are
// expected to be added via mutations as evolution progresses.
func NewGenome(id int, param *Param) *Genome {
	// number of nodes and connections including bias
	numNodes := param.NumSensors + param.NumOutputs
	nodes := make([]*NodeGene, 0, numNodes)

	for i := 0; i < param.NumSensors; i++ {
		nodes = append(nodes, NewNodeGene(i, "sensor", Identity()))
	}
	for i := param.NumSensors; i < numNodes; i++ {
		nodes = append(nodes, NewNodeGene(i, "output", Sigmoid()))
	}

	return &Genome{
		id:      id,
		param:   param,
		nodes:   nodes,
		conns:   make([]*ConnGene, 0),
		fitness: 0.0,
	}
}

// ID returns the genome's ID.
func (g *Genome) ID() int {
	return g.id
}

// NumHiddenNodes returns the number of hidden nodes in the genome.
func (g *Genome) NumHiddenNodes() int {
	return len(g.nodes) - (g.param.NumSensors + g.param.NumOutputs)
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
		id:    g.id,
		param: g.param,
		// deep copy of nodes
		nodes: func() []*NodeGene {
			nodes := make([]*NodeGene, 0, len(g.nodes))
			for _, node := range g.nodes {
				nodes = append(nodes, node.Copy())
			}
			return nodes
		}(),
		// deep copy of connections
		conns: func() []*ConnGene {
			conns := make([]*ConnGene, 0, len(g.conns))
			for _, conn := range g.conns {
				conns = append(conns, conn.Copy())
			}
			return conns
		}(),
		fitness: g.fitness,
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
func (g *Genome) Mutate() {
	// mutation by adding a new node; available only if there is at
	// least one connection in the genome.
	if rand.Float64() < g.param.MutAddNodeRate {
		g.mutateAddNode()
	}
	// mutation by adding a new connection.
	if rand.Float64() < g.param.MutAddConnRate {
		g.mutateAddConn()
	}
	// mutate connections
	for i := range g.conns {
		g.conns[i].mutate(g.param.MutWeightRate)
	}
}

// mutateAddNode mutates the genome by adding a node between a
// connection of two nodes. After the newly added node splits the
// existing connection, two new connections will be added with weights
// of 1.0 and the original connection's weight, in order to prevent
// sudden changes in the network's performance.
func (g *Genome) mutateAddNode() {
	if len(g.conns) > 0 {
		ci := rand.Intn(len(g.conns))
		oldIn := g.conns[ci].In()
		oldOut := g.conns[ci].Out()

		// Create a new node that will be placed between a connection of
		// two nodes.
		newNode := NewNodeGene(globalInnovNum, "hidden", Sigmoid())
		g.nodes = append(g.nodes, newNode)
		globalInnovNum++

		// The first connection that will be created by spliting an existing
		// connection will have a weight of 1.0, and will be connected from
		// the in-node of the existing node to the newly created node.
		newConn1 := NewConnGene(globalInnovNum, oldIn, newNode.id, 1.0)
		globalInnovNum++

		// The second new connection will have the same weight as the existing
		// connection, in order to prevent sudden changes after the mutation, and
		// will be connected from the new node to the out-node of the existing
		// connection.
		newConn2 := NewConnGene(globalInnovNum, newNode.id, oldOut, g.conns[ci].weight)
		globalInnovNum++

		g.conns = append(g.conns, newConn1)
		g.conns = append(g.conns, newConn2)

		// The original connection gene is now disabled.
		g.conns[ci].switchConn()
	}
}

// mutateAddConn mutates the genome by adding a connection between
// two nodes. A new connection can be added from a node to itself.
func (g *Genome) mutateAddConn() {
	// The in-node of the connection to be added can be selected
	// randomly from any node genes.
	in := rand.Intn(len(g.nodes))

	// The out-node can only be randomly selected from nodes that are
	// not sensor nodes.
	out := rand.Intn(len(g.nodes[g.param.NumSensors:])) + g.param.NumSensors

	// Search for a connection gene that has the same in-node and out-node.
	for _, conn := range g.conns {
		if in == conn.in && out == conn.out {
			// then, do nothing and return.
			return
		}
	}

	// A new connection gene with a random weight is added between the
	// selected nodes.
	newConn := NewConnGene(globalInnovNum, in, out, rand.NormFloat64())
	g.conns = append(g.conns, newConn)
	globalInnovNum++
}
