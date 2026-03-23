import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { ref } from 'vue';
import StudioSectorLinksView from './StudioSectorLinksView.vue';
import { findInBody, setupModalTestCleanup } from '../../../test-utils/studio-resource-helpers';

// Mock all three stores the view depends on.
vi.mock('../../../stores/mechaSectors', () => ({
  useMechaSectorsStore: vi.fn(() => ({
    sectors: [],
    loading: false,
    error: null,
    fetchSectors: vi.fn()
  }))
}));

vi.mock('../../../stores/mechaSectorLinks', () => ({
  useMechaSectorLinksStore: vi.fn(() => ({
    sectorLinks: [],
    loading: false,
    error: null,
    pageNumber: 1,
    hasMore: false,
    fetchSectorLinks: vi.fn(),
    createSectorLink: vi.fn(),
    updateSectorLink: vi.fn(),
    deleteSectorLink: vi.fn()
  }))
}));

vi.mock('../../../stores/games', () => ({
  useGamesStore: vi.fn(() => ({
    selectedGame: ref(null)
  }))
}));

describe('StudioSectorLinksView', () => {
  const modalCleanup = setupModalTestCleanup();

  const setupStoreMocks = async (selectedGame = null, sectors = [], sectorLinks = []) => {
    const { useGamesStore } = await import('../../../stores/games');
    const { useMechaSectorsStore } = await import('../../../stores/mechaSectors');
    const { useMechaSectorLinksStore } = await import('../../../stores/mechaSectorLinks');

    useGamesStore.mockReturnValue({
      selectedGame: ref(selectedGame)
    });
    useMechaSectorsStore.mockReturnValue({
      sectors,
      loading: false,
      error: null,
      fetchSectors: vi.fn()
    });
    useMechaSectorLinksStore.mockReturnValue({
      sectorLinks,
      loading: false,
      error: null,
      pageNumber: 1,
      hasMore: false,
      fetchSectorLinks: vi.fn(),
      createSectorLink: vi.fn(),
      updateSectorLink: vi.fn(),
      deleteSectorLink: vi.fn()
    });
  };

  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
    modalCleanup.beforeEach();
  });

  afterEach(() => {
    modalCleanup.afterEach();
  });

  it('shows prompt if no game is selected', () => {
    const wrapper = mount(StudioSectorLinksView);
    expect(wrapper.text()).toContain('Select a game to manage sector links.');
    expect(wrapper.find('.game-table-section').exists()).toBe(false);
  });

  it('renders sector links table when a game is selected', async () => {
    await setupStoreMocks({ id: 'game1', name: 'Test Mech Game' });
    const wrapper = mount(StudioSectorLinksView);

    expect(wrapper.find('.game-table-section').exists()).toBe(true);
    expect(wrapper.find('h2').text()).toBe('Sector Links');
  });

  it('resolves sector names from the sectors store', async () => {
    const sectors = [
      { id: 'sec1', name: 'Drop Zone' },
      { id: 'sec2', name: 'Citadel' }
    ];
    const sectorLinks = [
      {
        id: 'link1',
        from_mecha_sector_id: 'sec1',
        to_mecha_sector_id: 'sec2',
        cover_modifier: 2
      }
    ];
    await setupStoreMocks({ id: 'game1', name: 'Test Mech Game' }, sectors, sectorLinks);

    const wrapper = mount(StudioSectorLinksView);
    await wrapper.vm.$nextTick();

    const resourceTable = wrapper.findComponent({ name: 'ResourceTable' });
    expect(resourceTable.exists()).toBe(true);

    const rows = resourceTable.props('rows');
    expect(rows).toHaveLength(1);
    expect(rows[0].from_sector_name).toBe('Drop Zone');
    expect(rows[0].to_sector_name).toBe('Citadel');
  });

  it('opens create modal when create button is clicked', async () => {
    await setupStoreMocks({ id: 'game1', name: 'Test Mech Game' });
    const wrapper = mount(StudioSectorLinksView);

    const createButton = wrapper.find('button');
    await createButton.trigger('click');

    expect(wrapper.vm.showModal).toBe(true);
    expect(wrapper.vm.modalMode).toBe('create');
  });

  it('renders create modal with correct title', async () => {
    await setupStoreMocks({ id: 'game1', name: 'Test Mech Game' });
    const wrapper = mount(StudioSectorLinksView);

    wrapper.vm.showModal = true;
    wrapper.vm.modalMode = 'create';
    await wrapper.vm.$nextTick();

    expect(findInBody('.modal h2').textContent).toBe('Create Sector Link');
  });
});
