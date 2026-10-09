import { Car, CircleDollarSign, MapPin, Navigation } from "lucide-react";
import { formatCoord } from "../lib/format";
import type { FareEstimate, Location, Match, Ride } from "../types";
import { StatusLine } from "./StatusLine";

type RouteSummaryProps = {
  pickup: Location;
  dropoff: Location;
  fare: FareEstimate | null;
  ride: Ride | null;
  match: Match | null;
};

export function RouteSummary({ pickup, dropoff, fare, ride, match }: RouteSummaryProps) {
  return (
    <section className="panel compact route-summary">
      <div className="panel-title">
        <h3>Route summary</h3>
        <span>{ride?.status ?? "planning"}</span>
      </div>
      <div className="route-summary-grid">
        <StatusLine icon={MapPin} label="Pickup" value={formatCoord(pickup)} />
        <StatusLine icon={Navigation} label="Dropoff" value={formatCoord(dropoff)} />
        <StatusLine icon={CircleDollarSign} label="Fare" value={fare ? `${fare.amount.toFixed(2)} ${fare.currency}` : "-"} />
        <StatusLine icon={Car} label="Driver" value={ride?.driver_id || match?.driver_id || "-"} />
      </div>
    </section>
  );
}
