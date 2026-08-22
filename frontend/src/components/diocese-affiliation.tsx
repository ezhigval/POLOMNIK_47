import { siteConfig } from "@/lib/site-config";

type DioceseAffiliationProps = {
  className?: string;
  textClassName?: string;
  linkClassName?: string;
};

function dioceseSiteHost(): string {
  try {
    return new URL(siteConfig.parentOrganization.url).hostname.replace(/^www\./, "");
  } catch {
    return siteConfig.parentOrganization.url;
  }
}

export function DioceseAffiliation({
  className,
  textClassName,
  linkClassName,
}: DioceseAffiliationProps) {
  return (
    <div className={className}>
      <p className={textClassName}>
        Паломническая служба является структурным подразделением Тихвинской епархии.
      </p>
      <a
        href={siteConfig.parentOrganization.url}
        target="_blank"
        rel="noopener noreferrer"
        className={linkClassName ?? "mt-2 inline-flex items-center underline underline-offset-2"}
      >
        Перейти на сайт епархии
        <span className="whitespace-nowrap"> ({dioceseSiteHost()})</span>
        <span aria-hidden="true"> ↗</span>
      </a>
    </div>
  );
}
