import { test, expect } from './base-test'

test.describe('WhatsApp Integration', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/whatsapp')
    await page.waitForLoadState('networkidle')
  })

  test('displays whatsapp page heading', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: /whatsapp/i })).toBeVisible()
  })

  test('shows templates tab or section', async ({ page }) => {
    const templates = page.locator('text=/template/i')
    await expect(templates.first()).toBeVisible({ timeout: 10000 })
  })

  test('has new template button', async ({ page }) => {
    const addBtn = page.locator('button:has-text("New Template"), button:has-text("Add"), button:has-text("Create")')
    await expect(addBtn.first()).toBeVisible()
  })

  test('clicking new template shows inline form', async ({ page }) => {
    await page.locator('button:has-text("New Template"), button:has-text("Add"), button:has-text("Create")').first().click()
    await page.waitForTimeout(500)
    // Inline form (not dialog)
    const form = page.locator('form')
    await expect(form.first()).toBeVisible({ timeout: 5000 })
  })

  test('template form has name and body fields', async ({ page }) => {
    await page.locator('button:has-text("New Template"), button:has-text("Add"), button:has-text("Create")').first().click()
    await page.waitForTimeout(500)
    await expect(page.locator('text=/name/i').first()).toBeVisible()
    await expect(page.locator('text=/body|message|content/i').first()).toBeVisible()
  })

  test('shows message stats or info', async ({ page }) => {
    const stats = page.locator('text=/sent|delivered|message|template|log/i')
    await expect(stats.first()).toBeVisible({ timeout: 10000 })
  })

  test('has send message tab', async ({ page }) => {
    const sendTab = page.locator('button:has-text("Send Message")')
    await expect(sendTab).toBeVisible({ timeout: 5000 })
  })

  test('has automation rules tab', async ({ page }) => {
    const rulesTab = page.locator('button:has-text("Automation Rules")')
    await expect(rulesTab).toBeVisible({ timeout: 5000 })
  })

  test('send message tab shows form when clicked', async ({ page }) => {
    const sendTab = page.locator('button:has-text("Send Message")')
    if (await sendTab.isVisible({ timeout: 3000 })) {
      await sendTab.click()
      await page.waitForTimeout(500)
      const content = page.locator('text=/phone|recipient|message|send/i')
      await expect(content.first()).toBeVisible({ timeout: 5000 })
    }
  })
})
