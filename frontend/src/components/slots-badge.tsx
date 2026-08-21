import { getSlotsAvailability, slotsLabel } from "@/lib/format";

type SlotsBadgeProps = {
  slotsLeft: number;
};

export function SlotsBadge({ slotsLeft }: SlotsBadgeProps) {
  const availability = getSlotsAvailability(slotsLeft);

  const styles = {
    available: "bg-emerald-50 text-emerald-800 ring-emerald-200",
    low: "bg-amber-50 text-amber-900 ring-amber-200",
    sold_out: "bg-stone-100 text-stone-600 ring-stone-200",
  };

  return (
    <span
      className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset ${styles[availability]}`}
    >
      {slotsLabel(slotsLeft)}
    </span>
  );
}
