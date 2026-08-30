import { expect, test, type Page } from "@playwright/test";
import { randomUUID } from "node:crypto";

async function dismissCookieBanner(page: Page) {
  const banner = page.getByRole("dialog", { name: "Настройки cookie" });
  if (await banner.isVisible().catch(() => false)) {
    await banner.getByRole("button", { name: "Только необходимые" }).click();
    await expect(banner).toBeHidden();
  }
}

async function acceptRequiredConsents(page: Page) {
  await dismissCookieBanner(page);
  const form = page.locator("form").first();
  await form.locator('input[name="consent_personal_data"]').evaluate((el: HTMLInputElement) => el.click());
  await form.locator('input[name="consent_terms"]').evaluate((el: HTMLInputElement) => el.click());
}

function uniqueCredentials() {
  const suffix = randomUUID().replace(/-/g, "").slice(0, 12);
  return {
    name: "E2E User",
    email: `e2e-${suffix}@example.com`,
    phone: `+79${suffix.slice(0, 9)}`,
    password: "e2epassword123",
  };
}

async function expectBelow(above: { y: number; height: number }, below: { y: number }) {
  expect(below.y).toBeGreaterThan(above.y + above.height / 2);
}

test("login keeps phone call collapsed and social buttons under submit", async ({ page }) => {
  await page.goto("/account/login");
  await dismissCookieBanner(page);

  const submit = page.getByRole("button", { name: "Войти" });
  const social = page.getByText("Или войти через");
  const callBlock = page.locator("details").filter({ hasText: "Вход по звонку" });

  await expect(submit).toBeVisible();
  await expect(social).toBeVisible();
  await expect(callBlock).toBeVisible();
  await expect(callBlock).not.toHaveAttribute("open");
  await expect(page.getByLabel("Телефон аккаунта")).toBeHidden();

  const submitBox = await submit.boundingBox();
  const socialBox = await social.boundingBox();
  expect(submitBox).toBeTruthy();
  expect(socialBox).toBeTruthy();
  await expectBelow(submitBox!, socialBox!);

  await callBlock.locator("summary").click();
  await expect(callBlock).toHaveAttribute("open", "");
  await expect(page.getByLabel("Телефон аккаунта")).toBeVisible();
});

test("register keeps phone call collapsed and social buttons under submit", async ({ page }) => {
  await page.goto("/account/register");
  await dismissCookieBanner(page);

  const submit = page.getByRole("button", { name: "Создать аккаунт" });
  const social = page.getByText("Или войти через");
  const callBlock = page.locator("details").filter({ hasText: "Подтверждение телефона звонком" });

  await expect(submit).toBeVisible();
  await expect(social).toBeVisible();
  await expect(callBlock).toBeVisible();
  await expect(callBlock).not.toHaveAttribute("open");

  const submitBox = await submit.boundingBox();
  const socialBox = await social.boundingBox();
  expect(submitBox).toBeTruthy();
  expect(socialBox).toBeTruthy();
  await expectBelow(submitBox!, socialBox!);

  await callBlock.locator("summary").click();
  await expect(callBlock).toHaveAttribute("open", "");
});

test("user can register and see empty trips page", async ({ page }) => {
  const user = uniqueCredentials();

  await page.goto("/account/register");
  await page.getByLabel(/Имя и фамилия/i).fill(user.name);
  await page.getByLabel(/^Телефон/i).fill(user.phone);
  await page.getByLabel(/^Email/i).fill(user.email);
  await page.getByLabel(/^Пароль/i).fill(user.password);
  await acceptRequiredConsents(page);
  await page.getByRole("button", { name: "Создать аккаунт" }).click();

  await expect(page).toHaveURL(/\/account\/trips/);
  await expect(page.getByRole("heading", { name: "Мои поездки" })).toBeVisible();
  await expect(page.getByText("Пока нет заявок")).toBeVisible();
});

test("user can login after registration", async ({ page }) => {
  const user = uniqueCredentials();

  await page.goto("/account/register");
  await page.getByLabel(/Имя и фамилия/i).fill(user.name);
  await page.getByLabel(/^Телефон/i).fill(user.phone);
  await page.getByLabel(/^Email/i).fill(user.email);
  await page.getByLabel(/^Пароль/i).fill(user.password);
  await acceptRequiredConsents(page);
  await page.getByRole("button", { name: "Создать аккаунт" }).click();
  await expect(page).toHaveURL(/\/account\/trips/);

  await page.request.post("/api/auth/logout");
  await page.context().clearCookies();

  await page.goto("/account/login");
  await page.getByLabel(/Телефон или email/i).fill(user.email);
  await page.getByLabel(/^Пароль/i).fill(user.password);
  await page.getByRole("button", { name: "Войти" }).click();

  await expect(page).toHaveURL(/\/account\/trips/);
  await expect(
    page.getByRole("navigation", { name: "Разделы личного кабинета" }).getByRole("link", { name: "Избранное" }),
  ).toBeVisible();
});
