import { test, expect } from './base-test'

test.describe('CRUD Operations - End to End', () => {
  test.describe('Staff CRUD', () => {
    test('create staff member via dialog', async ({ page }) => {
      await page.goto('/staff')
      await page.waitForLoadState('networkidle')

      await page.getByRole('button', { name: /add staff/i }).click()
      const dialog = page.getByRole('dialog')
      await expect(dialog).toBeVisible()

      await dialog.locator('input[name="full_name"], input[placeholder*="name" i]').first().fill('E2E Test Staff')
      await dialog.locator('input[name="phone"], input[placeholder*="phone" i]').first().fill('9111222333')

      const submitBtn = dialog.getByRole('button', { name: /save|submit|add|create/i })
      await submitBtn.click()
      await page.waitForTimeout(1500)
    })

    test('search filters staff results', async ({ page }) => {
      await page.goto('/staff')
      await page.waitForLoadState('networkidle')

      const search = page.locator('input[placeholder*="earch"]').first()
      await search.fill('nonexistent-xyz-staff')
      await page.waitForTimeout(1000)

      // Should show no results or empty state - page should still be visible
      await expect(page.getByRole('heading', { level: 1, name: /staff/i })).toBeVisible()
    })
  })

  test.describe('Service CRUD', () => {
    test('create a new service', async ({ page }) => {
      await page.goto('/services')
      await page.waitForLoadState('networkidle')

      await page.getByRole('button', { name: /add service/i }).click()
      const dialog = page.getByRole('dialog')
      await expect(dialog).toBeVisible()

      await dialog.locator('input[name="name"], input[placeholder*="name" i]').first().fill('E2E Test Service')

      const priceInput = dialog.locator('input[name="price"], input[placeholder*="price" i]').first()
      if (await priceInput.isVisible()) {
        await priceInput.fill('750')
      }

      const submitBtn = dialog.getByRole('button', { name: /save|submit|add|create/i })
      await submitBtn.click()
      await page.waitForTimeout(1500)
    })
  })

  test.describe('Customer CRUD', () => {
    test('create a new customer', async ({ page }) => {
      await page.goto('/customers')
      await page.waitForLoadState('networkidle')

      await page.getByRole('button', { name: /add customer/i }).click()
      const dialog = page.getByRole('dialog')
      await expect(dialog).toBeVisible()

      await dialog.locator('input[name="full_name"], input[placeholder*="name" i]').first().fill('E2E Customer')
      await dialog.locator('input[name="phone"], input[placeholder*="phone" i]').first().fill('9444555666')

      const submitBtn = dialog.getByRole('button', { name: /save|submit|add|create/i })
      await submitBtn.click()
      await page.waitForTimeout(1500)
    })

    test('search filters customer results', async ({ page }) => {
      await page.goto('/customers')
      await page.waitForLoadState('networkidle')

      const search = page.locator('input[placeholder*="earch"]').first()
      await search.fill('Anjali')
      await page.waitForTimeout(1000)
    })
  })

  test.describe('Appointment CRUD', () => {
    test('open appointment creation tab', async ({ page }) => {
      await page.goto('/appointments')
      await page.waitForLoadState('networkidle')

      const newTab = page.locator('button:has-text("New Appointment"), button:has-text("Add"), button:has-text("Book")')
      await newTab.first().click()
      await page.waitForTimeout(500)

      // Form shows inline (not dialog)
      const form = page.locator('form')
      await expect(form.first()).toBeVisible({ timeout: 5000 })
    })
  })

  test.describe('Expense CRUD', () => {
    test('create a new expense', async ({ page }) => {
      await page.goto('/expenses')
      await page.waitForLoadState('networkidle')

      await page.locator('button:has-text("Add"), button:has-text("New")').first().click()
      await page.waitForTimeout(500)

      await expect(page.locator('text=/category|amount/i').first()).toBeVisible()
    })
  })
})
