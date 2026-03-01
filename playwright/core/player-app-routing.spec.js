import { test, expect } from '@playwright/test'

test.describe('Player App Routing', () => {
  test.skip(!!process.env.CI, 'Requires production build served by backend')

  test('serves player app at /player/ routes via backend', async ({ page }) => {
    // Use backend URL to test production build routing
    await page.goto('http://localhost:8080/player/game-subscription-instances/test-id/login/test-token')
    
    // Should load the player app (use first() to handle Vue's wrapper div)
    await expect(page.locator('#player-app').first()).toBeVisible()
    
    // Should show support footer
    await expect(page.locator('.player-support-footer')).toBeVisible()
    await expect(page.getByText('Need help?')).toBeVisible()
    await expect(page.getByText('support@playbymail.games')).toBeVisible()
  })

  test('player app has separate bundle from main app', async ({ page }) => {
    // Use backend URL to test production build routing
    await page.goto('http://localhost:8080/player/game-subscription-instances/test-id/login/test-token')
    
    // Wait for page to load
    await page.waitForLoadState('domcontentloaded')
    
    // Check that player bundle is loaded (look for player-*.js pattern)
    const html = await page.content()
    expect(html).toMatch(/assets\/player-[a-zA-Z0-9_-]+\.js/)
    expect(html).toContain('player-app')
  })

  test('main app still works at root via backend', async ({ page }) => {
    // Use backend URL to test production build routing
    await page.goto('http://localhost:8080/')
    
    // Should load the main app
    await expect(page.locator('#app').first()).toBeVisible()
    
    // Should not have player support footer
    await expect(page.locator('.player-support-footer')).not.toBeVisible()
  })
})
