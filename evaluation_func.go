package neat

import (
	"log"
	"math"
)

// EvaluationFunc is a type of function that evaluates an argument neural
// network and returns a its fitness (performance) score.
type EvaluationFunc func(*NeuralNetwork) float64

// XORTest returns an XOR test as an evaluation function.
func XORTest() EvaluationFunc {
	return func(n *NeuralNetwork) float64 {
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
	}
}
