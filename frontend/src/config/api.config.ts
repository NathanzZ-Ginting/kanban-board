export const API_CONFIG = {
  AUTH_SERVICE: process.env.NEXT_PUBLIC_AUTH_SERVICE_URL || 'http://localhost:8001',
  USER_SERVICE: process.env.NEXT_PUBLIC_USER_SERVICE_URL || 'http://localhost:8002',
  BOARD_SERVICE: process.env.NEXT_PUBLIC_BOARD_SERVICE_URL || 'http://localhost:8003',
  TASK_SERVICE: process.env.NEXT_PUBLIC_TASK_SERVICE_URL || 'http://localhost:8004',
  TIMEOUT: 10000,
}

export const API_ENDPOINTS = {
  AUTH: {
    REGISTER: '/api/v1/auth/register',
    LOGIN: '/api/v1/auth/login',
    LOGOUT: '/api/v1/auth/logout',
    REFRESH: '/api/v1/auth/refresh',
  },
  USER: {
    PROFILE: '/api/v1/users/profile',
    UPDATE: '/api/v1/users/profile',
    LIST: '/api/v1/users',
  },
  BOARD: {
    LIST: '/api/v1/boards',
    CREATE: '/api/v1/boards',
    GET: (id: string) => `/api/v1/boards/${id}`,
    UPDATE: (id: string) => `/api/v1/boards/${id}`,
    DELETE: (id: string) => `/api/v1/boards/${id}`,
  },
  TASK: {
    LIST: '/api/v1/tasks',
    CREATE: '/api/v1/tasks',
    GET: (id: string) => `/api/v1/tasks/${id}`,
    UPDATE: (id: string) => `/api/v1/tasks/${id}`,
    DELETE: (id: string) => `/api/v1/tasks/${id}`,
    BY_BOARD: (boardId: string) => `/api/v1/boards/${boardId}/tasks`,
  },
}
