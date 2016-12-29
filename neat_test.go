package neat

import (
	"fmt"
	"log"
	"testing"
)

func Status(n *NEAT) {
	fmt.Printf("")
}

func TestNEAT(t *testing.T) {
	p := NewParam("example_param.np")
	t := &Toolbox{
		NEATSet(),
		DirectCompare(),
		XORTest(),
	}

	Init(p, t)

	fmt.Printf("=== Creating NEAT ===\n")
	n, err := New()
	if err != nil {
		log.Fatal(err)
	}

}
