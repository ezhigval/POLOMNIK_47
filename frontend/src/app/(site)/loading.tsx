export default function Loading() {
  return (
    <div className="mx-auto max-w-6xl px-4 py-8">
      <div className="h-[420px] animate-pulse rounded-3xl bg-stone-200" />
      <div className="mt-12 grid gap-4 sm:grid-cols-2">
        {Array.from({ length: 4 }).map((_, index) => (
          <div key={index} className="h-40 animate-pulse rounded-2xl bg-stone-200" />
        ))}
      </div>
      <div className="mt-16 h-64 animate-pulse rounded-2xl bg-stone-200" />
    </div>
  );
}
