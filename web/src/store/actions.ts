import type { WalletsResponse } from '@/pb/wallet_pb'
import type { IState } from '@/types/store'
import { Centrifuge } from 'centrifuge'

export type IActions = {
  initApp(wallets: WalletsResponse): void
  unlockWallet(): void
} & ThisType<IState & IActions>

const actions: IActions = {
  initApp(wallets) {
    this.appInitialized = true
    this.wallets = wallets
  },
  unlockWallet () {
   
  },
}

export default actions
