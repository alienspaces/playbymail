// StudioItemsView.test.js
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioItemsView from './StudioItemsView.vue';
import {
  createStudioResourceMountHelper,
  setupGamesStore,
  setupStore,
  waitForVueUpdate
} from '../../../test-utils/studio-resource-helpers';

const mockItems = [
  { id: 1, name: 'Sword', description: 'Sharp', is_starting_item: true, created_at: '2024-07-10T12:00:00Z' },
  { id: 2, name: 'Shield', description: 'Protective', is_starting_item: false, created_at: '2024-07-10T12:00:00Z' }
];

vi.mock('../../../api/items', () => ({
  fetchItems: vi.fn(async () => mockItems),
  createItem: vi.fn(),
  updateItem: vi.fn(),
  deleteItem: vi.fn()
}));

describe('StudioItemsView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  const mountWithRealComponents = createStudioResourceMountHelper(StudioItemsView);

  it('shows prompt if no game is selected', () => {
    const wrapper = shallowMount(StudioItemsView, {
      props: { gameId: null }
    });
    expect(wrapper.text()).toContain('Select a game to manage items.');
  });

  it('renders table headers and data when items are loaded', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await waitForVueUpdate();

    const ths = wrapper.findAll('th');
    const headerTexts = ths.map(th => th.text());
    expect(headerTexts).toContain('Name');
    expect(headerTexts).toContain('Description');
    expect(headerTexts).toContain('Starting Item');
    expect(headerTexts).toContain('Actions');

    const tds = wrapper.findAll('td');
    const cellTexts = tds.map(td => td.text());
    expect(cellTexts).toContain('Sword');
    expect(cellTexts).toContain('Sharp');
    expect(cellTexts).toContain('Yes');
    expect(cellTexts).toContain('No');
  });

  it('shows loading state', async () => {
    await setupGamesStore();
    await setupStore('items', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });

  it('includes is_starting_item field in form', async () => {
    await setupGamesStore();
    const wrapper = mountWithRealComponents();
    await wrapper.vm.openCreate();
    await wrapper.vm.$nextTick();

    const fields = wrapper.vm.fields;
    const startingItemField = fields.find(f => f.key === 'is_starting_item');
    expect(startingItemField).toBeDefined();
    expect(startingItemField.label).toBe('Starting Item');
    expect(startingItemField.type).toBe('checkbox');
  });

  it('formats is_starting_item boolean values correctly', async () => {
    await setupGamesStore();
    await setupStore('items', {
      items: [
        { id: 1, name: 'Sword', description: 'Sharp', is_starting_item: true },
        { id: 2, name: 'Shield', description: 'Protective', is_starting_item: false }
      ]
    });
    const wrapper = mountWithRealComponents();
    await wrapper.vm.$nextTick();

    const formatted = wrapper.vm.formattedItems;
    expect(formatted[0].is_starting_item).toBe('Yes');
    expect(formatted[1].is_starting_item).toBe('No');
  });
}); 