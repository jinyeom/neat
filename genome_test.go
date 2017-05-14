package neat

import (
	"fmt"
	"math/rand"
	"testing"
)

func GenomeUnitTest() {
	fmt.Println("===== Genome Unit Test =====")

	fmt.Println("=Testing creating a new genome...")
	g0 := NewGenome(0, 3, 5)
	fmt.Println(g0.String())

	fmt.Println("=Testing mutation...")
	Mutate(g0, 1.0, 1.0, 1.0)
	fmt.Println(g0.String())

	fmt.Println("=Testing crossover...")

	// parent 1
	g1 := NewGenome(1, 3, 1)
	Mutate(g1, 1.0, 1.0, 1.0)
	Mutate(g1, 1.0, 1.0, 1.0)
	fmt.Println("Parent 1:")
	fmt.Println(g1.String())

	// parent 2
	g2 := NewGenome(2, 3, 1)
	Mutate(g2, 1.0, 1.0, 1.0)
	Mutate(g2, 1.0, 1.0, 1.0)
	fmt.Println("Parent 2:")
	fmt.Println(g2.String())

	// child
	g3 := Crossover(3, g1, g2)
	fmt.Println("Child:")
	fmt.Println(g3.String())
}

func TestGenome(t *testing.T) {
	rand.Seed(0)
	GenomeUnitTest()
}
