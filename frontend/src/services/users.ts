import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import ky from 'ky'

import type { RootState } from '../store'
import { BASE_URL } from '../utils/config'

import type {
  ProofCheckPayload,
  ProofCheckRequest,
  ProofCheckResponse,
  ProofPayloadResponse,
  SBTToken,
  StoredReward,
  User,
} from './types'


// Define a service using a base URL and expected endpoints
export const userApi = createApi({
  baseQuery: fetchBaseQuery({
    baseUrl: BASE_URL,
    fetchFn: (...args) => ky(...args),
    prepareHeaders: (headers, { getState }) => {
      const token = (getState() as RootState).authReducer.token

      if (token) {
        headers.set('Authorization', `Bearer ${token}`)
      }
      return headers
    },
  }),
  endpoints: (builder) => ({
    checkProofInBackend: builder.mutation<
      ProofCheckResponse,
      ProofCheckRequest
    >({
      invalidatesTags: ['User'],
      query: ({ profItemReply, account }) => ({
        body: {
          address: account.address,
          network: account.chain,
          proof: profItemReply.proof,
        } as ProofCheckPayload,
        method: 'POST',
        url: `/v1/ton-connect/check-proof`,
      }),
    }),

    getAchievements: builder.query<
      {
        achievements: StoredReward[]
        count: number
      },
      {
        start: number
        end: number
      }
    >({
      providesTags: ['SBT'],
      query: ({ start, end }) =>
        `/v1/incoming-achievements?_start=${start}&_end=${end}`,
    }),

    checkAuthTelegram: builder.mutation<void, { auth_obj: string }>({
      query: ({ auth_obj }) => ({
        url: `/v1/telegram/check_authorization`,
        method: 'POST',
        body: {
          auth_obj,
        },
      }),
      invalidatesTags: ['Account'],
    }),

    getMyAccount: builder.query<User, void>({
      providesTags: ['Account'],
      query: () => `/v1/my-account`,
    }),

    getNftsByUsername: builder.query<
      { nfts: SBTToken[]; count: number },
      { username: string; start: number; end: number }
    >({
      providesTags: ['User'],
      query: ({ username, start, end }) =>
        `/v1/nfts/${username}?_start=${start}&_end=${end}`,
    }),

    getProofPayload: builder.query<ProofPayloadResponse, void>({
      query: () => `/v1/ton-connect/generate-payload`,
    }),

    getTopUsers: builder.query<
      User[],
      {
        start: number
        end: number
      }
    >({
      providesTags: ['User'],
      query: ({ start, end }) => `/v1/users?_start=${start}&_end=${end}`,
    }),

    getUserByUsername: builder.query<{ user: User }, { username: string }>({
      providesTags: ['User'],
      query: ({ username }) => `/v1/users/${username}`,
    }),

    pinNft: builder.mutation<void, { id: number }>({
      invalidatesTags: ['User'],
      query: ({ id }) => ({
        method: 'PUT',
        url: `/v1/nft/${id}/pin`,
      }),
    }),

    unlinkAccount: builder.mutation<void, { provider: string }>({
      invalidatesTags: ['Account'],
      query: ({ provider }) => ({
        method: 'DELETE',
        url: `/v1/unlink/${provider}`,
      }),
    }),

    updateAchievement: builder.mutation<
      void,
      {
        id: number
        approved_by_user: boolean
      }
    >({
      invalidatesTags: ['SBT'],
      query: ({ id, approved_by_user }) => ({
        body: {
          approved_by_user,
        },
        method: 'PUT',
        url: `/v1/incoming-achievements/${id}`,
      }),
    }),

    updateUser: builder.mutation<User, Partial<User>>({
      invalidatesTags: ['Account'],
      query: ({ id, ...patch }) => ({
        body: patch,
        method: 'PATCH',
        url: `/v1/update/users`,
      }),
    }),

    uploadImage: builder.mutation<{ url: string }, FormData>({
      query: (file) => ({
        body: file,
        method: 'POST',
        url: `/v1/admin/media/upload`,
      }),
      invalidatesTags: ['Account'],
    }),
  }),
  reducerPath: 'userApi',
  tagTypes: ['User', 'Account', 'SBT'],
})

// Export hooks for usage in functional components, which are
// auto-generated based on the defined endpoints
export const {
  useGetUserByUsernameQuery,
  useGetProofPayloadQuery,
  useCheckProofInBackendMutation,
  useGetAchievementsQuery,
  useGetMyAccountQuery,
  useUploadImageMutation,
  useCheckAuthTelegramMutation,
  useUpdateUserMutation,
  useGetTopUsersQuery,
  usePinNftMutation,
  useUpdateAchievementMutation,
  useGetNftsByUsernameQuery,
  useUnlinkAccountMutation,
} = userApi
