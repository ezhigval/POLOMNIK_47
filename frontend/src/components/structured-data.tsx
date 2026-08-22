import { contactEmail, contactPhone } from "@/lib/contact";
import { absoluteUrl, siteConfig } from "@/lib/site-config";
import type { Tour } from "@/lib/api/tours";
import { getSlotsAvailability } from "@/lib/format";

function JsonLd({ data }: { data: Record<string, unknown> }) {
  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(data) }}
    />
  );
}

export function StructuredData() {
  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "TravelAgency",
    "@id": `${siteConfig.url}/#organization`,
    name: siteConfig.name,
    legalName: siteConfig.fullName,
    alternateName: siteConfig.fullName,
    description: siteConfig.description,
    url: siteConfig.url,
    logo: absoluteUrl("/opengraph-image"),
    image: absoluteUrl("/opengraph-image"),
    areaServed: {
      "@type": "AdministrativeArea",
      name: "Россия",
    },
    email: contactEmail,
    telephone: contactPhone,
    sameAs: [siteConfig.parentOrganization.url],
    parentOrganization: {
      "@type": "Organization",
      name: siteConfig.parentOrganization.name,
      url: siteConfig.parentOrganization.url,
    },
    potentialAction: {
      "@type": "ReserveAction",
      target: absoluteUrl("/search"),
    },
  };

  return <JsonLd data={jsonLd} />;
}

export function TourStructuredData({ tour }: { tour: Tour }) {
  const availability = getSlotsAvailability(tour.slots_left);
  const offerAvailability =
    availability === "sold_out"
      ? "https://schema.org/SoldOut"
      : "https://schema.org/InStock";

  const images = tour.images?.filter(Boolean).slice(0, 5) ?? [];
  const jsonLd: Record<string, unknown> = {
    "@context": "https://schema.org",
    "@type": "TouristTrip",
    name: tour.title,
    description: tour.description?.split("\n")[0] || tour.title,
    url: absoluteUrl(`/tours/${tour.id}`),
    touristType: "Паломник",
    offers: {
      "@type": "Offer",
      price: tour.price,
      priceCurrency: tour.currency || "RUB",
      availability: offerAvailability,
      url: absoluteUrl(`/tours/${tour.id}`),
      ...(tour.date_start ? { validFrom: tour.date_start } : {}),
    },
  };
  if (images.length > 0) {
    jsonLd.image = images;
  }
  if (tour.location) {
    jsonLd.itinerary = { "@type": "Place", name: tour.location };
  }

  return <JsonLd data={jsonLd} />;
}
