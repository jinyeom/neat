package neat

// Phenotype is an interface of
type Phenotype interface {
	// ForwardPropagate takes inputs, propagates signals and returns
	// output in a slice of float64.
	ForwardPropagate(input []float64) []float64
}

// Node implements a node in a phenotype network; it includes a node ID,
// its activation function, and a signal value that the node holds.
type Node struct {
	nid       int               // node ID
	connNodes []*Node           // nodes connected to this node
	weights   map[*Node]float64 // connection weights mapping
	signal    float64           // stored activation signal
	afn       *ActivationFunc   // activation function
}

// NewNode encodes the arguement node gene, and creates a new node.
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
	for i, node := range n.connNodes {
		sum += node.signal * n.weights[node]
	}
	n.signal = n.afn.fn(sum)
	return n.signal
}

// NN implements Phenotype interface as a neural network.
type NN struct {
	neurons []*Node
}

// NewNN encodes a genome into a neural network (phenotype).
func NewNN(g *Genome) *NN {
	nodes := make([]*Node, len(g.nodes))
	for i := range nodes {
		nodes[i] = NewNode(g.nodes[i])
	}
	// connect nodes
	for _, conn := range g.conns {
	}
	return &NN{
		neurons: nodes,
	}
}

// ForwardPropagate
func (n *NN) ForwardPropagate(intput []float64) []float64 {

}
