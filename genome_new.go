package neat

import (
	"math/rand"
)

type NodeType int

const (
	Sensor NodeType = iota
	Output
	Hidden
)

type ActivationFunc int

const (
	Linear ActivationFunc = iota
	Sigmoid
	Tanh
	ReLU
	Sine
	Gaussian
)

type NodeGene struct {
	nodeId     int
	nodeType   NodeType
	activation ActivationFunc
}

func NewNodeGene(nodeId int, nodeType NodeType, activation ActivationFunc) *NodeGene {
	return &NodeGene{
		nodeId:     nodeId,
		nodeType:   nodeType,
		activation: activation,
	}
}

func (n *NodeGene) Id() int {
	return n.nodeId
}

func (n *NodeGene) Type() NodeType {
	return n.nodeType
}

func (n *NodeGene) Activation() ActivationFunc {
	return n.activation
}

func (n *NodeGene) String() string {
	str := fmt.Sprintf("Node(%d, ", n.nodeId)
	// node type
	switch n.nodeType {
	case Sensor:
		str += "Sensor, "
	case Output:
		str += "Output, "
	case Hidden:
		str += "Hidden, "
	}
	// activation function
	switch n.activation {
	case Linear:
		str += "Linear"
	case Sigmoid:
		str += "Sigmoid"
	case Tanh:
		str += "Tanh"
	case ReLU:
		str += "ReLU"
	case Sine:
		str += "Sine"
	case Gaussian:
		str += "Gaussian"
	}
	str += ")"
	return str
}

type ConnectionGene struct {
	innovId   int       // innovation number
	src       *NodeGene // source node
	dst       *NodeGene // destination node
	weight    float64   // connection weight
	expressed bool      // connection expression
}

func NewConnectionGene(innovID int, src, dst *NodeGene) *ConnectionGene {
	return &ConnectionGene{
		innovID:   innovID,
		src:       src,
		dst:       dst,
		weight:    randWeight(),
		expressed: true,
	}
}

func randWeight() float64 {
	return rand.NormFloat64() * 3.0
}

func (c *ConnectionGene) Id() int {
	return c.innovId
}

func (c *ConnectionGene) Src() *NodeGene {
	return c.src
}

func (c *ConnectionGene) Dst() *NodeGene {
	return c.dst
}

func (c *ConnectionGene) Weight() float64 {
	return c.weight
}

func (c *ConnectionGene) Expressed() bool {
	return c.expressed
}

func (c *ConnectionGene) Enable() {
	c.expressed = true
}

func (c *ConnectionGene) Disable() {
	c.expressed = false
}

func (c *ConnectionGene) String() string {
	srcId := c.src.Id()
	dstId := c.dst.Id()
	if c.expressed {
		return fmt.Sprintf("[%d]--{%f}-->[%d]", srcId, c.weight, dstId)
	}
	return fmt.Sprintf("[%d]--/ /-->[%d]", srcId, dstId)
}

type Genome struct {
	id              int
	nodeGenes       []*NodeGene
	connectionGenes []*ConnectionGene
}

func NewGenome(id int) *Genome {
	return &Genome{
		id:              id,
		nodeGenes:       make([]*NodeGene),
		connectionGenes: make([]*ConnectionGene),
	}
}

func (g *Genome) String() string {
	numNodes := len(g.nodeGenes)
	numConnections := len(g.conncectionGenes)
	strs := make([]string, 0, numNodes+numConnections)
	for i := 0; i < numNodes; i++ {
		strs = append(strs, g.nodeGenes[i].String())
	}
	for i := 0; i < numConnections; i++ {
		strs = append(strs, g.connectionGenes[i].String())
	}
	return strings.Join(strs, "\n")
}

func (g *Genome) NodeGenes() []*NodeGene {
	return g.nodeGenes
}

func (g *Genome) ConnenctionGenes() []*ConnectionGene {
	return g.connectionGenes
}

func (g *Genome) PushNode() {
	g.nodeGene = append(g.nodeGene, n)
}
