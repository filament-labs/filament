import type { Centrifuge } from 'centrifuge'

export interface IState {
  appInitialized: boolean
  centrifuge: Centrifuge | null
}
