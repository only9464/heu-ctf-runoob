<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { useAuth } from '../../composables/useAuth'

const route = useRoute()
const router = useRouter()
const { user, logout } = useAuth()

const navItems = computed(() => [
  { to: '/admin', label: '后台总览' },
  { to: '/admin/students', label: '学生资料' },
  { to: '/admin/messages', label: '留言列表' },
])

async function handleLogout() {
  await logout()
  await router.push('/admin/login')
}
</script>

<template>
  <div class="portal-shell">
    <aside class="portal-sidebar portal-sidebar--admin">
      <div class="brand-card">
        <p class="brand-tag">Admin Console</p>
        <h1>校园餐厅后台</h1>
        <p class="brand-description">系统管理员/教师用于查看学生资料、留言与重置环境。</p>
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
        <p class="section-label">当前管理员</p>
        <h2>{{ user.displayName }}</h2>
        <p>{{ user.username }}</p>
        <p>角色：{{ user.role }}</p>
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
}

.portal-sidebar--admin {
  background: rgba(27, 13, 28, 0.96);
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
  color: #f59e0b;
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
  border-color: rgba(245, 158, 11, 0.65);
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

