import type { Operation } from '../operations';

export interface CalculateResponse {
  operation: Operation;
  a: number;
  b?: number;
  result: number;
}

/** Thrown for both transport failures and error responses from the API. */
export class ApiError extends Error {
  readonly status?: number;

  constructor(message: string, status?: number) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

function baseUrl(): string {
  return import.meta.env.VITE_API_URL ?? 'http://localhost:8080';
}

/**
 * Calls POST /api/calculate. `b` is omitted for unary operations, which the
 * backend requires to be absent rather than zero.
 */
export async function calculate(
  operation: Operation,
  a: number,
  b?: number,
): Promise<CalculateResponse> {
  const body: Record<string, unknown> = { operation, a };
  if (b !== undefined) {
    body.b = b;
  }

  let response: Response;
  try {
    response = await fetch(`${baseUrl()}/api/calculate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
  } catch {
    throw new ApiError('Cannot reach the calculator service. Is the backend running?');
  }

  const payload = await readJson(response);

  if (!response.ok) {
    // The backend reports failures as {"error": "..."}.
    const reported = payload && typeof payload.error === 'string' ? payload.error : null;
    throw new ApiError(
      reported ?? `Request failed with status ${response.status}`,
      response.status,
    );
  }

  if (!payload || typeof payload.result !== 'number') {
    throw new ApiError('Received an unexpected response from the service.');
  }

  return payload as unknown as CalculateResponse;
}

async function readJson(response: Response): Promise<Record<string, unknown> | null> {
  try {
    return (await response.json()) as Record<string, unknown>;
  } catch {
    return null;
  }
}
