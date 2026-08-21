export type TripSearchValues = {
  from?: string;
  destination?: string;
  q?: string;
  date_from?: string;
  date_to?: string;
  people?: string;
  page?: string;
};

export type TourFilterValues = TripSearchValues & {
  price_min?: string;
  price_max?: string;
  location?: string;
  is_hot?: string;
  min_slots?: string;
};

export function parseTourFilters(
  searchParams: Record<string, string | string[] | undefined>,
): TourFilterValues {
  const value = (key: keyof TourFilterValues) => {
    const raw = searchParams[key];
    if (Array.isArray(raw)) {
      return raw[0] ?? "";
    }
    return raw ?? "";
  };

  const q = value("q") || destinationToQuery(value("destination"));

  return {
    from: value("from"),
    destination: value("destination"),
    q,
    date_from: value("date_from"),
    date_to: value("date_to"),
    people: value("people"),
    price_min: value("price_min"),
    price_max: value("price_max"),
    location: value("location"),
    is_hot: value("is_hot"),
    min_slots: value("min_slots") || value("people"),
    page: value("page"),
  };
}

export function toTourQueryParams(filters: TourFilterValues): Record<string, string> {
  const params: Record<string, string> = {};
  const apiKeys: (keyof TourFilterValues)[] = [
    "q",
    "date_from",
    "date_to",
    "price_min",
    "price_max",
    "location",
    "is_hot",
    "min_slots",
    "page",
  ];

  for (const key of apiKeys) {
    const value = filters[key];
    if (value) {
      params[key] = value;
    }
  }

  return params;
}

export function toSearchParams(filters: TripSearchValues): URLSearchParams {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value && key !== "page") {
      params.set(key, value);
    }
  }
  return params;
}

export function hasActiveFilters(filters: TourFilterValues): boolean {
  return Object.entries(filters).some(
    ([key, value]) => !["page", "from", "destination", "people"].includes(key) && Boolean(value),
  );
}

function destinationToQuery(destination: string): string {
  if (!destination) {
    return "";
  }
  const map: Record<string, string> = {
    optina: "Оптина",
    diveevo: "Дивеево",
    valaam: "Валаам",
    solovki: "Солов",
  };
  return map[destination] ?? destination;
}
