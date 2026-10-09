import { RefreshCcw } from "lucide-react";
import { JsonBlock } from "../components/JsonBlock";
import type { Payment, Ride } from "../types";

type LiveStatePanelProps = {
  ride: Ride | null;
  payment: Payment | null;
  onRefreshPayment: () => void;
};

export function LiveStatePanel({ ride, payment, onRefreshPayment }: LiveStatePanelProps) {
  return (
    <section className="panel wide-panel">
      <div className="panel-title">
        <h2>Live state</h2>
        <button className="icon-only" onClick={onRefreshPayment} disabled={!payment} title="Refresh payment">
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
