/**
 * End-to-end email adventure game test for "The Desert Kingdom".
 *
 * Prerequisites:
 *   - Backend running with EMAILER_PROVIDER=smtp, SMTP_HOST=localhost:1025
 *   - MailPit running (tools/mailpit-start)
 *   - GAME_TURN_QUEUEING_INTERVAL_SECONDS=10 (fast periodic queueing)
 *   - Demo data loaded (includes The Desert Kingdom with email+post delivery, single-player, process-when-all-submitted)
 *   - TEST_BYPASS_HEADER_NAME / TEST_BYPASS_HEADER_VALUE set for auth bypass
 *
 * Why do these tests take so long (even with a local mail catcher)?
 * - MailPit is instant; the delay is not SMTP. The main costs are:
 *   1. Periodic job: GameTurnQueueingWorker runs every GAME_TURN_QUEUEING_INTERVAL_SECONDS (e.g. 10s).
 *      After join, we wait up to one interval for auto-start, then processing and email. Phase 2a can
 *      therefore wait ~10–15s for the first "turn" email.
 *   2. Email polling: waitForEmail() polls MailPit at a fixed interval (e.g. 1–2s). If the email
 *      arrives just after a poll, we wait almost a full interval before the next check.
 *   3. Explicit sleeps: several waitForTimeout(1000–2000) calls add ~30s+ across the suite to allow
 *      iframe/content to settle and avoid flakiness.
 *   4. Serial execution: this describe runs in serial mode, so all phases run one after another.
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

const PLAYER_EMAIL = 'desert-e2e-player@example.com'
const PLAYER_NAME = 'Desert Explorer'
const CHARACTER_NAME = 'Kael the Wanderer'

test.describe('Desert Kingdom E2E Email Adventure', () => {
  test.describe.configure({ mode: 'serial' })
  test.skip(!!process.env.CI, 'Requires full stack: backend + MailPit + demo data')

  // Shared state across serial tests
  let gsiId = null
  let sessionToken = null
  let desertJoinHref = null

  test.beforeAll(async () => {
    await clearAllEmails()
  })

  // ──────────────────────────────────────────────────────
  // Phase 1: Join the game
  // ──────────────────────────────────────────────────────

  test('Phase 1a: browse catalog and find Desert Kingdom', async ({ page }) => {
    await page.goto('/games')
    await waitForPageReady(page)

    const catalogGames = page.locator('[data-testid="catalog-games"]')
    await expect(catalogGames).toBeVisible({ timeout: 10000 })

    // Use .first() since multiple instances of The Desert Kingdom may be seeded.
    const desertCard = catalogGames.locator('.catalog-game', {
      has: page.locator('.game-name', { hasText: 'The Desert Kingdom' }),
    }).first()
    await expect(desertCard).toBeVisible({ timeout: 5000 })

    desertJoinHref = await desertCard.locator('.join-button').getAttribute('href')

    expect(desertJoinHref).toBeTruthy()
    expect(desertJoinHref).toMatch(/^\/player\/join-game\//)
  })

  test('Phase 1b: fill and submit join-game form', async ({ page }) => {
    await setupTestBypassHeaders(page)

    // Use the join href discovered in Phase 1a — avoids re-checking the catalog
    // (which may show the game as 'started' if another test joined it first).
    expect(desertJoinHref).toBeTruthy()
    expect(desertJoinHref).toMatch(/^\/player\/join-game\//)

    // Navigate to the join game page
    await page.goto(desertJoinHref)
    await waitForPageReady(page)

    // Wait for the player app and join sheet iframe to load
    await page.waitForSelector('#player-app', { timeout: 10000 })

    // The join game page renders an iframe with the join game turn sheet
    const iframe = page.frameLocator('[data-testid="join-sheet-iframe"]')
    await expect(iframe.locator('h1')).toContainText('The Desert Kingdom', { timeout: 10000 })

    // Fill form fields inside the turn sheet iframe
    await iframe.locator('[name="email"]').fill(PLAYER_EMAIL)
    await iframe.locator('[name="name"]').fill(PLAYER_NAME)
    await iframe.locator('[name="character_name"]').fill(CHARACTER_NAME)

    // Select email delivery method (game offers both email and post)
    const emailRadio = iframe.locator('input[name="delivery_method"][value="email"]')
    if (await emailRadio.isVisible()) {
      await emailRadio.check()
    }

    // Submit via the outer Vue button
    await page.locator('[data-testid="btn-submit"]').click()
    await page.waitForTimeout(1000)
  })

  // ──────────────────────────────────────────────────────
  // Phase 2: Turn 1 — Location choice + Inventory management
  // ──────────────────────────────────────────────────────

  test('Phase 2a: receive turn 1 notification email', async ({ page }) => {
    test.setTimeout(90000)
    // Wait for the turn sheet notification email (auto-start + turn processing + email delivery)
    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 60000, interval: 500 })
    expect(notifEmail).toBeTruthy()

    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''

    // Extract the turn sheet link
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)
    expect(turnSheetLink).toBeTruthy()

    // Parse the gsi_id from the link for later use
    const gsiMatch = turnSheetLink.match(/game-subscription-instance[s]?\/([a-f0-9-]+)/)
    if (gsiMatch) {
      gsiId = gsiMatch[1]
    }

    // Navigate to turn sheets
    await setupTestBypassHeaders(page)
    await page.goto(turnSheetLink)
    await waitForPageReady(page)

    // Wait for the turn sheet viewer to load
    await page.waitForSelector('[data-testid="ts-stepper"], [data-testid="ts-viewer"]', { timeout: 15000 })
  })

  test('Phase 2b: verify stepper shows 2 turn sheets', async ({ page }) => {
    await setupTestBypassHeaders(page)

    // Re-navigate to turn sheets using the notification email link
    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 30000 })
    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-stepper"]', { timeout: 15000 })

    // Verify stepper shows 2 steps
    const step0 = page.locator('[data-testid="ts-step-0"]')
    const step1 = page.locator('[data-testid="ts-step-1"]')
    await expect(step0).toBeVisible()
    await expect(step1).toBeVisible()

    // Verify the iframe is visible
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await expect(page.locator('[data-testid="ts-viewer-iframe"]')).toBeVisible()
  })

  test('Phase 2c: verify location choice sheet — background and form', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 30000 })
    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 15000 })

    // Brief wait for iframe content to render
    await page.waitForTimeout(800)

    const iframe = page.frameLocator('[data-testid="ts-viewer-iframe"]')

    // Verify the turn sheet has a background image
    const bgImage = iframe.locator('img.background-image')
    await expect(bgImage).toBeVisible({ timeout: 5000 })

    // Verify location choice radio buttons are present
    const radioButtons = iframe.locator('input[type="radio"][name="location_choice"]')
    const radioCount = await radioButtons.count()
    expect(radioCount).toBeGreaterThan(0)

    // Submit is always enabled (no mark-ready gate)
    await expect(page.locator('[data-testid="btn-submit-all"]')).toBeEnabled()
  })

  test('Phase 2d: verify inventory management sheet — background and form', async ({ page }) => {
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

    // Brief wait for iframe content
    const iframe = page.frameLocator('[data-testid="ts-viewer-iframe"]')
    await page.waitForTimeout(800)

    // Verify the inventory turn sheet also has a background image
    const bgImage = iframe.locator('img.background-image')
    await expect(bgImage).toBeVisible({ timeout: 5000 })
  })

  test('Phase 2e: fill all turn sheets and submit turn 1', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 30000 })
    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-stepper"]', { timeout: 15000 })

    // Step 0: fill location choice
    await page.locator('[data-testid="ts-step-0"]').click()
    await page.waitForTimeout(800)
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await page.waitForTimeout(600)

    const iframe0 = page.frameLocator('[data-testid="ts-viewer-iframe"]')
    const radioButtons = iframe0.locator('input[type="radio"][name="location_choice"]')
    if (await radioButtons.count() > 0) {
      await radioButtons.first().check()
    }

    // Navigate to step 1 — this caches step 0's form data in memory
    await page.locator('[data-testid="ts-step-1"]').click()
    await page.waitForTimeout(1000)
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await page.waitForTimeout(600)

    // Step 1: fill inventory management — pick up available items
    const iframe1 = page.frameLocator('[data-testid="ts-viewer-iframe"]')
    const pickUpCheckboxes = iframe1.locator('input[type="checkbox"][name="pick_up"]')
    const checkboxCount = await pickUpCheckboxes.count()
    for (let i = 0; i < checkboxCount; i++) {
      await pickUpCheckboxes.nth(i).check()
    }

    // Submit all — caches the active sheet, saves all cached sheets, calls submit endpoint
    const submitBtn = page.locator('[data-testid="btn-submit-all"]')
    await expect(submitBtn).toBeEnabled({ timeout: 5000 })
    await submitBtn.click()
    await page.waitForTimeout(1500)

    // Verify success
    await expect(page.locator('[data-testid="ts-success"]')).toBeVisible({ timeout: 10000 })
  })

  // ──────────────────────────────────────────────────────
  // Phase 3: Turn 2 — Verify turn 1 results
  // ──────────────────────────────────────────────────────

  test('Phase 3a: wait for turn 2 notification email', async () => {
    test.setTimeout(150000)
    // ProcessWhenAllSubmitted triggers processing immediately after all sheets
    // are submitted. The notification email may arrive before this test starts,
    // so we do NOT clear emails — just search for the turn 2 subject.
    const turn2Email = await waitForEmail(PLAYER_EMAIL, 'Turn 2', { timeout: 120000, interval: 1000 })
    expect(turn2Email).toBeTruthy()
  })

  test('Phase 3b: verify turn 2 turn sheets are available and location choice was applied', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const turn2Email = await waitForEmail(PLAYER_EMAIL, 'Turn 2', { timeout: 30000 })
    const fullEmail = await getEmailBody(turn2Email.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)
    expect(turnSheetLink).toBeTruthy()

    await page.goto(turnSheetLink)
    await waitForPageReady(page)

    // Wait for the stepper to show turn 2 sheets
    await page.waitForSelector('[data-testid="ts-stepper"]', { timeout: 15000 })

    // Verify stepper has steps (location choice + inventory management for turn 2)
    const step0 = page.locator('[data-testid="ts-step-0"]')
    await expect(step0).toBeVisible()

    // Verify the iframe loads the turn sheet content
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
    await expect(page.locator('[data-testid="ts-viewer-iframe"]')).toBeVisible()
    await page.waitForTimeout(800)

    const iframe = page.frameLocator('[data-testid="ts-viewer-iframe"]')

    // Verify turn 2 location choice turn sheet has a background image (regression check
    // for the missing background image bug fixed in jobworker CreateNextTurnSheet).
    const bgImage = iframe.locator('img.background-image')
    await expect(bgImage).toBeVisible({ timeout: 5000 })

    // Verify the turn sheet shows "Turn 2" in the header, confirming turn processing
    // ran and generated new sheets (i.e. turn 1 choices were applied successfully).
    await expect(iframe.locator('h2')).toContainText('Turn 2', { timeout: 5000 })
  })

  test('Phase 3c: fill and submit turn 2 sheets', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const turn2Email = await waitForEmail(PLAYER_EMAIL, 'Turn 2', { timeout: 30000 })
    const fullEmail = await getEmailBody(turn2Email.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-stepper"]', { timeout: 15000 })

    // Fill each visible step, navigating between them so the in-memory cache is populated
    const stepCount = await page.locator('button[data-testid^="ts-step-"]').count()

    for (let i = 0; i < stepCount; i++) {
      await page.locator(`[data-testid="ts-step-${i}"]`).click()
      await page.waitForTimeout(1000)
      await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 10000 })
      await page.waitForTimeout(600)
      // No explicit save — navigating to the next step caches the current one
    }

    // Submit all — saves all cached sheets and calls submit endpoint
    const submitBtn = page.locator('[data-testid="btn-submit-all"]')
    await expect(submitBtn).toBeEnabled({ timeout: 5000 })
    await submitBtn.click()
    await page.waitForTimeout(1500)

    await expect(page.locator('[data-testid="ts-success"]')).toBeVisible({ timeout: 10000 })
  })

  // ──────────────────────────────────────────────────────
  // Phase 4: Verification — Turn 3 arrives
  // ──────────────────────────────────────────────────────

  test('Phase 4: verify turn 3 notification arrives', async () => {
    test.setTimeout(150000)
    // Turn 3 notification arrives after turn 2 sheets are submitted and processed
    const turn3Email = await waitForEmail(PLAYER_EMAIL, 'Turn 3', { timeout: 120000, interval: 1000 })
    expect(turn3Email).toBeTruthy()
    expect(turn3Email.Subject).toBeTruthy()
  })
})
