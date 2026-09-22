import { describe, expect, it } from 'vitest';
import { describeExpression, formatResult } from './format';
import { getOperation } from './operations';

describe('formatResult', () => {
  it('renders integers without decimals', () => {
    expect(formatResult(6)).toBe('6');
    expect(formatResult(-4)).toBe('-4');
    expect(formatResult(0)).toBe('0');
  });

  it('keeps meaningful decimals', () => {
    expect(formatResult(2.5)).toBe('2.5');
    expect(formatResult(1 / 3)).toBe('0.333333333333');
  });

  it('trims floating-point noise', () => {
    expect(formatResult(0.1 + 0.2)).toBe('0.3');
  });
});

describe('describeExpression', () => {
  it('describes binary operations with their symbol', () => {
    expect(describeExpression(getOperation('add'), 4, 2)).toBe('4 + 2');
    expect(describeExpression(getOperation('divide'), 10, 4)).toBe('10 \u00f7 4');
  });

  it('describes unary operations with a single operand', () => {
    expect(describeExpression(getOperation('sqrt'), 81)).toBe('\u221a81');
  });

  it('spells out power and percentage', () => {
    expect(describeExpression(getOperation('power'), 2, 8)).toBe('2 ^ 8');
    expect(describeExpression(getOperation('percentage'), 10, 200)).toBe('10% of 200');
  });
});
