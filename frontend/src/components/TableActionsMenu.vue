<template>
  <div class="table-actions-menu" ref="menuContainer">
    <button 
      class="menu-trigger" 
      @click="toggleMenu"
      :aria-label="'Actions for row'"
      type="button"
    >
      <HandDrawnIcon type="more-vertical" color="black" />
    </button>
    <div v-if="isOpen" class="menu-dropdown">
      <button
        v-for="action in actions"
        :key="action.key"
        @click="handleAction(action)"
        class="menu-item"
        :class="{ 'menu-item-danger': action.danger }"
      >
        {{ action.label }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import HandDrawnIcon from './HandDrawnIcon.vue';

defineProps({
  actions: {
    type: Array,
    required: true,
    validator: (actions) => {
      return actions.every(action => 
        action.key && action.label && typeof action.handler === 'function'
      );
    }
  }
});

const isOpen = ref(false);
const menuContainer = ref(null);

function toggleMenu() {
  isOpen.value = !isOpen.value;
}

function handleAction(action) {
  action.handler();
  isOpen.value = false;
}

function handleClickOutside(event) {
  if (menuContainer.value && !menuContainer.value.contains(event.target)) {
    isOpen.value = false;
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
.table-actions-menu {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.menu-trigger {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text);
  transition: opacity 0.2s;
  width: 100%;
  height: 100%;
}

.menu-trigger:hover {
  opacity: 0.7;
}

.menu-dropdown {
  position: absolute;
  right: 0;
  top: 100%;
  margin-top: 4px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  min-width: 120px;
  z-index: 1000;
  overflow: hidden;
}

.menu-item {
  display: block;
  width: 100%;
  padding: 8px 12px;
  text-align: left;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-button);
  transition: all 0.2s;
  font-size: 14px;
  font-weight: var(--font-weight-semibold);
}

.menu-item:hover {
  background-color: var(--color-button);
  color: var(--color-text-light);
}

.menu-item-danger {
  color: var(--color-danger);
}

.menu-item-danger:hover {
  background-color: var(--color-danger);
  color: var(--color-text-light);
}
</style>

