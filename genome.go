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

	// initialize sensor nodes and output nodes
	nodes := make([]*NodeGene, 0, numNodes)
	conns := make([]*ConnGene, 0, numConns)
	for i := 0; i < numSensor; i++ {
		nodes = append(nodes, NewNode(i, "sensor"))
	}
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

type NodeGene struct {
	nid int // node ID
}

// NewNodeGene creates a new node gene with the given NID.
func NewNodeGene(nid int, ntype string) *NodeGene {
	return &NodeGene{
		nid:   nid,
		ntype: ntype,
	}
}

type ConnGene struct {
	innov   int // innovation number
	in, out int // in node and out node
}
