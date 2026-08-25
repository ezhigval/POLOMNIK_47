import { ImageResponse } from "next/og";
import { siteConfig } from "@/lib/site-config";

export const runtime = "edge";
export const alt = `${siteConfig.name} — ${siteConfig.tagline}`;
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default function OpenGraphImage() {
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
          fontFamily: "Georgia, 'Times New Roman', serif",
        }}
      >
        <div
          style={{
            display: "flex",
            fontSize: 28,
            letterSpacing: "0.18em",
            textTransform: "uppercase",
            color: "#b8d0e6",
          }}
        >
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
              fontFamily: "system-ui, sans-serif",
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
            fontFamily: "system-ui, sans-serif",
          }}
        >
          tikhvin-palomnik.ru
        </div>
      </div>
    ),
    { ...size },
  );
}
