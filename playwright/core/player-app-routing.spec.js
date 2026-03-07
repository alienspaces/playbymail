import { test, expect } from '@playwright/test'

test.describe('Player App Routing', () => {
  test.skip(!!process.env.CI, 'Requires production build served by backend')

  test('serves player app at /player/ routes via backend', async ({ page }) => {
    await page.goto('http://localhost:8080/player/game-subscription-instances/test-id/turn-sheets/test-token')

    await expect(page.locator('#player-app').first()).toBeVisible()

    await expect(page.locator('.player-support-footer')).toBeVisible()
    await expect(page.getByText('Need help?')).toBeVisible()
    await expect(page.getByText('support@playbymail.games')).toBeVisible()
  })

  test('player app has separate bundle from main app', async ({ request }) => {
    const response = await request.get('http://localhost:8080/player/game-subscription-instances/test-id/turn-sheets/test-token')
    const html = await response.text()
    expect(html).toMatch(/assets\/player-[a-zA-Z0-9_-]+\.js/)
    expect(html).toContain('player-app')
  })

  test('main app still works at root via backend', async ({ page }) => {
    await page.goto('http://localhost:8080/')

    await expect(page.locator('#app').first()).toBeVisible()

    await expect(page.locator('.player-support-footer')).not.toBeVisible()
  })
})
