import { isUploadedLogo } from "@/lib/brand";
import { fetchUploadedLogo } from "@/lib/site-logo";

export async function GET(request: Request) {
  const file = new URL(request.url).searchParams.get("file") ?? "";
  if (!isUploadedLogo(file)) return new Response(null, { status: 404 });
  const image = await fetchUploadedLogo(file);
  if (!image.ok) return new Response(null, { status: image.status });
  return new Response(image.body, { headers: {
    "Content-Type": "image/png",
    "Cache-Control": "public, max-age=31536000, immutable",
    "X-Content-Type-Options": "nosniff",
  } });
}
