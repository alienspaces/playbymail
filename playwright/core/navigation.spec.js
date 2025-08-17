import { test, expect } from '@playwright/test'
import { 
  navigateTo, 
  waitForPageReady, 
  checkPageTitle, 
  checkPageURL,
  checkElementVisible,
  takeScreenshot
} from '../utils/test-helpers.js'

test.describe('Core Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Clear any existing state before each test
    await page.context().clearCookies()
    await page.context().clearStorageState()
  })

  test.describe('Public Pages', () => {
    test('should load home page successfully', async ({ page }) => {
      await navigateTo(page, '/')
      
      // Check page title
      await checkPageTitle(page, 'Play by Mail')
      
      // Check page URL
      await checkPageURL(page, '/')
      
      // Verify main content is visible
      await checkElementVisible(page, '#app')
      
      // Take screenshot for visual verification
      await takeScreenshot(page, 'home-page')
    })

    test('should load FAQ page', async ({ page }) => {
      await navigateTo(page, '/faq')
      
      await checkPageTitle(page, 'FAQ')
      await checkPageURL(page, '/faq')
      await checkElementVisible(page, '#app')
    })

    test('should load about page', async ({ page }) => {
      await navigateTo(page, '/about')
      
      await checkPageTitle(page, 'About')
      await checkPageURL(page, '/about')
      await checkElementVisible(page, '#app')
    })
  })

  test.describe('Authentication Pages', () => {
    test('should load login page', async ({ page }) => {
      await navigateTo(page, '/login')
      
      await checkPageTitle(page, 'Sign In')
      await checkPageURL(page, '/login')
      
      // Check for login form elements
      await checkElementVisible(page, 'input[type="email"]')
      await checkElementVisible(page, 'button:has-text("Send Code")')
    })

    test('should load verification page with email parameter', async ({ page }) => {
      await navigateTo(page, '/verify?email=test@example.com')
      
      await checkPageTitle(page, 'Verify')
      await checkPageURL(page, /\/verify/)
      
      // Check for verification form elements
      await checkElementVisible(page, 'input[type="text"]')
      await checkElementVisible(page, 'button:has-text("Verify")')
    })
  })

  test.describe('Protected Pages (Unauthenticated)', () => {
    test('should load studio page with unauthenticated content', async ({ page }) => {
      await navigateTo(page, '/studio')
      
      await checkPageTitle(page, 'Studio')
      await checkPageURL(page, '/studio')
      
      // Should show unauthenticated studio content
      await checkElementVisible(page, 'body')
      
      // Take screenshot to verify unauthenticated state
      await takeScreenshot(page, 'studio-unauthenticated')
    })

    test('should load admin page with unauthenticated content', async ({ page }) => {
      await navigateTo(page, '/admin')
      
      await checkPageTitle(page, 'Admin')
      await checkPageURL(page, '/admin')
      
      // Should show unauthenticated admin content
      await checkElementVisible(page, 'body')
      
      await takeScreenshot(page, 'admin-unauthenticated')
    })

    test('should load account page with error for unauthenticated users', async ({ page }) => {
      await navigateTo(page, '/account')
      
      await checkPageTitle(page, 'Account')
      await checkPageURL(page, '/account')
      
      // Should show error or unauthenticated content
      await checkElementVisible(page, 'body')
      
      await takeScreenshot(page, 'account-unauthenticated')
    })
  })

  test.describe('Navigation Behavior', () => {
    test('should handle direct URL navigation', async ({ page }) => {
      // Test direct navigation to various pages
      const testPages = ['/', '/faq', '/about', '/login', '/studio', '/admin']
      
      for (const path of testPages) {
        await navigateTo(page, path)
        await checkPageURL(page, path)
        
        // Verify page loads without errors
        await checkElementVisible(page, 'body')
        
        // Small delay between navigations
        await page.waitForTimeout(500)
      }
    })

    test('should handle browser back/forward navigation', async ({ page }) => {
      // Navigate to multiple pages
      await navigateTo(page, '/')
      await navigateTo(page, '/faq')
      await navigateTo(page, '/about')
      
      // Go back
      await page.goBack()
      await checkPageURL(page, '/faq')
      
      // Go back again
      await page.goBack()
      await checkPageURL(page, '/')
      
      // Go forward
      await page.goForward()
      await checkPageURL(page, '/faq')
    })
  })

  test.describe('Page Load Performance', () => {
    test('should load pages within reasonable time', async ({ page }) => {
      const startTime = Date.now()
      
      await navigateTo(page, '/')
      
      const loadTime = Date.now() - startTime
      
      // Pages should load within 5 seconds
      expect(loadTime).toBeLessThan(5000)
      
      console.log(`Home page loaded in ${loadTime}ms`)
    })

    test('should handle slow network gracefully', async ({ page }) => {
      // Simulate slow network
      await page.route('**/*', route => {
        // Add small delay to all requests
        setTimeout(() => route.continue(), 100)
      })
      
      const startTime = Date.now()
      await navigateTo(page, '/')
      const loadTime = Date.now() - startTime
      
      // Should still load within reasonable time even with delays
      expect(loadTime).toBeLessThan(10000)
      
      console.log(`Page loaded with network delays in ${loadTime}ms`)
    })
  })
})
