<template>
  <v-sheet
    class="mx-auto mb-8"
    :class="variant"
    elevation="0"
    :height="sheetHeight"
    :width="sheetWidth"
    :rounded="sheetRadius"
  >
    <svg
      fill="white"
      :height="svgHeight"
      viewBox="0 0 24 20"
      :width="svgWidth"
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

  const sheetHeight = ref(props.size == "small" ? 60 : 100)
  const sheetWidth = ref(props.size == "small" ? 60 : 100)
  const svgHeight = ref(props.size == "small" ? 30 : 40)
  const svgWidth = ref(props.size == "small" ? 30 : 50)
  const sheetRadius = ref(props.size == "small" ? "lg" : "xl")//ref(props.size == "small" ? 12 : 24)
</script>
<style scoped>
.v-sheet {
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.v-sheet.normal {
  background-color: rgba(255,255,255,0.2);
}

.v-sheet.inversed {
  background-image: linear-gradient(135deg, rgb(0, 212, 255) 0%, rgb(0, 102, 255) 100%) !important;
  background-color: red;
}
</style>
