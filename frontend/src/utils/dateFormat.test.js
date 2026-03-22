import { describe, it, expect } from 'vitest'
import { formatDateTime, formatDateOnly, formatDeadline } from './dateFormat'

const FIXED_DATE = '2026-03-22T14:30:00Z'

describe('dateFormat', () => {
  describe('formatDateTime', () => {
    it('returns N/A for null', () => {
      expect(formatDateTime(null)).toBe('N/A')
    })

    it('returns N/A for undefined', () => {
      expect(formatDateTime(undefined)).toBe('N/A')
    })

    it('returns N/A for empty string', () => {
      expect(formatDateTime('')).toBe('N/A')
    })

    it('formats a UTC ISO string into a human-readable date+time', () => {
      const result = formatDateTime(FIXED_DATE, { timezone: 'UTC', locale: 'en-US' })
      expect(result).toMatch(/2026/)
      expect(result).toMatch(/Mar/)
      expect(result).toMatch(/22/)
      expect(result).toMatch(/\d{2}:\d{2}/)
    })

    it('respects the supplied timezone', () => {
      // UTC+10 should shift the hour forward
      const utcResult = formatDateTime(FIXED_DATE, { timezone: 'UTC', locale: 'en-US' })
      const sydneyResult = formatDateTime(FIXED_DATE, {
        timezone: 'Australia/Sydney',
        locale: 'en-US',
      })
      expect(utcResult).not.toBe(sydneyResult)
    })

    it('uses browser timezone when none supplied', () => {
      // Should not throw and should return a non-empty string
      const result = formatDateTime(FIXED_DATE)
      expect(result).toBeTruthy()
      expect(result).not.toBe('N/A')
    })

    it('includes both date and time components', () => {
      const result = formatDateTime(FIXED_DATE, { timezone: 'UTC', locale: 'en-US' })
      // Should contain a time separator (colon for HH:MM)
      expect(result).toMatch(/\d{1,2}:\d{2}/)
    })
  })

  describe('formatDateOnly', () => {
    it('returns N/A for null', () => {
      expect(formatDateOnly(null)).toBe('N/A')
    })

    it('returns N/A for empty string', () => {
      expect(formatDateOnly('')).toBe('N/A')
    })

    it('formats a UTC ISO string to date only', () => {
      const result = formatDateOnly(FIXED_DATE, { timezone: 'UTC', locale: 'en-US' })
      expect(result).toMatch(/2026/)
      expect(result).toMatch(/Mar/)
      expect(result).toMatch(/22/)
    })

    it('does not include a time component', () => {
      const result = formatDateOnly(FIXED_DATE, { timezone: 'UTC', locale: 'en-US' })
      // Date-only should not contain AM/PM or HH:MM pattern beyond years
      expect(result).not.toMatch(/\d{1,2}:\d{2}/)
    })

    it('respects the supplied timezone when computing the date', () => {
      // A UTC time near midnight may land on different dates in different timezones
      const nearMidnight = '2026-03-22T23:30:00Z'
      const utcResult = formatDateOnly(nearMidnight, { timezone: 'UTC', locale: 'en-US' })
      const nzResult = formatDateOnly(nearMidnight, {
        timezone: 'Pacific/Auckland',
        locale: 'en-US',
      })
      // NZ is UTC+13 in March, so 23:30 UTC = 12:30 next day NZ
      expect(utcResult).not.toBe(nzResult)
    })
  })

  describe('formatDeadline', () => {
    it('returns N/A for null', () => {
      expect(formatDeadline(null)).toBe('N/A')
    })

    it('returns N/A for empty string', () => {
      expect(formatDeadline('')).toBe('N/A')
    })

    it('returns Overdue for a past date', () => {
      const past = new Date(Date.now() - 60 * 60 * 1000).toISOString()
      expect(formatDeadline(past)).toBe('Overdue')
    })

    it('returns Today for a deadline within 24 hours', () => {
      const soon = new Date(Date.now() + 60 * 60 * 1000).toISOString()
      const result = formatDeadline(soon)
      expect(result).toMatch(/^Today,/)
    })

    it('includes time in the Today label', () => {
      const soon = new Date(Date.now() + 60 * 60 * 1000).toISOString()
      const result = formatDeadline(soon)
      expect(result).toMatch(/\d{1,2}:\d{2}/)
    })

    it('returns Tomorrow for a deadline between 24 and 48 hours away', () => {
      const tomorrow = new Date(Date.now() + 30 * 60 * 60 * 1000).toISOString()
      const result = formatDeadline(tomorrow)
      expect(result).toMatch(/^Tomorrow,/)
    })

    it('includes time in the Tomorrow label', () => {
      const tomorrow = new Date(Date.now() + 30 * 60 * 60 * 1000).toISOString()
      const result = formatDeadline(tomorrow)
      expect(result).toMatch(/\d{1,2}:\d{2}/)
    })

    it('returns a full date+time for deadlines more than 48 hours away', () => {
      const future = new Date(Date.now() + 72 * 60 * 60 * 1000).toISOString()
      const result = formatDeadline(future, { timezone: 'UTC', locale: 'en-US' })
      expect(result).not.toMatch(/^Today/)
      expect(result).not.toMatch(/^Tomorrow/)
      expect(result).not.toBe('Overdue')
      expect(result).toMatch(/\d{4}/)
    })

    it('passes timezone through to the formatted output', () => {
      const future = new Date(Date.now() + 72 * 60 * 60 * 1000).toISOString()
      const utcResult = formatDeadline(future, { timezone: 'UTC', locale: 'en-US' })
      const sydneyResult = formatDeadline(future, {
        timezone: 'Australia/Sydney',
        locale: 'en-US',
      })
      expect(utcResult).not.toBe(sydneyResult)
    })
  })
})
