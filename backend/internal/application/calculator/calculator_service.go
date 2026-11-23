package calculator

import (
	"fmt"

	"sezzle/backend/internal/domain/calculation"
)

// Service coordinates the calculate use case.
type Service struct {
	strategyResolver CalculationStrategyResolver
}

// NewService creates the use case with its required output port.
func NewService(strategyResolver CalculationStrategyResolver) *Service {
	return &Service{strategyResolver: strategyResolver}
}

// Calculate validates the request, resolves the appropriate Strategy and validates its result.
func (s *Service) Calculate(operation string, operands []float64) (float64, error) {
	if err := calculation.ValidateOperandValues(operands); err != nil {
		return 0, err
	}

	strategy, ok := s.strategyResolver.Resolve(operation)
	if !ok {
		return 0, &calculation.CalculationError{
			Code:    calculation.ErrorCodeUnknownOperation,
			Message: fmt.Sprintf("The %q operation is not supported.", operation),
		}
	}

	result, err := strategy.Execute(operands)
	if err != nil {
		return 0, err
	}
	if err := calculation.ValidateFiniteResult(result); err != nil {
		return 0, err
	}
	return result, nil
}

var _ CalculateUseCase = (*Service)(nil)
