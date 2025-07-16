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
  margin: 2rem auto;
  padding: 2.5rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.home-view h1 {
  margin-bottom: 1.5rem;
  font-size: 2.2rem;
  text-align: left;
}
.home-view p {
  margin-bottom: 1.25rem;
  font-size: 1.1rem;
  text-align: left;
}
.home-view ul {
  text-align: left;
  margin: 2rem 0 2rem 2rem;
  padding-left: 0.5rem;
}
.home-view li {
  margin-bottom: 1.5rem;
  font-size: 1.05rem;
}
.studio-cta {
  display: inline-block;
  margin-top: 2rem;
  padding: 0.75rem 2rem;
  background: #2c3e50;
  color: #fff;
  border-radius: 4px;
  text-decoration: none;
  font-weight: bold;
  transition: background 0.2s;
  text-align: center;
}
.studio-cta:hover {
  background: #1a232b;
}
</style>
