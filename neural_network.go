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
	"math"
	"sort"
)

// Neuron is an implementation of a single neuron of a neural network.
type Neuron struct {
	ID          int                 // neuron ID
	Type        string              // neuron type
	Signal      float64             // signal held by this neuron
	Delta       float64             // signal held by this neuron
	SynapsesIn  map[*Neuron]float64 // synapse from input neurons
	SynapsesOut map[*Neuron]float64 // synapse from output neurons
	Activation  *ActivationFunc     // activation function

	activated  bool // true if it has been activated
	propagated bool // true if it has backpropagated
}

// NewNeuron returns a new instance of neuron, given a node gene.
func NewNeuron(nodeGene *NodeGene) *Neuron {
	return &Neuron{
		ID:          nodeGene.ID,
		Type:        nodeGene.Type,
		Signal:      0.0,
		Delta:       0.0,
		SynapsesIn:  make(map[*Neuron]float64),
		SynapsesOut: make(map[*Neuron]float64),
		Activation:  nodeGene.Activation,
		activated:   false,
		propagated:  false,
	}
}

// String returns the string representation of Neuron.
func (n *Neuron) String() string {
	if len(n.SynapsesIn) == 0 {
		return fmt.Sprintf("[%s(%d, %s)]", n.Type, n.ID, n.Activation.Name)
	}
	str := fmt.Sprintf("[%s(%d, %s)] (\n", n.Type, n.ID, n.Activation.Name)
	for neuron, weight := range n.SynapsesIn {
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
	if n.activated || len(n.SynapsesIn) == 0 {
		return n.Signal
	}
	n.activated = true

	inputSum := 0.0
	for neuron, weight := range n.SynapsesIn {
		inputSum += neuron.Activate() * weight
	}
	n.Signal = n.Activation.Fn(inputSum)
	return n.Signal
}

// Propagate retrieves deltas from neurons that are connected from this neuron
// and return its delta.
func (n *Neuron) Propagate() float64 {
	// if the neuron's already activated, or it isn't connected from any neurons,
	// return its current signal.
	if n.propagated || len(n.SynapsesOut) == 0 {
		return n.Delta
	}
	n.propagated = true

	errorSum := 0.0
	for neuron, weight := range n.SynapsesOut {
		errorSum += neuron.Propagate() * weight
	}
	n.Delta = n.Activation.DFn(n.Signal) * errorSum
	return n.Delta
}

// UpdateWeights updates each synapse weight, given an argument learning rate
// (alpha). This method must be called after the neuron's delta value is already
// updated.
func (n *Neuron) UpdateWeights(lr float64) {
	for neuron, weight := range n.SynapsesIn {
		updated := weight - lr*neuron.Signal*n.Delta
		n.SynapsesIn[neuron] = updated
		neuron.SynapsesOut[n] = updated
	}
}

// NeuralNetwork is an implementation of the phenotype neural network that is
// decoded from a genome.
type NeuralNetwork struct {
	InputNeurons  []*Neuron // input neurons
	OutputNeurons []*Neuron // output neurons
	Neurons       []*Neuron // all neurons in the network
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
					neurons[out].SynapsesIn[neurons[in]] = connGene.Weight
					neurons[in].SynapsesOut[neurons[out]] = connGene.Weight
				}
			}
		}
	}
	return &NeuralNetwork{inputNeurons, outputNeurons, neurons}
}

// String returns the string representation of NeuralNetwork.
func (n *NeuralNetwork) String() string {
	str := fmt.Sprintf("NeuralNetwork(%d, %d):\n",
		len(n.InputNeurons), len(n.OutputNeurons))
	for _, neuron := range n.Neurons {
		str += neuron.String() + "\n"
	}
	return str[:len(str)-1]
}

// FeedForward propagates inputs signals from input neurons to output neurons,
// and return output signals.
func (n *NeuralNetwork) FeedForward(inputs []float64) ([]float64, error) {
	if len(inputs) != len(n.InputNeurons) {
		errStr := "Invalid number of inputs: %d != %d"
		return nil, fmt.Errorf(errStr, len(n.InputNeurons), len(inputs))
	}

	// register sensor inputs
	for i, neuron := range n.InputNeurons {
		neuron.Signal = inputs[i]
	}

	// recursively propagate from input neurons to output neurons
	outputs := make([]float64, 0, len(n.OutputNeurons))
	for _, neuron := range n.OutputNeurons {
		outputs = append(outputs, neuron.Activate())
	}

	// reset all neurons
	for _, neuron := range n.Neurons {
		neuron.activated = false
	}

	return outputs, nil
}

// Backprop
func (n *NeuralNetwork) Backprop(inputs, targets []float64,
	learningRate float64) (float64, error) {
	if len(targets) != len(n.OutputNeurons) {
		errStr := "Invalid number of outputs %d != %d"
		return -1.0, fmt.Errorf(errStr, len(n.OutputNeurons), len(targets))
	}
	outputs, err := n.FeedForward(inputs)
	if err != nil {
		return -1.0, err
	}

	// mean square error
	mse := 0.0

	// compute delta values
	for i, neuron := range n.OutputNeurons {
		outputErr := math.Pow(outputs[i]-targets[i], 2.0)
		neuron.Delta = outputErr * neuron.Activation.DFn(outputs[i])
		mse += outputErr
	}

	// recursively propagate delta from output to input neurons
	for _, neuron := range n.InputNeurons {
		neuron.Propagate()
	}
	for _, neuron := range n.Neurons {
		neuron.UpdateWeights(learningRate)
	}

	return mse / float64(len(n.OutputNeurons)), nil
}

// CPPN is an alias type of NeuralNetwork; there is no functional difference
// between CPPN and NeuralNetwork, i.e., CPPN type is purely symbolic.
type CPPN *NeuralNetwork
