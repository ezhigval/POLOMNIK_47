import Link from "next/link";

type PageIntroProps = {
  backHref?: string;
  backLabel?: string;
  title: string;
  description?: string;
};

export function PageIntro({
  backHref = "/",
  backLabel = "← На главную",
  title,
  description,
}: PageIntroProps) {
  return (
    <div className="space-y-3">
      <Link href={backHref} className="inline-block text-sm text-stone-500 transition hover:text-brand-800">
        {backLabel}
      </Link>
      <div>
        <h1 className="font-display text-3xl font-semibold text-stone-900 sm:text-4xl">{title}</h1>
        {description ? <p className="mt-2 max-w-2xl text-sm text-stone-600 sm:text-base">{description}</p> : null}
      </div>
    </div>
  );
}
