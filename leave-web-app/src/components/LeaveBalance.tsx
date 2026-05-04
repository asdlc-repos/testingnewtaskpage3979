import { useEffect, useState } from 'react'
import { getLeaveBalance } from '../api'
import type { LeaveBalance as LeaveBalanceType } from '../types'

interface Props {
  userId: string
}

export default function LeaveBalance({ userId }: Props) {
  const [balance, setBalance] = useState<LeaveBalanceType | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true)
    getLeaveBalance(userId)
      .then(setBalance)
      .catch(() => setError('Failed to load leave balance'))
      .finally(() => setLoading(false))
  }, [userId])

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow p-5 animate-pulse">
        <div className="h-5 bg-gray-200 rounded w-32 mb-3" />
        <div className="flex gap-6">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-10 w-20 bg-gray-200 rounded" />
          ))}
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700 text-sm">
        {error}
      </div>
    )
  }

  if (!balance) return null

  const usedPct = balance.entitlement > 0
    ? Math.round((balance.used / balance.entitlement) * 100)
    : 0

  return (
    <div className="bg-white rounded-lg shadow p-5">
      <h2 className="text-lg font-semibold text-gray-700 mb-4">Leave Balance</h2>
      <div className="grid grid-cols-3 gap-4 mb-4">
        <div className="text-center">
          <p className="text-3xl font-bold text-blue-600">{balance.balance}</p>
          <p className="text-xs text-gray-500 mt-1">Available</p>
        </div>
        <div className="text-center">
          <p className="text-3xl font-bold text-gray-700">{balance.entitlement}</p>
          <p className="text-xs text-gray-500 mt-1">Annual Entitlement</p>
        </div>
        <div className="text-center">
          <p className="text-3xl font-bold text-orange-500">{balance.used}</p>
          <p className="text-xs text-gray-500 mt-1">Used</p>
        </div>
      </div>
      <div className="w-full bg-gray-200 rounded-full h-2">
        <div
          className="bg-blue-500 h-2 rounded-full transition-all"
          style={{ width: `${usedPct}%` }}
        />
      </div>
      <p className="text-xs text-gray-400 mt-1 text-right">{usedPct}% used</p>
    </div>
  )
}
