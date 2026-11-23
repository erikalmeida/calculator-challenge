import { requiresExactlyTwoOperands, type CalculatorAction, type CalculatorState } from './calculator.types'

export const initialCalculatorState: CalculatorState = {
  operation: 'add',
  operands: ['', ''],
  status: 'idle',
  result: null,
  errorMessage: null,
  noticeMessage: null,
}

export function calculatorReducer(
  state: CalculatorState,
  action: CalculatorAction,
): CalculatorState {
  switch (action.type) {
    case 'operationChanged': {
      const requiresTwoOperands = requiresExactlyTwoOperands(action.operation)
      const removesExtraOperands = requiresTwoOperands && state.operands.length > 2
      return {
        ...state,
        operation: action.operation,
        operands: requiresTwoOperands ? state.operands.slice(0, 2) : state.operands,
        status: 'idle',
        result: null,
        errorMessage: null,
        noticeMessage: removesExtraOperands
          ? `${action.operation === 'percentage' ? 'Percentage' : 'Division'} supports only two operands. Extra operands were removed.`
          : null,
      }
    }
    case 'operandChanged': {
      const operands = [...state.operands]
      operands[action.index] = action.value
      return { ...state, operands, status: 'idle', result: null, errorMessage: null, noticeMessage: null }
    }
    case 'operandAdded':
      return { ...state, operands: [...state.operands, ''], status: 'idle', result: null, errorMessage: null, noticeMessage: null }
    case 'operandRemoved':
      return {
        ...state,
        operands: state.operands.filter((_, index) => index !== action.index),
        status: 'idle',
        result: null,
        errorMessage: null,
        noticeMessage: null,
      }
    case 'calculationStarted':
      return { ...state, status: 'loading', result: null, errorMessage: null, noticeMessage: null }
    case 'calculationSucceeded':
      return { ...state, status: 'success', result: action.result, errorMessage: null, noticeMessage: null }
    case 'calculationFailed':
      return { ...state, status: 'error', result: null, errorMessage: action.errorMessage, noticeMessage: null }
    case 'formReset':
      return { ...initialCalculatorState }
  }
}
