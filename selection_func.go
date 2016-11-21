package neat

// SelectionFunc is a type of function that selects a genome
// from the argument pool of genomes.
type SelectionFunc func([]*Genome) *Genome
