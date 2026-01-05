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
        console.log('Request with token:', config.method?.toUpperCase(), config.url)
      } else {
        console.warn('No access token found in localStorage')
      }
      return config
    },
    (error) => {
      console.error('Request error:', error)
      return Promise.reject(error)
    }
  )

  // Response interceptor
  instance.interceptors.response.use(
    (response) => {
      console.log('Response received:', response.status, response.config.url)
      return response
    },
    async (error) => {
      const originalRequest = error.config

      console.error('Response error:', {
        status: error.response?.status,
        url: error.config?.url,
        data: error.response?.data
      })

      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true

        try {
          const refreshToken = localStorage.getItem('refreshToken')
          
          if (!refreshToken) {
            console.error('No refresh token found')
            localStorage.removeItem('accessToken')
            localStorage.removeItem('refreshToken')
            window.location.href = '/login'
            return Promise.reject(error)
          }

          console.log('Attempting to refresh token...')
          const response = await axios.post(
            `${API_CONFIG.AUTH_SERVICE}/api/v1/auth/refresh`,
            { refreshToken }
          )

          const { accessToken } = response.data
          localStorage.setItem('accessToken', accessToken)
          console.log('Token refreshed successfully')

          originalRequest.headers.Authorization = `Bearer ${accessToken}`
          return instance(originalRequest)
        } catch (refreshError) {
          console.error('Token refresh failed:', refreshError)
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
