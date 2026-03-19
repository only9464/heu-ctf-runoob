<script setup lang="ts">
import type { StudentProfile } from '../../types'

defineProps<{
  profile: StudentProfile | null
}>()
</script>

<template>
  <article class="card">
    <div class="section-header">
      <h2>学生资料</h2>
      <span class="muted">JWT 越权目标</span>
    </div>

    <div v-if="profile" class="detail-grid">
      <div>
        <span class="label">姓名</span>
        <strong>{{ profile.displayName }}</strong>
      </div>
      <div>
        <span class="label">用户名</span>
        <strong>{{ profile.username }}</strong>
      </div>
      <div>
        <span class="label">学号</span>
        <strong>{{ profile.studentNo }}</strong>
      </div>
      <div>
        <span class="label">校园卡号</span>
        <strong>{{ profile.campusCardNo }}</strong>
      </div>
      <div>
        <span class="label">校园卡余额</span>
        <strong>{{ profile.campusCardBalance.toFixed(3) }} 元</strong>
      </div>
      <div>
        <span class="label">银行卡号</span>
        <strong>{{ profile.bankCardNo }}</strong>
      </div>
      <div>
        <span class="label">银行卡余额</span>
        <strong>{{ profile.bankBalance.toFixed(3) }} 元</strong>
      </div>
      <div>
        <span class="label">手机号</span>
        <strong>{{ profile.phone }}</strong>
      </div>
      <div>
        <span class="label">宿舍</span>
        <strong>{{ profile.dormitory }}</strong>
      </div>
    </div>
    <p v-else class="muted">暂无资料</p>
  </article>
</template>

<style scoped>
.card {
  padding: 1rem 1.1rem;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: var(--color-panel);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.label {
  display: block;
  margin-bottom: 0.35rem;
  color: var(--color-muted);
  font-size: 0.85rem;
}

@media (max-width: 720px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
