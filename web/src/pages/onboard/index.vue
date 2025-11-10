<route>
  {
    "meta": {
      "layout": "landing"
    }
  }
</route>
<template>
  <div id="onboarding-page" class="page fill-height">
    <v-container class="d-flex flex-column justify-center fill-height">
      <div class="content-container d-flex flex-column justify-center text-center">
        <Logo variant="inversed" />
        <page-title 
          title="Filament"
          subtitle="The most secure way to store and manage your Filecoin. 
          Get started by creating a new wallet or importing an existing one."
          class="subtitle-gray"
        />
        <Btn    
          data-gradient="primary"
          class="text-white mb-4"
          size="x-large"
          rounded="sm"
        >
          Create New Wallet
        </Btn>
        <Btn 
          size="x-large" 
          rounded="sm"
          variant="outlined"
          to="/onboard/import"
        >
          Import Existing Wallet
        </Btn>
      </div>
    </v-container>
  </div>
</template>
<style scoped>
.page {
  background-color: rgb(var(--v-theme-background));
}

.content-container {
  max-width: 40%;
}
</style>