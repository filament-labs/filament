<template>
  <v-btn 
    v-bind="$attrs"
    class="btn"
  >
    <slot />
  </v-btn>
</template>
<script lang="ts" setup>
  defineOptions({ inheritAttrs: false })
</script>
<style lang="scss" scoped>
.btn:not([color])[data-gradient]{
  &[data-gradient="primary"] {
    background: var(--v-secondary-gradient);
  }
}
</style>