export interface User {
  id: string
  email: string
  username: string
  firstName: string
  lastName: string
  avatar?: string
  role: 'admin' | 'user'
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Column {
  id: number
  boardId: number
  name: string
  position: number
  color: string
  createdAt: string
  updatedAt: string
}

export interface Board {
  id: string
  name: string
  description?: string
  ownerId: string
  isPublic?: boolean
  color?: string
  columns?: Column[]
  members?: BoardMember[]
  createdAt: string
  updatedAt: string
}

export interface BoardMember {
  id: number
  boardId: number
  userId: number
  role: 'owner' | 'admin' | 'member'
  createdAt: string
  updatedAt: string
}

export interface Task {
  id: string
  title: string
  description?: string
  boardId: number
  columnId: number
  position: number
  assigneeId?: number
  creatorId: number
  priority: 'low' | 'medium' | 'high' | 'urgent'
  dueDate?: string
  labels?: Label[]
  comments?: Comment[]
  createdAt: string
  updatedAt: string
}

export interface Label {
  id: number
  boardId: number
  name: string
  color: string
}

export interface Comment {
  id: number
  taskId: number
  userId: number
  content: string
  createdAt: string
  updatedAt: string
}

export interface AuthTokens {
  accessToken: string
  refreshToken: string
}

export interface ApiResponse<T = any> {
  data?: T
  error?: string
  message?: string
}
