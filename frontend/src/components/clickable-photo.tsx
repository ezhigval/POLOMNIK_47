"use client";

import { useState } from "react";
import { PhotoLightbox } from "@/components/photo-lightbox";

type ClickablePhotoProps = {
  src: string;
  alt?: string;
  className?: string;
  loading?: "eager" | "lazy";
};

export function ClickablePhoto({ src, alt = "", className, loading = "lazy" }: ClickablePhotoProps) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <button type="button" className="block w-full cursor-zoom-in" onClick={() => setOpen(true)}>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img src={src} alt={alt} className={className} loading={loading} />
      </button>
      <PhotoLightbox src={src} alt={alt} open={open} onClose={() => setOpen(false)} />
    </>
  );
}
