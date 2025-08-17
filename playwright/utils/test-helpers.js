/**
 * Test helper utilities for Playwright tests
 */

/**
 * Wait for an element to be visible and stable
 */
export async function waitForElement(page, selector, timeout = 5000) {
  const element = page.locator(selector)
  await element.waitFor({ state: 'visible', timeout })
  return element
}

/**
 * Wait for page to be fully loaded and stable
 */
export async function waitForPageReady(page) {
  await page.waitForLoadState('networkidle')
  await page.waitForTimeout(500) // Small delay for any final animations
}

/**
 * Fill a form field safely
 */
export async function safeFill(page, selector, value) {
  const field = page.locator(selector)
  if (await field.isVisible()) {
    await field.fill(value)
    return true
  }
  return false
}

/**
 * Click a button safely
 */
export async function safeClick(page, selector) {
  const button = page.locator(selector)
  if (await button.isVisible()) {
    await button.click()
    return true
  }
  return false
}

/**
 * Check if user is authenticated by looking for auth indicators
 */
export async function isAuthenticated(page) {
  // Look for common authenticated user indicators
  const authIndicators = [
    '.user-menu',
    '.profile',
    '[data-testid="user-menu"]',
    '.logout',
    'button:has-text("Logout")'
  ]
  
  for (const selector of authIndicators) {
    if (await page.locator(selector).isVisible()) {
      return true
    }
  }
  
  return false
}

/**
 * Navigate to a page and wait for it to be ready
 */
export async function navigateTo(page, path) {
  await page.goto(path)
  await waitForPageReady(page)
}
