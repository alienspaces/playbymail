<template>
  <div class="faq-view">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon type="question" color="blue" class="title-icon" />
      Frequently Asked Questions
    </h1>
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
import HandDrawnIcon from '../components/HandDrawnIcon.vue';
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
</style> 