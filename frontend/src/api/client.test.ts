import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ApiError, calculate } from './client';

function mockResponse(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

const fetchMock = vi.fn();

beforeEach(() => {
  fetchMock.mockReset();
  vi.stubGlobal('fetch', fetchMock);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('calculate', () => {
  it('posts the operation and both operands', async () => {
    fetchMock.mockResolvedValue(
      mockResponse({ operation: 'add', a: 4, b: 2, result: 6 }),
    );

    const response = await calculate('add', 4, 2);

    expect(response.result).toBe(6);
    expect(fetchMock).toHaveBeenCalledTimes(1);

    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe('http://localhost:8080/api/calculate');
    expect(init.method).toBe('POST');
    expect(JSON.parse(init.body)).toEqual({ operation: 'add', a: 4, b: 2 });
  });

  it('omits b for unary operations', async () => {
    fetchMock.mockResolvedValue(mockResponse({ operation: 'sqrt', a: 81, result: 9 }));

    await calculate('sqrt', 81);

    const [, init] = fetchMock.mock.calls[0];
    expect(JSON.parse(init.body)).toEqual({ operation: 'sqrt', a: 81 });
  });

  it('raises the error message returned by the API', async () => {
    fetchMock.mockResolvedValue(
      mockResponse({ error: 'division by zero is not allowed' }, 422),
    );

    await expect(calculate('divide', 10, 0)).rejects.toThrowError(
      new ApiError('division by zero is not allowed', 422),
    );
  });

  it('falls back to the status code when the body has no error field', async () => {
    fetchMock.mockResolvedValue(mockResponse({}, 500));

    await expect(calculate('add', 1, 2)).rejects.toThrowError(
      /Request failed with status 500/,
    );
  });

  it('reports an unreachable backend', async () => {
    fetchMock.mockRejectedValue(new TypeError('Failed to fetch'));

    await expect(calculate('add', 1, 2)).rejects.toThrowError(/Cannot reach/);
  });

  it('rejects a success response without a numeric result', async () => {
    fetchMock.mockResolvedValue(mockResponse({ operation: 'add', a: 1, b: 2 }));

    await expect(calculate('add', 1, 2)).rejects.toThrowError(/unexpected response/);
  });
});
