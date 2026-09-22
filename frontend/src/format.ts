import type { OperationMeta } from './operations';

/**
 * Formats a result for display, trimming binary floating-point noise
 * (0.1 + 0.2 would otherwise render as 0.30000000000000004) while keeping
 * legitimately long decimals intact.
 */
export function formatResult(value: number): string {
  if (Number.isInteger(value)) {
    return String(value);
  }
  return String(Number(value.toPrecision(12)));
}

/** Human-readable echo of the calculation that produced the current result. */
export function describeExpression(
  operation: OperationMeta,
  a: number,
  b?: number,
): string {
  if (operation.unary) {
    return `${operation.symbol}${a}`;
  }
  if (operation.id === 'percentage') {
    return `${a}% of ${b}`;
  }
  if (operation.id === 'power') {
    return `${a} ^ ${b}`;
  }
  return `${a} ${operation.symbol} ${b}`;
}
