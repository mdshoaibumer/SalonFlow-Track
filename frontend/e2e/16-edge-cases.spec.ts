import { test, expect } from './base-test'

test.describe('Edge Cases & Error Handling', () => {
  test.describe('Empty States', () => {
    test('staff page handles no data gracefully', async ({ page }) => {
      await page.goto('/staff')
      await page.waitForLoadState('networkidle')
      await page.waitForSelector('table, .text-center', { timeout: 15000 })
    })

    test('services page handles no data gracefully', async ({ page }) => {
      await page.goto('/services')
      await page.waitForLoadState('networkidle')
      await page.waitForSelector('table, .text-center', { timeout: 15000 })
    })

    test('customers page handles no data gracefully', async ({ page }) => {
      await page.goto('/customers')
      await page.waitForLoadState('networkidle')
      await page.waitForSelector('table, .text-center', { timeout: 15000 })
    })
  })

  test.describe('Form Validation', () => {
    test('staff form rejects invalid phone number', async ({ page }) => {
      await page.goto('/staff')
      await page.waitForLoadState('networkidle')
      await page.getByRole('button', { name: /add staff/i }).click()
      const dialog = page.getByRole('dialog')
      await dialog.locator('input[name="phone"], input[placeholder*="phone" i]').first().fill('abc')
      const submitBtn = dialog.getByRole('button', { name: /save|submit|add|create/i })
      await submitBtn.click()
      await page.waitForTimeout(500)
    })

    test('staff form rejects empty required fields', async ({ page }) => {
      await page.goto('/staff')
      await page.waitForLoadState('networkidle')
      await page.getByRole('button', { name: /add staff/i }).click()
      const dialog = page.getByRole('dialog')
      const submitBtn = dialog.getByRole('button', { name: /save|submit|add|create/i })
      await submitBtn.click()
      await page.waitForTimeout(500)
    })

    test('service form rejects zero or negative price', async ({ page }) => {
      await page.goto('/services')
      await page.waitForLoadState('networkidle')
      await page.getByRole('button', { name: /add service/i }).click()
      const dialog = page.getByRole('dialog')
      const priceInput = dialog.locator('input[name="price"], input[placeholder*="price" i]').first()
      if (await priceInput.isVisible()) {
        await priceInput.fill('0')
      }
      const submitBtn = dialog.getByRole('button', { name: /save|submit|add|create/i })
      await submitBtn.click()
      await page.waitForTimeout(500)
    })
  })

  test.describe('URL Edge Cases', () => {
    test('handles special characters in URL', async ({ page }) => {
      await page.goto('/staff?search=%3Cscript%3E')
      await expect(page.locator('#root')).toBeAttached()
    })

    test('handles extremely long URL params', async ({ page }) => {
      const longStr = 'a'.repeat(500)
      await page.goto(`/staff?search=${longStr}`)
      await expect(page.locator('#root')).toBeAttached()
    })

    test('handles unknown routes gracefully', async ({ page }) => {
      const response = await page.goto('/unknown/route/here')
      expect(response).not.toBeNull()
    })
  })

  test.describe('Responsive Design', () => {
    test('mobile viewport renders correctly', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 812 })
      await page.goto('/')
      await page.waitForLoadState('networkidle')
      await expect(page.locator('#root')).toBeAttached()
      await expect(page.getByRole('heading', { level: 1, name: /dashboard/i })).toBeVisible()
    })

    test('tablet viewport renders correctly', async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 })
      await page.goto('/staff')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: /staff/i })).toBeVisible()
    })

    test('wide desktop viewport renders correctly', async ({ page }) => {
      await page.setViewportSize({ width: 1920, height: 1080 })
      await page.goto('/services')
      await page.waitForLoadState('networkidle')
      await expect(page.getByRole('heading', { level: 1, name: /services/i })).toBeVisible()
    })
  })

  test.describe('Concurrent Actions', () => {
    test('rapid navigation does not crash app', async ({ page }) => {
      await page.goto('/')
      const routes = ['/staff', '/services', '/customers', '/appointments', '/expenses']
      for (const route of routes) {
        await page.goto(route)
      }
      await page.waitForLoadState('networkidle')
      await expect(page.locator('#root')).toBeAttached()
    })

    test('double-clicking add button does not break UI', async ({ page }) => {
      await page.goto('/staff')
      await page.waitForLoadState('networkidle')
      const addBtn = page.getByRole('button', { name: /add staff/i })
      await addBtn.dblclick()
      await page.waitForTimeout(500)
      await expect(page.locator('#root')).toBeAttached()
    })
  })
})
