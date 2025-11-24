import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PageHeader from './PageHeader.vue'

describe('PageHeader', () => {
  it('renders title correctly', () => {
    const wrapper = mount(PageHeader, {
      props: {
        title: 'Test Page'
      }
    })
    
    expect(wrapper.find('h1').text()).toBe('Test Page')
  })

  it('renders action button when actionText is provided', () => {
    const wrapper = mount(PageHeader, {
      props: {
        title: 'Test Page',
        actionText: 'Create New'
      }
    })
    
    const button = wrapper.find('button')
    expect(button.exists()).toBe(true)
    expect(button.text()).toBe('Create New')
  })

  it('does not render action button when actionText is not provided', () => {
    const wrapper = mount(PageHeader, {
      props: {
        title: 'Test Page'
      }
    })
    
    expect(wrapper.find('button').exists()).toBe(false)
  })

  it('emits action event when button is clicked', async () => {
    const wrapper = mount(PageHeader, {
      props: {
        title: 'Test Page',
        actionText: 'Create New'
      }
    })
    
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toBeTruthy()
  })

  it('renders subtitle when provided', () => {
    const wrapper = mount(PageHeader, {
      props: {
        title: 'Test Page',
        subtitle: 'Additional context'
      }
    })

    const subtitle = wrapper.find('.page-subtitle')
    expect(subtitle.exists()).toBe(true)
    expect(subtitle.text()).toBe('Additional context')
  })
}) 