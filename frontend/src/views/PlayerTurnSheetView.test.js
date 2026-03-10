import { describe, it, expect, vi, beforeEach } from 'vitest'
import { nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import PlayerTurnSheetView from './PlayerTurnSheetView.vue'

const mockGetGameSubscriptionInstanceTurnSheets = vi.fn()
const mockGetGameSubscriptionInstanceTurnSheetHTML = vi.fn()
const mockSaveGameSubscriptionInstanceTurnSheet = vi.fn()
const mockSubmitGameSubscriptionInstanceTurnSheets = vi.fn()
const mockVerifyGameSubscriptionToken = vi.fn()
const mockRequestNewTurnSheetToken = vi.fn()

vi.mock('../api/player', () => ({
  verifyGameSubscriptionToken: (...args) => mockVerifyGameSubscriptionToken(...args),
  requestNewTurnSheetToken: (...args) => mockRequestNewTurnSheetToken(...args),
  getGameSubscriptionInstanceTurnSheets: (...args) => mockGetGameSubscriptionInstanceTurnSheets(...args),
  getGameSubscriptionInstanceTurnSheetHTML: (...args) => mockGetGameSubscriptionInstanceTurnSheetHTML(...args),
  saveGameSubscriptionInstanceTurnSheet: (...args) => mockSaveGameSubscriptionInstanceTurnSheet(...args),
  submitGameSubscriptionInstanceTurnSheets: (...args) => mockSubmitGameSubscriptionInstanceTurnSheets(...args),
  downloadGameSubscriptionInstanceTurnSheetPDF: vi.fn(),
  uploadGameSubscriptionInstanceTurnSheetScan: vi.fn(),
}))

const mockSetSessionToken = vi.fn()

vi.mock('../stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    sessionToken: '',
    setSessionToken: mockSetSessionToken,
  })),
}))

let routeParams = { game_subscription_instance_id: 'gsi-abc-123' }

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({
    params: routeParams,
  })),
  useRouter: vi.fn(() => ({ push: vi.fn() })),
}))

const mockSheets = [
  { id: 'ts-1', sheet_type: 'adventure_game_location_choice', turn_number: 1, is_completed: false, scanned_data: null },
  { id: 'ts-2', sheet_type: 'adventure_game_inventory_management', turn_number: 1, is_completed: false, scanned_data: null },
]

describe('PlayerTurnSheetView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    routeParams = { game_subscription_instance_id: 'gsi-abc-123' }
    mockGetGameSubscriptionInstanceTurnSheetHTML.mockResolvedValue('<p>Turn sheet</p>')
  })

  describe('loading and auth states', () => {
    it('shows loading state while fetching', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockReturnValue(new Promise(() => { }))
      const wrapper = mount(PlayerTurnSheetView)
      await nextTick()
      expect(wrapper.find('[data-testid="ts-loading"]').exists()).toBe(true)
    })

    it('shows empty state when no turn sheets', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: [] })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()
      expect(wrapper.find('[data-testid="ts-empty"]').exists()).toBe(true)
    })

    it('shows load error when fetch fails', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockRejectedValue(new Error('Server error'))
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()
      expect(wrapper.find('[data-testid="ts-load-error"]').exists()).toBe(true)
      expect(wrapper.text()).toContain('Server error')
    })
  })

  describe('token auto-verify', () => {
    beforeEach(() => {
      routeParams = {
        game_subscription_instance_id: 'gsi-abc-123',
        turn_sheet_token: 'token-xyz',
      }
    })

    it('verifies token and loads turn sheets on success', async () => {
      mockVerifyGameSubscriptionToken.mockResolvedValue('session-abc')
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()
      expect(mockVerifyGameSubscriptionToken).toHaveBeenCalledWith('gsi-abc-123', 'token-xyz')
      expect(mockSetSessionToken).toHaveBeenCalledWith('session-abc')
      expect(wrapper.find('[data-testid="ts-viewer"]').exists()).toBe(true)
    })

    it('shows expired-token UI when verify fails', async () => {
      mockVerifyGameSubscriptionToken.mockRejectedValue(new Error('Token expired'))
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()
      expect(wrapper.find('[data-testid="ts-token-expired"]').exists()).toBe(true)
    })

    it('requests new link on expired-token form submit', async () => {
      mockVerifyGameSubscriptionToken.mockRejectedValue(new Error('Token expired'))
      mockRequestNewTurnSheetToken.mockResolvedValue()
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()
      await wrapper.find('[data-testid="expired-email-input"]').setValue('player@example.com')
      await wrapper.find('[data-testid="expired-request-btn"]').trigger('submit')
      await flushPromises()
      expect(mockRequestNewTurnSheetToken).toHaveBeenCalledWith('gsi-abc-123', 'player@example.com')
      expect(wrapper.find('[data-testid="expired-success"]').exists()).toBe(true)
    })
  })

  describe('stepper and viewer', () => {
    it('renders stepper with correct number of steps', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()
      expect(wrapper.find('[data-testid="ts-stepper"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="ts-step-0"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="ts-step-1"]').exists()).toBe(true)
    })

    it('loads HTML for the first sheet on mount', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()
      expect(mockGetGameSubscriptionInstanceTurnSheetHTML).toHaveBeenCalledWith('gsi-abc-123', 'ts-1')
      expect(wrapper.find('[data-testid="ts-viewer-iframe"]').exists()).toBe(true)
    })

    it('submit-all button is disabled when no sheets are marked ready', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()
      const btn = wrapper.find('[data-testid="btn-submit-all"]')
      expect(btn.exists()).toBe(true)
      expect(btn.attributes('disabled')).toBeDefined()
    })

    it('mark-ready toggles ready state', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      const markBtn = wrapper.find('[data-testid="btn-mark-ready"]')
      expect(markBtn.text()).toBe('Mark Ready')

      await markBtn.trigger('click')
      await nextTick()
      expect(wrapper.find('[data-testid="btn-mark-ready"]').text()).toBe('Unmark Ready')

      await wrapper.find('[data-testid="btn-mark-ready"]').trigger('click')
      await nextTick()
      expect(wrapper.find('[data-testid="btn-mark-ready"]').text()).toBe('Mark Ready')
    })

    it('shows success state after submission', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      mockSubmitGameSubscriptionInstanceTurnSheets.mockResolvedValue({ submitted_count: 2, total_count: 2 })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      // Mark both sheets ready
      await wrapper.find('[data-testid="btn-mark-ready"]').trigger('click')
      await wrapper.find('[data-testid="ts-step-1"]').trigger('click')
      await nextTick()
      await wrapper.find('[data-testid="btn-mark-ready"]').trigger('click')
      await nextTick()

      await wrapper.find('[data-testid="btn-submit-all"]').trigger('click')
      await flushPromises()

      expect(mockSubmitGameSubscriptionInstanceTurnSheets).toHaveBeenCalledWith('gsi-abc-123')
      expect(wrapper.find('[data-testid="ts-success"]').exists()).toBe(true)
    })

    it('shows submit error when submission fails', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      mockSubmitGameSubscriptionInstanceTurnSheets.mockRejectedValue(new Error('Server error'))
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      // Mark both sheets ready
      await wrapper.find('[data-testid="btn-mark-ready"]').trigger('click')
      await wrapper.find('[data-testid="ts-step-1"]').trigger('click')
      await nextTick()
      await wrapper.find('[data-testid="btn-mark-ready"]').trigger('click')
      await nextTick()

      await wrapper.find('[data-testid="btn-submit-all"]').trigger('click')
      await flushPromises()

      expect(wrapper.find('[data-testid="submit-error"]').exists()).toBe(true)
      expect(wrapper.text()).toContain('Server error')
    })

    it('calls saveGameSubscriptionInstanceTurnSheet when save button is clicked', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      mockSaveGameSubscriptionInstanceTurnSheet.mockResolvedValue({})
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      await wrapper.find('[data-testid="btn-save-sheet"]').trigger('click')
      await flushPromises()

      expect(mockSaveGameSubscriptionInstanceTurnSheet).toHaveBeenCalledWith('gsi-abc-123', 'ts-1', expect.any(Object))
    })

    it('navigates between sheets with prev/next buttons', async () => {
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="btn-prev-sheet"]').attributes('disabled')).toBeDefined()
      expect(wrapper.find('[data-testid="btn-next-sheet"]').attributes('disabled')).toBeUndefined()

      await wrapper.find('[data-testid="btn-next-sheet"]').trigger('click')
      await flushPromises()

      expect(mockGetGameSubscriptionInstanceTurnSheetHTML).toHaveBeenCalledWith('gsi-abc-123', 'ts-2')
      expect(wrapper.find('[data-testid="btn-next-sheet"]').attributes('disabled')).toBeDefined()
      expect(wrapper.find('[data-testid="btn-prev-sheet"]').attributes('disabled')).toBeUndefined()
    })

    it('only shows current turn sheets (latest turn, not completed)', async () => {
      const mixed = [
        { id: 'ts-old', sheet_type: 'adventure_game_location_choice', turn_number: 1, is_completed: true, scanned_data: null },
        { id: 'ts-new-1', sheet_type: 'adventure_game_location_choice', turn_number: 2, is_completed: false, scanned_data: null },
        { id: 'ts-new-2', sheet_type: 'adventure_game_inventory_management', turn_number: 2, is_completed: false, scanned_data: null },
      ]
      mockGetGameSubscriptionInstanceTurnSheets.mockResolvedValue({ turn_sheets: mixed })
      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="ts-step-0"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="ts-step-1"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="ts-step-2"]').exists()).toBe(false)
    })
  })
})
