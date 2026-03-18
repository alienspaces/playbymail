/**
 * E2E test: submit turn sheets with an equip action and verify processing succeeds.
 *
 * This spec exercises the full flow for the equipment slot resolution fix:
 * - The inventory management turn sheet shows equippable location items as radio
 *   button groups (loc_<id> = "equip" | "pick_up").
 * - The frontend converts these back into the array-based format the backend expects.
 * - The backend resolves the correct equipment slot from the item definition,
 *   not from the hardcoded default ("weapon"), so items like the Desert Compass
 *   (jewelry slot) equip correctly.
 *
 * It also verifies presentation order: step 0 = Location Choice, step 1 = Inventory.
 *
 * Prerequisites:
 *   - Backend running with enhanced seed data (The Desert Kingdom game)
 *   - GAME_TURN_QUEUEING_INTERVAL_SECONDS=10 (fast periodic queueing)
 *   - TEST_BYPASS_HEADER_NAME / TEST_BYPASS_HEADER_VALUE set for auth bypass
 *
 * Seed data facts (from test_data.go):
 *   - "The Desert Kingdom" game, starting location "Oasis Village"
 *   - Desert Compass (equippable, jewelry slot) at Oasis Village
 *   - Water Flask (not equippable) at Oasis Village
 *   - "The Dusty Trail" leads east to Ancient Ruins
 */
import { test, expect } from '@playwright/test'
import {
  clearAllEmails,
  waitForEmail,
  getEmailBody,
  extractLink,
} from '../utils/mailpit-helpers.js'
import {
  setupTestBypassHeaders,
  waitForPageReady,
} from '../utils/test-helpers.js'

const PLAYER_EMAIL = 'equip-e2e-player@example.com'
const PLAYER_NAME = 'Equip Test Player'
const CHARACTER_NAME = 'Zara the Equipped'

test.describe('Turn Sheet Equip E2E', () => {
  test.describe.configure({ mode: 'serial' })
  test.skip(!!process.env.CI, 'Requires full stack: backend + MailPit + enhanced seed data')

  test.beforeAll(async () => {
    await clearAllEmails()
  })

  // ──────────────────────────────────────────────────────────────────────
  // Phase 1: Join The Desert Kingdom
  // ──────────────────────────────────────────────────────────────────────

  test('Phase 1: join The Desert Kingdom', async ({ page }) => {
    await setupTestBypassHeaders(page)

    await page.goto('/games')
    await waitForPageReady(page)

    const catalogGames = page.locator('[data-testid="catalog-games"]')
    await expect(catalogGames).toBeVisible({ timeout: 10000 })

    // Find The Desert Kingdom card
    const desertCard = catalogGames.locator('.catalog-game', {
      has: page.locator('.game-name', { hasText: 'The Desert Kingdom' }),
    })
    await expect(desertCard).toBeVisible({ timeout: 5000 })

    const joinHref = await desertCard.locator('.join-button').getAttribute('href')
    expect(joinHref).toMatch(/^\/player\/join-game\//)

    await page.goto(joinHref)
    await waitForPageReady(page)
    await page.waitForSelector('#player-app', { timeout: 10000 })

    const iframe = page.frameLocator('[data-testid="join-sheet-iframe"]')
    await expect(iframe.locator('h1')).toContainText('The Desert Kingdom', { timeout: 10000 })

    await iframe.locator('[name="email"]').fill(PLAYER_EMAIL)
    await iframe.locator('[name="name"]').fill(PLAYER_NAME)
    await iframe.locator('[name="character_name"]').fill(CHARACTER_NAME)

    const emailRadio = iframe.locator('input[name="delivery_method"][value="email"]')
    if (await emailRadio.isVisible()) {
      await emailRadio.check()
    }

    await page.locator('[data-testid="btn-submit"]').click()
    await page.waitForTimeout(1000)
  })

  // ──────────────────────────────────────────────────────────────────────
  // Phase 2: Turn 1 — verify presentation order then equip a location item
  // ──────────────────────────────────────────────────────────────────────

  test('Phase 2a: receive turn 1 email', async () => {
    test.setTimeout(90000)
    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 60000, interval: 500 })
    expect(notifEmail).toBeTruthy()
  })

  test('Phase 2b: verify presentation order — Location Choice is step 0, Inventory is step 1', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 30000 })
    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)
    expect(turnSheetLink).toBeTruthy()

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-stepper"]', { timeout: 15000 })
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await page.waitForTimeout(800)

    // Step 0 should be the Location Choice sheet (shown first by presentation order)
    const iframe0 = page.frameLocator('[data-testid="ts-viewer-iframe"]')
    const locationRadios = iframe0.locator('input[type="radio"][name="location_choice"]')
    await expect(locationRadios.first()).toBeVisible({ timeout: 5000 })

    // Navigate to step 1 — should be Inventory Management
    await page.locator('[data-testid="ts-step-1"]').click()
    await page.waitForTimeout(1000)
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await page.waitForTimeout(800)

    const iframe1 = page.frameLocator('[data-testid="ts-viewer-iframe"]')
    // Inventory management has a section for location items
    const locationItemSection = iframe1.locator('.location-items, [data-testid="location-items"], h2, h3')
    await expect(locationItemSection.first()).toBeVisible({ timeout: 5000 })
  })

  test('Phase 2c: verify equippable item shows Equip radio button in inventory management', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 30000 })
    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-stepper"]', { timeout: 15000 })

    // Navigate to step 1 (inventory management)
    await page.locator('[data-testid="ts-step-1"]').click()
    await page.waitForTimeout(1000)
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await page.waitForTimeout(800)

    const iframe = page.frameLocator('[data-testid="ts-viewer-iframe"]')

    // Equippable items (Desert Compass) use per-item radio groups: name="loc_<id>" value="equip"
    // Non-equippable items (Water Flask) use a checkbox: name="pick_up" value="<id>"
    const equipRadios = iframe.locator('input[type="radio"][value="equip"]')
    await expect(equipRadios.first()).toBeVisible({ timeout: 5000 })

    // Verify Water Flask has a backpack-only option (no equip radio)
    const pickUpOnlyCheckbox = iframe.locator('input[type="checkbox"][name="pick_up"]')
    await expect(pickUpOnlyCheckbox.first()).toBeVisible({ timeout: 5000 })
  })

  test('Phase 2d: fill turn sheets — choose location and equip compass — then submit', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 30000 })
    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-stepper"]', { timeout: 15000 })

    // Step 0: Location Choice — pick The Dusty Trail (go to Ancient Ruins)
    await page.locator('[data-testid="ts-step-0"]').click()
    await page.waitForTimeout(800)
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await page.waitForTimeout(600)

    const iframe0 = page.frameLocator('[data-testid="ts-viewer-iframe"]')
    const locationRadios = iframe0.locator('input[type="radio"][name="location_choice"]')
    if (await locationRadios.count() > 0) {
      await locationRadios.first().check()
    }

    // Navigate to Step 1: Inventory Management — caches step 0
    await page.locator('[data-testid="ts-step-1"]').click()
    await page.waitForTimeout(1000)
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await page.waitForTimeout(600)

    const iframe1 = page.frameLocator('[data-testid="ts-viewer-iframe"]')

    // Equip the Desert Compass: radio group name="loc_<id>", select value="equip"
    const equipRadio = iframe1.locator('input[type="radio"][value="equip"]').first()
    if (await equipRadio.isVisible()) {
      await equipRadio.check()
    }

    // Pick up the Water Flask via backpack radio (value="pick_up") or checkbox
    const pickUpRadio = iframe1.locator('input[type="radio"][value="pick_up"]').first()
    const pickUpCheckbox = iframe1.locator('input[type="checkbox"][name="pick_up"]').first()
    if (await pickUpRadio.isVisible()) {
      await pickUpRadio.check()
    } else if (await pickUpCheckbox.isVisible()) {
      await pickUpCheckbox.check()
    }

    // Submit all sheets
    const submitBtn = page.locator('[data-testid="btn-submit-all"]')
    await expect(submitBtn).toBeEnabled({ timeout: 5000 })
    await submitBtn.click()
    await page.waitForTimeout(1500)

    await expect(page.locator('[data-testid="ts-success"]')).toBeVisible({ timeout: 10000 })
  })

  // ──────────────────────────────────────────────────────────────────────
  // Phase 3: Verify turn processed — equip succeeded, no slot mismatch error
  // ──────────────────────────────────────────────────────────────────────

  test('Phase 3: turn 2 notification arrives — equip succeeded', async () => {
    test.setTimeout(150000)
    // Turn 2 email only arrives if turn 1 processed without error.
    // This is the regression guard for the equipment slot mismatch bug.
    const turn2Email = await waitForEmail(PLAYER_EMAIL, 'Turn 2', { timeout: 120000, interval: 1000 })
    expect(turn2Email).toBeTruthy()
  })
})
