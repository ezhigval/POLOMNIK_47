import Link from "next/link";

export default function NotFound() {
  return (
    <div className="mx-auto flex min-h-[60vh] max-w-lg flex-col items-center justify-center px-4 py-16 text-center">
      <p className="text-sm font-medium uppercase tracking-widest text-brand-700">404</p>
      <h1 className="mt-3 font-display text-3xl font-semibold text-stone-900">Страница не найдена</h1>
      <p className="mt-3 text-stone-600">
        Возможно, тур уже завершился или ссылка устарела. Посмотрите актуальный каталог.
      </p>
      <div className="mt-8 flex flex-wrap justify-center gap-3">
        <Link href="/" className="btn-primary">
          На главную
        </Link>
        <Link href="/search" className="btn-secondary">
          Туры
        </Link>
        <Link href="/support" className="btn-secondary">
          Поддержка
        </Link>
      </div>
    </div>
  );
}
