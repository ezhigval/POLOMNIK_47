"use client";

import { useState } from "react";
import { PhotoLightbox } from "@/components/photo-lightbox";

type NewsPhotoStripProps = {
  srcs: string[];
};

export function NewsPhotoStrip({ srcs }: NewsPhotoStripProps) {
  const [activeSrc, setActiveSrc] = useState<string | null>(null);

  return (
    <>
      <div className="overflow-hidden rounded-3xl border border-stone-200 bg-white shadow-sm">
        <div className="flex flex-col">
          {srcs.map((src, index) => (
            <button
              key={src}
              type="button"
              onClick={() => setActiveSrc(src)}
              className="block w-full bg-stone-100 text-left"
            >
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src={src} alt="" className="block h-auto w-full" loading={index === 0 ? "eager" : "lazy"} />
            </button>
          ))}
        </div>
      </div>
      {activeSrc ? (
        <PhotoLightbox src={activeSrc} alt="" open={Boolean(activeSrc)} onClose={() => setActiveSrc(null)} />
      ) : null}
    </>
  );
}
