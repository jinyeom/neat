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
	"math"
	"math/rand"
	"sort"
)

// Genome is an implementation of genotype of an evolving network; it is
// consist of a genome ID, number of sensors, number of outputs, number
// of hidden nodes, nodes and connections within the network, and the
// genome's fitness value.
type Genome struct {
	gid int // genome ID
	sid int // species ID

	nodes []*NodeGene // collection of node genes
	conns []*ConnGene // collection of connection genes

	fitness float64 // fitness value of the genome
}

// NewGenome creates a new genome in its initial state with only fully-connected
// sensor nodes and output nodes, and no hidden nodes.
func NewGenome(gid int) *Genome {
	// initial number of nodes and connections
	numNodes := param.NumSensors + param.NumOutputs
	numConns := param.NumSensors * param.NumOutputs
	nodes := make([]*NodeGene, 0, numNodes)
	conns := make([]*ConnGene, 0, numConns)

	for i := 0; i < param.NumSensors; i++ {
		n := NewNodeGene(i, "sensor", Identity())
		nodes = append(nodes, n)
	}
	for i := param.NumSensors; i < numNodes; i++ {
		n := NewNodeGene(i, "output", Sigmoid())
		nodes = append(nodes, n)
		// connect the new output node to all input nodes
		for j := 0; j < param.NumSensors; j++ {
			innov := innovations[[2]int{nodes[j].nid, n.nid}]
			if innov == 0 {
				innov = globalInnovNum
				// register the new connection innovation
				innovations[[2]int{nodes[j].nid, n.nid}] = innov
				globalInnovNum++
			}
			c := NewConnGene(innov, nodes[j].nid, n.nid, rand.NormFloat64())
			conns = append(conns, c)
		}
	}

	return &Genome{
		gid:     gid,
		sid:     0,
		nodes:   nodes,
		conns:   conns,
		fitness: 0.0,
	}
}

// GID returns the genome's ID.
func (g *Genome) GID() int {
	return g.gid
}

// SID returns the genome's species ID.
func (g *Genome) SID() int {
	return g.sid
}

// Nodes returns all nodes in the genome.
func (g *Genome) Nodes() []*NodeGene {
	return g.nodes
}

// Conns returns all connections in the genome.
func (g *Genome) Conns() []*ConnGene {
	return g.conns
}

// Node returns a node gene with the argument nid; returns nil if
// a node with the nid doesn't exist.
func (g *Genome) Node(nid int) *NodeGene {
	i := sort.Search(len(g.nodes), func(i int) bool {
		return g.nodes[i].nid >= nid
	})

	if i < len(g.nodes) && g.nodes[i].nid == nid {
		return g.nodes[i]
	}
	return nil
}

// Conn returns a connection gene with the arguement innovation
// number; returns nil if a connection with the innovation number
// doesn't exist.
func (g *Genome) Conn(innov int) *ConnGene {
	i := sort.Search(len(g.conns), func(i int) bool {
		return g.conns[i].innov >= innov
	})

	if i < len(g.conns) && g.conns[i].innov == innov {
		return g.conns[i]
	}
	return nil
}

// Distance returns the compatibility distance between this genome
// and the argument genome. The compatibility distance is a measurement
// of two genomes' compatibility for speciating them.
func (g *Genome) Distance(g1 *Genome) float64 {
	numDisjoint := 0     // number of disjoint genes
	numExcess := 0       // number of excess genes
	numMatch := 0        // number of matching genes
	avgWeightDiff := 0.0 // average weight differences of matching genes

	small := g  // genome with smaller max innov
	large := g1 // genome with larger max innov

	maxSmallInnov := small.conns[len(small.conns)-1].innov
	maxLargeInnov := large.conns[len(large.conns)-1].innov

	if maxSmallInnov > maxLargeInnov {
		small, large = large, small
		maxSmallInnov, maxLargeInnov = maxLargeInnov, maxSmallInnov
	}

	// try innovation numbers from 1 to the small genome's largest
	// innovation numbers to count the number of disjoint genes
	for i := 1; i <= maxSmallInnov; i++ {
		sc := small.Conn(i)
		lc := large.Conn(i)

		if sc != nil && lc != nil {
			avgWeightDiff += math.Abs(sc.weight - lc.weight)
			numMatch++
		} else if (sc != nil && lc == nil) || (sc == nil && lc != nil) {
			numDisjoint++
		}
	}

	// get average difference if the number of matching genes is
	// larger than 0
	if numMatch != 0 {
		avgWeightDiff /= float64(numMatch)
	}

	// count excess genes
	for i := maxSmallInnov + 1; i <= maxLargeInnov; i++ {
		if large.Conn(i) != nil {
			numExcess++
		}
	}

	return (param.CoeffExcess * float64(numExcess)) +
		(param.CoeffDisjoint * float64(numDisjoint)) +
		(param.CoeffWeight * avgWeightDiff)
}

// Crossover returns a child genome created by crossover operation
// between this genome and other genome provided as an argument; since
// the two parent genomes have to be in the same species, it is assumed
// that the child's species id is the same as one of the two parents'
// (g0 in this implementation).
func Crossover(g0, g1 *Genome, gid int) *Genome {
	child := &Genome{
		gid:     gid,
		sid:     g0.sid,
		nodes:   make([]*NodeGene, 0, len(g0.nodes)),
		conns:   make([]*ConnGene, 0, len(g0.conns)),
		fitness: 0.0,
	}

	small := g0 // genome with smaller max innov
	large := g1 // genome with larger max innov

	// sort connections by innovation numbers
	sort.Sort(byInnov(small.conns))
	sort.Sort(byInnov(large.conns))

	maxSmallInnov := small.conns[len(small.conns)-1].innov
	maxLargeInnov := large.conns[len(large.conns)-1].innov

	if maxSmallInnov > maxLargeInnov {
		small, large = large, small
		maxSmallInnov, maxLargeInnov = maxLargeInnov, maxSmallInnov
	}

	for i := 1; i <= maxSmallInnov; i++ {
		sc := small.Conn(i)
		lc := large.Conn(i)

		// matching/disjoint genes
		if sc != nil && lc != nil {
			if rand.Float64() < 0.5 {
				child.copyConn(small, sc)
			} else {
				child.copyConn(large, lc)
			}
		} else if sc != nil && lc == nil {
			child.copyConn(small, sc)
		} else if sc == nil && lc != nil {
			child.copyConn(large, lc)
		}
	}

	// excess genes
	for i := maxSmallInnov + 1; i <= maxLargeInnov; i++ {
		lc := large.Conn(i)
		if lc != nil {
			child.copyConn(large, lc)
		}
	}

	return child
}

// copyConn is a helper function of Crossover which copies a connection from
// other genome to this genome, and nodes that are connected by this
// connection, accordingly.
func (g *Genome) copyConn(g0 *Genome, c *ConnGene) {
	if g.Conn(c.innov) == nil {
		g.conns = append(g.conns, c)
	}
	if g.Node(c.in) == nil {
		g.nodes = append(g.nodes, g0.Node(c.in))
	}
	if g.Node(c.out) == nil {
		g.nodes = append(g.nodes, g0.Node(c.out))
	}
	sort.Sort(byNID(g.nodes))
	sort.Sort(byInnov(g.conns))
}

// Mutate mutates the genome by adding a node, adding a connection,
// and by mutating connections' weights.
func (g *Genome) Mutate() {
	// mutation by adding a new node; available only if there is at
	// least one connection in the genome.
	if rand.Float64() < param.MutAddNodeRate {
		g.mutateAddNode()
	}
	// mutation by adding a new connection.
	if rand.Float64() < param.MutAddConnRate {
		g.mutateAddConn()
	}
	// mutate connections
	for i := range g.conns {
		g.conns[i].mutate(param.MutWeightRate)
	}

	// sort nodes and connections
	sort.Sort(byNID(g.nodes))
	sort.Sort(byInnov(g.conns))
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

		// Create a new node that will be placed between a connection
		newNode := NewNodeGene(len(g.nodes), "hidden", Sigmoid())
		g.nodes = append(g.nodes, newNode)

		// The first connection that will be created by spliting an existing
		// connection will have a weight of 1.0, and will be connected from
		// the in-node of the existing node to the newly created node.
		innov := innovations[[2]int{oldIn, newNode.nid}]
		if innov == 0 {
			innov = globalInnovNum
			// register the new connection innovation
			innovations[[2]int{oldIn, newNode.nid}] = innov
			globalInnovNum++
		}
		newConn1 := NewConnGene(innov, oldIn, newNode.nid, 1.0)

		// The second new connection will have the same weight as the existing
		// connection, in order to prevent sudden changes after the mutation, and
		// will be connected from the new node to the out-node of the existing
		// connection.
		innov = innovations[[2]int{newNode.nid, oldOut}]
		if innov == 0 {
			innov = globalInnovNum
			// register the new connection innovation
			innovations[[2]int{newNode.nid, oldOut}] = innov
			globalInnovNum++
		}
		newConn2 := NewConnGene(innov, newNode.nid, oldOut, g.conns[ci].weight)

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
	out := rand.Intn(len(g.nodes[param.NumSensors:])) + param.NumSensors

	// Search for a connection gene that has the same in-node and out-node.
	for _, conn := range g.conns {
		if in == conn.in && out == conn.out {
			// then, do nothing and return.
			return
		}
	}

	// A new connection gene with a random weight is added between the
	// selected nodes. If the connection innovation already exists, use
	// the same innovation number as before; use global innovation number,
	// otherwise.
	innov := innovations[[2]int{g.nodes[in].nid, g.nodes[out].nid}]
	if innov == 0 {
		innov = globalInnovNum
		// register the new connection innovation
		innovations[[2]int{g.nodes[in].nid, g.nodes[out].nid}] = innov
		globalInnovNum++
	}
	g.conns = append(g.conns, NewConnGene(innov, in, out, rand.NormFloat64()))
}
