<template>
  <div>
    <div v-if="!isLoggedIn" class="entry-card management-entry">
      <h1 class="hand-drawn-title">
        <HandDrawnIcon type="clipboard" color="blue" class="title-icon" />
        Game Management
      </h1>
      <p>
        This area is for game managers to subscribe to or purchase games, configure runtime parameters, review automated
        turn processing, attend to player queries, and provide assistance. The platform automatically processes player
        turns, but game managers provide human oversight and support.
      </p>
      <ul>
        <li>Subscribe to or purchase games to run</li>
        <li>Configure game runtime parameters and settings</li>
        <li>Review automated turn processing and outcomes</li>
        <li>Attend to player queries and provide assistance</li>
        <li>Review potential issues highlighted by the platform</li>
        <li>Monitor game states and player progress</li>
      </ul>
      <p><em>Sign in to access game management tools.</em></p>
      <div class="entry-cta">
        <router-link to="/login" class="cta-button">Sign In to Game Management</router-link>
      </div>
    </div>
    <div v-else class="entry-card management-entry">
      <div class="entry-help">
        <h2>Game Management Help</h2>
        <ul>
          <li>Subscribe to or purchase games you want to run for players.</li>
          <li>Configure game runtime parameters like turn frequency and player limits.</li>
          <li>Review automated turn processing results and player submissions.</li>
          <li>Attend to player queries and provide personalized assistance.</li>
          <li>Review potential issues flagged by the platform for human oversight.</li>
          <li>Monitor game states, player progress, and overall game health.</li>
        </ul>
      </div>
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
import { computed } from 'vue';
import HandDrawnIcon from '../components/HandDrawnIcon.vue';

const authStore = useAuthStore();
const isLoggedIn = computed(() => !!authStore.sessionToken);
</script>
