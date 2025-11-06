import type { IState } from '@/types/store'

const getters = {
  isAppInitialized(this: IState) {
    return this.appInitialized
  },
}


export default getters
