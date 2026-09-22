export type ParsedOperand =
  | { ok: true; value: number }
  | { ok: false; error: string };

// Accepts integers and decimals, with an optional leading minus: -12, 3.5, .5
const NUMBER_PATTERN = /^-?(\d+(\.\d+)?|\.\d+)$/;

/**
 * Validates raw text from an operand field. The error messages are fragments
 * ("is required") so callers can prefix the field name.
 */
export function parseOperand(raw: string): ParsedOperand {
  const trimmed = raw.trim();

  if (trimmed === '') {
    return { ok: false, error: 'is required' };
  }
  if (!NUMBER_PATTERN.test(trimmed)) {
    return { ok: false, error: 'must be a valid number' };
  }

  const value = Number(trimmed);
  if (!Number.isFinite(value)) {
    return { ok: false, error: 'is out of range' };
  }

  return { ok: true, value };
}
