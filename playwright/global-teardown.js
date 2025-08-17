async function globalTeardown(config) {
  console.log('🧹 Cleaning up Playwright test environment...')
  
  // Any global cleanup can go here
  // For now, just log completion
  
  console.log('✅ Global teardown complete')
}

module.exports = globalTeardown
