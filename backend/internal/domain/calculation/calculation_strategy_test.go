package calculation

import (
	"errors"
	"testing"
)

func TestOperationStrategies(t *testing.T) {
	tests := []struct {
		name     string
		strategy CalculationStrategy
		operands []float64
		want     float64
	}{
		{name: "add", strategy: AdditionStrategy{}, operands: []float64{1, 2, 3}, want: 6},
		{name: "subtract", strategy: SubtractionStrategy{}, operands: []float64{10, 3, 2}, want: 5},
		{name: "multiply", strategy: MultiplicationStrategy{}, operands: []float64{2, 3, 4}, want: 24},
		{name: "divide", strategy: DivisionStrategy{}, operands: []float64{8, 2}, want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.strategy.Execute(tt.operands)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationStrategyErrors(t *testing.T) {
	tests := []struct {
		name     string
		strategy CalculationStrategy
		operands []float64
		wantCode string
	}{
		{name: "add arity", strategy: AdditionStrategy{}, operands: []float64{1}, wantCode: ErrorCodeInvalidArity},
		{name: "subtract arity", strategy: SubtractionStrategy{}, operands: []float64{1}, wantCode: ErrorCodeInvalidArity},
		{name: "multiply arity", strategy: MultiplicationStrategy{}, operands: []float64{1}, wantCode: ErrorCodeInvalidArity},
		{name: "divide arity", strategy: DivisionStrategy{}, operands: []float64{8, 2, 1}, wantCode: ErrorCodeInvalidArity},
		{name: "divide by zero", strategy: DivisionStrategy{}, operands: []float64{8, 0}, wantCode: ErrorCodeDivisionByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.strategy.Execute(tt.operands)
			var calculationError *CalculationError
			if !errors.As(err, &calculationError) {
				t.Fatalf("Execute() error = %T, want *CalculationError", err)
			}
			if calculationError.Code != tt.wantCode {
				t.Fatalf("error code = %q, want %q", calculationError.Code, tt.wantCode)
			}
		})
	}
}
