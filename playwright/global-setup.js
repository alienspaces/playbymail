const { chromium } = require('@playwright/test')

async function globalSetup(config) {
  console.log('🚀 Setting up Playwright test environment...')
  
  // Check if backend is accessible
  const browser = await chromium.launch()
  const page = await browser.newPage()
  
  try {
    // Test backend connectivity
    await page.goto('http://localhost:8080/health')
    const status = await page.textContent('body')
    
    if (status && status.includes('ok')) {
      console.log('✅ Backend is running and healthy')
    } else {
      console.log('⚠️  Backend is running but health check failed')
    }
  } catch (error) {
    console.log('❌ Backend is not accessible - tests may fail')
    console.log('   Make sure to run: ./tools/start-backend')
  }
  
  await browser.close()
  
  // Create necessary directories
  const fs = require('fs')
  const path = require('path')
  
  const dirs = [
    'playwright/screenshots',
    'playwright/videos',
    'test-results'
  ]
  
  dirs.forEach(dir => {
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true })
      console.log(`📁 Created directory: ${dir}`)
    }
  })
  
  console.log('✅ Global setup complete')
}

module.exports = globalSetup
