package calculator

// CalculateUseCase defines the application capability consumed by inbound adapters.
type CalculateUseCase interface {
	Calculate(operation string, operands []float64) (float64, error)
}
