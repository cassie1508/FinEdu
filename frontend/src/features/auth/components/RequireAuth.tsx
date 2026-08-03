import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useSession } from '../../../lib/useSession'

export function RequireAuth({ children }: { children: ReactNode }) {
  const session = useSession()

  if (session === undefined) return null
  if (!session) return <Navigate to="/login" replace />

  return children
}
