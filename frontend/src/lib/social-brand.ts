import "server-only";
import { readFile } from "node:fs/promises";
import path from "node:path";
import { defaultLogo } from "./brand";
import { fetchUploadedLogo, getSiteLogo } from "./site-logo";

export async function socialBrand() {
  const selected = await getSiteLogo();
  const upload = selected !== defaultLogo ? await fetchUploadedLogo(selected).catch(() => null) : null;
  const publicDir = path.join(process.cwd(), "public");
  const [logo, regular, bold] = await Promise.all([
    upload?.ok ? upload.arrayBuffer().then(bytes => Buffer.from(bytes)) : readFile(path.join(publicDir, defaultLogo)),
    readFile(path.join(publicDir, "fonts", "Arial.ttf")),
    readFile(path.join(publicDir, "fonts", "Arial-Bold.ttf")),
  ]);
  return {
    image: `data:image/png;base64,${logo.toString("base64")}`,
    fonts: [
      { name: "Arial", data: regular, weight: 400 as const },
      { name: "Arial", data: bold, weight: 700 as const },
    ],
  };
}
