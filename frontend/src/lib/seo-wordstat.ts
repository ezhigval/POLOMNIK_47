/**
 * Wordstat cluster → live page mapping (v4 stage 10).
 * Frequencies and copy edits — only in owner cabinet / CMS (see docs/SEO_WORDSTAT.md).
 * Code reads meta from existing fields; does not invent keywords or counts.
 */

export type WordstatMetaSource =
  | "site_home_cms"
  | "site_config"
  | "search_page"
  | "news_index"
  | "news_article"
  | "tour_card";

export type WordstatCluster = {
  /** Stable id for owner checklist in SEO_WORDSTAT.md */
  id: string;
  /** Phrases to verify in Яндекс.Wordstat (owner fills frequency column in doc) */
  queries: string[];
  /** Existing URLs that should reflect the cluster when owner updates CMS/tour/news text */
  landing: Array<{
    path: string;
    metaSource: WordstatMetaSource;
    /** Optional slug for /news/{slug} or /tours/{slug} */
    slug?: string;
  }>;
};

/** Hypothesis clusters from SEO_WORDSTAT.md — not measured frequencies. */
export const wordstatClusters: WordstatCluster[] = [
  {
    id: "pilgrimage-general",
    queries: [
      "паломничество",
      "паломнические туры",
      "паломнические поездки",
      "паломнические туры из Санкт-Петербурга",
      "паломнические туры из Питера",
    ],
    landing: [
      { path: "/", metaSource: "site_home_cms" },
      { path: "/search", metaSource: "search_page" },
    ],
  },
  {
    id: "tikhvin",
    queries: [
      "Тихвин",
      "паломничество в Тихвин",
      "Тихвинский монастырь",
      "Тихвинская икона",
    ],
    landing: [
      { path: "/tours/tikhvin-path", metaSource: "tour_card", slug: "tikhvin-path" },
      { path: "/news/ikona-v-moskvu", metaSource: "news_article", slug: "ikona-v-moskvu" },
    ],
  },
  {
    id: "shrines-eparhia",
    queries: ["святыни", "святыни Тихвинской епархии"],
    landing: [
      {
        path: "/news/svyatyni-tikhvinskoy-eparhii",
        metaSource: "news_article",
        slug: "svyatyni-tikhvinskoy-eparhii",
      },
    ],
  },
  {
    id: "service-eparhia",
    queries: ["паломническая служба", "паломническая служба Тихвин", "епархия"],
    landing: [{ path: "/", metaSource: "site_config" }],
  },
];

/** Paths indexed for sitemap/SEO audit — no query-string filters. */
export function wordstatLandingPaths(): string[] {
  const paths = new Set<string>();
  for (const cluster of wordstatClusters) {
    for (const page of cluster.landing) {
      paths.add(page.path);
    }
  }
  return [...paths].sort();
}
