/*


compare_func.go implementation of comparison function type.

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

// compareFunc is a type of function that returns -1, 0, or 1,
// by comparing the two argument genomes' fitnesses: if the first
// argument genome has better fitness, return -1; return 0 if the
// two genomes have the same fitness; and 1, if the second genome
// has better fitness. This function type is not available
// externally; the only ways to aqcuire a comparison function are via
// DirectCompare() and InverseCompare().
type compareFunc func(g0, g1 *Genome) int

// DirectCompare returns a comparison function in which a genome's
// fitness value and its evolutionary advantage are directly related.
// In other words, the higher a fitness value, the better.
func DirectCompare() compareFunc {
	return func(g0, g1 *Genome) int {
		if g0.fitness > g1.fitness {
			return 1
		} else if g0.fitness == g1.fitness {
			return 0
		} else {
			return -1
		}
	}
}

// InverseCompare returns a comparison function in which a genome's
// fitness value and its evolutionary advantage are inversely related.
// In other words, the lower a fitness value, the better.
func InverseCompare() compareFunc {
	return func(g0, g1 *Genome) int {
		if g0.fitness < g1.fitness {
			return 1
		} else if g0.fitness == g1.fitness {
			return 0
		} else {
			return -1
		}
	}
}
