import { notFound, redirect } from "next/navigation";
import { CmsPageRenderer } from "@/components/cms/cms-page-renderer";
import { getPublishedPage } from "@/lib/api/cms";

type PageProps = {
  params: Promise<{ slug: string }>;
};

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
