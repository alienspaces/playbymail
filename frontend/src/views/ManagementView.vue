<template>
  <div class="admin-view">
    <h1>Game Management</h1>
    <p>
      This section provides tools for managing running games with automated turn processing and human oversight. The platform automatically handles routine turn processing, while game managers provide human intelligence for complex situations and player support.
    </p>
    <ul>
      <li>Subscribe to and configure games for automated processing</li>
      <li>Review automated turn processing results and player submissions</li>
      <li>Attend to player queries and provide personalized assistance</li>
      <li>Review potential issues flagged by the platform for human oversight</li>
      <li>Monitor game states, player progress, and overall game health</li>
      <li>Manage print/mail operations and player communications</li>
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
  margin: var(--space-md) auto var(--space-lg);
  padding: var(--space-lg);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.admin-view h1 {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-xl);
  text-align: left;
}
.admin-view p {
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-md);
  text-align: left;
}
.admin-view ul {
  margin: var(--space-md) 0 var(--space-md) var(--space-lg);
  padding-left: var(--space-sm);
  text-align: left;
}
.admin-view li {
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-sm);
}
</style> 