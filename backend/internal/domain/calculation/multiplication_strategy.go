package calculation

// MultiplicationStrategy implements multiplication for two or more operands.
type MultiplicationStrategy struct{}

func (MultiplicationStrategy) Execute(operands []float64) (float64, error) {
	if len(operands) < 2 {
		return 0, NewInvalidOperandCountError("multiply", "al menos 2 operandos")
	}

	result := 1.0
	for _, operand := range operands {
		result *= operand
	}
	return result, nil
}
