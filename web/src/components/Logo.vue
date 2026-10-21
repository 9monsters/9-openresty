<script setup lang="ts">
// The mark is a pair of config braces on a glass tile: this product is a tree
// of nginx blocks, and the tile echoes the concentric corners of the shell.
const isCollapsed = inject<Ref<boolean>>('sidebarCollapsed', ref(false))
</script>

<template>
  <div class="logo">
    <span class="mark" aria-hidden="true">
      <svg viewBox="0 0 32 32" width="32" height="32" fill="none">
        <defs>
          <linearGradient id="orx-tile" x1="4" y1="2" x2="28" y2="30" gradientUnits="userSpaceOnUse">
            <stop stop-color="#3AA0FF" />
            <stop offset="1" stop-color="#5B5BD6" />
          </linearGradient>
        </defs>
        <rect x="0.5" y="0.5" width="31" height="31" rx="9.5" fill="url(#orx-tile)" />
        <rect x="0.5" y="0.5" width="31" height="31" rx="9.5" stroke="rgba(255,255,255,.45)" stroke-width="1" />
        <path
          d="M13.4 8.6c-2 0-2.4 1.1-2.4 2.6v2.1c0 1.3-.5 2-1.7 2.7 1.2.7 1.7 1.4 1.7 2.7v2.1c0 1.5.4 2.6 2.4 2.6"
          stroke="#fff" stroke-width="2.1" stroke-linecap="round" stroke-linejoin="round"
        />
        <path
          d="M18.6 8.6c2 0 2.4 1.1 2.4 2.6v2.1c0 1.3.5 2 1.7 2.7-1.2.7-1.7 1.4-1.7 2.7v2.1c0 1.5-.4 2.6-2.4 2.6"
          stroke="#fff" stroke-width="2.1" stroke-linecap="round" stroke-linejoin="round"
        />
      </svg>
    </span>
    <span v-if="!isCollapsed" class="word">
      <span class="word-name">OpenResty</span>
      <span class="word-role">Admin</span>
    </span>
  </div>
</template>

<style scoped>
.logo {
  flex: none;
  display: flex;
  align-items: center;
  gap: 11px;
  height: 46px;
  padding: 0 10px;
  overflow: hidden;
}

.mark {
  flex: none;
  display: grid;
  place-items: center;
  width: 32px;
  height: 32px;
  border-radius: 9.5px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.16), 0 6px 14px -6px rgba(0, 0, 0, 0.4);
}

.word { display: flex; flex-direction: column; line-height: 1.15; white-space: nowrap; }
.word-name { font-size: 15px; font-weight: 650; letter-spacing: -0.016em; color: var(--label); }
.word-role { font-size: 12px; font-weight: 500; color: var(--label-2); }
</style>
