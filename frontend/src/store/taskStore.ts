import { create } from 'zustand'
import { Task } from '@/types'

interface TaskState {
  tasks: Task[]
  setTasks: (tasks: Task[]) => void
  addTask: (task: Task) => void
  updateTask: (id: string, task: Partial<Task>) => void
  deleteTask: (id: string) => void
  getTasksByColumn: (columnId: number) => Task[]
  getTasksByBoard: (boardId: number) => Task[]
}

export const useTaskStore = create<TaskState>((set, get) => ({
  tasks: [],
  setTasks: (tasks) => set({ tasks }),
  addTask: (task) => set((state) => ({ tasks: [...state.tasks, task] })),
  updateTask: (id, taskData) =>
    set((state) => ({
      tasks: state.tasks.map((t) => (t.id === id ? { ...t, ...taskData } : t)),
    })),
  deleteTask: (id) =>
    set((state) => ({ tasks: state.tasks.filter((t) => t.id !== id) })),
  getTasksByColumn: (columnId) => {
    return get().tasks.filter((t) => t.columnId === columnId)
  },
  getTasksByBoard: (boardId) => {
    return get().tasks.filter((t) => t.boardId === boardId)
  },
}))
