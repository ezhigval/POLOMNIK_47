import Link from "next/link";
import type { TourFilterValues } from "@/lib/tour-filters";

type PaginationProps = {
  meta: {
    page: number;
    limit: number;
    total: number;
    has_next: boolean;
  };
  filters: TourFilterValues;
  basePath?: string;
};

function buildPageHref(filters: TourFilterValues, page: number, basePath: string) {
  const params = new URLSearchParams();

  for (const [key, value] of Object.entries(filters)) {
    if (value) {
      params.set(key, value);
    }
  }

  if (page > 1) {
    params.set("page", String(page));
  } else {
    params.delete("page");
  }

  const query = params.toString();
  return query ? `${basePath}?${query}` : basePath;
}

export function Pagination({ meta, filters, basePath = "/search" }: PaginationProps) {
  if (meta.total <= meta.limit) {
    return null;
  }

  const prevPage = meta.page > 1 ? meta.page - 1 : null;
  const nextPage = meta.has_next ? meta.page + 1 : null;

  return (
    <nav className="mt-6 flex items-center justify-between gap-4" aria-label="Пагинация туров">
      <span className="text-sm text-stone-500">
        Страница {meta.page}, всего {meta.total}
      </span>
      <div className="flex gap-3">
        {prevPage ? (
          <Link
            href={buildPageHref(filters, prevPage, basePath)}
            className="btn-secondary px-4 py-2 text-sm"
          >
            Назад
          </Link>
        ) : null}
        {nextPage ? (
          <Link
            href={buildPageHref(filters, nextPage, basePath)}
            className="btn-primary px-4 py-2 text-sm"
          >
            Дальше
          </Link>
        ) : null}
      </div>
    </nav>
  );
}
