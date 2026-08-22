import { ImageResponse } from "next/og";
import { siteConfig } from "@/lib/site-config";

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
          justifyContent: "flex-end",
          padding: 64,
          background: "linear-gradient(135deg, #042f2e 0%, #115e59 45%, #0f766e 100%)",
          color: "white",
        }}
      >
        <div
          style={{
            fontSize: 28,
            letterSpacing: 4,
            textTransform: "uppercase",
            opacity: 0.85,
            marginBottom: 16,
          }}
        >
          Паломническая служба
        </div>
        <div style={{ fontSize: 72, fontWeight: 600, lineHeight: 1.05, maxWidth: 900 }}>
          {siteConfig.name}
        </div>
        <div style={{ fontSize: 32, marginTop: 24, opacity: 0.92, maxWidth: 800 }}>
          {siteConfig.tagline}
        </div>
      </div>
    ),
    { ...size },
  );
}
