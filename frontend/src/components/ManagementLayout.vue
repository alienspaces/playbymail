<!--
  ManagementLayout.vue
  Layout component for the game management interface.
-->
<template>
  <div class="management-layout">
    <!-- Management Header -->
    <div class="management-header">
      <h1>Game Management</h1>
      <div class="header-actions">
        <button class="help-btn" @click="showHelp = true">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17h-2v-2h2v2zm2.07-7.75l-.9.92C13.45 12.9 13 13.5 13 15h-2v-.5c0-1.1.45-2.1 1.17-2.83l1.24-1.26c.37-.36.59-.86.59-1.41 0-1.1-.9-2-2-2s-2 .9-2 2H8c0-2.21 1.79-4 4-4s4 1.79 4 4c0 .88-.36 1.68-.93 2.25z"/>
          </svg>
          Help
        </button>
      </div>
    </div>

    <!-- Main content area -->
    <div class="management-body">
      <section class="management-content">
        <router-view />
      </section>
    </div>

    <!-- Help Panel -->
    <div v-if="showHelp" class="help-panel-overlay" @click.self="showHelp = false">
      <div class="help-panel">
        <div class="help-header">
          <h2>Game Management Help</h2>
          <button class="close-btn" @click="showHelp = false">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </div>
        <div class="help-content">
          <h3>Managing Game Instances</h3>
          <ul>
            <li><strong>Games Dashboard:</strong> View all games you have access to and their current instances</li>
            <li><strong>Create Instances:</strong> Start new game sessions for players to join</li>
            <li><strong>Monitor Progress:</strong> Track turn processing, deadlines, and player activity</li>
            <li><strong>Runtime Controls:</strong> Start, pause, resume, or cancel game instances</li>
            <li><strong>Player Support:</strong> Review submissions and provide assistance</li>
          </ul>
          
          <h3>Game Types</h3>
          <ul>
            <li><strong>Adventure Games:</strong> Exploration and story-driven experiences</li>
            <li><strong>Future Types:</strong> Economic, sports, and other game types coming soon</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const showHelp = ref(false);
</script>

<style scoped>
.management-layout {
  min-height: 100vh;
  background: var(--color-bg-light);
}

.management-header {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border);
  padding: var(--space-md) var(--space-xl);
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.management-header h1 {
  margin: 0;
  font-size: var(--font-size-xl);
  color: var(--color-text);
}

.header-actions {
  display: flex;
  gap: var(--space-md);
}

.help-btn {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  background: var(--color-primary);
  color: var(--color-text-light);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-sm);
  transition: background 0.2s;
}

.help-btn:hover {
  background: var(--color-primary-dark);
}

.nav-icon {
  width: 16px;
  height: 16px;
}

.management-body {
  padding: var(--space-xl);
}

.management-content {
  max-width: 1200px;
  margin: 0 auto;
}

/* Help Panel */
.help-panel-overlay {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 400px;
  background: rgba(0,0,0,0.5);
  z-index: 1000;
  display: flex;
  justify-content: flex-end;
}

.help-panel {
  width: 100%;
  background: var(--color-bg);
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
}

.help-header {
  padding: var(--space-lg);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.help-header h2 {
  margin: 0;
  font-size: var(--font-size-lg);
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: var(--space-sm);
  border-radius: var(--radius-sm);
  transition: background 0.2s;
}

.close-btn:hover {
  background: var(--color-bg-light);
}

.close-btn svg {
  width: 20px;
  height: 20px;
  color: var(--color-text);
}

.help-content {
  padding: var(--space-lg);
  flex: 1;
  overflow-y: auto;
}

.help-content h3 {
  margin: var(--space-lg) 0 var(--space-md) 0;
  font-size: var(--font-size-md);
  color: var(--color-text);
}

.help-content ul {
  margin: 0 0 var(--space-lg) var(--space-lg);
  padding: 0;
}

.help-content li {
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-sm);
  line-height: 1.4;
}

.help-content strong {
  color: var(--color-text);
}
</style> 