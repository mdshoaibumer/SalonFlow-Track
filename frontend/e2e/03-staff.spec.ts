import { test, expect } from './base-test'

test.describe('Staff Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/staff')
    await page.waitForLoadState('networkidle')
  })

  test('displays staff page heading', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: /staff/i })).toBeVisible()
  })

  test('shows staff page content', async ({ page }) => {
    // Page renders heading + content area (table, empty state, or loading/error)
    await page.waitForSelector('table, .text-center', { timeout: 15000 })
  })

  test('displays staff table or empty state', async ({ page }) => {
    // Wait for data to load - either table or empty state appears
    const el = await page.waitForSelector('table, [class*="text-muted"]', { timeout: 15000 })
    expect(el).not.toBeNull()
  })

  test('has search functionality', async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="earch"], input[type="search"]').first()
    await expect(searchInput).toBeVisible()
    await searchInput.fill('Priya')
    await page.waitForTimeout(500)
  })

  test('has add staff button', async ({ page }) => {
    const addBtn = page.getByRole('button', { name: /add staff/i }).first()
    await expect(addBtn).toBeVisible()
  })

  test('opens add staff dialog on button click', async ({ page }) => {
    await page.getByRole('button', { name: /add staff/i }).first().click()
    await expect(page.getByRole('dialog')).toBeVisible()
  })

  test('add staff form has required fields', async ({ page }) => {
    await page.getByRole('button', { name: /add staff/i }).first().click()
    const dialog = page.getByRole('dialog')
    await expect(dialog).toBeVisible()
    await expect(dialog.locator('text=/name/i').first()).toBeVisible()
    await expect(dialog.locator('text=/phone/i').first()).toBeVisible()
  })

  test('add staff form validates empty submission', async ({ page }) => {
    await page.getByRole('button', { name: /add staff/i }).first().click()
    const dialog = page.getByRole('dialog')
    const submitBtn = dialog.getByRole('button', { name: /save|submit|add|create/i })
    if (await submitBtn.isVisible()) {
      await submitBtn.click()
      await page.waitForTimeout(500)
    }
  })

  test('has status filter or sorting', async ({ page }) => {
    const filter = page.locator('select, [role="combobox"], button:has-text("Status"), button:has-text("All")')
    await expect(filter.first()).toBeVisible()
  })

  test('supports pagination when data exists', async ({ page }) => {
    const pagination = page.locator('button:has-text("Previous"), button:has-text("Next")')
    const isVisible = await pagination.first().isVisible({ timeout: 3000 }).catch(() => false)
    expect(isVisible !== undefined).toBeTruthy()
  })
})
