import { useState } from 'react'
import { submitLeaveRequest } from '../api'

interface Props {
  userId: string
  onSuccess: () => void
}

export default function LeaveRequestForm({ userId, onSuccess }: Props) {
  const today = new Date().toISOString().split('T')[0]
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')
  const [reason, setReason] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [validationError, setValidationError] = useState('')

  function validate() {
    if (!startDate) return 'Start date is required'
    if (!endDate) return 'End date is required'
    if (endDate < startDate) return 'End date must be on or after start date'
    if (!reason.trim()) return 'Reason is required'
    return ''
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const err = validate()
    if (err) {
      setValidationError(err)
      return
    }
    setValidationError('')
    setError('')
    setSubmitting(true)
    try {
      await submitLeaveRequest({ userId, startDate, endDate, reason: reason.trim() })
      setStartDate('')
      setEndDate('')
      setReason('')
      onSuccess()
    } catch {
      setError('Failed to submit leave request. Please try again.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="bg-white rounded-lg shadow p-5">
      <h2 className="text-lg font-semibold text-gray-700 mb-4">Submit Leave Request</h2>
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 rounded-md px-4 py-2 text-sm mb-4">
          {error}
        </div>
      )}
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
            <input
              type="date"
              min={today}
              value={startDate}
              onChange={(e) => { setStartDate(e.target.value); setValidationError('') }}
              className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">End Date</label>
            <input
              type="date"
              min={startDate || today}
              value={endDate}
              onChange={(e) => { setEndDate(e.target.value); setValidationError('') }}
              className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
            />
          </div>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Reason</label>
          <textarea
            rows={3}
            value={reason}
            onChange={(e) => { setReason(e.target.value); setValidationError('') }}
            placeholder="Briefly describe the reason for your leave..."
            className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm resize-none"
          />
        </div>
        {validationError && (
          <p className="text-red-500 text-sm">{validationError}</p>
        )}
        <button
          type="submit"
          disabled={submitting}
          className="bg-blue-600 text-white px-5 py-2 rounded-md hover:bg-blue-700 transition-colors font-medium text-sm disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {submitting ? 'Submitting...' : 'Submit Request'}
        </button>
      </form>
    </div>
  )
}
