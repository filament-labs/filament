import type { IState } from '@/types/store'

const state: () => IState = () => ({
  appInitialized: false,
  centrifuge: null,
})

export default state
