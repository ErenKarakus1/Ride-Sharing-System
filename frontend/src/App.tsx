import { useEffect, useMemo, useState } from "react";
import { Sidebar } from "./components/Sidebar";
import { useLocalStorage } from "./hooks/useLocalStorage";
import { useNotifications } from "./hooks/useNotifications";
import { ApiError, createApi, request } from "./lib/api";
import { initialActionStates } from "./lib/actionState";
import { capitalize } from "./lib/format";
import { pageFromPath, pathForPage } from "./lib/routes";
import { AuthPage } from "./pages/AuthPage";
import { DriverPage } from "./pages/DriverPage";
import { NotificationsPage } from "./pages/NotificationsPage";
import { RideDetailsPage } from "./pages/RideDetailsPage";
import { RiderPage } from "./pages/RiderPage";
import { API_BASE, driverStart, initialDropoff, initialPickup, sampleAccount } from "./config";
import type { AccountForm, ActionKey, FareEstimate, Location, Match, Payment, Ride, Role, Session, Sessions } from "./types";
import type { Page } from "./lib/routes";

export function App() {
  const [role, setRole] = useState<Role>("rider");
  const [page, setPage] = useState<Page>(() => pageFromPath(window.location.pathname));
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

  useEffect(() => {
    const onPopState = () => setPage(pageFromPath(window.location.pathname));
    window.addEventListener("popstate", onPopState);
    return () => window.removeEventListener("popstate", onPopState);
  }, []);

  useEffect(() => {
    const session = sessions[role];
    if (!session && page !== "auth") {
      navigate("auth");
      return;
    }
    if (session && page === "auth") {
      navigate(role);
    }
  }, [page, role, sessions]);

  useEffect(() => {
    if (!ride || ride.driver_id || payment?.status !== "authorized") return;

    let stopped = false;
    const intervalID = window.setInterval(async () => {
      try {
        const updated = await riderApi.get<Ride>(`/api/v1/rides/${ride.id}`);
        if (stopped) return;
        setRide(updated);
        if (isRideMatched(updated)) {
          setActions((current) => ({
            ...current,
            matching: { loading: false, message: "Driver matched", error: "" },
          }));
          window.clearInterval(intervalID);
        }
      } catch {
        window.clearInterval(intervalID);
      }
    }, 2000);

    return () => {
      stopped = true;
      window.clearInterval(intervalID);
    };
  }, [payment?.status, ride?.driver_id, ride?.id, riderApi]);

  function navigate(nextPage: Page) {
    setPage(nextPage);
    window.history.pushState({}, "", pathForPage(nextPage));
  }

  function changeRole(nextRole: Role) {
    setRole(nextRole);
    navigate(sessions[nextRole] ? nextRole : "auth");
  }

  function signOut() {
    setSessions((current) => ({ ...current, [role]: null }));
    navigate("auth");
  }

  async function run<T>(key: ActionKey, action: () => Promise<T>, successMessage: string, sessionRole?: Role) {
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
      const expiredSession = sessionRole && err instanceof ApiError && err.status === 401;
      const message = expiredSession ? "Session expired. Sign in again." : errorMessage(err);
      if (expiredSession) {
        setSessions((current) => ({ ...current, [sessionRole]: null }));
      }

      setActions((current) => ({
        ...current,
        [key]: {
          loading: false,
          message: "",
          error: message,
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
      navigate(response.role);
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
      navigate(response.role);
      return response;
    }, `${capitalize(role)} signed in`);
  }

  async function estimateFare() {
    const response = await run(
      "fare",
      () => riderApi.post<FareEstimate>("/api/v1/fare-estimates", { pickup, dropoff }),
      "Fare estimated",
      "rider",
    );
    if (response) setFare(response);
  }

  async function createRide() {
    const response = await run(
      "ride",
      () => riderApi.post<Ride>("/api/v1/rides", { pickup, dropoff }),
      "Ride requested",
      "rider",
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
      "rider",
    );
    if (response) {
      setPayment(response);
      await refreshRideUntilMatched(ride.id);
    }
  }

  async function refreshRideUntilMatched(rideID: string) {
    setActions((current) => ({
      ...current,
      matching: { loading: true, message: "", error: "" },
    }));

    for (let attempt = 0; attempt < 20; attempt += 1) {
      await sleep(500);
      try {
        const updated = await riderApi.get<Ride>(`/api/v1/rides/${rideID}`);
        setRide(updated);
        if (isRideMatched(updated)) {
          setActions((current) => ({
            ...current,
            matching: { loading: false, message: "Driver matched", error: "" },
          }));
          return;
        }
      } catch {
        return;
      }
    }

    setActions((current) => ({
      ...current,
      matching: { loading: false, message: "Still waiting for driver match", error: "" },
    }));
  }

  async function refreshRide() {
    if (!ride) return;
    const response = await run("ride", () => riderApi.get<Ride>(`/api/v1/rides/${ride.id}`), "Ride refreshed", "rider");
    if (response) setRide(response);
  }

  async function refreshPayment() {
    if (!payment) return;
    const response = await run(
      "payment",
      () => riderApi.get<Payment>(`/api/v1/payments/${payment.id}`),
      "Payment refreshed",
      "rider",
    );
    if (response) setPayment(response);
  }

  async function updateDriverLocation() {
    await run(
      "driver-location",
      () => driverApi.put<void>(`/api/v1/drivers/${sessions.driver?.user_id}/location`, driverLocation),
      "Driver location updated",
      "driver",
    );
  }

  async function setDriverAvailable() {
    await run(
      "driver-location",
      () => driverApi.post<void>(`/api/v1/drivers/${sessions.driver?.user_id}/available`, {}),
      "Driver is available",
      "driver",
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
      "rider",
    );
    if (response) setMatch(response);
  }

  async function acceptRide() {
    if (!ride) return;
    const response = await run(
      "lifecycle",
      () => driverApi.post<Ride>(`/api/v1/rides/${ride.id}/accept`, {}),
      "Ride accepted",
      "driver",
    );
    if (response) setRide(response);
  }

  async function startRide() {
    if (!ride) return;
    const response = await run(
      "lifecycle",
      () => driverApi.post<Ride>(`/api/v1/rides/${ride.id}/start`, {}),
      "Ride started",
      "driver",
    );
    if (response) setRide(response);
  }

  async function completeRide() {
    if (!ride) return;
    const response = await run(
      "lifecycle",
      () => driverApi.post<Ride>(`/api/v1/rides/${ride.id}/complete`, {}),
      "Ride completed",
      "driver",
    );
    if (response) setRide(response);
  }

  function resetTrip() {
    setFare(null);
    setRide(null);
    setMatch(null);
    setPayment(null);
    setActions(initialActionStates);
  }

  return (
    <main className="app-shell">
      <Sidebar
        role={role}
        page={page}
        activeForm={activeForm}
        sessions={sessions}
        authLoading={actions.auth.loading}
        onRoleChange={changeRole}
        onPageChange={navigate}
        onFormChange={setActiveForm}
        onRegister={register}
        onLogin={login}
        onSignOut={signOut}
        onResetTrip={resetTrip}
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

        {page === "auth" && (
          <AuthPage
            sessions={sessions}
            action={actions.auth}
            fare={fare}
            ride={ride}
            match={match}
            payment={payment}
          />
        )}
        {page === "rider" && (
          <RiderPage
            isRider={isRider}
            action={actions.ride.error || actions.ride.message || actions.ride.loading ? actions.ride : actions.fare}
            matchingAction={actions.matching}
            pickup={pickup}
            dropoff={dropoff}
            fare={fare}
            ride={ride}
            match={match}
            payment={payment}
            onPickupChange={setPickup}
            onDropoffChange={setDropoff}
            onEstimateFare={estimateFare}
            onCreateRide={createRide}
            onAuthorizePayment={authorizePayment}
            onFindDriver={findDriver}
            onRefreshRide={refreshRide}
          />
        )}
        {page === "driver" && (
          <DriverPage
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
        )}
        {page === "ride" && (
          <RideDetailsPage
            ride={ride}
            payment={payment}
            paymentAction={actions.payment}
            onRefreshPayment={refreshPayment}
          />
        )}
        {page === "notifications" && <NotificationsPage notifications={notifications} />}
      </section>
    </main>
  );
}

function errorMessage(err: unknown) {
  if (err instanceof Error) return err.message;
  return "Unexpected error";
}

function sleep(ms: number) {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

function isRideMatched(ride: Ride) {
  return Boolean(ride.driver_id || ride.status === "accepted" || ride.status === "started" || ride.status === "completed");
}
