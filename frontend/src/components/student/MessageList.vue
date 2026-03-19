<script setup lang="ts">
import type { MessageItem } from '../../types'

defineProps<{
  items: MessageItem[]
}>()
</script>

<template>
  <article class="card">
    <div class="section-header">
      <h2>留言板</h2>
      <span class="muted">{{ items.length }} 条</span>
    </div>

    <div class="message-list">
      <article v-for="item in items" :key="item.id" class="message-item">
        <div class="section-header">
          <strong>{{ item.studentName }}</strong>
          <span class="muted">{{ item.createdAt }}</span>
        </div>
        <p>{{ item.content }}</p>
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

.message-item p {
  margin-bottom: 0;
  white-space: pre-wrap;
}
</style>

