import { test, expect } from '@playwright/test'

test.describe('Home Page', () => {
  test('should load home page for unauthenticated users', async ({ page }) => {
    // Navigate to home page
    await page.goto('/')
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle')
    
    // Verify the page loaded successfully - fix the title pattern
    await expect(page).toHaveTitle(/Play by Mail/i)
    
    // Take a screenshot for debugging
    await page.screenshot({ path: 'playwright/home-page.png' })
  })

  test('should show unauthenticated user content', async ({ page }) => {
    await page.goto('/')
    
    // Wait for any dynamic content to load
    await page.waitForLoadState('networkidle')
    
    // Verify the page has some content - use a more specific selector
    const mainContent = page.locator('#app').first()
    await expect(mainContent).toBeVisible()
    
    // Also check for some specific content to ensure the page loaded properly
    const pageContent = page.locator('body')
    await expect(pageContent).toContainText(/Play by Mail|Welcome|Game/i)
    
    // Check that we're not redirected to login (meaning we can access public pages)
    await expect(page).not.toHaveURL(/\/login/)
  })

  test('should show unauthenticated content on studio page', async ({ page }) => {
    // Studio page should be accessible but show unauthenticated content
    await page.goto('/studio')
    
    // Wait for page to load
    await page.waitForLoadState('networkidle')
    
    // Should show StudioEntryView for unauthenticated users
    const pageContent = page.locator('body')
    
    // The page should load without redirecting to login
    await expect(page).toHaveURL('/studio')
    
    // Should show some studio-related content (even if limited)
    await expect(pageContent).toContainText(/Studio|Game|Design/i)
  })
})
