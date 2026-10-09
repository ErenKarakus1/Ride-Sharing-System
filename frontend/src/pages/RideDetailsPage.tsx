import { LiveStatePanel } from "../features/LiveStatePanel";
import type { ActionState, Payment, Ride } from "../types";

type RideDetailsPageProps = {
  ride: Ride | null;
  payment: Payment | null;
  paymentAction: ActionState;
  onRefreshPayment: () => void;
};

export function RideDetailsPage({ ride, payment, paymentAction, onRefreshPayment }: RideDetailsPageProps) {
  return (
    <section className="page-grid single">
      <LiveStatePanel ride={ride} payment={payment} paymentAction={paymentAction} onRefreshPayment={onRefreshPayment} />
    </section>
  );
}
