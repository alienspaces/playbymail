<template>
  <div class="admin-view">
    <h1>Game Management & Administration</h1>
    <p>
      This section will provide tools for managing running games, processing mailed-in turns, printing and mailing results, and tracking player progress.
    </p>
    <ul>
      <li>Manage running games</li>
      <li>Process and scan mailed-in turns</li>
      <li>Print and mail results to players</li>
      <li>Track game state and player progress</li>
      <li>Administer player accounts and game data</li>
    </ul>
    <p><em>Coming soon!</em></p>
  </div>
</template>

<script setup>
// No logic needed for static admin view
</script>

<style scoped>
.admin-view {
  max-width: 700px;
  margin: var(--space-lg) auto;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.admin-view h1 {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-xl);
  text-align: left;
}
.admin-view p {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-md);
  text-align: left;
}
.admin-view ul {
  margin: var(--space-lg) 0 var(--space-lg) var(--space-lg);
  padding-left: var(--space-sm);
  text-align: left;
}
.admin-view li {
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-sm);
}
</style> 