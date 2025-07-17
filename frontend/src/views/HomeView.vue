<template>
  <div class="home-view">
    <h1>Welcome to PlayByMail</h1>
    <p>
      <strong>PlayByMail</strong> is a platform for designing, managing, and running play-by-mail games—games played entirely through physical mail, not online.
    </p>
    <ul>
      <li><strong>What is Play by Mail?</strong><br>
        Play by mail is a way to play games through the postal service. Players send in their moves by mail and receive results by mail—no online play.
      </li>
      <li><strong>What does playbymail.games provide?</strong><br>
        Tools for game designers and administrators to create, manage, and run play-by-mail games, including turn processing, printing, and mailing support.
      </li>
      <li><strong>Who is this site for?</strong><br>
        Game designers and administrators who want to run play-by-mail games. Players use the system only to get help, not to play.
      </li>
    </ul>
  </div>
</template>

<script setup>
// No logic needed for static home view
</script>

<style scoped>
.home-view {
  max-width: 700px;
  margin: var(--space-lg) auto;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
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
.home-view ul {
  text-align: left;
  margin: var(--space-lg) 0 var(--space-lg) var(--space-lg);
  padding-left: var(--space-sm);
}
.home-view li {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-sm);
}
.studio-cta {
  display: inline-block;
  margin-top: var(--space-lg);
  padding: var(--space-md) var(--space-lg);
  background: #2c3e50;
  color: var(--color-text-light);
  border-radius: var(--radius-sm);
  text-decoration: none;
  font-weight: var(--font-weight-bold);
  transition: background 0.2s;
  text-align: center;
}
.studio-cta:hover {
  background: #1a232b;
}
</style>
