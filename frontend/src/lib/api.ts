import { API_BASE } from "../config";

type RequestOptions = {
  method?: "GET" | "POST" | "PUT";
  token?: string;
  body?: unknown;
};

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
  const payload = text ? JSON.parse(text) : null;

  if (!response.ok) {
    throw new Error(payload?.error ?? `Request failed with ${response.status}`);
  }

  return payload as T;
}
