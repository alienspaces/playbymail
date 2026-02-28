import { test, expect } from '@playwright/test'
import {
  TEST_BYPASS_HEADER_NAME,
  TEST_BYPASS_HEADER_VALUE,
  navigateTo,
  waitForPageReady,
  fillFormField,
  safeClick,
  isAuthenticated,
} from '../utils/test-helpers.js'

const TEST_EMAIL = 'playwright-auth-test@example.com'

test.describe('Authentication Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
    await page.evaluate(() => localStorage.clear())
  })

  test.describe('Login Page', () => {
    test('displays login form with email input and submit button', async ({ page }) => {
      await navigateTo(page, '/login')
      await expect(page).toHaveURL('/login')
      await expect(page.locator('[data-testid="email-input"]')).toBeVisible()
      await expect(page.locator('[data-testid="login-submit"]')).toBeVisible()
    })

    test('stays on login page when submitting empty email', async ({ page }) => {
      await navigateTo(page, '/login')
      await safeClick(page, '[data-testid="login-submit"]')
      await expect(page).toHaveURL('/login')
    })

    test('submits email and redirects to verification page', async ({ page }) => {
      await page.context().setExtraHTTPHeaders({
        [TEST_BYPASS_HEADER_NAME]: TEST_BYPASS_HEADER_VALUE,
      })
      await navigateTo(page, '/login')
      await fillFormField(page, '[data-testid="email-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="login-submit"]')
      await expect(page).toHaveURL(/\/verify/)
    })
  })

  test.describe('Verification Page', () => {
    test('displays verification form', async ({ page }) => {
      await navigateTo(page, `/verify?email=${encodeURIComponent(TEST_EMAIL)}`)
      await expect(page.locator('[data-testid="verify-code-input"]')).toBeVisible()
      await expect(page.locator('[data-testid="verify-submit"]')).toBeVisible()
    })

    test('shows error for invalid verification code', async ({ page }) => {
      await page.context().setExtraHTTPHeaders({
        [TEST_BYPASS_HEADER_NAME]: TEST_BYPASS_HEADER_VALUE,
      })
      await navigateTo(page, `/verify?email=${encodeURIComponent(TEST_EMAIL)}`)
      await fillFormField(page, '[data-testid="verify-code-input"]', 'WRONGCODE')
      await safeClick(page, '[data-testid="verify-submit"]')
      await page.waitForTimeout(2000)
      await expect(page).toHaveURL(/\/verify/)
    })

    test('redirects to login when email is missing', async ({ page }) => {
      await navigateTo(page, '/verify')
      await expect(page).toHaveURL('/login')
    })
  })

  test.describe('Full Authentication Cycle', () => {
    test('login -> verify with bypass -> authenticated session', async ({ page }) => {
      await page.context().setExtraHTTPHeaders({
        [TEST_BYPASS_HEADER_NAME]: TEST_BYPASS_HEADER_VALUE,
      })

      // Step 1: Navigate to login
      await navigateTo(page, '/login')
      await fillFormField(page, '[data-testid="email-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="login-submit"]')

      // Step 2: Verify with email as code (bypass mode)
      await expect(page).toHaveURL(/\/verify/)
      await fillFormField(page, '[data-testid="verify-code-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="verify-submit"]')

      // Step 3: Should redirect to home and be authenticated
      await page.waitForURL('/', { timeout: 10000 })
      await waitForPageReady(page)

      const authenticated = await isAuthenticated(page)
      expect(authenticated).toBe(true)
    })

    test('authenticated state persists across page reload', async ({ page }) => {
      await page.context().setExtraHTTPHeaders({
        [TEST_BYPASS_HEADER_NAME]: TEST_BYPASS_HEADER_VALUE,
      })

      // Authenticate
      await navigateTo(page, '/login')
      await fillFormField(page, '[data-testid="email-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="login-submit"]')
      await expect(page).toHaveURL(/\/verify/)
      await fillFormField(page, '[data-testid="verify-code-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="verify-submit"]')
      await page.waitForURL('/', { timeout: 10000 })

      // Reload page
      await page.reload()
      await waitForPageReady(page)

      const authenticated = await isAuthenticated(page)
      expect(authenticated).toBe(true)
    })

    test('sign out clears authenticated state', async ({ page }) => {
      await page.context().setExtraHTTPHeaders({
        [TEST_BYPASS_HEADER_NAME]: TEST_BYPASS_HEADER_VALUE,
      })

      // Authenticate
      await navigateTo(page, '/login')
      await fillFormField(page, '[data-testid="email-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="login-submit"]')
      await expect(page).toHaveURL(/\/verify/)
      await fillFormField(page, '[data-testid="verify-code-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="verify-submit"]')
      await page.waitForURL('/', { timeout: 10000 })
      await waitForPageReady(page)

      // Sign out
      const signOutLink = page.locator('[data-testid="sign-out-link"]')
      if (await signOutLink.isVisible({ timeout: 3000 }).catch(() => false)) {
        await signOutLink.click()
        await waitForPageReady(page)
        const authenticated = await isAuthenticated(page)
        expect(authenticated).toBe(false)
      }
    })
  })

  test.describe('Error Handling', () => {
    test('handles network errors on login gracefully', async ({ page }) => {
      await navigateTo(page, '/login')
      await page.route('**/api/v1/request-auth', (route) => route.abort('failed'))
      await fillFormField(page, '[data-testid="email-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="login-submit"]')
      await page.waitForTimeout(2000)
      await expect(page.locator('[data-testid="login-error"]')).toBeVisible()
    })

    test('handles server errors on login gracefully', async ({ page }) => {
      await navigateTo(page, '/login')
      await page.route('**/api/v1/request-auth', (route) =>
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal Server Error' }),
        })
      )
      await fillFormField(page, '[data-testid="email-input"]', TEST_EMAIL)
      await safeClick(page, '[data-testid="login-submit"]')
      await page.waitForTimeout(2000)
      await expect(page.locator('[data-testid="login-error"]')).toBeVisible()
    })
  })
})
