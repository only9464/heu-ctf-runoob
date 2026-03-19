<script setup lang="ts">
import { reactive, shallowRef } from 'vue'
import { RouterLink, useRouter } from 'vue-router'

import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { loginAdmin, loading, errorMessage } = useAuth()

const form = reactive({
  username: 'admin',
  password: '',
})
const showPassword = shallowRef(false)

async function handleSubmit() {
  try {
    await loginAdmin(form.username, form.password)
    await router.push('/admin')
  } catch {
    // handled in composable
  }
}
</script>

<template>
  <div class="login-page login-page--admin">
    <section class="login-hero panel">
      <p class="eyebrow eyebrow--warning">Admin Console</p>
      <h1>校园餐厅后台入口</h1>
      <p class="lead">
        后台用于查看学生资料、留言列表和重置演示环境。系统中仅有一个管理员账号。
      </p>
      <ul class="hero-list">
        <li>管理员登录接口本身无 SQL 注入</li>
        <li>唯一管理员账号保留弱口令</li>
        <li>后台留言列表会触发学生提交的存储型 XSS</li>
      </ul>
    </section>

    <section class="login-form panel">
      <div class="section-header">
        <div>
          <p class="eyebrow eyebrow--warning">Admin Login</p>
          <h2>系统管理员/教师登录</h2>
        </div>
        <RouterLink class="text-link" to="/student/login">返回学生入口</RouterLink>
      </div>

      <form class="stack" @submit.prevent="handleSubmit">
        <label class="field">
          <span>管理员用户名</span>
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
        <button type="submit" :disabled="loading">{{ loading ? '登录中...' : '登录后台' }}</button>
      </form>

      <p class="muted">演示环境中唯一管理员账号用户名固定为 <code>admin</code>。</p>
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
    radial-gradient(circle at top left, rgba(245, 158, 11, 0.18), transparent 34%),
    radial-gradient(circle at bottom right, rgba(244, 63, 94, 0.12), transparent 34%),
    var(--color-bg);
}

.login-page--admin .panel {
  background: rgba(24, 14, 29, 0.94);
}

.panel {
  padding: 2rem;
  border: 1px solid var(--color-border);
  border-radius: 28px;
}

.lead {
  color: var(--color-muted);
}

.hero-list,
.stack {
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
