"use client";

import Image from "next/image";
import { useEffect, useState } from "react";
import { heroBackgroundImages } from "@/lib/tour-cover";
import { siteConfig } from "@/lib/site-config";

const INTERVAL_MS = 7000;

export function HeroBackground() {
  const [index, setIndex] = useState(0);

  useEffect(() => {
    const media = window.matchMedia("(prefers-reduced-motion: reduce)");
    if (media.matches) {
      return;
    }

    const id = window.setInterval(() => {
      setIndex((current) => (current + 1) % heroBackgroundImages.length);
    }, INTERVAL_MS);

    return () => window.clearInterval(id);
  }, []);

  return (
    <div className="absolute inset-0">
      {heroBackgroundImages.map((src, i) => (
        <Image
          key={src}
          src={src}
          alt={i === 0 ? `${siteConfig.name} — паломнические поездки` : ""}
          fill
          priority={i === 0}
          sizes="100vw"
          className={`object-cover transition-opacity duration-1000 ease-in-out ${
            i === index ? "opacity-100" : "opacity-0"
          }`}
        />
      ))}
      <div className="absolute inset-0 bg-gradient-to-r from-brand-950/78 via-brand-900/48 to-brand-800/28" />
      <div className="absolute inset-0 bg-gradient-to-t from-brand-950/45 via-transparent to-brand-950/15" />
    </div>
  );
}
