<template>
  <div class="studio-entry">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon type="pencil" color="blue" class="title-icon" />
      Game Designer Studio
    </h1>
    <p>
      The Game Designer Studio is your workspace for creating, editing, and managing play-by-mail games. Here, designers can:
    </p>
    <ul>
      <li>Create and organize games</li>
      <li>Design locations, items, and creatures</li>
      <li>Link locations and define requirements</li>
      <li>Assign items and creatures to locations</li>
      <li>Test and iterate on adventure game designs</li>
    </ul>
    <p><em>Sign in to access the full studio features.</em></p>
    <div class="studio-cta">
      <router-link to="/login" class="cta-button">Sign In to Studio</router-link>
    </div>
  </div>
</template>

<script setup>
import HandDrawnIcon from '../components/HandDrawnIcon.vue';
</script>

<style scoped>
.studio-entry {
  max-width: 900px;
  margin: var(--space-lg) auto;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.studio-entry h1 {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-xl);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  position: relative;
  color: var(--color-text);
}

.hand-drawn-title {
  position: relative;
}

.title-icon {
  font-size: 0.8em;
  margin-right: 0.2em;
}
.studio-entry p {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-md);
}
.studio-entry ul {
  margin-bottom: var(--space-lg);
  margin-left: var(--space-lg);
  padding-left: var(--space-sm);
}
.studio-entry li {
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-sm);
}
.studio-entry em {
  display: block;
  margin-top: var(--space-lg);
  font-size: var(--font-size-sm);
}
.studio-cta {
  text-align: center;
  margin-top: var(--space-lg);
}
.cta-button {
  display: inline-block;
  padding: var(--space-md) var(--space-xl);
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  border-radius: var(--radius-sm);
  text-decoration: none;
  font-weight: var(--font-weight-bold);
  transition: all 0.2s;
}
.cta-button:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}

/* Mobile responsive breakpoint */
@media (max-width: 900px) {
  .studio-entry {
    padding: var(--space-md);
  }
}
</style> 