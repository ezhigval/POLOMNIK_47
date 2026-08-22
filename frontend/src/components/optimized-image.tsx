import Image, { type ImageProps } from "next/image";

const OPTIMIZED_HOSTS = new Set(["images.unsplash.com", "api.tikhvin-palomnik.ru"]);

export function isOptimizableImageUrl(url: string): boolean {
  try {
    return OPTIMIZED_HOSTS.has(new URL(url).hostname);
  } catch {
    return false;
  }
}

type OptimizedImageProps = Omit<ImageProps, "src"> & {
  src: string;
  fallbackClassName?: string;
};

export function OptimizedImage({
  src,
  alt,
  className,
  fallbackClassName,
  priority,
  fill,
  sizes,
  width,
  height,
}: OptimizedImageProps) {
  if (!isOptimizableImageUrl(src)) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={src}
        alt={alt}
        className={fallbackClassName ?? className}
        loading={priority ? "eager" : "lazy"}
      />
    );
  }

  if (fill) {
    return (
      <Image
        src={src}
        alt={alt}
        fill
        priority={priority}
        sizes={sizes}
        className={className}
      />
    );
  }

  return (
    <Image
      src={src}
      alt={alt}
      width={width ?? 1200}
      height={height ?? 800}
      priority={priority}
      sizes={sizes}
      className={className}
    />
  );
}
