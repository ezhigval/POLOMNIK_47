import { expect, test } from "@playwright/test";

/** Tour with future dates in dev seed (Валаам, октябрь 2026) */
const BOOKABLE_TOUR_ID = "33333333-4444-4444-4444-444444444444";

test("guest can browse search and submit a booking", async ({ page }) => {
  await page.goto("/search");
  await expect(page.getByRole("heading", { name: "Расписание", level: 1 })).toBeVisible();

  await page.goto(`/tours/${BOOKABLE_TOUR_ID}`);
  await expect(page.getByRole("heading", { level: 1 })).toContainText("Валаам");

  const bookingForm = page.locator("#booking-form");
  await bookingForm.scrollIntoViewIfNeeded();
  await bookingForm.getByLabel(/Имя и фамилия/i).fill("E2E Guest");
  await bookingForm.getByLabel(/^Телефон/i).fill("+79991112233");
  const cookieBanner = page.getByRole("dialog", { name: "Настройки cookie" });
  if (await cookieBanner.isVisible().catch(() => false)) {
    await cookieBanner.getByRole("button", { name: "Только необходимые" }).click();
    await expect(cookieBanner).toBeHidden();
  }
  await bookingForm.locator('input[name="consent_personal_data"]').evaluate((el: HTMLInputElement) => {
    el.scrollIntoView({ block: "center" });
    el.click();
  });
  await bookingForm.getByRole("button", { name: "Отправить заявку" }).click();

  await expect(page.getByRole("heading", { name: "Заявка отправлена" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Другие туры" })).toBeVisible();
});

test("404 tour page offers search navigation", async ({ page }) => {
  await page.goto("/tours/00000000-0000-0000-0000-000000000000");
  await expect(page.getByRole("heading", { name: "Тур не найден" })).toBeVisible();
  await page.getByRole("main").getByRole("link", { name: "Туры" }).click();
  await expect(page).toHaveURL(/\/search/);
});
