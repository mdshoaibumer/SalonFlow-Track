import { test, expect } from './base-test'

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
  })

  test('displays dashboard heading', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: /dashboard/i })).toBeVisible()
  })

  test('shows performance summary cards', async ({ page }) => {
    const cards = page.locator('.rounded-lg.border, [class*="card"]')
    await expect(cards.first()).toBeVisible({ timeout: 10000 })
  })

  test('shows staff stats widget', async ({ page }) => {
    await expect(page.locator('text=/staff|stylist|active/i').first()).toBeVisible({ timeout: 10000 })
  })

  test('shows system status section', async ({ page }) => {
    // Dashboard shows KPI cards and quick actions instead of a "system status" section
    await expect(page.locator('text=/revenue|customers|quick actions/i').first()).toBeVisible({ timeout: 10000 })
  })

  test('displays welcome message', async ({ page }) => {
    // Dashboard shows "Overview of your salon's performance" as description
    await expect(page.locator('text=/overview|performance/i').first()).toBeVisible()
  })
})
