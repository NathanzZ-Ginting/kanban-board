import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'
import { API_CONFIG } from '@/config/api.config'

const createAxiosInstance = (baseURL: string): AxiosInstance => {
  const instance = axios.create({
    baseURL,
    timeout: API_CONFIG.TIMEOUT,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  // Request interceptor
  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('accessToken')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    },
    (error) => Promise.reject(error)
  )

  // Response interceptor
  instance.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config

      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true

        try {
          const refreshToken = localStorage.getItem('refreshToken')
          const response = await axios.post(
            `${API_CONFIG.AUTH_SERVICE}/api/v1/auth/refresh`,
            { refreshToken }
          )

          const { accessToken } = response.data
          localStorage.setItem('accessToken', accessToken)

          originalRequest.headers.Authorization = `Bearer ${accessToken}`
          return instance(originalRequest)
        } catch (refreshError) {
          localStorage.removeItem('accessToken')
          localStorage.removeItem('refreshToken')
          window.location.href = '/login'
          return Promise.reject(refreshError)
        }
      }

      return Promise.reject(error)
    }
  )

  return instance
}

export const authApi = createAxiosInstance(API_CONFIG.AUTH_SERVICE)
export const userApi = createAxiosInstance(API_CONFIG.USER_SERVICE)
export const boardApi = createAxiosInstance(API_CONFIG.BOARD_SERVICE)
export const taskApi = createAxiosInstance(API_CONFIG.TASK_SERVICE)
