import type { LucideIcon } from "lucide-react";

type MetricProps = {
  icon: LucideIcon;
  label: string;
  value: string;
};

export function Metric({ icon: Icon, label, value }: MetricProps) {
  return (
    <article className="metric">
      <Icon size={19} />
      <span>{label}</span>
      <strong>{value}</strong>
    </article>
  );
}
