import { NextRequest, NextResponse } from "next/server";
import { exportManagementBookingsCSV } from "@/lib/api/management";

export async function GET(request: NextRequest) {
  const params = request.nextUrl.searchParams;
  const upstream = await exportManagementBookingsCSV({
    status: params.get("status") ?? undefined,
    date_from: params.get("date_from") ?? undefined,
    date_to: params.get("date_to") ?? undefined,
    tour_id: params.get("tour_id") ?? undefined,
  }).catch(() => null);

  if (!upstream) {
    return NextResponse.json({ error: "API недоступен" }, { status: 503 });
  }
  if (upstream.status === 401) {
    return NextResponse.redirect(new URL("/management/login", request.url));
  }
  if (!upstream.ok) {
    const payload = await upstream.json().catch(() => null);
    return NextResponse.json(
      { error: payload?.error?.message ?? "Не удалось выгрузить заявки" },
      { status: upstream.status },
    );
  }

  const body = await upstream.arrayBuffer();
  return new NextResponse(body, {
    status: 200,
    headers: {
      "Content-Type": "text/csv; charset=utf-8",
      "Content-Disposition": `attachment; filename="bookings.csv"`,
    },
  });
}
