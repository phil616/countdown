<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="pa-6 rounded-xl" elevation="8">
          <v-card-text class="text-center mb-4">
            <v-icon size="64" color="primary" class="mb-4">mdi-timer-cog-outline</v-icon>
            <h1 class="text-h4 font-weight-bold text-primary">计时器</h1>
            <p class="text-body-2 text-medium-emphasis mt-1">运维计时管理系统</p>
          </v-card-text>

          <v-form @submit.prevent="handleLogin" ref="formRef">
            <v-text-field
              v-model="username"
              label="用户名"
              prepend-inner-icon="mdi-account"
              :rules="[rules.required]"
              :disabled="loading"
            />
            <v-text-field
              v-model="password"
              label="密码"
              prepend-inner-icon="mdi-lock"
              :type="showPassword ? 'text' : 'password'"
              :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
              @click:append-inner="showPassword = !showPassword"
              :rules="[rules.required]"
              :disabled="loading"
            />

            <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mb-4">
              {{ error }}
            </v-alert>

            <v-btn
              type="submit"
              color="primary"
              block
              size="large"
              :loading="loading"
              class="mt-2"
            >
              登录
            </v-btn>
          </v-form>

          <!-- OAuth 登录分割线 + 按钮（仅在后端启用 OAuth 时显示） -->
          <template v-if="oauthEnabled">
            <div class="d-flex align-center my-4">
              <v-divider />
              <span class="mx-3 text-caption text-medium-emphasis text-no-wrap">或</span>
              <v-divider />
            </div>

            <v-btn
              variant="outlined"
              block
              size="large"
              prepend-icon="mdi-shield-account-outline"
              :loading="oauthLoading"
              @click="handleOAuthLogin"
            >
              使用 OAuth 单点登录
            </v-btn>
          </template>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'

const router = useRouter()
const auth = useAuthStore()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

const oauthEnabled = ref(false)
const oauthLoading = ref(false)

const rules = {
  required: (v: string) => !!v || '此项为必填',
}

onMounted(async () => {
  try {
    const res = await authApi.oauthConfig()
    oauthEnabled.value = res.data?.enabled ?? false
  } catch {
    oauthEnabled.value = false
  }
})

async function handleLogin() {
  if (!username.value || !password.value) return
  loading.value = true
  error.value = ''
  try {
    await auth.login(username.value, password.value)
    router.push('/dashboard')
  } catch (e: any) {
    error.value = e.response?.data?.message || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}

function handleOAuthLogin() {
  oauthLoading.value = true
  // 直接跳转到后端 OAuth 授权接口，后端负责重定向到身份提供商
  window.location.href = '/api/v1/auth/oauth/login'
}
</script>
