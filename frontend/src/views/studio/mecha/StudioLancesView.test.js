import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioLancesView from './StudioLancesView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockLances = [
  {
    id: 'lance1',
    name: 'Falcon Lance',
    description: 'Fast scouting lance.',
    account_user_id: 'user-abc',
    created_at: '2024-07-10T12:00:00Z'
  }
];

vi.mock('../../../api/mechaLances', () => ({
  fetchLances: vi.fn(async () => ({ data: mockLances, hasMore: false })),
  createLance: vi.fn(),
  updateLance: vi.fn(),
  deleteLance: vi.fn()
}));

describe('StudioLancesView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioLancesView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioLancesView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage lances.');
  });

  it('renders table headers and data when lances are loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Description');
    expect(headerTexts).toContain('Account User ID');
    expect(headerTexts).toContain('Actions');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Falcon Lance');
    expect(cellTexts).toContain('Fast scouting lance.');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('mechaLances', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
});
