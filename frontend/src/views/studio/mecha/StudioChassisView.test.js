import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioChassisView from './StudioChassisView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockChassis = [
  {
    id: 'chassis1',
    name: 'Locust',
    description: 'A fast light mech.',
    chassis_class: 'light',
    armor_points: 56,
    structure_points: 24,
    heat_capacity: 16,
    speed: 8,
    created_at: '2024-07-10T12:00:00Z'
  }
];

vi.mock('../../../api/mechaChassis', () => ({
  fetchChassis: vi.fn(async () => ({ data: mockChassis, hasMore: false })),
  createChassis: vi.fn(),
  updateChassis: vi.fn(),
  deleteChassis: vi.fn()
}));

describe('StudioChassisView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioChassisView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioChassisView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage chassis.');
  });

  it('renders table headers and data when chassis are loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Class');
    expect(headerTexts).toContain('Armor');
    expect(headerTexts).toContain('Structure');
    expect(headerTexts).toContain('Heat Cap.');
    expect(headerTexts).toContain('Speed');
    expect(headerTexts).toContain('Actions');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Locust');
    expect(cellTexts).toContain('light');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('mechaChassis', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
});
