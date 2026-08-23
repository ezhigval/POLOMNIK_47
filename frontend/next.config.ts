import os from "node:os";
import type { NextConfig } from "next";

function lanDevHosts(): string[] {
  const hosts = new Set<string>();
  for (const addrs of Object.values(os.networkInterfaces())) {
    for (const addr of addrs ?? []) {
      const family = String(addr.family);
      if ((family === "IPv4" || family === "4") && !addr.internal) {
        hosts.add(addr.address);
      }
    }
  }
  return [...hosts];
}

function buildContentSecurityPolicy(): string {
  return [
    "default-src 'self'",
    "script-src 'self' 'unsafe-inline' 'unsafe-eval' https://mc.yandex.ru https://mc.yandex.com https://www.googletagmanager.com https://www.google-analytics.com",
    "style-src 'self' 'unsafe-inline'",
    "img-src 'self' data: blob: http://localhost:8080 http://127.0.0.1:8080 https://images.unsplash.com https://*.googleusercontent.com https://mc.yandex.ru https://mc.yandex.com https:",
    "font-src 'self' data:",
    [
      "connect-src 'self' ws: wss:",
      "http://localhost:8080",
      "http://127.0.0.1:8080",
      "https:",
      "https://mc.yandex.ru",
      "https://mc.yandex.com",
      "https://www.google-analytics.com",
      "https://region1.google-analytics.com",
      "https://www.googletagmanager.com",
    ].join(" "),
    "form-action 'self'",
    "frame-src 'self' https://mc.yandex.ru https://mc.yandex.com",
    "frame-ancestors 'self'",
    "base-uri 'self'",
  ].join("; ");
}

const securityHeaders = [
  { key: "X-DNS-Prefetch-Control", value: "on" },
  { key: "X-Frame-Options", value: "SAMEORIGIN" },
  { key: "X-Content-Type-Options", value: "nosniff" },
  { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
  {
    key: "Permissions-Policy",
    value: "camera=(), microphone=(), geolocation=(), interest-cohort=()",
  },
  {
    key: "Content-Security-Policy",
    value: buildContentSecurityPolicy(),
  },
];

const nextConfig: NextConfig = {
  output: "standalone",
  allowedDevOrigins: lanDevHosts(),
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "images.unsplash.com",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "api.tikhvin-palomnik.ru",
        pathname: "/**",
      },
    ],
  },
  async rewrites() {
    const internal = process.env.API_INTERNAL_URL?.replace(/\/$/, "");
    if (!internal) {
      return [];
    }
    return [
      {
        source: "/api/v1/:path*",
        destination: `${internal}/:path*`,
      },
    ];
  },
  async headers() {
    const headers = [...securityHeaders];
    if (process.env.NODE_ENV === "production") {
      headers.push({
        key: "Strict-Transport-Security",
        value: "max-age=63072000; includeSubDomains; preload",
      });
    }
    return [
      {
        source: "/(.*)",
        headers,
      },
    ];
  },
};

export default nextConfig;
