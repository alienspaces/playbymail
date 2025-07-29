#!/usr/bin/env node

import { execSync } from 'child_process'
import { writeFileSync } from 'fs'

// Get the current commit short ref
const getCommitRef = () => {
  try {
    return execSync('git rev-parse --short HEAD', { encoding: 'utf8' }).trim()
  } catch {
    console.warn('Could not get git commit ref, using "dev"')
    return 'dev'
  }
}

// Get current date and time
const now = new Date()
const buildDate = now.toISOString()
const buildTime = now.toISOString()

// Create environment variables
const envVars = {
  VITE_COMMIT_REF: getCommitRef(),
  VITE_BUILD_DATE: buildDate,
  VITE_BUILD_TIME: buildTime
}

// Write to .env.build file
const envContent = Object.entries(envVars)
  .map(([key, value]) => `${key}=${value}`)
  .join('\n')

writeFileSync('.env.build', envContent)

console.log('Build info generated:')
console.log(`Commit: ${envVars.VITE_COMMIT_REF}`)
console.log(`Build Date: ${envVars.VITE_BUILD_DATE}`)
console.log(`Build Time: ${envVars.VITE_BUILD_TIME}`)