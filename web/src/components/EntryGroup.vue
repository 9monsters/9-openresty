<script setup lang="ts">
import { useRouter } from 'vue-router'

export interface NavEntry {
  title: string
  description: string
  path: string
  icon: string
}

defineProps<{ entries: NavEntry[] }>()

const router = useRouter()
</script>

<template>
  <div class="ios-group">
    <button
      v-for="entry in entries"
      :key="entry.path"
      type="button"
      class="ios-row"
      @click="router.push(entry.path)"
    >
      <span class="ios-row-icon" :class="entry.icon" aria-hidden="true" />
      <span class="entry-copy">
        <span class="ios-headline">{{ entry.title }}</span>
        <span class="ios-footnote entry-desc">{{ entry.description }}</span>
      </span>
      <span class="ios-chevron i-tabler:chevron-right" aria-hidden="true" />
    </button>
  </div>
</template>

<style scoped>
.ios-row { width: 100%; text-align: left; background: transparent; border: 0; font: inherit; }
.ios-row + .ios-row { border-top: 0.5px solid var(--separator); }

.entry-copy {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.entry-desc {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.38;
}
</style>
