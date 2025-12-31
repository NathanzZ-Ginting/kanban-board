'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import {
  ArrowLeft,
  Plus,
  MoreHorizontal,
  Clock,
  X,
  Loader2,
  GripVertical,
} from 'lucide-react';
import { useAuthStore } from '@/store/authStore';
import { useBoardStore } from '@/store/boardStore';
import { useTaskStore } from '@/store/taskStore';
import { boardService } from '@/services/boardService';
import { taskService, CreateTaskData } from '@/services/taskService';
import { Task, Column } from '@/types';

const priorityColors: Record<string, string> = {
  high: 'bg-red-100 text-red-700',
  medium: 'bg-yellow-100 text-yellow-700',
  low: 'bg-green-100 text-green-700',
  urgent: 'bg-purple-100 text-purple-700',
};

const defaultColumnColors: string[] = [
  'bg-gray-400',
  'bg-blue-500',
  'bg-yellow-500',
  'bg-green-500',
  'bg-purple-500',
  'bg-pink-500',
];

interface BoardPageProps {
  params: { id: string };
}

export default function BoardPage({ params }: BoardPageProps) {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const { currentBoard, setCurrentBoard } = useBoardStore();
  const { tasks, setTasks, addTask } = useTaskStore();

  const [columns, setColumns] = useState<Column[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Create Task Modal State
  const [isCreateTaskModalOpen, setIsCreateTaskModalOpen] = useState(false);
  const [isCreatingTask, setIsCreatingTask] = useState(false);
  const [selectedColumnId, setSelectedColumnId] = useState<number | null>(null);
  const [newTaskTitle, setNewTaskTitle] = useState('');
  const [newTaskDescription, setNewTaskDescription] = useState('');
  const [newTaskPriority, setNewTaskPriority] = useState<Task['priority']>('medium');
  const [newTaskDueDate, setNewTaskDueDate] = useState('');

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, router]);

  // Fetch board and tasks from API
  useEffect(() => {
    const fetchData = async () => {
      if (!isAuthenticated || !params.id) return;
      
      try {
        setIsLoading(true);
        setError(null);
        
        // Fetch board details (includes columns)
        const boardData = await boardService.getBoard(params.id);
        if (boardData) {
          setCurrentBoard(boardData);
          // Sort columns by position
          const sortedColumns = (boardData.columns || []).sort(
            (a, b) => a.position - b.position
          );
          setColumns(sortedColumns);
          
          // Set default selected column for task creation
          if (sortedColumns.length > 0 && !selectedColumnId) {
            setSelectedColumnId(sortedColumns[0].id);
          }
        }
        
        // Fetch tasks for this board
        const tasksResponse = await taskService.getTasksByBoard(params.id);
        if (tasksResponse.tasks) {
          setTasks(tasksResponse.tasks);
        }
      } catch (err) {
        console.error('Failed to fetch data:', err);
        setError('Failed to load board data');
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, [isAuthenticated, params.id, setCurrentBoard, setTasks]);

  const openCreateTaskModal = (columnId: number) => {
    setSelectedColumnId(columnId);
    setIsCreateTaskModalOpen(true);
  };

  const closeCreateTaskModal = () => {
    setIsCreateTaskModalOpen(false);
    setNewTaskTitle('');
    setNewTaskDescription('');
    setNewTaskPriority('medium');
    setNewTaskDueDate('');
  };

  const handleCreateTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTaskTitle.trim() || !selectedColumnId) return;

    try {
      setIsCreatingTask(true);
      
      const taskData: CreateTaskData = {
        title: newTaskTitle.trim(),
        description: newTaskDescription.trim() || '',
        boardId: parseInt(params.id),
        columnId: selectedColumnId,
        priority: newTaskPriority,
      };
      
      // Only add dueDate if it has a value
      if (newTaskDueDate) {
        taskData.dueDate = newTaskDueDate;
      }

      console.log('Creating task with data:', taskData);
      const response = await taskService.createTask(taskData);

      if (response.data) {
        addTask(response.data);
        closeCreateTaskModal();
      }
    } catch (err) {
      console.error('Failed to create task:', err);
      setError('Failed to create task');
    } finally {
      setIsCreatingTask(false);
    }
  };

  // Group tasks by column
  const getTasksByColumnId = (columnId: number) => {
    return tasks.filter(task => task.columnId === columnId).sort((a, b) => a.position - b.position);
  };

  const getColumnColor = (index: number, columnColor?: string) => {
    if (columnColor && columnColor !== '#e5e7eb') {
      return `bg-[${columnColor}]`;
    }
    return defaultColumnColors[index % defaultColumnColors.length];
  };

  if (!isAuthenticated) {
    return null;
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50 flex items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-4">{error}</p>
          <Link href="/dashboard" className="text-indigo-600 hover:underline">
            Back to Dashboard
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-20">
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Link
                href="/dashboard"
                className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <ArrowLeft className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="text-xl font-bold text-gray-900">
                  {currentBoard?.name || 'Untitled Board'}
                </h1>
                <p className="text-sm text-gray-500">
                  {currentBoard?.description || 'No description'}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-3">
              <button 
                onClick={() => columns.length > 0 ? openCreateTaskModal(columns[0].id) : null}
                disabled={columns.length === 0}
                className="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Plus className="w-4 h-4" />
                Add Task
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Kanban Board */}
      <div className="p-6 overflow-x-auto">
        <div className="flex gap-6 min-w-max">
          {columns.length === 0 ? (
            <div className="text-center py-12 w-full">
              <p className="text-gray-500 mb-4">No columns yet. Create your first column to get started!</p>
            </div>
          ) : (
            columns.map((column, index) => {
              const columnTasks = getTasksByColumnId(column.id);
              
              return (
                <div
                  key={column.id}
                  className="w-80 flex-shrink-0"
                >
                  {/* Column Header */}
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className={`w-3 h-3 rounded-full ${getColumnColor(index, column.color)}`} />
                      <h2 className="font-semibold text-gray-900">{column.name}</h2>
                      <span className="text-sm text-gray-500 bg-gray-100 px-2 py-0.5 rounded-full">
                        {columnTasks.length}
                      </span>
                    </div>
                    <button className="p-1 text-gray-400 hover:text-gray-600 rounded">
                      <MoreHorizontal className="w-5 h-5" />
                    </button>
                  </div>

                  {/* Tasks */}
                  <div className="space-y-3">
                    {columnTasks.map((task) => (
                      <div
                        key={task.id}
                        className="bg-white rounded-xl p-4 shadow-sm border border-gray-100 hover:shadow-md hover:border-indigo-200 transition-all cursor-pointer group"
                      >
                        {/* Drag Handle */}
                        <div className="flex items-start gap-2">
                          <GripVertical className="w-4 h-4 text-gray-300 opacity-0 group-hover:opacity-100 transition-opacity mt-0.5 cursor-grab" />
                          <div className="flex-1">
                            {/* Title */}
                            <h3 className="font-medium text-gray-900 mb-1">{task.title}</h3>

                            {/* Description */}
                            {task.description && (
                              <p className="text-sm text-gray-500 mb-3 line-clamp-2">
                                {task.description}
                              </p>
                            )}

                            {/* Footer */}
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-3 text-xs text-gray-500">
                                {task.dueDate && (
                                  <span className="flex items-center gap-1">
                                    <Clock className="w-3.5 h-3.5" />
                                    {new Date(task.dueDate).toLocaleDateString('en-US', {
                                      month: 'short',
                                      day: 'numeric',
                                    })}
                                  </span>
                                )}
                                <span
                                  className={`px-2 py-0.5 rounded-full ${priorityColors[task.priority] || 'bg-gray-100 text-gray-700'}`}
                                >
                                  {task.priority}
                                </span>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}

                    {/* Add Task Button */}
                    <button 
                      onClick={() => openCreateTaskModal(column.id)}
                      className="w-full flex items-center justify-center gap-2 py-3 text-sm text-gray-500 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl border-2 border-dashed border-gray-200 hover:border-indigo-300 transition-all"
                    >
                      <Plus className="w-4 h-4" />
                      Add Task
                    </button>
                  </div>
                </div>
              );
            })
          )}

          {/* Add Column Button */}
          <div className="w-80 flex-shrink-0">
            <button className="w-full flex items-center justify-center gap-2 py-4 text-sm font-medium text-gray-500 hover:text-indigo-600 bg-white/50 hover:bg-indigo-50 rounded-xl border-2 border-dashed border-gray-200 hover:border-indigo-300 transition-all">
              <Plus className="w-4 h-4" />
              Add Column
            </button>
          </div>
        </div>
      </div>

      {/* Create Task Modal */}
      {isCreateTaskModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div 
            className="absolute inset-0 bg-black/50" 
            onClick={closeCreateTaskModal}
          />
          <div className="relative bg-white rounded-xl shadow-xl w-full max-w-md mx-4 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Create New Task</h2>
              <button
                onClick={closeCreateTaskModal}
                className="p-1 text-gray-400 hover:text-gray-600 rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleCreateTask}>
              <div className="space-y-4">
                <div>
                  <label htmlFor="taskTitle" className="block text-sm font-medium text-gray-700 mb-1">
                    Task Title *
                  </label>
                  <input
                    id="taskTitle"
                    type="text"
                    value={newTaskTitle}
                    onChange={(e) => setNewTaskTitle(e.target.value)}
                    placeholder="Enter task title"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                    required
                  />
                </div>

                <div>
                  <label htmlFor="taskDescription" className="block text-sm font-medium text-gray-700 mb-1">
                    Description
                  </label>
                  <textarea
                    id="taskDescription"
                    value={newTaskDescription}
                    onChange={(e) => setNewTaskDescription(e.target.value)}
                    placeholder="Enter task description (optional)"
                    rows={3}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label htmlFor="taskPriority" className="block text-sm font-medium text-gray-700 mb-1">
                      Priority
                    </label>
                    <select
                      id="taskPriority"
                      value={newTaskPriority}
                      onChange={(e) => setNewTaskPriority(e.target.value as Task['priority'])}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                    >
                      <option value="low">Low</option>
                      <option value="medium">Medium</option>
                      <option value="high">High</option>
                      <option value="urgent">Urgent</option>
                    </select>
                  </div>

                  <div>
                    <label htmlFor="taskColumn" className="block text-sm font-medium text-gray-700 mb-1">
                      Column
                    </label>
                    <select
                      id="taskColumn"
                      value={selectedColumnId || ''}
                      onChange={(e) => setSelectedColumnId(parseInt(e.target.value))}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                    >
                      {columns.map((col) => (
                        <option key={col.id} value={col.id}>
                          {col.name}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div>
                  <label htmlFor="taskDueDate" className="block text-sm font-medium text-gray-700 mb-1">
                    Due Date
                  </label>
                  <input
                    id="taskDueDate"
                    type="date"
                    value={newTaskDueDate}
                    onChange={(e) => setNewTaskDueDate(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  />
                </div>
              </div>

              <div className="flex justify-end gap-3 mt-6">
                <button
                  type="button"
                  onClick={closeCreateTaskModal}
                  className="px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isCreatingTask || !newTaskTitle.trim()}
                  className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  {isCreatingTask && <Loader2 className="w-4 h-4 animate-spin" />}
                  Create Task
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
