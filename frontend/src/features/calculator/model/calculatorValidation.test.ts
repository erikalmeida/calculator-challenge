import { describe, expect, it } from 'vitest'

import { MAX_OPERANDS } from './calculator.types'
import {
  canAddOperand,
  isAllowedNumericInput,
  supportsMultipleOperands,
  validateCalculationInput,
} from './calculatorValidation'

describe('calculatorValidation', () => {
  it('accepts allowed intermediate input and rejects unsupported characters', () => {
    expect(isAllowedNumericInput('-12.5')).toBe(true)
    expect(isAllowedNumericInput('-')).toBe(true)
    expect(isAllowedNumericInput('-12,5')).toBe(false)
    expect(isAllowedNumericInput('1e3')).toBe(false)
    expect(isAllowedNumericInput('12a')).toBe(false)
  })

  it('accepts decimal values that match the calculator contract', () => {
    expect(validateCalculationInput('add', ['2.5', '3'])).toEqual({
      isValid: true,
      operands: [2.5, 3],
    })
  })

  it('rejects empty values before calling the API', () => {
    expect(validateCalculationInput('multiply', ['2', ''])).toEqual({
      isValid: false,
      errorMessage: 'Enter a value for every operand.',
    })
  })

  it('requires exactly two operands for division', () => {
    expect(validateCalculationInput('divide', ['10', '2', '1'])).toEqual({
      isValid: false,
      errorMessage: 'Division requires exactly two operands.',
    })
  })

  it('requires exactly two operands for percentage', () => {
    expect(validateCalculationInput('percentage', ['200', '15', '5'])).toEqual({
      isValid: false,
      errorMessage: 'Percentage requires exactly two operands.',
    })
    expect(supportsMultipleOperands('percentage')).toBe(false)
  })

  it('rejects unsupported precision and values outside the numeric range', () => {
    expect(validateCalculationInput('add', ['0.00000000001', '2'])).toEqual({
      isValid: false,
      errorMessage: 'Operands can contain at most 10 decimal places.',
    })
    expect(validateCalculationInput('add', ['100000000000001', '2'])).toEqual({
      isValid: false,
      errorMessage: 'Each operand must be zero or be within the supported numeric range.',
    })
  })

  it('enforces the twenty operand limit', () => {
    expect(validateCalculationInput('add', Array.from({ length: MAX_OPERANDS + 1 }, () => '1'))).toEqual({
      isValid: false,
      errorMessage: 'You can add up to 20 operands.',
    })
    expect(canAddOperand('add', MAX_OPERANDS - 1)).toBe(true)
    expect(canAddOperand('add', MAX_OPERANDS)).toBe(false)
    expect(supportsMultipleOperands('divide')).toBe(false)
  })
})
