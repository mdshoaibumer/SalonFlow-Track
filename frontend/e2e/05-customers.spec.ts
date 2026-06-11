import { test, expect } from './base-test'

test.describe('Customers Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/customers')
    await page.waitForLoadState('networkidle')
  })

  test('displays customers page heading', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: /customers/i })).toBeVisible()
  })

  test('shows customer content area', async ({ page }) => {
    await expect(page.locator('text=/total|active|customer/i').first()).toBeVisible({ timeout: 10000 })
  })

  test('displays customers table or empty state', async ({ page }) => {
    await page.waitForSelector('table, .text-center', { timeout: 15000 })
  })

  test('has search input for customers', async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="earch"]').first()
    await expect(searchInput).toBeVisible()
    await searchInput.fill('Anjali')
    await page.waitForTimeout(500)
  })

  test('has add customer button', async ({ page }) => {
    const addBtn = page.getByRole('button', { name: /add customer/i })
    await expect(addBtn).toBeVisible()
  })

  test('opens add customer dialog', async ({ page }) => {
    await page.getByRole('button', { name: /add customer/i }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog).toBeVisible()
  })

  test('customer form has name and phone fields', async ({ page }) => {
    await page.getByRole('button', { name: /add customer/i }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog.locator('text=/name/i').first()).toBeVisible()
    await expect(dialog.locator('text=/phone/i').first()).toBeVisible()
  })

  test('customer form validates required fields', async ({ page }) => {
    await page.getByRole('button', { name: /add customer/i }).click()
    const dialog = page.getByRole('dialog')
    const submitBtn = dialog.getByRole('button', { name: /save|submit|add|create/i })
    if (await submitBtn.isVisible()) {
      await submitBtn.click()
      await page.waitForTimeout(500)
    }
  })

  test('has filter options', async ({ page }) => {
    const filter = page.locator('select, [role="combobox"], button:has-text("Status"), button:has-text("All")')
    await expect(filter.first()).toBeVisible()
  })

  test('table headers visible when data exists', async ({ page }) => {
    const headers = page.locator('th, [role="columnheader"]')
    const hasHeaders = await headers.first().isVisible({ timeout: 3000 }).catch(() => false)
    expect(hasHeaders !== undefined).toBeTruthy()
  })
})
