package neat

import (
	"fmt"
	"math/rand"
	"testing"
)

func NeuralNetworkUnitTest() {
	fmt.Println("===== Neural Network Unit Test =====")

	g0 := NewGenome(0, 3, 1, 0.0)

	for i := 0; i < 3; i++ {
		g0.MutatePerturb(1.0)
		g0.MutateAddNode(1.0, ActivationSet["sigmoid"])
		g0.MutateAddConn(1.0)
	}

	n0 := NewNeuralNetwork(g0)
	fmt.Println(n0.String())

	fmt.Println("\x1b[32m=Testing feedforward...\x1b[0m")

	inputs := []float64{rand.NormFloat64(), rand.NormFloat64(), 1.0}
	fmt.Println("inputs:", inputs)

	outputs, err := n0.FeedForward(inputs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("outputs:", outputs)
}

func TestNeuralNetwork(t *testing.T) {
	rand.Seed(0)
	NeuralNetworkUnitTest()
}
