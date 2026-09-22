import { describe, expect, it } from 'vitest';
import { parseOperand } from './validation';

describe('parseOperand', () => {
  it('accepts integers, decimals and negatives', () => {
    expect(parseOperand('42')).toEqual({ ok: true, value: 42 });
    expect(parseOperand('3.5')).toEqual({ ok: true, value: 3.5 });
    expect(parseOperand('-12')).toEqual({ ok: true, value: -12 });
    expect(parseOperand('.5')).toEqual({ ok: true, value: 0.5 });
    expect(parseOperand('0')).toEqual({ ok: true, value: 0 });
  });

  it('ignores surrounding whitespace', () => {
    expect(parseOperand('  7  ')).toEqual({ ok: true, value: 7 });
  });

  it('rejects empty input', () => {
    expect(parseOperand('')).toEqual({ ok: false, error: 'is required' });
    expect(parseOperand('   ')).toEqual({ ok: false, error: 'is required' });
  });

  it('rejects values that are not plain numbers', () => {
    for (const input of ['abc', '1,5', '5px', '--3', '1.2.3', '1e5', '-']) {
      expect(parseOperand(input)).toEqual({ ok: false, error: 'must be a valid number' });
    }
  });
});
