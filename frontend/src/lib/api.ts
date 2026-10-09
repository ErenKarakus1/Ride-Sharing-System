import { API_BASE } from "../config";

type RequestOptions = {
  method?: "GET" | "POST" | "PUT";
  token?: string;
  body?: unknown;
};

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export function createApi(token?: string) {
  return {
    get: <T>(path: string) => request<T>(path, { token }),
    post: <T>(path: string, body: unknown) => request<T>(path, { method: "POST", token, body }),
    put: <T>(path: string, body: unknown) => request<T>(path, { method: "PUT", token, body }),
  };
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    method: options.method ?? "GET",
    headers: {
      "Content-Type": "application/json",
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {}),
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const text = await response.text();
  const payload = parsePayload(text);

  if (!response.ok) {
    throw new ApiError(errorMessage(payload, text, response.status), response.status);
  }

  return payload as T;
}

function parsePayload(text: string) {
  if (!text) return null;

  try {
    return JSON.parse(text) as unknown;
  } catch {
    return null;
  }
}

function errorMessage(payload: unknown, text: string, status: number) {
  if (payload && typeof payload === "object" && "error" in payload) {
    const error = (payload as { error?: unknown }).error;
    if (typeof error === "string" && error.trim()) return error;
  }

  if (text.trim()) return text;
  return `Request failed with ${status}`;
}
