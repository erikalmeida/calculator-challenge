// Package calculation contains the business rules of the calculator.
package calculation

import (
	"fmt"
	"math"
)

const (
	ErrorCodeUnknownOperation  = "unknown_operation"
	ErrorCodeInvalidArity      = "invalid_arity"
	ErrorCodeDivisionByZero    = "division_by_zero"
	ErrorCodeNonFiniteResult   = "non_finite_result"
	ErrorCodeValidation        = "validation_error"
	ErrorCodeTooManyOperands   = "too_many_operands"
	ErrorCodeOperandOutOfRange = "operand_out_of_range"
)

// CalculationError represents a business error that can be translated by an inbound adapter.
type CalculationError struct {
	Code    string
	Message string
}

func (e *CalculationError) Error() string { return e.Message }

// NewInvalidOperandCountError creates the standard error for an operation with a wrong number of operands.
func NewInvalidOperandCountError(operation, requirement string) *CalculationError {
	return &CalculationError{
		Code:    ErrorCodeInvalidArity,
		Message: fmt.Sprintf("The %q operation requires %s.", operation, requirement),
	}
}

// ValidateOperandValues rejects an oversized list and values outside the calculator contract.
func ValidateOperandValues(operands []float64) error {
	if len(operands) > MaxOperands {
		return &CalculationError{
			Code:    ErrorCodeTooManyOperands,
			Message: fmt.Sprintf("A calculation can contain at most %d operands.", MaxOperands),
		}
	}

	for _, operand := range operands {
		if err := ValidateOperandValue(operand); err != nil {
			return err
		}
	}
	return nil
}

// ValidateOperandValue rejects values that cannot be safely calculated or presented.
func ValidateOperandValue(operand float64) error {
	if math.IsNaN(operand) || math.IsInf(operand, 0) {
		return &CalculationError{Code: ErrorCodeValidation, Message: "Operands must be finite numbers."}
	}

	absOperand := math.Abs(operand)
	if absOperand > MaxOperandAbsoluteValue || (absOperand > 0 && absOperand < MinNonZeroOperandAbsoluteValue) {
		return &CalculationError{
			Code:    ErrorCodeOperandOutOfRange,
			Message: "Each operand must be zero or be within the supported numeric range.",
		}
	}
	return nil
}

// ValidateFiniteResult rejects overflow and other non-finite calculation results.
func ValidateFiniteResult(result float64) error {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return &CalculationError{Code: ErrorCodeNonFiniteResult, Message: "The calculation produced a non-finite result."}
	}
	return nil
}
