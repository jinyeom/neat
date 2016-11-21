package neat

import (
	"testing"
)

func TestNEAT(t *testing.T) {
	n, err := New(&Config{
		NumSensors: 3,
		NumOutputs: 2,
	})
	if err != nil {
		panic(err)
	}
	n.Run()
}
