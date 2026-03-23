import { get, post, put, patch, del } from './client'
import type { Unit, UnitSummary, UnitLog } from '@/types'

export const unitApi = {
  list: (params?: Record<string, any>) => get<Unit[]>('/units', params),

  create: (data: any) => post<Unit>('/units', data),

  getById: (id: string) => get<Unit>(`/units/${id}`),

  update: (id: string, data: any) => put<Unit>(`/units/${id}`, data),

  patchUpdate: (id: string, data: any) => patch<Unit>(`/units/${id}`, data),

  delete: (id: string) => del(`/units/${id}`),

  updateStatus: (id: string, status: string) =>
    patch<Unit>(`/units/${id}/status`, { status }),

  /** 将单元归属到指定项目 */
  assignToProject: (id: string, projectId: string) =>
    patch<Unit>(`/units/${id}`, { project_id: projectId }),

  /** 将单元从所属项目中移除（project_id 置为 null） */
  removeFromProject: (id: string) =>
    patch<Unit>(`/units/${id}`, { clear_project: true }),

  step: (id: string, direction: 'up' | 'down', note?: string) =>
    post<Unit>(`/units/${id}/step`, { direction, note }),

  setValue: (id: string, value: number, note?: string) =>
    put<Unit>(`/units/${id}/value`, { value, note }),

  getLogs: (id: string, params?: Record<string, any>) =>
    get<UnitLog[]>(`/units/${id}/logs`, params),

  getSummary: () => get<UnitSummary>('/units/summary'),
}
