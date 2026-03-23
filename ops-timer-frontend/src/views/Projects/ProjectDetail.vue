<template>
  <div v-if="project">
    <div class="d-flex align-center mb-4">
      <v-btn icon="mdi-arrow-left" variant="text" @click="$router.back()" />
      <v-icon v-if="project.icon" class="ml-2 mr-2">{{ project.icon }}</v-icon>
      <h2 class="text-h5 font-weight-bold">{{ project.title }}</h2>
      <v-spacer />
      <v-chip :color="getStatusColor(project.status)">{{ getStatusLabel(project.status) }}</v-chip>
    </div>

    <v-row class="mb-4" v-if="project.unit_stats">
      <v-col cols="6" sm="3">
        <v-card class="rounded-lg text-center pa-4" variant="tonal" color="success">
          <div class="text-h5 font-weight-bold">{{ project.unit_stats.active_count }}</div>
          <div class="text-caption">活跃单元</div>
        </v-card>
      </v-col>
      <v-col cols="6" sm="3">
        <v-card class="rounded-lg text-center pa-4" variant="tonal" color="info">
          <div class="text-h5 font-weight-bold">{{ project.unit_stats.completed_count }}</div>
          <div class="text-caption">已完成</div>
        </v-card>
      </v-col>
      <v-col cols="6" sm="3">
        <v-card class="rounded-lg text-center pa-4" variant="tonal">
          <div class="text-h5 font-weight-bold">{{ project.unit_stats.total_count }}</div>
          <div class="text-caption">总计</div>
        </v-card>
      </v-col>
    </v-row>

    <v-card v-if="project.description" class="rounded-lg mb-4">
      <v-card-title>项目介绍</v-card-title>
      <v-divider />
      <v-card-text class="md-preview-wrap">
        <MdPreview :model-value="project.description" preview-only />
      </v-card-text>
    </v-card>

    <!-- 关联单元管理 -->
    <v-card class="rounded-lg">
      <v-card-title class="d-flex align-center">
        关联单元
        <v-chip size="small" class="ml-2" variant="tonal">{{ units.length }}</v-chip>
        <v-spacer />
        <v-btn
          size="small"
          color="primary"
          variant="tonal"
          prepend-icon="mdi-link-plus"
          @click="openAddDialog"
        >
          添加单元
        </v-btn>
      </v-card-title>
      <v-divider />

      <v-table v-if="units.length > 0">
        <thead>
          <tr>
            <th>标题</th>
            <th>类型</th>
            <th>状态</th>
            <th>优先级</th>
            <th>时间 / 进度</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="unit in units" :key="unit.id">
            <td>
              <router-link :to="`/units/${unit.id}`" class="text-decoration-none font-weight-medium">
                {{ unit.title }}
              </router-link>
            </td>
            <td>
              <v-chip size="x-small" variant="tonal">{{ getUnitTypeLabel(unit.type) }}</v-chip>
            </td>
            <td>
              <v-chip size="x-small" :color="getStatusColor(unit.status)" variant="flat">
                {{ getStatusLabel(unit.status) }}
              </v-chip>
            </td>
            <td>
              <v-chip size="x-small" :color="getPriorityColor(unit.priority)" variant="tonal">
                {{ getPriorityLabel(unit.priority) }}
              </v-chip>
            </td>
            <td class="text-body-2 text-medium-emphasis">
              <span v-if="unit.type === 'time_countdown'">
                {{ unit.remaining_seconds && unit.remaining_seconds > 0 ? formatDuration(unit.remaining_seconds) : '已超期' }}
              </span>
              <span v-else-if="unit.type === 'time_countup' && unit.elapsed_seconds">
                {{ formatDuration(unit.elapsed_seconds) }}
              </span>
              <span v-else-if="unit.type === 'count_countdown'">
                {{ unit.current_value || 0 }} / {{ unit.target_value }} {{ unit.unit_label }}
              </span>
              <span v-else>{{ unit.current_value || 0 }} {{ unit.unit_label }}</span>
            </td>
            <td class="text-right">
              <v-btn
                icon="mdi-link-off"
                size="small"
                variant="text"
                color="error"
                :loading="removingId === unit.id"
                title="从项目中移除"
                @click="removeUnit(unit)"
              />
            </td>
          </tr>
        </tbody>
      </v-table>

      <v-card-text v-else class="text-center text-medium-emphasis py-8">
        <v-icon size="48" color="grey-lighten-1" class="mb-3">mdi-link-off</v-icon>
        <div>该项目暂无关联单元</div>
        <v-btn color="primary" variant="tonal" class="mt-3" prepend-icon="mdi-link-plus" @click="openAddDialog">
          添加单元
        </v-btn>
      </v-card-text>
    </v-card>

    <!-- 添加单元对话框 -->
    <v-dialog v-model="showAddDialog" max-width="720" scrollable>
      <v-card class="rounded-lg">
        <v-card-title class="d-flex align-center">
          选择要加入「{{ project.title }}」的单元
          <v-spacer />
          <v-btn icon="mdi-close" variant="text" size="small" @click="showAddDialog = false" />
        </v-card-title>
        <v-divider />
        <v-card-text style="max-height: 480px; overflow-y: auto">
          <v-text-field
            v-model="addSearch"
            placeholder="搜索单元标题..."
            prepend-inner-icon="mdi-magnify"
            density="compact"
            hide-details
            clearable
            class="mb-3"
          />

          <div v-if="loadingAll" class="text-center py-6">
            <v-progress-circular indeterminate size="32" />
          </div>
          <div v-else-if="filteredAvailableUnits.length === 0" class="text-center text-medium-emphasis py-6">
            {{ addSearch ? '无匹配结果' : '所有单元均已在此项目中' }}
          </div>
          <v-list v-else select-strategy="independent" v-model:selected="selectedUnitIds" lines="two">
            <v-list-item
              v-for="unit in filteredAvailableUnits"
              :key="unit.id"
              :value="unit.id"
              rounded="lg"
              class="mb-1"
            >
              <template #prepend="{ isSelected }">
                <v-checkbox-btn :model-value="isSelected" color="primary" />
              </template>
              <v-list-item-title class="font-weight-medium">{{ unit.title }}</v-list-item-title>
              <v-list-item-subtitle>
                <v-chip size="x-small" variant="tonal" class="mr-1">{{ getUnitTypeLabel(unit.type) }}</v-chip>
                <v-chip size="x-small" :color="getStatusColor(unit.status)" variant="flat">{{ getStatusLabel(unit.status) }}</v-chip>
                <span v-if="unit.project_id" class="ml-2 text-caption text-warning">
                  当前所属：{{ getProjectTitle(unit.project_id) }}
                </span>
              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-divider />
        <v-card-actions>
          <span class="text-caption text-medium-emphasis ml-2">
            已选 {{ selectedUnitIds.length }} 个
          </span>
          <v-spacer />
          <v-btn variant="text" @click="showAddDialog = false">取消</v-btn>
          <v-btn
            color="primary"
            :loading="adding"
            :disabled="selectedUnitIds.length === 0"
            @click="confirmAddUnits"
          >
            确认添加
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 移除确认对话框 -->
    <v-dialog v-model="showRemoveConfirm" max-width="420" persistent>
      <v-card class="rounded-lg">
        <v-card-title>移除单元</v-card-title>
        <v-card-text>
          确定要将「<strong>{{ removingUnit?.title }}</strong>」从本项目中移除吗？该单元本身不会被删除。
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showRemoveConfirm = false">取消</v-btn>
          <v-btn color="error" :loading="!!removingId" @click="confirmRemove">确认移除</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
  <div v-else class="text-center py-12"><v-progress-circular indeterminate /></div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import type { Project, Unit } from '@/types'
import { projectApi } from '@/api/projects'
import { unitApi } from '@/api/units'
import {
  getStatusColor, getStatusLabel, getPriorityColor, getPriorityLabel,
  getUnitTypeLabel, formatDuration,
} from '@/utils/time'

const route = useRoute()
const project = ref<Project | null>(null)
const units = ref<Unit[]>([])
const allUnits = ref<Unit[]>([])
const loadingAll = ref(false)

// 添加单元
const showAddDialog = ref(false)
const addSearch = ref('')
const selectedUnitIds = ref<string[]>([])
const adding = ref(false)

// 移除单元
const showRemoveConfirm = ref(false)
const removingUnit = ref<Unit | null>(null)
const removingId = ref<string | null>(null)

async function fetchProject() {
  const resp = await projectApi.getById(route.params.id as string)
  project.value = resp.data
}

async function fetchUnits() {
  const resp = await projectApi.getUnits(route.params.id as string, { page_size: 200 })
  units.value = resp.data || []
}

async function fetchAllUnits() {
  loadingAll.value = true
  try {
    const resp = await unitApi.list({ page_size: 200 })
    allUnits.value = resp.data || []
  } finally {
    loadingAll.value = false
  }
}

// 项目 id → 标题的快速查找
function getProjectTitle(projectId: string): string {
  // 当前项目
  if (project.value && project.value.id === projectId) return project.value.title
  return projectId
}

// 已在当前项目中的 unit id 集合
const currentProjectUnitIds = computed(() => new Set(units.value.map(u => u.id)))

// 可加入当前项目的单元：排除已经在本项目中的
const availableUnits = computed(() =>
  allUnits.value.filter(u => !currentProjectUnitIds.value.has(u.id))
)

const filteredAvailableUnits = computed(() => {
  if (!addSearch.value) return availableUnits.value
  const q = addSearch.value.toLowerCase()
  return availableUnits.value.filter(u => u.title.toLowerCase().includes(q))
})

async function openAddDialog() {
  selectedUnitIds.value = []
  addSearch.value = ''
  showAddDialog.value = true
  await fetchAllUnits()
}

async function confirmAddUnits() {
  if (selectedUnitIds.value.length === 0) return
  adding.value = true
  try {
    await Promise.all(
      selectedUnitIds.value.map(id =>
        unitApi.assignToProject(id, route.params.id as string)
      )
    )
    showAddDialog.value = false
    await Promise.all([fetchProject(), fetchUnits()])
  } finally {
    adding.value = false
  }
}

function removeUnit(unit: Unit) {
  removingUnit.value = unit
  showRemoveConfirm.value = true
}

async function confirmRemove() {
  if (!removingUnit.value) return
  removingId.value = removingUnit.value.id
  try {
    await unitApi.removeFromProject(removingUnit.value.id)
    showRemoveConfirm.value = false
    await Promise.all([fetchProject(), fetchUnits()])
  } finally {
    removingId.value = null
    removingUnit.value = null
  }
}

onMounted(async () => {
  await fetchProject()
  fetchUnits()
})
</script>

<style>
.md-preview-wrap .md-editor-preview-wrapper {
  padding: 0;
  background: transparent;
}
.md-preview-wrap .md-editor-preview {
  font-size: 0.875rem;
  line-height: 1.6;
  background: transparent;
}
</style>
