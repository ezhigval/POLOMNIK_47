import { ImageResponse } from "next/og";
import { siteConfig } from "@/lib/site-config";

import { socialBrand } from "@/lib/social-brand";

export const dynamic = "force-dynamic";

export const alt = `${siteConfig.name} — ${siteConfig.tagline}`;
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default async function OpenGraphImage() {
  const brand = await socialBrand();

  return new ImageResponse(
    (
      <div
        style={{
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
          padding: "64px",
          background: "linear-gradient(135deg, #132c3e 0%, #21475f 48%, #2a5673 100%)",
          color: "#f8fafc",
          fontFamily: "Arial",
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            fontSize: 28,
            letterSpacing: "0.18em",
            textTransform: "uppercase",
            color: "#b8d0e6",
          }}
        >
          <img src={brand.image} width={110} height={110} alt="" style={{ marginRight: 24, borderRadius: 10, backgroundColor: "white" }} />
          {siteConfig.tagline}
        </div>
        <div style={{ display: "flex", flexDirection: "column", gap: 20 }}>
          <div style={{ display: "flex", fontSize: 44, fontWeight: 700, lineHeight: 1.15, maxWidth: 1040 }}>
            {siteConfig.name}
          </div>
          <div
            style={{
              display: "flex",
              maxWidth: 880,
              fontSize: 28,
              lineHeight: 1.35,
              color: "#dde9f4",
              fontFamily: "Arial",
            }}
          >
            {siteConfig.description}
          </div>
        </div>
        <div
          style={{
            display: "flex",
            fontSize: 24,
            color: "#b8d0e6",
            fontFamily: "Arial",
          }}
        >
          tikhvin-palomnik.ru
        </div>
      </div>
    ),
    { ...size, fonts: brand.fonts },
  );
}
