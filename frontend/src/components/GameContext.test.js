import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import GameContext from './GameContext.vue'

describe('GameContext', () => {
  it('renders game name when provided', () => {
    const wrapper = mount(GameContext, {
      props: {
        gameName: 'Test Game'
      }
    })
    
    expect(wrapper.find('p').text()).toBe('Game: Test Game')
  })

  it('does not render when gameName is empty', () => {
    const wrapper = mount(GameContext, {
      props: {
        gameName: ''
      }
    })
    
    expect(wrapper.find('p').exists()).toBe(false)
  })

  it('does not render when gameName is not provided', () => {
    const wrapper = mount(GameContext)
    
    expect(wrapper.find('p').exists()).toBe(false)
  })
}) 