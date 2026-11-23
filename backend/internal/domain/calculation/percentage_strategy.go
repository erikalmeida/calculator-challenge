package calculation

// PercentageStrategy calculates a percentage of an amount using amount * percentage / 100.
type PercentageStrategy struct{}

func (PercentageStrategy) Execute(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, NewInvalidOperandCountError("percentage", "exactly 2 operands")
	}

	return operands[0] * operands[1] / 100, nil
}
