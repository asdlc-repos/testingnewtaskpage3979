import axios from 'axios'
import type { LeaveRequest, LeaveBalance, LeaveRequestPayload, ApprovalPayload } from './types'

const w = window as Window & typeof globalThis & {
  RUNTIME_LEAVE_SERVICE_URL?: string
  RUNTIME_USER_SERVICE_URL?: string
}

const LEAVE_API_BASE: string =
  w.RUNTIME_LEAVE_SERVICE_URL ||
  (import.meta as { env: Record<string, string> }).env.VITE_LEAVE_API_BASE ||
  '/api/leave'

const USER_API_BASE: string =
  w.RUNTIME_USER_SERVICE_URL ||
  (import.meta as { env: Record<string, string> }).env.VITE_USER_API_BASE ||
  '/api/user'

const leaveClient = axios.create({ baseURL: LEAVE_API_BASE })
const userClient = axios.create({ baseURL: USER_API_BASE })

export async function getLeaveRequests(userId: string): Promise<LeaveRequest[]> {
  const res = await leaveClient.get<LeaveRequest[]>('/api/v1/requests', {
    params: { userId },
  })
  return res.data
}

export async function submitLeaveRequest(payload: LeaveRequestPayload): Promise<LeaveRequest> {
  const res = await leaveClient.post<LeaveRequest>('/api/v1/requests', payload)
  return res.data
}

export async function getPendingRequests(managerId: string): Promise<LeaveRequest[]> {
  const res = await leaveClient.get<LeaveRequest[]>('/api/v1/requests', {
    params: { managerId, status: 'pending' },
  })
  return res.data
}

export async function approveRequest(id: string, payload?: ApprovalPayload): Promise<LeaveRequest> {
  const res = await leaveClient.put<LeaveRequest>(`/api/v1/requests/${id}/approve`, payload || {})
  return res.data
}

export async function rejectRequest(id: string, payload?: ApprovalPayload): Promise<LeaveRequest> {
  const res = await leaveClient.put<LeaveRequest>(`/api/v1/requests/${id}/reject`, payload || {})
  return res.data
}

export async function getLeaveBalance(userId: string): Promise<LeaveBalance> {
  const res = await userClient.get<LeaveBalance>(`/api/v1/users/${userId}/balance`)
  return res.data
}
