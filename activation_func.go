/*


activation_func.go implementation of activation functions used in a network.

@licstart   The following is the entire license notice for
the Go code in this page.

Copyright (C) 2017 Jin Yeom

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
	"math"
)

var (
	ActivationSet = map[string]*ActivationFunc{
		"identity": Identity(),
		"sigmoid":  Sigmoid(),
	}
)

// ActivationFunc is a wrapper type for activation functions.
//
// While in this package, usage of this type is limited to only two functions,
// identity and sigmoid, as they are the only activation functions used in NEAT,
// more variety of functions are included in neat/cppn package for applications
// of CPPN-NEAT.
type ActivationFunc struct {
	Name string                  `json:"name"` // name of the function
	Fn   func(x float64) float64 `json:"-"`    // activation function
}

// Identity returns the identity function as an activation
// function. This function is only used for sensor nodes.
func Identity() *ActivationFunc {
	return &ActivationFunc{
		Name: "Identity",
		Fn: func(x float64) float64 {
			return x
		},
	}
}

// Sigmoid returns the sigmoid function as an activation function.
func Sigmoid() *ActivationFunc {
	return &ActivationFunc{
		Name: "Sigmoid",
		Fn: func(x float64) float64 {
			return 1.0 / (1.0 + math.Exp(-x))
		},
	}
}
