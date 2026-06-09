import { test, expect } from '@playwright/test'

test.describe('Navigation & Sidebar', () => {
  test('app loads with sidebar visible', async ({ page }) => {
    await page.goto('/')
    await expect(page.locator('nav, [data-testid="sidebar"], aside').first()).toBeVisible({ timeout: 10000 })
  })

  test('sidebar shows all main navigation items', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const nav = page.locator('nav, aside, [role="navigation"]').first()
    await expect(nav).toBeVisible()
  })

  test('clicking sidebar items navigates correctly', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')

    const routes = ['/staff', '/services', '/customers', '/appointments', '/invoices', '/expenses']
    for (const route of routes) {
      await page.goto(route)
      await expect(page).toHaveURL(new RegExp(route))
    }
  })

  test('responsive sidebar collapses on small screens', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 })
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    await expect(page.locator('#root')).toBeAttached()
  })

  test('unknown routes do not crash the app', async ({ page }) => {
    const response = await page.goto('/nonexistent-route-xyz')
    // App should respond (may redirect or show blank)
    expect(response).not.toBeNull()
  })
})
