<script setup lang="ts">
import { onMounted, shallowRef } from 'vue'
import { useRouter } from 'vue-router'

import AdminShell from '../components/layout/AdminShell.vue'
import { useAuth } from '../composables/useAuth'
import { apiFetch } from '../lib/api'
import type { AdminSummary } from '../types'

const router = useRouter()
const { logout } = useAuth()

const summary = shallowRef<AdminSummary | null>(null)
const loading = shallowRef(false)
const errorMessage = shallowRef('')
const successMessage = shallowRef('')

async function loadSummary() {
  errorMessage.value = ''
  try {
    const data = await apiFetch<{ summary: AdminSummary }>('/admin/summary')
    summary.value = data.summary
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '加载后台总览失败'
  }
}

async function resetDemo() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const data = await apiFetch<{ message: string }>('/admin/reset-demo', {
      method: 'POST',
    })
    successMessage.value = `${data.message}，当前会话将退出。`
    await logout()
    await router.push('/admin/login')
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '重置失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadSummary)
</script>

<template>
  <AdminShell>
    <section class="page-section">
      <div class="page-header">
        <div>
          <p class="eyebrow eyebrow--warning">Admin Overview</p>
          <h1>后台总览</h1>
          <p class="page-description">
            管理员账号弱口令为唯一认证漏洞点；后台负责查看学生数据和留言列表。
          </p>
        </div>
      </div>

      <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
      <p v-if="successMessage" class="info-text">{{ successMessage }}</p>

      <div v-if="summary" class="stat-grid">
        <article class="stat-card">
          <p class="muted">学生账号</p>
          <strong>{{ summary.studentCount }}</strong>
        </article>
        <article class="stat-card">
          <p class="muted">留言总数</p>
          <strong>{{ summary.messageCount }}</strong>
        </article>
        <article class="stat-card">
          <p class="muted">充值记录</p>
          <strong>{{ summary.rechargeCount }}</strong>
        </article>
      </div>

      <article class="card">
        <div class="section-header">
          <div>
            <h2>演示环境控制</h2>
            <p class="muted">默认管理员用户名固定为 {{ summary?.adminUsername ?? 'admin' }}</p>
          </div>
          <button type="button" class="ghost-button" :disabled="loading" @click="resetDemo">
            {{ loading ? '重置中...' : '重置演示环境' }}
          </button>
        </div>
      </article>
    </section>
  </AdminShell>
</template>

<style scoped>
.stat-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.stat-card,
.card {
  padding: 1rem 1.1rem;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: var(--color-panel);
}

.stat-card strong {
  display: block;
  margin-top: 0.5rem;
  font-size: 2rem;
}

@media (max-width: 960px) {
  .stat-grid {
    grid-template-columns: 1fr;
  }
}
</style>

