import { ImageResponse } from "next/og";
import { ApiError } from "@/lib/api/client";
import { getCachedTour } from "@/lib/api/tour-page";
import { formatDateRange, formatPrice } from "@/lib/format";
import { tourShowsPrice } from "@/lib/api/tours";
import { siteConfig } from "@/lib/site-config";

export const size = { width: 1200, height: 630 };
export const contentType = "image/png";
export const runtime = "nodejs";

type OgImageProps = {
  params: Promise<{ id: string }>;
};

export default async function TourOpenGraphImage({ params }: OgImageProps) {
  const { id } = await params;

  let title = "Паломнический тур";
  let subtitle = siteConfig.name;
  let meta = siteConfig.tagline;

  try {
    const tour = await getCachedTour(id);
    title = tour.title;
    subtitle = tour.location;
    const metaParts: string[] = [];
    if (tour.is_regular) {
      metaParts.push("Регулярный тур");
    } else {
      metaParts.push(formatDateRange(tour.date_start, tour.date_end));
    }
    if (tourShowsPrice(tour)) {
      metaParts.push(formatPrice(tour.price, tour.currency));
    }
    meta = metaParts.filter(Boolean).join(" · ");
  } catch (error) {
    if (!(error instanceof ApiError && error.status === 404)) {
      throw error;
    }
  }

  return new ImageResponse(
    (
      <div
        style={{
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          justifyContent: "flex-end",
          padding: 64,
          background: "linear-gradient(135deg, #132c3e 0%, #2a5673 45%, #3a6e91 100%)",
          color: "white",
        }}
      >
        <div
          style={{
            fontSize: 22,
            lineHeight: 1.25,
            opacity: 0.85,
            marginBottom: 16,
            maxWidth: 980,
          }}
        >
          {siteConfig.name}
        </div>
        <div style={{ fontSize: 64, fontWeight: 600, lineHeight: 1.05, maxWidth: 980 }}>
          {title}
        </div>
        <div style={{ fontSize: 30, marginTop: 20, opacity: 0.92 }}>{subtitle}</div>
        <div style={{ fontSize: 26, marginTop: 12, opacity: 0.85 }}>{meta}</div>
      </div>
    ),
    { ...size },
  );
}
