import { describe, expect, it } from 'vitest'

import { calculatorReducer, initialCalculatorState } from './calculatorReducer'

describe('calculatorReducer', () => {
  it('warns and trims extra operands when switching to division', () => {
    const state = {
      ...initialCalculatorState,
      operation: 'add' as const,
      operands: ['1', '2', '3'],
    }

    expect(calculatorReducer(state, { type: 'operationChanged', operation: 'divide' })).toMatchObject({
      operation: 'divide',
      operands: ['1', '2'],
      noticeMessage: 'Division supports only two operands. Extra operands were removed.',
    })
  })

  it('updates operands and resets calculator state transitions', () => {
    const changed = calculatorReducer(initialCalculatorState, {
      type: 'operandChanged',
      index: 0,
      value: '12',
    })
    const added = calculatorReducer(changed, { type: 'operandAdded' })
    const removed = calculatorReducer(added, { type: 'operandRemoved', index: 2 })
    const loading = calculatorReducer(removed, { type: 'calculationStarted' })
    const success = calculatorReducer(loading, { type: 'calculationSucceeded', result: 12 })
    const failure = calculatorReducer(success, {
      type: 'calculationFailed',
      errorMessage: 'Calculator service unavailable.',
    })

    expect(changed.operands).toEqual(['12', ''])
    expect(added.operands).toHaveLength(3)
    expect(removed.operands).toHaveLength(2)
    expect(loading.status).toBe('loading')
    expect(success).toMatchObject({ status: 'success', result: 12, errorMessage: null })
    expect(failure).toMatchObject({ status: 'error', result: null, errorMessage: 'Calculator service unavailable.' })
    expect(calculatorReducer(failure, { type: 'formReset' })).toEqual(initialCalculatorState)
  })
})
