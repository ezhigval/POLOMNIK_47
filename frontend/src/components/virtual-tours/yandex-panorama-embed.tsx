type YandexPanoramaEmbedProps = {
  src: string;
  title: string;
  className?: string;
};

export function YandexPanoramaEmbed({ src, title, className = "" }: YandexPanoramaEmbedProps) {
  return (
    <div className={`relative overflow-hidden rounded-[1.5rem] border border-stone-200 bg-stone-100 shadow-inner ${className}`}>
      <div className="pointer-events-none absolute inset-x-0 top-0 z-10 h-12 bg-gradient-to-b from-stone-900/10 to-transparent" />
      <iframe
        title={title}
        src={src}
        className="block h-[420px] w-full border-0 md:h-[500px]"
        loading="lazy"
        allow="fullscreen; accelerometer; gyroscope; picture-in-picture"
        allowFullScreen
      />
    </div>
  );
}
