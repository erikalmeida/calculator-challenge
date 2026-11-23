import { afterEach, describe, expect, it, vi } from 'vitest'

import { CalculatorApiError, calculateWithApi } from './calculatorApi'

describe('calculateWithApi', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('sends the calculator contract and returns a successful result', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ operation: 'multiply', operands: [2, 3], result: 6 }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )
    vi.stubGlobal('fetch', fetchMock)

    await expect(calculateWithApi({ operation: 'multiply', operands: [2, 3] })).resolves.toEqual({
      operation: 'multiply',
      operands: [2, 3],
      result: 6,
    })

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation: 'multiply', operands: [2, 3] }),
      signal: undefined,
    })
  })

  it('uses the backend error message when the calculation is rejected', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({ error: { code: 'division_by_zero', message: 'The divisor cannot be zero.' } }),
          { status: 400, headers: { 'Content-Type': 'application/json' } },
        ),
      ),
    )

    await expect(calculateWithApi({ operation: 'divide', operands: [4, 0] })).rejects.toThrow(
      'The divisor cannot be zero.',
    )
  })

  it('reports network failures without exposing low-level errors', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))

    await expect(calculateWithApi({ operation: 'add', operands: [2, 3] })).rejects.toMatchObject<Partial<CalculatorApiError>>({
      code: 'network_error',
      message: 'Unable to connect to the calculator service.',
      status: null,
    })
  })

  it('rejects an invalid successful API response', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ operation: 'add', operands: [2, 3] }), {
          status: 200,
          headers: { 'Content-Type': 'application/json', 'X-Request-ID': 'request-123' },
        }),
      ),
    )

    await expect(calculateWithApi({ operation: 'add', operands: [2, 3] })).rejects.toMatchObject<Partial<CalculatorApiError>>({
      code: 'invalid_api_response',
      requestID: 'request-123',
      status: 200,
    })
  })
})
