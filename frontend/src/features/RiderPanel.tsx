import { Car, CheckCircle2, CircleDollarSign, RefreshCcw } from "lucide-react";
import { CoordinateEditor } from "../components/CoordinateEditor";
import { DataStrip } from "../components/DataStrip";
import { IconButton } from "../components/IconButton";
import type { ActionState, FareEstimate, Location, Match, Ride } from "../types";

type RiderPanelProps = {
  isRider: boolean;
  action: ActionState;
  pickup: Location;
  dropoff: Location;
  fare: FareEstimate | null;
  ride: Ride | null;
  match: Match | null;
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
        <CoordinateEditor title="Pickup" value={pickup} onChange={onPickupChange} />
        <CoordinateEditor title="Dropoff" value={dropoff} onChange={onDropoffChange} />
      </div>

      <div className="button-row">
        <IconButton icon={CircleDollarSign} label="Estimate" onClick={onEstimateFare} disabled={!isRider || action.loading} />
        <IconButton icon={Car} label="Request ride" onClick={onCreateRide} disabled={!isRider || action.loading} />
        <IconButton icon={CheckCircle2} label="Authorize" onClick={onAuthorizePayment} disabled={!isRider || action.loading} />
        <IconButton icon={RefreshCcw} label="Refresh" onClick={onRefreshRide} disabled={!ride || action.loading} />
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
  );
}
