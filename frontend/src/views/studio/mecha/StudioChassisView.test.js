import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mount, shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import StudioChassisView from './StudioChassisView.vue';
import {
  createStudioResourceMountHelper,
  findAllInBody,
  findInBody,
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

vi.mock('../../../api/mechaGameChassis', () => ({
  fetchMechaGameChassis: vi.fn(async () => ({ data: mockChassis, hasMore: false })),
  createMechaGameChassis: vi.fn(),
  updateMechaGameChassis: vi.fn(),
  deleteMechaGameChassis: vi.fn()
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
    await setupStore('mechaGameChassis', { loading: true, error: null });
    const wrapper = mountWithRealComponents();
    expect(wrapper.html()).toContain('Loading...');
  });
});

describe('StudioChassisView create modal slot fields', () => {
  // These tests fully mount the component so the Teleport'd modal is actually
  // rendered into document.body. They verify the new slot selects behave as
  // designers expect: default to medium-class values, auto-update when the
  // chassis class changes, and stop auto-updating as soon as the user touches
  // any slot value.

  beforeEach(() => {
    setActivePinia(createPinia());
    document.body.innerHTML = '';
  });

  afterEach(() => {
    document.body.innerHTML = '';
  });

  async function openCreateModal() {
    await setupGamesStore();
    const wrapper = mount(StudioChassisView, { attachTo: document.body });
    await waitForVueUpdate();
    const actionButton = wrapper.find('[data-testid="page-header-action"]');
    expect(actionButton.exists()).toBe(true);
    await actionButton.trigger('click');
    await waitForVueUpdate();
    return wrapper;
  }

  function slotSelects() {
    const selects = findAllInBody('.modal select');
    const byLabel = {};
    selects.forEach((sel) => {
      const label = sel.parentElement.querySelector('label')?.textContent?.trim() || '';
      byLabel[label] = sel;
    });
    return {
      small: byLabel['Small Slots *'],
      medium: byLabel['Medium Slots *'],
      large: byLabel['Large Slots *'],
      chassisClass: byLabel['Chassis Class *'],
    };
  }

  it('renders three slot selects defaulting to medium-class values', async () => {
    await openCreateModal();
    const { small, medium, large } = slotSelects();
    expect(small).toBeTruthy();
    expect(medium).toBeTruthy();
    expect(large).toBeTruthy();
    expect(small.value).toBe('2');
    expect(medium.value).toBe('2');
    expect(large.value).toBe('1');
  });

  it('reapplies class defaults when chassis class changes and slots are untouched', async () => {
    const wrapper = await openCreateModal();
    const selects = slotSelects();

    selects.chassisClass.value = 'light';
    selects.chassisClass.dispatchEvent(new Event('change'));
    await waitForVueUpdate();

    const updated = slotSelects();
    expect(updated.small.value).toBe('2');
    expect(updated.medium.value).toBe('1');
    expect(updated.large.value).toBe('0');

    updated.chassisClass.value = 'assault';
    updated.chassisClass.dispatchEvent(new Event('change'));
    await waitForVueUpdate();

    const assaultSelects = slotSelects();
    expect(assaultSelects.small.value).toBe('2');
    expect(assaultSelects.medium.value).toBe('3');
    expect(assaultSelects.large.value).toBe('3');

    wrapper.unmount();
  });

  it('stops auto-updating slots once the designer has edited one', async () => {
    const wrapper = await openCreateModal();
    const selects = slotSelects();

    selects.large.value = '4';
    selects.large.dispatchEvent(new Event('change'));
    await waitForVueUpdate();

    selects.chassisClass.value = 'light';
    selects.chassisClass.dispatchEvent(new Event('change'));
    await waitForVueUpdate();

    const after = slotSelects();
    expect(after.small.value).toBe('2');
    expect(after.medium.value).toBe('2');
    expect(after.large.value).toBe('4');

    wrapper.unmount();
  });

  it('includes slot fields when submitting the form', async () => {
    const apiModule = await import('../../../api/mechaGameChassis');
    apiModule.createMechaGameChassis.mockClear();
    apiModule.createMechaGameChassis.mockResolvedValue({ id: 'new-id' });

    const wrapper = await openCreateModal();
    const selects = slotSelects();
    const nameInput = findInBody('.modal input[maxlength="100"]');
    nameInput.value = 'Test Chassis';
    nameInput.dispatchEvent(new Event('input'));

    selects.chassisClass.value = 'heavy';
    selects.chassisClass.dispatchEvent(new Event('change'));
    await waitForVueUpdate();

    const form = findInBody('.modal form');
    form.dispatchEvent(new Event('submit'));
    await waitForVueUpdate();

    expect(apiModule.createMechaGameChassis).toHaveBeenCalledTimes(1);
    const [, payload] = apiModule.createMechaGameChassis.mock.calls[0];
    expect(payload.chassis_class).toBe('heavy');
    expect(payload.small_slots).toBe(2);
    expect(payload.medium_slots).toBe(2);
    expect(payload.large_slots).toBe(2);

    wrapper.unmount();
  });
});
