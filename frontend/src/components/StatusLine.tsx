import type { LucideIcon } from "lucide-react";

type StatusLineProps = {
  icon: LucideIcon;
  label: string;
  value: string;
};

export function StatusLine({ icon: Icon, label, value }: StatusLineProps) {
  return (
    <div className="status-line">
      <Icon size={16} />
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}
