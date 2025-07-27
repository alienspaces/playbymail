<template>
  <div class="home-view card">
    <h1>Welcome to PlayByMail</h1>
    <p>
      <strong>PlayByMail</strong> is a platform for designing, managing, and running play-by-mail games—games played entirely through physical mail, not online.
    </p>
    
    <div class="who-section">
      <h2>Who is this site for?</h2>
      <p>Our platform is designed for:</p>
      <div class="user-types">
        <div class="user-type">
          <h3>Game Designers</h3>
          <p>Create and publish play-by-mail games using our design tools.</p>
        </div>
        <div class="user-type">
          <h3>Game Managers</h3>
          <p>Run games, review player turn submissions and outcomes, attend to player queries, provide assistance, and review potential issues highlighted by the platform.</p>
        </div>
        <div class="user-type">
          <h3>Players</h3>
          <p>Join games, manage your account, and seek help.</p>
        </div>
      </div>
    </div>
    
    <div class="divider"></div>
    
    <div class="sections">
      <div class="section">
        <h2>Game Designer Studio</h2>
        <p>Create, edit, and manage play-by-mail games. Design locations, items, creatures, and game mechanics.</p>
        <router-link to="/studio" class="section-link">Go to Studio</router-link>
      </div>
      
      <div class="section">
        <h2>Game Management</h2>
        <p>Subscribe to or purchase games, configure runtime parameters, review game states, process player turns, and provide player help.</p>
        <router-link to="/admin" class="section-link">Go to Game Management</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
// No logic needed for static home view
</script>

<style scoped>
.home-view {
  max-width: 900px;
  margin: var(--space-lg) auto;
  padding: var(--space-xl);
}
.home-view h1 {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-xl);
  text-align: left;
}
.home-view p {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-md);
  text-align: left;
}
.who-section {
  margin: var(--space-lg) 0;
}
.who-section h2 {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-lg);
  color: var(--color-primary);
  text-align: left;
}
.user-types {
  margin-top: var(--space-md);
  margin-bottom: var(--space-md);
}
.user-type {
  margin-bottom: var(--space-lg);
  padding: var(--space-md);
  background: var(--color-bg-alt);
  border-left: 4px solid var(--color-primary);
  border-radius: var(--radius-sm);
}
.user-type:last-child {
  margin-bottom: 0;
}
.user-type h3 {
  margin: 0 0 var(--space-sm) 0;
  font-size: var(--font-size-md);
  color: var(--color-primary);
  font-weight: var(--font-weight-bold);
}
.user-type p {
  margin: 0;
  font-size: var(--font-size-sm);
  line-height: 1.5;
}
.divider {
  height: 1px;
  background: var(--color-border);
  margin: var(--space-md) 0;
}
.sections {
  margin: var(--space-lg) 0;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-xl);
}
.section {
  padding: var(--space-lg);
  background: var(--color-bg-alt);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 200px;
}
.section h2 {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-lg);
  color: var(--color-primary);
}
.section p {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-md);
  line-height: 1.5;
}
.section-link {
  display: inline-block;
  padding: var(--space-md) var(--space-lg);
  background: var(--color-primary);
  color: var(--color-text-light);
  border-radius: var(--radius-sm);
  text-decoration: none;
  font-weight: var(--font-weight-bold);
  transition: background 0.2s;
  margin-top: auto;
}
.section-link:hover {
  background: var(--color-primary-dark);
}
@media (max-width: 768px) {
  .sections {
    grid-template-columns: 1fr;
    gap: var(--space-lg);
  }
}
</style>
