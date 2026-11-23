export const calculatorOperations = ['add', 'subtract', 'multiply', 'divide', 'percentage'] as const

export const MAX_OPERANDS = 20
export const MAX_OPERAND_ABSOLUTE_VALUE = 100_000_000_000_000
export const MIN_NON_ZERO_OPERAND_ABSOLUTE_VALUE = 0.0000000001
export const MAX_DECIMAL_PLACES = 10
export const MAX_SIGNIFICANT_DIGITS = 15
export const MAX_OPERAND_INPUT_LENGTH = 17

export type CalculatorOperation = (typeof calculatorOperations)[number]

export function requiresExactlyTwoOperands(operation: CalculatorOperation): boolean {
  return operation === 'divide' || operation === 'percentage'
}

export type CalculationStatus = 'idle' | 'loading' | 'success' | 'error'

export interface CalculateRequest {
  operation: CalculatorOperation
  operands: number[]
}

export interface CalculateResponse {
  operation: CalculatorOperation
  operands: number[]
  result: number
}

export interface CalculatorState {
  operation: CalculatorOperation
  operands: string[]
  status: CalculationStatus
  result: number | null
  errorMessage: string | null
  noticeMessage: string | null
}

export type CalculatorAction =
  | { type: 'operationChanged'; operation: CalculatorOperation }
  | { type: 'operandChanged'; index: number; value: string }
  | { type: 'operandAdded' }
  | { type: 'operandRemoved'; index: number }
  | { type: 'calculationStarted' }
  | { type: 'calculationSucceeded'; result: number }
  | { type: 'calculationFailed'; errorMessage: string }
  | { type: 'formReset' }
