import type { Metadata } from "next";
import { notFound, redirect } from "next/navigation";
import { CmsPageRenderer } from "@/components/cms/cms-page-renderer";
import { getPublishedPage } from "@/lib/api/cms";
import { siteConfig } from "@/lib/site-config";

type PageProps = {
  params: Promise<{ slug: string }>;
};

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug } = await params;
  if (slug === "home") {
    return {};
  }
  const page = await getPublishedPage(slug);
  if (!page) {
    return { title: "Страница" };
  }
  const title = page.meta_title?.trim() || page.title;
  const description = page.meta_description?.trim() || siteConfig.description;
  const canonical = page.path?.startsWith("/") && page.path !== "/" ? page.path : `/pages/${page.slug}`;
  return {
    title,
    description,
    alternates: { canonical },
    openGraph: {
      title,
      description,
      url: canonical,
    },
  };
}

export default async function CmsPublicPage({ params }: PageProps) {
  const { slug } = await params;

  if (slug === "home") {
    redirect("/");
  }

  const page = await getPublishedPage(slug);
  if (!page || page.blocks.length === 0) {
    notFound();
  }

  return <CmsPageRenderer blocks={page.blocks} />;
}
