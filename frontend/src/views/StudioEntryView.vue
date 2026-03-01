<template>
  <div class="entry-card studio-entry">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon type="pencil" color="blue" class="title-icon" />
      Game Designer Studio
    </h1>

    <p class="intro">Create the next generation of play-by-mail campaigns.</p>

    <p>
      The Game Designer Studio is an evolving toolkit dedicated exclusively to building PBM campaigns.
    </p>
    <p>
      Game designers will be able to design games for all of the game types the platform supports, from adventure games
      and strategy campaigns to sports simulations and more.
    </p>

    <section class="studio-purpose">
      <h2>What the Studio Is Being Designed For</h2>
      <p>Our goal is to provide creators with tools to:</p>
      <ul>
        <li>Outline branching stories, strategic scenarios, and player roles</li>
        <li>Define locations, maps, units, encounters, and events</li>
        <li>Import custom artwork for maps, tokens, handouts, and other game materials</li>
        <li>Link regions and scenes with travel rules, conditions, and objectives</li>
        <li>Playtest turn flows before publishing your game to players</li>
      </ul>
      <p>As features come online, you’ll be able to build complete PBM campaigns in one place.</p>
    </section>

    <p class="cta-note">
      Sign in to access current tools and track new features as they become available.
    </p>
    <div class="entry-cta">
      <router-link :to="{ path: '/login', query: { redirect: '/studio' } }" class="cta-button" data-testid="studio-sign-in">Sign In to
        Studio</router-link>
    </div>
  </div>
</template>

<script setup>
import HandDrawnIcon from '../components/HandDrawnIcon.vue';
</script>
