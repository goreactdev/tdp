import { createSlice } from '@reduxjs/toolkit'
import type { FetchBaseQueryError } from '@reduxjs/toolkit/dist/query'

import type { User } from '../../services/types'
import { userApi } from '../../services/users'
import type { RootState } from '../../store'

interface AuthState {
  token: string | null
  isAuthenticated: boolean
  user: User | null
  errors: FetchBaseQueryError | undefined
}

const initialState: AuthState = {
  errors: undefined,
  isAuthenticated: false,
  token: '',
  user: null,
}

export const authSlice = createSlice({
  extraReducers: (builder) => {
    builder.addMatcher(
      userApi.endpoints.checkProofInBackend.matchRejected,
      (state, { payload: errors }) => {
        state.errors = errors
        state.isAuthenticated = false
        state.user = null
      }
    )

    builder.addMatcher(
      userApi.endpoints.checkProofInBackend.matchPending,
      (state) => {
        state.errors = undefined
        state.isAuthenticated = false
        state.user = null
      }
    )

    builder.addMatcher(
      userApi.endpoints.checkProofInBackend.matchFulfilled,
      (state, { payload }) => {
        state.token = payload.token
        state.isAuthenticated = true
        state.user = payload.user
      }
    )

    builder.addMatcher(
      userApi.endpoints.updateUser.matchFulfilled,
      (state, { payload }) => {
        state.user = payload
      }
    )
  },
  initialState,
  name: 'auth',
  reducers: {
    logout: (state) => {
      state.token = null
      state.isAuthenticated = false
      state.user = null
    },
  },
})

export default authSlice.reducer

export const { logout } = authSlice.actions

// Selectors
export const selectCurrentUser = (state: RootState) => state.authReducer.user
