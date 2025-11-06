<route>
  {
    "meta": {
      "layout": "landing"
    }
  }
</route>
<template>
  <div id="splash-page" class="page fill-height">
    <v-container class="d-flex flex-column justify-center fill-height">
      <Logo />
      <h1 class="title">Filament</h1>
      <p class="subtitle">Your Secure, Self-Custodial Filecoin Wallet</p>
      <!--<v-progress-circular
        color="white"
        indeterminate
        :size="80"
        width="2"
      />-->
    </v-container>
  </div>
</template>
<style scoped>
#splash-page {
  background: linear-gradient(135deg, #00D4FF 0%, #0066FF 100%);
}

.title {
  color: rgb(255 255 255 / 0.95);
  font-size: 2.25rem;
  line-height: 2.5rem;
  margin-bottom: 1rem;
}

.subtitle {
  color: rgb(255 255 255 / 0.9);
  font-size: 1.25rem;
  line-height: 1.75rem;
  margin-bottom: 2rem;
}

.v-progress-circular {
  margin: 20px 0 100px 0;
}
</style>
