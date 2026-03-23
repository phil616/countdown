import { defineStore } from 'pinia'
import { ref } from 'vue'
import { notificationApi } from '@/api/notifications'

export const useNotificationStore = defineStore('notifications', () => {
  const unreadCount = ref(0)

  async function fetchUnreadCount() {
    try {
      const resp = await notificationApi.unreadCount()
      unreadCount.value = resp.data.count
    } catch {
      // ignore
    }
  }

  let intervalId: number | null = null

  function startPolling() {
    fetchUnreadCount()
    intervalId = window.setInterval(fetchUnreadCount, 60000)
  }

  function stopPolling() {
    if (intervalId !== null) {
      clearInterval(intervalId)
      intervalId = null
    }
  }

  return { unreadCount, fetchUnreadCount, startPolling, stopPolling }
})
