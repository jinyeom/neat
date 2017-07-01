package neat

import (
	"fmt"
	"math/rand"
	"testing"
)

func NEATUnitTest() {
	fmt.Println("===== NEAT Unit Test =====")

	fmt.Println("\x1b[32m=Testing config JSON file import...\x1b[0m")
	configXOR, err := NewConfigJSON("config_xor.json")
	if err != nil {
		fmt.Println("\x1b[31mFAIL\x1b[0m")
	}
	configXOR.Summarize()

	fmt.Println("\x1b[32m=Testing NEAT with XOR test...\x1b[0m")
	n0 := New(configXOR, XORTest())
	n0.Run()

	nn := NewNeuralNetwork(n0.Best)
	output, _ := nn.FeedForward([]float64{1.0, 1.0, 1.0})
	fmt.Println(output)
	output, _ = nn.FeedForward([]float64{1.0, 0.0, 1.0})
	fmt.Println(output)
	output, _ = nn.FeedForward([]float64{1.0, 1.0, 0.0})
	fmt.Println(output)
	output, _ = nn.FeedForward([]float64{1.0, 0.0, 0.0})
	fmt.Println(output)

	fmt.Println("\x1b[32m=Testing NEAT with pole balancing test...\x1b[0m")
	configPole, err := NewConfigJSON("config_pole_balancing.json")
	if err != nil {
		fmt.Println("\x1b[31mFAIL\x1b[0m")
	}
	configPole.Summarize()
	n1 := New(configPole, PoleBalancingTest(false, 120000))
	n1.Run()

	fmt.Println("\x1b[32m=Testing NEAT with pole balancing (random)...\x1b[0m")
	configPole.Summarize()
	n2 := New(configPole, PoleBalancingTest(true, 120000))
	n2.Run()
}

func TestNEAT(t *testing.T) {
	rand.Seed(0)
	NEATUnitTest()
}
