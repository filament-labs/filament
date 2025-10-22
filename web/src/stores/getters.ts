import type { IState } from '@/types/store'

const getters = {
  isAppInitialized: (state: IState) => state.appInitialized,
}
export default getters
