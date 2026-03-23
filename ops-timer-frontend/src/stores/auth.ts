import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<User | null>(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const resp = await authApi.login(username, password)
    token.value = resp.data.token
    user.value = resp.data.user
    localStorage.setItem('token', resp.data.token)
    localStorage.setItem('user', JSON.stringify(resp.data.user))
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // ignore
    }
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function fetchProfile() {
    const resp = await authApi.getProfile()
    user.value = resp.data
    localStorage.setItem('user', JSON.stringify(resp.data))
  }

  /** OAuth 回调后直接用 token 建立会话 */
  async function loginWithToken(jwt: string) {
    token.value = jwt
    localStorage.setItem('token', jwt)
    // 拉取用户信息
    try {
      const resp = await authApi.getProfile()
      user.value = resp.data
      localStorage.setItem('user', JSON.stringify(resp.data))
    } catch {
      // token 可能已失效，清空
      token.value = ''
      localStorage.removeItem('token')
      throw new Error('获取用户信息失败，token 可能无效')
    }
  }

  return { token, user, isLoggedIn, login, loginWithToken, logout, fetchProfile }
})
