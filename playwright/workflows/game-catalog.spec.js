import { test, expect } from '@playwright/test'
import {
  navigateTo,
  checkPageTitle,
  checkPageURL,
  checkElementVisible,
  takeScreenshot
} from '../utils/test-helpers.js'

test.describe('Game Catalog', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
  })

  test('displays published games with details', async ({ page }) => {
    await navigateTo(page, '/games')
    await checkPageTitle(page, 'Play by Mail')
    await checkPageURL(page, '/games')

    const catalogGames = page.locator('[data-testid="catalog-games"]')
    await expect(catalogGames).toBeVisible({ timeout: 10000 })

    const cards = catalogGames.locator('.catalog-game')
    await expect(cards).toHaveCount(await cards.count())
    expect(await cards.count()).toBeGreaterThanOrEqual(1)

    const firstCard = cards.first()
    await expect(firstCard.locator('.game-name')).toBeVisible()
    await expect(firstCard.locator('.game-name')).not.toBeEmpty()
    const badgeText = await firstCard.locator('.badge').textContent()
    expect(['Adventure', 'Mecha']).toContain(badgeText?.trim())
    await expect(firstCard.locator('.turn-duration')).toBeVisible()
    await expect(firstCard.locator('.join-button')).toBeVisible()
    await expect(firstCard.locator('.join-button')).toHaveAttribute('href', /^\/player\/join-game\//)

    await takeScreenshot(page, 'game-catalog')
  })

  test('does not show loading or empty state when games exist', async ({ page }) => {
    await navigateTo(page, '/games')

    const catalogGames = page.locator('[data-testid="catalog-games"]')
    await expect(catalogGames).toBeVisible({ timeout: 10000 })

    await expect(page.locator('[data-testid="catalog-loading"]')).not.toBeVisible()
    await expect(page.locator('[data-testid="catalog-empty"]')).not.toBeVisible()
  })
})

test.describe('Join Game Turn Sheet', () => {
  test('delivery method toggle shows and hides address fields', async ({ page }) => {
    await navigateTo(page, '/games')

    const catalogGames = page.locator('[data-testid="catalog-games"]')
    await expect(catalogGames).toBeVisible({ timeout: 10000 })

    // Find The Desert Kingdom card (it has email + post delivery options).
    // Use .first() since multiple instances may exist in seed data.
    const desertCard = catalogGames.locator('.catalog-game', {
      has: page.locator('.game-name', { hasText: 'The Desert Kingdom' })
    }).first()
    await expect(desertCard).toBeVisible()

    // Click its Join Game link
    const joinHref = await desertCard.locator('.join-button').getAttribute('href')
    await page.goto(joinHref)
    await page.waitForLoadState('load')

    // Wait for the turn sheet iframe to load
    const iframe = page.frameLocator('[data-testid="join-sheet-iframe"]')
    await expect(iframe.locator('h1')).toContainText('The Desert Kingdom', { timeout: 10000 })

    // Verify delivery method radio buttons are present
    const emailRadio = iframe.locator('input[name="delivery_method"][value="email"]')
    const postRadio = iframe.locator('input[name="delivery_method"][value="post"]')
    await expect(emailRadio).toBeVisible()
    await expect(postRadio).toBeVisible()

    // Address fields should be hidden initially (script runs toggle(false) on load)
    const addressSection = iframe.locator('#postal-address-fields')
    await expect(addressSection).toBeHidden()

    // Select "By Post" -- address fields should appear
    await postRadio.check()
    await expect(addressSection).toBeVisible()
    await expect(iframe.locator('#postal_address_line1')).toBeVisible()
    await expect(iframe.locator('#country')).toBeVisible()

    // Switch to "By Email" -- address fields should hide again
    await emailRadio.check()
    await expect(addressSection).toBeHidden()

    await takeScreenshot(page, 'join-game-delivery-toggle')
  })
})
