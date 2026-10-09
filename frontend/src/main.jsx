import React, { useEffect, useMemo, useRef, useState } from "react";
import { createRoot } from "react-dom/client";
import {
  Activity,
  Bell,
  Car,
  CheckCircle2,
  CircleDollarSign,
  Clock3,
  Crosshair,
  LogIn,
  LogOut,
  MapPin,
  Navigation,
  Radio,
  RefreshCcw,
  Route,
  ShieldCheck,
  UserPlus,
} from "lucide-react";
import "./styles.css";

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8088";
const WS_BASE = API_BASE.replace(/^http/, "ws");
const PASSWORD = "Password12345";

const initialPickup = { latitude: 41.015137, longitude: 28.97953 };
const initialDropoff = { latitude: 41.043, longitude: 29.009 };
const driverStart = { latitude: 41.016, longitude: 28.981 };

function App() {
  const [role, setRole] = useState("rider");
  const [session, setSession] = useLocalStorage("ride-sharing-session", null);
  const [riderForm, setRiderForm] = useState(sampleAccount("rider"));
  const [driverForm, setDriverForm] = useState(sampleAccount("driver"));
  const [pickup, setPickup] = useState(initialPickup);
  const [dropoff, setDropoff] = useState(initialDropoff);
  const [driverLocation, setDriverLocation] = useState(driverStart);
  const [fare, setFare] = useState(null);
  const [ride, setRide] = useState(null);
  const [match, setMatch] = useState(null);
  const [payment, setPayment] = useState(null);
  const [notifications, setNotifications] = useState([]);
  const [status, setStatus] = useState("Ready");
  const [error, setError] = useState("");
  const wsRef = useRef(null);

  const activeForm = role === "rider" ? riderForm : driverForm;
  const setActiveForm = role === "rider" ? setRiderForm : setDriverForm;
  const authenticatedRole = session?.role;
  const token = session?.access_token;

  const api = useMemo(() => createApi(token), [token]);

  useEffect(() => {
    if (!token) {
      wsRef.current?.close();
      wsRef.current = null;
      return;
    }

    const socket = new WebSocket(`${WS_BASE}/ws/notifications?token=${encodeURIComponent(token)}`);
    wsRef.current = socket;

    socket.onopen = () => pushNotification(setNotifications, "Notifications connected");
    socket.onmessage = (event) => pushNotification(setNotifications, event.data);
    socket.onerror = () => pushNotification(setNotifications, "Notification stream error");

    return () => socket.close();
  }, [token]);

  async function run(action, successMessage) {
    setError("");
    setStatus("Working...");
    try {
      const result = await action();
      setStatus(successMessage);
      return result;
    } catch (err) {
      setError(err.message);
      setStatus("Needs attention");
      return null;
    }
  }

  async function register() {
    await run(async () => {
      const response = await request("/api/v1/auth/register", {
        method: "POST",
        body: activeForm,
      });
      setSession(response);
      return response;
    }, `${capitalize(role)} registered`);
  }

  async function login() {
    await run(async () => {
      const response = await request("/api/v1/auth/login", {
        method: "POST",
        body: { email: activeForm.email, password: activeForm.password },
      });
      setSession(response);
      return response;
    }, `${capitalize(role)} signed in`);
  }

  async function estimateFare() {
    const response = await run(
      () => api.post("/api/v1/fare-estimates", { pickup, dropoff }),
      "Fare estimated",
    );
    if (response) setFare(response);
  }

  async function createRide() {
    const response = await run(
      () => api.post("/api/v1/rides", { pickup, dropoff }),
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
      setError("Create a ride and estimate fare first.");
      return;
    }

    const response = await run(
      () =>
        api.post("/api/v1/payments/authorize", {
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
    const response = await run(() => api.get(`/api/v1/rides/${ride.id}`), "Ride refreshed");
    if (response) setRide(response);
  }

  async function refreshPayment() {
    if (!payment) return;
    const response = await run(() => api.get(`/api/v1/payments/${payment.id}`), "Payment refreshed");
    if (response) setPayment(response);
  }

  async function updateDriverLocation() {
    await run(
      () => api.put(`/api/v1/drivers/${session?.user_id}/location`, driverLocation),
      "Driver location updated",
    );
  }

  async function setDriverAvailable() {
    await run(
      () => api.post(`/api/v1/drivers/${session?.user_id}/available`, {}),
      "Driver is available",
    );
  }

  async function findDriver() {
    if (!ride) {
      setError("Create a ride first.");
      return;
    }

    const response = await run(
      () =>
        api.post("/api/v1/matches", {
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
      () => api.post(`/api/v1/rides/${ride.id}/accept`, {}),
      "Ride accepted",
    );
    if (response) setRide(response);
  }

  async function startRide() {
    if (!ride) return;
    const response = await run(
      () => api.post(`/api/v1/rides/${ride.id}/start`, {}),
      "Ride started",
    );
    if (response) setRide(response);
  }

  async function completeRide() {
    if (!ride) return;
    const response = await run(
      () => api.post(`/api/v1/rides/${ride.id}/complete`, {}),
      "Ride completed",
    );
    if (response) setRide(response);
  }

  const isRider = authenticatedRole === "rider";
  const isDriver = authenticatedRole === "driver";

  return (
    <main className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <div className="brand-mark">
            <Route size={24} />
          </div>
          <div>
            <strong>Ride Sharing</strong>
            <span>Backend cockpit</span>
          </div>
        </div>

        <div className="role-switch" aria-label="Role">
          <button className={role === "rider" ? "active" : ""} onClick={() => setRole("rider")}>
            Rider
          </button>
          <button className={role === "driver" ? "active" : ""} onClick={() => setRole("driver")}>
            Driver
          </button>
        </div>

        <section className="panel compact">
          <h2>{capitalize(role)} access</h2>
          <Field label="Email">
            <input
              value={activeForm.email}
              onChange={(event) => setActiveForm({ ...activeForm, email: event.target.value })}
            />
          </Field>
          <Field label="Display name">
            <input
              value={activeForm.display_name}
              onChange={(event) =>
                setActiveForm({ ...activeForm, display_name: event.target.value })
              }
            />
          </Field>
          <Field label="Phone">
            <input
              value={activeForm.phone_number}
              onChange={(event) =>
                setActiveForm({ ...activeForm, phone_number: event.target.value })
              }
            />
          </Field>
          <Field label="Password">
            <input
              type="password"
              value={activeForm.password}
              onChange={(event) => setActiveForm({ ...activeForm, password: event.target.value })}
            />
          </Field>
          <div className="button-row">
            <IconButton icon={UserPlus} label="Register" onClick={register} />
            <IconButton icon={LogIn} label="Login" onClick={login} variant="secondary" />
          </div>
        </section>

        <section className="panel compact session-panel">
          <h2>Session</h2>
          <StatusLine icon={ShieldCheck} label="Role" value={session?.role ?? "Signed out"} />
          <StatusLine icon={Activity} label="User" value={session?.user_id ?? "-"} />
          <button className="wide ghost" onClick={() => setSession(null)}>
            <LogOut size={16} /> Sign out
          </button>
        </section>
      </aside>

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
          <section className="panel">
            <div className="panel-title">
              <h2>Rider flow</h2>
              <span>{isRider ? "active" : "needs rider token"}</span>
            </div>

            <div className="coordinate-grid">
              <CoordinateEditor title="Pickup" value={pickup} onChange={setPickup} />
              <CoordinateEditor title="Dropoff" value={dropoff} onChange={setDropoff} />
            </div>

            <div className="button-row">
              <IconButton icon={CircleDollarSign} label="Estimate" onClick={estimateFare} disabled={!isRider} />
              <IconButton icon={Car} label="Request ride" onClick={createRide} disabled={!isRider} />
              <IconButton
                icon={CheckCircle2}
                label="Authorize"
                onClick={authorizePayment}
                disabled={!isRider}
              />
              <IconButton icon={RefreshCcw} label="Refresh" onClick={refreshRide} disabled={!ride} />
            </div>

            <DataStrip
              items={[
                ["Fare", fare ? `${fare.amount.toFixed(2)} ${fare.currency}` : "-"],
                ["Distance", fare ? `${fare.distance_km.toFixed(2)} km` : "-"],
                ["Ride", ride?.id ?? "-"],
                ["Driver", ride?.driver_id || match?.driver_id || "-"],
              ]}
            />
          </section>

          <section className="panel">
            <div className="panel-title">
              <h2>Driver flow</h2>
              <span>{isDriver ? "active" : "needs driver token"}</span>
            </div>

            <CoordinateEditor title="Driver location" value={driverLocation} onChange={setDriverLocation} />

            <div className="button-row">
              <IconButton
                icon={Crosshair}
                label="Update"
                onClick={updateDriverLocation}
                disabled={!isDriver}
              />
              <IconButton icon={Radio} label="Available" onClick={setDriverAvailable} disabled={!isDriver} />
              <IconButton icon={Navigation} label="Match" onClick={findDriver} disabled={!ride} />
              <IconButton icon={CheckCircle2} label="Accept" onClick={acceptRide} disabled={!isDriver || !ride} />
              <IconButton icon={Car} label="Start" onClick={startRide} disabled={!isDriver || !ride} />
              <IconButton icon={CheckCircle2} label="Complete" onClick={completeRide} disabled={!isDriver || !ride} />
            </div>

            <DataStrip
              items={[
                ["Matched driver", match?.driver_id ?? "-"],
                ["Latitude", match ? match.latitude.toFixed(5) : "-"],
                ["Longitude", match ? match.longitude.toFixed(5) : "-"],
                ["Payment ID", payment?.id ?? "-"],
              ]}
            />
          </section>

          <section className="panel wide-panel">
            <div className="panel-title">
              <h2>Live state</h2>
              <button className="icon-only" onClick={refreshPayment} disabled={!payment} title="Refresh payment">
                <RefreshCcw size={18} />
              </button>
            </div>
            <div className="state-grid">
              <JsonBlock title="Ride" value={ride} />
              <JsonBlock title="Payment" value={payment} />
            </div>
          </section>

          <section className="panel notifications">
            <div className="panel-title">
              <h2>Notifications</h2>
              <Bell size={18} />
            </div>
            <div className="notification-list">
              {notifications.length === 0 ? (
                <p className="muted">No messages yet.</p>
              ) : (
                notifications.map((item, index) => <p key={`${item}-${index}`}>{item}</p>)
              )}
            </div>
          </section>
        </section>
      </section>
    </main>
  );
}

function createApi(token) {
  return {
    get: (path) => request(path, { token }),
    post: (path, body) => request(path, { method: "POST", token, body }),
    put: (path, body) => request(path, { method: "PUT", token, body }),
  };
}

async function request(path, options = {}) {
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

  return payload;
}

function sampleAccount(role) {
  return {
    email: `${role}.${Date.now()}@example.com`,
    password: PASSWORD,
    display_name: role === "rider" ? "Eren Rider" : "Eren Driver",
    phone_number: role === "rider" ? "+905551000001" : "+905551000002",
    role,
  };
}

function Field({ label, children }) {
  return (
    <label className="field">
      <span>{label}</span>
      {children}
    </label>
  );
}

function CoordinateEditor({ title, value, onChange }) {
  return (
    <div className="coordinate-editor">
      <h3>{title}</h3>
      <Field label="Latitude">
        <input
          type="number"
          step="0.000001"
          value={value.latitude}
          onChange={(event) => onChange({ ...value, latitude: Number(event.target.value) })}
        />
      </Field>
      <Field label="Longitude">
        <input
          type="number"
          step="0.000001"
          value={value.longitude}
          onChange={(event) => onChange({ ...value, longitude: Number(event.target.value) })}
        />
      </Field>
    </div>
  );
}

function IconButton({ icon: Icon, label, variant = "primary", ...props }) {
  return (
    <button className={`action ${variant}`} {...props}>
      <Icon size={17} />
      {label}
    </button>
  );
}

function Metric({ icon: Icon, label, value }) {
  return (
    <article className="metric">
      <Icon size={19} />
      <span>{label}</span>
      <strong>{value}</strong>
    </article>
  );
}

function StatusLine({ icon: Icon, label, value }) {
  return (
    <div className="status-line">
      <Icon size={16} />
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function DataStrip({ items }) {
  return (
    <div className="data-strip">
      {items.map(([label, value]) => (
        <div key={label}>
          <span>{label}</span>
          <strong>{value}</strong>
        </div>
      ))}
    </div>
  );
}

function JsonBlock({ title, value }) {
  return (
    <div className="json-block">
      <h3>{title}</h3>
      <pre>{value ? JSON.stringify(value, null, 2) : "null"}</pre>
    </div>
  );
}

function useLocalStorage(key, initialValue) {
  const [value, setValue] = useState(() => {
    const stored = localStorage.getItem(key);
    return stored ? JSON.parse(stored) : initialValue;
  });

  useEffect(() => {
    if (value === null) {
      localStorage.removeItem(key);
      return;
    }
    localStorage.setItem(key, JSON.stringify(value));
  }, [key, value]);

  return [value, setValue];
}

function pushNotification(setNotifications, message) {
  setNotifications((items) => [message, ...items].slice(0, 8));
}

function formatCoord(location) {
  return `${location.latitude.toFixed(3)}, ${location.longitude.toFixed(3)}`;
}

function capitalize(value) {
  return `${value.charAt(0).toUpperCase()}${value.slice(1)}`;
}

createRoot(document.getElementById("root")).render(<App />);
