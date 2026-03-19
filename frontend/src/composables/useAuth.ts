import { computed, reactive, shallowRef } from 'vue'

import { apiFetch } from '../lib/api'
import type { AuthUser } from '../types'

const STORAGE_TOKEN_KEY = 'canteen-ctf-token'
const STORAGE_USER_KEY = 'canteen-ctf-user'

const token = shallowRef<string>(localStorage.getItem(STORAGE_TOKEN_KEY) ?? '')
const user = shallowRef<AuthUser | null>(readStoredUser())
const loading = shallowRef(false)
const errorMessage = shallowRef('')

export const authState = reactive({
  get token() {
    return token.value
  },
  get user() {
    return user.value
  },
  clear() {
    token.value = ''
    user.value = null
    errorMessage.value = ''
    localStorage.removeItem(STORAGE_TOKEN_KEY)
    localStorage.removeItem(STORAGE_USER_KEY)
  },
})

function readStoredUser(): AuthUser | null {
  const rawValue = localStorage.getItem(STORAGE_USER_KEY)
  if (!rawValue) {
    return null
  }

  try {
    return JSON.parse(rawValue) as AuthUser
  } catch {
    return null
  }
}

export function useAuth() {
  const isAuthenticated = computed(() => Boolean(token.value))
  const isStudent = computed(() => user.value?.role === 'student')
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function loginStudent(username: string, password: string) {
    await login('/student/login', username, password)
  }

  async function loginAdmin(username: string, password: string) {
    await login('/admin/login', username, password)
  }

  async function login(path: string, username: string, password: string) {
    loading.value = true
    errorMessage.value = ''
    try {
      const data = await apiFetch<{ token: string; user: AuthUser }>(path, {
        method: 'POST',
        body: JSON.stringify({ username, password }),
        skipAuthRedirect: true,
      })

      token.value = data.token
      user.value = data.user
      localStorage.setItem(STORAGE_TOKEN_KEY, data.token)
      localStorage.setItem(STORAGE_USER_KEY, JSON.stringify(data.user))
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : '登录失败'
      throw error
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    if (token.value) {
      try {
        await apiFetch('/auth/logout', {
          method: 'POST',
        })
      } catch {
        // ignore
      }
    }

    authState.clear()
  }

  return {
    user,
    loading,
    errorMessage,
    isAuthenticated,
    isStudent,
    isAdmin,
    loginStudent,
    loginAdmin,
    logout,
  }
}
