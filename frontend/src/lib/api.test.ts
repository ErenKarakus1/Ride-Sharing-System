import { afterEach, describe, expect, it, vi } from "vitest";
import { createApi, request } from "./api";

const fetchMock = vi.fn();

globalThis.fetch = fetchMock;

describe("api client", () => {
  afterEach(() => {
    fetchMock.mockReset();
  });

  it("sends JSON requests through the configured API base URL", async () => {
    fetchMock.mockResolvedValueOnce(response({ id: "ride-1" }));

    const result = await request<{ id: string }>("/api/v1/rides", {
      method: "POST",
      body: { pickup: { lat: 41, lng: 29 } },
    });

    expect(result).toEqual({ id: "ride-1" });
    expect(fetchMock).toHaveBeenCalledWith("http://localhost:8088/api/v1/rides", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ pickup: { lat: 41, lng: 29 } }),
    });
  });

  it("adds bearer tokens for authenticated calls", async () => {
    fetchMock.mockResolvedValueOnce(response({ status: "accepted" }));

    await createApi("token-123").post("/api/v1/rides/ride-1/accept", {});

    expect(fetchMock).toHaveBeenCalledWith(
      "http://localhost:8088/api/v1/rides/ride-1/accept",
      expect.objectContaining({
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer token-123",
        },
      }),
    );
  });

  it("returns null for empty successful responses", async () => {
    fetchMock.mockResolvedValueOnce(response(null, { text: "" }));

    await expect(createApi("token-123").put("/api/v1/drivers/driver-1/location", {})).resolves.toBeNull();
  });

  it("throws backend error messages", async () => {
    fetchMock.mockResolvedValueOnce(response({ error: "ride is already completed" }, { ok: false, status: 409 }));

    await expect(createApi("token-123").post("/api/v1/rides/ride-1/start", {})).rejects.toThrow(
      "ride is already completed",
    );
  });
});

function response(payload: unknown, options: { ok?: boolean; status?: number; text?: string } = {}) {
  return {
    ok: options.ok ?? true,
    status: options.status ?? 200,
    text: () => Promise.resolve(options.text ?? JSON.stringify(payload)),
  };
}
