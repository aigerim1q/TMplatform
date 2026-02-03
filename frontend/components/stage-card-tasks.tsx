"use client"

import { useState, useEffect } from "react"
import { 
  Plus, 
  FileText, 
  Upload, 
  Loader2, 
  CheckCircle2, 
  Circle
} from "lucide-react"
import { fetchJson, getApiBaseUrl, getAccessToken } from "@/lib/api"

type Task = {
  id: number
  title: string
  status: string
  priority: number
  files?: any[] // Simplified
}

type Stage = {
  id: number
  title: string
  description?: string | null
  status: string
}

export default function StageCardWithTasks({ stage, projectName }: { stage: Stage, projectName: string }) {
  const [tasks, setTasks] = useState<Task[]>([])
  const [itemsLoading, setItemsLoading] = useState(false)
  
  // New Task State
  const [newTaskTitle, setNewTaskTitle] = useState("")
  const [isAdding, setIsAdding] = useState(false)

  // File Upload State
  const [uploadingTaskId, setUploadingTaskId] = useState<number | null>(null)

  useEffect(() => {
    loadTasks()
  }, [stage.id])

  async function loadTasks() {
    setItemsLoading(true)
    try {
      const data = await fetchJson<Task[]>(`/stages/${stage.id}/tasks`)
      setTasks(data || [])
    } catch (e) {
      console.error("Failed to load tasks", e)
    } finally {
      setItemsLoading(false)
    }
  }

  async function handleCreateTask(e: React.FormEvent) {
    e.preventDefault()
    if (!newTaskTitle.trim()) return

    try {
      const newTask = await fetchJson<Task>("/tasks", {
        method: "POST",
        body: JSON.stringify({
          stage_id: stage.id,
          title: newTaskTitle,
          priority: 1,
          description: "Created via frontend"
        })
      })
      setTasks([...tasks, newTask])
      setNewTaskTitle("")
      setIsAdding(false)
    } catch (e) {
      alert("Error creating task")
    }
  }

  async function handleFileUpload(taskId: number, file: File) {
    setUploadingTaskId(taskId)
    try {
      const formData = new FormData()
      formData.append("file", file)

      const token = getAccessToken()
      const res = await fetch(`${getApiBaseUrl()}/tasks/${taskId}/files`, {
        method: "POST",
        headers: token ? { "Authorization": `Bearer ${token}` } : {},
        body: formData
      })
      
      if (!res.ok) throw new Error("Upload failed")
      
      alert("File uploaded!")
      // In a real app, reload task details to see the file
    } catch (e) {
      alert("Upload error")
    } finally {
      setUploadingTaskId(null)
    }
  }

  return (
    <div className="min-w-[320px] max-w-[340px] shrink-0 space-y-4 rounded-3xl border border-gray-100 bg-[#FAFAFA] p-6 shadow-sm transition hover:shadow-md flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between">
        <span className="text-xs font-semibold text-gray-500 truncate max-w-[150px]">{projectName}</span>
        <span className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-bold ${
            stage.status === 'done' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'
        }`}>
            {stage.status}
        </span>
      </div>
      
      <div>
        <h3 className="text-lg font-bold leading-tight">{stage.title}</h3>
        <p className="mt-2 text-sm text-gray-500 line-clamp-2">
            {stage.description || "No description"}
        </p>
      </div>

      <hr className="border-gray-100" />

      {/* Tasks List */}
      <div className="flex-1 space-y-3 min-h-[100px]">
        {itemsLoading && <div className="text-sm text-gray-400">Loading tasks...</div>}
        
        {!itemsLoading && tasks.length === 0 && (
            <div className="text-sm text-gray-400 italic">Нет задач</div>
        )}

        {tasks.map(task => (
            <div key={task.id} className="group rounded-xl bg-white p-3 shadow-sm border border-gray-100">
                <div className="flex items-start justify-between gap-2">
                    <div className="flex gap-2">
                        {task.status === 'done' ? (
                            <CheckCircle2 className="h-4 w-4 text-green-500 mt-1" />
                        ) : (
                            <Circle className="h-4 w-4 text-gray-300 mt-1" />
                        )}
                        <span className="text-sm font-medium text-gray-800">{task.title}</span>
                    </div>
                </div>

                {/* File Upload Trigger */}
                <div className="mt-2 flex items-center justify-end">
                    <label className="cursor-pointer rounded-full p-1.5 hover:bg-gray-100 text-gray-400 hover:text-black transition">
                        {uploadingTaskId === task.id ? (
                            <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                            <Upload className="h-4 w-4" />
                        )}
                        <input 
                            type="file" 
                            className="hidden" 
                            onChange={(e) => {
                                if (e.target.files?.[0]) handleFileUpload(task.id, e.target.files[0])
                            }}
                        />
                    </label>
                </div>
            </div>
        ))}
      </div>

      {/* Add Task */}
      {isAdding ? (
          <form onSubmit={handleCreateTask} className="mt-4">
              <input 
                autoFocus
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm outline-none focus:border-black"
                placeholder="Task title..."
                value={newTaskTitle}
                onChange={e => setNewTaskTitle(e.target.value)}
              />
              <div className="mt-2 flex gap-2">
                  <button type="submit" className="flex-1 rounded-lg bg-black py-1.5 text-xs font-bold text-white">Save</button>
                  <button type="button" onClick={() => setIsAdding(false)} className="flex-1 rounded-lg bg-gray-100 py-1.5 text-xs font-bold text-gray-600">Cancel</button>
              </div>
          </form>
      ) : (
          <button 
            onClick={() => setIsAdding(true)}
            className="mt-2 flex w-full items-center justify-center gap-2 rounded-xl border border-dashed border-gray-300 py-2 text-sm font-medium text-gray-500 hover:border-black hover:text-black transition"
          >
            <Plus className="h-4 w-4" />
            Add Task
          </button>
      )}
    </div>
  )
}
