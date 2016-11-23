/*


gene.go implementation of node and connection genes in NEAT.

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

// NodeGene is an implementation of each node within a genome.
// Each node includes a node ID (NID), a node type (NType), and
// a pointer to an activation function.
type NodeGene struct {
	id    int             // node ID
	ntype string          // node type
	afn   *ActivationFunc // activation function
}

// NewNodeGene creates a new node gene with the given NID, node type, and
// a pointer to an activation function.
func NewNodeGene(id int, ntype string, afn *ActivationFunc) *NodeGene {
	return &NodeGene{
		id:    id,
		ntype: ntype,
		afn:   afn,
	}
}

// ID returns the node's node ID.
func (n *NodeGene) ID() int {
	return n.id
}

// NType returns the node's node type (NType).
func (n *NodeGene) NType() string {
	return n.ntype
}

// Afn returns the node's activation function.
func (n *NodeGene) Afn() *ActivationFunc {
	return n.afn
}

// Copy returns a deep copy of this gene.
func (n *NodeGene) Copy() *NodeGene {
	return &NodeGene{
		id:    n.id,
		ntype: n.ntype,
		afn:   n.afn,
	}
}

// ConnGene is an implementation of each connection within a genome.
// It represents a connection between an in-node and an out-node;
// it contains an innovation number and nids of the in-node and the
// out-node, whether if the connection is disabled, and the weight
// of the connection.
type ConnGene struct {
	innov    int     // innovation number
	in       int     // NID of in-node
	out      int     // NID of out-node
	disabled bool    // whether if the connection is true
	weight   float64 // weight of connection
}

// NewConnGene creates a new connection gene with the given innovation
// number, the in-node NID, and the out-node NID.
func NewConnGene(innov, in, out int, weight float64) *ConnGene {
	return &ConnGene{
		innov:    innov,
		in:       in,
		out:      out,
		disabled: false,
		weight:   weight,
	}
}

// Innov returns the connection's innovation number.
func (c *ConnGene) Innov() int {
	return c.innov
}

// In returns the NID of in-node of the connection.
func (c *ConnGene) In() int {
	return c.in
}

// Out returns the NID of out-node of the connection.
func (c *ConnGene) Out() int {
	return c.out
}

// IsDisabled indicates whether the connection is disabled.
func (c *ConnGene) IsDisabled() bool {
	return c.disabled
}

// Weight returns the connection's weight.
func (c *ConnGene) Weight() float64 {
	return c.weight
}

// Copy returns a deep copy of this gene.
func (c *ConnGene) Copy() *ConnGene {
	return &ConnGene{
		innov:    c.innov,
		in:       c.in,
		out:      c.out,
		disabled: c.disabled,
		weight:   c.weight,
	}
}

// mutate mutates the connection weight.
func (c *ConnGene) mutate(mutWeightRate float64) {
	if rand.Float64() < mutWeightRate {
		c.weight += rand.NormFloat64()
	}
}

// switchConn enables the connection if it is disabled; disables
// the connection otherwise.
func (c *ConnGene) switchConn() {
	c.disabled = !c.disabled
}
