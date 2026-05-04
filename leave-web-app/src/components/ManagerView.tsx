import { useState, useEffect, useCallback } from 'react'
import { getPendingRequests, approveRequest, rejectRequest } from '../api'
import type { LeaveRequest } from '../types'

interface Props {
  managerId: string
}

export default function ManagerView({ managerId }: Props) {
  const [requests, setRequests] = useState<LeaveRequest[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [actionInProgress, setActionInProgress] = useState<string | null>(null)
  const [comments, setComments] = useState<Record<string, string>>({})
  const [actionError, setActionError] = useState<Record<string, string>>({})

  const fetchPending = useCallback(() => {
    setLoading(true)
    setError('')
    getPendingRequests(managerId)
      .then(setRequests)
      .catch(() => setError('Failed to load pending requests'))
      .finally(() => setLoading(false))
  }, [managerId])

  useEffect(() => {
    fetchPending()
  }, [fetchPending])

  async function handleAction(id: string, action: 'approve' | 'reject') {
    setActionInProgress(id + action)
    setActionError((prev) => ({ ...prev, [id]: '' }))
    try {
      const comment = comments[id] || undefined
      if (action === 'approve') {
        await approveRequest(id, comment ? { comment } : undefined)
      } else {
        await rejectRequest(id, comment ? { comment } : undefined)
      }
      setRequests((prev) => prev.filter((r) => r.id !== id))
    } catch {
      setActionError((prev) => ({ ...prev, [id]: `Failed to ${action} request` }))
    } finally {
      setActionInProgress(null)
    }
  }

  return (
    <div className="bg-white rounded-lg shadow p-5">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-700">Pending Approval Requests</h2>
        <button
          onClick={fetchPending}
          className="text-sm text-blue-600 hover:underline"
        >
          Refresh
        </button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 rounded-md px-4 py-2 text-sm mb-4">
          {error}
        </div>
      )}

      {loading ? (
        <div className="space-y-2 animate-pulse">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-20 bg-gray-100 rounded" />
          ))}
        </div>
      ) : requests.length === 0 ? (
        <p className="text-gray-400 text-sm text-center py-8">No pending requests.</p>
      ) : (
        <div className="space-y-4">
          {requests.map((req) => (
            <div key={req.id} className="border border-gray-200 rounded-lg p-4">
              <div className="flex items-start justify-between">
                <div className="space-y-1">
                  <p className="text-sm font-medium text-gray-800">
                    Employee: <span className="text-blue-600">{req.userId}</span>
                  </p>
                  <p className="text-sm text-gray-600">
                    {req.startDate} &rarr; {req.endDate}
                  </p>
                  <p className="text-sm text-gray-500">{req.reason}</p>
                </div>
                <span className="px-2 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 capitalize">
                  {req.status}
                </span>
              </div>
              <div className="mt-3">
                <input
                  type="text"
                  placeholder="Optional comment..."
                  value={comments[req.id] || ''}
                  onChange={(e) =>
                    setComments((prev) => ({ ...prev, [req.id]: e.target.value }))
                  }
                  className="w-full border border-gray-200 rounded-md px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
                />
              </div>
              {actionError[req.id] && (
                <p className="text-red-500 text-xs mt-1">{actionError[req.id]}</p>
              )}
              <div className="mt-3 flex gap-2">
                <button
                  onClick={() => handleAction(req.id, 'approve')}
                  disabled={actionInProgress !== null}
                  className="bg-green-600 text-white px-4 py-1.5 rounded-md text-sm hover:bg-green-700 transition-colors disabled:opacity-50"
                >
                  {actionInProgress === req.id + 'approve' ? 'Approving...' : 'Approve'}
                </button>
                <button
                  onClick={() => handleAction(req.id, 'reject')}
                  disabled={actionInProgress !== null}
                  className="bg-red-600 text-white px-4 py-1.5 rounded-md text-sm hover:bg-red-700 transition-colors disabled:opacity-50"
                >
                  {actionInProgress === req.id + 'reject' ? 'Rejecting...' : 'Reject'}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
