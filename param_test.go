package neat

import (
	"fmt"
	"testing"
)

func TestParam(t *testing.T) {
	fmt.Printf("=== Test Param ===\n")
	p, err := NewParam("example_param.np")
	if err != nil {
		panic(err)
	}
	if err = p.IsValid(); err != nil {
		panic(err)
	}
  fmt.Printf("Number of sensors: %d\n", p.NumSensors)
  fmt.Printf("Number of outputs: %d\n", p.NumOutputs)
}
