type NewsPhotoStripProps = {
  srcs: string[];
};

export function NewsPhotoStrip({ srcs }: NewsPhotoStripProps) {
  return (
    <div className="overflow-hidden rounded-3xl border border-stone-200 bg-white shadow-sm">
      <div className="flex flex-col">
        {srcs.map((src, index) => (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            key={src}
            src={src}
            alt=""
            className="block h-auto w-full"
            loading={index === 0 ? "eager" : "lazy"}
          />
        ))}
      </div>
    </div>
  );
}
