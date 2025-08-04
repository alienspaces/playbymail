// StudioLocationsView.test.js
// This test file follows the same pattern as StudioItemsView.test.js and StudioCreaturesView.test.js.
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioLocationsView from './StudioLocationsView.vue';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';

vi.mock('../../../api/locations', () => ({
  fetchLocations: vi.fn(async () => [
    { id: 'loc1', name: 'Cave', description: 'Dark cave', created_at: '2024-07-10T12:00:00Z' }
  ]),
  createLocation: vi.fn(),
  updateLocation: vi.fn(),
  deleteLocation: vi.fn()
}));

describe('StudioLocationsView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  function mountWithRealComponents() {
    return shallowMount(StudioLocationsView, {
      global: {
        stubs: {
          ResourceTable: false,
          ResourceModalForm: false
        },
        components: { ResourceTable, ResourceModalForm }
      }
    });
  }

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioLocationsView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage locations.');
  });

  it('renders table headers and data when locations are loaded', async () => {
    const { useGamesStore } = await import('../../../stores/games');
    const gamesStore = useGamesStore();
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game' };
    const wrapper = mountWithRealComponents();
    await new Promise(r => setTimeout(r));
    // Check table headers
    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Description');
    expect(headerTexts).toContain('Actions');
    // Check table row data
    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Cave');
    expect(cellTexts).toContain('Dark cave');
  });

  it('shows loading state', async () => {
    const { useLocationsStore } = await import('../../../stores/locations');
    const { useGamesStore } = await import('../../../stores/games');
    const gamesStore = useGamesStore();
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game' };
    const store = useLocationsStore();
    store.loading = true;
    store.error = null;
    const wrapper = mountWithRealComponents();
    // Look for loading indicator
    expect(wrapper.html()).toContain('Loading...');
  });
}); 