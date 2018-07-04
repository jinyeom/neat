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

// Neuron is an implementation of a single neuron of a neural network.
type Neuron struct {
	ID         int                 // neuron ID
	Type       string              // neuron type
	Signal     float64             // signal held by this neuron
	Synapses   map[*Neuron]float64 // synapse from input neurons
	Activation *ActivationFunc     // activation function

	activated bool // true if it has been activated
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
		// github/@jesuicamille: if Activation is <nil>, print N/A
		switch n.Activation {
		case nil:
			return fmt.Sprintf("[%s(%d, %s)]", n.Type, n.ID, "N/A")
		default:
			return fmt.Sprintf("[%s(%d, %s)]", n.Type, n.ID, n.Activation.Name)
		}
	}

	var str string

	// github/@jesuicamille: if Activation is <nil>, print N/A
	switch n.Activation {
	case nil:
		str = fmt.Sprintf("[%s(%d, %s)] (\n", n.Type, n.ID, "N/A")
	default:
		str = fmt.Sprintf("[%s(%d, %s)] (\n", n.Type, n.ID, n.Activation.Name)
	}

	for neuron, weight := range n.Synapses {
		// github/@jesuicamille: if Activation is <nil>, print N/A
		switch neuron.Activation {
		case nil:
			str += fmt.Sprintf("  <--{%.3f}--[%s(%d, %s)]\n",
				weight, neuron.Type, neuron.ID, "N/A")
		default:
			str += fmt.Sprintf("  <--{%.3f}--[%s(%d, %s)]\n",
				weight, neuron.Type, neuron.ID, neuron.Activation.Name)
		}
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
	Neurons []*Neuron // all neurons in the network

	inputNeurons  []*Neuron // input neurons
	outputNeurons []*Neuron // output neurons
}

// NewNeuralNetwork returns a new instance of NeuralNetwork given a genome to
// decode from.
func NewNeuralNetwork(g *Genome) *NeuralNetwork {
	sort.Slice(g.NodeGenes, func(i, j int) bool {
		return g.NodeGenes[i].ID < g.NodeGenes[j].ID
	})

	inputNeurons := make([]*Neuron, 0, len(g.NodeGenes))
	outputNeurons := make([]*Neuron, 0, len(g.NodeGenes))
	neurons := make([]*Neuron, 0, len(g.NodeGenes))

	for _, nodeGene := range g.NodeGenes {
		neuron := NewNeuron(nodeGene)

		// record input and output neurons separately
		if nodeGene.Type == "input" {
			inputNeurons = append(inputNeurons, neuron)
		} else if nodeGene.Type == "output" {
			outputNeurons = append(outputNeurons, neuron)
		}

		neurons = append(neurons, neuron)
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
	return &NeuralNetwork{neurons, inputNeurons, outputNeurons}
}

// String returns the string representation of NeuralNetwork.
func (n *NeuralNetwork) String() string {
	str := fmt.Sprintf("NeuralNetwork(%d, %d):\n",
		len(n.inputNeurons), len(n.outputNeurons))
	for _, neuron := range n.Neurons {
		str += neuron.String() + "\n"
	}
	return str[:len(str)-1]
}

// FeedForward propagates inputs signals from input neurons to output neurons,
// and return output signals.
func (n *NeuralNetwork) FeedForward(inputs []float64) ([]float64, error) {
	if len(inputs) != len(n.inputNeurons) {
		errStr := "Invalid number of inputs: %d != %d"
		return nil, fmt.Errorf(errStr, len(n.inputNeurons), len(inputs))
	}

	// register sensor inputs
	for i, neuron := range n.inputNeurons {
		neuron.Signal = inputs[i]
	}

	// recursively propagate from input neurons to output neurons
	outputs := make([]float64, 0, len(n.outputNeurons))
	for _, neuron := range n.outputNeurons {
		outputs = append(outputs, neuron.Activate())
	}

	// reset all neurons
	for _, neuron := range n.Neurons {
		neuron.Signal = 0.0
		neuron.activated = false
	}

	return outputs, nil
}
