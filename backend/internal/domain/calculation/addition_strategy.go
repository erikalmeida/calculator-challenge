package calculation

// AdditionStrategy implements addition for two or more operands.
type AdditionStrategy struct{}

func (AdditionStrategy) Execute(operands []float64) (float64, error) {
	if len(operands) < 2 {
		return 0, NewInvalidOperandCountError("add", "at least 2 operands")
	}

	var result float64
	for _, operand := range operands {
		result += operand
	}
	return result, nil
}
