import { test, expect } from '@playwright/test'

test.describe('Join Game Flow', () => {
  test.skip(!!process.env.CI, 'Requires production build served by backend')

  test('unauthenticated user is redirected to login from join game page', async ({ page }) => {
    await page.goto('http://localhost:8080/player/join-game/test-subscription-id')

    await page.waitForURL('**/login?redirect=**', { timeout: 5000 })

    expect(page.url()).toContain('/login')
    expect(page.url()).toContain('redirect=%2Fplayer%2Fjoin-game%2Ftest-subscription-id')
  })

  test('join game page uses player app bundle not main app bundle', async ({ page }) => {
    await page.goto('http://localhost:8080/player/join-game/test-subscription-id')
    await page.waitForLoadState('domcontentloaded')

    const html = await page.content()
    expect(html).toMatch(/assets\/player-[a-zA-Z0-9_-]+\.js/)
    expect(html).toContain('player-app')
  })

  test('catalog join game link points to player app route', async ({ page }) => {
    await page.goto('http://localhost:8080/games')
    await page.waitForLoadState('networkidle')

    const joinButtons = page.locator('[data-testid^="join-button-"]')
    const count = await joinButtons.count()

    if (count > 0) {
      const href = await joinButtons.first().getAttribute('href')
      expect(href).toMatch(/^\/player\/join-game\//)
    } else {
      await expect(page.locator('[data-testid="catalog-empty"]')).toBeVisible()
    }
  })
})
