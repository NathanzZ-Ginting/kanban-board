import { boardApi } from '@/utils/axios'
import { API_ENDPOINTS } from '@/config/api.config'
import { Board, Column } from '@/types'

export interface CreateBoardData {
  name: string
  description?: string
  isPublic?: boolean
  color?: string
}

export interface BoardsListResponse {
  boards: Board[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export interface CreateColumnData {
  name: string
  position?: number
  color?: string
}

export const boardService = {
  async getBoards(): Promise<BoardsListResponse> {
    const response = await boardApi.get(API_ENDPOINTS.BOARD.LIST)
    return response.data
  },

  async getBoard(id: string): Promise<Board> {
    const response = await boardApi.get(API_ENDPOINTS.BOARD.GET(id))
    return response.data
  },

  async createBoard(data: CreateBoardData): Promise<{ message: string; data: Board }> {
    const response = await boardApi.post(API_ENDPOINTS.BOARD.CREATE, data)
    return response.data
  },

  async updateBoard(id: string, data: Partial<Board>): Promise<{ message: string; data: Board }> {
    const response = await boardApi.put(API_ENDPOINTS.BOARD.UPDATE(id), data)
    return response.data
  },

  async deleteBoard(id: string): Promise<{ message: string }> {
    const response = await boardApi.delete(API_ENDPOINTS.BOARD.DELETE(id))
    return response.data
  },

  async createColumn(boardId: string, data: CreateColumnData): Promise<{ message: string; data: Column }> {
    const response = await boardApi.post(`/api/v1/boards/${boardId}/columns`, data)
    return response.data
  },
}
