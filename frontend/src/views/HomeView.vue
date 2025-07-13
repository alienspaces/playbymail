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
  max-width: 600px;
  margin: 2rem auto;
  padding: 2rem;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  text-align: center;
}
.home-view h1 {
  margin-bottom: 1rem;
}
.home-view ul {
  text-align: left;
  margin: 2rem auto 2rem 2rem;
  padding-left: 1.5rem;
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
}
.studio-cta:hover {
  background: #1a232b;
}
</style>
