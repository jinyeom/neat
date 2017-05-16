package neat

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"testing"
)

func NEATUnitTest() {
	fmt.Println("===== NEAT Unit Test =====")

	fmt.Println("\x1b[32m=Testing config JSON file import...\x1b[0m")
	config, err := NewConfigJSON("config.json")
	if err != nil {
		fmt.Println("\x1b[31mFAIL\x1b[0m")
	}
	config.Summarize()

	fmt.Println("\x1b[32m=Testing creating and running NEAT...\x1b[0m")
	New(config, func(n *NeuralNetwork) float64 {
		score := 0.0

		inputs := make([]float64, 3)
		inputs[0] = 1.0 // bias

		// 0 xor 0
		inputs[1] = 0.0
		inputs[2] = 0.0
		output, err := n.FeedForward(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 0.0), 2.0)

		// 0 xor 1
		inputs[1] = 0.0
		inputs[2] = 1.0
		output, err = n.FeedForward(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 1.0), 2.0)

		// 1 xor 0
		inputs[1] = 1.0
		inputs[2] = 0.0
		output, err = n.FeedForward(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 1.0), 2.0)

		// 1 xor 1
		inputs[1] = 1.0
		inputs[2] = 1.0
		output, err = n.FeedForward(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 0.0), 2.0)

		return score
	}).Run(true)

	//fmt.Println("=Testing evaluation in sequence...")
	//for _, genome := range n.Population {
	//	Mutate(genome, 1.0, 1.0, 1.0)
	//	Mutate(genome, 1.0, 1.0, 1.0)
	//}
	//n.evaluateSequential()
	//for _, genome := range n.Population {
	//	fmt.Printf("Genome %d fitness: %.3f\n", genome.ID, genome.Fitness)
	//}
}

func TestNEAT(t *testing.T) {
	rand.Seed(0)
	NEATUnitTest()
}
