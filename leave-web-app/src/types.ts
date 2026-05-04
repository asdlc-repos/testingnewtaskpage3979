export interface LeaveRequest {
  id: string
  userId: string
  startDate: string
  endDate: string
  reason: string
  status: 'pending' | 'approved' | 'rejected'
  managerComment?: string
  createdAt?: string
}

export interface LeaveBalance {
  userId: string
  balance: number
  entitlement: number
  used: number
}

export interface User {
  id: string
  name: string
  email: string
  role: 'employee' | 'manager'
  managerId?: string
}

export interface ApprovalPayload {
  comment?: string
}

export interface LeaveRequestPayload {
  userId: string
  startDate: string
  endDate: string
  reason: string
}
