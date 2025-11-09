import type { IState } from '@/types/store'

const state: () => IState = () => ({
  appInitialized: false,
})

export default state
