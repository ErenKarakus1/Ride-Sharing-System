import { Car, CheckCircle2, Crosshair, Navigation, Radio } from "lucide-react";
import { CoordinateEditor } from "../components/CoordinateEditor";
import { DataStrip } from "../components/DataStrip";
import { FlowSteps } from "../components/FlowSteps";
import { IconButton } from "../components/IconButton";
import type { ActionState, Location, Match, Payment, Ride } from "../types";

type DriverPanelProps = {
  isDriver: boolean;
  locationAction: ActionState;
  matchingAction: ActionState;
  lifecycleAction: ActionState;
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
  locationAction,
  matchingAction,
  lifecycleAction,
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
        <span>{panelStatus(isDriver, locationAction, matchingAction, lifecycleAction)}</span>
      </div>

      <CoordinateEditor title="Driver location" value={driverLocation} onChange={onDriverLocationChange} />

      <FlowSteps
        steps={[
          { label: "Update location", done: Boolean(locationAction.message), active: isDriver && !locationAction.message },
          { label: "Go available", done: locationAction.message === "Driver is available", active: Boolean(locationAction.message) },
          { label: "Match ride", done: Boolean(match), active: Boolean(ride && !match) },
          { label: "Accept", done: ride?.status === "accepted" || ride?.status === "started" || ride?.status === "completed", active: Boolean(match && ride?.status === "requested") },
          { label: "Start", done: ride?.status === "started" || ride?.status === "completed", active: ride?.status === "accepted" },
          { label: "Complete", done: ride?.status === "completed", active: ride?.status === "started" },
        ]}
      />

      <div className="button-row">
        <IconButton icon={Crosshair} label="Update" onClick={onUpdateDriverLocation} disabled={!isDriver || locationAction.loading} />
        <IconButton icon={Radio} label="Available" onClick={onSetDriverAvailable} disabled={!isDriver || locationAction.loading} />
        <IconButton icon={Navigation} label="Match" onClick={onFindDriver} disabled={!ride || matchingAction.loading} />
        <IconButton icon={CheckCircle2} label="Accept" onClick={onAcceptRide} disabled={!isDriver || !ride || lifecycleAction.loading} />
        <IconButton icon={Car} label="Start" onClick={onStartRide} disabled={!isDriver || !ride || lifecycleAction.loading} />
        <IconButton icon={CheckCircle2} label="Complete" onClick={onCompleteRide} disabled={!isDriver || !ride || lifecycleAction.loading} />
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

function panelStatus(isDriver: boolean, ...actions: ActionState[]) {
  const active = actions.find((action) => action.loading || action.error || action.message);
  if (active?.loading) return "working";
  if (active?.error) return active.error;
  if (active?.message) return active.message;
  return isDriver ? "active" : "needs driver token";
}
