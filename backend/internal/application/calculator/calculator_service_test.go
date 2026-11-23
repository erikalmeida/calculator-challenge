package calculator

import (
	"errors"
	"math"
	"testing"

	"sezzle/backend/internal/domain/calculation"
)

type calculationStrategyResolverStub map[string]calculation.CalculationStrategy

type nonFiniteResultStrategy struct{}

func (nonFiniteResultStrategy) Execute(_ []float64) (float64, error) {
	return math.Inf(1), nil
}

func (r calculationStrategyResolverStub) Resolve(name string) (calculation.CalculationStrategy, bool) {
	strategy, ok := r[name]
	return strategy, ok
}

func TestCalculate(t *testing.T) {
	service := NewService(calculationStrategyResolverStub{"add": calculation.AdditionStrategy{}, "multiply": calculation.MultiplicationStrategy{}})

	result, err := service.Calculate("add", []float64{2, 3})
	if err != nil || result != 5 {
		t.Fatalf("Calculate(add) = (%v, %v), want (5, nil)", result, err)
	}
}

func TestCalculateErrors(t *testing.T) {
	service := NewService(calculationStrategyResolverStub{
		"multiply":          calculation.MultiplicationStrategy{},
		"non_finite_result": nonFiniteResultStrategy{},
	})

	tests := []struct {
		name      string
		operation string
		operands  []float64
		wantCode  string
	}{
		{name: "unknown operation", operation: "modulo", operands: []float64{4, 2}, wantCode: calculation.ErrorCodeUnknownOperation},
		{name: "non finite operand", operation: "multiply", operands: []float64{math.Inf(1), 1}, wantCode: calculation.ErrorCodeValidation},
		{name: "non finite result", operation: "non_finite_result", operands: []float64{2, 2}, wantCode: calculation.ErrorCodeNonFiniteResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Calculate(tt.operation, tt.operands)
			var calculationError *calculation.CalculationError
			if !errors.As(err, &calculationError) {
				t.Fatalf("Calculate() error = %T, want *calculation.CalculationError", err)
			}
			if calculationError.Code != tt.wantCode {
				t.Fatalf("error code = %q, want %q", calculationError.Code, tt.wantCode)
			}
		})
	}
}
