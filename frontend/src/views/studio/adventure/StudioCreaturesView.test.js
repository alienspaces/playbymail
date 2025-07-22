// StudioCreaturesView.test.js
// This test file follows the same pattern as StudioLocationsView.test.js and StudioItemsView.test.js.
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioCreaturesView from './StudioCreaturesView.vue';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';

vi.mock('../../../api/creatures', () => ({
  fetchCreatures: vi.fn(async () => [
    { id: 1, name: 'Goblin', description: 'Small and green', created_at: '2024-07-10T12:00:00Z' }
  ]),
  createCreature: vi.fn(),
  updateCreature: vi.fn(),
  deleteCreature: vi.fn()
}));

describe('StudioCreaturesView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  function mountWithRealComponents() {
    return shallowMount(StudioCreaturesView, {
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
    const wrapper = shallowMount(StudioCreaturesView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage creatures.');
  });

  it('renders table headers when creatures are loaded', async () => {
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
    expect(headerTexts).toContain('Created');
    // Check table row data
    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Goblin');
    expect(cellTexts).toContain('Small and green');
  });

  it('shows loading state', async () => {
    const { useCreaturesStore } = await import('../../../stores/creatures');
    const { useGamesStore } = await import('../../../stores/games');
    const gamesStore = useGamesStore();
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game' };
    const store = useCreaturesStore();
    store.loading = true;
    store.error = null;
    const wrapper = mountWithRealComponents();
    // Look for loading indicator
    expect(wrapper.html()).toContain('Loading...');
  });
}); 