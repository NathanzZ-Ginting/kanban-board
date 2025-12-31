import axios from 'axios';
import { API_CONFIG, API_ENDPOINTS } from '@/config/api.config';
import { User, AuthTokens, ApiResponse } from '@/types';

const authApi = axios.create({
  baseURL: API_CONFIG.AUTH_SERVICE,
  timeout: API_CONFIG.TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface RegisterData {
  email: string;
  username: string;
  password: string;
  firstName?: string;
  lastName?: string;
}

export interface LoginData {
  email: string;
  password: string;
}

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

export const authService = {
  async register(data: RegisterData): Promise<AuthResponse> {
    const response = await authApi.post(API_ENDPOINTS.AUTH.REGISTER, data);
    return response.data;
  },

  async login(data: LoginData): Promise<AuthResponse> {
    const response = await authApi.post(API_ENDPOINTS.AUTH.LOGIN, data);
    return response.data;
  },

  async logout(refreshToken: string): Promise<void> {
    await authApi.post(API_ENDPOINTS.AUTH.LOGOUT, { refreshToken });
  },

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    const response = await authApi.post(API_ENDPOINTS.AUTH.REFRESH, { refreshToken });
    return response.data;
  },
};
