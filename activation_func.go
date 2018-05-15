// activation_func.go implementation of activation functions used in a network.
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
	"math"
)

var (
	// ActivationSet is a set of functions that can be used as activation
	// functions by neurons.
	ActivationSet = map[string]*ActivationFunc{
		"linear":   Linear(),
		"sigmoid":  Sigmoid(),
		"tanh":     Tanh(),
		"sin":      Sin(),
		"cos":      Cos(),
		"relu":     ReLU(),
		"log":      Log(),
		"exp":      Exp(),
		"abs":      Abs(),
		"square":   Square(),
		"cube":     Cube(),
		"gaussian": Gaussian(0.0, 1.0),
	}
)

// ActivationFunc is a wrapper type for activation functions.
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

// Tanh returns the hyperbolic tangent function as an activation function.
func Tanh() *ActivationFunc {
	return &ActivationFunc{
		Name: "Tanh",
		Fn:   math.Tanh,
	}
}

// Sin returns the sin function as an activation function.
func Sin() *ActivationFunc {
	return &ActivationFunc{
		Name: "Sine",
		Fn:   math.Sin,
	}
}

// Cos returns the cosine function as an activation function.
func Cos() *ActivationFunc {
	return &ActivationFunc{
		Name: "Cosine",
		Fn:   math.Cos,
	}
}

// ReLU returns a rectifier linear unit as an activation function.
func ReLU() *ActivationFunc {
	return &ActivationFunc{
		Name: "ReLU",
		Fn: func(x float64) float64 {
			return math.Max(x, 0.0)
		},
	}
}

// Log returns the log function as an activation function.
func Log() *ActivationFunc {
	return &ActivationFunc{
		Name: "Log",
		Fn:   math.Log,
	}
}

// Exp returns the exponential function as an activation function.
func Exp() *ActivationFunc {
	return &ActivationFunc{
		Name: "Exp",
		Fn:   math.Exp,
	}
}

// Abs returns the absolute value function as an activation function.
func Abs() *ActivationFunc {
	return &ActivationFunc{
		Name: "Abs",
		Fn:   math.Abs,
	}
}

// Square returns the square function as an activation function.
func Square() *ActivationFunc {
	return &ActivationFunc{
		Name: "Square",
		Fn: func(x float64) float64 {
			return x * x
		},
	}
}

// Cube returns the cube function as an activation function.
func Cube() *ActivationFunc {
	return &ActivationFunc{
		Name: "Cube",
		Fn: func(x float64) float64 {
			return x * x * x
		},
	}
}

// Gaussian returns the Gaussian function as an activation function, given a
// mean and a standard deviation.
func Gaussian(mean, stdev float64) *ActivationFunc {
	return &ActivationFunc{
		Name: "Gaussian",
		Fn: func(x float64) float64 {
			return 1.0 / (stdev * math.Sqrt(2*math.Pi)) *
				math.Exp(math.Pow((x-mean)/stdev, 2.0)/-2.0)
		},
	}
}
