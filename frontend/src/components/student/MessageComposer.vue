<script setup lang="ts">
import { reactive } from 'vue'

const props = defineProps<{
  loading: boolean
}>()

const emit = defineEmits<{
  submit: [content: string]
}>()

const form = reactive({
  content: '',
})

function handleSubmit() {
  emit('submit', form.content)
  form.content = ''
}
</script>

<template>
  <article class="card">
    <div class="section-header">
      <h2>提交留言</h2>
      <span class="muted">管理员后台将不安全渲染</span>
    </div>

    <form class="stack" @submit.prevent="handleSubmit">
      <label class="field">
        <span>留言内容</span>
        <textarea v-model="form.content" rows="5" placeholder="输入对食堂的建议或意见"></textarea>
      </label>
      <button type="submit" :disabled="props.loading || !form.content.trim()">
        {{ props.loading ? '提交中...' : '提交留言' }}
      </button>
    </form>
  </article>
</template>

<style scoped>
.card {
  padding: 1rem 1.1rem;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: var(--color-panel);
}

.stack,
.field {
  display: grid;
  gap: 0.75rem;
}
</style>

