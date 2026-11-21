import { logoSrc } from "@/lib/brand";
import { getSiteLogo } from "@/lib/site-logo";

export async function GET() {
  return new Response(null, { status: 307, headers: { Location: logoSrc(await getSiteLogo()), "Cache-Control": "no-store" } });
}
