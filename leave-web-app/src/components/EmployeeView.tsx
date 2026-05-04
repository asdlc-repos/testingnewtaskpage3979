import { useState, useEffect, useCallback } from 'react'
import { getLeaveRequests } from '../api'
import type { LeaveRequest } from '../types'
import LeaveBalance from './LeaveBalance'
import LeaveRequestForm from './LeaveRequestForm'
import LeaveRequestsTable from './LeaveRequestsTable'

interface Props {
  userId: string
}

export default function EmployeeView({ userId }: Props) {
  const [requests, setRequests] = useState<LeaveRequest[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const fetchRequests = useCallback(() => {
    setLoading(true)
    setError('')
    getLeaveRequests(userId)
      .then(setRequests)
      .catch(() => setError('Failed to load leave requests'))
      .finally(() => setLoading(false))
  }, [userId])

  useEffect(() => {
    fetchRequests()
  }, [fetchRequests])

  return (
    <div className="space-y-6">
      <LeaveBalance userId={userId} />
      <LeaveRequestForm userId={userId} onSuccess={fetchRequests} />
      <LeaveRequestsTable requests={requests} loading={loading} error={error} />
    </div>
  )
}
