import type { CalculationStatus } from '../model/calculator.types'

interface CalculationResultProps {
  errorMessage: string | null
  noticeMessage: string | null
  result: number | null
  status: CalculationStatus
}

export function CalculationResult({ errorMessage, noticeMessage, result, status }: CalculationResultProps) {
  if (status === 'idle' && !noticeMessage) {
    return null
  }

  if (noticeMessage) {
    return (
      <section aria-live="polite" className="calculation-feedback calculation-feedback-loading">
        <p>{noticeMessage}</p>
      </section>
    )
  }

  if (status === 'loading') {
    return (
      <section aria-live="polite" className="calculation-feedback calculation-feedback-loading">
        <p>Calculating…</p>
      </section>
    )
  }

  if (status === 'error') {
    return (
      <section aria-live="assertive" className="calculation-feedback calculation-feedback-error" role="alert">
        <p>{errorMessage}</p>
      </section>
    )
  }

  return (
    <section aria-live="polite" className="calculation-feedback calculation-feedback-success">
      <p>
        Result: <strong>{result}</strong>
      </p>
    </section>
  )
}
