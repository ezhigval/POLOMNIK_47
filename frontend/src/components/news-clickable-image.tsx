"use client";

import { useState } from "react";
import { PhotoLightbox } from "@/components/photo-lightbox";

type NewsClickableImageProps = {
  src: string;
  alt: string;
  className?: string;
};

export function NewsClickableImage({ src, alt, className }: NewsClickableImageProps) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <button type="button" onClick={() => setOpen(true)} className="block w-full overflow-hidden bg-stone-100 text-left">
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img src={src} alt={alt} className={className} loading="lazy" />
      </button>
      <PhotoLightbox src={src} alt={alt} open={open} onClose={() => setOpen(false)} />
    </>
  );
}
