"use client";

import Link from "next/link";
import { useEffect, useRef } from "react";
import styles from "./home-spaces.module.css";

const modes = [
  { href: "/", label: "Путеводитель", guide: true },
  { href: "/?mode=trips", label: "Подобрать поездку", guide: false },
];

export function HomeSpaceSwitcher({ guide }: { guide: boolean }) {
  const slot = useRef<HTMLDivElement>(null);
  const nav = useRef<HTMLElement>(null);

  useEffect(() => {
    const anchor = slot.current;
    const control = nav.current;
    if (!anchor || !control) return;

    let motion: Animation | undefined;
    const header = document.querySelector("header");
    let observer: IntersectionObserver;

    function observe() {
      if (!anchor || !control) return;
      observer?.disconnect();
      const headerHeight = header?.getBoundingClientRect().height ?? 0;
      observer = new IntersectionObserver(([entry]) => {
        const docked = entry.boundingClientRect.bottom <= headerHeight;
        if (control.hasAttribute("data-docked") === docked) return;
        const from = control.getBoundingClientRect();
        motion?.cancel();
        control.toggleAttribute("data-docked", docked);
        if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
        const to = control.getBoundingClientRect();
        if (from.bottom < 0 || !to.width || !to.height) return;
        motion = control.animate([
          { transform: `translate(${from.left - to.left}px, ${from.top - to.top}px) scale(${from.width / to.width}, ${from.height / to.height})` },
          { transform: "none" },
        ], { duration: 320, easing: "cubic-bezier(.22,.75,.24,1)" });
      }, { rootMargin: `-${headerHeight}px 0px 0px` });
      observer.observe(anchor);
    }

    observe();
    const size = new ResizeObserver(observe);
    if (header) size.observe(header);
    return () => {
      size.disconnect();
      observer.disconnect();
      motion?.cancel();
    };
  }, []);

  return (
    <div className={styles.slot} ref={slot}>
      <nav ref={nav} className={styles.switcher} aria-label="Режим главной страницы" data-guide={guide}>
        <span aria-hidden className={styles.slider} />
        {modes.map((mode) => (
          <Link key={mode.href} href={mode.href} aria-current={mode.guide === guide ? "page" : undefined}
            onNavigate={() => window.scrollTo({top: 0, behavior: "instant"})}>
            <span className={mode.guide ? styles.book : styles.ticket} aria-hidden />
            {mode.label}
          </Link>
        ))}
      </nav>
    </div>
  );
}
