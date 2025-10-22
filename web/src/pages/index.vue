<template>
  <component :is="ActivePage" />
</template>

<script lang="ts" setup>
  import { useStore } from '@/stores/app'
  import Overview from './overview.vue'
  import Splash from './splash.vue'

  const store = useStore()

  const ActivePage = computed(() => {
    return store.isAppInitialized ? Overview : Splash
  })
</script>
