import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioEquipmentView from './StudioEquipmentView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockEquipment = [
  {
    id: 'eq1',
    name: 'Double Heat Sink',
    description: 'Additional cooling.',
    effect_kind: 'heat_sink',
    mount_size: 'small',
    magnitude: 4,
    heat_cost: 0,
    created_at: '2024-07-10T12:00:00Z'
  },
  {
    id: 'eq2',
    name: 'Targeting Computer Mk II',
    description: 'Improves targeting solutions.',
    effect_kind: 'targeting_computer',
    mount_size: 'medium',
    magnitude: 10,
    heat_cost: 2,
    created_at: '2024-07-10T12:00:00Z'
  }
];

vi.mock('../../../api/mechaGameEquipment', () => ({
  fetchMechaGameEquipment: vi.fn(async () => ({ data: mockEquipment, hasMore: false })),
  createMechaGameEquipment: vi.fn(),
  updateMechaGameEquipment: vi.fn(),
  deleteMechaGameEquipment: vi.fn()
}));

describe('StudioEquipmentView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioEquipmentView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioEquipmentView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage equipment.');
  });

  it('renders table headers and data when equipment is loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Effect');
    expect(headerTexts).toContain('Magnitude');
    expect(headerTexts).toContain('Heat');
    expect(headerTexts).toContain('Mount');
    expect(headerTexts).toContain('Actions');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Double Heat Sink');
    expect(cellTexts).toContain('Targeting Computer Mk II');
    // Pretty-printed effect kind
    expect(cellTexts).toContain('Heat Sink');
    expect(cellTexts).toContain('Targeting Computer');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('mechaGameEquipment', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
});
