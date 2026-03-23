<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4" class="text-center">
        <v-card class="pa-8 rounded-xl" elevation="8">
          <!-- 加载中 -->
          <template v-if="status === 'loading'">
            <v-progress-circular indeterminate color="primary" size="64" class="mb-6" />
            <p class="text-body-1 text-medium-emphasis">正在验证 OAuth 身份，请稍候…</p>
          </template>

          <!-- 成功 -->
          <template v-else-if="status === 'success'">
            <v-icon size="64" color="success" class="mb-4">mdi-check-circle-outline</v-icon>
            <h2 class="text-h5 font-weight-bold mb-2">登录成功</h2>
            <p class="text-body-2 text-medium-emphasis">正在跳转到控制台…</p>
          </template>

          <!-- 失败 -->
          <template v-else>
            <v-icon size="64" color="error" class="mb-4">mdi-alert-circle-outline</v-icon>
            <h2 class="text-h5 font-weight-bold mb-3">OAuth 登录失败</h2>
            <v-alert type="error" variant="tonal" class="mb-6 text-left">
              {{ errorMsg }}
            </v-alert>
            <v-btn color="primary" variant="outlined" prepend-icon="mdi-arrow-left" @click="goLogin">
              返回登录页
            </v-btn>
          </template>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

type Status = 'loading' | 'success' | 'error'
const status = ref<Status>('loading')
const errorMsg = ref('')

onMounted(async () => {
  const token = route.query.token as string | undefined
  const error = route.query.error as string | undefined

  if (error) {
    status.value = 'error'
    errorMsg.value = decodeURIComponent(error)
    return
  }

  if (!token) {
    status.value = 'error'
    errorMsg.value = '未收到有效的登录凭证，请重试'
    return
  }

  try {
    await auth.loginWithToken(token)
    status.value = 'success'
    setTimeout(() => router.replace('/dashboard'), 1000)
  } catch (e: any) {
    status.value = 'error'
    errorMsg.value = e.message || 'OAuth 会话建立失败，请重试'
  }
})

function goLogin() {
  router.replace('/login')
}
</script>
