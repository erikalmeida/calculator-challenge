interface OperandInputListProps {
  operands: string[]
  canAddOperand: boolean
  canRemoveOperand: boolean
  hasReachedOperandLimit: boolean
  disabled: boolean
  onOperandChange: (index: number, value: string) => void
  onOperandRemove: (index: number) => void
  onOperandAdd: () => void
}

export function OperandInputList({
  operands,
  canAddOperand,
  canRemoveOperand,
  disabled,
  hasReachedOperandLimit,
  onOperandChange,
  onOperandRemove,
  onOperandAdd,
}: OperandInputListProps) {
  return (
    <fieldset className="operand-list">
      <legend>Operands</legend>
      <p className="form-hint">Use integers or decimals with an optional leading minus sign.</p>

      {operands.map((operand, index) => {
        const inputId = `operand-${index + 1}`
        return (
          <div className="operand-row" key={inputId}>
            <div className="operand-field">
              <label className="operand-label" htmlFor={inputId}>
                Operand {index + 1}
              </label>
              <input
                aria-label={`Operand ${index + 1}`}
                className="operand-input"
                disabled={disabled}
                id={inputId}
                inputMode="decimal"
                maxLength={17}
                onChange={(event) => onOperandChange(index, event.target.value)}
                placeholder="Example: 12.5"
                type="text"
                value={operand}
              />
            </div>
            {canRemoveOperand && (
              <button
                className="button button-tertiary"
                disabled={disabled}
                onClick={() => onOperandRemove(index)}
                type="button"
              >
                Remove
              </button>
            )}
          </div>
        )
      })}

      {canAddOperand && (
        <div className="operand-actions">
          <button className="button button-secondary" disabled={disabled || hasReachedOperandLimit} onClick={onOperandAdd} type="button">
            Add operand
          </button>
          {hasReachedOperandLimit && <p className="form-hint">You can add up to 20 operands.</p>}
        </div>
      )}
    </fieldset>
  )
}
