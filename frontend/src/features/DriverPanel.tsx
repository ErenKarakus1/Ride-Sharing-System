import { Car, CheckCircle2, Crosshair, Navigation, Radio } from "lucide-react";
import { CoordinateEditor } from "../components/CoordinateEditor";
import { DataStrip } from "../components/DataStrip";
import { IconButton } from "../components/IconButton";
import type { Location, Match, Payment, Ride } from "../types";

type DriverPanelProps = {
  isDriver: boolean;
  driverLocation: Location;
  ride: Ride | null;
  match: Match | null;
  payment: Payment | null;
  onDriverLocationChange: (location: Location) => void;
  onUpdateDriverLocation: () => void;
  onSetDriverAvailable: () => void;
  onFindDriver: () => void;
  onAcceptRide: () => void;
  onStartRide: () => void;
  onCompleteRide: () => void;
};

export function DriverPanel({
  isDriver,
  driverLocation,
  ride,
  match,
  payment,
  onDriverLocationChange,
  onUpdateDriverLocation,
  onSetDriverAvailable,
  onFindDriver,
  onAcceptRide,
  onStartRide,
  onCompleteRide,
}: DriverPanelProps) {
  return (
    <section className="panel">
      <div className="panel-title">
        <h2>Driver flow</h2>
        <span>{isDriver ? "active" : "needs driver token"}</span>
      </div>

      <CoordinateEditor title="Driver location" value={driverLocation} onChange={onDriverLocationChange} />

      <div className="button-row">
        <IconButton icon={Crosshair} label="Update" onClick={onUpdateDriverLocation} disabled={!isDriver} />
        <IconButton icon={Radio} label="Available" onClick={onSetDriverAvailable} disabled={!isDriver} />
        <IconButton icon={Navigation} label="Match" onClick={onFindDriver} disabled={!ride} />
        <IconButton icon={CheckCircle2} label="Accept" onClick={onAcceptRide} disabled={!isDriver || !ride} />
        <IconButton icon={Car} label="Start" onClick={onStartRide} disabled={!isDriver || !ride} />
        <IconButton icon={CheckCircle2} label="Complete" onClick={onCompleteRide} disabled={!isDriver || !ride} />
      </div>

      <DataStrip
        items={[
          ["Matched driver", match?.driver_id ?? "-"],
          ["Latitude", match ? match.latitude.toFixed(5) : "-"],
          ["Longitude", match ? match.longitude.toFixed(5) : "-"],
          ["Payment ID", payment?.id ?? "-"],
        ]}
      />
    </section>
  );
}
