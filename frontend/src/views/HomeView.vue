<template>
  <div class="home-view card">
    <div class="hero-section">
      <div class="hero-decoration">
        <svg class="hero-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
        </svg>
      </div>
      <h1>Welcome to PlayByMail</h1>
      <p>
        <strong>PlayByMail</strong> is a platform for designing, managing, and running play-by-mail games—games played entirely through physical mail, not online.
      </p>
    </div>
    
    <div class="who-section">
      <h2>Who is this site for?</h2>
      <p>Our platform is designed for:</p>
      <div class="user-types">
        <div class="user-type">
          <div class="user-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
          </div>
          <div class="user-content">
            <h3>Game Designers</h3>
            <p>Create and publish play-by-mail games using our design tools.</p>
          </div>
        </div>
        <div class="user-type">
          <div class="user-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
            </svg>
          </div>
          <div class="user-content">
            <h3>Game Managers</h3>
            <p>Run games, review player turn submissions and outcomes, attend to player queries, provide assistance, and review potential issues highlighted by the platform.</p>
          </div>
        </div>
        <div class="user-type">
          <div class="user-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
            </svg>
          </div>
          <div class="user-content">
            <h3>Players</h3>
            <p>Join games, manage your account, and seek help.</p>
          </div>
        </div>
      </div>
    </div>
    
    <div class="divider"></div>
    
    <div class="sections">
      <div class="section">
        <div class="section-icon">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
        </div>
        <h2>Game Designer Studio</h2>
        <p>Create, edit, and manage play-by-mail games. Design locations, items, creatures, and game mechanics.</p>
        <router-link to="/studio" class="section-link">Go to Studio</router-link>
      </div>
      
      <div class="section">
        <div class="section-icon">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
          </svg>
        </div>
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
.hero-section {
  position: relative;
  margin-bottom: var(--space-lg);
  text-align: center;
}
.hero-decoration {
  position: absolute;
  top: -50px;
  left: 50%;
  transform: translateX(-50%);
  width: 100px;
  height: 100px;
  background: var(--color-primary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: -1;
}
.hero-icon {
  width: 50px;
  height: 50px;
  color: var(--color-text-light);
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
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.user-content {
  flex: 1;
}
.user-type:last-child {
  margin-bottom: 0;
}
.user-type h3 {
  margin: 0 0 var(--space-sm) 0;
  font-size: var(--font-size-md);
  color: var(--color-primary-light);
  font-weight: var(--font-weight-bold);
}
.user-type p {
  margin: 0;
  font-size: var(--font-size-sm);
  line-height: 1.5;
}
.user-icon {
  width: 30px;
  height: 30px;
  color: var(--color-primary);
  flex-shrink: 0;
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
  position: relative;
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
.section-icon {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  width: 40px;
  height: 40px;
  background: var(--color-primary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
}

.section-icon svg {
  width: 24px;
  height: 24px;
  color: var(--color-text-light);
}
@media (max-width: 768px) {
  .sections {
    grid-template-columns: 1fr;
    gap: var(--space-lg);
  }
}
</style>
