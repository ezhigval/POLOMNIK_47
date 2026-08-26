import { contactEmail, contactPhone } from "@/lib/contact";
import { absoluteUrl, siteConfig } from "@/lib/site-config";
import type { Tour } from "@/lib/api/tours";
import { canBookTour } from "@/lib/api/tours";
import { tourPath, tourSeoDescription } from "@/lib/tour-path";
import type { NewsArticle } from "@/lib/news";

export function JsonLd({ data }: { data: Record<string, unknown> }) {
  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(data) }}
    />
  );
}

export function StructuredData() {
  const organizationId = `${siteConfig.url}/#organization`;
  const websiteId = `${siteConfig.url}/#website`;

  const jsonLd = {
    "@context": "https://schema.org",
    "@graph": [
      {
        "@type": "TravelAgency",
        "@id": organizationId,
        name: siteConfig.name,
        legalName: siteConfig.fullName,
        alternateName: siteConfig.fullName,
        description: siteConfig.description,
        url: siteConfig.url,
        logo: absoluteUrl("/opengraph-image"),
        image: absoluteUrl("/opengraph-image"),
        inLanguage: "ru-RU",
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
      },
      {
        "@type": "WebSite",
        "@id": websiteId,
        url: siteConfig.url,
        name: siteConfig.name,
        inLanguage: "ru-RU",
        publisher: { "@id": organizationId },
        potentialAction: {
          "@type": "SearchAction",
          target: {
            "@type": "EntryPoint",
            urlTemplate: `${absoluteUrl("/search")}?q={search_term_string}`,
          },
          "query-input": "required name=search_term_string",
        },
      },
    ],
  };

  return <JsonLd data={jsonLd} />;
}

export function TourStructuredData({ tour }: { tour: Tour }) {
  const bookingClosed = !canBookTour(tour);
  const offerAvailability = bookingClosed
    ? "https://schema.org/SoldOut"
    : "https://schema.org/InStock";
  const url = absoluteUrl(tourPath(tour));

  const images = tour.images?.filter(Boolean).slice(0, 5) ?? [];
  const regular = Boolean(tour.is_regular);
  const jsonLd: Record<string, unknown> = {
    "@context": "https://schema.org",
    "@type": "TouristTrip",
    name: tour.title,
    description: tourSeoDescription(tour),
    url,
    touristType: "Паломник",
    provider: {
      "@id": `${siteConfig.url}/#organization`,
    },
  };
  if (!regular && tour.price != null) {
    jsonLd.offers = {
      "@type": "Offer",
      price: tour.price,
      priceCurrency: tour.currency || "RUB",
      availability: offerAvailability,
      url,
      ...(tour.date_start ? { validFrom: tour.date_start } : {}),
    };
  }
  if (images.length > 0) {
    jsonLd.image = images;
  }
  if (tour.location) {
    jsonLd.itinerary = { "@type": "Place", name: tour.location };
  }

  return <JsonLd data={jsonLd} />;
}

export function NewsArticleStructuredData({ article }: { article: NewsArticle }) {
  const url = absoluteUrl(`/news/${article.slug}`);
  const jsonLd: Record<string, unknown> = {
    "@context": "https://schema.org",
    "@type": "NewsArticle",
    headline: article.title,
    description: article.excerpt,
    datePublished: article.date,
    url,
    inLanguage: "ru-RU",
    mainEntityOfPage: url,
    author: {
      "@id": `${siteConfig.url}/#organization`,
    },
    publisher: {
      "@id": `${siteConfig.url}/#organization`,
    },
  };
  if (article.image) {
    jsonLd.image = [absoluteUrl(article.image)];
  }
  return <JsonLd data={jsonLd} />;
}
