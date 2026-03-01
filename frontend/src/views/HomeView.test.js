import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import HomeView from './HomeView.vue'

// Mock router-link component
const RouterLinkStub = {
  template: '<a :href="to"><slot /></a>',
  props: ['to']
}

// Mock router for testing - not used in current tests but available for future use
// const router = createRouter({
//   history: createWebHistory(),
//   routes: [
//     { path: '/', component: HomeView },
//     { path: '/studio', component: { template: '<div>Studio</div>' } },
//     { path: '/admin', component: { template: '<div>Admin</div>' } }
//   ]
// })

describe('HomeView', () => {
  const mountWithStubs = () => mount(HomeView, {
    global: {
      stubs: {
        'router-link': RouterLinkStub
      }
    }
  })

  it('renders the main heading', () => {
    const wrapper = mountWithStubs()
    expect(wrapper.find('h1').text()).toBe('Welcome to PlayByMail')
  })

  it('renders the main description', () => {
    const wrapper = mountWithStubs()
    const description = wrapper.find('p')
    expect(description.text()).toContain('PlayByMail')
    expect(description.text()).toContain('turn-based games')
  })

  it('renders the who section with three user types', () => {
    const wrapper = mountWithStubs()
    
    expect(wrapper.find('.who-section').exists()).toBe(true)
    expect(wrapper.find('.who-section h2').text()).toBe('Who is this platform for?')
    
    const userTypes = wrapper.findAll('.user-type')
    expect(userTypes).toHaveLength(3)
    
    // Check user type headings
    const headings = wrapper.findAll('.user-type h3')
    expect(headings[0].text()).toBe('Game Designers')
    expect(headings[1].text()).toBe('Game Managers')
    expect(headings[2].text()).toBe('Players')
  })

  it('renders the sections with navigation links', () => {
    const wrapper = mountWithStubs()

    const sections = wrapper.findAll('.section')
    expect(sections).toHaveLength(3)

    // Check section headings
    const sectionHeadings = wrapper.findAll('.section h2')
    expect(sectionHeadings[0].text()).toBe('Play a Game')
    expect(sectionHeadings[1].text()).toBe('Game Designer Studio')
    expect(sectionHeadings[2].text()).toBe('Game Management')

    // Check navigation links
    const links = wrapper.findAll('.section-link')
    expect(links).toHaveLength(3)
    expect(links[0].text()).toBe('Browse Games')
    expect(links[1].text()).toBe('Go to Studio')
    expect(links[2].text()).toBe('Go to Game Management')
  })

  it('has correct navigation link destinations', () => {
    const wrapper = mountWithStubs()

    expect(wrapper.find('a[href="/games"]').exists()).toBe(true)
    expect(wrapper.find('a[href="/studio"]').exists()).toBe(true)
    expect(wrapper.find('a[href="/admin"]').exists()).toBe(true)
  })

  it('renders the divider between sections', () => {
    const wrapper = mountWithStubs()
    expect(wrapper.find('.divider').exists()).toBe(true)
  })

  it('has the correct CSS classes for styling', () => {
    const wrapper = mountWithStubs()
    
    expect(wrapper.find('.home-view').exists()).toBe(true)
    expect(wrapper.find('.card').exists()).toBe(true)
    expect(wrapper.find('.who-section').exists()).toBe(true)
    expect(wrapper.find('.user-types').exists()).toBe(true)
    expect(wrapper.find('.sections').exists()).toBe(true)
  })

  it('renders user type descriptions correctly', () => {
    const wrapper = mountWithStubs()
    
    const userTypeDescriptions = wrapper.findAll('.user-type p')
    expect(userTypeDescriptions).toHaveLength(3)
    
    expect(userTypeDescriptions[0].text()).toContain('Create and publish')
    expect(userTypeDescriptions[1].text()).toContain('Run games')
    expect(userTypeDescriptions[2].text()).toContain('Join games')
  })

  it('renders section descriptions correctly', () => {
    const wrapper = mountWithStubs()

    const sectionDescriptions = wrapper.findAll('.section p')
    expect(sectionDescriptions).toHaveLength(3)

    expect(sectionDescriptions[0].text()).toContain('Browse available play-by-mail games')
    expect(sectionDescriptions[1].text()).toContain('Create, edit, and manage')
    expect(sectionDescriptions[2].text()).toContain('Subscribe to or purchase')
  })

  it('has proper semantic structure', () => {
    const wrapper = mountWithStubs()

    // Check for proper heading hierarchy
    expect(wrapper.find('h1').exists()).toBe(true)
    expect(wrapper.findAll('h2')).toHaveLength(5) // genres section + who section + 3 sections
    expect(wrapper.findAll('h3')).toHaveLength(3) // 3 user types
  })

  it('renders a Browse Games link to the game catalog', () => {
    const wrapper = mountWithStubs()

    const catalogLink = wrapper.find('[data-testid="home-catalog-link"]')
    expect(catalogLink.exists()).toBe(true)
    expect(catalogLink.text()).toBe('Browse Games')
  })

  it('renders without any JavaScript errors', () => {
    const wrapper = mountWithStubs()
    expect(wrapper.vm).toBeDefined()
  })
}) 