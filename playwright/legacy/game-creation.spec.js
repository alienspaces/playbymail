import { test, expect } from '@playwright/test'

test.describe('Game Creation Flow', () => {
  test('should show unauthenticated content on studio page', async ({ page }) => {
    // Studio page should be accessible but show unauthenticated content
    await page.goto('/studio')
    
    // Wait for page to load
    await page.waitForLoadState('networkidle')
    
    // Should show StudioEntryView for unauthenticated users
    // Look for content that indicates this is the unauthenticated studio view
    const pageContent = page.locator('body')
    
    // The page should load without redirecting to login
    await expect(page).toHaveURL('/studio')
    
    // Should show some studio-related content (even if limited)
    await expect(pageContent).toContainText(/Studio|Game|Design/i)
  })

  test('should show login form for unauthenticated users', async ({ page }) => {
    await page.goto('/login')
    
    await page.waitForLoadState('networkidle')
    
    // Check for login form elements
    const emailInput = page.locator('input[type="email"], input[name="email"]')
    const submitButton = page.locator('button[type="submit"], button:has-text("Send Code")')
    
    await expect(emailInput).toBeVisible()
    await expect(submitButton).toBeVisible()
    
    // Check for some login-specific content
    const pageContent = page.locator('body')
    await expect(pageContent).toContainText(/Sign in|Email|Send Code/i)
  })

  test('should handle login form submission', async ({ page }) => {
    await page.goto('/login')
    
    // Fill in email
    const emailInput = page.locator('input[type="email"], input[name="email"]')
    await emailInput.fill('test@example.com')
    
    // Submit form
    const submitButton = page.locator('button[type="submit"], button:has-text("Send Code")')
    await submitButton.click()
    
    // Should redirect to verification page
    await expect(page).toHaveURL(/\/verify/)
    
    // Verify page shows email
    const emailDisplay = page.locator('.email, [data-testid="email"]')
    if (await emailDisplay.isVisible()) {
      await expect(emailDisplay).toContainText('test@example.com')
    }
  })
})
