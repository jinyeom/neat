package neat

import (
	"fmt"
	"testing"
)

func TestGenome(t *testing.T) {
	fmt.Printf("=== Creating a Genome ===\n")
	g := NewGenome(3, 2)
	// check if the genome is initialized correctly
	fmt.Printf("Nodes:\n")
	nodes := g.Nodes()
	for _, n := range nodes {
		fmt.Printf("NID %d (%s)\n", n.NID(), n.NType())
	}
	fmt.Printf("Connections:\n")
	conns := g.Conns()
	for _, c := range conns {
		fmt.Printf("Innov %d: (%d) -> (%d)\n", c.Innov(), c.In(), c.Out())
	}
}
