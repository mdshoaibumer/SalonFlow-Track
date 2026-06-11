import { test, expect } from './base-test'

test.describe('Desktop App - Core Functionality', () => {
  test('app loads with React frontend', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('domcontentloaded')
    await expect(page.locator('#root')).toBeAttached()
  })

  test('sidebar navigation is visible', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const nav = page.locator('nav, aside, [role="navigation"]').first()
    await expect(nav).toBeVisible({ timeout: 10000 })
  })

  test('dashboard page loads', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    // Dashboard should show heading or KPI cards
    await expect(page.locator('h1, h2').first()).toBeVisible({ timeout: 10000 })
  })

  test('can navigate to Staff page', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    // Expand Management section if collapsed
    const managementBtn = page.locator('button:has-text("Management")')
    if (await managementBtn.isVisible()) await managementBtn.click()
    const staffLink = page.locator('a[href="/staff"]').first()
    await staffLink.click()
    await page.waitForLoadState('networkidle')
    const heading = page.getByRole('heading', { name: /staff/i }).first()
    await expect(heading).toBeVisible({ timeout: 10000 })
  })

  test('can navigate to Services page', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    // Expand Management section if collapsed
    const managementBtn = page.locator('button:has-text("Management")')
    if (await managementBtn.isVisible()) await managementBtn.click()
    const link = page.locator('a[href="/services"]').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    const heading = page.getByRole('heading', { name: /service/i }).first()
    await expect(heading).toBeVisible({ timeout: 10000 })
  })

  test('can navigate to Customers page', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    // Expand Management section if collapsed
    const managementBtn = page.locator('button:has-text("Management")')
    if (await managementBtn.isVisible()) await managementBtn.click()
    const link = page.locator('a[href="/customers"]').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    const heading = page.getByRole('heading', { name: /customer/i }).first()
    await expect(heading).toBeVisible({ timeout: 10000 })
  })
})

test.describe('Desktop App - Page Layout', () => {
  test('title contains SalonFlow', async ({ page }) => {
    await page.goto('/')
    const title = await page.title()
    expect(title).toContain('SalonFlow')
  })

  test('desktop layout has sidebar navigation', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const sidebar = page.locator('nav, aside').first()
    await expect(sidebar).toBeVisible()
  })
})

test.describe('Desktop App - CRUD Operations', () => {
  test('can open Add Staff dialog', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    // Expand Management section if collapsed
    const managementBtn = page.locator('button:has-text("Management")')
    if (await managementBtn.isVisible()) await managementBtn.click()
    const staffLink = page.locator('a[href="/staff"]').first()
    await staffLink.click()
    await page.waitForLoadState('networkidle')
    
    const addButton = page.getByRole('button', { name: /add staff/i }).first()
    await expect(addButton).toBeVisible({ timeout: 10000 })
    await addButton.click()
    
    // Dialog or form should appear
    const dialog = page.locator('[role="dialog"], form').first()
    await expect(dialog).toBeVisible({ timeout: 5000 })
  })

  test('can open Add Service dialog', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    // Expand Management section if collapsed
    const managementBtn = page.locator('button:has-text("Management")')
    if (await managementBtn.isVisible()) await managementBtn.click()
    const link = page.locator('a[href="/services"]').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    
    const addButton = page.getByRole('button', { name: /add service/i })
    await expect(addButton).toBeVisible({ timeout: 10000 })
    await addButton.click()
    
    const dialog = page.locator('[role="dialog"], form').first()
    await expect(dialog).toBeVisible({ timeout: 5000 })
  })

  test('can open Add Customer dialog', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    // Expand Management section if collapsed
    const managementBtn = page.locator('button:has-text("Management")')
    if (await managementBtn.isVisible()) await managementBtn.click()
    const link = page.locator('a[href="/customers"]').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    
    const addButton = page.getByRole('button', { name: /add customer/i })
    await expect(addButton).toBeVisible({ timeout: 10000 })
    await addButton.click()
    
    const dialog = page.locator('[role="dialog"], form').first()
    await expect(dialog).toBeVisible({ timeout: 5000 })
  })
})
