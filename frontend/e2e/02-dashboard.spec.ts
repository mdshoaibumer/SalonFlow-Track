import { test, expect } from '@playwright/test'

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
    await expect(page.locator('text=/system status/i').first()).toBeVisible({ timeout: 10000 })
  })

  test('displays welcome message', async ({ page }) => {
    await expect(page.locator('text=/welcome to salonflow/i')).toBeVisible()
  })
})
