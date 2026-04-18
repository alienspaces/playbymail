import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FieldHint from './FieldHint.vue'

describe('FieldHint', () => {
  it('renders slot content inside a paragraph with the field-hint class', () => {
    const wrapper = mount(FieldHint, {
      slots: { default: 'Sector hops allowed per turn (1–10).' }
    })

    const p = wrapper.find('p')
    expect(p.exists()).toBe(true)
    expect(p.classes()).toContain('field-hint')
    expect(p.text()).toBe('Sector hops allowed per turn (1–10).')
  })

  it('renders nothing visible when no slot content is provided', () => {
    const wrapper = mount(FieldHint)
    expect(wrapper.find('p').text()).toBe('')
  })

  it('forwards fallthrough attributes such as id to the root element', () => {
    const wrapper = mount(FieldHint, {
      attrs: { id: 'speed-hint' },
      slots: { default: 'Help text.' }
    })
    expect(wrapper.find('p').attributes('id')).toBe('speed-hint')
  })
})
