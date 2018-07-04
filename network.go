package neat

type Network struct {
	signals []float64
}

// Note from github/@jesuiscamille: initially, there was no return value required,
// but there was one returned in the code. I added the required return value
// second note: *new([]float64) was originally make([]float64)
func NewNetwork(g *Genome) *Network {
	var signals []float64
	return &Network{
		signals: signals,
	}
}
