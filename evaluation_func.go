/*


evaluation_func.go implementation of the evaluation function.

@licstart   The following is the entire license notice for
the Go code in this page.

Copyright (C) 2016 jin yeom, whitewolf.studio

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

As additional permission under GNU GPL version 3 section 7, you
may distribute non-source (e.g., minimized or compacted) forms of
that code without the copy of the GNU GPL normally required by
section 4, provided you include this license notice and a URL
through which recipients can access the Corresponding Source.

@licend    The above is the entire license notice
for the Go code in this page.


*/

package neat

import (
	"log"
	"math"
)

// EvaluationFunc is a type of function that evaluates a network
// and returns a float64 number as a fitness score.
type EvaluationFunc func(n *Network) float64

// XORTest returns an evaluation function in which a neural network is
// evaluated for its ability to compute XOR operations; the evaluation
// function returns its overall error as the network's fitness, which
// means a network's score and its fitness are inversely related.
func XORTest() EvaluationFunc {
	return func(n *Network) float64 {
		score := 0.0

		inputs := make([]float64, 3)
		inputs[0] = 1.0 // bias

		// 0 xor 0
		inputs[1] = 0.0
		inputs[2] = 0.0
		output, err := n.ForwardPropagate(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 0.0), 2.0)

		// 0 xor 1
		inputs[1] = 0.0
		inputs[2] = 1.0
		output, err = n.ForwardPropagate(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 1.0), 2.0)

		// 1 xor 0
		inputs[1] = 1.0
		inputs[2] = 0.0
		output, err = n.ForwardPropagate(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 1.0), 2.0)

		// 1 xor 1
		inputs[1] = 1.0
		inputs[2] = 1.0
		output, err = n.ForwardPropagate(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 0.0), 2.0)

		return score
	}
}
