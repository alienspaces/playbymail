#!/usr/bin/env node

/* eslint-env node */
 

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

// Set environment variables for the current process
 
Object.entries(envVars).forEach(([key, value]) => {
   
  process.env[key] = value
})

// Write to .env.build file for reference
const envContent = Object.entries(envVars)
  .map(([key, value]) => `${key}=${value}`)
  .join('\n')

writeFileSync('.env.build', envContent)

console.log('Build info generated:')
console.log(`Commit: ${envVars.VITE_COMMIT_REF}`)
console.log(`Build Date: ${envVars.VITE_BUILD_DATE}`)
console.log(`Build Time: ${envVars.VITE_BUILD_TIME}`)