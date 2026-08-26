import { SmartBackButton } from "@/components/smart-back-button";

type ReturnToListLinkProps = {
  fallbackHref: string;
  label: string;
};

/** Bottom-of-page «return to list» — same smart back as the header control. */
export function ReturnToListLink({ fallbackHref, label }: ReturnToListLinkProps) {
  return (
    <div className="border-t border-stone-100 pt-6">
      <SmartBackButton fallbackHref={fallbackHref} label={label} />
    </div>
  );
}
