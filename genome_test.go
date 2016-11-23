package neat

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenome(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// Test creating a new genome
	fmt.Printf("=== Creating a Genome ===\n")
	g := NewGenome(0, &Param{
		NumSensors:     3,
		NumOutputs:     2,
		PopulationSize: 1,
		MutAddNodeRate: 0.1,
		MutAddConnRate: 0.1,
		MutWeightRate:  0.1,
	})
	fmt.Printf("GID: %d\n", g.ID())
	fmt.Printf("Nodes:\n")
	nodes := g.Nodes()
	for _, n := range nodes {
		fmt.Printf("NID %d (%s) - %s\n", n.ID(), n.NType(), n.Afn().Name())
	}
	fmt.Printf("Connections:\n")
	conns := g.Conns()
	for _, c := range conns {
		if c.disabled {
			fmt.Printf("DISABLED ")
		} else {
			fmt.Printf("ENABLED  ")
		}
		fmt.Printf("Innov %d (%f): (%d) -> (%d)\n", c.Innov(), c.Weight(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)
	fmt.Println()

	// Test mutation by adding a node
	fmt.Printf("=== Mutation by Adding Node ===\n")
	g.mutateAddNode()
	fmt.Printf("Nodes after mutation:\n")
	nodes = g.Nodes()
	for _, n := range nodes {
		fmt.Printf("NID %d (%s) - %s\n", n.ID(), n.NType(), n.Afn().Name())
	}
	fmt.Printf("Connections after mutation:\n")
	conns = g.Conns()
	for _, c := range conns {
		if c.disabled {
			fmt.Printf("DISABLED ")
		} else {
			fmt.Printf("ENABLED  ")
		}
		fmt.Printf("Innov %d (%f): (%d) -> (%d)\n", c.Innov(), c.Weight(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)
	fmt.Println()

	// Test mutation by adding a new connection
	fmt.Printf("=== Mutation by Adding Connection ===\n")
	g.mutateAddConn()
	fmt.Printf("Nodes after mutation:\n")
	nodes = g.Nodes()
	for _, n := range nodes {
		fmt.Printf("NID %d (%s) - %s\n", n.ID(), n.NType(), n.Afn().Name())
	}
	fmt.Printf("Connections after mutation:\n")
	conns = g.Conns()
	for _, c := range conns {
		if c.disabled {
			fmt.Printf("DISABLED ")
		} else {
			fmt.Printf("ENABLED  ")
		}
		fmt.Printf("Innov %d (%f): (%d) -> (%d)\n", c.Innov(), c.Weight(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)
	fmt.Println()

	// Test overall mutation of a genome
	fmt.Printf("Overall mutation\n")
	g.Mutate()
	fmt.Printf("Nodes after mutation:\n")
	nodes = g.Nodes()
	for _, n := range nodes {
		fmt.Printf("NID %d (%s) - %s\n", n.ID(), n.NType(), n.Afn().Name())
	}
	fmt.Printf("Connections after mutation:\n")
	conns = g.Conns()
	for _, c := range conns {
		if c.disabled {
			fmt.Printf("DISABLED ")
		} else {
			fmt.Printf("ENABLED  ")
		}
		fmt.Printf("Innov %d (%f): (%d) -> (%d)\n", c.Innov(), c.Weight(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)

}
