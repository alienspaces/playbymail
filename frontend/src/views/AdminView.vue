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
  margin: 2rem auto;
  padding: 2.5rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.admin-view h1 {
  margin-bottom: 1.5rem;
  font-size: 2.2rem;
  text-align: left;
}
.admin-view p {
  margin-bottom: 1.25rem;
  font-size: 1.1rem;
  text-align: left;
}
.admin-view ul {
  margin: 1.5rem 0 1.5rem 2rem;
  padding-left: 0.5rem;
  text-align: left;
}
.admin-view li {
  margin-bottom: 0.75rem;
  font-size: 1.05rem;
}
</style> 