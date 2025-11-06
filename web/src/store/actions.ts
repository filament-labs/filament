import type { IState } from '@/types/store'
import { Centrifuge } from 'centrifuge'

export type IActions = {
  init(): void
} & ThisType<IState & IActions>

const actions: IActions = {
  init () {
    if (this.centrifuge) {
      return
    }

    this.centrifuge = new Centrifuge('ws://localhost:8000/connection/websocket')
    this.centrifuge.on('connected', function(){
      console.log("connected")
    })

    this.centrifuge.on('disconnected', function(ctx){
      console.log("disconnected--", ctx.reason)
    })
    this.centrifuge.connect()
  },
}

export default actions
