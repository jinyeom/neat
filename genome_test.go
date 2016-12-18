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

	p := &Param{
		NumSensors:     3,
		NumOutputs:     2,
		PopulationSize: 1,
		MutAddNodeRate: 0.1,
		MutAddConnRate: 0.1,
		MutWeightRate:  0.1,
		CoeffExcess:    0.3,
		CoeffDisjoint:  0.3,
		CoeffWeight:    0.3,
	}

	// Test creating a new genome
	fmt.Printf("=== Creating a Genome ===\n")
	g := NewGenome(0, p)
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

	// Test overall mutation of a genome
	fmt.Printf("Overall mutation\n")
	g.Mutate()
	genomeStatus(g)

	// Test compatibility
	fmt.Printf("=== Compatibility distance test ===\n")
	g0 := NewGenome(1, p)
	g1 := NewGenome(2, p)
	genomeStatus(g0)
	genomeStatus(g1)
	fmt.Printf("Compatibility of g0 and g1: %f\n", g0.Compatibility(g1))
}
