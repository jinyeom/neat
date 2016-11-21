package neat

// EvaluationFunc is a type of function that evaluates a genome
// and returns a float64 number as a fitness score.
type EvaluationFunc func(g *Genome) float64
