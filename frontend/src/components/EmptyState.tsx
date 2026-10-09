import type { LucideIcon } from "lucide-react";

type EmptyStateProps = {
  icon: LucideIcon;
  title: string;
  detail: string;
};

export function EmptyState({ icon: Icon, title, detail }: EmptyStateProps) {
  return (
    <div className="empty-state">
      <Icon size={18} />
      <strong>{title}</strong>
      <span>{detail}</span>
    </div>
  );
}
