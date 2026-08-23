import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { AuthAlert } from "@/components/auth/auth-alert";
import { fetchAuthMethods, fetchCurrentUser, fetchMyIdentities, type AuthMethods } from "@/lib/api/auth";
import { getAuthToken } from "@/lib/auth/session";
import { userContactLine } from "@/lib/site-nav";

export const metadata: Metadata = {
  title: "Профиль",
};

const PROVIDER_LABELS: Record<string, string> = {
  yandex: "Яндекс",
  vk: "VK",
  max: "Max",
  telegram: "Telegram",
};

const LINKABLE = [
  { key: "yandex" as const, label: "Яндекс", href: "/api/auth/social/yandex" },
  { key: "vk" as const, label: "VK", href: "/api/auth/social/vk" },
  { key: "max" as const, label: "Max", href: "/api/auth/social/max" },
];

const KEPT_LABELS: Record<string, string> = {
  email: "почта",
  phone: "телефон",
  name: "имя",
};

type PageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

function firstString(value: string | string[] | undefined): string {
  return typeof value === "string" ? value : "";
}

export default async function AccountProfilePage({ searchParams }: PageProps) {
  const token = await getAuthToken();
  if (!token) {
    redirect("/account/login?returnUrl=%2Faccount");
  }

  const params = await searchParams;

  let user;
  try {
    user = await fetchCurrentUser(token);
  } catch {
    redirect("/account/login?returnUrl=%2Faccount");
  }

  const [identities, methods] = await Promise.all([
    fetchMyIdentities(token).catch(() => []),
    fetchAuthMethods().catch((): AuthMethods => ({})),
  ]);
  const linkedProviders = new Set(identities.map((item) => item.provider));

  const linked = firstString(params.linked) === "1";
  const merged = firstString(params.merged) === "1";
  const kept = firstString(params.kept)
    .split(",")
    .map((item) => KEPT_LABELS[item.trim()])
    .filter(Boolean);

  let notice: string | null = null;
  if (merged) {
    notice =
      kept.length > 0
        ? `Входы объединены. Различавшиеся поля не меняли: ${kept.join(", ")}.`
        : "Входы объединены. Различавшиеся поля не меняли.";
  } else if (linked) {
    notice = "Вход привязан к этому аккаунту.";
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-display text-3xl font-semibold text-stone-900">Профиль</h1>
        <p className="mt-2 text-sm text-stone-600">Данные аккаунта. UUID нужен для назначения роли в админке.</p>
      </div>

      {notice ? <AuthAlert message={notice} /> : null}

      <dl className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
        <div>
          <dt className="text-xs uppercase tracking-wide text-stone-500">Имя</dt>
          <dd className="mt-1 text-stone-900">{user.name || "—"}</dd>
        </div>
        <div>
          <dt className="text-xs uppercase tracking-wide text-stone-500">Контакт</dt>
          <dd className="mt-1 text-stone-900">{userContactLine(user)}</dd>
        </div>
        <div>
          <dt className="text-xs uppercase tracking-wide text-stone-500">UUID пользователя</dt>
          <dd className="mt-1 break-all font-mono text-sm text-stone-900">{user.id}</dd>
        </div>
      </dl>

      <section className="space-y-3 rounded-2xl border border-stone-200 bg-white p-5">
        <h2 className="font-display text-xl font-semibold text-stone-900">Входы</h2>
        <p className="text-sm text-stone-600">
          Привязка соцсети к уже открытому кабинету объединяет заявки и избранное в этот аккаунт.
        </p>
        {identities.length > 0 ? (
          <ul className="space-y-2 text-sm text-stone-800">
            {identities.map((identity) => (
              <li key={`${identity.provider}:${identity.subject}`}>
                {PROVIDER_LABELS[identity.provider] ?? identity.provider}
                <span className="ml-2 font-mono text-xs text-stone-500">{identity.subject}</span>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-sm text-stone-500">Пока нет привязанных соцсетей.</p>
        )}
        <div className="flex flex-wrap gap-2">
          {LINKABLE.map((item) => {
            if (linkedProviders.has(item.key)) {
              return null;
            }
            const available = Boolean(methods[item.key]?.available);
            if (!available) {
              return (
                <button
                  key={item.key}
                  type="button"
                  disabled
                  className="rounded-xl border border-stone-200 bg-stone-50 px-3 py-2 text-sm text-stone-400"
                  title={methods[item.key]?.message || "Пока что недоступно, используйте другой вариант."}
                >
                  {item.label}
                </button>
              );
            }
            return (
              <a key={item.key} href={item.href} className="btn-secondary text-sm">
                Привязать {item.label}
              </a>
            );
          })}
        </div>
        {!linkedProviders.has("telegram") ? (
          methods.telegram?.available && methods.telegram.username ? (
            <p className="text-xs text-stone-500">
              Telegram Login: виджет для @{methods.telegram.username} (домен в BotFather). Обработчик тот же, с cookie
              сессии.
            </p>
          ) : (
            <p className="text-xs text-stone-400">
              Telegram Login: {methods.telegram?.message || "пока недоступно"}
            </p>
          )
        ) : null}
      </section>
    </div>
  );
}
