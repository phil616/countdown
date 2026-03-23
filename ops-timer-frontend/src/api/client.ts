import axios from 'axios'
import type { ApiResponse } from '@/types'
import router from '@/router'

const client = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export async function get<T>(url: string, params?: any): Promise<ApiResponse<T>> {
  const { data } = await client.get<ApiResponse<T>>(url, { params })
  return data
}

export async function post<T>(url: string, body?: any): Promise<ApiResponse<T>> {
  const { data } = await client.post<ApiResponse<T>>(url, body)
  return data
}

export async function put<T>(url: string, body?: any): Promise<ApiResponse<T>> {
  const { data } = await client.put<ApiResponse<T>>(url, body)
  return data
}

export async function patch<T>(url: string, body?: any): Promise<ApiResponse<T>> {
  const { data } = await client.patch<ApiResponse<T>>(url, body)
  return data
}

export async function del<T>(url: string): Promise<any> {
  const { data } = await client.delete<ApiResponse<T>>(url)
  return data
}

export default client
