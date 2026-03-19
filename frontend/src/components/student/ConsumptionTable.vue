<script setup lang="ts">
import type { MealRecord } from '../../types'

defineProps<{
  items: MealRecord[]
}>()
</script>

<template>
  <article class="card">
    <div class="section-header">
      <h2>消费记录</h2>
      <span class="muted">{{ items.length }} 条</span>
    </div>

    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>窗口</th>
            <th>菜品</th>
            <th>金额</th>
            <th>时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.id">
            <td>{{ item.canteenWindow }}</td>
            <td>{{ item.dishName }}</td>
            <td>{{ item.amount.toFixed(2) }} 元</td>
            <td>{{ item.consumedAt }}</td>
          </tr>
          <tr v-if="!items.length">
            <td colspan="4" class="empty-cell">暂无消费记录</td>
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

