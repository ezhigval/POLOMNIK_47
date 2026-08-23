const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";
const ADMIN_TOKEN = process.env.ADMIN_TOKEN ?? "dev-admin-token";
const BITRIX_INBOUND_TOKEN = process.env.BITRIX_INBOUND_TOKEN ?? "dev-bitrix-inbound";
const MOCK_BITRIX_URL = process.env.MOCK_BITRIX_URL ?? "http://localhost:8091";
const MOCK_ONEC_URL = process.env.MOCK_ONEC_URL ?? "http://localhost:8092";

// Seed tour with future dates (Валаам, октябрь 2026 — backend/seeds/dev.sql)
const SEED_TOUR_ID = "33333333-4444-4444-4444-444444444444";

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

function findRef(refs, predicate) {
  return refs.find(predicate);
}

async function main() {
  const suffix = Date.now();

  const tour = await request(`/tours/${SEED_TOUR_ID}`);
  assert(tour.response.ok, `seed tour ${SEED_TOUR_ID} not found — run docker compose up with seed`);

  const booking = await request("/bookings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      tour_id: SEED_TOUR_ID,
      name: "Integration Tester",
      phone: `+7999${String(suffix).slice(-7)}`,
      email: `integration-${suffix}@example.com`,
      people_count: 2,
      comment: "integration smoke",
    }),
  });
  assert(booking.response.status === 201, "create booking failed");
  assert(
    booking.body.data.integration_status === "synced",
    `expected bitrix synced, got ${booking.body.data.integration_status}`,
  );

  const bookingId = booking.body.data.booking_id;

  const refsAfterCreate = await request("/management/integration-references", {
    headers: adminHeaders,
  });
  assert(refsAfterCreate.response.ok, "list integration references failed");
  const refs = refsAfterCreate.body.data ?? [];

  const bitrixDeal = findRef(
    refs,
    (r) =>
      r.local_entity_id === bookingId &&
      r.external_system === "bitrix24" &&
      r.external_entity_type === "deal",
  );
  assert(bitrixDeal, "bitrix24 deal reference missing");
  assert(bitrixDeal.sync_status === "synced", `bitrix deal not synced: ${bitrixDeal.sync_status}`);
  assert(bitrixDeal.external_entity_id, "bitrix deal external id empty");

  const onecCounterparty = findRef(
    refs,
    (r) =>
      r.local_entity_id === bookingId &&
      r.external_system === "onec" &&
      r.external_entity_type === "counterparty",
  );
  assert(onecCounterparty, "1c counterparty reference missing");
  assert(onecCounterparty.sync_status === "synced", "1c counterparty not synced");

  const dealId = bitrixDeal.external_entity_id;

  const stageUpdate = await fetch(`${MOCK_BITRIX_URL}/debug/deals/${dealId}/stage`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ stage_id: "PREPARATION" }),
  });
  assert(stageUpdate.ok, "mock-bitrix stage update failed");

  const webhook = await request(`/webhooks/bitrix/deal?token=${BITRIX_INBOUND_TOKEN}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      event: "ONCRMDEALUPDATE",
      data: { FIELDS: { ID: dealId, STAGE_ID: "PREPARATION" } },
    }),
  });
  assert(webhook.response.status === 204, "bitrix inbound webhook failed");

  const bookingAfterWebhook = await request(`/management/bookings/${bookingId}`, {
    headers: adminHeaders,
  });
  assert(bookingAfterWebhook.response.ok, "get booking after webhook failed");
  assert(
    bookingAfterWebhook.body.data.status === "CONTACTED",
    `expected CONTACTED after bitrix webhook, got ${bookingAfterWebhook.body.data.status}`,
  );

  const confirmed = await request(`/management/bookings/${bookingId}/status`, {
    method: "PATCH",
    headers: adminHeaders,
    body: JSON.stringify({ status: "CONFIRMED" }),
  });
  assert(confirmed.response.ok, "confirm booking failed");

  const refsAfterConfirm = await request("/management/integration-references", {
    headers: adminHeaders,
  });
  const refs2 = refsAfterConfirm.body.data ?? [];
  const onecOrder = findRef(
    refs2,
    (r) =>
      r.local_entity_id === bookingId &&
      r.external_system === "onec" &&
      r.external_entity_type === "order",
  );
  assert(onecOrder, "1c order export reference missing after CONFIRMED");
  assert(onecOrder.sync_status === "synced", "1c order not synced");

  const onecDocs = await fetch(`${MOCK_ONEC_URL}/accounting/debug/documents`);
  assert(onecDocs.ok, "mock-onec debug endpoint failed");
  const docs = await onecDocs.json();
  assert(docs.bookings?.[bookingId], "mock-onec has no booking document");

  console.log("Integration smoke test passed");
}

main().catch((error) => {
  console.error("Integration smoke test failed:", error.message);
  process.exit(1);
});
