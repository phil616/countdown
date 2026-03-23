import { get, post, patch, del } from './client'
import type { Notification, UnreadCount } from '@/types'

export const notificationApi = {
  list: (params?: Record<string, any>) =>
    get<Notification[]>('/notifications', params),

  markAsRead: (id: string) =>
    patch(`/notifications/${id}/read`),

  markAllAsRead: () => post('/notifications/read-all'),

  unreadCount: () => get<UnreadCount>('/notifications/unread-count'),

  delete: (id: string) => del(`/notifications/${id}`),
}
