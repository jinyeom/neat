// evaluation_func.go implementation of evaluation functions of a network.
//
// Copyright (C) 2017  Jin Yeom
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package neat

import (
	"log"
	"math"
	"math/rand"
)

// EvaluationFunc is a type of function that evaluates an argument neural
// network and returns a its fitness (performance) score.
type EvaluationFunc func(*NeuralNetwork) float64

// XORTest returns an XOR test as an evaluation function. The fitness is
// measured with the total error, which should be minimized.
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

// PoleBalancingTest returns the pole balancing task as an evaluation function.
// The fitness is measured with how long the network can balanced the pole,
// given a max time. Suggested max time is 120000 ticks.
func PoleBalancingTest(randomStart bool, maxTime int) EvaluationFunc {
	// physics constants
	xLim := 2.4      // x position limit [-2.4, 2.4]
	dxLim := 1.0     // x velocity limit [-1.0, 1.0]
	thLim := 0.2     // theta limit [-0.2, 0.2]
	dthLim := 1.5    // angular velocity limit [-1.5, 1.5]
	gravity := 9.8   // gravity constant
	cartMass := 1.0  // mass of the cart
	poleMass := 0.1  // mass of the pole
	length := 0.5    // half length of pole
	forceMag := 10.0 // force applied to the cart
	tau := 0.02      // seconds between state updates

	totalMass := cartMass + poleMass
	poleMassLength := poleMass * length

	cartpole := func(action bool, inputs []float64) []float64 {
		force := forceMag
		if action {
			force = -forceMag
		}

		cosTh := math.Cos(inputs[2])
		sinTh := math.Sin(inputs[3])
		tmp := (force + poleMassLength*inputs[3]*inputs[3]*sinTh) / totalMass

		// angular acceleration
		ath := (gravity*sinTh - cosTh*tmp) /
			(length * (4.0/3.0 - poleMass*cosTh*cosTh/totalMass))

		// x acceleration
		ax := tmp - poleMassLength*ath*cosTh/totalMass

		return []float64{
			inputs[0] + tau*inputs[1],
			inputs[1] + tau*ax,
			inputs[2] + tau*inputs[3],
			inputs[3] + tau*ath,
		}
	}

	return func(n *NeuralNetwork) float64 {
		inputs := make([]float64, 4)
		if randomStart {
			inputs[0] = float64(rand.Int31()%4800)/1000.0 - xLim
			inputs[1] = float64(rand.Int31()%2000)/1000.0 - dxLim
			inputs[2] = float64(rand.Int31()%400)/1000.0 - thLim
			inputs[3] = float64(rand.Int31()%3000)/1000.0 - dthLim
		}

		for i := 0; i < maxTime; i++ {
			outputs, err := n.FeedForward(inputs)
			if err != nil {
				panic(err)
			}

			// update the next inputs; if the cart moves out of bound (xLim), or the
			// pole falls beyond the limit (thLim), return the time.
			inputs = cartpole(outputs[0] <= outputs[1], inputs)
			if math.Abs(inputs[0]) > xLim || math.Abs(inputs[2]) > thLim {
				return float64(i)
			}
		}
		return float64(maxTime)
	}
}
