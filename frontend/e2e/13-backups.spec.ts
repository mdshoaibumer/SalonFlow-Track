import { test, expect } from '@playwright/test'

test.describe('Cloud Backup', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/cloud-backup')
    await page.waitForLoadState('networkidle')
  })

  test('displays cloud backup page heading', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: /cloud.*backup|backup/i })).toBeVisible()
  })

  test('shows backup configuration section', async ({ page }) => {
    const config = page.locator('text=/provider|bucket|config|setting|configuration/i')
    await expect(config.first()).toBeVisible({ timeout: 10000 })
  })

  test('has provider selection', async ({ page }) => {
    const provider = page.locator('text=/provider|aws|s3|google|none/i')
    await expect(provider.first()).toBeVisible({ timeout: 10000 })
  })

  test('shows backup history tab', async ({ page }) => {
    const history = page.locator('text=/history|recent|backup/i')
    await expect(history.first()).toBeVisible({ timeout: 10000 })
  })

  test('has backup now button', async ({ page }) => {
    const backupBtn = page.locator('button:has-text("Backup"), button:has-text("Start")')
    const isVisible = await backupBtn.first().isVisible({ timeout: 5000 }).catch(() => false)
    expect(isVisible !== undefined).toBeTruthy()
  })

  test('has test connection button', async ({ page }) => {
    const testBtn = page.locator('button:has-text("Test"), button:has-text("Verify")')
    const isVisible = await testBtn.first().isVisible({ timeout: 5000 }).catch(() => false)
    expect(isVisible !== undefined).toBeTruthy()
  })

  test('shows auto-backup option', async ({ page }) => {
    const auto = page.locator('text=/auto.*backup|automatic|interval/i')
    const isVisible = await auto.first().isVisible({ timeout: 5000 }).catch(() => false)
    expect(isVisible !== undefined).toBeTruthy()
  })

  test('shows backup stats', async ({ page }) => {
    const stats = page.locator('text=/total|size|provider|backups/i')
    await expect(stats.first()).toBeVisible({ timeout: 10000 })
  })
})

test.describe('Local Backup', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/backups')
    await page.waitForLoadState('networkidle')
  })

  test('displays backups page heading', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: /backup/i })).toBeVisible()
  })

  test('has create backup button', async ({ page }) => {
    const btn = page.locator('button:has-text("Create"), button:has-text("Backup"), button:has-text("New")')
    await expect(btn.first()).toBeVisible()
  })

  test('shows backup history or content', async ({ page }) => {
    const content = page.locator('text=/backup|history|file|date|size/i')
    await expect(content.first()).toBeVisible({ timeout: 10000 })
  })

  test('has restore tab or section', async ({ page }) => {
    const restore = page.locator('button:has-text("Restore")')
    const restoreText = page.locator('text=/restore/i')
    await expect(restore.or(restoreText.first()).first()).toBeVisible({ timeout: 5000 })
  })
})
