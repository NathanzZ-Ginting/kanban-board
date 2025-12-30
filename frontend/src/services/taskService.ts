import { taskApi } from '@/utils/axios'
import { API_ENDPOINTS } from '@/config/api.config'
import { Task, ApiResponse } from '@/types'

export interface CreateTaskData {
  title: string
  description?: string
  boardId: string
  assigneeId?: string
  status: Task['status']
  priority: Task['priority']
  dueDate?: string
}

export const taskService = {
  async getTasks(): Promise<ApiResponse<Task[]>> {
    const response = await taskApi.get(API_ENDPOINTS.TASK.LIST)
    return response.data
  },

  async getTask(id: string): Promise<ApiResponse<Task>> {
    const response = await taskApi.get(API_ENDPOINTS.TASK.GET(id))
    return response.data
  },

  async getTasksByBoard(boardId: string): Promise<ApiResponse<Task[]>> {
    const response = await taskApi.get(API_ENDPOINTS.TASK.BY_BOARD(boardId))
    return response.data
  },

  async createTask(data: CreateTaskData): Promise<ApiResponse<Task>> {
    const response = await taskApi.post(API_ENDPOINTS.TASK.CREATE, data)
    return response.data
  },

  async updateTask(id: string, data: Partial<Task>): Promise<ApiResponse<Task>> {
    const response = await taskApi.put(API_ENDPOINTS.TASK.UPDATE(id), data)
    return response.data
  },

  async deleteTask(id: string): Promise<ApiResponse> {
    const response = await taskApi.delete(API_ENDPOINTS.TASK.DELETE(id))
    return response.data
  },
}
