package neat

// EvaluationFunc is a type of function that evaluates an argument neural
// network and returns a its fitness (performance) score.
type EvaluationFunc func(*NeuralNetwork) float64
