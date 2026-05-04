import { useState } from 'react'

interface Props {
  onLogin: (userId: string, isManager: boolean) => void
}

export default function LoginForm({ onLogin }: Props) {
  const [userId, setUserId] = useState('')
  const [isManager, setIsManager] = useState(false)
  const [error, setError] = useState('')

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!userId.trim()) {
      setError('User ID is required')
      return
    }
    onLogin(userId.trim(), isManager)
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-8 w-full max-w-md">
      <div className="text-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Leave Management</h1>
        <p className="text-gray-500 mt-1">Sign in to continue</p>
      </div>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="userId" className="block text-sm font-medium text-gray-700 mb-1">
            User ID
          </label>
          <input
            id="userId"
            type="text"
            value={userId}
            onChange={(e) => { setUserId(e.target.value); setError('') }}
            placeholder="e.g. user-123"
            className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
        </div>
        <div className="flex items-center gap-2">
          <input
            id="isManager"
            type="checkbox"
            checked={isManager}
            onChange={(e) => setIsManager(e.target.checked)}
            className="h-4 w-4 text-blue-600 border-gray-300 rounded"
          />
          <label htmlFor="isManager" className="text-sm text-gray-700">
            I am a manager
          </label>
        </div>
        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors font-medium"
        >
          Sign In
        </button>
      </form>
    </div>
  )
}
