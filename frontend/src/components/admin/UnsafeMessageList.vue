<script setup lang="ts">
import type { MessageItem } from '../../types'

defineProps<{
  items: MessageItem[]
  deletingId: number | null
}>()

const emit = defineEmits<{
  delete: [id: number]
}>()
</script>

<template>
  <article class="card">
    <div class="section-header">
      <h2>后台留言列表</h2>
      <span class="muted">管理员查看时触发存储型 XSS</span>
    </div>

    <div class="message-list">
      <article v-for="item in items" :key="item.id" class="message-item">
        <div class="section-header">
          <div>
            <strong>{{ item.studentName }}</strong>
            <span class="muted muted--block">{{ item.createdAt }}</span>
          </div>
          <button
            type="button"
            class="ghost-button"
            :disabled="deletingId === item.id"
            @click="emit('delete', item.id)"
          >
            {{ deletingId === item.id ? '删除中...' : '删除留言' }}
          </button>
        </div>
        <!-- INTENTIONALLY UNSAFE: teaching environment keeps stored XSS here. -->
        <div class="unsafe-html" v-html="item.content"></div>
      </article>
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

.message-list {
  display: grid;
  gap: 0.85rem;
}

.message-item {
  padding: 0.9rem;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.03);
}

.muted--block {
  display: block;
  margin-top: 0.2rem;
}

.unsafe-html {
  margin-top: 0.75rem;
  color: var(--color-text);
}
</style>
