import { test, expect } from '@playwright/test'

test.describe('Admin Loading Bug', () => {
  test('should reproduce loading stuck issue when returning to admin dashboard', async ({ page }) => {
    // This test reproduces the exact issue: admin dashboard gets stuck on "Loading games..."
    // when returning from game instances management
    
    // Step 1: Navigate to admin dashboard
    await page.goto('/admin')
    await page.waitForLoadState('networkidle')
    
    // Check what's displayed
    const initialContent = await page.textContent('body')
    console.log('Initial admin page content:', initialContent)
    
    // If showing unauthenticated view, we need to authenticate first
    if (initialContent.includes('Sign In to Game Management')) {
      console.log('User not authenticated, need to login first')
      
      // For now, let's just document what we see
      await expect(page.getByText('Sign In to Game Management')).toBeVisible()
      
      // The real issue happens when:
      // 1. User is authenticated and sees games dashboard
      // 2. User navigates to manage instances
      // 3. User returns to admin dashboard
      // 4. Dashboard gets stuck on "Loading games..."
      
      console.log('This test needs authentication to reproduce the loading bug')
      return
    }
    
    // If authenticated, proceed with the navigation test
    if (initialContent.includes('Games & Instances')) {
      console.log('User is authenticated, proceeding with navigation test')
      
      // Wait for games to load
      if (initialContent.includes('Loading games...')) {
        await expect(page.getByText('Loading games...')).not.toBeVisible({ timeout: 10000 })
      }
      
      // Look for manage button
      const buttons = await page.locator('button').allTextContents()
      console.log('Available buttons:', buttons)
      
      if (buttons.includes('Manage Instances')) {
        // Click manage instances
        const manageButton = page.getByRole('button', { name: 'Manage Instances' }).first()
        await manageButton.click()
        
        // Should navigate to game instances page
        await expect(page).toHaveURL(/\/admin\/games\/.*\/instances/)
        await page.waitForLoadState('networkidle')
        
        // Now navigate back to admin dashboard
        await page.goto('/admin')
        await page.waitForLoadState('networkidle')
        
        // Check if we're stuck on loading
        const returnContent = await page.textContent('body')
        console.log('Content after returning to admin:', returnContent)
        
        // This is the bug: should NOT be stuck on "Loading games..."
        if (returnContent.includes('Loading games...')) {
          console.log('BUG REPRODUCED: Dashboard stuck on "Loading games..."')
          
          // Wait a bit to see if it resolves
          await page.waitForTimeout(5000)
          
          const finalContent = await page.textContent('body')
          if (finalContent.includes('Loading games...')) {
            throw new Error('Dashboard stuck on loading state - this is the bug!')
          }
        } else {
          console.log('No loading issue detected')
        }
      } else {
        console.log('No manage button found')
      }
    }
  })
})
