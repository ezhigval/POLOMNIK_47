import { CreateReviewForm } from "@/components/management/create-review-form";
import { ManagementPanel } from "@/components/management/management-panel";
import { StatusBadge } from "@/components/management/status-badge";
import {
  approveReviewAction,
  deleteReviewAction,
  rejectReviewAction,
} from "@/app/management/actions";
import { listManagementReviews, listManagementTours } from "@/lib/api/management";
import { buildTourTitleMap, tourTitle } from "@/lib/tour-title-map";

export default async function ManagementReviewsPage() {
  const [reviews, tours] = await Promise.all([listManagementReviews(), listManagementTours()]);
  const tourNames = buildTourTitleMap(tours);

  return (
    <div className="grid gap-8 lg:grid-cols-[1.2fr_1fr]">
      <section className="space-y-4">
        {reviews.length === 0 ? (
          <ManagementPanel>
            <div className="px-5 py-12 text-center text-stone-500">Отзывов пока нет.</div>
          </ManagementPanel>
        ) : (
          reviews.map((review) => (
            <article key={review.id} className="rounded-2xl border border-stone-200 bg-white p-5">
              <div className="mb-3 flex flex-wrap items-start justify-between gap-3">
                <div>
                  <h2 className="font-semibold text-stone-900">{review.client_name}</h2>
                  <p className="text-sm text-stone-500">Тур: {tourTitle(tourNames, review.tour_id)}</p>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-amber-600" aria-label={`Рейтинг ${review.rating}`}>
                    {"★".repeat(review.rating)}
                    <span className="text-stone-300">{"★".repeat(5 - review.rating)}</span>
                  </span>
                  <StatusBadge variant={review.is_approved ? "success" : "warning"}>
                    {review.is_approved ? "Одобрен" : "На модерации"}
                  </StatusBadge>
                </div>
              </div>
              <p className="mb-4 text-sm leading-6 text-stone-700">{review.text}</p>
              <div className="flex flex-wrap gap-2">
                {!review.is_approved ? (
                  <form action={approveReviewAction}>
                    <input type="hidden" name="id" value={review.id} />
                    <button type="submit" className="btn-primary">
                      Одобрить
                    </button>
                  </form>
                ) : null}
                <form action={rejectReviewAction}>
                  <input type="hidden" name="id" value={review.id} />
                  <button type="submit" className="btn-secondary">
                    Отклонить
                  </button>
                </form>
                <form action={deleteReviewAction}>
                  <input type="hidden" name="id" value={review.id} />
                  <button type="submit" className="btn-danger">
                    Удалить
                  </button>
                </form>
              </div>
            </article>
          ))
        )}
      </section>

      <CreateReviewForm tours={tours} />
    </div>
  );
}
