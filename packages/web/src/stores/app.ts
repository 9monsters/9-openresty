import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', () => {
  const message = ref<any>(null)
  const modal = ref<any>(null)
  const notification = ref<any>(null)

  function initApp(msg: any, mod: any, notif: any) {
    message.value = msg
    modal.value = mod
    notification.value = notif
  }

  return { message, modal, notification, initApp }
})