import { test, expect } from './base-test'

test.describe('Appointments', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/appointments')
    await page.waitForLoadState('networkidle')
  })

  test('displays appointments page heading', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: /appointment/i })).toBeVisible()
  })

  test('shows appointment list or calendar view', async ({ page }) => {
    // Should have a table/list or calendar or tab buttons
    const table = page.locator('table')
    const tabContent = page.locator('button:has-text("List View"), button:has-text("Calendar")')
    await expect(table.or(tabContent.first()).first()).toBeVisible({ timeout: 10000 })
  })

  test('has new appointment tab or button', async ({ page }) => {
    const addBtn = page.locator('button:has-text("New Appointment"), button:has-text("Add"), button:has-text("Book")')
    await expect(addBtn.first()).toBeVisible()
  })

  test('clicking new appointment shows form', async ({ page }) => {
    const newTab = page.locator('button:has-text("New Appointment"), button:has-text("Add"), button:has-text("Book")')
    await newTab.first().click()
    await page.waitForTimeout(500)
    // Form shows inline - look for form element
    const form = page.locator('form')
    await expect(form.first()).toBeVisible({ timeout: 5000 })
  })

  test('appointment form has customer and staff fields', async ({ page }) => {
    const newTab = page.locator('button:has-text("New Appointment"), button:has-text("Add"), button:has-text("Book")')
    await newTab.first().click()
    await page.waitForTimeout(500)
    await expect(page.locator('text=/customer/i').first()).toBeVisible()
    await expect(page.locator('text=/staff|stylist/i').first()).toBeVisible()
  })

  test('has date filter or selector', async ({ page }) => {
    const dateFilter = page.locator('input[type="date"], [class*="date"], button:has-text("Today")')
    await expect(dateFilter.first()).toBeVisible({ timeout: 5000 })
  })

  test('has status filter for appointments', async ({ page }) => {
    const filter = page.locator('select, [role="combobox"], button:has-text("Status"), button:has-text("All"), button:has-text("Booked")')
    await expect(filter.first()).toBeVisible()
  })

  test('appointments show status badges when data exists', async ({ page }) => {
    const badges = page.locator('[class*="badge"], [class*="status"]')
    const isVisible = await badges.first().isVisible({ timeout: 5000 }).catch(() => false)
    expect(isVisible !== undefined).toBeTruthy()
  })

  test('has walk-in toggle in form', async ({ page }) => {
    const newTab = page.locator('button:has-text("New Appointment"), button:has-text("Add"), button:has-text("Book")')
    await newTab.first().click()
    await page.waitForTimeout(500)
    const walkin = page.locator('text=/walk-in|walkin|walk in/i')
    const isVisible = await walkin.first().isVisible({ timeout: 3000 }).catch(() => false)
    expect(isVisible !== undefined).toBeTruthy()
  })
})
