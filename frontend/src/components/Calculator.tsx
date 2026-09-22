import { useState, type FormEvent } from 'react';
import { ApiError, calculate } from '../api/client';
import { describeExpression } from '../format';
import { getOperation, type Operation } from '../operations';
import { parseOperand } from '../validation';
import { Display } from './Display';
import { OperationPad } from './OperationPad';

export function Calculator() {
  const [operation, setOperation] = useState<Operation>('add');
  const [a, setA] = useState('');
  const [b, setB] = useState('');
  const [result, setResult] = useState<number | null>(null);
  const [expression, setExpression] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  const meta = getOperation(operation);
  const isUnary = meta.unary === true;

  function selectOperation(next: Operation) {
    setOperation(next);
    // A result belongs to the operation that produced it.
    setResult(null);
    setExpression(null);
    setError(null);
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    setResult(null);
    setExpression(null);

    const parsedA = parseOperand(a);
    if (!parsedA.ok) {
      setError(`First value ${parsedA.error}.`);
      return;
    }

    let operandB: number | undefined;
    if (!isUnary) {
      const parsedB = parseOperand(b);
      if (!parsedB.ok) {
        setError(`Second value ${parsedB.error}.`);
        return;
      }
      operandB = parsedB.value;
    }

    setPending(true);
    try {
      const response = await calculate(operation, parsedA.value, operandB);
      setResult(response.result);
      setExpression(describeExpression(meta, parsedA.value, operandB));
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : 'Something went wrong. Please try again.',
      );
    } finally {
      setPending(false);
    }
  }

  function handleClear() {
    setA('');
    setB('');
    setResult(null);
    setExpression(null);
    setError(null);
  }

  return (
    <form className="calculator" onSubmit={handleSubmit} noValidate>
      <Display result={result} error={error} expression={expression} />

      <OperationPad value={operation} onChange={selectOperation} disabled={pending} />

      <p className="calculator__hint">{meta.hint}</p>

      <div className="fields">
        <label className="field">
          <span className="field__label">First value (a)</span>
          <input
            className="field__input"
            type="text"
            inputMode="decimal"
            autoComplete="off"
            value={a}
            onChange={(event) => setA(event.target.value)}
          />
        </label>

        <label className="field" hidden={isUnary}>
          <span className="field__label">Second value (b)</span>
          <input
            className="field__input"
            type="text"
            inputMode="decimal"
            autoComplete="off"
            value={b}
            onChange={(event) => setB(event.target.value)}
            disabled={isUnary}
          />
        </label>
      </div>

      <div className="actions">
        <button className="button button--primary" type="submit" disabled={pending}>
          {pending ? 'Calculating...' : 'Calculate'}
        </button>
        <button className="button" type="button" onClick={handleClear} disabled={pending}>
          Clear
        </button>
      </div>
    </form>
  );
}
