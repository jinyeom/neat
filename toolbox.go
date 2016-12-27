package neat

import (
	"errors"
)

// Toolbox is a container that consists of all functions that are
// utilized for NEAT operations, such as activation, selection, or
// evaluation.
type Toolbox struct {
	Activation ActivationSet
	Comparison CompareFunc
	Selection  SelectFunc
	Evaluation EvaluationFunc
}

// IsValid checks whether this tool box is valid and returns an error
// if there's a function that is not initialized.
func (t *Toolbox) IsValid() error {
	if t.Activation == nil {
		return errors.New("activation set not initialized")
	}
	if t.Comparison == nil {
		return errors.New("comparison not initialized")
	}
	if t.Selection == nil {
		return errors.New("selection not initialized")
	}
	if t.Evaluation == nil {
		return errors.New("evaluation not initialized")
	}
	return nil
}
