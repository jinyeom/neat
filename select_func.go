/*


select_func.go implementation of selection function type.

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

// SelectFunc is a type of function that selects a genome
// from the argument pool of genomes.
type SelectFunc func([]*Genome) *Genome

// TSelect() returns a selection function that performs
// Tournament Selection given a comparison function.
func TSelect() SelectFunc {
	return func(c CompareFunc, p []*Genome) *Genome {
		popSize := len(p)
		best := rand.Intn(popSize)
		for i := 0; i < popSize; i++ {
			next := rand.Intn(popSize)
			if c(p[next], p[best]) == 1 {
				best = next
			}
		}
		return p[best]
	}
}

// FPSelect() returns a selection function that performs
// Fitness-Proportionate Selection (not recommended).
func FPSelect() SelectFunc {
	return func(c CompareFunc, p []*Genome) *Genome {
		popSize := len(p)
		best := p[rand.Intn(popSize)]
		bestScore := best.fitness
		for i := 0; i < popSize; i++ {
			if c(p[i], best) == 1 {
				best = p[i]
				bestScore = p[i].fitness
			}
		}
		// stochastic acceptance
		for {
			i := rand.Intn(popSize)
			r := p[i].fitness / bestScore
			if rand.Float64() < r {
				return p[i]
			}
		}
	}
}
