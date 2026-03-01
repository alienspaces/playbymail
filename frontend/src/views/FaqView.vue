<template>
  <div class="entry-card faq-view">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon type="question" color="blue" class="title-icon" />
      Frequently Asked Questions
    </h1>
    <ul class="faq-list">
      <li>
        <strong>What is Play by Mail?</strong>
        <p>Play by Mail is a style of turn-based gaming where players submit their moves and receive results by email
          or physical mail, rather than playing in real time online.</p>
      </li>
      <li>
        <strong>What does playbymail.games provide?</strong>
        <p>We provide a comprehensive platform for creating, managing, and playing turn-based games via email or
          physical mail. The platform includes automated turn processing, game design tools, and player management
          systems.</p>
      </li>
      <li>
        <strong>Is this a real-time online game?</strong>
        <p>No. This is not a real-time game platform. Games are turn-based and driven by the game's turn schedule—players
          submit moves and receive results by email or physical mail each turn, not in real time.</p>
      </li>
      <li>
        <strong>How does turn processing work?</strong>
        <p>Players submit their turn sheets online via email link or by returning printed forms through the mail. The
          platform automatically processes submissions, and game managers review results, attend to player queries, and
          handle any issues flagged for human oversight.</p>
      </li>
      <li>
        <strong>How do I get started as a Game Designer?</strong>
        <p>Sign up, access the Game Designer Studio, and create your first game. The studio provides tools for designing
          locations, items, creatures, and game mechanics.</p>
      </li>
      <li>
        <strong>How do I get started as a Game Manager?</strong>
        <p>Sign up, access the Game Management area, subscribe to games you want to run, and configure them for
          automated processing with your oversight.</p>
      </li>
      <li>
        <strong>How do I get started as a Player?</strong>
        <p>Sign up, browse available games, join the ones you're interested in, and start receiving game materials by
          email or physical mail.</p>
      </li>
    </ul>
  </div>
</template>

<script setup>
import HandDrawnIcon from '../components/HandDrawnIcon.vue';
</script>

<style scoped>
.faq-list {
  list-style: none;
  padding: 0;
  margin-left: 0;
}

.faq-list>li {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-sm);
  text-align: left;
}

.faq-list strong {
  display: block;
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-md);
}

.faq-list p {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-md);
  text-align: left;
}
</style>