<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12" sm="6" md="3">
        <v-card class="rounded-lg" color="primary" variant="flat">
          <v-card-text class="text-white">
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ summary?.total_active || 0 }}</div>
                <div class="text-body-2 mt-1">活跃单元</div>
              </div>
              <v-icon size="48" class="opacity-50">mdi-timer</v-icon>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card class="rounded-lg" color="error" variant="flat">
          <v-card-text class="text-white">
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ (summary?.expiring_count || 0) + (summary?.expired_count || 0) }}</div>
                <div class="text-body-2 mt-1">告警单元</div>
              </div>
              <v-icon size="48" class="opacity-50">mdi-alert</v-icon>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card class="rounded-lg" color="warning" variant="flat">
          <v-card-text class="text-white">
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ pendingTodos }}</div>
                <div class="text-body-2 mt-1">待办事项</div>
              </div>
              <v-icon size="48" class="opacity-50">mdi-checkbox-marked-outline</v-icon>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card class="rounded-lg" color="info" variant="flat">
          <v-card-text class="text-white">
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ notificationStore.unreadCount }}</div>
                <div class="text-body-2 mt-1">未读通知</div>
              </div>
              <v-icon size="48" class="opacity-50">mdi-bell</v-icon>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="7">
        <v-card class="rounded-lg">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" color="error">mdi-alert-circle</v-icon>
            告警单元
            <v-spacer />
            <v-btn variant="text" size="small" to="/units">查看全部</v-btn>
          </v-card-title>
          <v-divider />
          <v-list v-if="alertUnits.length > 0">
            <v-list-item
              v-for="unit in alertUnits"
              :key="unit.id"
              :to="`/units/${unit.id}`"
              class="py-3"
            >
              <template v-slot:prepend>
                <v-icon :color="getAlertColor(unit)">{{ getUnitTypeIcon(unit.type) }}</v-icon>
              </template>
              <v-list-item-title>{{ unit.title }}</v-list-item-title>
              <v-list-item-subtitle>
                <span v-if="unit.type === 'time_countdown' && unit.remaining_seconds !== undefined">
                  {{ unit.remaining_seconds <= 0 ? '已超期' : `距到期还剩 ${formatDuration(unit.remaining_seconds)}` }}
                </span>
                <span v-else-if="unit.type === 'count_countdown' && unit.target_value">
                  已完成 {{ unit.current_value || 0 }} / {{ unit.target_value }} {{ unit.unit_label }}
                </span>
              </v-list-item-subtitle>
              <template v-slot:append>
                <v-chip :color="getAlertColor(unit)" size="small" variant="flat">
                  {{ getAlertLabel(unit) }}
                </v-chip>
              </template>
            </v-list-item>
          </v-list>
          <v-card-text v-else class="text-center text-medium-emphasis py-8">
            暂无告警单元
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="5">
        <v-card class="rounded-lg">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" color="warning">mdi-checkbox-marked-outline</v-icon>
            待办事项
            <v-spacer />
            <v-btn variant="text" size="small" to="/todos">查看全部</v-btn>
          </v-card-title>
          <v-divider />
          <v-list v-if="recentTodos.length > 0">
            <v-list-item
              v-for="todo in recentTodos"
              :key="todo.id"
              class="py-2"
            >
              <template v-slot:prepend>
                <v-checkbox-btn
                  :model-value="todo.status === 'done'"
                  @update:model-value="toggleTodo(todo)"
                  color="success"
                />
              </template>
              <v-list-item-title :class="{ 'text-decoration-line-through text-medium-emphasis': todo.status === 'done' }">
                {{ todo.title }}
              </v-list-item-title>
              <v-list-item-subtitle v-if="todo.due_date">
                截止: {{ formatDate(todo.due_date) }}
              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
          <v-card-text v-else class="text-center text-medium-emphasis py-8">
            暂无待办事项
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Unit, UnitSummary, Todo } from '@/types'
import { unitApi } from '@/api/units'
import { todoApi } from '@/api/todos'
import { useNotificationStore } from '@/stores/notifications'
import { formatDuration, formatDate, getUnitTypeIcon } from '@/utils/time'

const notificationStore = useNotificationStore()

const summary = ref<UnitSummary | null>(null)
const alertUnits = ref<Unit[]>([])
const recentTodos = ref<Todo[]>([])
const pendingTodos = ref(0)

function getAlertColor(unit: Unit): string {
  if (unit.type === 'time_countdown' && unit.remaining_seconds !== undefined) {
    if (unit.remaining_seconds <= 0) return 'error'
    if (unit.remaining_seconds <= 86400) return 'error'
    if (unit.remaining_seconds <= 7 * 86400) return 'warning'
  }
  return 'warning'
}

function getAlertLabel(unit: Unit): string {
  if (unit.type === 'time_countdown' && unit.remaining_seconds !== undefined) {
    if (unit.remaining_seconds <= 0) return '已超期'
    if (unit.remaining_seconds <= 86400) return '紧急'
    return '即将到期'
  }
  return '注意'
}

async function toggleTodo(todo: Todo) {
  const newStatus = todo.status === 'done' ? 'pending' : 'done'
  await todoApi.updateStatus(todo.id, newStatus)
  todo.status = newStatus as any
}

onMounted(async () => {
  try {
    const [summaryResp, unitsResp, todosResp] = await Promise.all([
      unitApi.getSummary(),
      unitApi.list({ status: 'active', sort_by: 'priority', sort_order: 'desc', page_size: 10 }),
      todoApi.list({ status: 'pending', page_size: 10 }),
    ])
    summary.value = summaryResp.data
    alertUnits.value = (unitsResp.data || []).filter((u: Unit) => {
      if (u.type === 'time_countdown' && u.remaining_seconds !== undefined) {
        return u.remaining_seconds <= 7 * 86400
      }
      return false
    })
    recentTodos.value = todosResp.data || []
    pendingTodos.value = todosResp.meta?.total || 0
  } catch {
    // ignore
  }
})
</script>
