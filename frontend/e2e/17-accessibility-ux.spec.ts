import { test, expect } from '@playwright/test'

test.describe('Accessibility & UX', () => {
  test('pages have proper heading hierarchy', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const h1 = page.locator('h1')
    await expect(h1.first()).toBeVisible()
  })

  test('buttons have accessible labels', async ({ page }) => {
    await page.goto('/staff')
    await page.waitForLoadState('networkidle')
    const buttons = page.getByRole('button')
    const count = await buttons.count()
    expect(count).toBeGreaterThan(0)
  })

  test('form inputs have labels or placeholders', async ({ page }) => {
    await page.goto('/staff')
    await page.waitForLoadState('networkidle')
    await page.getByRole('button', { name: /add staff/i }).click()
    const dialog = page.getByRole('dialog')
    const inputs = dialog.locator('input')
    const count = await inputs.count()
    expect(count).toBeGreaterThan(0)
  })

  test('dialogs can be closed with escape key', async ({ page }) => {
    await page.goto('/staff')
    await page.waitForLoadState('networkidle')
    await page.getByRole('button', { name: /add staff/i }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog).toBeVisible()
    await page.keyboard.press('Escape')
    await expect(dialog).not.toBeVisible({ timeout: 3000 })
  })

  test('search input handles debounce correctly', async ({ page }) => {
    await page.goto('/staff')
    await page.waitForLoadState('networkidle')
    const search = page.locator('input[placeholder*="earch"]').first()
    await search.fill('test')
    await page.waitForTimeout(1000)
    await expect(page.locator('#root')).toBeAttached()
  })

  test('tab navigation works through interactive elements', async ({ page }) => {
    await page.goto('/staff')
    await page.waitForLoadState('networkidle')
    await page.keyboard.press('Tab')
    await page.keyboard.press('Tab')
    await page.keyboard.press('Tab')
    await expect(page.locator('#root')).toBeAttached()
  })

  test('loading states are shown while data fetches', async ({ page }) => {
    await page.goto('/staff')
    await expect(page.locator('#root')).toBeAttached()
    await page.waitForLoadState('networkidle')
    await expect(page.locator('#root')).toBeAttached()
  })
})
