<script setup lang="ts">
import { reactive } from 'vue'

const props = defineProps<{
  loading: boolean
}>()

const emit = defineEmits<{
  submit: [amount: number]
}>()

const form = reactive({
  amount: 10,
})

function handleSubmit() {
  emit('submit', form.amount)
}
</script>

<template>
  <article class="card">
    <div class="section-header">
      <h2>银行卡转入校园卡</h2>
      <span class="muted">故意保留非法金额漏洞</span>
    </div>

    <form class="stack" @submit.prevent="handleSubmit">
      <label class="field">
        <span>转入金额（系统本应只接受正整数元）</span>
        <input v-model.number="form.amount" type="number" step="0.001" />
      </label>
      <button type="submit" :disabled="props.loading">
        {{ props.loading ? '转入中...' : '提交转入' }}
      </button>
    </form>

    <p class="muted">可测试：0.001、-1、0.5</p>
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
