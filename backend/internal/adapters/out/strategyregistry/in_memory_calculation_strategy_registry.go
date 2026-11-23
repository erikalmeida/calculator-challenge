// Package strategyregistry provides the in-memory implementation of the resolver port.
package strategyregistry

import "sezzle/backend/internal/domain/calculation"

// InMemoryCalculationStrategyRegistry binds public operation names to their Strategy implementations.
type InMemoryCalculationStrategyRegistry struct {
	strategies map[string]calculation.CalculationStrategy
}

// NewInMemoryCalculationStrategyRegistry creates the registry required by the current take-home scope.
func NewInMemoryCalculationStrategyRegistry() *InMemoryCalculationStrategyRegistry {
	return &InMemoryCalculationStrategyRegistry{strategies: map[string]calculation.CalculationStrategy{
		"add":      calculation.AdditionStrategy{},
		"subtract": calculation.SubtractionStrategy{},
		"multiply": calculation.MultiplicationStrategy{},
		"divide":   calculation.DivisionStrategy{},
	}}
}

// Resolve satisfies calculator.CalculationStrategyResolver.
func (r *InMemoryCalculationStrategyRegistry) Resolve(name string) (calculation.CalculationStrategy, bool) {
	strategy, ok := r.strategies[name]
	return strategy, ok
}
