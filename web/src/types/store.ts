import type { WalletsResponse } from '@/pb/wallet_pb'

export interface IState {
  appInitialized: boolean
  wallets: WalletsResponse
}
