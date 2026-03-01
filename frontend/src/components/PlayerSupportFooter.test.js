import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PlayerSupportFooter from './PlayerSupportFooter.vue'

describe('PlayerSupportFooter', () => {
  it('renders support email link', () => {
    const wrapper = mount(PlayerSupportFooter)
    const link = wrapper.find('a[href="mailto:support@playbymail.games"]')
    expect(link.exists()).toBe(true)
    expect(link.text()).toBe('support@playbymail.games')
  })

  it('renders copyright with current year', () => {
    const wrapper = mount(PlayerSupportFooter)
    const currentYear = new Date().getFullYear()
    expect(wrapper.text()).toContain(`© ${currentYear} PlayByMail`)
  })

  it('renders help text', () => {
    const wrapper = mount(PlayerSupportFooter)
    expect(wrapper.text()).toContain('Need help?')
  })

  it('has correct footer styling class', () => {
    const wrapper = mount(PlayerSupportFooter)
    expect(wrapper.find('.player-support-footer').exists()).toBe(true)
  })
})
