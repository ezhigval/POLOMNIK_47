import Link from "next/link";

export default function TourNotFound() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-20 text-center">
      <h1 className="mb-3 text-2xl font-semibold">Тур не найден</h1>
      <p className="mb-6 text-stone-600">Возможно, тур уже недоступен или ссылка устарела.</p>
      <div className="flex flex-wrap justify-center gap-3">
        <Link href="/search" className="btn-primary">
          Туры
        </Link>
        <Link href="/" className="btn-secondary">
          На главную
        </Link>
      </div>
    </div>
  );
}
