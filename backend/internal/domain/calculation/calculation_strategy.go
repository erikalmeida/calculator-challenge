package calculation

// CalculationStrategy is the Strategy port for one mathematical operation.
// Implementations are pure: they know neither HTTP nor the application service.
type CalculationStrategy interface {
	Execute(operands []float64) (float64, error)
}
