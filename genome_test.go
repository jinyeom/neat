package neat

import (
	"fmt"
	"math/rand"
	"testing"
)

func GenomeUnitTest() {
	fmt.Println("===== Genome Unit Test =====")

	fmt.Println("\x1b[32m=Testing creating a new genome...\x1b[0m")
	g0 := NewGenome(0, 3, 5, 0.0)
	fmt.Println(g0.String())

	fmt.Println("\x1b[32m=Testing mutation...\x1b[0m")
	Mutate(g0, 1.0, 1.0, 1.0)
	fmt.Println(g0.String())

	fmt.Println("\x1b[32m=Testing crossover...\x1b[0m")

	g1 := NewGenome(1, 3, 1, 0.0)
	Mutate(g1, 1.0, 1.0, 1.0)
	Mutate(g1, 1.0, 1.0, 1.0)
	fmt.Println("Parent 1:")
	fmt.Println(g1.String())

	g2 := NewGenome(2, 3, 1, 0.0)
	Mutate(g2, 1.0, 1.0, 1.0)
	Mutate(g2, 1.0, 1.0, 1.0)
	fmt.Println("Parent 2:")
	fmt.Println(g2.String())

	g3 := Crossover(3, g1, g2, 0.0)
	fmt.Println("Child:")
	fmt.Println(g3.String())

	fmt.Println("\x1b[32m=Testing compatibility distance...\x1b[0m")
	g4 := NewGenome(4, 3, 1, 0.0)
	g5 := NewGenome(5, 3, 1, 0.0)

	// before mutation (they should be fairly compatible)
	fmt.Println(g4.String())
	fmt.Println(g5.String())
	fmt.Printf("Compatibility distance: %f\n", Compatibility(g4, g5, 1.0, 1.0))

	// after 1 mutation (should be less compatible)
	Mutate(g4, 1.0, 1.0, 1.0)
	Mutate(g5, 1.0, 1.0, 1.0)

	fmt.Println(g4.String())
	fmt.Println(g5.String())
	fmt.Printf("Compatibility distance: %f\n", Compatibility(g4, g5, 1.0, 1.0))

	// after 2 mutation (should be less compatible)
	Mutate(g4, 1.0, 1.0, 1.0)
	Mutate(g5, 1.0, 1.0, 1.0)

	fmt.Println(g4.String())
	fmt.Println(g5.String())
	fmt.Printf("Compatibility distance: %f\n", Compatibility(g4, g5, 1.0, 1.0))

	// after 3 mutation (should be less compatible)
	Mutate(g4, 1.0, 1.0, 1.0)
	Mutate(g5, 1.0, 1.0, 1.0)

	fmt.Println(g4.String())
	fmt.Println(g5.String())
	fmt.Printf("Compatibility distance: %f\n", Compatibility(g4, g5, 1.0, 1.0))

	/*
		fmt.Println("\x1b[32m=Testing JSON export...\x1b[0m")
		if err := g1.ExportJSON(); err != nil {
			log.Fatal(err)
		}
	*/

}

func TestGenome(t *testing.T) {
	rand.Seed(0)
	GenomeUnitTest()
}
