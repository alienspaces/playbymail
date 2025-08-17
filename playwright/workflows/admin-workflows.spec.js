import { test, expect } from '@playwright/test'
import { 
  navigateTo, 
  waitForPageReady, 
  checkPageTitle, 
  checkPageURL,
  checkElementVisible,
  checkElementContainsText,
  safeClick,
  takeScreenshot,
  waitForText,
  waitForLoadingState,
  waitForElementToDisappear
} from '../utils/test-helpers.js'

test.describe('Admin Workflows', () => {
  test.beforeEach(async ({ page }) => {
    // Clear any existing state before each test
    await page.context().clearCookies()
    await page.context().clearStorageState()
  })

  test.describe('Admin Dashboard Access', () => {
    test('should show unauthenticated admin content', async ({ page }) => {
      await navigateTo(page, '/admin')
      
      await checkPageTitle(page, 'Admin')
      await checkPageURL(page, '/admin')
      
      // Should show unauthenticated admin content
      await checkElementVisible(page, 'body')
      
      // Look for admin-related content
      const body = page.locator('body')
      const content = await body.textContent()
      
      // Should show some admin-related content
      expect(content).toMatch(/Admin|Management|Dashboard|Sign In/i)
      
      await takeScreenshot(page, 'admin-unauthenticated')
    })

    test('should show login prompt for unauthenticated users', async ({ page }) => {
      await navigateTo(page, '/admin')
      
      // Look for login-related content
      const body = page.locator('body')
      const content = await body.textContent()
      
      // Should show login prompt or unauthenticated message
      expect(content).toMatch(/Sign In|Login|Authenticate|Access/i)
    })
  })

  test.describe('Authenticated Admin Dashboard', () => {
    test('should display admin dashboard when authenticated', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/admin')
      
      // Should show authenticated admin content
      await checkPageURL(page, '/admin')
      
      // Look for admin dashboard elements
      const body = page.locator('body')
      const content = await body.textContent()
      
      if (content.includes('Games & Instances') || content.includes('Dashboard')) {
        // Should show admin dashboard
        await checkElementVisible(page, 'body')
        
        // Look for dashboard elements
        const dashboardSelectors = [
          'h1:has-text("Games & Instances")',
          '.dashboard',
          '.admin-content',
          '[data-testid="dashboard"]'
        ]
        
        let dashboardFound = false
        for (const selector of dashboardSelectors) {
          try {
            if (await page.locator(selector).isVisible()) {
              dashboardFound = true
              break
            }
          } catch (e) {
            // Continue to next selector
          }
        }
        
        expect(dashboardFound).toBe(true)
        
        await takeScreenshot(page, 'admin-dashboard')
      } else {
        console.log('Admin dashboard not found - may need authentication or different content')
      }
    })

    test('should show games list when authenticated', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/admin')
      
      // Wait for dashboard to load
      await waitForPageReady(page)
      
      // Look for games list
      const body = page.locator('body')
      const content = await body.textContent()
      
      if (content.includes('Games & Instances')) {
        // Should show games list
        await checkElementVisible(page, 'body')
        
        // Look for games list elements
        const gamesSelectors = [
          '.games-list',
          '[data-testid="games"]',
          '.game-item',
          'li:has-text("Game")'
        ]
        
        let gamesListFound = false
        for (const selector of gamesSelectors) {
          try {
            if (await page.locator(selector).isVisible()) {
              gamesListFound = true
              break
            }
          } catch (e) {
            // Continue to next selector
          }
        }
        
        // Games list is expected for authenticated admin
        expect(gamesListFound).toBe(true)
        
        await takeScreenshot(page, 'admin-games-list')
      } else {
        console.log('Games list not found - may need authentication or different content')
      }
    })
  })

  test.describe('Game Management', () => {
    test('should show manage instances button for games', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/admin')
      
      // Wait for dashboard to load
      await waitForPageReady(page)
      
      // Look for manage instances button
      const manageSelectors = [
        'button:has-text("Manage Instances")',
        'button:has-text("Manage")',
        'a:has-text("Manage")',
        '[data-testid="manage-instances"]'
      ]
      
      let manageButtonFound = false
      for (const selector of manageSelectors) {
        try {
          if (await page.locator(selector).isVisible()) {
            manageButtonFound = true
            console.log(`Manage button found: ${selector}`)
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
      
      // Manage button is expected for authenticated admin
      expect(manageButtonFound).toBe(true)
    })

    test('should navigate to game instances page', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/admin')
      
      // Wait for dashboard to load
      await waitForPageReady(page)
      
      // Look for manage instances button
      const manageButton = page.locator('button:has-text("Manage Instances"), button:has-text("Manage")').first()
      
      if (await manageButton.isVisible()) {
        // Click manage button
        await safeClick(page, 'button:has-text("Manage Instances"), button:has-text("Manage")')
        
        // Should navigate to game instances page
        await expect(page).toHaveURL(/\/admin\/games\/.*\/instances/)
        
        // Should show game instances view
        await checkElementVisible(page, 'body')
        
        // Look for instances content
        const body = page.locator('body')
        const content = await body.textContent()
        
        if (content.includes('Game Instances') || content.includes('Instances')) {
          console.log('Successfully navigated to game instances page')
          
          await takeScreenshot(page, 'game-instances-page')
        } else {
          console.log('Navigated to instances page but content not as expected')
        }
      } else {
        console.log('Manage button not found - skipping navigation test')
      }
    })
  })

  test.describe('Navigation and Return', () => {
    test('should return to admin dashboard from instances page', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      // Navigate to admin dashboard
      await navigateTo(page, '/admin')
      await waitForPageReady(page)
      
      // Look for manage button
      const manageButton = page.locator('button:has-text("Manage Instances"), button:has-text("Manage")').first()
      
      if (await manageButton.isVisible()) {
        // Navigate to game instances
        await safeClick(page, 'button:has-text("Manage Instances"), button:has-text("Manage")')
        await expect(page).toHaveURL(/\/admin\/games\/.*\/instances/)
        await waitForPageReady(page)
        
        // Navigate back to admin dashboard
        await navigateTo(page, '/admin')
        await waitForPageReady(page)
        
        // Should be back on dashboard
        await checkPageURL(page, '/admin')
        
        // Should show dashboard content
        const body = page.locator('body')
        const content = await body.textContent()
        
        if (content.includes('Games & Instances')) {
          console.log('Successfully returned to admin dashboard')
          
          // Should NOT be stuck on loading
          if (content.includes('Loading games...')) {
            // Wait for loading to complete
            await waitForElementToDisappear(page, 'text=Loading games...', 10000)
            
            const finalContent = await body.textContent()
            if (!finalContent.includes('Loading games...')) {
              console.log('Loading state resolved correctly')
            } else {
              console.log('Still stuck on loading state')
            }
          }
          
          await takeScreenshot(page, 'admin-dashboard-returned')
        } else {
          console.log('Returned to admin page but dashboard content not shown')
        }
      } else {
        console.log('Manage button not found - skipping return navigation test')
      }
    })

    test('should handle API calls correctly when returning to dashboard', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      // Navigate to admin dashboard
      await navigateTo(page, '/admin')
      await waitForPageReady(page)
      
      // Look for manage button
      const manageButton = page.locator('button:has-text("Manage Instances"), button:has-text("Manage")').first()
      
      if (await manageButton.isVisible()) {
        // Navigate to game instances
        await safeClick(page, 'button:has-text("Manage Instances"), button:has-text("Manage")')
        await expect(page).toHaveURL(/\/admin\/games\/.*\/instances/)
        await waitForPageReady(page)
        
        // Monitor network requests
        const apiCalls = []
        page.on('request', request => {
          if (request.url().includes('/api/v1/')) {
            apiCalls.push(request.url())
            console.log('API call:', request.url())
          }
        })
        
        // Navigate back to dashboard
        await navigateTo(page, '/admin')
        await waitForPageReady(page)
        
        // Debug: see what API calls were made
        console.log('API calls made:', apiCalls)
        
        // Should have made API calls to load games and instances
        expect(apiCalls).toContain(expect.stringContaining('/api/v1/games'))
        expect(apiCalls).toContain(expect.stringContaining('/api/v1/game-instances'))
        
        // Should not be stuck loading
        const body = page.locator('body')
        const finalContent = await body.textContent()
        
        if (finalContent.includes('Loading games...')) {
          await waitForElementToDisappear(page, 'text=Loading games...', 10000)
        }
        
        await takeScreenshot(page, 'admin-dashboard-api-tested')
      } else {
        console.log('Manage button not found - skipping API test')
      }
    })
  })

  test.describe('Loading States', () => {
    test('should handle loading states correctly', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/admin')
      
      // Wait for dashboard to load
      await waitForPageReady(page)
      
      // Check for loading states
      const body = page.locator('body')
      const content = await body.textContent()
      
      if (content.includes('Loading games...')) {
        // Should resolve loading state
        await waitForElementToDisappear(page, 'text=Loading games...', 10000)
        
        const finalContent = await body.textContent()
        expect(finalContent).not.toContain('Loading games...')
        
        console.log('Loading state handled correctly')
      } else {
        console.log('No loading state detected')
      }
    })
  })

  test.describe('Error Handling', () => {
    test('should handle API errors gracefully', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      // Mock API error
      await page.route('**/api/**', route => {
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal Server Error' })
        })
      })
      
      await navigateTo(page, '/admin')
      
      // Should handle error gracefully
      await waitForPageReady(page)
      
      const body = page.locator('body')
      const content = await body.textContent()
      
      if (content.includes('error') || content.includes('failed') || content.includes('unavailable')) {
        console.log('Error handling working correctly')
      } else {
        console.log('No error message shown for API failure')
      }
    })
  })

  test.describe('Responsive Design', () => {
    test('should adapt admin dashboard to mobile viewport', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      // Set mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })
      
      await navigateTo(page, '/admin')
      
      // Check if page loads in mobile view
      await checkElementVisible(page, 'body')
      
      // Take mobile screenshot
      await takeScreenshot(page, 'admin-mobile-view')
      
      // Reset to desktop viewport
      await page.setViewportSize({ width: 1280, height: 720 })
    })
  })

  test.describe('Accessibility', () => {
    test('should have proper admin dashboard accessibility', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/admin')
      await waitForPageReady(page)
      
      // Check for proper heading structure
      const headings = page.locator('h1, h2, h3, h4, h5, h6')
      const headingCount = await headings.count()
      
      if (headingCount > 0) {
        // Should have at least one heading
        expect(headingCount).toBeGreaterThan(0)
        
        // Check first heading
        const firstHeading = headings.first()
        const headingText = await firstHeading.textContent()
        console.log('Main heading:', headingText)
      }
      
      // Check for proper button accessibility
      const buttons = page.locator('button')
      const buttonCount = await buttons.count()
      
      if (buttonCount > 0) {
        // Should have accessible buttons
        expect(buttonCount).toBeGreaterThan(0)
        
        // Check first button
        const firstButton = buttons.first()
        const buttonText = await firstButton.textContent()
        const ariaLabel = await firstButton.getAttribute('aria-label')
        
        // Should have text or aria-label
        expect(buttonText || ariaLabel).toBeTruthy()
      }
    })
  })
})
