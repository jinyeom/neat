package neat

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func PrintNetwork(n *Network) {
	for _, node := range n.nodes {
		fmt.Printf("Node %d (%s) connected from: %d\n",
			node.nid, node.afn.name, func() []int {
				nids := make([]int, len(node.connNodes))
				for j := range nids {
					nids[j] = node.connNodes[j].nid
				}
				return nids
			}(),
		)
	}
	fmt.Println()
}

func TestNetwork(t *testing.T) {
	// Test creating a new network
	fmt.Printf("=== Creating a Network ===\n")
	g := NewGenome(0)
	for i := 0; i < 20; i++ {
		g.Mutate()
	}
	genomeStatus(g)

	n := NewNetwork(g)
	PrintNetwork(n)

	// Test forward propagating
	fmt.Printf("=== Forward Propagating ===\n")
	inputs := make([]float64, param.NumSensors)
	for i := range inputs {
		inputs[i] = rand.Float64()
	}
	outputs, err := n.ForwardPropagate(inputs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("output: %f\n", outputs)
}
