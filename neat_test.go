package neat

import (
	"testing"
)

func TestNEAT(t *testing.T) {
	n, err := New(&Param{
		NumSensors:     3,
		NumOutputs:     2,
		PopulationSize: 50,
		CrossoverRate:  0.1,
		MutAddNodeRate: 0.1,
		MutAddConnRate: 0.1,
		MutWeightRate:  0.1,
	}, nil)
	if err != nil {
		panic(err)
	}
	n.Run()
}
