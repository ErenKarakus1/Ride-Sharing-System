import { describe, expect, it } from "vitest";
import { pageFromPath, pathForPage } from "./routes";

describe("routes", () => {
  it("maps known paths to pages", () => {
    expect(pageFromPath("/driver")).toBe("driver");
    expect(pageFromPath("/notifications/extra")).toBe("notifications");
  });

  it("falls back unknown paths to auth", () => {
    expect(pageFromPath("/")).toBe("auth");
    expect(pageFromPath("/unknown")).toBe("auth");
  });

  it("builds page paths", () => {
    expect(pathForPage("rider")).toBe("/rider");
  });
});
