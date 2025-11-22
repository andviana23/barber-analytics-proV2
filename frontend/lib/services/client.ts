import { z } from 'zod';

const DEFAULT_RETRIES = 3;
const RETRY_BASE_MS = 100;

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

export type FetcherConfig = {
  baseUrl?: string;
  getAuthToken?: () => string | undefined;
  defaultHeaders?: Record<string, string>;
  retries?: number;
};

export type HttpError = {
  status: number;
  code: string;
  message: string;
  details?: Record<string, unknown>;
};

export class ApiClient {
  private readonly baseUrl: string;
  private readonly getAuthToken?: () => string | undefined;
  private readonly defaultHeaders: Record<string, string>;
  private readonly retries: number;

  constructor(config?: FetcherConfig) {
    this.baseUrl = config?.baseUrl ?? process.env.NEXT_PUBLIC_API_URL ?? '/api/v1';
    this.getAuthToken = config?.getAuthToken;
    this.defaultHeaders = config?.defaultHeaders ?? { 'Content-Type': 'application/json' };
    this.retries = config?.retries ?? DEFAULT_RETRIES;
  }

  async request<T>({
    method,
    path,
    body,
    schema,
    signal,
  }: {
    method: HttpMethod;
    path: string;
    body?: unknown;
    schema?: z.ZodType<T>;
    signal?: AbortSignal;
  }): Promise<T> {
    const headers: Record<string, string> = { ...this.defaultHeaders };
    const token = this.getAuthToken?.();
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    const url = `${this.baseUrl}${path}`;
    const fetcher = async () =>
      fetch(url, {
        method,
        headers,
        body: body ? JSON.stringify(body) : undefined,
        signal,
      });

    const response = await this.retry(fetcher, this.retries);
    const json = await this.safeJson(response);

    if (!response.ok) {
      const err = this.toHttpError(response.status, json);
      throw err;
    }

    if (!schema) {
      return json as T;
    }

    return schema.parse(json);
  }

  private async retry(fetcher: () => Promise<Response>, retries: number): Promise<Response> {
    let attempt = 0;
    let lastError: unknown;

    while (attempt <= retries) {
      try {
        const res = await fetcher();
        if (res.status >= 500 && attempt < retries) {
          attempt += 1;
          await this.delay(RETRY_BASE_MS * 2 ** (attempt - 1));
          continue;
        }
        return res;
      } catch (err) {
        lastError = err;
        if (attempt >= retries) {
          throw err;
        }
        attempt += 1;
        await this.delay(RETRY_BASE_MS * 2 ** (attempt - 1));
      }
    }

    // Should never reach here
    throw lastError ?? new Error('Erro desconhecido ao executar requisição');
  }

  private async safeJson(response: Response): Promise<unknown> {
    try {
      return await response.json();
    } catch {
      return {};
    }
  }

  private toHttpError(status: number, payload: unknown): HttpError {
    const code = (payload as { error?: string })?.error ?? 'unknown_error';
    const message = (payload as { message?: string })?.message ?? 'Erro inesperado';
    const details = (payload as { details?: Record<string, unknown> })?.details;
    return { status, code, message, details };
  }

  private delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}

export const apiClient = new ApiClient();
