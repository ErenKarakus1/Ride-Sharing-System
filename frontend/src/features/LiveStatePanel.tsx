import { RefreshCcw } from "lucide-react";
import { JsonBlock } from "../components/JsonBlock";
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
      <div className="state-grid">
        <JsonBlock title="Ride" value={ride} />
        <JsonBlock title="Payment" value={payment} />
      </div>
    </section>
  );
}
