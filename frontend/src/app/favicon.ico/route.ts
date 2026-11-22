import { logoSrc } from "@/lib/brand";
import { getSiteLogo } from "@/lib/site-logo";

// Resolve the logo at request time, not during the production image build.
export const dynamic = "force-dynamic";

export async function GET() {
  return new Response(null, { status: 307, headers: { Location: logoSrc(await getSiteLogo()), "Cache-Control": "no-store" } });
}
