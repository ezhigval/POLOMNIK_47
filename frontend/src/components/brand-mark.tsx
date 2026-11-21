import Image from "next/image";
import { logoSrc } from "@/lib/brand";
import { getSiteLogo } from "@/lib/site-logo";

export async function BrandMark({ className = "size-14 sm:size-16" }: { className?: string }) {
  return <Image unoptimized src={logoSrc(await getSiteLogo())} alt="" width={64} height={64}
    className={`shrink-0 rounded-lg bg-white object-contain ${className}`} />;
}
