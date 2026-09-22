import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ApiError, calculate } from '../api/client';
import { Calculator } from './Calculator';

vi.mock('../api/client', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../api/client')>();
  return { ...actual, calculate: vi.fn() };
});

const calculateMock = vi.mocked(calculate);

beforeEach(() => {
  calculateMock.mockReset();
});

async function fillOperands(a: string, b?: string) {
  const user = userEvent.setup();
  await user.type(screen.getByLabelText('First value (a)'), a);
  if (b !== undefined) {
    await user.type(screen.getByLabelText('Second value (b)'), b);
  }
  return user;
}

describe('Calculator', () => {
  it('sends the selected operation and shows the result', async () => {
    calculateMock.mockResolvedValue({ operation: 'add', a: 4, b: 2, result: 6 });
    render(<Calculator />);

    const user = await fillOperands('4', '2');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(calculateMock).toHaveBeenCalledWith('add', 4, 2);
    expect(await screen.findByText('6')).toBeInTheDocument();
    expect(screen.getByText('4 + 2')).toBeInTheDocument();
  });

  it('switches operation before calculating', async () => {
    calculateMock.mockResolvedValue({ operation: 'multiply', a: 3, b: 5, result: 15 });
    render(<Calculator />);

    const user = await fillOperands('3', '5');
    await user.click(screen.getByRole('button', { name: 'Multiply' }));
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(calculateMock).toHaveBeenCalledWith('multiply', 3, 5);
    expect(await screen.findByText('15')).toBeInTheDocument();
  });

  it('hides the second operand for square root and omits it from the request', async () => {
    calculateMock.mockResolvedValue({ operation: 'sqrt', a: 81, result: 9 });
    render(<Calculator />);

    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'Square root' }));

    expect(screen.getByLabelText('Second value (b)')).not.toBeVisible();

    await user.type(screen.getByLabelText('First value (a)'), '81');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(calculateMock).toHaveBeenCalledWith('sqrt', 81, undefined);
    expect(await screen.findByText('9')).toBeInTheDocument();
  });

  it('rejects an empty operand without calling the API', async () => {
    render(<Calculator />);

    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('First value is required.');
    expect(calculateMock).not.toHaveBeenCalled();
  });

  it('rejects a non-numeric operand without calling the API', async () => {
    render(<Calculator />);

    const user = await fillOperands('12', 'abc');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Second value must be a valid number.',
    );
    expect(calculateMock).not.toHaveBeenCalled();
  });

  it('shows the error message returned by the API', async () => {
    calculateMock.mockRejectedValue(new ApiError('division by zero is not allowed', 422));
    render(<Calculator />);

    const user = await fillOperands('10', '0');
    await user.click(screen.getByRole('button', { name: 'Divide' }));
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'division by zero is not allowed',
    );
  });

  it('falls back to a generic message for unexpected failures', async () => {
    calculateMock.mockRejectedValue(new Error('boom'));
    render(<Calculator />);

    const user = await fillOperands('1', '2');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Something went wrong');
  });

  it('clears operands and the previous result', async () => {
    calculateMock.mockResolvedValue({ operation: 'add', a: 4, b: 2, result: 6 });
    render(<Calculator />);

    const user = await fillOperands('4', '2');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));
    expect(await screen.findByText('6')).toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: 'Clear' }));

    expect(screen.getByLabelText('First value (a)')).toHaveValue('');
    expect(screen.getByLabelText('Second value (b)')).toHaveValue('');
    expect(screen.queryByText('6')).not.toBeInTheDocument();
  });

  it('drops a stale result when the operation changes', async () => {
    calculateMock.mockResolvedValue({ operation: 'add', a: 4, b: 2, result: 6 });
    render(<Calculator />);

    const user = await fillOperands('4', '2');
    await user.click(screen.getByRole('button', { name: 'Calculate' }));
    expect(await screen.findByText('6')).toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: 'Subtract' }));

    expect(screen.queryByText('6')).not.toBeInTheDocument();
  });
});
