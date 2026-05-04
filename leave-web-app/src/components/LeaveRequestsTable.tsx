import type { LeaveRequest } from '../types'

interface Props {
  requests: LeaveRequest[]
  loading: boolean
  error: string
}

const statusClasses: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800',
}

export default function LeaveRequestsTable({ requests, loading, error }: Props) {
  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow p-5">
        <h2 className="text-lg font-semibold text-gray-700 mb-4">My Leave Requests</h2>
        <div className="space-y-2 animate-pulse">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-10 bg-gray-100 rounded" />
          ))}
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow p-5">
        <h2 className="text-lg font-semibold text-gray-700 mb-4">My Leave Requests</h2>
        <div className="bg-red-50 border border-red-200 text-red-700 rounded-md px-4 py-2 text-sm">
          {error}
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow p-5">
      <h2 className="text-lg font-semibold text-gray-700 mb-4">My Leave Requests</h2>
      {requests.length === 0 ? (
        <p className="text-gray-400 text-sm text-center py-6">No leave requests found.</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-gray-500 border-b">
                <th className="pb-2 pr-4">Start Date</th>
                <th className="pb-2 pr-4">End Date</th>
                <th className="pb-2 pr-4">Reason</th>
                <th className="pb-2 pr-4">Status</th>
                <th className="pb-2">Comment</th>
              </tr>
            </thead>
            <tbody>
              {requests.map((req) => (
                <tr key={req.id} className="border-b last:border-0 hover:bg-gray-50">
                  <td className="py-3 pr-4">{req.startDate}</td>
                  <td className="py-3 pr-4">{req.endDate}</td>
                  <td className="py-3 pr-4 max-w-xs truncate" title={req.reason}>
                    {req.reason}
                  </td>
                  <td className="py-3 pr-4">
                    <span
                      className={`px-2 py-0.5 rounded-full text-xs font-medium capitalize ${
                        statusClasses[req.status] || 'bg-gray-100 text-gray-700'
                      }`}
                    >
                      {req.status}
                    </span>
                  </td>
                  <td className="py-3 text-gray-500 text-xs">
                    {req.managerComment || '—'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
