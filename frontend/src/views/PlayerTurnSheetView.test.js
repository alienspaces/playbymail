import { describe, it, expect, vi, beforeEach } from 'vitest'
import { nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import PlayerTurnSheetView from './PlayerTurnSheetView.vue'

const mockGetGSITurnSheets = vi.fn()
const mockSubmitGSITurnSheets = vi.fn()
const mockDownloadGSITurnSheetPDF = vi.fn()

vi.mock('../api/player', () => ({
  getGSITurnSheets: (...args) => mockGetGSITurnSheets(...args),
  submitGSITurnSheets: (...args) => mockSubmitGSITurnSheets(...args),
  downloadGSITurnSheetPDF: (...args) => mockDownloadGSITurnSheetPDF(...args),
}))

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({
    params: { game_subscription_instance_id: 'gsi-abc-123' },
  })),
  useRouter: vi.fn(() => ({ push: vi.fn() })),
}))

const mockSheets = [
  { id: 'ts-1', sheet_type: 'location_choice', turn_number: 1, is_completed: false },
  { id: 'ts-2', sheet_type: 'inventory_management', turn_number: 1, is_completed: false },
]

describe('PlayerTurnSheetView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading state while fetching', async () => {
    mockGetGSITurnSheets.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(PlayerTurnSheetView)
    await nextTick()

    expect(wrapper.find('[data-testid="ts-loading"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(false)
  })

  it('shows turn sheet list after successful load', async () => {
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    expect(wrapper.find('[data-testid="ts-loading"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="ts-title"]').text()).toBe('Your Turn Sheets')
    expect(wrapper.find('[data-testid="ts-card-ts-1"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="ts-card-ts-2"]').exists()).toBe(true)
  })

  it('shows empty state when no turn sheets', async () => {
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: [] })

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    expect(wrapper.find('[data-testid="ts-empty"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="btn-submit"]').exists()).toBe(false)
  })

  it('shows load error when fetch fails', async () => {
    mockGetGSITurnSheets.mockRejectedValue(new Error('Session expired'))

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    expect(wrapper.find('[data-testid="ts-load-error"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Session expired')
    expect(wrapper.find('[data-testid="link-browse-games"]').exists()).toBe(true)
  })

  it('shows pending and completed status for each sheet', async () => {
    const sheets = [
      { id: 'ts-1', sheet_type: 'location_choice', turn_number: 1, is_completed: false },
      { id: 'ts-2', sheet_type: 'inventory_management', turn_number: 1, is_completed: true },
    ]
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: sheets })

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    expect(wrapper.find('[data-testid="ts-status-ts-1"]').text()).toBe('Pending')
    expect(wrapper.find('[data-testid="ts-status-ts-2"]').text()).toBe('Completed')
  })

  it('submit button is disabled when all sheets are already completed', async () => {
    const completedSheets = mockSheets.map((s) => ({ ...s, is_completed: true }))
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: completedSheets })

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    expect(wrapper.find('[data-testid="btn-submit"]').attributes('disabled')).toBeDefined()
    expect(wrapper.find('[data-testid="btn-submit"]').text()).toBe('Already submitted')
  })

  it('submit button is enabled when there are pending sheets', async () => {
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    expect(wrapper.find('[data-testid="btn-submit"]').attributes('disabled')).toBeUndefined()
  })

  it('shows success state after submission', async () => {
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
    mockSubmitGSITurnSheets.mockResolvedValue({ submitted_count: 2, total_count: 2 })

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-submit"]').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="ts-success"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="link-browse-more"]').attributes('href')).toBe('/games')
    expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(false)
  })

  it('shows submit error when submission fails', async () => {
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
    mockSubmitGSITurnSheets.mockRejectedValue(new Error('Server error'))

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-submit"]').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="submit-error"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Server error')
    expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(true)
  })

  it('renders a Download PDF button for each turn sheet', async () => {
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    expect(wrapper.find('[data-testid="btn-download-ts-1"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="btn-download-ts-2"]').exists()).toBe(true)
  })

  it('calls downloadGSITurnSheetPDF when Download PDF is clicked', async () => {
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

    const mockBlob = new Blob(['%PDF'], { type: 'application/pdf' })
    mockDownloadGSITurnSheetPDF.mockResolvedValue({
      blob: () => Promise.resolve(mockBlob),
    })

    // jsdom doesn't support URL.createObjectURL; mock it
    URL.createObjectURL = vi.fn(() => 'blob:test-url')
    URL.revokeObjectURL = vi.fn()

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-download-ts-1"]').trigger('click')
    await flushPromises()

    expect(mockDownloadGSITurnSheetPDF).toHaveBeenCalledWith('gsi-abc-123', 'ts-1')
  })

  it('calls submitGSITurnSheets with the correct gsi id', async () => {
    mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
    mockSubmitGSITurnSheets.mockResolvedValue({ submitted_count: 2, total_count: 2 })

    const wrapper = mount(PlayerTurnSheetView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-submit"]').trigger('click')
    await flushPromises()

    expect(mockSubmitGSITurnSheets).toHaveBeenCalledWith('gsi-abc-123')
  })
})
