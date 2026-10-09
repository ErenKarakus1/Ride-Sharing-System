import type { ButtonHTMLAttributes } from "react";
import type { LucideIcon } from "lucide-react";

type IconButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  icon: LucideIcon;
  label: string;
  loading?: boolean;
  loadingLabel?: string;
  variant?: "primary" | "secondary";
};

export function IconButton({ icon: Icon, label, loading = false, loadingLabel, variant = "primary", disabled, ...props }: IconButtonProps) {
  return (
    <button className={`action ${variant}`} disabled={disabled || loading} {...props}>
      <Icon size={17} />
      {loading ? loadingLabel ?? "Working..." : label}
    </button>
  );
}
