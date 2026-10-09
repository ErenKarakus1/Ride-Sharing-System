import { describe, expect, it } from "vitest";
import { capitalize, formatCoord } from "./format";

describe("format helpers", () => {
  it("formats coordinates for compact display", () => {
    expect(formatCoord({ latitude: 41.015137, longitude: 28.97953 })).toBe("41.015, 28.980");
  });

  it("capitalizes role labels", () => {
    expect(capitalize("driver")).toBe("Driver");
  });
});
