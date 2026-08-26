/** Noscript pixel from the official Yandex tag (SSR HTML). JS tag stays in Analytics after cookie consent. */
export function YandexMetrikaPixel() {
  const ymId = process.env.NEXT_PUBLIC_YM_ID?.trim();
  if (!ymId) {
    return null;
  }

  return (
    <noscript>
      <div>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={`https://mc.yandex.ru/watch/${ymId}`}
          alt=""
          style={{ position: "absolute", left: "-9999px" }}
        />
      </div>
    </noscript>
  );
}
