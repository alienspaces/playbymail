import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
import StudioItemPlacementsView from './StudioItemPlacementsView.vue'
import { findInBody, setupModalTestCleanup } from '../../../test-utils/studio-resource-helpers'

// Mock the stores
vi.mock('../../../stores/items', () => ({
  useItemsStore: vi.fn(() => ({
    items: [],
    loading: false,
    error: null,
    fetchItems: vi.fn()
  }))
}))

vi.mock('../../../stores/locations', () => ({
  useLocationsStore: vi.fn(() => ({
    locations: [],
    loading: false,
    error: null,
    fetchLocations: vi.fn()
  }))
}))

vi.mock('../../../stores/itemPlacements', () => ({
  useItemPlacementsStore: vi.fn(() => ({
    itemPlacements: [],
    loading: false,
    error: null,
    createItemPlacement: vi.fn(),
    updateItemPlacement: vi.fn(),
    deleteItemPlacement: vi.fn(),
    fetchItemPlacements: vi.fn()
  }))
}))

vi.mock('../../../stores/games', () => ({
  useGamesStore: vi.fn(() => ({
    selectedGame: ref(null)
  }))
}))

describe('StudioItemPlacementsView', () => {
  let pinia
  const modalCleanup = setupModalTestCleanup()

  const setupStoreMocks = async (selectedGame = null) => {
    const { useGamesStore } = await import('../../../stores/games')
    const { useItemsStore } = await import('../../../stores/items')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useItemPlacementsStore } = await import('../../../stores/itemPlacements')

    useGamesStore.mockReturnValue({
      selectedGame: ref(selectedGame)
    })
    useItemsStore.mockReturnValue({
      items: [],
      loading: false,
      error: null,
      fetchItems: vi.fn()
    })
    useLocationsStore.mockReturnValue({
      locations: [],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })
    useItemPlacementsStore.mockReturnValue({
      itemPlacements: [],
      loading: false,
      error: null,
      createItemPlacement: vi.fn(),
      updateItemPlacement: vi.fn(),
      deleteItemPlacement: vi.fn(),
      fetchItemPlacements: vi.fn()
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
    const wrapper = mount(StudioItemPlacementsView)

    expect(wrapper.text()).toContain('Please select or create a game to manage item placements.')
    expect(wrapper.find('.game-table-section').exists()).toBe(false)
  })

  it('renders item placements table when game is selected', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    const contextLabel = wrapper.find('.game-context-label')
    const contextName = wrapper.find('.game-context-name')
    expect(contextLabel.text()).toBe('Game:')
    expect(contextName.text()).toBe('Test Game')
    expect(wrapper.find('h2').text()).toBe('Item Placements')
  })

  it('renders create button when game is selected', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    const createButton = wrapper.find('button')
    expect(createButton.text()).toBe('Create Item Placement')
  })

  it('renders ResourceTable with correct props', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    const resourceTable = wrapper.findComponent({ name: 'ResourceTable' })
    expect(resourceTable.exists()).toBe(true)
    expect(resourceTable.props('columns')).toEqual([
      { key: 'item_name', label: 'Item' },
      { key: 'location_name', label: 'Location' },
      { key: 'initial_count', label: 'Count' },
      { key: 'created_at', label: 'Created' }
    ])
  })

  it('opens create modal when create button is clicked', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    const createButton = wrapper.find('button')
    await createButton.trigger('click')

    expect(wrapper.vm.showItemPlacementModal).toBe(true)
    expect(wrapper.vm.itemPlacementModalMode).toBe('create')
  })

  it('renders create modal with correct title', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    // Open modal
    wrapper.vm.showItemPlacementModal = true
    wrapper.vm.itemPlacementModalMode = 'create'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal h2').textContent).toBe('Create Item Placement')
  })

  it('renders edit modal with correct title', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    // Open modal in edit mode
    wrapper.vm.showItemPlacementModal = true
    wrapper.vm.itemPlacementModalMode = 'edit'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal h2').textContent).toBe('Edit Item Placement')
  })

  it('renders item and location select options', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    // Override the items and locations for this specific test
    const { useItemsStore } = await import('../../../stores/items')
    const { useLocationsStore } = await import('../../../stores/locations')

    useItemsStore.mockReturnValue({
      items: [
        { id: 1, name: 'Sword' },
        { id: 2, name: 'Shield' }
      ],
      loading: false,
      error: null,
      fetchItems: vi.fn()
    })
    useLocationsStore.mockReturnValue({
      locations: [
        { id: 1, name: 'Cave' },
        { id: 2, name: 'Forest' }
      ],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })

    const wrapper = mount(StudioItemPlacementsView)

    // Open modal
    wrapper.vm.showItemPlacementModal = true
    await wrapper.vm.$nextTick()

    const itemSelect = findInBody('#adventure_game_item_id')
    const locationSelect = findInBody('#adventure_game_location_id')

    expect(itemSelect).toBeTruthy()
    expect(locationSelect).toBeTruthy()

    const itemOptions = itemSelect.querySelectorAll('option')
    const locationOptions = locationSelect.querySelectorAll('option')

    expect(itemOptions).toHaveLength(3) // placeholder + 2 items
    expect(locationOptions).toHaveLength(3) // placeholder + 2 locations
  })

  it('closes modal when cancel button is clicked', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    // Open modal
    wrapper.vm.showItemPlacementModal = true
    await wrapper.vm.$nextTick()

    const cancelButton = findInBody('button[type="button"]')
    cancelButton.click()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showItemPlacementModal).toBe(false)
  })

  it('displays error message in modal', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    // Open modal and set error
    wrapper.vm.showItemPlacementModal = true
    wrapper.vm.itemPlacementModalError = 'Validation failed'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal .error').textContent).toBe('Validation failed')
  })

  it('renders delete confirmation modal', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    // Open delete confirmation
    wrapper.vm.showItemPlacementDeleteConfirm = true
    await wrapper.vm.$nextTick()

    // Check that ConfirmationModal is rendered with correct props
    const confirmationModal = wrapper.findComponent({ name: 'ConfirmationModal' })
    expect(confirmationModal.exists()).toBe(true)
    expect(confirmationModal.props('visible')).toBe(true)
    expect(confirmationModal.props('title')).toBe('Delete Item Placement')
    expect(confirmationModal.props('message')).toBe('Are you sure you want to delete this item placement?')
  })

  it('has correct CSS classes for styling', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioItemPlacementsView)

    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    expect(wrapper.find('.game-context-name').exists()).toBe(true)
    expect(wrapper.findComponent({ name: 'PageHeader' }).exists()).toBe(true)
  })

  it('watches for selectedGame changes', async () => {
    const selectedGameRef = ref(null)
    const { useGamesStore } = await import('../../../stores/games')
    const { useItemsStore } = await import('../../../stores/items')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useItemPlacementsStore } = await import('../../../stores/itemPlacements')

    useGamesStore.mockReturnValue({
      selectedGame: selectedGameRef
    })
    useItemsStore.mockReturnValue({
      items: [],
      loading: false,
      error: null,
      fetchItems: vi.fn()
    })
    useLocationsStore.mockReturnValue({
      locations: [],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })
    useItemPlacementsStore.mockReturnValue({
      itemPlacements: [],
      loading: false,
      error: null,
      createItemPlacement: vi.fn(),
      updateItemPlacement: vi.fn(),
      deleteItemPlacement: vi.fn(),
      fetchItemPlacements: vi.fn()
    })

    const wrapper = mount(StudioItemPlacementsView)

    // Initially no game selected
    expect(wrapper.text()).toContain('Please select or create a game')

    // Change selected game
    selectedGameRef.value = { id: 1, name: 'New Game' }
    await wrapper.vm.$nextTick()

    const contextName = wrapper.find('.game-context-name')
    expect(contextName.text()).toBe('New Game')
  })
})

