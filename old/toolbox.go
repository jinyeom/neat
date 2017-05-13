/*


toolbox.go implementation of tool box in NEAT.

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
	"errors"
)

// Toolbox is a container that consists of all functions that are
// utilized for NEAT operations, such as activation, selection, or
// evaluation. Toolbox needs be to initialized prior to creating
// a NEAT struct.
type Toolbox struct {
	Activation ActivationSet
	Comparison CompareFunc
	Evaluation EvaluationFunc
}

// IsValid checks whether this tool box is valid and returns an error
// if there's a function that is not initialized.
func (t *Toolbox) IsValid() error {
	if t.Activation == nil {
		return errors.New("activation set not initialized")
	}
	if t.Comparison == nil {
		return errors.New("comparison not initialized")
	}
	if t.Evaluation == nil {
		return errors.New("evaluation not initialized")
	}
	return nil
}
