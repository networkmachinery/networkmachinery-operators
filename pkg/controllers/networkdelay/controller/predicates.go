package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var generationChangedPredicate = predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		return e.MetaOld.GetGeneration() != e.MetaNew.GetGeneration()
	},
}

// GenerationChangedPredicate is a predicate for generation changes.
func GenerationChangedPredicate() predicate.Predicate {
	return generationChangedPredicate
}
