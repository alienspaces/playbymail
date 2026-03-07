<template>
  <div id="player-app">
    <router-view />
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
</style>
