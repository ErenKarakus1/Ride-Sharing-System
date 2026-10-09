import type { Location } from "../types";

export function formatCoord(location: Location) {
  return `${location.latitude.toFixed(3)}, ${location.longitude.toFixed(3)}`;
}

export function capitalize(value: string) {
  return `${value.charAt(0).toUpperCase()}${value.slice(1)}`;
}
