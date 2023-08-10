import type { Account, TonProofItemReplySuccess } from '@tonconnect/ui-react'

export interface User {
  id: number
  first_name: string
  last_name: string
  username: string
  raw_address: string
  friendly_address: string
  job?: string
  bio?: string
  languages?: Array<string>
  certifications?: Array<string>
  avatar_url?: string
  awards_count: number
  messages_count: number
  rating: number
  last_award_at?: number
  linked_accounts?: Array<LinkedAccount>
  created_at: number
  updated_at: number
  version: number
}

export interface ProofCheckPayload {
  address: string
  network: string
  proof: {
    timestamp: number
    domain: {
      lengthBytes: number
      value: string
    }
    payload: string
    signature: string
  }
}

export interface ProofCheckRequest {
  profItemReply: TonProofItemReplySuccess
  account: Account
}

export interface ProofCheckResponse {
  user: User
  token: string
}

export interface ProofPayloadResponse {
  payload: string
}

export interface LinkedAccount {
  id: number
  user_id: number
  provider: string
  avatar_url: string
  login: string
  access_token?: string
  created_at: number
  updated_at: number
  version: number
}

export interface SBTToken {
  id: number
  raw_address: string
  friendly_address: string
  collection_id: number
  content_uri: string
  raw_owner_address: string
  friendly_owner_address: string
  name: string
  is_pinned: boolean
  description: string
  image: string
  content_json: string
  weight: number
  index: number
  created_at: string
  updated_at: string
  version: number
}

export interface StoredReward {
  achievement_id: number
  image_url: string
  name: string
  description: string
  weight: number
  approved: boolean
  processed: boolean
}
