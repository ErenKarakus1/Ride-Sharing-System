import { CircleDollarSign, Clock3, MapPin, Navigation } from "lucide-react";
import { Metric } from "../components/Metric";
import { RouteSummary } from "../components/RouteSummary";
import { RiderPanel } from "../features/RiderPanel";
import { formatCoord } from "../lib/format";
import type { ActionState, FareEstimate, Location, Match, Payment, Ride } from "../types";

type RiderPageProps = {
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

export function RiderPage(props: RiderPageProps) {
  return (
    <section className="page-grid">
      <section className="metrics">
        <Metric icon={MapPin} label="Pickup" value={formatCoord(props.pickup)} />
        <Metric icon={Navigation} label="Dropoff" value={formatCoord(props.dropoff)} />
        <Metric icon={Clock3} label="Ride status" value={props.ride?.status ?? "No ride"} />
        <Metric
          icon={CircleDollarSign}
          label="Fare"
          value={props.fare ? `${props.fare.amount.toFixed(2)} ${props.fare.currency}` : "No fare"}
        />
      </section>
      <RouteSummary pickup={props.pickup} dropoff={props.dropoff} fare={props.fare} ride={props.ride} match={props.match} />
      <RiderPanel {...props} />
    </section>
  );
}
