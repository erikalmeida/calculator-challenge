import {
  MAX_DECIMAL_PLACES,
  MAX_OPERAND_ABSOLUTE_VALUE,
  MAX_OPERAND_INPUT_LENGTH,
  MAX_OPERANDS,
  MAX_SIGNIFICANT_DIGITS,
  MIN_NON_ZERO_OPERAND_ABSOLUTE_VALUE,
  type CalculatorOperation,
} from './calculator.types'

export type ValidationResult =
  | { isValid: true; operands: number[] }
  | { isValid: false; errorMessage: string }

const intermediateNumberPattern = /^-?(?:0|[1-9]\d*)?(?:\.\d*)?$/
const completeDecimalNumberPattern = /^-?(?:0|[1-9]\d*)(?:\.\d+)?$/

export function isAllowedNumericInput(value: string): boolean {
  return (
    value.length <= MAX_OPERAND_INPUT_LENGTH &&
    value !== '.' &&
    value !== '-.' &&
    intermediateNumberPattern.test(value)
  )
}

export function supportsMultipleOperands(operation: CalculatorOperation): boolean {
  return operation !== 'divide'
}

export function canAddOperand(operation: CalculatorOperation, operandCount: number): boolean {
  return supportsMultipleOperands(operation) && operandCount < MAX_OPERANDS
}

export function validateCalculationInput(
  operation: CalculatorOperation,
  rawOperands: string[],
): ValidationResult {
  if (rawOperands.length > MAX_OPERANDS) {
    return { isValid: false, errorMessage: `You can add up to ${MAX_OPERANDS} operands.` }
  }

  if (operation === 'divide' && rawOperands.length !== 2) {
    return { isValid: false, errorMessage: 'Division requires exactly two operands.' }
  }

  if (operation !== 'divide' && rawOperands.length < 2) {
    return { isValid: false, errorMessage: 'This operation requires at least two operands.' }
  }

  for (const rawOperand of rawOperands) {
    const validationError = validateRawOperand(rawOperand)
    if (validationError) {
      return { isValid: false, errorMessage: validationError }
    }
  }

  return { isValid: true, operands: rawOperands.map(Number) }
}

function validateRawOperand(rawOperand: string): string | null {
  if (rawOperand === '') {
    return 'Enter a value for every operand.'
  }
  if (!completeDecimalNumberPattern.test(rawOperand)) {
    return 'Use digits, an optional leading minus sign, and a period decimal separator.'
  }
  if (rawOperand.length > MAX_OPERAND_INPUT_LENGTH) {
    return 'Each operand exceeds the maximum supported input length.'
  }

  const decimalPart = rawOperand.split('.')[1] ?? ''
  if (decimalPart.length > MAX_DECIMAL_PLACES) {
    return `Operands can contain at most ${MAX_DECIMAL_PLACES} decimal places.`
  }
  if (significantDigitCount(rawOperand) > MAX_SIGNIFICANT_DIGITS) {
    return `Operands can contain at most ${MAX_SIGNIFICANT_DIGITS} significant digits.`
  }

  const operand = Number(rawOperand)
  if (!Number.isFinite(operand)) {
    return 'Enter finite numbers only.'
  }

  const absoluteOperand = Math.abs(operand)
  if (
    absoluteOperand > MAX_OPERAND_ABSOLUTE_VALUE ||
    (absoluteOperand > 0 && absoluteOperand < MIN_NON_ZERO_OPERAND_ABSOLUTE_VALUE)
  ) {
    return 'Each operand must be zero or be within the supported numeric range.'
  }

  return null
}

function significantDigitCount(value: string): number {
  return value.replace(/[-.]/g, '').replace(/^0+/, '').length
}
