/*


network.go implementation of phenotype interface.

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

// Node implements a node in a phenotype network; it includes a node ID,
// its activation function, and a signal value that the node holds.
type Node struct {
	nid       int               // node ID
	connNodes []*Node           // nodes connected to this node
	weights   map[*Node]float64 // connection weights mapping
	signal    float64           // stored activation signal
	afn       *ActivationFunc   // activation function
}

// NewNode decodes the arguement node gene, and creates a new node.
func NewNode(n *NodeGene) *Node {
	return &Node{
		nid:       n.nid,
		connNodes: make([]*Node, 0),
		weights:   make(map[*Node]float64),
		signal:    0.0,
		afn:       n.afn,
	}
}

// Output sets and returns the signal of this node after it
// activates via its activation function.
func (n *Node) Output() float64 {
	sum := 0.0
	for _, node := range n.connNodes {
		sum += node.signal * n.weights[node]
	}
	n.signal = n.afn.fn(sum)
	return n.signal
}

// Network is the phenotype in NEAT, which decodes from a genome.
// A network can be used as a neural network, CPPN, etc.
type Network struct {
	nodes []*Node
}

// NewNetwork decodes a genome into a network (phenotype).
func NewNetwork(g *Genome) *Network {
	nodes := make([]*Node, len(g.nodes))
	for i := range g.nodes {
		nodes[i] = NewNode(g.nodes[i])
	}

	for _, conn := range g.conns {
		// connect the two nodes
		if !conn.disabled {
			// to be implemented
		}
	}

	return &Network{
		nodes: nodes,
	}
}

// ForwardPropagate
func (n *Network) ForwardPropagate(intput []float64) []float64 {

	// to be implemented

	return nil
}
