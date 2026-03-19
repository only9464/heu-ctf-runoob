export interface AuthUser {
  id: number
  username: string
  displayName: string
  role: 'student' | 'admin'
  studentNo: string
  campusCardNo: string
  bankCardNo: string
  campusCardBalance: number
  bankCardBalance: number
}

export interface StudentProfile {
  id: number
  username: string
  displayName: string
  studentNo: string
  campusCardNo: string
  bankCardNo: string
  campusCardBalance: number
  bankBalance: number
  phone: string
  dormitory: string
}

export interface MealRecord {
  id: number
  canteenWindow: string
  dishName: string
  amount: number
  consumedAt: string
}

export interface RechargeRecord {
  id: number
  rawAmount: number
  creditedAmount: number
  note: string
  createdAt: string
}

export interface MessageItem {
  id: number
  studentName: string
  content: string
  createdAt: string
}

export interface AdminSummary {
  studentCount: number
  messageCount: number
  rechargeCount: number
  adminUsername: string
}

export interface StudentListItem {
  id: number
  username: string
  displayName: string
  studentNo: string
  campusCardBalance: number
  bankBalance: number
  phone: string
  dormitory: string
}
