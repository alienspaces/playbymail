<template>
  <div class="entry-card studio-entry">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon type="pencil" color="blue" class="title-icon" />
      Game Designer Studio
    </h1>
    <p>
      The Game Designer Studio is your workspace for creating, editing, and managing play-by-mail games. Here, designers
      can:
    </p>
    <ul>
      <li>Create and organize games</li>
      <li>Design locations, items, and creatures</li>
      <li>Link locations and define requirements</li>
      <li>Assign items and creatures to locations</li>
      <li>Test and iterate on adventure game designs</li>
    </ul>
    <p><em>Sign in to access the full studio features.</em></p>
    <div class="entry-cta">
      <router-link to="/login" class="cta-button">Sign In to Studio</router-link>
    </div>
  </div>
</template>

<script setup>
import HandDrawnIcon from '../components/HandDrawnIcon.vue';
</script>
