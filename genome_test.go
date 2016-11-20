package neat

import (
	"fmt"
	"testing"
)

func TestGenome(t *testing.T) {
	fmt.Printf("=== Creating a Genome ===\n")
	g, err := NewGenome(0, 3, 2)
	if err != nil {
		panic(err)
	}
	// check if the genome is initialized correctly
	fmt.Printf("GID: %d\n", g.GID())
	fmt.Printf("Nodes:\n")
	nodes := g.Nodes()
	for _, n := range nodes {
		fmt.Printf("NID %d (%s) - %s\n", n.NID(), n.NType(), n.Afn().Name())
	}
	fmt.Printf("Connections:\n")
	conns := g.Conns()
	for _, c := range conns {
		if c.disabled {
			fmt.Printf("DISABLED ")
		} else {
			fmt.Printf("ENABLED ")
		}
		fmt.Printf("Innov %d: (%d) -> (%d)\n", c.Innov(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)

	fmt.Printf("=== Genome Mutation ===\n")

	fmt.Printf("Mutate Add Node\n")
	g.mutateAddNode()
	fmt.Printf("Nodes after mutation:\n")
	nodes = g.Nodes()
	for _, n := range nodes {
		fmt.Printf("NID %d (%s) - %s\n", n.NID(), n.NType(), n.Afn().Name())
	}
	fmt.Printf("Connections after mutation:\n")
	conns = g.Conns()
	for _, c := range conns {
		if c.disabled {
			fmt.Printf("DISABLED ")
		} else {
			fmt.Printf("ENABLED ")
		}
		fmt.Printf("Innov %d: (%d) -> (%d)\n", c.Innov(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)
}
