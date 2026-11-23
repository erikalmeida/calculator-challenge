import type { FormEvent } from 'react'

import { calculatorOperations, type CalculatorOperation, type CalculatorState } from '../model/calculator.types'
import { MAX_OPERANDS } from '../model/calculator.types'
import { supportsMultipleOperands } from '../model/calculatorValidation'
import { CalculationResult } from './CalculationResult'
import { OperandInputList } from './OperandInputList'

interface CalculatorFormProps {
  state: CalculatorState
  onCalculate: () => void
  onOperationChange: (operation: CalculatorOperation) => void
  onOperandAdd: () => void
  onOperandChange: (index: number, value: string) => void
  onOperandRemove: (index: number) => void
  onReset: () => void
}

const operationLabels: Record<CalculatorOperation, string> = {
  add: 'Addition',
  subtract: 'Subtraction',
  multiply: 'Multiplication',
  divide: 'Division',
  percentage: 'Percentage',
}

export function CalculatorForm({
  state,
  onCalculate,
  onOperationChange,
  onOperandAdd,
  onOperandChange,
  onOperandRemove,
  onReset,
}: CalculatorFormProps) {
  const isLoading = state.status === 'loading'
  const canManageMultipleOperands = supportsMultipleOperands(state.operation)
  const hasReachedOperandLimit = state.operands.length >= MAX_OPERANDS

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    onCalculate()
  }

  return (
    <form className="calculator-form" noValidate onSubmit={handleSubmit}>
      <div className="form-field">
        <label className="form-label" htmlFor="calculator-operation">
          Operation
        </label>
        <select
          className="operation-select"
          disabled={isLoading}
          id="calculator-operation"
          onChange={(event) => onOperationChange(event.target.value as CalculatorOperation)}
          value={state.operation}
        >
          {calculatorOperations.map((operation) => (
            <option key={operation} value={operation}>
              {operationLabels[operation]}
            </option>
          ))}
        </select>
      </div>

      <OperandInputList
        canAddOperand={canManageMultipleOperands}
        canRemoveOperand={canManageMultipleOperands && state.operands.length > 2}
        disabled={isLoading}
        hasReachedOperandLimit={hasReachedOperandLimit}
        onOperandAdd={onOperandAdd}
        onOperandChange={onOperandChange}
        onOperandRemove={onOperandRemove}
        operands={state.operands}
      />

      <div className="form-actions">
        <button className="button button-primary" disabled={isLoading} type="submit">
          {isLoading ? 'Calculating…' : 'Calculate'}
        </button>
        <button className="button button-secondary" disabled={isLoading} onClick={onReset} type="button">
          Clear
        </button>
      </div>

      <CalculationResult
        errorMessage={state.errorMessage}
        noticeMessage={state.noticeMessage}
        result={state.result}
        status={state.status}
      />
    </form>
  )
}
