import { Activity, LogIn, LogOut, RotateCcw, Route, ShieldCheck, UserPlus } from "lucide-react";
import type { AccountForm, Role, Sessions } from "../types";
import { capitalize } from "../lib/format";
import type { Page } from "../lib/routes";
import { Field } from "./Field";
import { IconButton } from "./IconButton";
import { StatusLine } from "./StatusLine";

type SidebarProps = {
  role: Role;
  page: Page;
  activeForm: AccountForm;
  sessions: Sessions;
  authLoading: boolean;
  onRoleChange: (role: Role) => void;
  onPageChange: (page: Page) => void;
  onFormChange: (form: AccountForm) => void;
  onRegister: () => void;
  onLogin: () => void;
  onSignOut: () => void;
  onResetTrip: () => void;
};

export function Sidebar({
  role,
  page,
  activeForm,
  sessions,
  authLoading,
  onRoleChange,
  onPageChange,
  onFormChange,
  onRegister,
  onLogin,
  onSignOut,
  onResetTrip,
}: SidebarProps) {
  const session = sessions[role];
  const navItems: Array<[Page, string]> =
    role === "rider"
      ? [
          ["auth", "Auth"],
          ["rider", "Rider"],
          ["ride", "Ride"],
          ["notifications", "Notifications"],
        ]
      : [
          ["auth", "Auth"],
          ["driver", "Driver"],
          ["ride", "Ride"],
          ["notifications", "Notifications"],
        ];

  return (
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
        <button className={role === "rider" ? "active" : ""} onClick={() => onRoleChange("rider")}>
          Rider {sessions.rider ? "on" : "off"}
        </button>
        <button className={role === "driver" ? "active" : ""} onClick={() => onRoleChange("driver")}>
          Driver {sessions.driver ? "on" : "off"}
        </button>
      </div>

      <nav className="app-nav" aria-label="Primary navigation">
        {navItems.map(([id, label]) => (
          <button key={id} className={page === id ? "active" : ""} onClick={() => onPageChange(id)}>
            {label}
          </button>
        ))}
      </nav>

      {page === "auth" && (
        <section className="panel compact">
          <h2>{capitalize(role)} access</h2>
          <Field label="Email">
            <input
              value={activeForm.email}
              onChange={(event) => onFormChange({ ...activeForm, email: event.target.value })}
            />
          </Field>
          <Field label="Display name">
            <input
              value={activeForm.display_name}
              onChange={(event) => onFormChange({ ...activeForm, display_name: event.target.value })}
            />
          </Field>
          <Field label="Phone">
            <input
              value={activeForm.phone_number}
              onChange={(event) => onFormChange({ ...activeForm, phone_number: event.target.value })}
            />
          </Field>
          <Field label="Password">
            <input
              type="password"
              value={activeForm.password}
              onChange={(event) => onFormChange({ ...activeForm, password: event.target.value })}
            />
          </Field>
          <div className="button-row">
            <IconButton icon={UserPlus} label="Register" loading={authLoading} loadingLabel="Registering..." onClick={onRegister} />
            <IconButton icon={LogIn} label="Login" loading={authLoading} loadingLabel="Signing in..." onClick={onLogin} variant="secondary" />
          </div>
        </section>
      )}

      <section className="panel compact session-panel">
        <h2>Session</h2>
        <StatusLine icon={ShieldCheck} label="Role" value={session?.role ?? "Signed out"} />
        <StatusLine icon={Activity} label="User" value={session?.user_id ?? "-"} />
        <div className="button-row stacked">
          <button className="wide ghost" onClick={onResetTrip}>
            <RotateCcw size={16} /> Reset trip
          </button>
          <button className="wide ghost" onClick={onSignOut} disabled={!session}>
            <LogOut size={16} /> Sign out
          </button>
        </div>
      </section>
    </aside>
  );
}
