<script setup lang="ts">
import { onMounted, shallowRef } from 'vue'

import StudentTable from '../components/admin/StudentTable.vue'
import AdminShell from '../components/layout/AdminShell.vue'
import { apiFetch } from '../lib/api'
import type { StudentListItem, StudentProfile } from '../types'

const students = shallowRef<StudentListItem[]>([])
const selectedId = shallowRef<number | null>(null)
const selectedStudent = shallowRef<StudentProfile | null>(null)
const errorMessage = shallowRef('')

async function loadStudents() {
  errorMessage.value = ''
  try {
    const data = await apiFetch<{ items: StudentListItem[] }>('/admin/students')
    students.value = data.items
    if (!selectedId.value && data.items.length > 0) {
      await selectStudent(data.items[0].id)
    }
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '加载学生列表失败'
  }
}

async function selectStudent(id: number) {
  selectedId.value = id
  try {
    const data = await apiFetch<{ student: StudentProfile }>(`/admin/students/${id}`)
    selectedStudent.value = data.student
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '加载学生详情失败'
  }
}

onMounted(loadStudents)
</script>

<template>
  <AdminShell>
    <section class="page-section">
      <div class="page-header">
        <div>
          <p class="eyebrow eyebrow--warning">Student Accounts</p>
          <h1>学生资料后台</h1>
          <p class="page-description">
            这里展示管理员能看到的全部学生双卡信息，可用于对照学生端 JWT 越权读取效果。
          </p>
        </div>
      </div>

      <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>

      <div class="content-grid">
        <StudentTable :items="students" :selected-id="selectedId" @select="selectStudent" />

        <article class="card">
          <div class="section-header">
            <h2>学生详情</h2>
            <span class="muted">后台视图</span>
          </div>

          <div v-if="selectedStudent" class="detail-grid">
            <div><span class="label">姓名</span><strong>{{ selectedStudent.displayName }}</strong></div>
            <div><span class="label">用户名</span><strong>{{ selectedStudent.username }}</strong></div>
            <div><span class="label">学号</span><strong>{{ selectedStudent.studentNo }}</strong></div>
            <div><span class="label">校园卡号</span><strong>{{ selectedStudent.campusCardNo }}</strong></div>
            <div><span class="label">校园卡余额</span><strong>{{ selectedStudent.campusCardBalance.toFixed(3) }} 元</strong></div>
            <div><span class="label">银行卡号</span><strong>{{ selectedStudent.bankCardNo }}</strong></div>
            <div><span class="label">银行卡余额</span><strong>{{ selectedStudent.bankBalance.toFixed(3) }} 元</strong></div>
            <div><span class="label">手机号</span><strong>{{ selectedStudent.phone }}</strong></div>
            <div><span class="label">宿舍</span><strong>{{ selectedStudent.dormitory }}</strong></div>
          </div>
          <p v-else class="muted">请选择一个学生。</p>
        </article>
      </div>
    </section>
  </AdminShell>
</template>

<style scoped>
.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(320px, 0.8fr);
  gap: 1rem;
}

.card {
  padding: 1rem 1.1rem;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: var(--color-panel);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.label {
  display: block;
  margin-bottom: 0.35rem;
  color: var(--color-muted);
  font-size: 0.85rem;
}

@media (max-width: 960px) {
  .content-grid,
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
