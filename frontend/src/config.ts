import type { Location, Role } from "./types";

export const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8088";
export const WS_BASE = API_BASE.replace(/^http/, "ws");

export const initialPickup: Location = { latitude: 41.015137, longitude: 28.97953 };
export const initialDropoff: Location = { latitude: 41.043, longitude: 29.009 };
export const driverStart: Location = { latitude: 41.016, longitude: 28.981 };

export function sampleAccount(role: Role) {
  return {
    email: `${role}.${Date.now()}@example.com`,
    password: "Password12345",
    display_name: role === "rider" ? "Eren Rider" : "Eren Driver",
    phone_number: role === "rider" ? "+905551000001" : "+905551000002",
    role,
  };
}
