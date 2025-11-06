import type { IState } from '@/types/store'
import { defineStore } from 'pinia'
import actions, { type IActions } from './actions'
import getters from './getters'
import state from './state'

export const useStore = defineStore('app', {
  state,
  actions,
  getters,
} as {
  state: () => IState
  actions: IActions
  getters: typeof getters
  //getters: Record<string, any>
})
