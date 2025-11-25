<template>
  <div class="entry-card studio-entry">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon type="pencil" color="blue" class="title-icon" />
      Game Designer Studio
    </h1>

    <p class="intro">Create the next generation of play-by-mail campaigns.</p>

    <p>
      The Game Designer Studio is an evolving toolkit dedicated exclusively to building PBM campaigns for the all of the
      game types the platform supports.
    </p>

    <section class="genre-callout">
      <h2>PBM Genres We Aim to Support</h2>
      <p>The studio is being built to support the full spectrum of classic and modern PBM game types, including:</p>
      <ul>
        <li>Strategy and wargames</li>
        <li>Fantasy empire-building campaigns</li>
        <li>Space 4X and science-fiction strategy</li>
        <li>Diplomacy and political negotiation games</li>
        <li>Hybrid strategy/RPG formats</li>
      </ul>
      <p>These genres shape the roadmap as we expand the platform’s capabilities.</p>
    </section>

    <section class="studio-purpose">
      <h2>What the Studio Is Being Designed For</h2>
      <p>Our goal is to provide creators with tools to:</p>
      <ul>
        <li>Outline branching stories, strategic scenarios, and player roles</li>
        <li>Define locations, maps, units, encounters, and events</li>
        <li>Import custom artwork for maps, tokens, handouts, and other physical mailers</li>
        <li>Link regions and scenes with travel rules, conditions, and objectives</li>
        <li>Playtest turn flows before publishing your game to players</li>
      </ul>
      <p>As features come online, you’ll be able to build complete PBM campaigns in one place.</p>
    </section>

    <p class="cta-note">
      Sign in to access current tools and track new features as they become available.
    </p>
    <div class="entry-cta">
      <router-link to="/login" class="cta-button">Sign In to Studio</router-link>
    </div>
  </div>
</template>

<script setup>
import HandDrawnIcon from '../components/HandDrawnIcon.vue';
</script>
