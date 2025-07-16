// StudioItemsView.test.js
// This test file follows the same pattern as StudioLocationsView.test.js and StudioCreaturesView.test.js.
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioItemsView from './StudioItemsView.vue';

vi.mock('../../../api/items', () => ({
  fetchItems: vi.fn(async () => [
    { id: 1, name: 'Sword', description: 'Sharp', created_at: '2024-07-10T12:00:00Z' }
  ]),
  createItem: vi.fn(),
  updateItem: vi.fn(),
  deleteItem: vi.fn()
}));

describe('StudioItemsView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioItemsView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage items.');
  });

  it('renders table headers when items are loaded', async () => {
    const { useGamesStore } = await import('../../../stores/games');
    const gamesStore = useGamesStore();
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game' };
    const wrapper = shallowMount(StudioItemsView);
    await new Promise(r => setTimeout(r));
    expect(wrapper.text()).toContain('Game Items');
    expect(wrapper.text()).toContain('Name');
    expect(wrapper.text()).toContain('Description');
    expect(wrapper.text()).toContain('Created');
    expect(wrapper.text()).toContain('Sword');
    expect(wrapper.text()).toContain('Sharp');
  });

  it('shows loading state', async () => {
    const { useItemsStore } = await import('../../../stores/items');
    const { useGamesStore } = await import('../../../stores/games');
    const gamesStore = useGamesStore();
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game' };
    const store = useItemsStore();
    store.loading = true;
    store.error = null;
    const wrapper = shallowMount(StudioItemsView);
    expect(wrapper.text()).toContain('Loading...');
  });
}); 