import { authState } from '../composables/useAuth'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api'

interface RequestOptions extends RequestInit {
  skipAuthRedirect?: boolean
}

export async function apiFetch<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers = new Headers(options.headers ?? {})
  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }

  if (authState.token) {
    headers.set('Authorization', `Bearer ${authState.token}`)
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  })

  if (response.status === 401 && !options.skipAuthRedirect) {
    const loginPath = window.location.pathname.startsWith('/admin') ? '/admin/login' : '/student/login'
    authState.clear()
    window.location.href = loginPath
    throw new Error('登录已失效，请重新登录')
  }

  const data = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new Error(data.error ?? '请求失败')
  }

  return data as T
}
