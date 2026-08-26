import { cookies } from "next/headers";
import { NextResponse } from "next/server";
import {
  addViewedTourId,
  parseViewedTourIds,
  VIEWED_TOURS_COOKIE,
  viewedToursCookieOptions,
} from "@/lib/viewed-tours";

export async function GET() {
  const cookieStore = await cookies();
  const ids = parseViewedTourIds(cookieStore.get(VIEWED_TOURS_COOKIE)?.value);
  return NextResponse.json({ data: ids });
}

export async function POST(request: Request) {
  let tourId = "";
  try {
    const body = await request.json();
    tourId = String(body.tour_id ?? "").trim();
  } catch {
    return NextResponse.json({ error: "Invalid body" }, { status: 400 });
  }

  if (!tourId) {
    return NextResponse.json({ error: "tour_id required" }, { status: 400 });
  }

  const cookieStore = await cookies();
  const current = parseViewedTourIds(cookieStore.get(VIEWED_TOURS_COOKIE)?.value);
  const next = addViewedTourId(current, tourId);

  const response = NextResponse.json({ data: next });
  response.cookies.set(VIEWED_TOURS_COOKIE, JSON.stringify(next), viewedToursCookieOptions());
  return response;
}
