import { cookies } from "next/headers";
import { getTour } from "@/lib/api/tours";
import type { Tour } from "@/lib/api/tours";
import { parseViewedTourIds, VIEWED_TOURS_COOKIE } from "@/lib/viewed-tours";

export async function getViewedTourIds(limit = 5): Promise<string[]> {
  const cookieStore = await cookies();
  return parseViewedTourIds(cookieStore.get(VIEWED_TOURS_COOKIE)?.value).slice(0, limit);
}

export async function loadViewedTours(limit = 5): Promise<Tour[]> {
  const ids = await getViewedTourIds(limit);
  if (ids.length === 0) {
    return [];
  }

  const tours = await Promise.all(
    ids.map(async (id) => {
      try {
        return await getTour(id);
      } catch {
        return null;
      }
    }),
  );

  return tours.filter((tour): tour is Tour => tour != null);
}
