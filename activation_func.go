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
		"identity": Identity(),
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
	DFn  func(x float64) float64 `json:"-"`    // derivative of activation
}

// Identity returns the identity function as an activation
// function. This function is only used for sensor nodes.
func Identity() *ActivationFunc {
	return &ActivationFunc{
		Name: "identity",
		Fn: func(x float64) float64 {
			return x
		},
		DFn: func(x float64) float64 {
			return 1.0
		},
	}
}

// Sigmoid returns the sigmoid function as an activation function.
func Sigmoid() *ActivationFunc {
	return &ActivationFunc{
		Name: "sigmoid",
		Fn: func(x float64) float64 {
			x = math.Max(-60.0, math.Min(60.0, x))
			return 1.0 / (1.0 + math.Exp(-x))
		},
		DFn: func(x float64) float64 {
			x = math.Max(-60.0, math.Min(60.0, x))
			sig := 1.0 / (1.0 + math.Exp(-x))
			return sig * (1.0 - sig)
		},
	}
}

// Tanh returns the hyperbolic tangent function as an activation function.
func Tanh() *ActivationFunc {
	return &ActivationFunc{
		Name: "tanh",
		Fn: func(x float64) float64 {
			x = math.Max(-60.0, math.Min(60.0, x))
			return math.Tanh(x)
		},
		DFn: func(x float64) float64 {
			x = math.Max(-60.0, math.Min(60.0, x))
			return 1.0 - math.Pow(math.Tanh(x), 2.0)
		},
	}
}

// Sin returns the sin function as an activation function.
func Sin() *ActivationFunc {
	return &ActivationFunc{
		Name: "sine",
		Fn:   math.Sin,
		DFn:  math.Cos,
	}
}

// Cos returns the cosine function as an activation function.
func Cos() *ActivationFunc {
	return &ActivationFunc{
		Name: "cosine",
		Fn:   math.Cos,
		DFn:  math.Sin,
	}
}

// ReLU returns a rectifier linear unit as an activation function.
func ReLU() *ActivationFunc {
	return &ActivationFunc{
		Name: "relu",
		Fn: func(x float64) float64 {
			return math.Max(x, 0.0)
		},
		DFn: func(x float64) float64 {
			if x > 0.0 {
				return 1.0
			}
			return 0.0
		},
	}
}

// Log returns the natural log function as an activation function.
func Log() *ActivationFunc {
	return &ActivationFunc{
		Name: "log",
		Fn:   math.Log,
		DFn: func(x float64) float64 {
			if x == 0.0 {
				return 0.0001
			}
			return 1.0 / x
		},
	}
}

// Exp returns the exponential function as an activation function.
func Exp() *ActivationFunc {
	return &ActivationFunc{
		Name: "exp",
		Fn:   math.Exp,
		DFn:  math.Exp,
	}
}

// Abs returns the absolute value function as an activation function.
func Abs() *ActivationFunc {
	return &ActivationFunc{
		Name: "abs",
		Fn:   math.Abs,
		DFn: func(x float64) float64 {
			if x >= 0.0 {
				return 1.0
			}
			return -1.0
		},
	}
}

// Square returns the square function as an activation function.
func Square() *ActivationFunc {
	return &ActivationFunc{
		Name: "square",
		Fn: func(x float64) float64 {
			return x * x
		},
		DFn: func(x float64) float64 {
			return 2.0 * x
		},
	}
}

// Cube returns the cube function as an activation function.
func Cube() *ActivationFunc {
	return &ActivationFunc{
		Name: "cube",
		Fn: func(x float64) float64 {
			return x * x * x
		},
		DFn: func(x float64) float64 {
			return 3 * x * x
		},
	}
}

// Gaussian returns the Gaussian function as an activation function, given a
// mean and a standard deviation.
func Gaussian(mu, sigma float64) *ActivationFunc {
	return &ActivationFunc{
		Name: "gaussian",
		Fn: func(x float64) float64 {
			x = math.Max(-3.4, math.Min(3.4, x))
			return (1.0 / math.Sqrt(2.0*sigma*math.Pi)) *
				math.Exp(-math.Pow(x-mu, 2.0)/(2.0*sigma*sigma))
		},
		DFn: func(x float64) float64 {
			x = math.Max(-3.4, math.Min(3.4, x))
			return (1.0 / math.Sqrt(2.0*sigma*math.Pi)) *
				math.Exp(-math.Pow(x-mu, 2.0)/(2.0*sigma*sigma)) *
				(mu - x) / (sigma * sigma)
		},
	}
}
