"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import Image from "next/image"
import { ArrowLeft, ArrowRight, Plus, Sparkles } from "lucide-react"
import Header from "@/components/header"
import { fetchJson } from "@/lib/api"

type Assignee = {
  id: number
  name: string
  email: string
  role: string
}

type Stage = {
  id: number
  title: string
}

type Project = {
  id: number
  name: string
  description?: string
  image_url?: string | null
  status: string
  start_date?: string
  end_date?: string
  assignees?: Assignee[]
}

type TaskSection = "urgent" | "my" | "subordinate" | null

type TaskSummary = {
  id: number
  title: string
  description?: string | null
  status: string
  priority: number
  due_date?: string | null
  stage_id: number
  stage_title: string
  project_id: number
  project_name: string
  project_image_url?: string | null
  assignee_name?: string | null
  author_name: string
}

export default function Dashboard() {
  const router = useRouter()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [urgentTasks, setUrgentTasks] = useState<TaskSummary[]>([])
  const [myTasks, setMyTasks] = useState<TaskSummary[]>([])
  const [subordinateTasks, setSubordinateTasks] = useState<TaskSummary[]>([])
  const [tasksLoading, setTasksLoading] = useState(false)
  const [activeTaskForm, setActiveTaskForm] = useState<TaskSection>(null)
  const [taskTitle, setTaskTitle] = useState("")
  const [taskDescription, setTaskDescription] = useState("")
  const [selectedProjectId, setSelectedProjectId] = useState<number | "">("")
  const [stages, setStages] = useState<Stage[]>([])
  const [selectedStageId, setSelectedStageId] = useState<number | "">("")
  const [assigneeId, setAssigneeId] = useState<number | "">("")
  const [dueDate, setDueDate] = useState("")
  const [priority, setPriority] = useState(1)
  const [submittingTask, setSubmittingTask] = useState(false)
  const [stageTitle, setStageTitle] = useState("")
  const [creatingStage, setCreatingStage] = useState(false)
  const [taskError, setTaskError] = useState<string | null>(null)
  const [taskSuccess, setTaskSuccess] = useState<string | null>(null)

  useEffect(() => {
    async function loadData() {
      try {
        // Add timestamp to prevent caching
        const res = await fetchJson<any>(`/projects?t=${Date.now()}`)
        console.log("Dashboard API response:", res)
        
        if (Array.isArray(res)) {
            setProjects(res)
        } else if (res && Array.isArray(res.items)) {
            setProjects(res.items)
        } else {
            console.warn("Unexpected projects response format:", res)
            setProjects([])
        }
      } catch (err) {
        console.error("Failed to fetch projects:", err)
      } finally {
        setLoading(false)
      }
    }
    loadData()
  }, [])

  const loadTasks = async () => {
    setTasksLoading(true)
    try {
      const [urgent, mine, sub] = await Promise.all([
        fetchJson<TaskSummary[]>("/tasks?scope=urgent"),
        fetchJson<TaskSummary[]>("/tasks?scope=my"),
        fetchJson<TaskSummary[]>("/tasks?scope=subordinate"),
      ])
      setUrgentTasks(Array.isArray(urgent) ? urgent : [])
      setMyTasks(Array.isArray(mine) ? mine : [])
      setSubordinateTasks(Array.isArray(sub) ? sub : [])
    } catch (err) {
      console.error("Failed to fetch tasks:", err)
      setUrgentTasks([])
      setMyTasks([])
      setSubordinateTasks([])
    } finally {
      setTasksLoading(false)
    }
  }

  useEffect(() => {
    loadTasks()
  }, [])

  useEffect(() => {
    async function loadStages(projectId: number) {
      try {
        const res = await fetchJson<any>(`/projects/${projectId}/stages`)
        const items = Array.isArray(res) ? res : res?.items || []
        setStages(items)
        if (items.length > 0) {
          setSelectedStageId(items[0].id)
        } else {
          setSelectedStageId("")
        }
      } catch (err) {
        console.error("Failed to fetch stages:", err)
        setStages([])
        setSelectedStageId("")
      }
    }

    if (selectedProjectId) {
      loadStages(Number(selectedProjectId))
    } else {
      setStages([])
      setSelectedStageId("")
    }
  }, [selectedProjectId])

  const calculateDaysLeft = (dateStr?: string) => {
    if (!dateStr) return null
    const end = new Date(dateStr)
    const now = new Date()
    const diffTime = end.getTime() - now.getTime()
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    return diffDays
  }

  const getTaskDaysLeft = (dateStr?: string | null) => {
    if (!dateStr) return null
    const end = new Date(dateStr)
    const now = new Date()
    const diffTime = end.getTime() - now.getTime()
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    return diffDays
  }

  const SectionHeader = ({ 
    title, 
    count, 
    color, 
    onAdd,
    isTask = false
  }: { 
    title: string, 
    count: number, 
    color: string, 
    onAdd: () => void,
    isTask?: boolean
  }) => (
    <div className="flex flex-wrap items-center gap-4">
      <div className={`h-4 w-4 rounded-full ${color}`} />
      <div className="rounded-full bg-black px-4 py-2 text-sm font-bold text-white transition hover:bg-neutral-800">
        {title}: {count}
      </div>
      <button 
        onClick={onAdd}
        className="flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-bold text-white transition hover:bg-neutral-800 hover:scale-105"
      >
        <Plus className="h-4 w-4" />
        {isTask ? "Добавить задачу" : "Добавить проект"}
      </button>
    </div>
  )

  const EmptyState = ({ text }: { text: string }) => (
    <div className="flex w-full items-center justify-center rounded-3xl border-2 border-dashed border-gray-200 p-12">
      <span className="rounded bg-blue-100 px-3 py-1 text-sm font-medium text-blue-800">
        {text}
      </span>
    </div>
  )

  const openTaskForm = (section: TaskSection) => {
    setActiveTaskForm(section)
    setTaskError(null)
    setTaskSuccess(null)
    setTaskTitle("")
    setTaskDescription("")
    setAssigneeId("")
    setDueDate("")
    setPriority(section === "urgent" ? 3 : 1)
    if (projects.length > 0) {
      setSelectedProjectId(projects[0].id)
    } else {
      setSelectedProjectId("")
    }
    setStageTitle("")
  }

  const activeProject = projects.find((p) => p.id === Number(selectedProjectId))
  const projectAssignees = activeProject?.assignees || []

  const TaskCreateCard = () => (
    <form
      onSubmit={async (e) => {
        e.preventDefault()
        setTaskError(null)
        setTaskSuccess(null)
        if (!selectedProjectId) {
          setTaskError("Выберите проект")
          return
        }
        if (!selectedStageId) {
          setTaskError("Выберите этап")
          return
        }
        if (!taskTitle.trim()) {
          setTaskError("Введите название задачи")
          return
        }

        try {
          setSubmittingTask(true)
          const payload: any = {
            stage_id: Number(selectedStageId),
            title: taskTitle.trim(),
            description: taskDescription.trim(),
            priority,
          }
          if (assigneeId) payload.assignee_id = Number(assigneeId)
          if (dueDate) payload.due_date = dueDate

          await fetchJson("/tasks", {
            method: "POST",
            body: JSON.stringify(payload),
          })

          setTaskSuccess("Задача создана")
          setTaskTitle("")
          setTaskDescription("")
          await loadTasks()
        } catch (err) {
          console.error(err)
          setTaskError("Не удалось создать задачу")
        } finally {
          setSubmittingTask(false)
        }
      }}
      className="w-full max-w-2xl rounded-3xl border border-gray-200 bg-white p-6 shadow-xl"
    >
      <div className="mb-4 flex items-center justify-between">
        <div className="text-sm font-semibold text-gray-700">Новая задача</div>
        <button
          type="button"
          onClick={() => setActiveTaskForm(null)}
          className="text-xs text-gray-400 hover:text-black"
        >
          Закрыть
        </button>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-2">
          <label className="text-xs font-semibold text-gray-600">Проект</label>
          <select
            value={selectedProjectId}
            onChange={(e) => setSelectedProjectId(e.target.value ? Number(e.target.value) : "")}
            className="w-full rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm"
          >
            {projects.length === 0 ? (
              <option value="">Нет проектов</option>
            ) : (
              projects.map((p) => (
                <option key={p.id} value={p.id}>{p.name}</option>
              ))
            )}
          </select>
        </div>
        <div className="space-y-2">
          <label className="text-xs font-semibold text-gray-600">Этап</label>
          {stages.length > 0 ? (
            <select
              value={selectedStageId}
              onChange={(e) => setSelectedStageId(e.target.value ? Number(e.target.value) : "")}
              className="w-full rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm"
            >
              {stages.map((s) => (
                <option key={s.id} value={s.id}>{s.title}</option>
              ))}
            </select>
          ) : (
            <div className="space-y-2">
              <div className="rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-700">
                В проекте нет этапов. Создайте первый этап.
              </div>
              <div className="flex gap-2">
                <input
                  value={stageTitle}
                  onChange={(e) => setStageTitle(e.target.value)}
                  className="w-full rounded-xl border border-gray-200 bg-white px-3 py-2 text-sm"
                  placeholder="Название этапа"
                />
                <button
                  type="button"
                  onClick={async () => {
                    if (!selectedProjectId) {
                      setTaskError("Выберите проект")
                      return
                    }
                    if (!stageTitle.trim()) {
                      setTaskError("Введите название этапа")
                      return
                    }
                    try {
                      setCreatingStage(true)
                      setTaskError(null)
                      const stage = await fetchJson<Stage>(`/projects/${selectedProjectId}/stages`, {
                        method: "POST",
                        body: JSON.stringify({ title: stageTitle.trim() }),
                      })
                      setStages([stage])
                      setSelectedStageId(stage.id)
                      setStageTitle("")
                    } catch (err) {
                      console.error(err)
                      setTaskError("Не удалось создать этап")
                    } finally {
                      setCreatingStage(false)
                    }
                  }}
                  className="shrink-0 rounded-xl bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-neutral-800 disabled:opacity-60"
                  disabled={creatingStage}
                >
                  {creatingStage ? "Создание..." : "Создать"}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="mt-4 space-y-2">
        <label className="text-xs font-semibold text-gray-600">Название</label>
        <input
          value={taskTitle}
          onChange={(e) => setTaskTitle(e.target.value)}
          className="w-full rounded-xl border border-gray-200 bg-white px-3 py-2 text-sm"
          placeholder="Например: Подготовить ТЗ"
        />
      </div>

      <div className="mt-4 space-y-2">
        <label className="text-xs font-semibold text-gray-600">Описание</label>
        <textarea
          value={taskDescription}
          onChange={(e) => setTaskDescription(e.target.value)}
          className="w-full rounded-xl border border-gray-200 bg-white px-3 py-2 text-sm"
          rows={3}
          placeholder="Кратко опишите задачу"
        />
      </div>

      <div className="mt-4 grid gap-4 md:grid-cols-3">
        <div className="space-y-2">
          <label className="text-xs font-semibold text-gray-600">Адресовано</label>
          {projectAssignees.length > 0 ? (
            <select
              value={assigneeId}
              onChange={(e) => setAssigneeId(e.target.value ? Number(e.target.value) : "")}
              className="w-full rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm"
            >
              <option value="">Вы</option>
              {projectAssignees.map((a) => (
                <option key={a.id} value={a.id}>{a.name || a.email}</option>
              ))}
            </select>
          ) : (
            <div className="rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-600">
              Вы
            </div>
          )}
        </div>
        <div className="space-y-2">
          <label className="text-xs font-semibold text-gray-600">Срок</label>
          <input
            type="date"
            value={dueDate}
            onChange={(e) => setDueDate(e.target.value)}
            className="w-full rounded-xl border border-gray-200 bg-white px-3 py-2 text-sm"
          />
        </div>
        <div className="space-y-2">
          <label className="text-xs font-semibold text-gray-600">Приоритет</label>
          <select
            value={priority}
            onChange={(e) => setPriority(Number(e.target.value))}
            className="w-full rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm"
          >
            <option value={1}>Низкий</option>
            <option value={2}>Средний</option>
            <option value={3}>Высокий</option>
          </select>
        </div>
      </div>

      {taskError && (
        <div className="mt-4 rounded-xl bg-red-50 px-4 py-2 text-sm text-red-600">{taskError}</div>
      )}
      {taskSuccess && (
        <div className="mt-4 rounded-xl bg-emerald-50 px-4 py-2 text-sm text-emerald-700">{taskSuccess}</div>
      )}

      <div className="mt-5 flex items-center gap-3">
        <button
          type="submit"
          disabled={submittingTask}
          className="rounded-full bg-black px-5 py-2 text-sm font-bold text-white transition hover:bg-neutral-800 disabled:opacity-60"
        >
          {submittingTask ? "Создание..." : "Создать задачу"}
        </button>
        <button
          type="button"
          onClick={() => setActiveTaskForm(null)}
          className="text-sm text-gray-500 hover:text-black"
        >
          Отмена
        </button>
      </div>
    </form>
  )

  const TaskCard = (task: TaskSummary) => {
    const daysLeft = getTaskDaysLeft(task.due_date || null)
    const badgeClass = daysLeft !== null && daysLeft < 0
      ? "bg-red-100 text-red-700"
      : "bg-emerald-100 text-emerald-700"
    const badgeText = daysLeft === null
      ? "Без срока"
      : daysLeft < 0
        ? `${daysLeft} дн. просрочка`
        : `${daysLeft} дн. осталось`

    return (
      <div key={task.id} className="min-w-[280px] max-w-[320px] shrink-0 rounded-2xl border border-gray-100 bg-white p-4 shadow-sm">
        <div className="flex items-start justify-between gap-2">
          <div className="text-xs text-gray-500">
            Проект: <span className="font-semibold text-gray-700">{task.project_name}</span>
          </div>
          <span className={`rounded-full px-2.5 py-0.5 text-[11px] font-semibold ${badgeClass}`}>{badgeText}</span>
        </div>
        <div className="mt-2 text-xs text-gray-400">Этап: {task.stage_title}</div>
        <div className="mt-2 text-sm font-semibold text-gray-900 line-clamp-2">{task.title}</div>
        {task.description && (
          <div className="mt-1 text-xs text-gray-500 line-clamp-2">{task.description}</div>
        )}
        <div className="mt-3 flex items-center justify-between">
          <div className="text-xs text-gray-500">Адресовано: {task.assignee_name || "Вы"}</div>
          <button
            onClick={() => router.push(`/project/${task.project_id}`)}
            className="rounded-full bg-black px-3 py-1.5 text-xs font-semibold text-white"
          >
            Открыть
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-white text-black transition-colors">
      <Header />
      <main className="container mx-auto max-w-7xl space-y-12 p-4 md:p-8">
        
        {/* Urgent Tasks - Static Empty */}
        <section className="space-y-6">
          <SectionHeader 
            title="Срочные задачи" 
            count={urgentTasks.length} 
            color="bg-green-500" 
            onAdd={() => openTaskForm("urgent")} 
            isTask
          />
          {activeTaskForm === "urgent" && (
            <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
              <TaskCreateCard />
            </div>
          )}
          {tasksLoading && (
            <div className="text-sm text-gray-400">Загрузка задач...</div>
          )}
          {!tasksLoading && urgentTasks.length === 0 && (
            <EmptyState text="Нет срочных задач" />
          )}
          {!tasksLoading && urgentTasks.length > 0 && (
            <div className="flex gap-4 overflow-x-auto pb-4 items-stretch scrollbar-hide">
              {urgentTasks.map((task) => TaskCard(task))}
            </div>
          )}
        </section>

        {/* My Tasks - Static Empty */}
        <section className="space-y-6">
          <SectionHeader 
            title="Мои задачи" 
            count={myTasks.length} 
            color="bg-red-500" 
            onAdd={() => openTaskForm("my")} 
            isTask
          />
          {activeTaskForm === "my" && (
            <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
              <TaskCreateCard />
            </div>
          )}
          {tasksLoading && (
            <div className="text-sm text-gray-400">Загрузка задач...</div>
          )}
          {!tasksLoading && myTasks.length === 0 && (
            <EmptyState text="Нет задач" />
          )}
          {!tasksLoading && myTasks.length > 0 && (
            <div className="flex gap-4 overflow-x-auto pb-4 items-stretch scrollbar-hide">
              {myTasks.map((task) => TaskCard(task))}
            </div>
          )}
        </section>

        {/* Projects - Fetched */}
        <section className="space-y-6">
          <SectionHeader 
            title="Проекты" 
            count={projects.length} 
            color="bg-yellow-500" 
            onAdd={() => router.push("/project/new")} 
          />
          
           <div className="flex gap-6 overflow-x-auto pb-4 items-stretch scrollbar-hide">
             {projects.map((project) => {
                const daysLeft = calculateDaysLeft(project.end_date)
                return (
                  <div
                    key={project.id}
                    onClick={() => router.push(`/project/${project.id}`)}
                    className="group relative min-w-[280px] max-w-[320px] shrink-0 overflow-hidden rounded-3xl bg-white shadow-lg transition-all hover:shadow-xl border border-neutral-100 cursor-pointer"
                  >
                     <div className="relative h-48 w-full overflow-hidden rounded-t-3xl bg-neutral-100">
                        {project.image_url ? (
                          <Image src={project.image_url} alt={project.name} fill unoptimized className="object-cover transition duration-700 group-hover:scale-105" />
                        ) : (
                           <div className="flex h-full items-center justify-center bg-stone-100 text-stone-300">
                             <span className="text-6xl font-black opacity-10 uppercase">{project.name.substring(0, 2)}</span>
                           </div>
                        )}
                        <div className="absolute left-4 top-4 rounded-full bg-black/70 px-4 py-1.5 text-xs font-bold text-white backdrop-blur-md">
                           {project.name}
                        </div>
                        {daysLeft !== null && (
                          <div className={`absolute right-4 top-4 rounded-full px-4 py-1.5 text-xs font-bold backdrop-blur-md ${daysLeft < 0 ? 'bg-red-500/90 text-white' : 'bg-emerald-500/90 text-white'}`}>
                            {daysLeft < 0 ? `${daysLeft} дн. просрочка` : `${daysLeft} дн. осталось`}
                          </div>
                        )}
                     </div>
                     
                     <div className="p-5">
                       <h3 className="mb-2 text-xl font-bold tracking-tight text-black">
                         {project.name}
                       </h3>
                       <p className="mb-6 text-sm text-gray-500 line-clamp-2 min-h-[40px]">
                          {project.description || "Нет описания"}
                       </p>
                     </div>
                  </div>
                )
             })}
          </div>
          
          {loading && (
            <div className="py-20 text-center text-gray-400">Loading projects...</div>
          )}

          {!loading && projects.length === 0 && (
             <EmptyState text="Нет проектов" />
          )}
        </section>

        {/* Subordinate Tasks - Static Empty */}
        <section className="space-y-6">
           <SectionHeader 
             title="Задачи подчинённых" 
             count={subordinateTasks.length} 
             color="bg-pink-500" 
             onAdd={() => openTaskForm("subordinate")} 
             isTask
           />
           {activeTaskForm === "subordinate" && (
             <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
               <TaskCreateCard />
             </div>
           )}
           {tasksLoading && (
             <div className="text-sm text-gray-400">Загрузка задач...</div>
           )}
           {!tasksLoading && subordinateTasks.length === 0 && (
             <EmptyState text="Нет задач подчинённых" />
           )}
           {!tasksLoading && subordinateTasks.length > 0 && (
             <div className="flex gap-4 overflow-x-auto pb-4 items-stretch scrollbar-hide">
               {subordinateTasks.map((task) => TaskCard(task))}
             </div>
           )}
        </section>

        {/* Subordinate Projects - Static Empty */}
        <section className="space-y-6">
           <SectionHeader 
             title="Проекты подчинённых" 
             count={0} 
             color="bg-amber-500" 
             onAdd={() => {}} 
           />
           <EmptyState text="Нет проектов подчинённых" />
        </section>

        <button className="fixed bottom-8 right-8 z-50 flex items-center gap-3 rounded-full bg-[#1a1a1a] px-6 py-4 text-white shadow-2xl transition hover:scale-105 hover:bg-black ring-4 ring-white/10">
          <Sparkles className="h-5 w-5 text-yellow-400 fill-yellow-400" />
          <span className="font-bold tracking-wide">AI: предложить план</span>
        </button>

      </main>
    </div>
  )
}
