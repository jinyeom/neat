// neural_network.go implementation of the neural network.
//
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

package neat

import (
	"fmt"
	"sort"
)

// Network defines a phenotype network that feed forwards inputs and returns
// outputs.
type Network interface {
	FeedForward([]float64) ([]float64, error)
}

// Neuron is an implementation of a single neuron of a neural network.
type Neuron struct {
	ID         int                 // neuron ID
	Type       string              // neuron type
	Signal     float64             // signal held by this neuron
	Synapses   map[*Neuron]float64 // synapse from input neurons
	Activation *ActivationFunc     // activation function
	activated  bool                // true if it has been activated
}

// NewNeuron returns a new instance of neuron, given a node gene.
func NewNeuron(nodeGene *NodeGene) *Neuron {
	return &Neuron{
		ID:         nodeGene.ID,
		Type:       nodeGene.Type,
		Signal:     0.0,
		Synapses:   make(map[*Neuron]float64),
		Activation: nodeGene.Activation,
		activated:  false,
	}
}

// String returns the string representation of Neuron.
func (n *Neuron) String() string {
	if len(n.Synapses) == 0 {
		return fmt.Sprintf("[%s(%d, %s)]", n.Type, n.ID, n.Activation.Name)
	}
	str := fmt.Sprintf("[%s(%d, %s)] (\n", n.Type, n.ID, n.Activation.Name)
	for neuron, weight := range n.Synapses {
		str += fmt.Sprintf("  <--{%.3f}--[%s(%d, %s)]\n",
			weight, neuron.Type, neuron.ID, neuron.Activation.Name)
	}
	return str + ")"
}

// Activate retrieves signal from neurons that are connected to this neuron and
// return its signal.
func (n *Neuron) Activate() float64 {
	// if the neuron's already activated, or it isn't connected from any neurons,
	// return its current signal.
	if n.activated || len(n.Synapses) == 0 {
		return n.Signal
	}
	n.activated = true

	inputSum := 0.0
	for neuron, weight := range n.Synapses {
		inputSum += neuron.Activate() * weight
	}
	n.Signal = n.Activation.Fn(inputSum)
	return n.Signal
}

// NeuralNetwork is an implementation of the phenotype neural network that is
// decoded from a genome.
type NeuralNetwork struct {
	NumInputs  int       // number of inputs
	NumOutputs int       // number of outputs
	Neurons    []*Neuron // neurons in the neural network
}

// NewNeuralNetwork returns a new instance of NeuralNetwork given a genome to
// decode from.
func NewNeuralNetwork(g *Genome) *NeuralNetwork {
	sort.Slice(g.NodeGenes, func(i, j int) bool {
		return g.NodeGenes[i].ID < g.NodeGenes[j].ID
	})

	numInputs := 0
	numOutputs := 0

	neurons := make([]*Neuron, 0, len(g.NodeGenes))
	for _, nodeGene := range g.NodeGenes {
		if nodeGene.Type == "input" {
			numInputs++
		} else if nodeGene.Type == "output" {
			numOutputs++
		}
		neurons = append(neurons, NewNeuron(nodeGene))
	}

	for _, connGene := range g.ConnGenes {
		if !connGene.Disabled {
			if in := sort.Search(len(neurons), func(i int) bool {
				return neurons[i].ID >= connGene.From
			}); in < len(neurons) && neurons[in].ID == connGene.From {
				if out := sort.Search(len(neurons), func(i int) bool {
					return neurons[i].ID >= connGene.To
				}); out < len(neurons) && neurons[out].ID == connGene.To {
					neurons[out].Synapses[neurons[in]] = connGene.Weight
				}
			}
		}
	}
	return &NeuralNetwork{numInputs, numOutputs, neurons}
}

// String returns the string representation of NeuralNetwork.
func (n *NeuralNetwork) String() string {
	str := fmt.Sprintf("NeuralNetwork(%d, %d):\n", n.NumInputs, n.NumOutputs)
	for _, neuron := range n.Neurons {
		str += neuron.String() + "\n"
	}
	return str[:len(str)-1]
}

// FeedForward propagates inputs signals from input neurons to output neurons,
// and return output signals.
func (n *NeuralNetwork) FeedForward(inputs []float64) ([]float64, error) {
	if len(inputs) != n.NumInputs {
		errStr := "Invalid number of inputs: %d != %d"
		return nil, fmt.Errorf(errStr, n.NumInputs, len(inputs))
	}

	// register sensor inputs
	for i := 0; i < n.NumInputs; i++ {
		n.Neurons[i].Signal = inputs[i]
	}

	// recursively propagate from input neurons to output neurons
	outputs := make([]float64, 0, n.NumOutputs)
	for i := n.NumInputs; i < n.NumInputs+n.NumOutputs; i++ {
		outputs = append(outputs, n.Neurons[i].Activate())
	}

	// reset all neurons
	for _, neuron := range n.Neurons {
		neuron.activated = false
	}

	return outputs, nil
}

// CPPN is an alias type of NeuralNetwork; there is no functional difference
// between CPPN and NeuralNetwork, i.e., CPPN type is purely symbolic.
type CPPN *NeuralNetwork
