<template>
  <p class="field-hint">
    <slot />
  </p>
</template>

<script setup>
// FieldHint renders secondary help text below a form field. It exists to keep
// the styling consistent across every designer studio form and to avoid
// copy-pasting the same scoped `.field-hint` rule into each view.
//
// Usage:
//   <label>Speed <span class="required">*</span></label>
//   <input v-model.number="form.speed" type="number" min="1" max="10" />
//   <FieldHint>Sector hops allowed per turn (1–10).</FieldHint>
//
// Attrs fall through to the root <p>, so `id` can be set and referenced via
// `aria-describedby` on the associated input.
</script>

<style scoped>
.field-hint {
  margin: var(--space-xs) 0 0;
  font-size: var(--font-size-sm, 0.875rem);
  color: var(--color-text-muted);
}
</style>
