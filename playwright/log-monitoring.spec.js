import { test, expect } from '@playwright/test'
import { extractVerificationToken } from './utils/log-parser.js'

test.describe('Log Monitoring', () => {
  test('should parse verification tokens from log patterns', () => {
    // Test with the expected log pattern
    const sampleLogs = 'verification token >ABC123< for account ID >12345<'
    const token = extractVerificationToken(sampleLogs, 'test@example.com')
    
    expect(token).toBe('ABC123')
  })

  test('should parse verification tokens from alternative patterns', () => {
    // Test with just the token pattern
    const sampleLogs = 'some log message verification token >XYZ789< another message'
    const token = extractVerificationToken(sampleLogs, 'test@example.com')
    
    expect(token).toBe('XYZ789')
  })

  test('should handle empty logs', () => {
    const token = extractVerificationToken('', 'test@example.com')
    expect(token).toBeNull()
  })

  test('should handle logs without tokens', () => {
    const logs = 'some random log message without verification tokens'
    const token = extractVerificationToken(logs, 'test@example.com')
    expect(token).toBeNull()
  })

  test('should handle null logs', () => {
    const token = extractVerificationToken(null, 'test@example.com')
    expect(token).toBeNull()
  })
})
