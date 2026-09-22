import { OPERATIONS, type Operation } from '../operations';

interface Props {
  value: Operation;
  onChange: (operation: Operation) => void;
  disabled?: boolean;
}

export function OperationPad({ value, onChange, disabled }: Props) {
  return (
    <div className="pad" role="group" aria-label="Operation">
      {OPERATIONS.map((operation) => (
        <button
          key={operation.id}
          type="button"
          className="pad__button"
          aria-pressed={operation.id === value}
          aria-label={operation.label}
          disabled={disabled}
          onClick={() => onChange(operation.id)}
        >
          <span className="pad__symbol" aria-hidden="true">
            {operation.symbol}
          </span>
          <span className="pad__label">{operation.label}</span>
        </button>
      ))}
    </div>
  );
}
