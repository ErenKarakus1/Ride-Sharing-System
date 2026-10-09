import { Car, CheckCircle2, CircleDollarSign, RefreshCcw } from "lucide-react";
import { CoordinateEditor } from "../components/CoordinateEditor";
import { DataStrip } from "../components/DataStrip";
import { IconButton } from "../components/IconButton";
import type { FareEstimate, Location, Match, Ride } from "../types";

type RiderPanelProps = {
  isRider: boolean;
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
        <span>{isRider ? "active" : "needs rider token"}</span>
      </div>

      <div className="coordinate-grid">
        <CoordinateEditor title="Pickup" value={pickup} onChange={onPickupChange} />
        <CoordinateEditor title="Dropoff" value={dropoff} onChange={onDropoffChange} />
      </div>

      <div className="button-row">
        <IconButton icon={CircleDollarSign} label="Estimate" onClick={onEstimateFare} disabled={!isRider} />
        <IconButton icon={Car} label="Request ride" onClick={onCreateRide} disabled={!isRider} />
        <IconButton icon={CheckCircle2} label="Authorize" onClick={onAuthorizePayment} disabled={!isRider} />
        <IconButton icon={RefreshCcw} label="Refresh" onClick={onRefreshRide} disabled={!ride} />
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
