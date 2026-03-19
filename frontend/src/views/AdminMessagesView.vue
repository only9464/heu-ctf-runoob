<script setup lang="ts">
import { onMounted, shallowRef } from 'vue'

import UnsafeMessageList from '../components/admin/UnsafeMessageList.vue'
import AdminShell from '../components/layout/AdminShell.vue'
import { apiFetch } from '../lib/api'
import type { MessageItem } from '../types'

const messages = shallowRef<MessageItem[]>([])
const deletingId = shallowRef<number | null>(null)
const errorMessage = shallowRef('')
const successMessage = shallowRef('')

async function loadMessages() {
  errorMessage.value = ''
  try {
    const data = await apiFetch<{ items: MessageItem[] }>('/admin/messages')
    messages.value = data.items
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '加载留言列表失败'
  }
}

async function deleteMessage(id: number) {
  deletingId.value = id
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const data = await apiFetch<{ message: string }>(`/admin/messages/${id}`, {
      method: 'DELETE',
    })
    messages.value = messages.value.filter((item) => item.id !== id)
    successMessage.value = `${data.message}。如果刚才是 XSS payload，可借此清理后台演示环境。`
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '删除留言失败'
  } finally {
    deletingId.value = null
  }
}

onMounted(loadMessages)
</script>

<template>
  <AdminShell>
    <section class="page-section">
      <div class="page-header">
        <div>
          <p class="eyebrow eyebrow--warning">Stored XSS</p>
          <h1>后台留言列表</h1>
          <p class="page-description">
            管理员端会直接渲染学生留言内容，因此学生提交的恶意 HTML/脚本会在这里执行；管理员也可直接删除留言。
          </p>
        </div>
      </div>

      <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
      <p v-if="successMessage" class="info-text">{{ successMessage }}</p>
      <UnsafeMessageList :items="messages" :deleting-id="deletingId" @delete="deleteMessage" />
    </section>
  </AdminShell>
</template>
