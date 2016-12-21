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
)

// ActivationSet is a collection of all activation functions.
var ActivationSet = []*ActivationFunc{
	Identity(),
	Sigmoid(),
}

// ActivationFunc is a function type of which the independence
// variable is a single float64 and its dependence variable is
// also a single float64.
type ActivationFunc struct {
	name string                  // name of the function
	fn   func(x float64) float64 // activation function
}

// Identity returns the identity function as an activation
// function.
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

// Name returns the activation function's name.
func (a *ActivationFunc) Name() string {
	return a.name
}

// Fn applies the activation function on a given value; returns its output.
func (a *ActivationFunc) Fn(x float64) float64 {
	return a.fn(x)
}
