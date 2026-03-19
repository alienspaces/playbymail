<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
      <div class="flow-modal">
        <div class="modal-header">
          <h2>State Flow — {{ objectName }}</h2>
          <button class="close-btn" @click="$emit('close')" aria-label="Close">✕</button>
        </div>

        <div class="modal-content">
          <div v-if="!hasTransitions" class="empty-state">
            No state transitions defined. Add effects with effect type <strong>change_state</strong> to see the flow.
          </div>
          <svg
            v-else
            :viewBox="`0 0 ${svgWidth} ${svgHeight}`"
            :width="svgWidth"
            :height="svgHeight"
            class="flow-svg"
            xmlns="http://www.w3.org/2000/svg"
          >
            <defs>
              <marker
                id="arrowhead"
                markerWidth="10"
                markerHeight="7"
                refX="9"
                refY="3.5"
                orient="auto"
              >
                <polygon points="0 0, 10 3.5, 0 7" class="arrowhead-fill" />
              </marker>
              <marker
                id="arrowhead-initial"
                markerWidth="10"
                markerHeight="7"
                refX="9"
                refY="3.5"
                orient="auto"
              >
                <polygon points="0 0, 10 3.5, 0 7" class="arrowhead-initial-fill" />
              </marker>
            </defs>

            <!-- Initial state marker (filled circle) -->
            <circle
              v-if="initialNodeLayout"
              :cx="initialNodeLayout.cx"
              :cy="initialNodeLayout.cy"
              r="10"
              class="start-circle"
            />

            <!-- Arrow from initial marker to first/initial state node -->
            <line
              v-if="initialArrow"
              :x1="initialArrow.x1"
              :y1="initialArrow.y1"
              :x2="initialArrow.x2"
              :y2="initialArrow.y2"
              class="transition-line initial-line"
              marker-end="url(#arrowhead-initial)"
            />

            <!-- Transition arrows -->
            <g v-for="(edge, i) in edgeLayouts" :key="i">
              <path
                :d="edge.path"
                class="transition-line"
                :class="{ 'self-loop': edge.isSelf }"
                marker-end="url(#arrowhead)"
                fill="none"
              />
              <rect
                :x="edge.labelX - edge.labelWidth / 2 - 4"
                :y="edge.labelY - 9"
                :width="edge.labelWidth + 8"
                height="18"
                rx="3"
                class="edge-label-bg"
              />
              <text
                :x="edge.labelX"
                :y="edge.labelY + 5"
                text-anchor="middle"
                class="edge-label"
              >{{ edge.label }}</text>
            </g>

            <!-- State nodes -->
            <g v-for="node in nodeLayouts" :key="node.id">
              <rect
                :x="node.x"
                :y="node.y"
                :width="NODE_W"
                :height="NODE_H"
                rx="8"
                class="state-node"
                :class="{ 'state-node-initial': node.isInitial }"
              />
              <text
                :x="node.x + NODE_W / 2"
                :y="node.y + NODE_H / 2 + 5"
                text-anchor="middle"
                class="state-label"
              >{{ node.name }}</text>
            </g>
          </svg>
        </div>

        <div class="modal-footer">
          <button class="close-btn-secondary" @click="$emit('close')">Close</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  visible: { type: Boolean, default: false },
  objectName: { type: String, default: '' },
  /** @type {import('../types').GameLocationObjectState[]} */
  states: { type: Array, default: () => [] },
  /** @type {import('../types').GameLocationObjectEffect[]} */
  effects: { type: Array, default: () => [] },
  initialStateId: { type: String, default: null },
});

defineEmits(['close']);

// ── Layout constants ──────────────────────────────────────────────────────────

const NODE_W = 160;
const NODE_H = 44;
const H_GAP = 100; // horizontal gap either side of the node column
const V_GAP = 80;  // vertical gap between nodes
const START_MARKER_OFFSET = 60; // space above first node for the start circle
const SVG_PADDING = 24;
// Extra horizontal space for self-loop arcs on the right side
const SELF_LOOP_W = 60;

// ── Derived data ──────────────────────────────────────────────────────────────

const transitions = computed(() =>
  props.effects.filter(
    (e) =>
      e.effect_type === 'change_state' &&
      e.result_adventure_game_location_object_state_id
  )
);

const hasTransitions = computed(() => transitions.value.length > 0);

const stateById = computed(() => {
  const map = {};
  for (const s of props.states) map[s.id] = s;
  return map;
});

// ── Node layout ───────────────────────────────────────────────────────────────

/**
 * Order states: initial state first if present, then remaining sorted by sort_order.
 */
const orderedStates = computed(() => {
  const sorted = [...props.states].sort((a, b) => (a.sort_order ?? 0) - (b.sort_order ?? 0));
  if (!props.initialStateId) return sorted;
  const idx = sorted.findIndex((s) => s.id === props.initialStateId);
  if (idx <= 0) return sorted;
  const moved = sorted.splice(idx, 1)[0];
  sorted.unshift(moved);
  return sorted;
});

const nodeLayouts = computed(() => {
  const cx = SVG_PADDING + H_GAP + NODE_W / 2;
  const x = cx - NODE_W / 2;
  return orderedStates.value.map((state, i) => ({
    id: state.id,
    name: state.name,
    x,
    y: SVG_PADDING + START_MARKER_OFFSET + i * (NODE_H + V_GAP),
    cx,
    cy: SVG_PADDING + START_MARKER_OFFSET + i * (NODE_H + V_GAP) + NODE_H / 2,
    isInitial: state.id === props.initialStateId,
  }));
});

const nodeMap = computed(() => {
  const m = {};
  for (const n of nodeLayouts.value) m[n.id] = n;
  return m;
});

// ── Initial state marker layout ───────────────────────────────────────────────

const initialNodeLayout = computed(() => {
  if (!props.initialStateId) {
    // Use the first node as the de-facto start
    return nodeLayouts.value.length > 0
      ? { cx: nodeLayouts.value[0].cx, cy: SVG_PADDING + START_MARKER_OFFSET / 2 }
      : null;
  }
  const node = nodeMap.value[props.initialStateId];
  if (!node) return null;
  return { cx: node.cx, cy: SVG_PADDING + START_MARKER_OFFSET / 2 };
});

const initialArrow = computed(() => {
  if (!initialNodeLayout.value) return null;
  const targetNode = props.initialStateId
    ? nodeMap.value[props.initialStateId]
    : nodeLayouts.value[0];
  if (!targetNode) return null;
  return {
    x1: initialNodeLayout.value.cx,
    y1: initialNodeLayout.value.cy + 10,
    x2: targetNode.cx,
    y2: targetNode.y,
  };
});

// ── Edge layout ───────────────────────────────────────────────────────────────

/**
 * Estimate text width (rough: ~7.5px per char at font-size 12).
 */
function textWidth(str) {
  return Math.max(str.length * 7.5, 40);
}

/**
 * For multiple edges between the same pair of nodes, offset them so labels
 * don't overlap. We track pairs and assign an index to each.
 */
const edgeLayouts = computed(() => {
  const pairCount = {};
  const edges = [];

  for (const t of transitions.value) {
    const fromId = t.required_adventure_game_location_object_state_id;
    const toId = t.result_adventure_game_location_object_state_id;

    const fromNode = fromId ? nodeMap.value[fromId] : nodeLayouts.value[0];
    const toNode = nodeMap.value[toId];
    if (!fromNode || !toNode) continue;

    const pairKey = `${fromNode.id}→${toNode.id}`;
    pairCount[pairKey] = (pairCount[pairKey] || 0) + 1;
    const edgeIndex = pairCount[pairKey] - 1;

    const isSelf = fromNode.id === toNode.id;
    const label = t.action_type || '';
    const lw = textWidth(label);

    let path, labelX, labelY;

    if (isSelf) {
      // Self-loop: arc to the right of the node
      const loopX = fromNode.x + NODE_W;
      const loopY = fromNode.cy;
      const loopRadius = 30 + edgeIndex * 16;
      const startX = loopX;
      const startY = loopY - 10;
      const endX = loopX;
      const endY = loopY + 10;
      path = `M ${startX} ${startY} C ${loopX + loopRadius} ${startY}, ${loopX + loopRadius} ${endY}, ${endX} ${endY}`;
      labelX = loopX + loopRadius + 4 + lw / 2;
      labelY = loopY;
    } else {
      // Straight arrow with horizontal offset for parallel edges
      const offset = (edgeIndex - Math.floor(pairCount[pairKey] / 2)) * 24;
      const x1 = fromNode.cx + offset;
      const y1 = fromNode.y + NODE_H;
      const x2 = toNode.cx + offset;
      const y2 = toNode.y;
      const midX = (x1 + x2) / 2;
      const midY = (y1 + y2) / 2;

      if (Math.abs(offset) < 1) {
        path = `M ${x1} ${y1} L ${x2} ${y2}`;
      } else {
        // Slight curve for offset edges
        path = `M ${x1} ${y1} Q ${midX + offset} ${midY} ${x2} ${y2}`;
      }
      labelX = midX + offset * 0.5;
      labelY = midY;
    }

    edges.push({ path, label, labelX, labelY, labelWidth: lw, isSelf });
  }

  return edges;
});

// ── SVG canvas dimensions ─────────────────────────────────────────────────────

const svgWidth = computed(() => {
  const hasSelfLoop = edgeLayouts.value.some((e) => e.isSelf);
  return SVG_PADDING * 2 + H_GAP * 2 + NODE_W + (hasSelfLoop ? SELF_LOOP_W + 40 : 0);
});

const svgHeight = computed(() => {
  const n = orderedStates.value.length;
  return SVG_PADDING * 2 + START_MARKER_OFFSET + n * NODE_H + Math.max(0, n - 1) * V_GAP;
});
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: var(--space-md);
}

.flow-modal {
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 100%;
  max-width: 640px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-md) var(--space-lg);
  border-bottom: 1px solid var(--color-border);
  background: var(--color-bg-light);
}

.modal-header h2 {
  margin: 0;
  font-size: var(--font-size-lg);
  color: var(--color-text);
}

.close-btn {
  background: none;
  border: none;
  font-size: var(--font-size-lg);
  cursor: pointer;
  color: var(--color-text-muted);
  padding: var(--space-xs);
  border-radius: var(--radius-sm);
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: var(--color-border);
  color: var(--color-text);
}

.modal-content {
  flex: 1;
  overflow: auto;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: var(--space-lg);
  min-height: 200px;
}

.empty-state {
  color: var(--color-text-muted);
  font-size: 0.9rem;
  text-align: center;
  padding: var(--space-xl);
  align-self: center;
}

.flow-svg {
  display: block;
  overflow: visible;
}

/* SVG element styles — using fill/stroke directly since CSS vars work inside SVG */
.state-node {
  fill: var(--color-surface, #fff);
  stroke: var(--color-border, #cbd5e1);
  stroke-width: 1.5;
}

.state-node-initial {
  fill: var(--color-primary-light, #eff6ff);
  stroke: var(--color-primary, #3b82f6);
  stroke-width: 2;
}

.state-label {
  fill: var(--color-text, #1e293b);
  font-size: 13px;
  font-weight: 600;
  font-family: inherit;
}

.start-circle {
  fill: var(--color-text, #1e293b);
}

.transition-line {
  stroke: var(--color-text-muted, #94a3b8);
  stroke-width: 1.5;
}

.initial-line {
  stroke: var(--color-text, #1e293b);
  stroke-width: 2;
}

.arrowhead-fill {
  fill: var(--color-text-muted, #94a3b8);
}

.arrowhead-initial-fill {
  fill: var(--color-text, #1e293b);
}

.edge-label-bg {
  fill: var(--color-bg, #fff);
  stroke: none;
}

.edge-label {
  fill: var(--color-text-secondary, #64748b);
  font-size: 11px;
  font-family: inherit;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  padding: var(--space-md) var(--space-lg);
  border-top: 1px solid var(--color-border);
  background: var(--color-bg-light);
}

.close-btn-secondary {
  padding: var(--space-sm) var(--space-md);
  background: var(--color-bg);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: var(--font-weight-bold);
}

.close-btn-secondary:hover {
  background: var(--color-border);
}
</style>
