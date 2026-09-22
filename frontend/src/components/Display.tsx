import { formatResult } from '../format';

interface Props {
  result: number | null;
  error: string | null;
  expression: string | null;
}

export function Display({ result, error, expression }: Props) {
  if (error) {
    return (
      <div className="display display--error" role="alert">
        {error}
      </div>
    );
  }

  return (
    <div className="display" aria-live="polite">
      {result === null ? (
        <span className="display__placeholder">Pick an operation and enter values</span>
      ) : (
        <>
          {expression && <span className="display__expression">{expression}</span>}
          <output className="display__value">{formatResult(result)}</output>
        </>
      )}
    </div>
  );
}
