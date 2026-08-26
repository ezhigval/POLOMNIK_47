"use client";

import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { updateCmsBlockAction } from "@/app/management/actions";
import type { CmsBlock, CmsBlockType } from "@/lib/api/cms";
import { featuredRoute } from "@/lib/featured-route";
import { FormError } from "@/components/form-error";

type BlockEditorProps = {
  block: CmsBlock;
  pageId: string;
  slug: string;
};

type StringRow = { key: string; value: string };

function stringValue(content: Record<string, unknown>, key: string) {
  const value = content[key];
  return typeof value === "string" ? value : "";
}

function boolValue(content: Record<string, unknown>, key: string, fallback = false) {
  const value = content[key];
  return typeof value === "boolean" ? value : fallback;
}

function stringArray(content: Record<string, unknown>, key: string): string[] {
  const value = content[key];
  return Array.isArray(value) ? value.map((item) => String(item ?? "")) : [];
}

function objectArray<T extends Record<string, string>>(
  content: Record<string, unknown>,
  key: string,
  keys: (keyof T)[],
): T[] {
  const value = content[key];
  if (!Array.isArray(value)) return [];
  return value.map((item) => {
    const row = {} as T;
    const source = item && typeof item === "object" ? (item as Record<string, unknown>) : {};
    for (const field of keys) {
      row[field] = String(source[String(field)] ?? "") as T[keyof T];
    }
    return row;
  });
}

export function BlockEditor({ block, pageId, slug }: BlockEditorProps) {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [isVisible, setIsVisible] = useState(block.is_visible);
  const [content, setContent] = useState<Record<string, unknown>>(block.content ?? {});

  const label = useMemo(() => blockTypeLabel(block.type), [block.type]);

  async function save(nextContent: Record<string, unknown>, nextVisible = isVisible) {
    setLoading(true);
    setError(null);
    try {
      await updateCmsBlockAction(
        block.id,
        { content: nextContent, is_visible: nextVisible },
        pageId,
        slug,
      );
      setContent(nextContent);
      setIsVisible(nextVisible);
      setOpen(false);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось сохранить блок");
    } finally {
      setLoading(false);
    }
  }

  async function onToggleVisible() {
    const next = !isVisible;
    setLoading(true);
    setError(null);
    try {
      await updateCmsBlockAction(block.id, { is_visible: next }, pageId, slug);
      setIsVisible(next);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось обновить видимость");
    } finally {
      setLoading(false);
    }
  }

  if (!open) {
    return (
      <div className="flex flex-wrap items-center gap-2">
        <button
          type="button"
          onClick={() => setOpen(true)}
          className="text-sm font-medium text-brand-800 hover:text-brand-900"
        >
          Редактировать
        </button>
        <button
          type="button"
          onClick={onToggleVisible}
          disabled={loading}
          className="text-sm text-stone-600 hover:text-stone-900"
        >
          {isVisible ? "Скрыть" : "Показать"}
        </button>
        <span className="text-xs text-stone-400">{label}</span>
        <FormError className="w-full">{error}</FormError>
      </div>
    );
  }

  return (
    <div className="mt-3 space-y-3 rounded-xl border border-stone-200 bg-stone-50 p-3">
      <div className="flex items-center justify-between gap-3">
        <p className="text-sm font-medium text-stone-900">{label}</p>
        <button type="button" onClick={() => setOpen(false)} className="text-sm text-stone-500">
          Закрыть
        </button>
      </div>

      <label className="flex items-center gap-2 text-sm">
        <input
          type="checkbox"
          checked={isVisible}
          onChange={(event) => setIsVisible(event.target.checked)}
          className="size-4"
        />
        Видимый
      </label>

      <BlockContentForm
        type={block.type}
        content={content}
        loading={loading}
        error={error}
        onCancel={() => setOpen(false)}
        onSave={(next) => save(next, isVisible)}
      />
    </div>
  );
}

type BlockFormFields = {
  content: Record<string, unknown>;
  loading: boolean;
  error: string | null;
  onCancel: () => void;
  onSave: (content: Record<string, unknown>) => void;
};

type BlockContentFormProps = BlockFormFields & {
  type: CmsBlockType;
};

function BlockContentForm({ type, content, loading, error, onCancel, onSave }: BlockContentFormProps) {
  if (type === "popular_destinations" || type === "testimonials") {
    return (
      <div className="space-y-3">
        <p className="text-sm text-stone-600">
          {type === "popular_destinations"
            ? "На сайте показываются туры с пометкой «Популярный». Если таких нет, блок скрывается."
            : "Виджет без настроек контента."}
        </p>
        <FormError>{error}</FormError>
        <div className="flex gap-2">
          <button type="button" disabled={loading} onClick={() => onSave({})} className="btn-primary">
            {loading ? "Сохраняем..." : "Сохранить"}
          </button>
          <button type="button" onClick={onCancel} className="btn-secondary">
            Отмена
          </button>
        </div>
      </div>
    );
  }

  if (type === "featured_route") {
    return (
      <FeaturedRouteForm content={content} loading={loading} error={error} onCancel={onCancel} onSave={onSave} />
    );
  }

  if (type === "hero") {
    return (
      <HeroForm content={content} loading={loading} error={error} onCancel={onCancel} onSave={onSave} />
    );
  }
  if (type === "about") {
    return (
      <AboutForm content={content} loading={loading} error={error} onCancel={onCancel} onSave={onSave} />
    );
  }
  if (type === "why_us") {
    return (
      <WhyUsForm content={content} loading={loading} error={error} onCancel={onCancel} onSave={onSave} />
    );
  }
  if (type === "how_it_works") {
    return (
      <HowItWorksForm content={content} loading={loading} error={error} onCancel={onCancel} onSave={onSave} />
    );
  }
  if (type === "faq") {
    return (
      <FaqForm content={content} loading={loading} error={error} onCancel={onCancel} onSave={onSave} />
    );
  }
  if (type === "cta") {
    return (
      <CtaForm content={content} loading={loading} error={error} onCancel={onCancel} onSave={onSave} />
    );
  }
  if (type === "rich_text") {
    return (
      <RichTextForm content={content} loading={loading} error={error} onCancel={onCancel} onSave={onSave} />
    );
  }

  return <p className="text-sm text-stone-500">Неизвестный тип блока.</p>;
}

function Field({
  label,
  value,
  onChange,
  multiline = false,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  multiline?: boolean;
}) {
  return (
    <label className="block text-sm">
      <span className="mb-1 block font-medium">{label}</span>
      {multiline ? (
        <textarea
          rows={3}
          value={value}
          onChange={(event) => onChange(event.target.value)}
          className="input-field"
        />
      ) : (
        <input
          value={value}
          onChange={(event) => onChange(event.target.value)}
          className="input-field"
        />
      )}
    </label>
  );
}

function FormActions({
  loading,
  error,
  onCancel,
}: {
  loading: boolean;
  error: string | null;
  onCancel: () => void;
}) {
  return (
    <>
      <FormError>{error}</FormError>
      <div className="flex flex-wrap gap-2">
        <button type="submit" disabled={loading} className="btn-primary">
          {loading ? "Сохраняем..." : "Сохранить"}
        </button>
        <button type="button" onClick={onCancel} className="btn-secondary">
          Отмена
        </button>
      </div>
    </>
  );
}

function StringListEditor({
  label,
  values,
  onChange,
}: {
  label: string;
  values: string[];
  onChange: (values: string[]) => void;
}) {
  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium">{label}</p>
        <button
          type="button"
          className="text-sm text-brand-800"
          onClick={() => onChange([...values, ""])}
        >
          + Добавить
        </button>
      </div>
      {values.map((value, index) => (
        <div key={index} className="flex gap-2">
          <input
            value={value}
            onChange={(event) => {
              const next = [...values];
              next[index] = event.target.value;
              onChange(next);
            }}
            className="input-field"
          />
          <button
            type="button"
            className="btn-danger shrink-0"
            onClick={() => onChange(values.filter((_, i) => i !== index))}
          >
            ×
          </button>
        </div>
      ))}
    </div>
  );
}

function ObjectListEditor({
  label,
  rows,
  fields,
  onChange,
}: {
  label: string;
  rows: Record<string, string>[];
  fields: StringRow[];
  onChange: (rows: Record<string, string>[]) => void;
}) {
  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium">{label}</p>
        <button
          type="button"
          className="text-sm text-brand-800"
          onClick={() => {
            const empty: Record<string, string> = {};
            for (const field of fields) empty[field.key] = "";
            onChange([...rows, empty]);
          }}
        >
          + Добавить
        </button>
      </div>
      {rows.map((row, index) => (
        <div key={index} className="space-y-2 rounded-lg border border-stone-200 bg-white p-3">
          {fields.map((field) => (
            <Field
              key={field.key}
              label={field.value}
              value={row[field.key] ?? ""}
              onChange={(value) => {
                const next = rows.map((item, i) =>
                  i === index ? { ...item, [field.key]: value } : item,
                );
                onChange(next);
              }}
              multiline={field.key === "description" || field.key === "answer"}
            />
          ))}
          <button
            type="button"
            className="btn-danger"
            onClick={() => onChange(rows.filter((_, i) => i !== index))}
          >
            Удалить
          </button>
        </div>
      ))}
    </div>
  );
}

function HeroForm({ content, loading, error, onCancel, onSave }: BlockFormFields) {
  const [eyebrow, setEyebrow] = useState(stringValue(content, "eyebrow"));
  const [title, setTitle] = useState(stringValue(content, "title"));
  const [subtitle, setSubtitle] = useState(stringValue(content, "subtitle"));

  return (
    <form
      className="space-y-3"
      onSubmit={(event) => {
        event.preventDefault();
        onSave({
          eyebrow,
          title,
          subtitle,
        });
      }}
    >
      <Field label="Надзаголовок" value={eyebrow} onChange={setEyebrow} />
      <Field label="Заголовок" value={title} onChange={setTitle} />
      <Field label="Подзаголовок" value={subtitle} onChange={setSubtitle} multiline />
      <FormActions loading={loading} error={error} onCancel={onCancel} />
    </form>
  );
}

function FeaturedRouteForm({ content, loading, error, onCancel, onSave }: BlockFormFields) {
  const [eyebrow, setEyebrow] = useState(stringValue(content, "eyebrow") || featuredRoute.eyebrow);
  const [title, setTitle] = useState(stringValue(content, "title") || featuredRoute.title);
  const [parentRoute, setParentRoute] = useState(stringValue(content, "parentRoute") || featuredRoute.parentRoute);
  const [duration, setDuration] = useState(stringValue(content, "duration") || featuredRoute.duration);
  const [region, setRegion] = useState(stringValue(content, "region") || featuredRoute.region);
  const [lead, setLead] = useState(stringValue(content, "lead") || featuredRoute.lead);
  const [body, setBody] = useState(stringValue(content, "body") || featuredRoute.body);
  const [ctaLabel, setCtaLabel] = useState(stringValue(content, "ctaLabel") || featuredRoute.ctaLabel);
  const [ctaHref, setCtaHref] = useState(stringValue(content, "ctaHref") || featuredRoute.ctaHref);
  const [secondaryCta, setSecondaryCta] = useState(stringValue(content, "secondaryCta") || featuredRoute.secondaryCta);
  const [secondaryHref, setSecondaryHref] = useState(stringValue(content, "secondaryHref") || featuredRoute.secondaryHref);
  const [days, setDays] = useState(featuredDaysForEditor(content));

  return (
    <form
      className="space-y-3"
      onSubmit={(event) => {
        event.preventDefault();
        onSave({
          eyebrow,
          title,
          parentRoute,
          duration,
          region,
          lead,
          body,
          ctaLabel,
          ctaHref,
          secondaryCta,
          secondaryHref,
          days: days.map((day) => ({
            title: day.title,
            points: day.points
              .split("\n")
              .map((point) => point.trim())
              .filter(Boolean),
          })),
        });
      }}
    >
      <Field label="Надзаголовок" value={eyebrow} onChange={setEyebrow} />
      <Field label="Заголовок" value={title} onChange={setTitle} />
      <Field label="Родительский маршрут" value={parentRoute} onChange={setParentRoute} />
      <Field label="Длительность" value={duration} onChange={setDuration} />
      <Field label="Регион" value={region} onChange={setRegion} />
      <Field label="Лид" value={lead} onChange={setLead} multiline />
      <Field label="Текст" value={body} onChange={setBody} multiline />
      <ObjectListEditor
        label="Программа (пункты дня — с новой строки)"
        rows={days}
        fields={[
          { key: "title", value: "День" },
          { key: "points", value: "Пункты" },
        ]}
        onChange={(rows) => setDays(rows as { title: string; points: string }[])}
      />
      <Field label="Основная кнопка" value={ctaLabel} onChange={setCtaLabel} />
      <Field label="Ссылка основной кнопки" value={ctaHref} onChange={setCtaHref} />
      <Field label="Вторая кнопка" value={secondaryCta} onChange={setSecondaryCta} />
      <Field label="Ссылка второй кнопки" value={secondaryHref} onChange={setSecondaryHref} />
      <FormActions loading={loading} error={error} onCancel={onCancel} />
    </form>
  );
}

function featuredDaysForEditor(content: Record<string, unknown>): { title: string; points: string }[] {
  const value = content.days;
  if (!Array.isArray(value) || value.length === 0) {
    return featuredRoute.days.map((day) => ({ title: day.title, points: day.points.join("\n") }));
  }
  return value.map((item) => {
    const source = item && typeof item === "object" ? (item as Record<string, unknown>) : {};
    const points = Array.isArray(source.points)
      ? source.points.map((point) => String(point ?? "")).join("\n")
      : String(source.points ?? "");
    return { title: String(source.title ?? ""), points };
  });
}

function AboutForm({ content, loading, error, onCancel, onSave }: BlockFormFields) {
  const [eyebrow, setEyebrow] = useState(stringValue(content, "eyebrow"));
  const [title, setTitle] = useState(stringValue(content, "title"));
  const [paragraphs, setParagraphs] = useState(stringArray(content, "paragraphs"));
  const [highlights, setHighlights] = useState(stringArray(content, "highlights"));
  const [showContacts, setShowContacts] = useState(boolValue(content, "showContacts", true));
  return (
    <form
      className="space-y-3"
      onSubmit={(event) => {
        event.preventDefault();
        onSave({ eyebrow, title, paragraphs, highlights, showContacts, stats });
      }}
    >
      <Field label="Надзаголовок" value={eyebrow} onChange={setEyebrow} />
      <Field label="Заголовок" value={title} onChange={setTitle} />
      <StringListEditor label="Абзацы" values={paragraphs} onChange={setParagraphs} />
      <StringListEditor label="Акценты" values={highlights} onChange={setHighlights} />
      <ObjectListEditor
        label="Цифры"
        rows={stats}
        fields={[
          { key: "value", value: "Значение" },
          { key: "label", value: "Подпись" },
        ]}
        onChange={(rows) => setStats(rows as { value: string; label: string }[])}
      />
      <label className="flex items-center gap-2 text-sm">
        <input
          type="checkbox"
          checked={showContacts}
          onChange={(event) => setShowContacts(event.target.checked)}
          className="size-4"
        />
        Показывать контакты
      </label>
      <FormActions loading={loading} error={error} onCancel={onCancel} />
    </form>
  );
}

function WhyUsForm({ content, loading, error, onCancel, onSave }: BlockFormFields) {
  const [eyebrow, setEyebrow] = useState(stringValue(content, "eyebrow"));
  const [title, setTitle] = useState(stringValue(content, "title"));
  const [description, setDescription] = useState(stringValue(content, "description"));
  const [items, setItems] = useState(
    objectArray<{ title: string; description: string; icon: string }>(content, "items", [
      "title",
      "description",
      "icon",
    ]),
  );
  const [stats, setStats] = useState(
    objectArray<{ value: string; label: string }>(content, "stats", ["value", "label"]),
  );

  return (
    <form
      className="space-y-3"
      onSubmit={(event) => {
        event.preventDefault();
        onSave({ eyebrow, title, description, items });
      }}
    >
      <Field label="Надзаголовок" value={eyebrow} onChange={setEyebrow} />
      <Field label="Заголовок" value={title} onChange={setTitle} />
      <Field label="Описание" value={description} onChange={setDescription} multiline />
      <ObjectListEditor
        label="Пункты"
        rows={items}
        fields={[
          { key: "title", value: "Заголовок" },
          { key: "description", value: "Описание" },
          { key: "icon", value: "Иконка (route/cross/shield/wallet)" },
        ]}
        onChange={(rows) => setItems(rows as { title: string; description: string; icon: string }[])}
      />
      <FormActions loading={loading} error={error} onCancel={onCancel} />
    </form>
  );
}

function HowItWorksForm({ content, loading, error, onCancel, onSave }: BlockFormFields) {
  const [eyebrow, setEyebrow] = useState(stringValue(content, "eyebrow"));
  const [title, setTitle] = useState(stringValue(content, "title"));
  const [description, setDescription] = useState(stringValue(content, "description"));
  const [ctaLabel, setCtaLabel] = useState(stringValue(content, "ctaLabel"));
  const [ctaHref, setCtaHref] = useState(stringValue(content, "ctaHref"));
  const [steps, setSteps] = useState(
    objectArray<{ title: string; description: string }>(content, "steps", ["title", "description"]),
  );

  return (
    <form
      className="space-y-3"
      onSubmit={(event) => {
        event.preventDefault();
        onSave({ eyebrow, title, description, steps, ctaLabel, ctaHref });
      }}
    >
      <Field label="Надзаголовок" value={eyebrow} onChange={setEyebrow} />
      <Field label="Заголовок" value={title} onChange={setTitle} />
      <Field label="Описание" value={description} onChange={setDescription} multiline />
      <ObjectListEditor
        label="Шаги"
        rows={steps}
        fields={[
          { key: "title", value: "Заголовок" },
          { key: "description", value: "Описание" },
        ]}
        onChange={(rows) => setSteps(rows as { title: string; description: string }[])}
      />
      <Field label="Текст кнопки" value={ctaLabel} onChange={setCtaLabel} />
      <Field label="Ссылка кнопки" value={ctaHref} onChange={setCtaHref} />
      <FormActions loading={loading} error={error} onCancel={onCancel} />
    </form>
  );
}

function FaqForm({ content, loading, error, onCancel, onSave }: BlockFormFields) {
  const [eyebrow, setEyebrow] = useState(stringValue(content, "eyebrow"));
  const [title, setTitle] = useState(stringValue(content, "title"));
  const [description, setDescription] = useState(stringValue(content, "description"));
  const [items, setItems] = useState(
    objectArray<{ question: string; answer: string }>(content, "items", ["question", "answer"]),
  );

  return (
    <form
      className="space-y-3"
      onSubmit={(event) => {
        event.preventDefault();
        onSave({ eyebrow, title, description, items });
      }}
    >
      <Field label="Надзаголовок" value={eyebrow} onChange={setEyebrow} />
      <Field label="Заголовок" value={title} onChange={setTitle} />
      <Field label="Описание" value={description} onChange={setDescription} multiline />
      <ObjectListEditor
        label="Вопросы"
        rows={items}
        fields={[
          { key: "question", value: "Вопрос" },
          { key: "answer", value: "Ответ" },
        ]}
        onChange={(rows) => setItems(rows as { question: string; answer: string }[])}
      />
      <FormActions loading={loading} error={error} onCancel={onCancel} />
    </form>
  );
}

function CtaForm({ content, loading, error, onCancel, onSave }: BlockFormFields) {
  const [title, setTitle] = useState(stringValue(content, "title"));
  const [subtitle, setSubtitle] = useState(stringValue(content, "subtitle"));
  const [button, setButton] = useState(stringValue(content, "button"));
  const [href, setHref] = useState(stringValue(content, "href"));

  return (
    <form
      className="space-y-3"
      onSubmit={(event) => {
        event.preventDefault();
        onSave({ title, subtitle, button, href });
      }}
    >
      <Field label="Заголовок" value={title} onChange={setTitle} />
      <Field label="Подзаголовок" value={subtitle} onChange={setSubtitle} multiline />
      <Field label="Кнопка" value={button} onChange={setButton} />
      <Field label="Ссылка" value={href} onChange={setHref} />
      <FormActions loading={loading} error={error} onCancel={onCancel} />
    </form>
  );
}

function RichTextForm({ content, loading, error, onCancel, onSave }: BlockFormFields) {
  const [eyebrow, setEyebrow] = useState(stringValue(content, "eyebrow"));
  const [title, setTitle] = useState(stringValue(content, "title"));
  const [body, setBody] = useState(stringValue(content, "body"));

  return (
    <form
      className="space-y-3"
      onSubmit={(event) => {
        event.preventDefault();
        onSave({ eyebrow, title, body });
      }}
    >
      <Field label="Надзаголовок" value={eyebrow} onChange={setEyebrow} />
      <Field label="Заголовок" value={title} onChange={setTitle} />
      <Field label="Текст" value={body} onChange={setBody} multiline />
      <FormActions loading={loading} error={error} onCancel={onCancel} />
    </form>
  );
}

function blockTypeLabel(type: CmsBlockType) {
  switch (type) {
    case "hero":
      return "Шапка";
    case "about":
      return "О службе";
    case "why_us":
      return "Почему мы";
    case "how_it_works":
      return "Как записаться";
    case "faq":
      return "Вопросы";
    case "cta":
      return "Баннер";
    case "rich_text":
      return "Текст";
    case "popular_destinations":
      return "Направления";
    case "testimonials":
      return "Отзывы";
    case "featured_route":
      return "Тихвинский путь";
    default:
      return type;
  }
}
