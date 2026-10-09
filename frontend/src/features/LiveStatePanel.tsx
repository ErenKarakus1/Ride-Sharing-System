import { CreditCard, MapPin, RefreshCcw, Route } from "lucide-react";
import { DataStrip } from "../components/DataStrip";
import { JsonBlock } from "../components/JsonBlock";
import { StatusLine } from "../components/StatusLine";
import { formatCoord } from "../lib/format";
import type { ActionState, Payment, Ride } from "../types";

type LiveStatePanelProps = {
  ride: Ride | null;
  payment: Payment | null;
  paymentAction: ActionState;
  onRefreshPayment: () => void;
};

export function LiveStatePanel({ ride, payment, paymentAction, onRefreshPayment }: LiveStatePanelProps) {
  return (
    <section className="panel wide-panel">
      <div className="panel-title">
        <h2>Live state</h2>
        <span>{paymentAction.loading ? "working" : paymentAction.error || paymentAction.message}</span>
        <button className="icon-only" onClick={onRefreshPayment} disabled={!payment || paymentAction.loading} title="Refresh payment">
          <RefreshCcw size={18} />
        </button>
      </div>

      <div className="summary-grid">
        <article className="summary-card">
          <div className="summary-heading">
            <Route size={18} />
            <h3>Ride</h3>
            <span className={ride ? `status-pill ${ride.status}` : "status-pill"}>{ride?.status ?? "none"}</span>
          </div>
          <StatusLine icon={MapPin} label="Pickup" value={ride ? formatCoord(ride.pickup) : "-"} />
          <StatusLine icon={MapPin} label="Dropoff" value={ride ? formatCoord(ride.dropoff) : "-"} />
          <DataStrip
            items={[
              ["Ride ID", ride?.id ?? "-"],
              ["Rider", ride?.rider_id ?? "-"],
              ["Driver", ride?.driver_id ?? "-"],
              ["Updated", ride?.updated_at ?? "-"],
            ]}
          />
        </article>

        <article className="summary-card">
          <div className="summary-heading">
            <CreditCard size={18} />
            <h3>Payment</h3>
            <span className={payment ? `status-pill ${payment.status}` : "status-pill"}>{payment?.status ?? "none"}</span>
          </div>
          <DataStrip
            items={[
              ["Payment ID", payment?.id ?? "-"],
              ["Amount", payment ? `${payment.amount.toFixed(2)} ${payment.currency}` : "-"],
              ["Rider", payment?.rider_id ?? "-"],
              ["Driver", payment?.driver_id ?? "-"],
            ]}
          />
        </article>
      </div>

      <div className="state-grid">
        <JsonBlock title="Ride" value={ride} />
        <JsonBlock title="Payment" value={payment} />
      </div>
    </section>
  );
}
