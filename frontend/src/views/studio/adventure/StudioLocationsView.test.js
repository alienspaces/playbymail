// StudioLocationsView.test.js
// This test file follows the same pattern as StudioItemsView.test.js and StudioCreaturesView.test.js.
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioLocationsView from './StudioLocationsView.vue';

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
    const wrapper = shallowMount(StudioLocationsView);
    await new Promise(r => setTimeout(r));
    expect(wrapper.text()).toContain('Game Locations');
    expect(wrapper.text()).toContain('Name');
    expect(wrapper.text()).toContain('Description');
    expect(wrapper.text()).toContain('Created');
    expect(wrapper.text()).toContain('Cave');
    expect(wrapper.text()).toContain('Dark cave');
  });

  it('shows loading state', async () => {
    const { useLocationsStore } = await import('../../../stores/locations');
    const { useGamesStore } = await import('../../../stores/games');
    const gamesStore = useGamesStore();
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game' };
    const store = useLocationsStore();
    store.loading = true;
    store.error = null;
    const wrapper = shallowMount(StudioLocationsView);
    expect(wrapper.text()).toContain('Loading...');
  });
}); 