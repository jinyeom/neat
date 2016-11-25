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
}
