import { useMemo, useState } from "react";
import { CircleDollarSign, Clock3, MapPin, Navigation } from "lucide-react";
import { Metric } from "./components/Metric";
import { Sidebar } from "./components/Sidebar";
import { DriverPanel } from "./features/DriverPanel";
import { LiveStatePanel } from "./features/LiveStatePanel";
import { NotificationsPanel } from "./features/NotificationsPanel";
import { RiderPanel } from "./features/RiderPanel";
import { useLocalStorage } from "./hooks/useLocalStorage";
import { useNotifications } from "./hooks/useNotifications";
import { createApi, request } from "./lib/api";
import { capitalize, formatCoord } from "./lib/format";
import { API_BASE, driverStart, initialDropoff, initialPickup, sampleAccount } from "./config";
import type { AccountForm, FareEstimate, Location, Match, Payment, Ride, Role, Session } from "./types";

export function App() {
  const [role, setRole] = useState<Role>("rider");
  const [session, setSession] = useLocalStorage<Session | null>("ride-sharing-session", null);
  const [riderForm, setRiderForm] = useState<AccountForm>(sampleAccount("rider"));
  const [driverForm, setDriverForm] = useState<AccountForm>(sampleAccount("driver"));
  const [pickup, setPickup] = useState<Location>(initialPickup);
  const [dropoff, setDropoff] = useState<Location>(initialDropoff);
  const [driverLocation, setDriverLocation] = useState<Location>(driverStart);
  const [fare, setFare] = useState<FareEstimate | null>(null);
  const [ride, setRide] = useState<Ride | null>(null);
  const [match, setMatch] = useState<Match | null>(null);
  const [payment, setPayment] = useState<Payment | null>(null);
  const [status, setStatus] = useState("Ready");
  const [error, setError] = useState("");

  const activeForm = role === "rider" ? riderForm : driverForm;
  const setActiveForm = role === "rider" ? setRiderForm : setDriverForm;
  const token = session?.access_token;
  const api = useMemo(() => createApi(token), [token]);
  const notifications = useNotifications(token);

  const isRider = session?.role === "rider";
  const isDriver = session?.role === "driver";

  async function run<T>(action: () => Promise<T>, successMessage: string) {
    setError("");
    setStatus("Working...");

    try {
      const result = await action();
      setStatus(successMessage);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unexpected error");
      setStatus("Needs attention");
      return null;
    }
  }

  async function register() {
    await run(async () => {
      const response = await request<Session>("/api/v1/auth/register", {
        method: "POST",
        body: activeForm,
      });
      setSession(response);
      return response;
    }, `${capitalize(role)} registered`);
  }

  async function login() {
    await run(async () => {
      const response = await request<Session>("/api/v1/auth/login", {
        method: "POST",
        body: { email: activeForm.email, password: activeForm.password },
      });
      setSession(response);
      return response;
    }, `${capitalize(role)} signed in`);
  }

  async function estimateFare() {
    const response = await run(
      () => api.post<FareEstimate>("/api/v1/fare-estimates", { pickup, dropoff }),
      "Fare estimated",
    );
    if (response) setFare(response);
  }

  async function createRide() {
    const response = await run(() => api.post<Ride>("/api/v1/rides", { pickup, dropoff }), "Ride requested");
    if (response) {
      setRide(response);
      setMatch(null);
      setPayment(null);
    }
  }

  async function authorizePayment() {
    if (!ride || !fare) {
      setError("Create a ride and estimate fare first.");
      return;
    }

    const response = await run(
      () =>
        api.post<Payment>("/api/v1/payments/authorize", {
          ride_id: ride.id,
          amount: fare.amount,
          currency: fare.currency,
        }),
      "Payment authorized",
    );
    if (response) setPayment(response);
  }

  async function refreshRide() {
    if (!ride) return;
    const response = await run(() => api.get<Ride>(`/api/v1/rides/${ride.id}`), "Ride refreshed");
    if (response) setRide(response);
  }

  async function refreshPayment() {
    if (!payment) return;
    const response = await run(() => api.get<Payment>(`/api/v1/payments/${payment.id}`), "Payment refreshed");
    if (response) setPayment(response);
  }

  async function updateDriverLocation() {
    await run(
      () => api.put<void>(`/api/v1/drivers/${session?.user_id}/location`, driverLocation),
      "Driver location updated",
    );
  }

  async function setDriverAvailable() {
    await run(() => api.post<void>(`/api/v1/drivers/${session?.user_id}/available`, {}), "Driver is available");
  }

  async function findDriver() {
    if (!ride) {
      setError("Create a ride first.");
      return;
    }

    const response = await run(
      () =>
        api.post<Match>("/api/v1/matches", {
          ride_id: ride.id,
          pickup,
          radius_km: 5,
          limit: 5,
        }),
      "Driver matched",
    );
    if (response) setMatch(response);
  }

  async function acceptRide() {
    if (!ride) return;
    const response = await run(() => api.post<Ride>(`/api/v1/rides/${ride.id}/accept`, {}), "Ride accepted");
    if (response) setRide(response);
  }

  async function startRide() {
    if (!ride) return;
    const response = await run(() => api.post<Ride>(`/api/v1/rides/${ride.id}/start`, {}), "Ride started");
    if (response) setRide(response);
  }

  async function completeRide() {
    if (!ride) return;
    const response = await run(() => api.post<Ride>(`/api/v1/rides/${ride.id}/complete`, {}), "Ride completed");
    if (response) setRide(response);
  }

  return (
    <main className="app-shell">
      <Sidebar
        role={role}
        activeForm={activeForm}
        session={session}
        onRoleChange={setRole}
        onFormChange={setActiveForm}
        onRegister={register}
        onLogin={login}
        onSignOut={() => setSession(null)}
      />

      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">Gateway {API_BASE}</p>
            <h1>Trip operations</h1>
          </div>
          <div className={`system-state ${error ? "error" : "ok"}`}>
            <span>{error || status}</span>
          </div>
        </header>

        <section className="metrics">
          <Metric icon={MapPin} label="Pickup" value={formatCoord(pickup)} />
          <Metric icon={Navigation} label="Dropoff" value={formatCoord(dropoff)} />
          <Metric icon={Clock3} label="Ride status" value={ride?.status ?? "No ride"} />
          <Metric icon={CircleDollarSign} label="Payment" value={payment?.status ?? "No payment"} />
        </section>

        <section className="flow-grid">
          <RiderPanel
            isRider={isRider}
            pickup={pickup}
            dropoff={dropoff}
            fare={fare}
            ride={ride}
            match={match}
            onPickupChange={setPickup}
            onDropoffChange={setDropoff}
            onEstimateFare={estimateFare}
            onCreateRide={createRide}
            onAuthorizePayment={authorizePayment}
            onRefreshRide={refreshRide}
          />
          <DriverPanel
            isDriver={isDriver}
            driverLocation={driverLocation}
            ride={ride}
            match={match}
            payment={payment}
            onDriverLocationChange={setDriverLocation}
            onUpdateDriverLocation={updateDriverLocation}
            onSetDriverAvailable={setDriverAvailable}
            onFindDriver={findDriver}
            onAcceptRide={acceptRide}
            onStartRide={startRide}
            onCompleteRide={completeRide}
          />
          <LiveStatePanel ride={ride} payment={payment} onRefreshPayment={refreshPayment} />
          <NotificationsPanel notifications={notifications} />
        </section>
      </section>
    </main>
  );
}
