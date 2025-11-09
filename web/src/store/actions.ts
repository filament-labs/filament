import { walletClient } from '@/api/client'
import type { WalletsResponse } from '@/pb/wallet_pb'
import type { IState } from '@/types/store'
import { Centrifuge } from 'centrifuge'

export type IActions = {
  initApp(): void
  unlockWallet(): void
} & ThisType<IState & IActions>

const actions: IActions = {
  async initApp() {
    const wallet = walletClient
    this.wallets = await wallet.getWallets({})
    this.appInitialized = true
  },
  unlockWallet () {
   
  },
}

export default actions
