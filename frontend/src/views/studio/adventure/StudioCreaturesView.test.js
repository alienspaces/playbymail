// StudioCreaturesView.test.js
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioCreaturesView from './StudioCreaturesView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockCreatures = [
  { id: 1, name: 'Goblin', description: 'Small and green', created_at: '2024-07-10T12:00:00Z' }
];

vi.mock('../../../api/creatures', () => ({
  fetchCreatures: vi.fn(async () => mockCreatures),
  createCreature: vi.fn(),
  updateCreature: vi.fn(),
  deleteCreature: vi.fn()
}));

describe('StudioCreaturesView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioCreaturesView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioCreaturesView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage creatures.');
  });

  it('renders table headers and data when creatures are loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Description');
    expect(headerTexts).toContain('Created');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Goblin');
    expect(cellTexts).toContain('Small and green');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('creatures', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
}); 