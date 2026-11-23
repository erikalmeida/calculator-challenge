package calculation

// SubtractionStrategy implements sequential subtraction for two or more operands.
type SubtractionStrategy struct{}

func (SubtractionStrategy) Execute(operands []float64) (float64, error) {
	if len(operands) < 2 {
		return 0, NewInvalidOperandCountError("subtract", "al menos 2 operandos")
	}

	result := operands[0]
	for _, operand := range operands[1:] {
		result -= operand
	}
	return result, nil
}
