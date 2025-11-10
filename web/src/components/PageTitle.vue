<template>
  <div class="page-title">
    <h2 class="mb-1">{{ title }}</h2>
    <p class="mb-4">{{ subtitle }}</p>
  </div>
</template>
<script lang="ts" setup>
  interface Props {
    title: string 
    subtitle?: string
  }

  defineProps<Props>()
</script>
<style scoped>
  .title-white h1 {
    color: #fff
  }

  .subtitle-white p {
    color: #fff;
  }

  .subtitle-gray p {
    color: var(--v-gray-600);
  }
</style>