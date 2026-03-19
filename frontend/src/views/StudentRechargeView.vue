<script setup lang="ts">
import { onMounted, shallowRef } from 'vue'

import StudentShell from '../components/layout/StudentShell.vue'
import RechargeForm from '../components/student/RechargeForm.vue'
import RechargeHistory from '../components/student/RechargeHistory.vue'
import { apiFetch } from '../lib/api'
import type { RechargeRecord, StudentProfile } from '../types'

const profile = shallowRef<StudentProfile | null>(null)
const history = shallowRef<RechargeRecord[]>([])
const loading = shallowRef(false)
const errorMessage = shallowRef('')
const successMessage = shallowRef('')

async function loadPage() {
  errorMessage.value = ''
  try {
    const [profileData, historyData] = await Promise.all([
      apiFetch<{ profile: StudentProfile }>('/student/profile'),
      apiFetch<{ items: RechargeRecord[] }>('/student/recharges'),
    ])
    profile.value = profileData.profile
    history.value = historyData.items
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '加载失败'
  }
}

async function submitRecharge(amount: number) {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const data = await apiFetch<{
      message: string
      rawAmount: number
      creditedAmount: number
      campusCardBalance: number
      bankBalance: number
      note: string
    }>('/student/recharge', {
      method: 'POST',
      body: JSON.stringify({ amount }),
    })
    successMessage.value = `${data.message}：原始 ${data.rawAmount}，校园卡入账 ${data.creditedAmount}，校园卡余额 ${data.campusCardBalance.toFixed(3)}，银行卡余额 ${data.bankBalance.toFixed(3)}。${data.note}`
    await loadPage()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '转入失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadPage)
</script>

<template>
  <StudentShell>
    <section class="page-section">
      <div class="page-header">
        <div>
          <p class="eyebrow">Recharge</p>
          <h1>校园卡余额充值</h1>
          <p class="page-description">
            每位学生都绑定了银行卡和校园卡，校内消费只扣校园卡；当前教学环境保留错误的转账入账逻辑。
          </p>
        </div>
      </div>

      <article v-if="profile" class="balance-card">
        <p class="muted">当前学生：{{ profile.displayName }}（{{ profile.studentNo }}）</p>
        <div class="balance-grid">
          <div>
            <span class="label">校园卡</span>
            <strong>{{ profile.campusCardNo }}</strong>
            <p>{{ profile.campusCardBalance.toFixed(3) }} 元</p>
          </div>
          <div>
            <span class="label">银行卡</span>
            <strong>{{ profile.bankCardNo }}</strong>
            <p>{{ profile.bankBalance.toFixed(3) }} 元</p>
          </div>
        </div>
      </article>

      <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
      <p v-if="successMessage" class="info-text">{{ successMessage }}</p>

      <div class="content-grid">
        <RechargeForm :loading="loading" @submit="submitRecharge" />
        <RechargeHistory :items="history" />
      </div>
    </section>
  </StudentShell>
</template>

<style scoped>
.balance-card {
  padding: 1rem 1.1rem;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: var(--color-panel);
}

.balance-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
  margin-top: 0.75rem;
}

.balance-grid strong {
  display: block;
  margin: 0.35rem 0;
  font-size: 1.1rem;
}

.balance-grid p {
  margin: 0;
  font-size: 1.65rem;
  font-weight: 700;
}

.label {
  color: var(--color-muted);
  font-size: 0.85rem;
}

.content-grid {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr);
  gap: 1rem;
}

@media (max-width: 960px) {
  .content-grid,
  .balance-grid {
    grid-template-columns: 1fr;
  }
}
</style>
