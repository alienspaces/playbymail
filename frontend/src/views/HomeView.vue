<template>
  <div class="home-view card">
    <div class="hero-section">
      <h1>Welcome to PlayByMail</h1>
      <p>
        <strong>PlayByMail</strong> is a platform for designing, managing, and running play-by-mail games—games played
        entirely through physical mail, not online.
      </p>
    </div>

    <div class="genres-section">
      <h2>Game Types We Support</h2>
      <p>The platform is being built to support the full spectrum of classic and modern play-by-mail game types,
        including:</p>
      <ul class="genres-list">
        <li>Strategy and wargames</li>
        <li>Fantasy empire-building campaigns</li>
        <li>Space 4X and science-fiction strategy</li>
        <li>Diplomacy and political negotiation games</li>
        <li>Hybrid strategy/RPG formats</li>
      </ul>
    </div>

    <div class="divider"></div>

    <div class="who-section">
      <h2>Who is this platform for?</h2>
      <p>This platform is designed for:</p>
      <div class="user-types">
        <div class="user-type">
          <div class="user-icon">
            <HandDrawnIcon type="pencil" color="blue" />
          </div>
          <div class="user-content">
            <h3>Game Designers</h3>
            <p>Create and publish play-by-mail games using our design tools.</p>
          </div>
        </div>
        <div class="user-type">
          <div class="user-icon">
            <HandDrawnIcon type="clipboard" color="blue" />
          </div>
          <div class="user-content">
            <h3>Game Managers</h3>
            <p>Run games, review player turn submissions and outcomes, attend to player queries, provide assistance, and
              review potential issues highlighted by the platform.</p>
          </div>
        </div>
        <div class="user-type">
          <div class="user-icon">
            <HandDrawnIcon type="person" color="blue" />
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
          <HandDrawnIcon type="pencil" color="blue" />
        </div>
        <h2>Game Designer Studio</h2>
        <p>Create, edit, and manage play-by-mail games. Design locations, items, creatures, and game mechanics.</p>
        <router-link to="/studio" class="section-link">Go to Studio</router-link>
      </div>

      <div class="section">
        <div class="section-icon">
          <HandDrawnIcon type="clipboard" color="blue" />
        </div>
        <h2>Game Management</h2>
        <p>Subscribe to or purchase games, configure runtime parameters, review game states, process player turns, and
          provide player help.</p>
        <router-link to="/admin" class="section-link">Go to Game Management</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import HandDrawnIcon from '../components/HandDrawnIcon.vue';
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

.genres-section {
  margin: var(--space-lg) 0;
}

.home-view .genres-section h2 {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-lg);
  color: var(--color-text) !important;
  text-align: left;
}

.genres-section p {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-md);
  text-align: left;
}

.genres-list {
  margin-left: var(--space-lg);
  padding-left: var(--space-sm);
  line-height: 1.8;
}

.genres-list li {
  margin-bottom: var(--space-xs);
  font-size: var(--font-size-md);
}

.who-section {
  margin: var(--space-lg) 0;
}

.home-view .who-section h2 {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-lg);
  color: var(--color-text) !important;
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

.home-view .user-type h3 {
  margin: 0 0 var(--space-sm) 0;
  font-size: var(--font-size-md);
  color: var(--color-text) !important;
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
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-icon :deep(.hand-drawn-icon) {
  font-size: 1.5em;
}

.user-icon :deep(.icon-svg) {
  width: 30px;
  height: 30px;
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

.home-view .sections .section h2 {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-lg);
  color: var(--color-text) !important;
}

.section p {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-md);
  line-height: 1.5;
}

.section-link {
  display: inline-block;
  padding: var(--space-md) var(--space-lg);
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  border-radius: var(--radius-sm);
  text-decoration: none;
  font-weight: var(--font-weight-bold);
  transition: all 0.2s;
  margin-top: auto;
}

.section-link:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}

.section-icon {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  width: 50px;
  height: 50px;
  background: var(--color-bg);
  border: 3px solid var(--color-pen-blue);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.section-icon :deep(.hand-drawn-icon) {
  font-size: 1.8em;
}

.section-icon :deep(.icon-svg) {
  width: 28px;
  height: 28px;
}

@media (max-width: 768px) {
  .sections {
    grid-template-columns: 1fr;
    gap: var(--space-lg);
  }
}
</style>
