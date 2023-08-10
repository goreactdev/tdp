import { AuthBindings } from "@refinedev/core";
import { axiosInstance } from "@refinedev/simple-rest";
import { Account, TonProofItemReplySuccess, useTonConnectUI } from "@tonconnect/ui-react";
import { API_URL } from "./App";

export const TOKEN_KEY = "refine-auth";
export const USER_KEY = "refine-user";
export const EXPIRATION_TIME = "refine-expiration"


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
  rating: number
  last_award_at?: number
  linked_accounts?: Array<LinkedAccount>
  created_at: number
  updated_at: number
  version: number
}


export interface LinkedAccount {
  id: number;
  user_id: number;
  provider: string;
  avatar_url: string;
  login: string;
  access_token?: string; 
  created_at: number;
  updated_at: number;
  version: number;
}


export interface ProofCheckResponse {
  user: User
  token: string
  expires: number;
}

export const authProvider: AuthBindings = {
  
  login: async ({ address, network, proof }) => {
    if (address && proof && network) {
      
      // do this with axios
      const response = await axiosInstance.post<ProofCheckResponse>(API_URL + "/v1/ton-connect/check-proof", {
        address: address,
        network: network,
        proof: proof,
      } as ProofCheckPayload);
      
      localStorage.setItem(TOKEN_KEY, response.data.token);

      localStorage.setItem(EXPIRATION_TIME, String(response.data.expires));

      localStorage.setItem(USER_KEY, JSON.stringify(response.data.user));

      axiosInstance.defaults.headers.common = {
        Authorization: `Bearer ${response.data.token}`,
    };


      return {
        success: true,
        redirectTo: "/",
      };
    }

    return {
      success: false,
      error: {
        name: "LoginError",
        message: "Invalid username or password",
      },
    };
  },
  logout: async () => {

    localStorage.removeItem(TOKEN_KEY);

    return {
      success: true,
      redirectTo: "/admin/login",
    };
  },
  check: async () => {
    const token = localStorage.getItem(TOKEN_KEY);

    const expires = localStorage.getItem(EXPIRATION_TIME);
    
    if (token && expires && (parseInt(expires) * 1000) > Date.now()) {
      return {
        authenticated: true,
      };
    }

    if (token) {
      localStorage.removeItem(TOKEN_KEY);

      localStorage.removeItem(EXPIRATION_TIME);

      localStorage.removeItem(USER_KEY);
      // remove everything related to ton-connect
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);

        if (key && key.includes("ton-connect")) {
          localStorage.removeItem(key);
        }
      }

      axiosInstance.defaults.headers.common = {};
   };

   return {
     authenticated: false,
   };

  },
  getIdentity: async () => {
    const token = localStorage.getItem(TOKEN_KEY);

    const response = await axiosInstance.get<User>(API_URL + "/v1/my-account", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (token) {
      return {
        id: 1,
        name: response.data.username,
        avatar: response.data.avatar_url,
      };
    }
    return null;
  },
  onError: async (error) => {
    console.error(error);
    return { error };
  },
};
