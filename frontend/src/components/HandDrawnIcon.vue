<template>
  <span class="hand-drawn-icon" :class="iconClass">
    <svg 
      v-if="type === 'checkbox-checked'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn checkbox with check - wobbly paths -->
      <path d="M4.5 4.5h15v15H4.5z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.5, 0.5;"/>
      <path d="M8 12.5l2.5 2.5 5-5" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
    </svg>
    <svg 
      v-else-if="type === 'checkbox-unchecked'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn checkbox empty - wobbly path -->
      <path d="M4.5 4.5h15v15H4.5z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.5, 0.5;"/>
    </svg>
    <svg 
      v-else-if="type === 'star'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn star - wobbly path -->
      <path d="M12 2.5l2.2 6.8h7.3l-5.8 4.3 2.3 7.2L12 15.8l-5.8 5.2 2.3-7.2L2.7 9.3h7.3z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.4, 0.4;"/>
    </svg>
    <svg 
      v-else-if="type === 'circle'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn circle - wobbly -->
      <circle cx="12" cy="12" r="9.5" 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              vector-effect="non-scaling-stroke"
              style="stroke-dasharray: 0.4, 0.4;"/>
    </svg>
    <svg 
      v-else-if="type === 'arrow'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn arrow - wobbly -->
      <path d="M5.5 12h13M12 5.5l6.5 6.5-6.5 6.5" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
    </svg>
    <svg 
      v-else-if="type === 'underline'" 
      viewBox="0 0 24 4" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg underline-svg"
    >
      <!-- Hand-drawn underline - wobbly -->
      <path d="M2 2.5c2 0 4 0 6 0s4 0 6 0 4 0 6 0" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.2, 0.2;"/>
    </svg>
    <svg 
      v-else-if="type === 'pencil'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn pencil - perfect for game design -->
      <path d="M17 3l4 4-9 9-5-5 9-9z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
      <path d="M3 21h4l9-9-4-4-9 9z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
    </svg>
    <svg 
      v-else-if="type === 'compass'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn drafting compass -->
      <circle cx="12" cy="12" r="8" 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              vector-effect="non-scaling-stroke"
              style="stroke-dasharray: 0.4, 0.4;"/>
      <path d="M12 4v8l6 6" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
      <circle cx="12" cy="12" r="1.5" fill="currentColor"/>
    </svg>
    <svg 
      v-else-if="type === 'lightbulb'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn lightbulb - ideas/creativity -->
      <path d="M9 21h6M12 3a5 5 0 0 1 5 5c0 2-1 3-1 3H8s-1-1-1-3a5 5 0 0 1 5-5z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
      <path d="M9 18v-3M15 18v-3" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.2, 0.2;"/>
    </svg>
    <svg 
      v-else-if="type === 'document'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn document with pen -->
      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
      <path d="M14 2v6h6M10 13h4M10 17h4" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.2, 0.2;"/>
    </svg>
    <svg 
      v-else-if="type === 'dice'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn dice - game element -->
      <rect x="4" y="4" width="16" height="16" rx="2" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.4, 0.4;"/>
      <circle cx="8" cy="8" r="1" fill="currentColor"/>
      <circle cx="12" cy="12" r="1" fill="currentColor"/>
      <circle cx="16" cy="16" r="1" fill="currentColor"/>
    </svg>
    <svg 
      v-else-if="type === 'clipboard'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn clipboard - perfect for game management/review -->
      <path d="M9 5H7a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2h-2" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
      <path d="M9 5a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v0a2 2 0 0 1-2 2h-2a2 2 0 0 1-2-2z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
      <path d="M9 12h6M9 16h6" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.2, 0.2;"/>
    </svg>
    <svg 
      v-else-if="type === 'shield'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn shield - protection/authority -->
      <path d="M12 2L4 5v6c0 5.55 3.84 10.74 8 12 4.16-1.26 8-6.45 8-12V5l-8-3z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.4, 0.4;"/>
    </svg>
    <svg 
      v-else-if="type === 'eye'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn eye - reviewing/overseeing -->
      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
      <circle cx="12" cy="12" r="3" 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              vector-effect="non-scaling-stroke"
              style="stroke-dasharray: 0.2, 0.2;"/>
    </svg>
    <svg 
      v-else-if="type === 'person'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn person/user - perfect for players -->
      <circle cx="12" cy="8" r="4" 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              vector-effect="non-scaling-stroke"
              style="stroke-dasharray: 0.3, 0.3;"/>
      <path d="M6 21v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
    </svg>
    <svg 
      v-else-if="type === 'game-piece'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn game piece/pawn -->
      <circle cx="12" cy="18" r="3" 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              vector-effect="non-scaling-stroke"
              style="stroke-dasharray: 0.3, 0.3;"/>
      <path d="M12 18V8M12 8c-2 0-3-1.5-3-3.5S10 2 12 2s3 1.5 3 3.5S14 8 12 8z" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
    </svg>
    <svg 
      v-else-if="type === 'trophy'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn trophy - achievement/participation -->
      <path d="M6 9H4a2 2 0 0 0-2 2v2a2 2 0 0 0 2 2h2M18 9h2a2 2 0 0 1 2 2v2a2 2 0 0 1-2 2h-2M6 9V7a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v2M6 9v6a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V9" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
      <path d="M12 15v4M9 19h6" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.2, 0.2;"/>
    </svg>
    <svg 
      v-else-if="type === 'question'" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      stroke-width="2.5"
      class="icon-svg"
    >
      <!-- Hand-drawn question mark - perfect for FAQ -->
      <circle cx="12" cy="12" r="10" 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              vector-effect="non-scaling-stroke"
              style="stroke-dasharray: 0.4, 0.4;"/>
      <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3M12 17h.01" 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            vector-effect="non-scaling-stroke"
            style="stroke-dasharray: 0.3, 0.3;"/>
    </svg>
    <svg 
      v-else-if="type === 'more-vertical'" 
      viewBox="0 0 24 24" 
      fill="currentColor" 
      class="icon-svg"
    >
      <!-- Hand-drawn vertical ellipsis - three dots stacked vertically -->
      <circle cx="12" cy="5" r="2" />
      <circle cx="12" cy="12" r="2" />
      <circle cx="12" cy="19" r="2" />
    </svg>
  </span>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  type: {
    type: String,
    required: true,
    validator: (value) => ['checkbox-checked', 'checkbox-unchecked', 'star', 'circle', 'arrow', 'underline', 'pencil', 'compass', 'lightbulb', 'document', 'dice', 'clipboard', 'shield', 'eye', 'person', 'game-piece', 'trophy', 'question', 'more-vertical'].includes(value)
  },
  color: {
    type: String,
    default: 'blue' // 'blue', 'black', 'red'
  }
});

const iconClass = computed(() => {
  return `icon-${props.color}`;
});
</script>

<style scoped>
.hand-drawn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  vertical-align: middle;
}

.icon-svg {
  width: 1em;
  height: 1em;
  /* Hand-drawn effect */
  filter: drop-shadow(0.5px 0.5px 0px currentColor);
  /* Slight rotation for hand-drawn feel - varies per icon */
  transform: rotate(-0.5deg);
}

.underline-svg {
  width: 100%;
  height: 0.3em;
  transform: rotate(0deg);
}

/* Color variants */
.icon-blue {
  color: var(--color-pen-blue);
}

.icon-black {
  color: var(--color-pen-black);
}

.icon-red {
  color: var(--color-pen-red);
}

/* Hand-drawn wobbly effect */
@keyframes hand-drawn-wobble {
  0%, 100% { transform: rotate(-0.5deg) translate(0, 0); }
  25% { transform: rotate(0.5deg) translate(0.3px, -0.3px); }
  50% { transform: rotate(-0.3deg) translate(-0.2px, 0.2px); }
  75% { transform: rotate(0.3deg) translate(0.2px, 0.1px); }
}

.hand-drawn-icon:hover .icon-svg {
  animation: hand-drawn-wobble 0.3s ease-in-out;
}
</style>
