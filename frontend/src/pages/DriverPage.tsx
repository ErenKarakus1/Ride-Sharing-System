import { Car, Clock3, MapPin, Navigation } from "lucide-react";
import { Metric } from "../components/Metric";
import { DriverPanel } from "../features/DriverPanel";
import { formatCoord } from "../lib/format";
import type { ActionState, Location, Match, Payment, Ride } from "../types";

type DriverPageProps = {
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

export function DriverPage(props: DriverPageProps) {
  return (
    <section className="page-grid">
      <section className="metrics">
        <Metric icon={MapPin} label="Driver location" value={formatCoord(props.driverLocation)} />
        <Metric icon={Navigation} label="Matched driver" value={props.match?.driver_id ?? "No match"} />
        <Metric icon={Clock3} label="Ride status" value={props.ride?.status ?? "No ride"} />
        <Metric icon={Car} label="Assigned driver" value={props.ride?.driver_id ?? "-"} />
      </section>
      <DriverPanel {...props} />
    </section>
  );
}
