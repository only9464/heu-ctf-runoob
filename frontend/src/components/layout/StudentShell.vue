<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { useAuth } from '../../composables/useAuth'

const route = useRoute()
const router = useRouter()
const { user, logout } = useAuth()

const navItems = computed(() => [
  { to: '/student', label: '个人中心' },
  { to: '/student/recharge', label: '余额充值' },
  { to: '/student/messages', label: '学生留言' },
])

async function handleLogout() {
  await logout()
  await router.push('/student/login')
}
</script>

<template>
  <div class="portal-shell">
    <aside class="portal-sidebar">
      <div class="brand-card">
        <p class="brand-tag">Campus Canteen</p>
        <h1>校园餐厅学生端</h1>
        <p class="brand-description">余额充值、消费记录、留言反馈都从这里进入。</p>
      </div>

      <nav class="nav-list">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="nav-link"
          :class="{ 'nav-link--active': route.path === item.to }"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <section v-if="user" class="user-card">
        <p class="section-label">当前学生</p>
        <h2>{{ user.displayName }}</h2>
        <p>{{ user.username }}</p>
        <p>学号：{{ user.studentNo }}</p>
        <button class="ghost-button" type="button" @click="handleLogout">退出登录</button>
      </section>
    </aside>

    <main class="portal-main">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.portal-shell {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  min-height: 100vh;
  background: var(--color-bg);
}

.portal-sidebar {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.5rem;
  border-right: 1px solid var(--color-border);
  background: rgba(11, 19, 39, 0.96);
}

.brand-card,
.user-card {
  padding: 1rem;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: var(--color-panel);
}

.brand-tag,
.section-label {
  margin: 0 0 0.6rem;
  color: var(--color-accent);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 0.8rem;
  font-weight: 700;
}

.brand-card h1,
.user-card h2 {
  margin: 0 0 0.75rem;
}

.brand-description,
.user-card p {
  margin: 0.25rem 0;
  color: var(--color-muted);
}

.nav-list {
  display: grid;
  gap: 0.75rem;
}

.nav-link {
  padding: 0.95rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: 16px;
  color: var(--color-text);
  text-decoration: none;
  background: rgba(255, 255, 255, 0.03);
}

.nav-link--active,
.nav-link:hover {
  border-color: rgba(96, 165, 250, 0.65);
}

.portal-main {
  padding: 2rem;
}

@media (max-width: 960px) {
  .portal-shell {
    grid-template-columns: 1fr;
  }

  .portal-sidebar {
    border-right: none;
    border-bottom: 1px solid var(--color-border);
  }
}
</style>
