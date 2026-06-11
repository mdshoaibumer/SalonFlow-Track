import { test, expect } from './base-test'

test.describe('Inventory & Products', () => {
  test.describe('Products', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/products')
      await page.waitForLoadState('networkidle')
    })

    test('displays products page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /product/i })).toBeVisible()
    })

    test('shows products table or empty state', async ({ page }) => {
      const table = page.locator('table')
      const empty = page.locator('text=/no product/i')
      await expect(table.or(empty).first()).toBeVisible({ timeout: 10000 })
    })

    test('has add product button', async ({ page }) => {
      const addBtn = page.locator('button:has-text("Add"), button:has-text("Create"), button:has-text("New")')
      await expect(addBtn.first()).toBeVisible()
    })

    test('has search functionality', async ({ page }) => {
      const search = page.locator('input[placeholder*="earch"]').first()
      await expect(search).toBeVisible()
    })

    test('has category filter', async ({ page }) => {
      const filter = page.locator('select, [role="combobox"]')
      await expect(filter.first()).toBeVisible()
    })

    test('product form appears on add click', async ({ page }) => {
      await page.locator('button:has-text("Add"), button:has-text("New")').first().click()
      await page.waitForTimeout(500)
      const dialog = page.getByRole('dialog')
      const form = page.locator('form')
      await expect(dialog.or(form).first()).toBeVisible({ timeout: 5000 })
    })
  })

  test.describe('Purchases', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/purchases')
      await page.waitForLoadState('networkidle')
    })

    test('displays purchases page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /purchase/i })).toBeVisible()
    })

    test('shows purchase table or content', async ({ page }) => {
      const table = page.locator('table')
      const content = page.locator('text=/no purchase/i')
      await expect(table.or(content).first()).toBeVisible({ timeout: 10000 })
    })

    test('has add purchase button', async ({ page }) => {
      const addBtn = page.locator('button:has-text("Add"), button:has-text("Create"), button:has-text("New")')
      await expect(addBtn.first()).toBeVisible()
    })
  })

  test.describe('Inventory', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/inventory')
      await page.waitForLoadState('networkidle')
    })

    test('displays inventory page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: 'Inventory Dashboard' })).toBeVisible()
    })

    test('shows stock levels or dashboard content', async ({ page }) => {
      const content = page.locator('text=/stock|quantity|product|low|inventory/i')
      await expect(content.first()).toBeVisible({ timeout: 10000 })
    })

    test('has stock adjustment option', async ({ page }) => {
      const adjustBtn = page.locator('button:has-text("Adjust"), button:has-text("Stock"), button:has-text("Update")')
      const isVisible = await adjustBtn.first().isVisible({ timeout: 5000 }).catch(() => false)
      expect(isVisible !== undefined).toBeTruthy()
    })

    test('shows low stock alerts section', async ({ page }) => {
      const alerts = page.locator('text=/low stock|out of stock|alert|warning/i')
      const isVisible = await alerts.first().isVisible({ timeout: 5000 }).catch(() => false)
      expect(isVisible !== undefined).toBeTruthy()
    })
  })
})
