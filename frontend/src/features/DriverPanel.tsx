import { Car, CheckCircle2, Crosshair, Navigation, Radio } from "lucide-react";
import { CoordinateEditor } from "../components/CoordinateEditor";
import { DataStrip } from "../components/DataStrip";
import { EmptyState } from "../components/EmptyState";
import { FlowSteps } from "../components/FlowSteps";
import { IconButton } from "../components/IconButton";
import { locationPresets } from "../config";
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
  const locationUpdated = Boolean(locationAction.message);
  const driverAvailable = locationAction.message === "Driver is available";
  const canUpdateLocation = isDriver && !locationUpdated;
  const canSetAvailable = isDriver && locationUpdated && !driverAvailable;
  const canAccept = isDriver && Boolean(match) && ride?.status === "requested";
  const canStart = isDriver && ride?.status === "accepted";
  const canComplete = isDriver && ride?.status === "started";
  const showMatchData = Boolean(match || payment);
  const matchedDriver = match?.driver_id ?? ride?.driver_id ?? "-";
  const canMatchRide = isDriver && driverAvailable && Boolean(ride) && !match && !ride?.driver_id;

  return (
    <section className="panel">
      <div className="panel-title">
        <h2>Driver flow</h2>
        <span>{panelStatus(isDriver, locationAction, matchingAction, lifecycleAction)}</span>
      </div>

      <CoordinateEditor
        title="Driver location"
        value={driverLocation}
        presets={locationPresets}
        onChange={onDriverLocationChange}
      />

      <FlowSteps
        steps={[
          { label: "Update location", done: Boolean(locationAction.message), active: isDriver && !locationAction.message },
          { label: "Go available", done: locationAction.message === "Driver is available", active: Boolean(locationAction.message) },
          { label: "Match ride", done: Boolean(match || ride?.driver_id), active: canMatchRide },
          { label: "Accept", done: ride?.status === "accepted" || ride?.status === "started" || ride?.status === "completed", active: Boolean(match && ride?.status === "requested") },
          { label: "Start", done: ride?.status === "started" || ride?.status === "completed", active: ride?.status === "accepted" },
          { label: "Complete", done: ride?.status === "completed", active: ride?.status === "started" },
        ]}
      />

      {isDriver && (
        <div className="button-row">
          {canUpdateLocation && (
            <IconButton icon={Crosshair} label="Update" loading={locationAction.loading} loadingLabel="Updating..." onClick={onUpdateDriverLocation} />
          )}
          {canSetAvailable && (
            <IconButton icon={Radio} label="Available" loading={locationAction.loading} loadingLabel="Saving..." onClick={onSetDriverAvailable} />
          )}
          {canAccept && (
            <IconButton icon={CheckCircle2} label="Accept" loading={lifecycleAction.loading} loadingLabel="Accepting..." onClick={onAcceptRide} />
          )}
          {canStart && (
            <IconButton icon={Car} label="Start" loading={lifecycleAction.loading} loadingLabel="Starting..." onClick={onStartRide} />
          )}
          {canComplete && (
            <IconButton icon={CheckCircle2} label="Complete" loading={lifecycleAction.loading} loadingLabel="Completing..." onClick={onCompleteRide} />
          )}
        </div>
      )}

      {showMatchData && (
        <DataStrip
          items={[
            ["Matched driver", matchedDriver],
            ["Latitude", match ? match.latitude.toFixed(5) : "-"],
            ["Longitude", match ? match.longitude.toFixed(5) : "-"],
            ["Payment ID", payment?.id ?? "-"],
          ]}
        />
      )}

      {!isDriver && (
        <EmptyState icon={Car} title="Driver session required" detail="Register or sign in as a driver before accepting trips." />
      )}
      {isDriver && !locationUpdated && (
        <EmptyState icon={Crosshair} title="Set driver position" detail="Update the driver's current location before marking them available." />
      )}
      {isDriver && locationUpdated && !driverAvailable && (
        <EmptyState icon={Radio} title="Location saved" detail="Mark the driver available so they can be matched to the active ride." />
      )}
      {isDriver && driverAvailable && !ride && (
        <EmptyState icon={Navigation} title="Waiting for a ride" detail="Update driver location, mark the driver available, then match an active ride." />
      )}
      {isDriver && driverAvailable && ride && !ride.driver_id && !match && (
        <EmptyState icon={Navigation} title="Finding driver" detail="Matching runs in the background after payment authorization." />
      )}
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
