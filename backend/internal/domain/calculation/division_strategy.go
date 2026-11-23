package calculation

// DivisionStrategy implements division for exactly two operands.
type DivisionStrategy struct{}

func (DivisionStrategy) Execute(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, NewInvalidOperandCountError("divide", "exactamente 2 operandos")
	}
	if operands[1] == 0 {
		return 0, &CalculationError{Code: ErrorCodeDivisionByZero, Message: "The divisor cannot be zero."}
	}
	return operands[0] / operands[1], nil
}
