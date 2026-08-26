"use client";

import Script from "next/script";
import { useEffect, useState } from "react";
import { analyticsConfig, isAnalyticsEnabled } from "@/lib/analytics";
import { allowsAnalyticsCookies, getCookieConsent, type CookieConsentChoice } from "@/lib/cookie-consent";

export function Analytics() {
  const [choice, setChoice] = useState<CookieConsentChoice | null>(null);

  useEffect(() => {
    function onChange(event: Event) {
      const detail = (event as CustomEvent<CookieConsentChoice>).detail;
      setChoice(detail ?? getCookieConsent());
    }
    void Promise.resolve().then(() => setChoice(getCookieConsent()));
    window.addEventListener("palomnik:cookie-consent-changed", onChange);
    return () => window.removeEventListener("palomnik:cookie-consent-changed", onChange);
  }, []);

  if (!isAnalyticsEnabled() || !allowsAnalyticsCookies(choice)) {
    return null;
  }

  const ymId = analyticsConfig.yandexMetrikaId;
  const gaId = analyticsConfig.googleAnalyticsId;
  const webvisor = analyticsConfig.yandexWebvisor;
  const clickmap = analyticsConfig.yandexClickmap;

  return (
    <>
      {ymId ? (
        <>
          <Script id="yandex-metrika" strategy="afterInteractive">
            {`
              (function(m,e,t,r,i,k,a){
                m[i]=m[i]||function(){(m[i].a=m[i].a||[]).push(arguments)};
                m[i].l=1*new Date();
                for (var j = 0; j < document.scripts.length; j++) {if (document.scripts[j].src === r) { return; }}
                k=e.createElement(t),a=e.getElementsByTagName(t)[0],k.async=1,k.src=r,a.parentNode.insertBefore(k,a)
              })(window, document,'script','https://mc.yandex.ru/metrika/tag.js?id=${ymId}', 'ym');
              ym(${ymId}, 'init', {
                ssr:true,
                webvisor:${webvisor ? "true" : "false"},
                clickmap:${clickmap ? "true" : "false"},
                ecommerce:"dataLayer",
                accurateTrackBounce:true,
                trackLinks:true
              });
            `}
          </Script>
        </>
      ) : null}

      {gaId ? (
        <>
          <Script src={`https://www.googletagmanager.com/gtag/js?id=${gaId}`} strategy="afterInteractive" />
          <Script id="google-analytics" strategy="afterInteractive">
            {`
              window.dataLayer = window.dataLayer || [];
              function gtag(){dataLayer.push(arguments);}
              gtag('js', new Date());
              gtag('config', '${gaId}');
            `}
          </Script>
        </>
      ) : null}
    </>
  );
}
