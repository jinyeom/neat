package neat

import (
	"testing"
)

func TestNEAT(t *testing.T) {
	n, err := New(&NEATConfig{
		numSensors: 3,
		numOutputs: 2,
	})
	if err != nil {
		panic(err)
	}
	n.Run()
}
