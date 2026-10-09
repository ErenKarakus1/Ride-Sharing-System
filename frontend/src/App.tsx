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
import { initialActionStates } from "./lib/actionState";
import { capitalize, formatCoord } from "./lib/format";
import { API_BASE, driverStart, initialDropoff, initialPickup, sampleAccount } from "./config";
import type { AccountForm, ActionKey, FareEstimate, Location, Match, Payment, Ride, Role, Session, Sessions } from "./types";

export function App() {
  const [role, setRole] = useState<Role>("rider");
  const [sessions, setSessions] = useLocalStorage<Sessions>("ride-sharing-sessions", {
    rider: null,
    driver: null,
  });
  const [riderForm, setRiderForm] = useState<AccountForm>(sampleAccount("rider"));
  const [driverForm, setDriverForm] = useState<AccountForm>(sampleAccount("driver"));
  const [pickup, setPickup] = useState<Location>(initialPickup);
  const [dropoff, setDropoff] = useState<Location>(initialDropoff);
  const [driverLocation, setDriverLocation] = useState<Location>(driverStart);
  const [fare, setFare] = useState<FareEstimate | null>(null);
  const [ride, setRide] = useState<Ride | null>(null);
  const [match, setMatch] = useState<Match | null>(null);
  const [payment, setPayment] = useState<Payment | null>(null);
  const [actions, setActions] = useState(initialActionStates);

  const activeForm = role === "rider" ? riderForm : driverForm;
  const setActiveForm = role === "rider" ? setRiderForm : setDriverForm;
  const activeSession = sessions[role];
  const riderApi = useMemo(() => createApi(sessions.rider?.access_token), [sessions.rider?.access_token]);
  const driverApi = useMemo(() => createApi(sessions.driver?.access_token), [sessions.driver?.access_token]);
  const notifications = useNotifications(sessions.rider?.access_token);

  const isRider = Boolean(sessions.rider);
  const isDriver = Boolean(sessions.driver);

  async function run<T>(key: ActionKey, action: () => Promise<T>, successMessage: string) {
    setActions((current) => ({
      ...current,
      [key]: { loading: true, message: "", error: "" },
    }));

    try {
      const result = await action();
      setActions((current) => ({
        ...current,
        [key]: { loading: false, message: successMessage, error: "" },
      }));
      return result;
    } catch (err) {
      setActions((current) => ({
        ...current,
        [key]: {
          loading: false,
          message: "",
          error: err instanceof Error ? err.message : "Unexpected error",
        },
      }));
      return null;
    }
  }

  async function register() {
    await run("auth", async () => {
      const response = await request<Session>("/api/v1/auth/register", {
        method: "POST",
        body: activeForm,
      });
      setSessions((current) => ({ ...current, [response.role]: response }));
      return response;
    }, `${capitalize(role)} registered`);
  }

  async function login() {
    await run("auth", async () => {
      const response = await request<Session>("/api/v1/auth/login", {
        method: "POST",
        body: { email: activeForm.email, password: activeForm.password },
      });
      setSessions((current) => ({ ...current, [response.role]: response }));
      return response;
    }, `${capitalize(role)} signed in`);
  }

  async function estimateFare() {
    const response = await run(
      "fare",
      () => riderApi.post<FareEstimate>("/api/v1/fare-estimates", { pickup, dropoff }),
      "Fare estimated",
    );
    if (response) setFare(response);
  }

  async function createRide() {
    const response = await run(
      "ride",
      () => riderApi.post<Ride>("/api/v1/rides", { pickup, dropoff }),
      "Ride requested",
    );
    if (response) {
      setRide(response);
      setMatch(null);
      setPayment(null);
    }
  }

  async function authorizePayment() {
    if (!ride || !fare) {
      setActions((current) => ({
        ...current,
        payment: { loading: false, message: "", error: "Create a ride and estimate fare first." },
      }));
      return;
    }

    const response = await run(
      "payment",
      () =>
        riderApi.post<Payment>("/api/v1/payments/authorize", {
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
    const response = await run("ride", () => riderApi.get<Ride>(`/api/v1/rides/${ride.id}`), "Ride refreshed");
    if (response) setRide(response);
  }

  async function refreshPayment() {
    if (!payment) return;
    const response = await run(
      "payment",
      () => riderApi.get<Payment>(`/api/v1/payments/${payment.id}`),
      "Payment refreshed",
    );
    if (response) setPayment(response);
  }

  async function updateDriverLocation() {
    await run(
      "driver-location",
      () => driverApi.put<void>(`/api/v1/drivers/${sessions.driver?.user_id}/location`, driverLocation),
      "Driver location updated",
    );
  }

  async function setDriverAvailable() {
    await run(
      "driver-location",
      () => driverApi.post<void>(`/api/v1/drivers/${sessions.driver?.user_id}/available`, {}),
      "Driver is available",
    );
  }

  async function findDriver() {
    if (!ride) {
      setActions((current) => ({
        ...current,
        matching: { loading: false, message: "", error: "Create a ride first." },
      }));
      return;
    }

    const response = await run(
      "matching",
      () =>
        riderApi.post<Match>("/api/v1/matches", {
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
    const response = await run(
      "lifecycle",
      () => driverApi.post<Ride>(`/api/v1/rides/${ride.id}/accept`, {}),
      "Ride accepted",
    );
    if (response) setRide(response);
  }

  async function startRide() {
    if (!ride) return;
    const response = await run(
      "lifecycle",
      () => driverApi.post<Ride>(`/api/v1/rides/${ride.id}/start`, {}),
      "Ride started",
    );
    if (response) setRide(response);
  }

  async function completeRide() {
    if (!ride) return;
    const response = await run(
      "lifecycle",
      () => driverApi.post<Ride>(`/api/v1/rides/${ride.id}/complete`, {}),
      "Ride completed",
    );
    if (response) setRide(response);
  }

  return (
    <main className="app-shell">
      <Sidebar
        role={role}
        activeForm={activeForm}
        sessions={sessions}
        onRoleChange={setRole}
        onFormChange={setActiveForm}
        onRegister={register}
        onLogin={login}
        onSignOut={() => setSessions((current) => ({ ...current, [role]: null }))}
      />

      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">Gateway {API_BASE}</p>
            <h1>Trip operations</h1>
          </div>
          <div className={`system-state ${actions.auth.error ? "error" : "ok"}`}>
            <span>{actions.auth.loading ? "Signing in..." : actions.auth.error || actions.auth.message || activeSession?.email || "Ready"}</span>
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
            action={actions.ride.error || actions.ride.message || actions.ride.loading ? actions.ride : actions.fare}
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
            locationAction={actions["driver-location"]}
            matchingAction={actions.matching}
            lifecycleAction={actions.lifecycle}
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
          <LiveStatePanel ride={ride} payment={payment} paymentAction={actions.payment} onRefreshPayment={refreshPayment} />
          <NotificationsPanel notifications={notifications} />
        </section>
      </section>
    </main>
  );
}
