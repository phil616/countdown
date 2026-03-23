import { get, post, put, patch, del } from './client'
import type { Todo, TodoGroup } from '@/types'

export const todoApi = {
  list: (params?: Record<string, any>) => get<Todo[]>('/todos', params),

  create: (data: any) => post<Todo>('/todos', data),

  getById: (id: string) => get<Todo>(`/todos/${id}`),

  update: (id: string, data: any) => put<Todo>(`/todos/${id}`, data),

  patchUpdate: (id: string, data: any) => patch<Todo>(`/todos/${id}`, data),

  delete: (id: string) => del(`/todos/${id}`),

  updateStatus: (id: string, status: string) =>
    patch<Todo>(`/todos/${id}/status`, { status }),

  batchAction: (action: 'complete' | 'delete', ids: string[]) =>
    post('/todos/batch', { action, ids }),

  listGroups: () => get<TodoGroup[]>('/todo-groups'),

  createGroup: (data: { name: string; color?: string; sort_order?: number }) =>
    post<TodoGroup>('/todo-groups', data),

  updateGroup: (id: string, data: any) =>
    put<TodoGroup>(`/todo-groups/${id}`, data),

  deleteGroup: (id: string) => del(`/todo-groups/${id}`),
}
