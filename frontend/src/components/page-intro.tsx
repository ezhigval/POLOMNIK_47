import { SmartBackButton } from "@/components/smart-back-button";

type PageIntroProps = {
  backHref?: string;
  backLabel?: string;
  showBack?: boolean;
  title: string;
  description?: string;
};

export function PageIntro({
  backHref = "/",
  backLabel = "На главную",
  showBack = true,
  title,
  description,
}: PageIntroProps) {
  return (
    <div className="space-y-3">
      {showBack ? <SmartBackButton fallbackHref={backHref} label={backLabel} /> : null}
      <div>
        <h1 className="font-display text-3xl font-semibold text-stone-900 sm:text-4xl">{title}</h1>
        {description ? <p className="mt-2 max-w-2xl text-sm text-stone-600 sm:text-base">{description}</p> : null}
      </div>
    </div>
  );
}
