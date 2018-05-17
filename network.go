package neat

type Network struct {
	signals []float64
}

func NewNetwork(g *Genome) {
	return &Network{
		signals: make([]float64),
	}
}
