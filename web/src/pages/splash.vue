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
      <div class="d-flex flex-column justify-center align-center text-center">
        <Logo />
        <PageTitle  
          title="Filament"
          subtitle="Your Secure, Self-Custodial Filecoin Wallet"
          class="title-white subtitle-white"
        />
        <v-progress-circular
          color="white"
          indeterminate
          :size="40"
          width="2"
        />
      </div>
    </v-container>
  </div>
</template>
<script setup lang="ts">
  import { onMounted, nextTick } from 'vue';
  import { walletClient } from '@/api/client';
  import { useStore } from '@/store/app';
  
  const router = useRouter()
  const route = useRoute()
  const store = useStore()
  const wallet = walletClient

  onMounted(async () => {
    await sleep(10000)
    await initApp()
    
  })

  const initApp = async () => {
    try {
      await store.initApp()
      

      if (store.wallets.wallets.length === 0) {
        await router.replace("/onboard" )
      } else {
        const redirect = typeof route.query.redirect === "string" ? route.query.redirect : undefined 
        if (redirect) {
          const resolved = router.resolve(redirect)
          if (resolved.matched && resolved.matched.length > 0) {
            await router.replace(redirect)
            return
          }
        }
        await router.replace("/overview")
      }
    } catch (err) {

    }
  }

  function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms))
  }
</script>
<style scoped>
#splash-page {
  background: var(--v-primary-gradient);
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
