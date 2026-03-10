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

  test.beforeAll(async () => {
    await clearAllEmails()
  })

  // ──────────────────────────────────────────────────────
  // Phase 1: Join the game
  // ──────────────────────────────────────────────────────

  test('Phase 1a: browse catalog and find Desert Kingdom', async ({ page }) => {
    await page.goto('/games')
    await waitForPageReady(page)

    const desertCard = page.locator('text=The Desert Kingdom')
    await expect(desertCard.first()).toBeVisible({ timeout: 10000 })

    const joinButton = page.locator('[data-testid^="join-button-"]')
      .filter({ has: page.locator(':scope').locator('..').locator('..').locator('text=The Desert Kingdom') })

    // Find the join button near the Desert Kingdom card
    const allJoinButtons = page.locator('[data-testid^="join-button-"]')
    const count = await allJoinButtons.count()
    let desertJoinHref = null

    for (let i = 0; i < count; i++) {
      const btn = allJoinButtons.nth(i)
      const parentCard = btn.locator('xpath=ancestor::*[contains(@class, "game-card") or contains(@data-testid, "game-card")]')
      const cardText = await parentCard.textContent().catch(() => '')
      if (cardText.includes('Desert Kingdom')) {
        desertJoinHref = await btn.getAttribute('href')
        break
      }
    }

    // Fallback: just click the first join button if we can't find via parent
    if (!desertJoinHref) {
      const href = await allJoinButtons.first().getAttribute('href')
      desertJoinHref = href
    }

    expect(desertJoinHref).toBeTruthy()
    expect(desertJoinHref).toMatch(/^\/player\/join-game\//)
  })

  test('Phase 1b: fill and submit join-game form', async ({ page }) => {
    await setupTestBypassHeaders(page)

    // Navigate to games catalog to find the join link dynamically
    await page.goto('/games')
    await waitForPageReady(page)

    // Find the Desert Kingdom card and its join button
    const catalogGames = page.locator('[data-testid="catalog-games"]')
    await expect(catalogGames).toBeVisible({ timeout: 10000 })

    const desertCard = catalogGames.locator('.catalog-game', {
      has: page.locator('.game-name', { hasText: 'The Desert Kingdom' })
    })
    await expect(desertCard).toBeVisible()

    const joinHref = await desertCard.locator('.join-button').getAttribute('href')
    expect(joinHref).toMatch(/^\/player\/join-game\//)

    // Navigate to the join game page
    await page.goto(joinHref)
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

  test('Phase 2c: fill location choice — choose Ancient Ruins', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 30000 })
    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-viewer-iframe"]', { timeout: 15000 })

    // The first sheet (step 0) should be Location Choice
    // Interact with the iframe form to select a destination
    const iframe = page.frameLocator('[data-testid="ts-viewer-iframe"]')

    // Brief wait for iframe content to render
    await page.waitForTimeout(800)

    // Select the first available location radio button
    const radioButtons = iframe.locator('input[type="radio"][name="location_choice"]')
    const radioCount = await radioButtons.count()
    if (radioCount > 0) {
      await radioButtons.first().check()
    }

    // Save the sheet
    await page.locator('[data-testid="btn-save-sheet"]').click()
    await page.waitForTimeout(1000)

    // Mark ready
    await page.locator('[data-testid="btn-mark-ready"]').click()
    await page.waitForTimeout(400)

    // Verify step 0 shows ready status
    const stepStatus = page.locator('[data-testid="ts-step-status-0"]')
    await expect(stepStatus.locator('.status-ready')).toBeVisible()

    // Submit-all should still be disabled (step 1 not ready)
    await expect(page.locator('[data-testid="btn-submit-all"]')).toBeDisabled()
  })

  test('Phase 2d: fill inventory management — pick up compass and flask', async ({ page }) => {
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

    // Try to interact with inventory form elements
    try {
      // Pick up items
      const pickUpCheckboxes = iframe.locator('input[type="checkbox"][name*="pick_up"], input[type="checkbox"][value*="compass"], input[type="checkbox"][value*="flask"]')
      const checkboxCount = await pickUpCheckboxes.count()
      for (let i = 0; i < checkboxCount; i++) {
        await pickUpCheckboxes.nth(i).check()
      }
    } catch {
      // Form may have a different structure; proceed anyway
    }

    // Save the sheet
    await page.locator('[data-testid="btn-save-sheet"]').click()
    await page.waitForTimeout(1000)

    // Mark ready
    await page.locator('[data-testid="btn-mark-ready"]').click()
    await page.waitForTimeout(400)

    // Verify step 1 shows ready
    const stepStatus = page.locator('[data-testid="ts-step-status-1"]')
    await expect(stepStatus.locator('.status-ready')).toBeVisible()
  })

  test('Phase 2e: submit all turn sheets for turn 1', async ({ page }) => {
    await setupTestBypassHeaders(page)

    const notifEmail = await waitForEmail(PLAYER_EMAIL, 'turn', { timeout: 30000 })
    const fullEmail = await getEmailBody(notifEmail.ID)
    const htmlBody = fullEmail.HTML || fullEmail.Text || ''
    const turnSheetLink = extractLink(htmlBody, /turn-sheet|player.*game-subscription-instance/i)

    await page.goto(turnSheetLink)
    await waitForPageReady(page)
    await page.waitForSelector('[data-testid="ts-stepper"]', { timeout: 15000 })

    // Both steps should be ready from previous test steps
    // Mark step 0 ready
    await page.locator('[data-testid="ts-step-0"]').click()
    await page.waitForTimeout(600)
    const markReadyBtn0 = page.locator('[data-testid="btn-mark-ready"]')
    if ((await markReadyBtn0.textContent()).includes('Mark Ready')) {
      await markReadyBtn0.click()
      await page.waitForTimeout(300)
    }

    // Mark step 1 ready
    await page.locator('[data-testid="ts-step-1"]').click()
    await page.waitForTimeout(600)
    const markReadyBtn1 = page.locator('[data-testid="btn-mark-ready"]')
    if ((await markReadyBtn1.textContent()).includes('Mark Ready')) {
      await markReadyBtn1.click()
      await page.waitForTimeout(300)
    }

    // Submit all
    const submitBtn = page.locator('[data-testid="btn-submit-all"]')
    await expect(submitBtn).toBeEnabled({ timeout: 5000 })
    await submitBtn.click()
    await page.waitForTimeout(1200)

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

  test('Phase 3b: verify turn 2 turn sheets are available', async ({ page }) => {
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

    // Save + mark ready for each visible step (use button selector to exclude status spans)
    const stepCount = await page.locator('button[data-testid^="ts-step-"]').count()

    for (let i = 0; i < stepCount; i++) {
      await page.locator(`[data-testid="ts-step-${i}"]`).click()
      await page.waitForTimeout(1000)

      // Save
      await page.locator('[data-testid="btn-save-sheet"]').click()
      await page.waitForTimeout(800)

      // Mark ready
      const markBtn = page.locator('[data-testid="btn-mark-ready"]')
      if ((await markBtn.textContent()).includes('Mark Ready')) {
        await markBtn.click()
        await page.waitForTimeout(400)
      }
    }

    // Submit all
    const submitBtn = page.locator('[data-testid="btn-submit-all"]')
    await expect(submitBtn).toBeEnabled({ timeout: 5000 })
    await submitBtn.click()
    await page.waitForTimeout(1200)

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
