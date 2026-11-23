import { useCallback, useReducer, useRef } from 'react'

import { CalculatorApiError, calculateWithApi } from '../api/calculatorApi'
import { calculatorReducer, initialCalculatorState } from '../model/calculatorReducer'
import type { CalculatorOperation } from '../model/calculator.types'
import {
  canAddOperand,
  isAllowedNumericInput,
  validateCalculationInput,
} from '../model/calculatorValidation'

const calculatorRequestTimeoutMilliseconds = 10_000

export function useCalculator() {
  const [state, dispatch] = useReducer(calculatorReducer, initialCalculatorState)
  const activeRequestController = useRef<AbortController | null>(null)
  const requestSequence = useRef(0)

  const changeOperation = useCallback((operation: CalculatorOperation) => {
    dispatch({ type: 'operationChanged', operation })
  }, [])

  const changeOperand = useCallback((index: number, value: string) => {
    if (isAllowedNumericInput(value)) {
      dispatch({ type: 'operandChanged', index, value })
    }
  }, [])

  const addOperand = useCallback(() => {
    if (canAddOperand(state.operation, state.operands.length)) {
      dispatch({ type: 'operandAdded' })
    }
  }, [state.operation, state.operands.length])

  const removeOperand = useCallback(
    (index: number) => {
      if (state.operands.length > 2) {
        dispatch({ type: 'operandRemoved', index })
      }
    },
    [state.operands.length],
  )

  const calculate = useCallback(async () => {
    const validation = validateCalculationInput(state.operation, state.operands)
    if (!validation.isValid) {
      dispatch({ type: 'calculationFailed', errorMessage: validation.errorMessage })
      return
    }

    activeRequestController.current?.abort()
    const requestController = new AbortController()
    activeRequestController.current = requestController
    const currentRequestSequence = requestSequence.current + 1
    requestSequence.current = currentRequestSequence
    const timeoutID = window.setTimeout(() => requestController.abort(), calculatorRequestTimeoutMilliseconds)

    dispatch({ type: 'calculationStarted' })
    try {
      const response = await calculateWithApi(
        { operation: state.operation, operands: validation.operands },
        requestController.signal,
      )
      if (requestSequence.current === currentRequestSequence) {
        dispatch({ type: 'calculationSucceeded', result: response.result })
      }
    } catch (error) {
      if (requestSequence.current === currentRequestSequence) {
        const errorMessage = error instanceof Error ? error.message : 'An unexpected error occurred.'
        logCalculatorError(error)
        dispatch({ type: 'calculationFailed', errorMessage })
      }
    } finally {
      window.clearTimeout(timeoutID)
      if (activeRequestController.current === requestController) {
        activeRequestController.current = null
      }
    }
  }, [state.operation, state.operands])

  const reset = useCallback(() => {
    requestSequence.current += 1
    activeRequestController.current?.abort()
    activeRequestController.current = null
    dispatch({ type: 'formReset' })
  }, [])

  return { state, changeOperation, changeOperand, addOperand, removeOperand, calculate, reset }
}

function logCalculatorError(error: unknown) {
  if (!import.meta.env.DEV) {
    return
  }

  if (error instanceof CalculatorApiError) {
    console.error('Calculator request failed', {
      code: error.code,
      status: error.status,
      requestID: error.requestID,
    })
    return
  }

  console.error('Calculator request failed with an unexpected client error')
}
