import type { ButtonHTMLAttributes } from "react";
import type { LucideIcon } from "lucide-react";

type IconButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  icon: LucideIcon;
  label: string;
  variant?: "primary" | "secondary";
};

export function IconButton({ icon: Icon, label, variant = "primary", ...props }: IconButtonProps) {
  return (
    <button className={`action ${variant}`} {...props}>
      <Icon size={17} />
      {label}
    </button>
  );
}
