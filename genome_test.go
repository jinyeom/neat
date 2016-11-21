package neat

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenome(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("=== Creating NEAT ===\n")
	_, err := New(&NEATConfig{
		numSensors: 3,
		numOutputs: 2,
		populationSize: 50,
		crossoverRate: 0.1,
		mutAddNodeRate: 0.1,
		mutAddConnRate: 0.1,
		mutWeightRate: 0.1,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("=== Creating a Genome ===\n")
	g, err := NewGenome(0, 3, 2)
	if err != nil {
		panic(err)
	}
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
			fmt.Printf("ENABLED  ")
		}
		fmt.Printf("Innov %d (%f): (%d) -> (%d)\n", c.Innov(), c.Weight(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)
	fmt.Println()

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
			fmt.Printf("ENABLED  ")
		}
		fmt.Printf("Innov %d (%f): (%d) -> (%d)\n", c.Innov(), c.Weight(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)
	fmt.Println()

	fmt.Printf("Mutate Add Connection\n")
	g.mutateAddConn()
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
			fmt.Printf("ENABLED  ")
		}
		fmt.Printf("Innov %d (%f): (%d) -> (%d)\n", c.Innov(), c.Weight(), c.In(), c.Out())
	}
	fmt.Printf("Current innovation number: %d\n", globalInnovNum)
	fmt.Println()
}
