package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("NEAT Copyright (c) 2016 by White Wolf Studio\n")

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Number of sensors: ")
	s := reader.ReadString('\n')

}
