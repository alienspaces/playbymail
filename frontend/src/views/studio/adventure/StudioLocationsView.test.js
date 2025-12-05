// StudioLocationsView.test.js
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioLocationsView from './StudioLocationsView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockLocations = [
  { id: 'loc1', name: 'Cave', description: 'Dark cave', created_at: '2024-07-10T12:00:00Z' }
];

vi.mock('../../../api/locations', () => ({
  fetchLocations: vi.fn(async () => mockLocations),
  createLocation: vi.fn(),
  updateLocation: vi.fn(),
  deleteLocation: vi.fn()
}));

describe('StudioLocationsView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioLocationsView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioLocationsView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage locations.');
  });

  it('renders table headers and data when locations are loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Description');
    expect(headerTexts).toContain('Actions');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Cave');
    expect(cellTexts).toContain('Dark cave');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('locations', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
}); 