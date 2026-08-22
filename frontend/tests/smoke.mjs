const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";
const ADMIN_TOKEN = process.env.ADMIN_TOKEN ?? "dev-admin-token";
const INTERNAL_API_SECRET = process.env.INTERNAL_API_SECRET ?? "dev-internal-api-secret";

const adminHeaders = {
  "Content-Type": "application/json",
  "X-Admin-Token": ADMIN_TOKEN,
};

async function request(path, init = {}) {
  const response = await fetch(`${API_URL}${path}`, init);
  const text = await response.text();
  let body;
  try {
    body = text ? JSON.parse(text) : null;
  } catch {
    body = text;
  }
  return { response, body };
}

function assert(condition, message) {
  if (!condition) {
    throw new Error(message);
  }
}

async function main() {
  const health = await request("/health");
  assert(health.response.ok, "health endpoint failed");

  const ready = await request("/health/ready");
  assert(ready.response.ok, "health ready endpoint failed");
  assert(ready.body?.data?.status === "ok", "health ready status is not ok");

  const oauthBlocked = await request("/auth/oauth", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      provider: "google",
      subject: "smoke-attacker",
      email: "attacker@example.com",
      name: "Attacker",
    }),
  });
  assert(oauthBlocked.response.status === 401, "oauth endpoint must reject requests without internal secret");

  const suffix = Date.now();
  const register = await request("/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name: "Smoke User",
      email: `smoke-${suffix}@example.com`,
      phone: `+7999${String(suffix).slice(-7)}`,
      password: "smokepassword123",
    }),
  });
  assert(register.response.status === 201, "register user failed");

  const authToken = register.body.data.token;
  assert(authToken, "register response must include token");

  const me = await request("/me", {
    headers: { Authorization: `Bearer ${authToken}` },
  });
  assert(me.response.ok, "authenticated /me failed");

  const favorites = await request("/me/favorites", {
    headers: { Authorization: `Bearer ${authToken}` },
  });
  assert(favorites.response.ok, "list favorites failed");
  assert(Array.isArray(favorites.body.data), "favorites data must be array");

  const login = await request("/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      login: `smoke-${suffix}@example.com`,
      password: "smokepassword123",
    }),
  });
  assert(login.response.status === 200, "login failed");

  const tourPayload = {
    slug: `smoke-${suffix}`,
    title: "Smoke Test Tour",
    description: "Automated smoke test tour",
    price: 12000,
    currency: "RUB",
    date_start: "2026-10-01",
    date_end: "2026-10-05",
    slots_total: 10,
    slots_left: 10,
    location: "Moscow",
    images: [],
    is_active: true,
    is_hot: true,
    overbooking_enabled: false,
  };

  const created = await request("/management/tours", {
    method: "POST",
    headers: adminHeaders,
    body: JSON.stringify(tourPayload),
  });
  assert(created.response.status === 201, "create tour failed");

  const tourId = created.body.data.id;
  const publicTour = await request(`/tours/${tourId}`);
  assert(publicTour.response.ok, "public tour endpoint failed");

  const booking = await request("/bookings", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    },
    body: JSON.stringify({
      tour_id: tourId,
      name: "Smoke Tester",
      phone: "+79990000000",
      people_count: 2,
    }),
  });
  assert(booking.response.status === 201, "create booking failed");
  assert(
    booking.body.data.integration_status === "not_configured",
    "expected not_configured integration status",
  );

  const bookingId = booking.body.data.booking_id;
  const statusUpdate = await request(`/management/bookings/${bookingId}/status`, {
    method: "PATCH",
    headers: adminHeaders,
    body: JSON.stringify({ status: "CONTACTED" }),
  });
  assert(statusUpdate.response.ok, "update booking status failed");
  assert(statusUpdate.body.data.status === "CONTACTED", "expected CONTACTED status");

  const myBookings = await request("/me/bookings", {
    headers: { Authorization: `Bearer ${authToken}` },
  });
  assert(myBookings.response.ok, "authenticated bookings list failed");
  assert(Array.isArray(myBookings.body.data), "bookings data must be array");

  const favoriteAdd = await request(`/me/favorites/${tourId}`, {
    method: "POST",
    headers: { Authorization: `Bearer ${authToken}` },
  });
  assert(favoriteAdd.response.status === 201, "add favorite failed");

  const cmsPages = await request("/management/cms/pages", {
    headers: adminHeaders,
  });
  assert(cmsPages.response.ok, "management cms pages must not 500 (missing SEO columns or empty CMS)");
  assert(Array.isArray(cmsPages.body?.data), "cms pages data must be array");

  const publicHome = await request("/pages/home");
  assert(
    publicHome.response.status === 200 || publicHome.response.status === 404,
    `public CMS home must not 500, got ${publicHome.response.status}`,
  );

  const integrationRefs = await request("/management/integration-references", {
    headers: adminHeaders,
  });
  assert(integrationRefs.response.ok, "list integration references failed");

  const outboxEvents = await request("/management/outbox-events", {
    headers: adminHeaders,
  });
  assert(outboxEvents.response.ok, "list outbox events failed");
  assert(Array.isArray(outboxEvents.body.data), "outbox events data must be array");

  const oauthAllowed = await request("/auth/oauth", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Internal-Secret": INTERNAL_API_SECRET,
    },
    body: JSON.stringify({
      provider: "google",
      subject: `google-smoke-${suffix}`,
      email: `google-${suffix}@example.com`,
      name: "Google Smoke",
    }),
  });
  assert(oauthAllowed.response.status === 200, "oauth with internal secret failed");

  console.log("Smoke test passed");
}

main().catch((error) => {
  console.error("Smoke test failed:", error.message);
  process.exit(1);
});
