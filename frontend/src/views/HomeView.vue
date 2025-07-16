<template>
  <div class="home-view">
    <h1>Game Designer Studio</h1>
    <p>
      Welcome to the PlayByMail Game Designer Studio! This is your hub for creating, managing, and editing games and all their components.
    </p>
    <ul>
      <li>List, create, edit, and delete <strong>Games</strong></li>
      <li>Manage <strong>Locations</strong> for adventure games</li>
      <li>Manage <strong>Items</strong> and <strong>Creatures</strong></li>
      <li>Link locations and define requirements</li>
      <li>Assign items and creatures to locations</li>
      <li>All features are designed for game designers (not players)</li>
    </ul>
    <router-link to="/studio" class="studio-cta">Enter the Studio</router-link>
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
  margin-bottom: 0.75rem;
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
