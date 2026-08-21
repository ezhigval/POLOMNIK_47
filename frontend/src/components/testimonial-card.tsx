export type Testimonial = {
  client_name: string;
  text: string;
  rating: number;
  tour_title: string;
};

export function StarRating({ rating }: { rating: number }) {
  return (
    <span className="text-amber-500" aria-label={`Оценка ${rating} из 5`}>
      {"★".repeat(rating)}
      <span className="text-stone-200">{"★".repeat(5 - rating)}</span>
    </span>
  );
}

export function TestimonialCard({ item }: { item: Testimonial }) {
  return (
    <blockquote className="flex flex-col rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
      <StarRating rating={item.rating} />
      <p className="my-4 flex-1 text-sm leading-7 text-stone-700">&ldquo;{item.text}&rdquo;</p>
      <footer>
        <cite className="not-italic">
          <span className="font-medium text-stone-900">{item.client_name}</span>
          <span className="mt-0.5 block text-xs text-stone-500">{item.tour_title}</span>
        </cite>
      </footer>
    </blockquote>
  );
}
