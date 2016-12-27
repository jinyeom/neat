# NEAT (NeuroEvolution of Augmenting Topologies)

[![GoDoc](https://godoc.org/github.com/whitewolf-studio/neat?status.svg)](https://godoc.org/github.com/whitewolf-studio/neat)

NEAT (NeuroEvolution of Augmenting Topologies) is a neuroevolution algorithm that
evolves not only neural networks' weights but also their topologies. This package
is created for optimization of neural networks for general purpose reinforcement
learning, given that the user can provide a clear evaluation function.

## Example

```go
package main

import (
    "github.com/whitewolf-studio/neat"
)

param := NewParam("xor_param.np")
toolbox := &Toolbox{
    Activation: neat.NEATSet(),
    Comparison: neat.DirectCompare(),
    Selection:  TSelect(DirectCompare()),
    Evaluation: XORTest(),
}

neat.Init(param, toolbox)

n := neat.New()
n.Run()
// for parallel computing,
// numProcs := 4
// n.RunParallel(numProcs)


