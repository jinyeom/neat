package neat

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// genomeStatus prints the current state of the argument genome, followed
// by the current global innovation number.
func genomeStatus(g *Genome) {
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
}

func TestGenome(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	Init(&Param{
		NumSensors:     3,
		NumOutputs:     2,
		PopulationSize: 1,
		MutAddNodeRate: 0.1,
		MutAddConnRate: 0.1,
		MutWeightRate:  0.1,
		CoeffExcess:    0.3,
		CoeffDisjoint:  0.3,
		CoeffWeight:    0.3,
	}, &Toolbox{
		Activation: NEATSet(),
		Comparison: DirectCompare(),
		Selection:  TSelect(DirectCompare()),
		Evaluation: XORTest(),
	})

	// Test creating a new genome
	fmt.Printf("=== Creating a Genome ===\n")
	g := NewGenome(0)
	genomeStatus(g)

	// Test mutation by adding a new connection
	fmt.Printf("=== Mutation by Adding Connection ===\n")
	g.mutateAddConn()
	genomeStatus(g)

	// Test mutation by adding a node
	fmt.Printf("=== Mutation by Adding Node ===\n")
	g.mutateAddNode()
	fmt.Printf("Nodes after mutation:\n")
	genomeStatus(g)

	// Test overall mutation of a genome (power test)
	fmt.Printf("Overall mutation (power test)\n")
	for i := 0; i < 20; i++ {
		g.Mutate()
	}
	genomeStatus(g)

	// Test compatibility
	fmt.Printf("=== Compatibility distance test ===\n")
	g0 := NewGenome(1)
	g1 := NewGenome(2)
	for i := 0; i < 5; i++ {
		g0.Mutate()
		g1.Mutate()
	}
	genomeStatus(g0)
	genomeStatus(g1)
	fmt.Printf("Compatibility distance of g0 and g1: %f\n", g0.Distance(g1))

	// Test crossover
	fmt.Printf("=== Crossover test ===\n")
	for i := 0; i < 10; i++ {
		g0.Mutate()
		g1.Mutate()
	}
	genomeStatus(g0)
	genomeStatus(g1)
	child := Crossover(g0, g1, 3)
	genomeStatus(child)

}
