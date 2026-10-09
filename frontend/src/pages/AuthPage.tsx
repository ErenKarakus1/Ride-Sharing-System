import { Activity, Car, CircleDollarSign, MapPinned, ShieldCheck } from "lucide-react";
import { DataStrip } from "../components/DataStrip";
import { Metric } from "../components/Metric";
import type { ActionState, FareEstimate, Match, Payment, Ride, Sessions } from "../types";

type AuthPageProps = {
  sessions: Sessions;
  action: ActionState;
  fare: FareEstimate | null;
  ride: Ride | null;
  match: Match | null;
  payment: Payment | null;
};

export function AuthPage({ sessions, action, fare, ride, match, payment }: AuthPageProps) {
  return (
    <section className="page-grid single">
      <section className="metrics">
        <Metric icon={ShieldCheck} label="Rider" value={sessions.rider ? "Signed in" : "Missing"} />
        <Metric icon={Activity} label="Driver" value={sessions.driver ? "Signed in" : "Missing"} />
        <Metric icon={Car} label="Ride" value={ride?.status ?? "No ride"} />
        <Metric icon={CircleDollarSign} label="Payment" value={payment?.status ?? "No payment"} />
      </section>

      <section className="panel">
        <div className="panel-title">
          <h2>Operations state</h2>
          <span>{action.loading ? "working" : action.error || action.message || "choose a role"}</span>
        </div>
        <DataStrip
          items={[
            ["Rider", sessions.rider ? sessions.rider.email : "Signed out"],
            ["Driver", sessions.driver ? sessions.driver.email : "Signed out"],
            ["Fare", fare ? `${fare.amount.toFixed(2)} ${fare.currency}` : "-"],
            ["Matched driver", match?.driver_id ?? "-"],
          ]}
        />
      </section>

      <section className="flow-grid">
        <article className="panel compact">
          <div className="panel-title">
            <h3>Session IDs</h3>
            <span>{sessions.rider && sessions.driver ? "ready" : "incomplete"}</span>
          </div>
          <DataStrip
            items={[
              ["Rider ID", sessions.rider?.user_id ?? "-"],
              ["Driver ID", sessions.driver?.user_id ?? "-"],
            ]}
          />
        </article>

        <article className="panel compact">
          <div className="panel-title">
            <h3>Trip checkpoint</h3>
            <span>{ride?.id ? "started" : "waiting"}</span>
          </div>
          <DataStrip
            items={[
              ["Ride ID", ride?.id ?? "-"],
              ["Ride status", ride?.status ?? "-"],
              ["Payment ID", payment?.id ?? "-"],
              ["Payment status", payment?.status ?? "-"],
            ]}
          />
        </article>
      </section>

      <section className="panel compact">
        <div className="panel-title">
          <h3>Flow readiness</h3>
          <span>{readinessLabel(sessions, fare, ride, match, payment)}</span>
        </div>
        <div className="timeline">
          <Step done={Boolean(sessions.rider)} label="Rider session" />
          <Step done={Boolean(sessions.driver)} label="Driver session" />
          <Step done={Boolean(fare)} label="Fare estimate" />
          <Step done={Boolean(ride)} label="Ride request" />
          <Step done={Boolean(match)} label="Driver match" />
          <Step done={payment?.status === "captured"} label="Captured payment" />
        </div>
      </section>
    </section>
  );
}

function Step({ done, label }: { done: boolean; label: string }) {
  return (
    <div className={done ? "step done" : "step"}>
      <MapPinned size={16} />
      <span>{label}</span>
    </div>
  );
}

function readinessLabel(
  sessions: Sessions,
  fare: FareEstimate | null,
  ride: Ride | null,
  match: Match | null,
  payment: Payment | null,
) {
  if (payment?.status === "captured") return "complete";
  if (match) return "driver matched";
  if (ride) return "ride active";
  if (fare) return "fare ready";
  if (sessions.rider && sessions.driver) return "sessions ready";
  return "needs sessions";
}
