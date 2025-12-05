// Test helpers for Studio resource view components (Items, Locations, Creatures, etc.)
import { vi } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import ResourceTable from '../components/ResourceTable.vue';
import ResourceModalForm from '../components/ResourceModalForm.vue';

/**
 * Creates a mounting helper for Studio resource views
 * @param {Component} Component - The Vue component to mount
 * @returns {Function} A function that returns a mounted component
 */
export function createStudioResourceMountHelper(Component) {
    return function mountWithRealComponents() {
        return shallowMount(Component, {
            global: {
                stubs: {
                    ResourceTable: false,
                    ResourceModalForm: false
                },
                components: { ResourceTable, ResourceModalForm }
            }
        });
    };
}

/**
 * Sets up Pinia stores for testing
 * @param {string} storeName - Name of the store (e.g., 'items', 'locations')
 * @param {Object} storeData - Data to set on the store
 * @returns {Object} The store instance
 */
export async function setupStore(storeName, storeData = {}) {
    const storeModule = await import(`../stores/${storeName}.js`);
    const capitalized = storeName.charAt(0).toUpperCase() + storeName.slice(1);
    const storeFunction = storeModule[`use${capitalized}Store`];
    const store = storeFunction();

    Object.assign(store, storeData);
    return store;
}

/**
 * Sets up games store with a selected game
 * @param {Object} game - Game object to select
 * @returns {Object} The games store instance
 */
export async function setupGamesStore(game = { id: 'game1', name: 'Test Game' }) {
    const { useGamesStore } = await import('../stores/games.js');
    const gamesStore = useGamesStore();
    gamesStore.selectedGame = game;
    return gamesStore;
}

/**
 * Waits for Vue to update (alternative to setTimeout)
 * @returns {Promise} A promise that resolves after Vue updates
 */
export function waitForVueUpdate() {
    return new Promise(resolve => setTimeout(resolve, 0));
}

// Note: setupStudioResourceTests was removed - use beforeEach directly in test files
// This keeps the test context clear and avoids issues with beforeEach scope

/**
 * Creates mock API functions for a resource
 * @param {string} resourceName - Name of the resource (e.g., 'items', 'locations')
 * @param {Array} mockData - Mock data to return from fetch
 * @returns {Object} Mock API functions
 */
export function createMockApi(resourceName, mockData = []) {
    const capitalized = resourceName.charAt(0).toUpperCase() + resourceName.slice(1);
    return {
        [`fetch${capitalized}`]: vi.fn(async () => mockData),
        [`create${capitalized.slice(0, -1)}`]: vi.fn(),
        [`update${capitalized.slice(0, -1)}`]: vi.fn(),
        [`delete${capitalized.slice(0, -1)}`]: vi.fn()
    };
}

/**
 * Helper to find elements in document body (where Teleport renders modals)
 * @param {string} selector - CSS selector
 * @returns {Element|null} The found element or null
 */
export function findInBody(selector) {
    return document.body.querySelector(selector);
}

/**
 * Helper to find all elements in document body (where Teleport renders modals)
 * @param {string} selector - CSS selector
 * @returns {NodeList} The found elements
 */
export function findAllInBody(selector) {
    return document.body.querySelectorAll(selector);
}

/**
 * Sets up document.body cleanup for modal tests
 * Call in beforeEach and afterEach
 */
export function setupModalTestCleanup() {
    return {
        beforeEach: () => {
            document.body.innerHTML = '';
        },
        afterEach: () => {
            document.body.innerHTML = '';
        }
    };
}

