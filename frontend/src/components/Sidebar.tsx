import { Activity, LogIn, LogOut, Route, ShieldCheck, UserPlus } from "lucide-react";
import type { AccountForm, Role, Sessions } from "../types";
import { capitalize } from "../lib/format";
import { Field } from "./Field";
import { IconButton } from "./IconButton";
import { StatusLine } from "./StatusLine";

type SidebarProps = {
  role: Role;
  activeForm: AccountForm;
  sessions: Sessions;
  onRoleChange: (role: Role) => void;
  onFormChange: (form: AccountForm) => void;
  onRegister: () => void;
  onLogin: () => void;
  onSignOut: () => void;
};

export function Sidebar({
  role,
  activeForm,
  sessions,
  onRoleChange,
  onFormChange,
  onRegister,
  onLogin,
  onSignOut,
}: SidebarProps) {
  const session = sessions[role];

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

      <section className="panel compact">
        <h2>{capitalize(role)} access</h2>
        <Field label="Email">
          <input value={activeForm.email} onChange={(event) => onFormChange({ ...activeForm, email: event.target.value })} />
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
          <IconButton icon={UserPlus} label="Register" onClick={onRegister} />
          <IconButton icon={LogIn} label="Login" onClick={onLogin} variant="secondary" />
        </div>
      </section>

      <section className="panel compact session-panel">
        <h2>Session</h2>
        <StatusLine icon={ShieldCheck} label="Role" value={session?.role ?? "Signed out"} />
        <StatusLine icon={Activity} label="User" value={session?.user_id ?? "-"} />
        <button className="wide ghost" onClick={onSignOut}>
          <LogOut size={16} /> Sign out
        </button>
      </section>
    </aside>
  );
}
