// Common test utilities for Playwright tests
import { expect } from '@playwright/test'

// Test bypass authentication configuration
// These must be set via environment variables - no defaults allowed
export const TEST_BYPASS_HEADER_NAME = process.env.TEST_BYPASS_HEADER_NAME
export const TEST_BYPASS_HEADER_VALUE = process.env.TEST_BYPASS_HEADER_VALUE

// Validate required environment variables at module load time
if (!TEST_BYPASS_HEADER_NAME) {
  throw new Error('TEST_BYPASS_HEADER_NAME environment variable is required')
}
if (!TEST_BYPASS_HEADER_VALUE) {
  throw new Error('TEST_BYPASS_HEADER_VALUE environment variable is required')
}

// Set up test bypass headers on a page context for API requests
export async function setupTestBypassHeaders(page) {
  await page.setExtraHTTPHeaders({
    [TEST_BYPASS_HEADER_NAME]: TEST_BYPASS_HEADER_VALUE
  })
}

export async function waitForElement(page, selector, timeout = 5000) {
  await page.waitForSelector(selector, { timeout })
  // Return locator that handles multiple matches
  return page.locator(selector).first()
}

export async function waitForPageReady(page) {
  // Wait for page to be fully loaded (avoid networkidle as it hangs with polling/websockets)
  await page.waitForLoadState('domcontentloaded')
  await page.waitForLoadState('load')
}

export async function safeFill(page, selector, value) {
  const element = await waitForElement(page, selector)
  await element.clear()
  await element.fill(value)
}

export async function safeClick(page, selector) {
  const element = await waitForElement(page, selector)
  await element.scrollIntoViewIfNeeded()
  await element.click()
}

export async function isAuthenticated(page) {
  // Check if user is authenticated by looking for authenticated content
  const body = page.locator('body')
  const content = await body.textContent()

  // Look for signs of authentication
  const authenticatedIndicators = [
    'Games & Instances',
    'Manage Instances',
    'Account Settings',
    'Sign Out'
  ]

  return authenticatedIndicators.some(indicator =>
    content.includes(indicator)
  )
}

export async function navigateTo(page, path) {
  await page.goto(path)
  await waitForPageReady(page)
}

export async function takeScreenshot(page, name) {
  await page.screenshot({
    path: `playwright/screenshots/${name}-${Date.now()}.png`,
    fullPage: true
  })
}

export async function waitForText(page, text, timeout = 5000) {
  await page.waitForFunction(
    (text) => document.body.textContent.includes(text),
    text,
    { timeout }
  )
}

export async function waitForElementToDisappear(page, selector, timeout = 5000) {
  await page.waitForSelector(selector, { state: 'hidden', timeout })
}

export async function getElementText(page, selector) {
  const element = await waitForElement(page, selector)
  return element.textContent()
}

export async function checkPageTitle(page, expectedTitle) {
  await expect(page).toHaveTitle(new RegExp(expectedTitle, 'i'))
}

export async function checkPageURL(page, expectedURL) {
  await expect(page).toHaveURL(expectedURL)
}

export async function checkElementVisible(page, selector) {
  const element = await waitForElement(page, selector)
  // Use .first() to handle multiple matching elements
  await expect(element.first()).toBeVisible()
}

export async function checkElementContainsText(page, selector, text) {
  const element = await waitForElement(page, selector)
  await expect(element).toContainText(text)
}

export async function checkButtonEnabled(page, selector) {
  const button = await waitForElement(page, selector)
  await expect(button).toBeEnabled()
}

export async function checkButtonDisabled(page, selector) {
  const button = await waitForElement(page, selector)
  await expect(button).toBeDisabled()
}

export async function fillFormField(page, fieldSelector, value) {
  await safeFill(page, fieldSelector, value)
}

export async function submitForm(page, formSelector) {
  const form = await waitForElement(page, formSelector)
  await form.submit()
}

export async function waitForLoadingState(page, loadingSelector, timeout = 10000) {
  // Wait for loading to appear
  await waitForElement(page, loadingSelector, 2000)

  // Wait for loading to disappear
  await waitForElementToDisappear(page, loadingSelector, timeout)
}

export async function checkErrorDisplayed(page, errorSelector, expectedError) {
  const errorElement = await waitForElement(page, errorSelector)
  await expect(errorElement).toBeVisible()

  if (expectedError) {
    await expect(errorElement).toContainText(expectedError)
  }
}

export async function checkSuccessMessage(page, successSelector, expectedMessage) {
  const successElement = await waitForElement(page, successSelector)
  await expect(successElement).toBeVisible()

  if (expectedMessage) {
    await expect(successElement).toContainText(expectedMessage)
  }
}
