package neat

// Species is an implementation of species, or niche for speciation of genomes
// that are differentiated by their toplogical differences, measured with
// compatibility distance. Each species is created with a new genome that is not
// compatible with other genomes in the population, i.e., when a genome is not
// compatible with any other species.
type Species struct {
	ID             int       // species ID
	Stagnation     int       // number of generations of stagnation
	Representative *Genome   // genome that represents this species
	Members        []*Genome // member genomes
}

// NewSpecies creates and returns a new instance of Species, given an initial
// genome that will also become the new species' representative.
func NewSpecies(id int, g *Genome) *Species {
	return &Species{
		ID:             id,
		Stagnation:     0,
		Representative: g,
		Members:        []*Genome{g},
	}
}

// Register adds an argument genome as a new member of this species; in addition,
// if the new member genome outperforms this species' representative, it becomes
// the new representative.
func (s *Species) Register(g *Genome, minimizeFitness bool) {
	s.Members = append(s.Members, g)
	if minimizeFitness {
		if g.Fitness < s.Representative.Fitness {
			s.Representative = g
		}
	} else {
		if g.Fitness > s.Representative.Fitness {
			s.Representative = g
		}
	}
}
