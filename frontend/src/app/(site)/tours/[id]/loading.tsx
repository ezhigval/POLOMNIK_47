export default function TourLoading() {
  return (
    <div className="mx-auto max-w-5xl px-4 py-10">
      <div className="mb-6 h-4 w-32 animate-pulse rounded bg-stone-200" />
      <div className="grid gap-8 lg:grid-cols-[2fr_1fr]">
        <div className="h-80 animate-pulse rounded-2xl bg-stone-200" />
        <div className="h-96 animate-pulse rounded-2xl bg-stone-200" />
      </div>
    </div>
  );
}
