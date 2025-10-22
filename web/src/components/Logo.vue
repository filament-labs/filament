<template>
  <v-sheet
    class="mx-auto mb-8"
    elevation="0"
    height="100"
    rounded="xl"
    width="100"
  >
    <svg
      fill="white"
      height="40"
      viewBox="0 0 24 20"
      width="50"
      xmlns="http://www.w3.org/2000/svg"
    >
      <polygon points="12,0 21,5 21,15 12,20 3,15 3,5" />
    </svg>
  </v-sheet>
</template>
<script lang="ts" setup>
  interface IProps {
    variant?: 'normal' | 'inversed'
    size?: 'normal' | 'small'
  }

  const props = withDefaults(defineProps<IProps>(), {
    variant: 'normal',
    size: 'normal',
  })

</script>
<style scoped>
.v-sheet {
  background-color: rgba(255,255,255,0.2);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
