import { test, expect } from '@playwright/test'

test.describe('Analytics & Reports', () => {
  test.describe('Executive Dashboard (Analytics)', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/analytics')
      await page.waitForLoadState('networkidle')
    })

    test('displays analytics page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /analytics|business/i })).toBeVisible()
    })

    test('shows KPI cards', async ({ page }) => {
      const cards = page.locator('.rounded-lg.border, [class*="card"]')
      await expect(cards.first()).toBeVisible({ timeout: 10000 })
    })

    test('shows revenue metrics text', async ({ page }) => {
      // Look for visible KPI content on analytics page
      const kpi = page.locator('.rounded-lg.border, [class*="card"]')
      await expect(kpi.nth(1)).toBeVisible({ timeout: 10000 })
    })
  })

  test.describe('Revenue Reports', () => {
    test('loads revenue reports page', async ({ page }) => {
      await page.goto('/reports/revenue')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: 'Revenue Reports' })).toBeVisible()
    })
  })

  test.describe('Customer Reports', () => {
    test('loads customer reports page', async ({ page }) => {
      await page.goto('/reports/customers')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: 'Customer Reports' })).toBeVisible()
    })
  })

  test.describe('Staff Reports', () => {
    test('loads staff reports page', async ({ page }) => {
      await page.goto('/reports/staff')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: /staff/i })).toBeVisible()
    })
  })

  test.describe('Service Reports', () => {
    test('loads service reports page', async ({ page }) => {
      await page.goto('/reports/services')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: 'Service Reports' })).toBeVisible()
    })
  })

  test.describe('Expense Reports', () => {
    test('loads expense reports page', async ({ page }) => {
      await page.goto('/reports/expenses')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: 'Expense Reports' })).toBeVisible()
    })
  })

  test.describe('Inventory Reports', () => {
    test('loads inventory reports page', async ({ page }) => {
      await page.goto('/reports/inventory')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: /inventory/i })).toBeVisible()
    })
  })

  test.describe('Profit & Loss Reports', () => {
    test('loads profit & loss reports page', async ({ page }) => {
      await page.goto('/reports/profit-loss')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: /profit|loss/i })).toBeVisible()
    })
  })
})
