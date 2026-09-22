export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage';

export interface OperationMeta {
  id: Operation;
  label: string;
  symbol: string;
  /** Short reminder of how the operands are combined. */
  hint: string;
  /** Unary operations only use the first operand. */
  unary?: boolean;
}

export const OPERATIONS: readonly OperationMeta[] = [
  { id: 'add', label: 'Add', symbol: '+', hint: 'a + b' },
  { id: 'subtract', label: 'Subtract', symbol: '\u2212', hint: 'a \u2212 b' },
  { id: 'multiply', label: 'Multiply', symbol: '\u00d7', hint: 'a \u00d7 b' },
  { id: 'divide', label: 'Divide', symbol: '\u00f7', hint: 'a \u00f7 b' },
  { id: 'power', label: 'Power', symbol: 'x\u02b8', hint: 'a to the power of b' },
  { id: 'sqrt', label: 'Square root', symbol: '\u221a', hint: '\u221aa', unary: true },
  { id: 'percentage', label: 'Percentage', symbol: '%', hint: 'a% of b' },
];

export function getOperation(id: Operation): OperationMeta {
  const match = OPERATIONS.find((operation) => operation.id === id);
  if (!match) {
    throw new Error(`Unknown operation: ${id}`);
  }
  return match;
}
