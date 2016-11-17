package neat

// Genome is an implementation of genotype of an evolving network;
// it includes NodeGenes and ConnGenes.
type Genome struct {
	numSensor int // number of sensor nodes
	numOutput int // number of output nodes
	numHidden int // number of hidden nodes

	numNodes int // total number of nodes
	numConns int // total number of connections

	nodes []*NodeGene // collection of node genes
	conns []*ConnGene // collection of connection genes
}

// NewGenome creates a new genome in its initial state, it is
// only consist of fully connected sensor nodes and output nodes.
func NewGenome(numSensor, numOutput int) *Genome {
	// number of nodes and connections including bias
	numNodes := numSensor + 1 + numOutput
	numConns := (numSensor + 1) * numOutput

	nodes := make([]*NodeGene, 0, numNodes)
	conns := make([]*ConnGene, 0, numConns)
	// sensor nodes
	for i := 0; i < numSensor; i++ {
		nodes = append(nodes, NewNode(i, "sensor"))
	}
	// output nodes and connections
	nodes = append(nodes, NewNode(numNodes-1, "bias"))
	for i := numSensor + 1; i < numNodes; i++ {
		nodes = append(nodes, NewNode(i, "output"))
		// connect from input nodes to this node
		for j := 0; j <= numSensor; j++ {
			conns = append(conns, NewConn(j, i))
		}
	}

	&Genome{
		numSensor: numSensor,
		numOutput: numOutput,
		numHidden: 0,
		numNodes:  numNodes,
		numConns:  numConns,
		nodes:     nodes,
		conns:     conns,
	}
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
