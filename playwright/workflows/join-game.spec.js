import { test, expect } from '@playwright/test'

test.describe('Join Game Flow', () => {
  test.skip(!!process.env.CI, 'Requires production build served by backend')

  test('player app loads at /player/join-game/:id route', async ({ page }) => {
    await page.goto('http://localhost:8080/player/join-game/test-subscription-id')

    await expect(page.locator('#player-app').first()).toBeVisible()

    await expect(page.locator('.player-support-footer')).toBeVisible()
    await expect(page.getByText('Need help?')).toBeVisible()
  })

  test('join game page shows loading then error for unknown subscription', async ({ page }) => {
    await page.goto('http://localhost:8080/player/join-game/00000000-0000-0000-0000-000000000000')

    await page.waitForLoadState('networkidle')

    const errorLocator = page.locator('[data-testid="join-load-error"]')
    const loadingLocator = page.locator('[data-testid="join-loading"]')

    await expect(loadingLocator).not.toBeVisible({ timeout: 5000 })

    await expect(errorLocator).toBeVisible()
    await expect(page.getByText('Browse other games')).toBeVisible()
  })

  test('join game page uses player app bundle not main app bundle', async ({ page }) => {
    await page.goto('http://localhost:8080/player/join-game/test-subscription-id')
    await page.waitForLoadState('domcontentloaded')

    const html = await page.content()
    expect(html).toMatch(/assets\/player-[a-zA-Z0-9_-]+\.js/)
    expect(html).toContain('player-app')
  })

  test('catalog join game link navigates to player app', async ({ page }) => {
    await page.goto('http://localhost:8080/games')
    await page.waitForLoadState('networkidle')

    const joinButtons = page.locator('[data-testid^="join-button-"]')
    const count = await joinButtons.count()

    if (count > 0) {
      const href = await joinButtons.first().getAttribute('href')
      expect(href).toMatch(/^\/player\/join-game\//)

      await page.goto(`http://localhost:8080${href}`)

      await expect(page.locator('#player-app').first()).toBeVisible()
      await expect(page.locator('.player-support-footer')).toBeVisible()
    } else {
      await expect(page.locator('[data-testid="catalog-empty"]')).toBeVisible()
    }
  })
})
