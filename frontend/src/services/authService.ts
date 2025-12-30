import { authApi } from '@/utils/axios'
import { API_ENDPOINTS } from '@/config/api.config'
import { User, AuthTokens, ApiResponse } from '@/types'

export interface RegisterData {
  email: string
  username: string
  password: string
  firstName: string
  lastName: string
}

export interface LoginData {
  email: string
  password: string
}

export const authService = {
  async register(data: RegisterData): Promise<ApiResponse<{ user: User; tokens: AuthTokens }>> {
    const response = await authApi.post(API_ENDPOINTS.AUTH.REGISTER, data)
    return response.data
  },

  async login(data: LoginData): Promise<ApiResponse<{ user: User; tokens: AuthTokens }>> {
    const response = await authApi.post(API_ENDPOINTS.AUTH.LOGIN, data)
    return response.data
  },

  async logout(): Promise<ApiResponse> {
    const response = await authApi.post(API_ENDPOINTS.AUTH.LOGOUT)
    return response.data
  },

  async refreshToken(refreshToken: string): Promise<ApiResponse<{ accessToken: string }>> {
    const response = await authApi.post(API_ENDPOINTS.AUTH.REFRESH, { refreshToken })
    return response.data
  },
}
