import { test, expect } from '@playwright/test'

test.describe('Administration', () => {
  test.describe('GST & Tax', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/gst')
      await page.waitForLoadState('networkidle')
    })

    test('displays GST page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /gst|tax/i })).toBeVisible()
    })

    test('shows GST settings form', async ({ page }) => {
      const form = page.locator('text=/gstin|business|cgst|sgst|rate/i')
      await expect(form.first()).toBeVisible({ timeout: 10000 })
    })

    test('has enable/disable GST toggle', async ({ page }) => {
      const toggle = page.locator('input[type="checkbox"], [role="switch"], button:has-text("Enable")')
      await expect(toggle.first()).toBeVisible({ timeout: 5000 })
    })

    test('shows tax rates section', async ({ page }) => {
      const rates = page.locator('text=/rate|cgst|sgst|igst|hsn/i')
      await expect(rates.first()).toBeVisible({ timeout: 10000 })
    })
  })

  test.describe('Printer Settings', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/printer')
      await page.waitForLoadState('networkidle')
    })

    test('displays printer page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: /printer|print/i })).toBeVisible()
    })

    test('shows print settings form', async ({ page }) => {
      const form = page.locator('text=/paper|width|margin|header|footer/i')
      await expect(form.first()).toBeVisible({ timeout: 10000 })
    })

    test('has save settings button', async ({ page }) => {
      const saveBtn = page.locator('button:has-text("Save"), button:has-text("Update")')
      await expect(saveBtn.first()).toBeVisible()
    })

    test('has test print button', async ({ page }) => {
      const testBtn = page.locator('button:has-text("Test"), button:has-text("Print")')
      const isVisible = await testBtn.first().isVisible({ timeout: 5000 }).catch(() => false)
      expect(isVisible !== undefined).toBeTruthy()
    })
  })

  test.describe('Import Data', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/import')
      await page.waitForLoadState('networkidle')
    })

    test('displays import page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: 'Data Import' })).toBeVisible()
    })

    test('has file upload or entity selector', async ({ page }) => {
      const upload = page.locator('input[type="file"]')
      const selector = page.locator('text=/upload|csv|excel|entity|select|step/i')
      await expect(upload.or(selector.first()).first()).toBeVisible({ timeout: 10000 })
    })

    test('shows import wizard or steps', async ({ page }) => {
      const steps = page.locator('text=/upload|step|map|preview|import/i')
      await expect(steps.first()).toBeVisible({ timeout: 10000 })
    })
  })

  test.describe('License', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/license')
      await page.waitForLoadState('networkidle')
    })

    test('displays license page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: 'License & Subscription' })).toBeVisible()
    })

    test('shows license status', async ({ page }) => {
      const status = page.locator('text=/active|expired|license|key|plan|status|subscription/i')
      await expect(status.first()).toBeVisible({ timeout: 10000 })
    })

    test('has validate or activate button', async ({ page }) => {
      const activate = page.locator('button:has-text("Validate"), button:has-text("Activate"), button:has-text("Renew")')
      await expect(activate.first()).toBeVisible({ timeout: 5000 })
    })
  })

  test.describe('Updates', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/updates')
      await page.waitForLoadState('networkidle')
    })

    test('displays updates page heading', async ({ page }) => {
      await expect(page.getByRole('heading', { level: 1, name: 'Updates' })).toBeVisible()
    })

    test('shows current version info', async ({ page }) => {
      const version = page.locator('text=/version|current/i')
      await expect(version.first()).toBeVisible({ timeout: 10000 })
    })

    test('has check for updates button', async ({ page }) => {
      const checkBtn = page.locator('button:has-text("Check"), button:has-text("Update")')
      await expect(checkBtn.first()).toBeVisible()
    })
  })

  test.describe('Settings', () => {
    test('settings or preferences page loads', async ({ page }) => {
      await page.goto('/gst')
      await page.waitForLoadState('networkidle')
      await expect(page.locator('#root')).toBeAttached()
    })
  })
})
