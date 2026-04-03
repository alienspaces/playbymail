import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioWeaponsView from './StudioWeaponsView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockWeapons = [
  {
    id: 'weapon1',
    name: 'Medium Laser',
    description: 'A reliable direct-fire weapon.',
    damage: 5,
    heat_cost: 3,
    range_band: 'medium',
    mount_size: 'small',
    created_at: '2024-07-10T12:00:00Z'
  }
];

vi.mock('../../../api/mechaWeapons', () => ({
  fetchWeapons: vi.fn(async () => ({ data: mockWeapons, hasMore: false })),
  createWeapon: vi.fn(),
  updateWeapon: vi.fn(),
  deleteWeapon: vi.fn()
}));

describe('StudioWeaponsView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioWeaponsView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioWeaponsView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage weapons.');
  });

  it('renders table headers and data when weapons are loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Damage');
    expect(headerTexts).toContain('Heat');
    expect(headerTexts).toContain('Range');
    expect(headerTexts).toContain('Mount');
    expect(headerTexts).toContain('Actions');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Medium Laser');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('mechaWeapons', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
});
