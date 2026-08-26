export const VIEWED_TOURS_COOKIE = "palomnik_viewed_tours";
export const VIEWED_TOURS_MAX = 20;

export function parseViewedTourIds(raw: string | undefined | null): string[] {
  if (!raw) {
    return [];
  }
  try {
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) {
      return [];
    }
    return parsed
      .filter((id): id is string => typeof id === "string" && id.length > 0)
      .slice(0, VIEWED_TOURS_MAX);
  } catch {
    return [];
  }
}

export function addViewedTourId(ids: string[], tourId: string): string[] {
  const trimmed = tourId.trim();
  if (!trimmed) {
    return ids;
  }
  const without = ids.filter((id) => id !== trimmed);
  return [trimmed, ...without].slice(0, VIEWED_TOURS_MAX);
}

export function viewedToursCookieOptions() {
  return {
    httpOnly: true,
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 60 * 60 * 24 * 7,
  };
}
