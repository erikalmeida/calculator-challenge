import { useCalculator } from '../hooks/useCalculator'
import { CalculatorForm } from './CalculatorForm'

export function CalculatorPage() {
  const calculator = useCalculator()

  return (
    <main className="calculator-page">
      <section aria-labelledby="calculator-title" className="calculator-card">
        <p className="calculator-eyebrow">Sezzle take-home</p>
        <h1 id="calculator-title">Calculator</h1>
        <p className="calculator-description">
          Choose an operation, enter your numbers, and see the result.
        </p>
        <CalculatorForm
          onCalculate={calculator.calculate}
          onOperandAdd={calculator.addOperand}
          onOperandChange={calculator.changeOperand}
          onOperandRemove={calculator.removeOperand}
          onOperationChange={calculator.changeOperation}
          onReset={calculator.reset}
          state={calculator.state}
        />
      </section>
    </main>
  )
}
