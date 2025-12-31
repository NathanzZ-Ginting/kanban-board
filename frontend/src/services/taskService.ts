import { taskApi } from '@/utils/axios'
import { API_ENDPOINTS } from '@/config/api.config'
import { Task } from '@/types'

export interface CreateTaskData {
  title: string
  description?: string
  boardId: number
  columnId: number
  priority: Task['priority']
  assigneeId?: number
  dueDate?: string
}

export interface UpdateTaskData {
  title?: string
  description?: string
  columnId?: number
  priority?: Task['priority']
  position?: number
  assigneeId?: number
  dueDate?: string
}

export interface MoveTaskData {
  columnId: number
  position: number
}

export interface TasksListResponse {
  tasks: Task[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export const taskService = {
  async getTasks(): Promise<TasksListResponse> {
    const response = await taskApi.get(API_ENDPOINTS.TASK.LIST)
    return response.data
  },

  async getTask(id: string): Promise<Task> {
    const response = await taskApi.get(API_ENDPOINTS.TASK.GET(id))
    return response.data
  },

  async getTasksByBoard(boardId: string): Promise<TasksListResponse> {
    const response = await taskApi.get(API_ENDPOINTS.TASK.BY_BOARD(boardId))
    return response.data
  },

  async createTask(data: CreateTaskData): Promise<{ message: string; data: Task }> {
    const response = await taskApi.post(API_ENDPOINTS.TASK.CREATE, data)
    return response.data
  },

  async updateTask(id: string, data: UpdateTaskData): Promise<{ message: string; data: Task }> {
    const response = await taskApi.put(API_ENDPOINTS.TASK.UPDATE(id), data)
    return response.data
  },

  async moveTask(id: string, data: MoveTaskData): Promise<{ message: string; data: Task }> {
    const response = await taskApi.put(`${API_ENDPOINTS.TASK.UPDATE(id)}/move`, data)
    return response.data
  },

  async deleteTask(id: string): Promise<{ message: string }> {
    const response = await taskApi.delete(API_ENDPOINTS.TASK.DELETE(id))
    return response.data
  },
}
