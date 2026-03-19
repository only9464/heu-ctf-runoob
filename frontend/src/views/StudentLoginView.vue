<script setup lang="ts">
import { reactive, shallowRef } from 'vue'
import { useRouter } from 'vue-router'

import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { loginStudent, loading, errorMessage } = useAuth()

const form = reactive({
  username: 'student001',
  password: 'Meal2026#001',
})
const showPassword = shallowRef(false)

const demoAccounts = [
  { username: 'student001', password: 'Meal2026#001' },
  { username: 'student002', password: 'Meal2026#002' },
]

async function handleSubmit() {
  try {
    await loginStudent(form.username, form.password)
    await router.push('/student')
  } catch {
    // handled in composable
  }
}

function fillDemoAccount(username: string, password: string) {
  form.username = username
  form.password = password
}
</script>

<template>
  <div class="login-page">
    <section class="login-hero panel">
      <p class="eyebrow">Student Portal</p>
      <h1>校园餐厅学生入口</h1>
      <p class="lead">
        这是食堂学生端：可查看个人余额、充值记录、消费流水和留言反馈。
      </p>
      <ul class="hero-list">
        <li>学生登录接口保留 SQL 注入</li>
        <li>个人资料与消费记录存在水平越权</li>
        <li>充值逻辑保留非法金额入账漏洞</li>
        <li>留言会在管理员后台触发存储型 XSS</li>
      </ul>
    </section>

    <section class="login-form panel">
      <div class="section-header">
        <div>
          <p class="eyebrow">Student Login</p>
          <h2>学生登录</h2>
        </div>
        <RouterLink class="text-link" to="/admin/login">前往管理员入口</RouterLink>
      </div>

      <form class="stack" @submit.prevent="handleSubmit">
        <label class="field">
          <span>用户名</span>
          <input v-model="form.username" autocomplete="username" />
        </label>

        <label class="field">
          <span>密码</span>
          <div class="password-row">
            <input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              autocomplete="current-password"
            />
            <button class="ghost-button password-toggle" type="button" @click="showPassword = !showPassword">
              {{ showPassword ? '隐藏密码' : '显示密码' }}
            </button>
          </div>
        </label>

        <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
        <button type="submit" :disabled="loading">{{ loading ? '登录中...' : '登录学生端' }}</button>
      </form>

      <div class="demo-grid">
        <button
          v-for="account in demoAccounts"
          :key="account.username"
          class="ghost-button"
          type="button"
          @click="fillDemoAccount(account.username, account.password)"
        >
          使用 {{ account.username }}
        </button>
      </div>

      <p class="muted">
        系统默认存在 6 个学生账号：student001 ~ student006。
      </p>
    </section>
  </div>
</template>

<style scoped>
.login-page {
  display: grid;
  grid-template-columns: 1.2fr minmax(320px, 430px);
  min-height: 100vh;
  gap: 2rem;
  padding: 2rem;
  background:
    radial-gradient(circle at top left, rgba(16, 185, 129, 0.18), transparent 34%),
    radial-gradient(circle at bottom right, rgba(59, 130, 246, 0.12), transparent 34%),
    var(--color-bg);
}

.panel {
  padding: 2rem;
  border: 1px solid var(--color-border);
  border-radius: 28px;
  background: rgba(10, 16, 30, 0.94);
}

.lead {
  color: var(--color-muted);
}

.hero-list,
.stack,
.demo-grid {
  display: grid;
  gap: 0.8rem;
}

.hero-list {
  margin-top: 1.5rem;
  padding-left: 1.2rem;
}

.field {
  display: grid;
  gap: 0.4rem;
}

.password-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 0.75rem;
  align-items: center;
}

.password-toggle {
  white-space: nowrap;
}

@media (max-width: 900px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .password-row {
    grid-template-columns: 1fr;
  }
}
</style>
