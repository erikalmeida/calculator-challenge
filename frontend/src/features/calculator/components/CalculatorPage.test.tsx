import { afterEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import { CalculatorPage } from './CalculatorPage'

describe('CalculatorPage', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('shows a validation message without sending an incomplete calculation', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)

    render(<CalculatorPage />)
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Enter a value for every operand.')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('submits valid operands and renders the API result', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ operation: 'add', operands: [2, 3], result: 5 }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )
    vi.stubGlobal('fetch', fetchMock)

    render(<CalculatorPage />)
    await user.type(screen.getByLabelText('Operand 1'), '2')
    await user.type(screen.getByLabelText('Operand 2'), '3')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByText(/Result:/)).toHaveTextContent('Result: 5')
    expect(fetchMock).toHaveBeenCalledOnce()
  })

  it('shows the backend error message when the API rejects a calculation', async () => {
    const user = userEvent.setup()
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({ error: { code: 'division_by_zero', message: 'The divisor cannot be zero.' } }),
          { status: 400, headers: { 'Content-Type': 'application/json', 'X-Request-ID': 'request-456' } },
        ),
      ),
    )

    render(<CalculatorPage />)
    await user.selectOptions(screen.getByLabelText('Operation'), 'divide')
    await user.type(screen.getByLabelText('Operand 1'), '2')
    await user.type(screen.getByLabelText('Operand 2'), '0')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('The divisor cannot be zero.')
  })

  it('disables adding operands after reaching the maximum', async () => {
    const user = userEvent.setup()
    render(<CalculatorPage />)

    const addOperandButton = screen.getByRole('button', { name: 'Add operand' })
    for (let index = 0; index < 18; index += 1) {
      await user.click(addOperandButton)
    }

    expect(screen.getAllByRole('textbox')).toHaveLength(20)
    expect(addOperandButton).toBeDisabled()
    expect(screen.getByText('You can add up to 20 operands.')).toBeInTheDocument()
  })
})
