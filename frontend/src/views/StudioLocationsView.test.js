// StudioLocationsView.test.js
// This test file follows the same pattern as StudioItemsView.test.js and StudioCreaturesView.test.js.
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioLocationsView from './StudioLocationsView.vue';

vi.mock('../api/locations', () => ({
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
    const wrapper = shallowMount(StudioLocationsView, {
      props: { gameId: 'game1' }
    });
    await new Promise(r => setTimeout(r));
    expect(wrapper.text()).toContain('Game Locations');
    expect(wrapper.text()).toContain('Name');
    expect(wrapper.text()).toContain('Description');
    expect(wrapper.text()).toContain('Created');
    expect(wrapper.text()).toContain('Cave');
    expect(wrapper.text()).toContain('Dark cave');
  });

  it('shows loading state', async () => {
    const { useLocationsStore } = await import('../stores/locations');
    const store = useLocationsStore();
    store.loading = true;
    store.error = null;
    const wrapper = shallowMount(StudioLocationsView, {
      props: { gameId: 'game1' }
    });
    expect(wrapper.text()).toContain('Loading...');
  });

  it('shows error state', async () => {
    const { useLocationsStore } = await import('../stores/locations');
    const store = useLocationsStore();
    // Stub loadLocations so it doesn't change state
    vi.spyOn(store, 'loadLocations').mockImplementation(() => {});
    store.loading = false;
    store.error = 'Something went wrong';
    const wrapper = shallowMount(StudioLocationsView, {
      props: { gameId: 'game1' }
    });
    expect(wrapper.text()).toContain('Error: Something went wrong');
  });
}); 