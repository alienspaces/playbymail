import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
import StudioLocationLinksView from './StudioLocationLinksView.vue'
import { findInBody, setupModalTestCleanup } from '../../../test-utils/studio-resource-helpers'

// Mock the stores
vi.mock('../../../stores/locations', () => ({
  useLocationsStore: vi.fn(() => ({
    locations: [],
    loading: false,
    error: null,
    fetchLocations: vi.fn()
  }))
}))

vi.mock('../../../stores/locationLinks', () => ({
  useLocationLinksStore: vi.fn(() => ({
    locationLinks: [],
    loading: false,
    error: null,
    createLocationLink: vi.fn(),
    updateLocationLink: vi.fn(),
    deleteLocationLink: vi.fn(),
    fetchLocationLinks: vi.fn()
  }))
}))

vi.mock('../../../stores/games', () => ({
  useGamesStore: vi.fn(() => ({
    selectedGame: ref(null)
  }))
}))

describe('StudioLocationLinksView', () => {
  let pinia
  const modalCleanup = setupModalTestCleanup()

  const setupStoreMocks = async (selectedGame = null) => {
    const { useGamesStore } = await import('../../../stores/games')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useLocationLinksStore } = await import('../../../stores/locationLinks')

    useGamesStore.mockReturnValue({
      selectedGame: ref(selectedGame)
    })
    useLocationsStore.mockReturnValue({
      locations: [],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })
    useLocationLinksStore.mockReturnValue({
      locationLinks: [],
      loading: false,
      error: null,
      createLocationLink: vi.fn(),
      updateLocationLink: vi.fn(),
      deleteLocationLink: vi.fn(),
      fetchLocationLinks: vi.fn()
    })
  }

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    vi.clearAllMocks()
    modalCleanup.beforeEach()
  })

  afterEach(() => {
    modalCleanup.afterEach()
  })

  it('renders prompt when no game is selected', () => {
    const wrapper = mount(StudioLocationLinksView)

    expect(wrapper.text()).toContain('Select a game to manage location links.')
    expect(wrapper.find('.game-table-section').exists()).toBe(false)
  })

  it('renders location links table when game is selected', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    const contextLabel = wrapper.find('.game-context-label')
    const contextName = wrapper.find('.game-context-name')
    expect(contextLabel.text()).toBe('Game:')
    expect(contextName.text()).toBe('Test Game')
    expect(wrapper.find('h2').text()).toBe('Location Links')
  })

  it('renders create button when game is selected', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    const createButton = wrapper.find('button')
    expect(createButton.text()).toBe('Create New Location Link')
  })

  it('renders ResourceTable with correct props', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    const resourceTable = wrapper.findComponent({ name: 'ResourceTable' })
    expect(resourceTable.exists()).toBe(true)
    expect(resourceTable.props('columns')).toEqual([
      { key: 'name', label: 'Link Name' },
      { key: 'from_location_name', label: 'From Location' },
      { key: 'to_location_name', label: 'To Location' },
      { key: 'description', label: 'Description' },
      { key: 'created_at', label: 'Created' }
    ])
  })

  it('opens create modal when create button is clicked', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    const createButton = wrapper.find('button')
    await createButton.trigger('click')

    expect(wrapper.vm.showModal).toBe(true)
    expect(wrapper.vm.modalMode).toBe('create')
  })

  it('renders create modal with correct title', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    // Open modal
    wrapper.vm.showModal = true
    wrapper.vm.modalMode = 'create'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal h2').textContent).toBe('Create Location Link')
  })

  it('renders edit modal with correct title', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    // Open modal in edit mode
    wrapper.vm.showModal = true
    wrapper.vm.modalMode = 'edit'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal h2').textContent).toBe('Edit Location Link')
  })

  it('renders location select options', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    // Override the locations for this specific test
    const { useLocationsStore } = await import('../../../stores/locations')

    useLocationsStore.mockReturnValue({
      locations: [
        { id: 1, name: 'Cave' },
        { id: 2, name: 'Forest' }
      ],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })

    const wrapper = mount(StudioLocationLinksView)

    // Open modal
    wrapper.vm.showModal = true
    await wrapper.vm.$nextTick()

    const fromSelect = findInBody('#from_adventure_game_location_id')
    const toSelect = findInBody('#to_adventure_game_location_id')

    expect(fromSelect).toBeTruthy()
    expect(toSelect).toBeTruthy()

    const fromOptions = fromSelect.querySelectorAll('option')
    const toOptions = toSelect.querySelectorAll('option')

    expect(fromOptions).toHaveLength(3) // placeholder + 2 locations
    expect(toOptions).toHaveLength(3) // placeholder + 2 locations
  })

  it('closes modal when cancel button is clicked', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    // Open modal
    wrapper.vm.showModal = true
    await wrapper.vm.$nextTick()

    const cancelButton = findInBody('button[type="button"]')
    cancelButton.click()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showModal).toBe(false)
  })

  it('displays error message in modal', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    // Open modal and set error
    wrapper.vm.showModal = true
    wrapper.vm.modalError = 'Validation failed'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal .error').textContent).toBe('Validation failed')
  })

  it('renders delete confirmation modal', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    // Set delete target and open confirmation
    wrapper.vm.deleteTarget = { id: 1, name: 'North Path' }
    wrapper.vm.showDeleteConfirm = true
    await wrapper.vm.$nextTick()

    // Check that ConfirmationModal is rendered with correct props
    const confirmationModal = wrapper.findComponent({ name: 'ConfirmationModal' })
    expect(confirmationModal.exists()).toBe(true)
    expect(confirmationModal.props('visible')).toBe(true)
    expect(confirmationModal.props('title')).toBe('Delete Location Link')
    expect(confirmationModal.props('message')).toContain("Are you sure you want to delete the link 'North Path'?")
  })

  it('has correct CSS classes for styling', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioLocationLinksView)

    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    expect(wrapper.find('.game-context-name').exists()).toBe(true)
    expect(wrapper.findComponent({ name: 'PageHeader' }).exists()).toBe(true)
  })

  it('watches for selectedGame changes', async () => {
    const selectedGameRef = ref(null)
    const { useGamesStore } = await import('../../../stores/games')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useLocationLinksStore } = await import('../../../stores/locationLinks')

    useGamesStore.mockReturnValue({
      selectedGame: selectedGameRef
    })
    useLocationsStore.mockReturnValue({
      locations: [],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })
    useLocationLinksStore.mockReturnValue({
      locationLinks: [],
      loading: false,
      error: null,
      createLocationLink: vi.fn(),
      updateLocationLink: vi.fn(),
      deleteLocationLink: vi.fn(),
      fetchLocationLinks: vi.fn()
    })

    const wrapper = mount(StudioLocationLinksView)

    // Initially no game selected
    expect(wrapper.text()).toContain('Select a game to manage location links')

    // Change selected game
    selectedGameRef.value = { id: 1, name: 'New Game' }
    await wrapper.vm.$nextTick()

    const contextName = wrapper.find('.game-context-name')
    expect(contextName.text()).toBe('New Game')
  })
})

