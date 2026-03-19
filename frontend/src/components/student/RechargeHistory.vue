<script setup lang="ts">
import type { RechargeRecord } from '../../types'

defineProps<{
  items: RechargeRecord[]
}>()
</script>

<template>
  <article class="card">
    <div class="section-header">
      <h2>充值记录</h2>
      <span class="muted">{{ items.length }} 条</span>
    </div>

    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>原始金额</th>
            <th>实际入账</th>
            <th>备注</th>
            <th>时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.id">
            <td>{{ item.rawAmount.toFixed(3) }}</td>
            <td>{{ item.creditedAmount.toFixed(3) }}</td>
            <td>{{ item.note }}</td>
            <td>{{ item.createdAt }}</td>
          </tr>
          <tr v-if="!items.length">
            <td colspan="4" class="empty-cell">暂无充值记录</td>
          </tr>
        </tbody>
      </table>
    </div>
  </article>
</template>

<style scoped>
.card {
  padding: 1rem 1.1rem;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: var(--color-panel);
}

.empty-cell {
  text-align: center;
  color: var(--color-muted);
}
</style>

