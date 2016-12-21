/*


activation_func.go implementation of activation functions used in a network.

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
	"math"
	"math/rand"
)

// ActivationSet is a collection of activation functions that are
// available within a network that is being evolved.
type ActivationSet []*ActivationFunc

// NEATSet returns a set of activation functions that are used in
// general NEAT algorithm, which only includes Sigmoid function.
func NEATSet() ActivationSet {
	return []*ActivationFunc{
		Sigmoid(),
	}
}

// CPPNSet returns a set of activation functions that are used in
// CPPN-NEAT algorithm, which includes all provided types of activation
// functions, except for Gaussian function, which can be further added
// by the user.
func CPPNSet() ActivationSet {
	return []*ActivationFunc{
		Tanh(),
		Sin(),
		Cos(),
		Sigmoid(),
		ReLU(),
		Log(),
		Exp(),
		Abs(),
		Square(),
		Cube(),
	}
}

// AddFunc adds a custom activation function to this set of functions.
func (a ActivationSet) AddFunc(afn *ActivationFunc) {
	a = append(a, afn)
}

// RandFunc returns a random activation function from this set.
func (a ActivationSet) RandFunc() *ActivationFunc {
	return a[rand.Intn(len(a))]
}

// ActivationFunc is a function type of which the independence
// variable is a single float64 and its dependence variable is
// also a single float64.
type ActivationFunc struct {
	name string                  // name of the function
	fn   func(x float64) float64 // activation function
}

// Identity returns the identity function as an activation
// function. This function is only used for sensor nodes.
func Identity() *ActivationFunc {
	return &ActivationFunc{
		name: "identity",
		fn: func(x float64) float64 {
			return x
		},
	}
}

// Sigmoid returns the sigmoid (or soft step) function as an
// activation function.
func Sigmoid() *ActivationFunc {
	return &ActivationFunc{
		name: "sigmoid",
		fn: func(x float64) float64 {
			return 1.0 / (1.0 + math.Exp(-x))
		},
	}
}

// Tanh returns the hyperbolic tangent function as an activation
// function.
func Tanh() *ActivationFunc {
	return &ActivationFunc{
		name: "tanh",
		fn:   math.Tanh,
	}
}

// Sin returns the sin function as an activation function.
func Sin() *ActivationFunc {
	return &ActivationFunc{
		name: "sin",
		fn:   math.Sin,
	}
}

// Cos returns the cosine function as an activation function.
func Cos() *ActivationFunc {
	return &ActivationFunc{
		name: "cos",
		fn:   math.Cos,
	}
}

// ReLU returns a rectifier linear unit as an activation function.
func ReLU() *ActivationFunc {
	return &ActivationFunc{
		name: "relu",
		fn: func(x float64) float64 {
			if x > 0.0 {
				return x
			}
			return 0.0
		},
	}
}

// Log returns the log function as an activation function.
func Log() *ActivationFunc {
	return &ActivationFunc{
		name: "log",
		fn:   math.Log,
	}
}

// Exp returns the exponential function as an activation function.
func Exp() *ActivationFunc {
	return &ActivationFunc{
		name: "exp",
		fn:   math.Exp,
	}
}

// Abs returns the absolute value function as an activation function.
func Abs() *ActivationFunc {
	return &ActivationFunc{
		name: "abs",
		fn:   math.Abs,
	}
}

// Square returns the square function as an activation function.
func Square() *ActivationFunc {
	return &ActivationFunc{
		name: "square",
		fn: func(x float64) float64 {
			return x * x
		},
	}
}

// Cube returns the cube function as an activation function.
func Cube() *ActivationFunc {
	return &ActivationFunc{
		name: "cube",
		fn: func(x float64) float64 {
			return x * x * x
		},
	}
}

// Gaussian returns the Gaussian distribution function as an activation
// function; this function is initially not included in CPPNSet, but
// the users can add this function with custom mu and sigma value.
func Gaussian(mu, sigma float64) *ActivationFunc {
	return &ActivationFunc{
		name: "gaussian",
		fn: func(x float64) float64 {
			return 1.0 / (sigma * math.Sqrt(2*math.Pi)) *
				math.Exp(math.Pow((x-mu)/sigma, 2.0)/-2.0)
		},
	}
}

// Name returns the activation function's name.
func (a *ActivationFunc) Name() string {
	return a.name
}

// Fn applies the activation function on a given value; returns its output.
func (a *ActivationFunc) Fn(x float64) float64 {
	return a.fn(x)
}
