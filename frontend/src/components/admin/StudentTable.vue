<script setup lang="ts">
import type { StudentListItem } from '../../types'

defineProps<{
  items: StudentListItem[]
  selectedId: number | null
}>()

const emit = defineEmits<{
  select: [id: number]
}>()
</script>

<template>
  <article class="card">
    <div class="section-header">
      <h2>学生账号</h2>
      <span class="muted">{{ items.length }} 人</span>
    </div>

    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>姓名</th>
            <th>用户名</th>
            <th>校园卡余额</th>
            <th>银行卡余额</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.id">
            <td>{{ item.id }}</td>
            <td>{{ item.displayName }}</td>
            <td>{{ item.username }}</td>
            <td>{{ item.campusCardBalance.toFixed(3) }}</td>
            <td>{{ item.bankBalance.toFixed(3) }}</td>
            <td>
              <button
                class="ghost-button"
                type="button"
                :class="{ 'ghost-button--active': selectedId === item.id }"
                @click="emit('select', item.id)"
              >
                查看详情
              </button>
            </td>
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

.ghost-button--active {
  border-color: rgba(245, 158, 11, 0.65);
}
</style>
