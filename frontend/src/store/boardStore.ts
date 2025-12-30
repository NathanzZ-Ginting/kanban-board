import { create } from 'zustand'
import { Board } from '@/types'

interface BoardState {
  boards: Board[]
  currentBoard: Board | null
  setBoards: (boards: Board[]) => void
  setCurrentBoard: (board: Board | null) => void
  addBoard: (board: Board) => void
  updateBoard: (id: string, board: Partial<Board>) => void
  deleteBoard: (id: string) => void
}

export const useBoardStore = create<BoardState>((set) => ({
  boards: [],
  currentBoard: null,
  setBoards: (boards) => set({ boards }),
  setCurrentBoard: (board) => set({ currentBoard: board }),
  addBoard: (board) =>
    set((state) => ({ boards: [...state.boards, board] })),
  updateBoard: (id, boardData) =>
    set((state) => ({
      boards: state.boards.map((b) =>
        b.id === id ? { ...b, ...boardData } : b
      ),
      currentBoard:
        state.currentBoard?.id === id
          ? { ...state.currentBoard, ...boardData }
          : state.currentBoard,
    })),
  deleteBoard: (id) =>
    set((state) => ({
      boards: state.boards.filter((b) => b.id !== id),
      currentBoard: state.currentBoard?.id === id ? null : state.currentBoard,
    })),
}))
