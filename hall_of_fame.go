package neat

// HallOfFame is a list of genomes that has the best performance throughout the
// process of evolution. The list is in increasing order of fitness.
type HallOfFame struct {
	Size        int
	BestGenomes []*Genome
}

// NewHallOfFame returns a new instance of HallOfFame. The list of genomes is
// initialized by placeholding genomes, each of which has -1 as its ID, and
// initial fitness score as its fitness.
func NewHallOfFame(size int, initFitness float64) *HallOfFame {
	return &HallOfFame{
		Size: size,
		BestGenomes: func() []*Genome {
			placeholders := make([]*Genome, size)
			for i := range placeholders {
				placeholders[i] = &Genome{
					ID:        -1,
					SpeciesID: -1,
					NodeGenes: nil,
					ConnGenes: nil,
					Fitness:   initFitness,
				}
			}
			return placeholders
		}(),
	}
}

// Update checks for the argument genome's place in the hall of fame. If it
// performs better than one of the genomes in the list, update the list.
func (h *HallOfFame) Update(g *Genome, comparison ComparisonFunc) {
	if comparison(g, h.BestGenomes[0]) {
		i := 1
		for i < h.Size {
			if !comparison(g, h.BestGenomes[i]) {
				h.BestGenomes[i-1] = g
				return
			}
			i++
		}
		h.BestGenomes[h.Size-1] = g
	}
}

// BestScore returns the current best score in the list.
func (h *HallOfFame) BestScore() float64 {
	return h.BestGenomes[h.Size-1].Fitness
}
