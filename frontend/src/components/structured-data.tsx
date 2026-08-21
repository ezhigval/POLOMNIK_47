import { contactEmail, contactPhone } from "@/lib/contact";
import { absoluteUrl, siteConfig } from "@/lib/site-config";

export function StructuredData() {
  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "TravelAgency",
    name: siteConfig.name,
    description: siteConfig.description,
    url: siteConfig.url,
    areaServed: siteConfig.region,
    email: contactEmail,
    telephone: contactPhone,
    sameAs: [],
    potentialAction: {
      "@type": "ReserveAction",
      target: absoluteUrl("/#tours"),
    },
  };

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
    />
  );
}
