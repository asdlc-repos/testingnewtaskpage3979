import { useState, useEffect } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Navbar from './components/Navbar'
import EmployeeView from './components/EmployeeView'
import ManagerView from './components/ManagerView'
import LoginForm from './components/LoginForm'

export default function App() {
  const [userId, setUserId] = useState<string>(() => localStorage.getItem('userId') || '')
  const [isManager, setIsManager] = useState<boolean>(
    () => localStorage.getItem('isManager') === 'true'
  )

  useEffect(() => {
    if (userId) {
      localStorage.setItem('userId', userId)
    } else {
      localStorage.removeItem('userId')
    }
  }, [userId])

  useEffect(() => {
    localStorage.setItem('isManager', String(isManager))
  }, [isManager])

  function handleLogin(id: string, manager: boolean) {
    setUserId(id)
    setIsManager(manager)
  }

  function handleLogout() {
    setUserId('')
    setIsManager(false)
    localStorage.removeItem('userId')
    localStorage.removeItem('isManager')
  }

  if (!userId) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <LoginForm onLogin={handleLogin} />
      </div>
    )
  }

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        <Navbar userId={userId} isManager={isManager} onLogout={handleLogout} />
        <main className="max-w-7xl mx-auto px-4 py-6">
          <Routes>
            <Route path="/" element={<Navigate to="/employee" replace />} />
            <Route path="/employee" element={<EmployeeView userId={userId} />} />
            {isManager && (
              <Route path="/manager" element={<ManagerView managerId={userId} />} />
            )}
            <Route path="*" element={<Navigate to="/employee" replace />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}
