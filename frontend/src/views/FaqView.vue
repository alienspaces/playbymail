<template>
  <div class="faq-view">
    <h1>Frequently Asked Questions</h1>
    <ul>
      <li>
        <strong>What is Play by Mail?</strong>
        <p>Play by Mail is a way to play games through physical mail. Players send in their moves by mail, and receive results by mail—no online play.</p>
      </li>
      <li>
        <strong>What does playbymail.games provide?</strong>
        <p>We provide a comprehensive platform for creating, managing, and playing play-by-mail games. The platform includes automated turn processing, game design tools, and player management systems.</p>
      </li>
      <li>
        <strong>Can I play games online here?</strong>
        <p>No. This platform is for designing and running games that are played entirely by mail. Players join games here but play offline via physical mail.</p>
      </li>
      <li>
        <strong>Who is this site for?</strong>
        <p>This platform serves three distinct user types:</p>
        <div class="user-types">
          <div class="user-type">
            <h3>Game Designers</h3>
            <p>Create and publish play-by-mail games using our design tools.</p>
          </div>
          <div class="user-type">
            <h3>Game Managers</h3>
            <p>Run games, review automated turn processing, attend to player queries, and provide assistance.</p>
          </div>
          <div class="user-type">
            <h3>Players</h3>
            <p>Join games, manage their accounts, and seek help.</p>
          </div>
        </div>
      </li>
      <li>
        <strong>How does turn processing work?</strong>
        <p>The platform automatically processes player turns submitted by mail. Game managers review the results, attend to player queries, and handle any issues flagged by the system for human oversight.</p>
      </li>
      <li>
        <strong>How do I get started as a Game Designer?</strong>
        <p>Sign up, access the Game Designer Studio, and create your first game. The studio provides tools for designing locations, items, creatures, and game mechanics.</p>
      </li>
      <li>
        <strong>How do I get started as a Game Manager?</strong>
        <p>Sign up, access the Game Management area, subscribe to games you want to run, and configure them for automated processing with your oversight.</p>
      </li>
      <li>
        <strong>How do I get started as a Player?</strong>
        <p>Sign up, browse available games, join the ones you're interested in, and start receiving game materials by mail.</p>
      </li>
    </ul>
  </div>
</template>

<script setup>
// No logic needed for static FAQ view
</script>

<style scoped>
.faq-view {
  max-width: 900px;
  margin: var(--space-lg) auto;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.faq-view h1 {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-xl);
  text-align: left;
}
.faq-view ul {
  list-style: none;
  padding: 0;
  margin-left: var(--space-lg);
}
.faq-view li {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-sm);
  text-align: left;
}
.faq-view strong {
  display: block;
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-md);
}
.faq-view p {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-md);
  text-align: left;
}
.faq-view ul ul {
  margin-top: var(--space-sm);
  margin-bottom: var(--space-sm);
}
.faq-view ul ul li {
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-sm);
}
.faq-view ul ul strong {
  display: inline;
  font-size: var(--font-size-sm);
  margin-bottom: 0;
}
.user-types {
  margin-top: var(--space-md);
  margin-bottom: var(--space-md);
}
.user-type {
  margin-bottom: var(--space-lg);
  padding: var(--space-md);
  background: #f8f9fa;
  border-left: 4px solid var(--color-primary);
  border-radius: var(--radius-sm);
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
</style> 