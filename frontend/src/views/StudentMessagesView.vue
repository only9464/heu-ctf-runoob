<script setup lang="ts">
import { onMounted, shallowRef } from 'vue'

import StudentShell from '../components/layout/StudentShell.vue'
import MessageComposer from '../components/student/MessageComposer.vue'
import MessageList from '../components/student/MessageList.vue'
import { apiFetch } from '../lib/api'
import type { MessageItem } from '../types'

const messages = shallowRef<MessageItem[]>([])
const loading = shallowRef(false)
const errorMessage = shallowRef('')
const successMessage = shallowRef('')

async function loadMessages() {
  errorMessage.value = ''
  try {
    const data = await apiFetch<{ items: MessageItem[] }>('/student/messages')
    messages.value = data.items
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '加载留言失败'
  }
}

async function submitMessage(content: string) {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const data = await apiFetch<{ message: string }>('/student/messages', {
      method: 'POST',
      body: JSON.stringify({ content }),
    })
    successMessage.value = `${data.message}。管理员查看后台留言列表时会直接渲染内容。`
    await loadMessages()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '提交留言失败'
  } finally {
    loading.value = false
  }
}

onMounted(loadMessages)
</script>

<template>
  <StudentShell>
    <section class="page-section">
      <div class="page-header">
        <div>
          <p class="eyebrow">Message Board</p>
          <h1>学生留言板</h1>
          <p class="page-description">
            学生端留言板正常显示文本，但后台会不安全渲染同一份留言内容。
          </p>
        </div>
      </div>

      <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
      <p v-if="successMessage" class="info-text">{{ successMessage }}</p>

      <div class="content-grid">
        <MessageComposer :loading="loading" @submit="submitMessage" />
        <MessageList :items="messages" />
      </div>
    </section>
  </StudentShell>
</template>

<style scoped>
.content-grid {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr);
  gap: 1rem;
}

@media (max-width: 960px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>

