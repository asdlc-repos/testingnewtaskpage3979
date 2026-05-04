import { Link, useLocation } from 'react-router-dom'

interface Props {
  userId: string
  isManager: boolean
  onLogout: () => void
}

export default function Navbar({ userId, isManager, onLogout }: Props) {
  const location = useLocation()

  return (
    <nav className="bg-blue-700 text-white shadow-md">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-6">
          <span className="font-bold text-lg">Leave Management</span>
          <Link
            to="/employee"
            className={`text-sm px-3 py-1 rounded-md transition-colors ${
              location.pathname === '/employee'
                ? 'bg-blue-900 font-semibold'
                : 'hover:bg-blue-600'
            }`}
          >
            My Leaves
          </Link>
          {isManager && (
            <Link
              to="/manager"
              className={`text-sm px-3 py-1 rounded-md transition-colors ${
                location.pathname === '/manager'
                  ? 'bg-blue-900 font-semibold'
                  : 'hover:bg-blue-600'
              }`}
            >
              Manager Dashboard
            </Link>
          )}
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="text-blue-200">
            {isManager ? 'Manager' : 'Employee'}: <strong className="text-white">{userId}</strong>
          </span>
          <button
            onClick={onLogout}
            className="bg-blue-800 hover:bg-blue-900 px-3 py-1 rounded-md transition-colors"
          >
            Sign Out
          </button>
        </div>
      </div>
    </nav>
  )
}
