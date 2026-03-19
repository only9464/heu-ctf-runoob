import { createRouter, createWebHistory } from 'vue-router'

import { authState } from '../composables/useAuth'
import AdminDashboardView from '../views/AdminDashboardView.vue'
import AdminLoginView from '../views/AdminLoginView.vue'
import AdminMessagesView from '../views/AdminMessagesView.vue'
import AdminStudentsView from '../views/AdminStudentsView.vue'
import StudentDashboardView from '../views/StudentDashboardView.vue'
import StudentLoginView from '../views/StudentLoginView.vue'
import StudentMessagesView from '../views/StudentMessagesView.vue'
import StudentRechargeView from '../views/StudentRechargeView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: () => {
        if (authState.user?.role === 'admin') {
          return '/admin'
        }
        if (authState.user?.role === 'student') {
          return '/student'
        }
        return '/student/login'
      },
    },
    {
      path: '/student/login',
      name: 'student-login',
      component: StudentLoginView,
      meta: { public: true, audience: 'student' },
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: AdminLoginView,
      meta: { public: true, audience: 'admin' },
    },
    {
      path: '/student',
      name: 'student-dashboard',
      component: StudentDashboardView,
      meta: { role: 'student' },
    },
    {
      path: '/student/recharge',
      name: 'student-recharge',
      component: StudentRechargeView,
      meta: { role: 'student' },
    },
    {
      path: '/student/messages',
      name: 'student-messages',
      component: StudentMessagesView,
      meta: { role: 'student' },
    },
    {
      path: '/admin',
      name: 'admin-dashboard',
      component: AdminDashboardView,
      meta: { role: 'admin' },
    },
    {
      path: '/admin/students',
      name: 'admin-students',
      component: AdminStudentsView,
      meta: { role: 'admin' },
    },
    {
      path: '/admin/messages',
      name: 'admin-messages',
      component: AdminMessagesView,
      meta: { role: 'admin' },
    },
  ],
})

router.beforeEach((to) => {
  const currentUser = authState.user

  if (to.meta.public) {
    if (!currentUser) {
      return true
    }

    return currentUser.role === 'admin' ? '/admin' : '/student'
  }

  if (!currentUser) {
    return to.meta.role === 'admin' ? '/admin/login' : '/student/login'
  }

  if (to.meta.role === 'admin' && currentUser.role !== 'admin') {
    return '/student'
  }

  if (to.meta.role === 'student' && currentUser.role !== 'student') {
    return '/admin'
  }

  return true
})

export default router
