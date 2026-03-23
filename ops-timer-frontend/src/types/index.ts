export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  meta?: PaginationMeta
  errors?: FieldError[]
}

export interface PaginationMeta {
  page: number
  page_size: number
  total: number
  total_pages: number
}

export interface FieldError {
  field: string
  message: string
}

export interface User {
  id: string
  username: string
  display_name: string
  email: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface TokenResponse {
  api_token: string
}

export interface Unit {
  id: string
  project_id: string | null
  title: string
  description: string
  type: UnitType
  status: UnitStatus
  priority: Priority
  tags: string[]
  color: string

  target_time?: string
  start_time?: string
  display_unit?: string
  remind_before_days?: number[]
  remind_after_days?: number[]

  current_value?: number
  target_value?: number
  step?: number
  unit_label?: string
  allow_exceed?: boolean
  remind_on_values?: number[]

  remaining_seconds?: number
  elapsed_seconds?: number
  progress?: number

  created_at: string
  updated_at: string
}

export type UnitType = 'time_countdown' | 'time_countup' | 'count_countdown' | 'count_countup'
export type UnitStatus = 'active' | 'paused' | 'completed' | 'archived'
export type Priority = 'low' | 'normal' | 'high' | 'critical'

export interface UnitSummary {
  total_active: number
  total_paused: number
  total_completed: number
  total_archived: number
  expiring_count: number
  expired_count: number
}

export interface UnitLog {
  id: string
  unit_id: string
  delta: number
  value_before: number
  value_after: number
  note: string
  operated_at: string
}

export interface Project {
  id: string
  title: string
  description: string
  status: string
  color: string
  icon: string
  sort_order: number
  created_at: string
  updated_at: string
  unit_stats?: ProjectUnitStats
}

export interface ProjectUnitStats {
  active_count: number
  expiring_count: number
  completed_count: number
  total_count: number
}

export interface Todo {
  id: string
  group_id: string | null
  title: string
  description: string
  status: TodoStatus
  priority: Priority
  due_date: string | null
  sort_order: number
  completed_at: string | null
  created_at: string
  updated_at: string
}

export type TodoStatus = 'pending' | 'in_progress' | 'done' | 'cancelled'

export interface TodoGroup {
  id: string
  name: string
  color: string
  sort_order: number
  todo_count: number
  created_at: string
  updated_at: string
}

export interface Notification {
  id: string
  unit_id: string
  level: 'info' | 'warning' | 'critical'
  message: string
  is_read: boolean
  triggered_at: string
  read_at: string | null
  unit_title?: string
}

export interface UnreadCount {
  count: number
}
