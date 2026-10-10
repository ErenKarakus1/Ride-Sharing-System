import { Car, CheckCircle2, CircleDollarSign, RefreshCcw } from "lucide-react";
import { CoordinateEditor } from "../components/CoordinateEditor";
import { DataStrip } from "../components/DataStrip";
import { EmptyState } from "../components/EmptyState";
import { FlowSteps } from "../components/FlowSteps";
import { IconButton } from "../components/IconButton";
import { locationPresets } from "../config";
import type { ActionState, FareEstimate, Location, Match, Payment, Ride } from "../types";

type RiderPanelProps = {
  isRider: boolean;
  action: ActionState;
  pickup: Location;
  dropoff: Location;
  fare: FareEstimate | null;
  ride: Ride | null;
  match: Match | null;
  payment: Payment | null;
  onPickupChange: (location: Location) => void;
  onDropoffChange: (location: Location) => void;
  onEstimateFare: () => void;
  onCreateRide: () => void;
  onAuthorizePayment: () => void;
  onRefreshRide: () => void;
};

export function RiderPanel({
  isRider,
  action,
  pickup,
  dropoff,
  fare,
  ride,
  match,
  payment,
  onPickupChange,
  onDropoffChange,
  onEstimateFare,
  onCreateRide,
  onAuthorizePayment,
  onRefreshRide,
}: RiderPanelProps) {
  return (
    <section className="panel">
      <div className="panel-title">
        <h2>Rider flow</h2>
        <span>{action.loading ? "working" : action.error || action.message || (isRider ? "active" : "needs rider token")}</span>
      </div>

      <div className="coordinate-grid">
        <CoordinateEditor title="Pickup" value={pickup} presets={locationPresets} onChange={onPickupChange} />
        <CoordinateEditor title="Dropoff" value={dropoff} presets={locationPresets} onChange={onDropoffChange} />
      </div>

      <FlowSteps
        steps={[
          { label: "Estimate fare", done: Boolean(fare), active: isRider && !fare },
          { label: "Request ride", done: Boolean(ride), active: Boolean(fare && !ride) },
          { label: "Authorize payment", done: Boolean(payment), active: Boolean(ride && !payment) },
          { label: "Match driver", done: Boolean(match), active: Boolean(ride && !match) },
          { label: "Complete trip", done: ride?.status === "completed", active: Boolean(match && ride?.status !== "completed") },
        ]}
      />

      <div className="button-row">
        <IconButton icon={CircleDollarSign} label="Estimate" loading={action.loading} loadingLabel="Estimating..." onClick={onEstimateFare} disabled={!isRider} />
        <IconButton icon={Car} label="Request ride" loading={action.loading} loadingLabel="Requesting..." onClick={onCreateRide} disabled={!isRider} />
        <IconButton icon={CheckCircle2} label="Authorize" loading={action.loading} loadingLabel="Authorizing..." onClick={onAuthorizePayment} disabled={!isRider} />
        <IconButton icon={RefreshCcw} label="Refresh" loading={action.loading} loadingLabel="Refreshing..." onClick={onRefreshRide} disabled={!ride} />
      </div>

      <DataStrip
        items={[
          ["Fare", fare ? `${fare.amount.toFixed(2)} ${fare.currency}` : "-"],
          ["Distance", fare ? `${fare.distance_km.toFixed(2)} km` : "-"],
          ["Ride", ride?.id ?? "-"],
          ["Payment", payment?.status ?? "-"],
        ]}
      />

      {!isRider && (
        <EmptyState icon={Car} title="Rider session required" detail="Register or sign in as a rider before running this flow." />
      )}
      {isRider && !ride && (
        <EmptyState icon={CircleDollarSign} title="Ready to plan" detail="Estimate a fare, then request a ride from the selected pickup and dropoff." />
      )}
    </section>
  );
}
