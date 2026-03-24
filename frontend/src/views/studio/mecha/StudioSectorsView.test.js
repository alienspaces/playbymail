import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioSectorsView from './StudioSectorsView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockSectors = [
  {
    id: 'sector1',
    name: 'Drop Zone',
    description: 'The staging area.',
    terrain_type: 'open',
    elevation: 0,
    is_starting_sector: true,
    created_at: '2024-07-10T12:00:00Z'
  },
  {
    id: 'sector2',
    name: 'Citadel',
    description: 'Fortified command hub.',
    terrain_type: 'urban',
    elevation: 2,
    is_starting_sector: false,
    created_at: '2024-07-10T12:00:00Z'
  }
];

vi.mock('../../../api/mechaSectors', () => ({
  fetchSectors: vi.fn(async () => ({ data: mockSectors, hasMore: false })),
  createSector: vi.fn(),
  updateSector: vi.fn(),
  deleteSector: vi.fn()
}));

describe('StudioSectorsView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioSectorsView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioSectorsView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage sectors.');
  });

  it('renders table headers and data when sectors are loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Description');
    expect(headerTexts).toContain('Starting Sector');
    expect(headerTexts).toContain('Actions');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Drop Zone');
    expect(cellTexts).toContain('The staging area.');
  });

  it('formats is_starting_sector boolean values correctly', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Yes');
    expect(cellTexts).toContain('No');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('mechaSectors', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
});
