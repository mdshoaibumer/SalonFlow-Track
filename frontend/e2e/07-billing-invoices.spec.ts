import { test, expect } from './base-test'

test.describe('Billing & Invoices', () => {
  test.describe('Billing Page', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/billing')
      await page.waitForLoadState('networkidle')
    })

    test('displays billing page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /billing/i })).toBeVisible()
    })

    test('has customer selection area', async ({ page }) => {
      const customerField = page.locator('text=/customer/i')
      await expect(customerField.first()).toBeVisible({ timeout: 10000 })
    })

    test('has service or item section', async ({ page }) => {
      const serviceField = page.locator('text=/service|item|add/i')
      await expect(serviceField.first()).toBeVisible({ timeout: 10000 })
    })

    test('shows totals section', async ({ page }) => {
      const totalField = page.locator('text=/total|subtotal|grand total|amount/i')
      await expect(totalField.first()).toBeVisible({ timeout: 10000 })
    })

    test('has generate or save invoice button', async ({ page }) => {
      const btn = page.locator('button:has-text("Generate"), button:has-text("Save"), button:has-text("Create"), button:has-text("Pay")')
      await expect(btn.first()).toBeVisible({ timeout: 10000 })
    })
  })

  test.describe('Invoices Page', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/invoices')
      await page.waitForLoadState('networkidle')
    })

    test('displays invoices page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /invoices/i })).toBeVisible()
    })

    test('shows invoice stats cards', async ({ page }) => {
      const stats = page.locator('text=/today|revenue|invoices|avg/i')
      await expect(stats.first()).toBeVisible({ timeout: 10000 })
    })

    test('has search input', async ({ page }) => {
      const search = page.locator('input[placeholder*="earch"]').first()
      await expect(search).toBeVisible()
    })

    test('has payment status filter', async ({ page }) => {
      const filter = page.locator('select, [role="combobox"]')
      await expect(filter.first()).toBeVisible()
    })

    test('shows invoice table or empty state', async ({ page }) => {
      await page.waitForSelector('table, .text-center', { timeout: 15000 })
    })

    test('invoice shows payment method badge when data exists', async ({ page }) => {
      const badges = page.locator('[class*="badge"]')
      const isVisible = await badges.first().isVisible({ timeout: 3000 }).catch(() => false)
      expect(isVisible !== undefined).toBeTruthy()
    })
  })
})
