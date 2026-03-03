<template>
  <router-view v-slot="{ Component }">
    <transition name="page" mode="out-in">
      <component :is="Component" />
    </transition>
  </router-view>
</template>

<script setup lang="ts">
import { useDark } from '@vueuse/core'

// Initialize theme globally to ensure persistence across all pages
useDark()
</script>