package neat

// Species is an implementation of species, or niche for speciation of genomes
// that are differentiated by their toplogical differences, measured with
// compatibility distance. Each species is created with a new genome that is not
// compatible with other genomes in the population, i.e., when a genome is not
// compatible with any other species.
type Species struct {
	ID             int       // species ID
	Stagnation     int       // number of generations of stagnation
	Representative *Genome   // genome that represents this species (permanent)
	Best           *Genome   // best performing genome
	BestFitness    float64   // best fitness score in this species
	Members        []*Genome // member genomes
}

// NewSpecies creates and returns a new instance of Species, given an initial
// genome that will also become the new species' representative.
func NewSpecies(id int, g *Genome) *Species {
	return &Species{
		ID:             id,
		Stagnation:     0,
		Representative: g,
		Best:           g,
		BestFitness:    g.Fitness,
		Members:        []*Genome{g},
	}
}

// Register adds an argument genome as a new member of this species; in
// addition, if the new member genome outperforms this species' best genome, it
// replaces the best genome in this species.
func (s *Species) Register(g *Genome, minimizeFitness bool) {
	s.Members = append(s.Members, g)
	if minimizeFitness {
		if g.Fitness < s.Best.Fitness {
			s.Best = g
			s.BestFitness = g.Fitness
			s.Stagnation = 0
		}
	} else {
		if g.Fitness > s.Best.Fitness {
			s.Best = g
			s.BestFitness = g.Fitness
			s.Stagnation = 0
		}
	}
}

// Flush empties the species membership, except for its representative.
func (s *Species) Flush() {
	s.Members = []*Genome{}
	s.Best = s.Representative
	s.BestFitness = s.Representative.Fitness
}
