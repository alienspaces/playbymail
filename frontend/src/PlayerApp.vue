<template>
  <div id="player-app">
    <main class="player-content">
      <router-view />
    </main>
    <PlayerSupportFooter />
  </div>
</template>

<script setup>
import PlayerSupportFooter from './components/PlayerSupportFooter.vue'
console.log('[PlayerApp] mounted, current URL:', window.location.href)
</script>

<style>
#player-app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/*
 * Block-level wrapper (no display: flex) so child views are in normal
 * block flow. margin: auto centering on views works correctly here
 * regardless of whether views explicitly set width: 100%.
 */
.player-content {
  flex: 1;
  width: 100%;
}
</style>
