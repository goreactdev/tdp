import { useMemo } from 'react'
import { useSelector } from 'react-redux'

import { selectCurrentUser } from '../features/auth/authSlice'

export const useMemoizedUser = () => {
  const user = useSelector(selectCurrentUser)
  return useMemo(() => ({ user }), [user])
}
