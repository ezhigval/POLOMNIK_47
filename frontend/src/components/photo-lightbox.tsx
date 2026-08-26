"use client";

import { useCallback, useEffect } from "react";

type PhotoLightboxProps = {
  src: string;
  alt: string;
  open: boolean;
  onClose: () => void;
};

export function PhotoLightbox({ src, alt, open, onClose }: PhotoLightboxProps) {
  const onKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    },
    [onClose],
  );

  useEffect(() => {
    if (!open) {
      return;
    }
    document.addEventListener("keydown", onKeyDown);
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKeyDown);
      document.body.style.overflow = "";
    };
  }, [open, onKeyDown]);

  if (!open) {
    return null;
  }

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-black/90 p-4"
      role="dialog"
      aria-modal="true"
      onClick={onClose}
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={src}
        alt={alt}
        className="max-h-full max-w-full object-contain"
        onClick={(event) => event.stopPropagation()}
      />
    </div>
  );
}
