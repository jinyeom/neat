![alt text](https://github.com/jinyeom/neat/blob/master/banner.png "neat")
[![GoDoc](https://godoc.org/github.com/jinyeom/neat?status.svg)](https://godoc.org/github.com/jinyeom/neat)

NEAT (NeuroEvolution of Augmenting Topologies) is a neuroevolution algorithm by 
Dr. Kenneth O. Stanley which evolves not only neural networks' weights but also their 
topologies. 

## Installation
To install `neat` run the following:

```bash
$ go get -u github.com/jinyeom/neat
```

## Usage

This NEAT package is as simple as plug and play. All you have to do is to create
a new instance of NEAT, given the configuration from a JSON file, for which the
template is provided below, and an evaluation method of a neural network, and 
run.

```json
{
	"numInputs": 3,
	"numOutputs": 1,
	"numGenerations": 50,
	"populationSize": 100,
	"initFitness": 9999.0,
	"minimizeFitness": true,
	"survivalRate": 0.5,
	"hallOfFameSize": 10,
	"ratePerturb": 0.2,
	"rateAddNode": 0.2,
	"rateAddConn": 0.2,
	"distanceThreshold": 20.0,
	"coeffUnmatching": 1.0,
	"coeffMatching": 1.0
}
```

Now that you have the configuration JSON file is ready as `config.json`, we can
start experiment with NEAT. Below is an example XOR experiment.

```go
package main

import (
	"log"
	"math"

	// Import NEAT package after installing the package through
	// the instruction provided above.
	"github.com/jinyeom/neat"
)

func main() {

	// First, create a new instance of Config from the JSON file created above.
	// If there's a file import error, the program will crash.
	config, err := neat.NewConfigJSON("config.json")
	if err != nil{
		log.Fatal(err)
	}

	// Then, we can define the evaluation function, which is a type of function
	// which takes a neural network, evaluates its performance, and returns some
	// score that indicates its performance. This score is essentially a genome's
	// fitness score. With the configuration and the evaluation function we
	// defined, we can create a new instance of NEAT and start the evolution 
	// process.
	//
	// By passing true as an argument of Run(), we can also observe the progress
	// by generations.
	neat.New(config, func(n *neat.NeuralNetwork) float64 {
		score := 0.0

		inputs := make([]float64, 3)
		inputs[0] = 1.0 // bias

		// 0 xor 0
		inputs[1] = 0.0
		inputs[2] = 0.0
		output, err := n.FeedForward(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 0.0), 2.0)

		// 0 xor 1
		inputs[1] = 0.0
		inputs[2] = 1.0
		output, err = n.FeedForward(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 1.0), 2.0)

		// 1 xor 0
		inputs[1] = 1.0
		inputs[2] = 0.0
		output, err = n.FeedForward(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 1.0), 2.0)

		// 1 xor 1
		inputs[1] = 1.0
		inputs[2] = 1.0
		output, err = n.FeedForward(inputs)
		if err != nil {
			log.Fatal(err)
		}
		score += math.Pow((output[0] - 0.0), 2.0)

		return score
	}).Run(true)
}

```

## License
This package is under GNU General Public License.
