package neat

// Genome is an implementation of genotype of an evolving network;
// it includes NodeGenes and ConnGenes.
type Genome struct {
	numSensors int // number of sensor nodes
	numOutputs int // number of output nodes
	numHidden  int // number of hidden nodes

	numNodes int // total number of nodes
	numConns int // total number of connections

	nodes []*NodeGene // collection of node genes
	conns []*ConnGene // collection of connection genes
}

// NewGenome creates a new genome in its initial state, it is
// only consist of fully connected sensor nodes and output nodes.
func NewGenome(numSensors, numOutputs int) *Genome {
	// number of nodes and connections including bias
	numNodes := numSensors + 1 + numOutputs
	numConns := (numSensors + 1) * numOutputs

	nodes := make([]*NodeGene, 0, numNodes)
	conns := make([]*ConnGene, 0, numConns)
	// sensor nodes
	for i := 0; i < numSensors; i++ {
		nodes = append(nodes, NewNode(i, "sensor"))
	}
	// output nodes and connections
	nodes = append(nodes, NewNode(numNodes-1, "bias"))
	for i := numSensor + 1; i < numNodes; i++ {
		nodes = append(nodes, NewNode(i, "output"))
		// connect from input nodes to this node
		for j := 0; j <= numSensors; j++ {
			conns = append(conns, NewConn(j, i))
		}
	}

	&Genome{
		numSensors: numSensors,
		numOutputs: numOutputs,
		numHidden:  0,
		numNodes:   numNodes,
		numConns:   numConns,
		nodes:      nodes,
		conns:      conns,
	}
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

// NumNodes returns the total number of nodes in the genome.
func (g *Genome) NumNodes() int {
	return g.numNodes
}

// NumConns returns the total number of connections in the genome.
func (g *Genome) NumConns() int {
	return g.numConns
}

// Nodes returns all nodes in the genome.
func (g *Genome) Nodes() []*NodeGene {
	return g.nodes
}

// Conns returns all connections in the genome.
func (g *Genome) Conns() []*ConnGene {
	return g.conns
}

// NodeGene is an implementation of each node within a genome.
// Each node includes a node ID (NID), a node type (NType), and
// a pointer to an activation function.
type NodeGene struct {
	nid   int             // node ID
	ntype string          // node type
	afn   *ActivationFunc // activation function
}

// NewNodeGene creates a new node gene with the given NID, node type, and
// a pointer to an activation function.
func NewNodeGene(nid int, ntype string, afn *ActivationFunc) *NodeGene {
	return &NodeGene{
		nid:   nid,
		ntype: ntype,
		afn:   afn,
	}
}

// NID returns the node's node ID (NID).
func (n *NodeGene) NID() int {
	return n.nid
}

// NType returns the node's node type (NType).
func (n *NodeGene) NType() int {
	return n.ntype
}

// Afn returns the node's activation function.
func (n *NodeGene) Afn() *ActivationFunc {
	return n.afn
}

// ConnGene is an implementation of each connection within a genome.
// It represents a connection between an in-node and an out-node;
// it contains an innovation number and nids of the in-node and the
// out-node.
type ConnGene struct {
	innov int // innovation number
	in    int // NID of in-node
	out   int // NID of out-node
}

// NewConnGene creates a new connection gene with the given innovation
// number, the in-node NID, and the out-node NID.
func NewConnGene(innov, in, out int) *ConnGene {
	return &ConnGene{
		innov: innov,
		in:    in,
		out:   out,
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
