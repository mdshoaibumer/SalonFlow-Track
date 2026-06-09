import { test, expect } from './desktop-fixtures'

test.describe('Desktop App - Core Functionality', () => {
  test('app window loads with React frontend', async ({ desktopPage: page }) => {
    // The desktop app should already have the frontend loaded
    await page.waitForLoadState('domcontentloaded')
    await expect(page.locator('#root')).toBeAttached()
  })

  test('sidebar navigation is visible', async ({ desktopPage: page }) => {
    await page.waitForLoadState('networkidle')
    const nav = page.locator('nav, aside, [role="navigation"]').first()
    await expect(nav).toBeVisible({ timeout: 10000 })
  })

  test('API server responds through desktop app', async ({ desktopPage: page }) => {
    // Navigate to dashboard and verify data loads from the embedded API
    const response = await page.request.get('http://localhost:8080/api/v1/health')
    expect(response.status()).toBe(200)
    const body = await response.json()
    expect(body.status).toBe('healthy')
    expect(body.version).toBe('0.2.0')
  })

  test('dashboard page loads', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const heading = page.getByRole('heading', { level: 1 }).first()
    await expect(heading).toBeVisible({ timeout: 10000 })
  })

  test('can navigate to Staff page', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const staffLink = page.locator('a[href*="staff"], button:has-text("Staff")').first()
    await staffLink.click()
    await page.waitForLoadState('networkidle')
    const heading = page.getByRole('heading', { name: /staff/i }).first()
    await expect(heading).toBeVisible({ timeout: 10000 })
  })

  test('can navigate to Services page', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const link = page.locator('a[href*="service"], button:has-text("Services")').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    const heading = page.getByRole('heading', { name: /service/i }).first()
    await expect(heading).toBeVisible({ timeout: 10000 })
  })

  test('can navigate to Customers page', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const link = page.locator('a[href*="customer"], button:has-text("Customers")').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    const heading = page.getByRole('heading', { name: /customer/i }).first()
    await expect(heading).toBeVisible({ timeout: 10000 })
  })

  test('can navigate to Appointments page', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const link = page.locator('a[href*="appointment"], button:has-text("Appointments")').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    const heading = page.getByRole('heading', { name: /appointment/i }).first()
    await expect(heading).toBeVisible({ timeout: 10000 })
  })
})

test.describe('Desktop App - Window Behavior', () => {
  test('window title is SalonFlow Track', async ({ desktopPage: page }) => {
    const title = await page.title()
    expect(title).toContain('SalonFlow')
  })

  test('responsive layout works in desktop window', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    // Desktop app window should be wide enough for sidebar
    const viewport = page.viewportSize()
    if (viewport) {
      expect(viewport.width).toBeGreaterThanOrEqual(1024)
    }
    const sidebar = page.locator('nav, aside').first()
    await expect(sidebar).toBeVisible()
  })
})

test.describe('Desktop App - CRUD Operations', () => {
  test('can open Add Staff dialog', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const staffLink = page.locator('a[href*="staff"], button:has-text("Staff")').first()
    await staffLink.click()
    await page.waitForLoadState('networkidle')
    
    const addButton = page.getByRole('button', { name: /add staff/i })
    await expect(addButton).toBeVisible({ timeout: 10000 })
    await addButton.click()
    
    // Dialog or form should appear
    const dialog = page.locator('[role="dialog"], form').first()
    await expect(dialog).toBeVisible({ timeout: 5000 })
  })

  test('can open Add Service dialog', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const link = page.locator('a[href*="service"], button:has-text("Services")').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    
    const addButton = page.getByRole('button', { name: /add service/i })
    await expect(addButton).toBeVisible({ timeout: 10000 })
    await addButton.click()
    
    const dialog = page.locator('[role="dialog"], form').first()
    await expect(dialog).toBeVisible({ timeout: 5000 })
  })

  test('can open Add Customer dialog', async ({ desktopPage: page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    const link = page.locator('a[href*="customer"], button:has-text("Customers")').first()
    await link.click()
    await page.waitForLoadState('networkidle')
    
    const addButton = page.getByRole('button', { name: /add customer/i })
    await expect(addButton).toBeVisible({ timeout: 10000 })
    await addButton.click()
    
    const dialog = page.locator('[role="dialog"], form').first()
    await expect(dialog).toBeVisible({ timeout: 5000 })
  })
})
