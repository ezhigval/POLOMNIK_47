import Link from "next/link";

type EmptyStateProps = {
  title: string;
  description: string;
  actionHref: string;
  actionLabel: string;
  secondaryHref?: string;
  secondaryLabel?: string;
};

export function EmptyState({
  title,
  description,
  actionHref,
  actionLabel,
  secondaryHref,
  secondaryLabel,
}: EmptyStateProps) {
  return (
    <div className="rounded-2xl border border-dashed border-stone-300 bg-white p-10 text-center">
      <h2 className="text-lg font-semibold text-stone-900">{title}</h2>
      <p className="mx-auto mt-2 max-w-md text-sm text-stone-600">{description}</p>
      <div className="mt-6 flex flex-wrap justify-center gap-3">
        <Link href={actionHref} className="btn-primary">
          {actionLabel}
        </Link>
        {secondaryHref && secondaryLabel ? (
          <Link href={secondaryHref} className="btn-secondary">
            {secondaryLabel}
          </Link>
        ) : null}
      </div>
    </div>
  );
}
