package neat

import (
	"fmt"
	"log"
	"testing"
)

func TestNEAT(t *testing.T) {
	p, err := NewParam("example_param.np")
	if err != nil {
		log.Fatal(err)
	}

	tb := &Toolbox{
		NEATSet(),
		DirectCompare(),
		XORTest(),
	}

	Init(p, tb)

	fmt.Printf("=== Creating NEAT ===\n")
	n, err := New()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("=== Evalutation of genomes ===\n")
	fmt.Printf("Before evaluation:\n")
	for i, genome := range n.population {
		fmt.Printf("GENOME %d: %f\n", i, genome.fitness)
	}
	n.evaluate()
	fmt.Printf("After evaluation:\n")
	for i, genome := range n.population {
		fmt.Printf("GENOME %d: %f\n", i, genome.fitness)
	}

	fmt.Printf("=== Speciation Test ===\n")
	n.speciate()
	for i, niche := range n.species {
		fmt.Printf("SPECIES %d:\n", i)
		fmt.Printf("Repr.: GENOME %d\n", niche.representative.gid)
		for j, member := range niche.members {
			fmt.Printf("Member %d: GENOME %d\n", j, member.gid)
		}
	}
}
