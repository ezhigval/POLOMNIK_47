import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

function uniqueCredentials() {
  const suffix = randomUUID().replace(/-/g, "").slice(0, 12);
  return {
    name: "E2E User",
    email: `e2e-${suffix}@example.com`,
    phone: `+79${suffix.slice(0, 9)}`,
    password: "e2epassword123",
  };
}

test("user can register and see empty trips page", async ({ page }) => {
  const user = uniqueCredentials();

  await page.goto("/account/register");
  await page.getByLabel(/Имя и фамилия/i).fill(user.name);
  await page.getByLabel(/^Телефон/i).fill(user.phone);
  await page.getByLabel(/^Email/i).fill(user.email);
  await page.getByLabel(/^Пароль/i).fill(user.password);
  await page.getByRole("checkbox", { name: /согласие на обработку моих персональных данных/i }).check();
  await page.getByRole("checkbox", { name: /принимаю условия/i }).check();
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
  await page.getByRole("checkbox", { name: /согласие на обработку моих персональных данных/i }).check();
  await page.getByRole("checkbox", { name: /принимаю условия/i }).check();
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
