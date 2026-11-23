import { calculatorOperations, type CalculateRequest, type CalculateResponse } from '../model/calculator.types'

interface ApiErrorResponse {
  error?: {
    code?: string
    message?: string
  }
}

export class CalculatorApiError extends Error {
  readonly code: string
  readonly requestID: string | null
  readonly status: number | null

  constructor({
    message,
    code = 'unknown_error',
    status = null,
    requestID = null,
  }: {
    message: string
    code?: string
    status?: number | null
    requestID?: string | null
  }) {
    super(message)
    this.name = 'CalculatorApiError'
    this.code = code
    this.status = status
    this.requestID = requestID
  }
}

export async function calculateWithApi(
  request: CalculateRequest,
  signal?: AbortSignal,
): Promise<CalculateResponse> {
  let response: Response

  try {
    response = await fetch('/api/v1/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
      signal,
    })
  } catch (error) {
    if (isAbortError(error)) {
      throw new CalculatorApiError({
        message: 'The calculator request timed out. Please try again.',
        code: 'request_timeout',
      })
    }
    throw new CalculatorApiError({
      message: 'Unable to connect to the calculator service.',
      code: 'network_error',
    })
  }

  const requestID = response.headers.get('X-Request-ID')
  const payload = (await response.json().catch(() => null)) as unknown

  if (!response.ok) {
    const apiError = payload as ApiErrorResponse | null
    throw new CalculatorApiError({
      message: apiError?.error?.message ?? 'Unable to complete the calculation.',
      code: apiError?.error?.code ?? 'api_error',
      status: response.status,
      requestID,
    })
  }

  if (!isCalculateResponse(payload)) {
    throw new CalculatorApiError({
      message: 'The calculator service returned an invalid response.',
      code: 'invalid_api_response',
      status: response.status,
      requestID,
    })
  }

  return payload
}

function isAbortError(error: unknown): boolean {
  return error instanceof Error && error.name === 'AbortError'
}

function isCalculateResponse(payload: unknown): payload is CalculateResponse {
  if (!payload || typeof payload !== 'object') {
    return false
  }

  const response = payload as Partial<CalculateResponse>
  return (
    typeof response.operation === 'string' &&
    calculatorOperations.includes(response.operation as CalculateResponse['operation']) &&
    Array.isArray(response.operands) &&
    response.operands.every((operand) => typeof operand === 'number' && Number.isFinite(operand)) &&
    typeof response.result === 'number' &&
    Number.isFinite(response.result)
  )
}
