const UUID_RE =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

export function isUuidParam(value: string): boolean {
  return UUID_RE.test(value.trim());
}

export function tourPath(tour: { id: string; slug?: string }): string {
  const key = tour.slug?.trim() || tour.id;
  return `/tours/${encodeURIComponent(key)}`;
}

export function tourSeoTitle(tour: { title: string; location?: string }): string {
  const location = tour.location?.trim();
  if (location) {
    return `${tour.title} — паломничество, ${location}`;
  }
  return tour.title;
}

export function tourSeoDescription(tour: {
  title: string;
  description?: string;
  location?: string;
  is_regular?: boolean;
}): string {
  const firstLine = tour.description
    ?.split("\n")
    .map((line) => line.trim())
    .find(Boolean);
  if (firstLine) {
    return firstLine.length > 180 ? `${firstLine.slice(0, 177)}…` : firstLine;
  }
  if (tour.location?.trim()) {
    if (tour.is_regular) {
      return `Паломнический тур «${tour.title}» — ${tour.location}. Запись на сайте.`;
    }
    return `Паломнический тур «${tour.title}» — ${tour.location}. Даты и стоимость на сайте.`;
  }
  if (tour.is_regular) {
    return `Паломнический тур «${tour.title}». Запись на сайте.`;
  }
  return `Паломнический тур «${tour.title}». Даты и стоимость на сайте.`;
}
