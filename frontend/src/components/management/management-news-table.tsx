"use client";

import { useState } from "react";
import { EditNewsForm } from "@/components/management/edit-news-form";
import { NewsPinStarButton } from "@/components/management/news-pin-star-button";
import {
  ManagementEmptyRow,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge } from "@/components/management/status-badge";
import { FormError } from "@/components/form-error";
import { deleteNewsAction } from "@/app/management/actions";
import type { ManagementNewsArticle } from "@/lib/api/management";
import { formatNewsDate } from "@/lib/news";

type ManagementNewsTableProps = {
  articles: ManagementNewsArticle[];
};

export function ManagementNewsTable({ articles }: ManagementNewsTableProps) {
  const [pinError, setPinError] = useState<string | null>(null);

  return (
    <div>
      <FormError className="px-4 pt-3">{pinError}</FormError>
      <ManagementTable>
        <ManagementTableHead>
          <ManagementTh>Закрепить</ManagementTh>
          <ManagementTh>Статья</ManagementTh>
          <ManagementTh>Дата</ManagementTh>
          <ManagementTh>Статус</ManagementTh>
          <ManagementTh />
        </ManagementTableHead>
        <tbody>
          {articles.length === 0 ? (
            <ManagementEmptyRow colSpan={5}>Статей пока нет.</ManagementEmptyRow>
          ) : (
            articles.map((article) => (
              <tr key={article.id} className="border-b border-stone-100 align-top last:border-0">
                <td className="px-4 py-4">
                  <NewsPinStarButton
                    id={article.id}
                    isPinned={article.is_pinned}
                    onError={setPinError}
                  />
                </td>
                <td className="px-4 py-4">
                  <div className="font-medium text-stone-900">{article.title}</div>
                  <div className="text-stone-500">{article.slug}</div>
                </td>
                <td className="px-4 py-4 whitespace-nowrap">{formatNewsDate(article.published_at)}</td>
                <td className="px-4 py-4">
                  <div className="flex flex-col gap-1">
                    <StatusBadge variant={article.is_published ? "success" : "neutral"}>
                      {article.is_published ? "Опубликована" : "Черновик"}
                    </StatusBadge>
                    {article.is_pinned ? <StatusBadge variant="brand">Закреплена</StatusBadge> : null}
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="space-y-2">
                    <EditNewsForm article={article} />
                    <form action={deleteNewsAction}>
                      <input type="hidden" name="id" value={article.id} />
                      <button type="submit" className="btn-danger">
                        Удалить
                      </button>
                    </form>
                  </div>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </ManagementTable>
    </div>
  );
}
