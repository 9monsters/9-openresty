<script setup lang="ts">
import { t } from '@/i18n'
// Placeholder for SSH host setup page — matches nginx-ui's Nginx Host SSH Setup
</script>

<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">{{ t('Host SSH Setup') }}</h3>
    <ACard>
      <AEmpty :description="t('Configure SSH access to the Nginx host for remote management.')" />
    </ACard>
  </div>
</template>