type BadgeVariant = "neutral" | "success" | "warning" | "danger" | "brand";

const variantClasses: Record<BadgeVariant, string> = {
  neutral: "bg-stone-100 text-stone-700",
  success: "bg-emerald-100 text-emerald-800",
  warning: "bg-amber-100 text-amber-800",
  danger: "bg-red-100 text-red-800",
  brand: "bg-brand-100 text-brand-900",
};

type StatusBadgeProps = {
  children: React.ReactNode;
  variant?: BadgeVariant;
};

export function StatusBadge({ children, variant = "neutral" }: StatusBadgeProps) {
  return (
    <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${variantClasses[variant]}`}>
      {children}
    </span>
  );
}

export function bookingStatusVariant(status: string): BadgeVariant {
  switch (status) {
    case "NEW":
      return "brand";
    case "CONTACTED":
      return "warning";
    case "CONFIRMED":
      return "success";
    case "COMPLETED":
      return "neutral";
    case "CANCELLED":
      return "danger";
    default:
      return "neutral";
  }
}

export function paymentStatusVariant(status: string): BadgeVariant {
  switch (status) {
    case "PAID":
      return "success";
    case "AWAITING_PAYMENT":
      return "warning";
    case "UNPAID":
      return "brand";
    case "NOT_REQUIRED":
      return "neutral";
    default:
      return "neutral";
  }
}

export function syncStatusVariant(status: string): BadgeVariant {
  switch (status) {
    case "synced":
    case "processed":
      return "success";
    case "pending":
      return "warning";
    case "failed":
      return "danger";
    default:
      return "neutral";
  }
}
