import { Activity, ShieldCheck } from "lucide-react";
import { DataStrip } from "../components/DataStrip";
import type { ActionState, Sessions } from "../types";

type AuthPageProps = {
  sessions: Sessions;
  action: ActionState;
};

export function AuthPage({ sessions, action }: AuthPageProps) {
  return (
    <section className="page-grid single">
      <section className="panel">
        <div className="panel-title">
          <h2>Access</h2>
          <span>{action.loading ? "working" : action.error || action.message || "choose a role"}</span>
        </div>
        <DataStrip
          items={[
            ["Rider", sessions.rider ? sessions.rider.email : "Signed out"],
            ["Driver", sessions.driver ? sessions.driver.email : "Signed out"],
            ["Rider ID", sessions.rider?.user_id ?? "-"],
            ["Driver ID", sessions.driver?.user_id ?? "-"],
          ]}
        />
      </section>

      <section className="metrics">
        <article className="metric">
          <ShieldCheck size={19} />
          <span>Rider session</span>
          <strong>{sessions.rider ? "Ready" : "Missing"}</strong>
        </article>
        <article className="metric">
          <Activity size={19} />
          <span>Driver session</span>
          <strong>{sessions.driver ? "Ready" : "Missing"}</strong>
        </article>
      </section>
    </section>
  );
}
