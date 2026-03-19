<script setup lang="ts">
import { onMounted, shallowRef } from 'vue'

import StudentShell from '../components/layout/StudentShell.vue'
import ConsumptionTable from '../components/student/ConsumptionTable.vue'
import ProfileCard from '../components/student/ProfileCard.vue'
import { authState } from '../composables/useAuth'
import { apiFetch } from '../lib/api'
import type { MealRecord, StudentProfile } from '../types'

const profile = shallowRef<StudentProfile | null>(null)
const consumptions = shallowRef<MealRecord[]>([])
const loading = shallowRef(false)
const errorMessage = shallowRef('')
const copyMessage = shallowRef('')

async function loadStudentData() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [profileData, consumptionData] = await Promise.all([
      apiFetch<{ profile: StudentProfile }>('/student/profile'),
      apiFetch<{ items: MealRecord[] }>('/student/consumptions'),
    ])
    profile.value = profileData.profile
    consumptions.value = consumptionData.items
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '查询失败'
  } finally {
    loading.value = false
  }
}

async function copyToken() {
  errorMessage.value = ''
  copyMessage.value = ''
  if (!authState.token) {
    errorMessage.value = '当前没有可复制的 JWT'
    return
  }

  try {
    await navigator.clipboard.writeText(authState.token)
    copyMessage.value = '已复制'
    window.setTimeout(() => {
      copyMessage.value = ''
    }, 1800)
  } catch {
    errorMessage.value = '复制失败，请手动复制'
  }
}

onMounted(loadStudentData)
</script>

<template>
  <StudentShell>
    <section class="page-section">
      <div class="page-header">
        <div>
          <p class="eyebrow">Student Center</p>
          <h1>个人中心与消费流水</h1>
          <p class="page-description">
            学生端不再允许手工指定查询对象，后端会直接信任 JWT 中的学号作为查询凭证；因此一旦伪造 JWT 中的学号，就能读取其他学生资料与消费记录。
          </p>
        </div>
      </div>

      <article class="card status-card">
        <div class="status-header">
          <div>
            <strong>当前 JWT</strong>
            <p class="muted status-description">可直接复制后用于本地伪造与调试。</p>
          </div>
          <div class="status-actions">
            <button type="button" class="ghost-button" :disabled="!authState.token" @click="copyToken">
              {{ copyMessage || '复制 JWT' }}
            </button>
            <button type="button" :disabled="loading" @click="loadStudentData">
              {{ loading ? '查询中...' : '重新加载资料与流水' }}
            </button>
          </div>
        </div>
        <code class="token-box">{{ authState.token || '当前未获取到 JWT' }}</code>
        <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
      </article>

      <div class="content-grid">
        <ProfileCard :profile="profile" />
        <ConsumptionTable :items="consumptions" />
      </div>
    </section>
  </StudentShell>
</template>

<style scoped>
.card {
  padding: 1rem 1.1rem;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: var(--color-panel);
}

.status-card {
  display: grid;
  gap: 0.75rem;
}

.status-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.status-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  justify-content: flex-end;
}

.status-description {
  margin: 0.35rem 0 0;
}

.token-box {
  display: block;
  max-width: 100%;
  padding: 0.9rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.03);
  overflow-wrap: anywhere;
  word-break: break-all;
  white-space: pre-wrap;
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

@media (max-width: 960px) {
  .status-header {
    flex-direction: column;
  }

  .status-actions {
    justify-content: flex-start;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
