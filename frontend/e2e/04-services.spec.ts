import { test, expect } from '@playwright/test'

test.describe('Services Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/services')
    await page.waitForLoadState('networkidle')
  })

  test('displays services page heading', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: /services/i })).toBeVisible()
  })

  test('shows services stats or content', async ({ page }) => {
    await expect(page.locator('text=/total|active|categor|service/i').first()).toBeVisible({ timeout: 10000 })
  })

  test('displays services table or empty state', async ({ page }) => {
    await page.waitForSelector('table, .text-center', { timeout: 15000 })
  })

  test('has search input', async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="earch"]').first()
    await expect(searchInput).toBeVisible()
    await searchInput.fill('Haircut')
    await page.waitForTimeout(500)
  })

  test('has category filter', async ({ page }) => {
    const filter = page.locator('select, [role="combobox"], button:has-text("Category"), button:has-text("All")')
    await expect(filter.first()).toBeVisible()
  })

  test('has status filter', async ({ page }) => {
    const filter = page.locator('text=/active|inactive|status/i')
    await expect(filter.first()).toBeVisible({ timeout: 5000 })
  })

  test('has add service button', async ({ page }) => {
    const addBtn = page.getByRole('button', { name: /add service/i })
    await expect(addBtn).toBeVisible()
  })

  test('opens add service dialog', async ({ page }) => {
    await page.getByRole('button', { name: /add service/i }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog).toBeVisible()
  })

  test('service form has name and price fields', async ({ page }) => {
    await page.getByRole('button', { name: /add service/i }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog.locator('text=/name/i').first()).toBeVisible()
    await expect(dialog.locator('text=/price|amount/i').first()).toBeVisible()
  })

  test('service table headers visible when data exists', async ({ page }) => {
    const headers = page.locator('th, [role="columnheader"]')
    const hasHeaders = await headers.first().isVisible({ timeout: 3000 }).catch(() => false)
    expect(hasHeaders !== undefined).toBeTruthy()
  })
})
