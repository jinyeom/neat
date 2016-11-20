package neat

import (
	"fmt"
	"testing"
)

func TestGenome(t *testing.T) {
	fmt.Printf("=== Creating a Genome ===\n")
	g := NewGenome(0, 3, 2)
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
		fmt.Printf("Innov %d: (%d) -> (%d)\n", c.Innov(), c.In(), c.Out())
	}
	fmt.Printf("Current node counter: %d\n", g.ncount)
	fmt.Printf("Current innovation number: %d\n", g.innov)

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
		fmt.Printf("Innov %d: (%d) -> (%d)\n", c.Innov(), c.In(), c.Out())
	}
	fmt.Printf("Current node counter: %d\n", g.ncount)
	fmt.Printf("Current innovation number: %d\n", g.innov)
}
