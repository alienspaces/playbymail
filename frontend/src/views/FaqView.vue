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
        <p>We provide tools for game designers and administrators to create, manage, and run play-by-mail games, including turn processing, printing, and mailing support.</p>
      </li>
      <li>
        <strong>Can I play games online here?</strong>
        <p>No. This platform is for designing and running games that are played entirely by mail.</p>
      </li>
      <li>
        <strong>Who is this site for?</strong>
        <p>Game designers and administrators who want to run play-by-mail games. Players use the system only to get help, not to play.</p>
      </li>
      <li>
        <strong>How do I get started?</strong>
        <p>Sign up, create a game in the Game Designer Studio, and follow the documentation to set up your adventure.</p>
      </li>
    </ul>
  </div>
</template>

<script setup>
// No logic needed for static FAQ view
</script>

<style scoped>
.faq-view {
  max-width: 700px;
  margin: 2rem auto;
  padding: 2.5rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.faq-view h1 {
  margin-bottom: 1.5rem;
  font-size: 2.2rem;
  text-align: left;
}
.faq-view ul {
  list-style: none;
  padding: 0;
  margin-left: 2rem;
}
.faq-view li {
  margin-bottom: 2rem;
  font-size: 1.05rem;
  text-align: left;
}
.faq-view strong {
  display: block;
  margin-bottom: 0.5rem;
}
.faq-view p {
  margin-bottom: 1.25rem;
  font-size: 1.1rem;
  text-align: left;
}
</style> 