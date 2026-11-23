// Package calculator contains the calculation use case and its ports.
package calculator

import "sezzle/backend/internal/domain/calculation"

// CalculationStrategyResolver is the output port the use case needs to find a Strategy.
// The application does not know whether this comes from memory, configuration or another source.
type CalculationStrategyResolver interface {
	Resolve(name string) (calculation.CalculationStrategy, bool)
}
