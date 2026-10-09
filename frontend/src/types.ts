export type Role = "rider" | "driver";

export type Location = {
  latitude: number;
  longitude: number;
};

export type LocationPreset = {
  label: string;
  location: Location;
};

export type AccountForm = {
  email: string;
  password: string;
  display_name: string;
  phone_number: string;
  role: Role;
};

export type Session = {
  access_token: string;
  user_id: string;
  email: string;
  role: Role;
};

export type Sessions = Record<Role, Session | null>;

export type ActionKey =
  | "auth"
  | "fare"
  | "ride"
  | "payment"
  | "driver-location"
  | "matching"
  | "lifecycle";

export type ActionState = {
  loading: boolean;
  message: string;
  error: string;
};

export type FareEstimate = {
  distance_km: number;
  duration_minutes: number;
  currency: string;
  amount: number;
};

export type Ride = {
  id: string;
  rider_id: string;
  driver_id?: string;
  pickup: Location;
  dropoff: Location;
  status: string;
  created_at?: string;
  updated_at?: string;
};

export type Match = {
  ride_id: string;
  driver_id: string;
  latitude: number;
  longitude: number;
};

export type Payment = {
  id: string;
  ride_id: string;
  rider_id: string;
  driver_id?: string;
  amount: number;
  currency: string;
  status: string;
  created_at?: string;
  updated_at?: string;
};

export type NotificationItem = {
  message: string;
  received_at: string;
};
