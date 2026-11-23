package calculation

import (
	"errors"
	"math"
	"testing"
)

func TestValidateOperandValues(t *testing.T) {
	if err := ValidateOperandValues([]float64{1, -2.5, 0}); err != nil {
		t.Fatalf("ValidateOperandValues() error = %v, want nil", err)
	}

	err := ValidateOperandValues([]float64{math.Inf(1)})
	assertCalculationErrorCode(t, err, ErrorCodeValidation)
}

func TestValidateOperandValuesRejectsOutOfRangeAndOversizedRequests(t *testing.T) {
	tests := []struct {
		name     string
		operands []float64
		wantCode string
	}{
		{name: "above maximum", operands: []float64{MaxOperandAbsoluteValue + 1}, wantCode: ErrorCodeOperandOutOfRange},
		{name: "below minimum non-zero value", operands: []float64{MinNonZeroOperandAbsoluteValue / 10}, wantCode: ErrorCodeOperandOutOfRange},
		{name: "more than twenty operands", operands: make([]float64, MaxOperands+1), wantCode: ErrorCodeTooManyOperands},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCalculationErrorCode(t, ValidateOperandValues(tt.operands), tt.wantCode)
		})
	}
}

func TestValidateResult(t *testing.T) {
	if err := ValidateFiniteResult(42); err != nil {
		t.Fatalf("ValidateFiniteResult() error = %v, want nil", err)
	}

	err := ValidateFiniteResult(math.NaN())
	assertCalculationErrorCode(t, err, ErrorCodeNonFiniteResult)
}

func assertCalculationErrorCode(t *testing.T, err error, wantCode string) {
	t.Helper()
	var calculationError *CalculationError
	if !errors.As(err, &calculationError) {
		t.Fatalf("error = %T, want *CalculationError", err)
	}
	if calculationError.Code != wantCode {
		t.Fatalf("error code = %q, want %q", calculationError.Code, wantCode)
	}
}
