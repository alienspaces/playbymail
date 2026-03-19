import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
import StudioLocationObjectsView from './StudioLocationObjectsView.vue'
import { findInBody, setupModalTestCleanup } from '../../../test-utils/studio-resource-helpers'

vi.mock('../../../stores/locations', () => ({
  useLocationsStore: vi.fn(() => ({
    locations: [],
    loading: false,
    error: null,
    fetchLocations: vi.fn(),
  })),
}))

vi.mock('../../../stores/locationObjects', () => ({
  useLocationObjectsStore: vi.fn(() => ({
    locationObjects: [],
    loading: false,
    error: null,
    fetchLocationObjects: vi.fn(),
    createLocationObject: vi.fn(),
    updateLocationObject: vi.fn(),
    deleteLocationObject: vi.fn(),
  })),
}))

vi.mock('../../../stores/games', () => ({
  useGamesStore: vi.fn(() => ({
    selectedGame: ref(null),
  })),
}))

describe('StudioLocationObjectsView', () => {
  let pinia
  const modalCleanup = setupModalTestCleanup()

  const setupStoreMocks = async (selectedGame = null) => {
    const { useGamesStore } = await import('../../../stores/games')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useLocationObjectsStore } = await import('../../../stores/locationObjects')

    useGamesStore.mockReturnValue({ selectedGame: ref(selectedGame) })
    useLocationsStore.mockReturnValue({
      locations: [],
      loading: false,
      error: null,
      fetchLocations: vi.fn(),
    })
    useLocationObjectsStore.mockReturnValue({
      locationObjects: [],
      loading: false,
      error: null,
      fetchLocationObjects: vi.fn(),
      createLocationObject: vi.fn(),
      updateLocationObject: vi.fn(),
      deleteLocationObject: vi.fn(),
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
    const wrapper = mount(StudioLocationObjectsView)

    expect(wrapper.text()).toContain('Please select or create a game to manage location objects.')
    expect(wrapper.find('.game-table-section').exists()).toBe(false)
  })

  it('renders location objects table when game is selected', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toBe('Location Objects')
  })

  it('renders create button when game is selected', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    const createButton = wrapper.find('button')
    expect(createButton.text()).toBe('Create Location Object')
  })

  it('renders ResourceTable with correct columns', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    const resourceTable = wrapper.findComponent({ name: 'ResourceTable' })
    expect(resourceTable.exists()).toBe(true)
    expect(resourceTable.props('columns')).toEqual([
      { key: 'name', label: 'Name' },
      { key: 'location_name', label: 'Location' },
      { key: 'initial_state_name', label: 'Initial State' },
      { key: 'is_hidden', label: 'Hidden' },
      { key: 'created_at', label: 'Created' },
    ])
  })

  it('opens create modal when create button is clicked', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    const createButton = wrapper.find('button')
    await createButton.trigger('click')

    expect(wrapper.vm.showModal).toBe(true)
    expect(wrapper.vm.modalMode).toBe('create')
  })

  it('renders create modal with correct title', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    wrapper.vm.showModal = true
    wrapper.vm.modalMode = 'create'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal h2').textContent).toBe('Create Location Object')
  })

  it('renders edit modal with correct title', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    wrapper.vm.showModal = true
    wrapper.vm.modalMode = 'edit'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal h2').textContent).toBe('Edit Location Object')
  })

  it('renders delete confirmation modal with correct props', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    wrapper.vm.showDeleteConfirm = true
    await wrapper.vm.$nextTick()

    const confirmationModal = wrapper.findComponent({ name: 'ConfirmationModal' })
    expect(confirmationModal.exists()).toBe(true)
    expect(confirmationModal.props('visible')).toBe(true)
    expect(confirmationModal.props('title')).toBe('Delete Location Object')
    expect(confirmationModal.props('message')).toBe('Are you sure you want to delete this location object?')
  })

  it('closes modal when cancel is triggered', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    wrapper.vm.showModal = true
    await wrapper.vm.$nextTick()

    const cancelButton = findInBody('button[type="button"]')
    cancelButton.click()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showModal).toBe(false)
  })

  it('displays error message in modal', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectsView)

    wrapper.vm.showModal = true
    wrapper.vm.modalError = 'Validation failed'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal .error').textContent).toBe('Validation failed')
  })

  it('watches for selectedGame changes', async () => {
    const selectedGameRef = ref(null)
    const { useGamesStore } = await import('../../../stores/games')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useLocationObjectsStore } = await import('../../../stores/locationObjects')

    useGamesStore.mockReturnValue({ selectedGame: selectedGameRef })
    useLocationsStore.mockReturnValue({
      locations: [],
      loading: false,
      error: null,
      fetchLocations: vi.fn(),
    })
    useLocationObjectsStore.mockReturnValue({
      locationObjects: [],
      loading: false,
      error: null,
      fetchLocationObjects: vi.fn(),
      createLocationObject: vi.fn(),
      updateLocationObject: vi.fn(),
      deleteLocationObject: vi.fn(),
    })

    const wrapper = mount(StudioLocationObjectsView)

    expect(wrapper.text()).toContain('Please select or create a game')

    selectedGameRef.value = { id: 'game-1', name: 'New Game' }
    await wrapper.vm.$nextTick()

    const contextName = wrapper.find('.game-context-name')
    expect(contextName.text()).toBe('New Game')
  })
})
