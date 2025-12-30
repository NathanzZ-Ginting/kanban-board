import { boardApi } from '@/utils/axios'
import { API_ENDPOINTS } from '@/config/api.config'
import { Board, ApiResponse } from '@/types'

export interface CreateBoardData {
  name: string
  description?: string
}

export const boardService = {
  async getBoards(): Promise<ApiResponse<Board[]>> {
    const response = await boardApi.get(API_ENDPOINTS.BOARD.LIST)
    return response.data
  },

  async getBoard(id: string): Promise<ApiResponse<Board>> {
    const response = await boardApi.get(API_ENDPOINTS.BOARD.GET(id))
    return response.data
  },

  async createBoard(data: CreateBoardData): Promise<ApiResponse<Board>> {
    const response = await boardApi.post(API_ENDPOINTS.BOARD.CREATE, data)
    return response.data
  },

  async updateBoard(id: string, data: Partial<Board>): Promise<ApiResponse<Board>> {
    const response = await boardApi.put(API_ENDPOINTS.BOARD.UPDATE(id), data)
    return response.data
  },

  async deleteBoard(id: string): Promise<ApiResponse> {
    const response = await boardApi.delete(API_ENDPOINTS.BOARD.DELETE(id))
    return response.data
  },
}
