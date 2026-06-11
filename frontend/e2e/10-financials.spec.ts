import { test, expect } from './base-test'

test.describe('Financial Modules', () => {
  test.describe('Expenses', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/expenses')
      await page.waitForLoadState('networkidle')
    })

    test('displays expenses page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /expenses/i })).toBeVisible()
    })

    test('shows expense stats or content', async ({ page }) => {
      await expect(page.locator('text=/expense|today|monthly|total/i').first()).toBeVisible({ timeout: 10000 })
    })

    test('has add expense button', async ({ page }) => {
      const addBtn = page.locator('button:has-text("Add"), button:has-text("Create"), button:has-text("New")')
      await expect(addBtn.first()).toBeVisible()
    })

    test('has category filter', async ({ page }) => {
      const filter = page.locator('select, [role="combobox"]')
      await expect(filter.first()).toBeVisible()
    })

    test('shows expense list or table', async ({ page }) => {
      const table = page.locator('table')
      const empty = page.locator('text=/no expense/i')
      await expect(table.or(empty).first()).toBeVisible({ timeout: 10000 })
    })

    test('expense form appears on add click', async ({ page }) => {
      await page.locator('button:has-text("Add"), button:has-text("New")').first().click()
      await page.waitForTimeout(500)
      await expect(page.locator('text=/category|amount|date/i').first()).toBeVisible()
    })
  })

  test.describe('Commissions', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/commissions')
      await page.waitForLoadState('networkidle')
    })

    test('displays commissions page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: 'Commissions' })).toBeVisible()
    })

    test('shows commission content', async ({ page }) => {
      const content = page.locator('text=/rule|commission|percentage|no commission/i')
      await expect(content.first()).toBeVisible({ timeout: 10000 })
    })

    test('has add rule button', async ({ page }) => {
      const addBtn = page.locator('button:has-text("Add"), button:has-text("Create"), button:has-text("New")')
      await expect(addBtn.first()).toBeVisible()
    })
  })

  test.describe('Salary', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/salary')
      await page.waitForLoadState('networkidle')
    })

    test('displays salary page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: 'Salary Management' })).toBeVisible()
    })

    test('shows salary content', async ({ page }) => {
      const content = page.locator('text=/salary|month|generate|staff/i')
      await expect(content.first()).toBeVisible({ timeout: 10000 })
    })

    test('has month or period selector', async ({ page }) => {
      const selector = page.locator('select, input[type="month"], [role="combobox"]')
      await expect(selector.first()).toBeVisible({ timeout: 5000 })
    })
  })

  test.describe('Advances', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/advances')
      await page.waitForLoadState('networkidle')
    })

    test('displays advances page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /advance/i })).toBeVisible()
    })

    test('has add advance button', async ({ page }) => {
      const addBtn = page.locator('button:has-text("Add"), button:has-text("Create"), button:has-text("New"), button:has-text("Give")')
      await expect(addBtn.first()).toBeVisible()
    })

    test('shows advance content', async ({ page }) => {
      const table = page.locator('table')
      const content = page.locator('text=/amount|staff|no advance/i')
      await expect(table.or(content.first()).first()).toBeVisible({ timeout: 10000 })
    })
  })

  test.describe('Profit & Loss', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/profit-loss')
      await page.waitForLoadState('networkidle')
    })

    test('displays profit & loss page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /profit|loss/i })).toBeVisible()
    })

    test('shows revenue and expense summary', async ({ page }) => {
      const content = page.locator('text=/revenue|expense|profit|loss|income/i')
      await expect(content.first()).toBeVisible({ timeout: 10000 })
    })
  })

  test.describe('Performance', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/performance')
      await page.waitForLoadState('networkidle')
    })

    test('displays performance page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: 'Staff Performance' })).toBeVisible()
    })

    test('shows staff performance content', async ({ page }) => {
      const content = page.locator('text=/staff|revenue|service|daily|monthly/i')
      await expect(content.first()).toBeVisible({ timeout: 10000 })
    })
  })
})
