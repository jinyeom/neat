package neat

// Species is an implementation of species, or niche for speciation of genomes
// that are differentiated by their toplogical differences, measured with
// compatibility distance. Each species is created with a new genome that is not
// compatible with other genomes in the population, i.e., when a genome is not
// compatible with any other species.
type Species struct {
	ID       int       // species ID
	Stagnant int       // number of generations of stagnation
	Members  []*Genome // member genomes
}
