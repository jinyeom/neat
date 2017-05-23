// Copyright (C) 2017  Jin Yeom
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

/*
Package neat provides an implementation of NeuroEvolution of Augmenting
Topologies (NEAT) as a plug-and-play framework, which can be used by simply
adding and appropriate configuration and an evaluation function. This package
can also be extended to applications of CPPN-NEAT by using neat/cppn package.

NEAT

NEAT (NeuroEvolution of Augmenting Topologies) is a neuroevolution algorithm
by Dr. Kenneth O. Stanley which evolves not only neural networks' weights but
also their topologies. This method starts the evolution process with genomes
with minimal structure, then complexifies the structure of each genome as it
progresses. You can read the original paper from here:
http://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf

Example

This NEAT package is as simple as plug and play. All you have to do is to create
a new instance of NEAT, given the configuration from a JSON file, for which the
template is provided below, and an evaluation method of a neural network, and
run.

  {
  	"experimentName": "XOR Test",
  	"verbose": true,
  	"numInputs": 3,
  	"numOutputs": 1,
  	"numGenerations": 50,
  	"populationSize": 100,
  	"initFitness": 9999.0,
  	"minimizeFitness": true,
  	"survivalRate": 0.5,
  	"stagnationLimit": 5,
  	"ratePerturb": 0.2,
  	"rateAddNode": 0.2,
  	"rateAddConn": 0.2,
  	"rateMutateChild": 0.5,
  	"distanceThreshold": 20.0,
  	"coeffUnmatching": 1.0,
  	"coeffMatching": 1.0
  }

Now that you have the configuration JSON file is ready as `config.json`, we can
start experiment with NEAT. Below is an example XOR experiment.

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
  	// score that indicates its performance. This score is essentially a
    // genome's fitness score. With the configuration and the evaluation
    // function we defined, we can create a new instance of NEAT and start the
    // evolution process.
  	neat.New(config, neat.XORTest()).Run()
  }
*/
package neat
