<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h2 class="text-h5 font-weight-bold">通知中心</h2>
      <v-spacer />
      <v-btn variant="outlined" prepend-icon="mdi-check-all" @click="markAllRead" :disabled="unreadOnly && notifications.length === 0">
        全部已读
      </v-btn>
    </div>

    <v-card class="rounded-lg mb-4 pa-3">
      <v-btn-toggle v-model="unreadOnly" variant="outlined" density="compact">
        <v-btn :value="false">全部</v-btn>
        <v-btn :value="true">未读</v-btn>
      </v-btn-toggle>
    </v-card>

    <v-card class="rounded-lg">
      <v-list v-if="notifications.length > 0">
        <v-list-item
          v-for="n in notifications"
          :key="n.id"
          :class="{ 'bg-grey-lighten-4': !n.is_read }"
          class="py-3"
        >
          <template v-slot:prepend>
            <v-icon :color="levelColor(n.level)">{{ levelIcon(n.level) }}</v-icon>
          </template>
          <v-list-item-title :class="{ 'font-weight-bold': !n.is_read }">
            {{ n.message }}
          </v-list-item-title>
          <v-list-item-subtitle>
            {{ formatDateTime(n.triggered_at) }}
            <span v-if="n.unit_title"> · {{ n.unit_title }}</span>
          </v-list-item-subtitle>
          <template v-slot:append>
            <v-btn v-if="!n.is_read" icon="mdi-check" size="small" variant="text" @click="markRead(n.id)" />
            <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="deleteNotif(n.id)" />
          </template>
        </v-list-item>
      </v-list>
      <v-card-text v-else class="text-center text-medium-emphasis py-8">
        暂无通知
      </v-card-text>
    </v-card>

    <div class="d-flex justify-center mt-4" v-if="totalPages > 1">
      <v-pagination v-model="page" :length="totalPages" @update:model-value="fetchNotifications" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { Notification } from '@/types'
import { notificationApi } from '@/api/notifications'
import { useNotificationStore } from '@/stores/notifications'
import { formatDateTime } from '@/utils/time'

const notificationStore = useNotificationStore()
const notifications = ref<Notification[]>([])
const page = ref(1)
const totalPages = ref(0)
const unreadOnly = ref(false)

function levelColor(level: string): string {
  return { info: 'info', warning: 'warning', critical: 'error' }[level] || 'info'
}

function levelIcon(level: string): string {
  return { info: 'mdi-information', warning: 'mdi-alert', critical: 'mdi-alert-circle' }[level] || 'mdi-information'
}

async function fetchNotifications() {
  try {
    const params: Record<string, any> = { page: page.value, page_size: 20 }
    if (unreadOnly.value) params.is_read = false
    const resp = await notificationApi.list(params)
    notifications.value = resp.data || []
    totalPages.value = resp.meta?.total_pages || 0
  } catch { /* ignore */ }
}

async function markRead(id: string) {
  await notificationApi.markAsRead(id)
  notificationStore.fetchUnreadCount()
  fetchNotifications()
}

async function markAllRead() {
  await notificationApi.markAllAsRead()
  notificationStore.fetchUnreadCount()
  fetchNotifications()
}

async function deleteNotif(id: string) {
  await notificationApi.delete(id)
  notificationStore.fetchUnreadCount()
  fetchNotifications()
}

watch(unreadOnly, () => { page.value = 1; fetchNotifications() })

onMounted(fetchNotifications)
</script>
