import { Breadcrumbs, type BreadcrumbItem } from "@/components/breadcrumbs";
import { SmartBackButton } from "@/components/smart-back-button";

type DetailPageNavProps = {
  fallbackHref: string;
  backLabel?: string;
  breadcrumbs: BreadcrumbItem[];
};

export function DetailPageNav({ fallbackHref, backLabel = "Назад", breadcrumbs }: DetailPageNavProps) {
  return (
    <div className="mb-6 space-y-3">
      <SmartBackButton fallbackHref={fallbackHref} label={backLabel} />
      <Breadcrumbs items={breadcrumbs} className="mb-0" />
    </div>
  );
}
