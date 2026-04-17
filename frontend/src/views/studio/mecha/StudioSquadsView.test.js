import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioSquadsView from './StudioSquadsView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockSquads = [
  {
    id: 'squad1',
    name: 'Falcon Squad',
    description: 'Fast scouting squad.',
    squad_type: 'opponent',
    created_at: '2024-07-10T12:00:00Z'
  }
];

vi.mock('../../../api/mechaSquads', () => ({
  fetchSquads: vi.fn(async () => ({ data: mockSquads, hasMore: false })),
  createSquad: vi.fn(),
  updateSquad: vi.fn(),
  deleteSquad: vi.fn()
}));

describe('StudioSquadsView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioSquadsView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioSquadsView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage squads.');
  });

  it('renders table headers and data when squads are loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Type');
    expect(headerTexts).toContain('Description');
    expect(headerTexts).toContain('Actions');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts.some(t => t.includes('Falcon Squad'))).toBe(true);
    expect(cellTexts).toContain(`Fast scouting squad.`);
    expect(cellTexts).toContain('Opponent');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('mechaSquads', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
});
