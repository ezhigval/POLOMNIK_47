import type { NextConfig } from "next";

const googleOAuthOrigins = [
  "https://accounts.google.com",
  "https://oauth2.googleapis.com",
  "https://www.googleapis.com",
];

function buildContentSecurityPolicy(): string {
  return [
    "default-src 'self'",
      "script-src 'self' 'unsafe-inline' 'unsafe-eval' https://mc.yandex.ru https://www.googletagmanager.com",
      "style-src 'self' 'unsafe-inline'",
      "img-src 'self' data: blob: http://localhost:8080 http://127.0.0.1:8080 https://images.unsplash.com https://*.googleusercontent.com https://mc.yandex.ru https:",
      "font-src 'self' data:",
      [
        "connect-src 'self'",
        "http://localhost:8080",
        "http://127.0.0.1:8080",
        "https:",
        "https://mc.yandex.ru",
        "https://www.google-analytics.com",
        "https://region1.google-analytics.com",
        ...googleOAuthOrigins,
      ].join(" "),
    `form-action 'self' ${googleOAuthOrigins[0]}`,
    `frame-src 'self' ${googleOAuthOrigins[0]}`,
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
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "images.unsplash.com",
        pathname: "/**",
      },
    ],
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
