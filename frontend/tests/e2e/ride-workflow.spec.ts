import { expect, test, type Page } from "@playwright/test";

const riderSession = {
  access_token: "rider-token",
  user_id: "rider-1",
  email: "rider@example.com",
  role: "rider",
};

const driverSession = {
  access_token: "driver-token",
  user_id: "driver-1",
  email: "driver@example.com",
  role: "driver",
};

const ride = {
  id: "ride-1",
  rider_id: "rider-1",
  pickup: { latitude: 41.0369, longitude: 28.985 },
  dropoff: { latitude: 41.043, longitude: 29.009 },
  status: "requested",
};

test("runs the main ride workflow with mocked backend responses", async ({ page }) => {
  await mockApi(page);

  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Trip operations" })).toBeVisible();
  await expect(page.getByText("needs sessions")).toBeVisible();

  await page.getByRole("button", { name: "Register" }).click();
  await expect(page).toHaveURL(/\/rider$/);
  await expect(page.getByText("Rider registered")).toBeVisible();
  await expect(page.getByText("rider-1")).toBeVisible();
  await expect(page.getByText("Rider session required")).toBeHidden();

  await page.getByRole("button", { name: "Taksim" }).first().click();
  await expect(page.getByLabel("Latitude").first()).toHaveValue("41.0369");

  await page.getByRole("button", { name: "Estimate" }).click();
  await expect(page.getByText("89.41 TRY").first()).toBeVisible();

  await page.getByRole("button", { name: "Request ride" }).click();
  await expect(page.getByText("Ride requested").first()).toBeVisible();

  await page.getByRole("button", { name: "Driver" }).first().click();
  await page.getByRole("button", { name: "Auth", exact: true }).click();
  await page.getByRole("button", { name: "Driver off" }).click();
  await page.getByRole("button", { name: "Register" }).click();
  await expect(page).toHaveURL(/\/driver$/);

  await page.getByRole("button", { name: "Update" }).click();
  await expect(page.getByText("Driver location updated")).toBeVisible();
  await page.getByRole("button", { name: "Available" }).click();
  await expect(page.getByText("Driver is available")).toBeVisible();

  await page.getByRole("button", { name: "Rider", exact: true }).click();
  await page.getByRole("button", { name: "Authorize" }).click();
  await expect(page.getByText("authorized")).toBeVisible();
  await page.getByRole("button", { name: "Refresh" }).click();
  await expect(page.getByText("accepted").first()).toBeVisible();

  await page.getByRole("button", { name: "Driver", exact: true }).click();
  await page.getByRole("button", { name: "Start" }).click();
  await expect(page.getByText("started").first()).toBeVisible();
  await page.getByRole("button", { name: "Complete" }).click();
  await expect(page.getByText("completed").first()).toBeVisible();

  await page.getByRole("button", { name: "Ride", exact: true }).click();
  await page.getByRole("button", { name: "Refresh payment" }).click();
  await expect(page.getByText("captured").first()).toBeVisible();
});

async function mockApi(page: Page) {
  let currentRide = ride;

  await page.route("**/api/v1/auth/register", async (route) => {
    const body = route.request().postDataJSON() as { role: "rider" | "driver" };
    await route.fulfill({ json: body.role === "driver" ? driverSession : riderSession });
  });

  await page.route("**/api/v1/fare-estimates", async (route) => {
    await route.fulfill({
      json: { distance_km: 8.2, duration_minutes: 22, currency: "TRY", amount: 89.41 },
    });
  });

  await page.route("**/api/v1/rides", async (route) => {
    currentRide = ride;
    await route.fulfill({ json: currentRide });
  });

  await page.route("**/api/v1/payments/authorize", async (route) => {
    currentRide = { ...ride, driver_id: "driver-1", status: "accepted" };
    await route.fulfill({
      json: {
        id: "payment-1",
        ride_id: "ride-1",
        rider_id: "rider-1",
        amount: 89.41,
        currency: "TRY",
        status: "authorized",
      },
    });
  });

  await page.route("**/api/v1/drivers/driver-1/location", async (route) => {
    await route.fulfill({ status: 204 });
  });

  await page.route("**/api/v1/drivers/driver-1/available", async (route) => {
    await route.fulfill({ status: 204 });
  });

  await page.route("**/api/v1/rides/ride-1", async (route) => {
    await route.fulfill({ json: currentRide });
  });

  await page.route("**/api/v1/rides/ride-1/start", async (route) => {
    currentRide = { ...ride, driver_id: "driver-1", status: "started" };
    await route.fulfill({ json: currentRide });
  });

  await page.route("**/api/v1/rides/ride-1/complete", async (route) => {
    currentRide = { ...ride, driver_id: "driver-1", status: "completed" };
    await route.fulfill({ json: currentRide });
  });

  await page.route("**/api/v1/payments/payment-1", async (route) => {
    await route.fulfill({
      json: {
        id: "payment-1",
        ride_id: "ride-1",
        rider_id: "rider-1",
        driver_id: "driver-1",
        amount: 89.41,
        currency: "TRY",
        status: "captured",
      },
    });
  });
}
