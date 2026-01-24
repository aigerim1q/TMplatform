"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import Image from "next/image"
import { 
  ArrowLeft, 
  ArrowRight, 
  Plus, 
  Maximize2,
  Clock,
  CheckCircle2,
  Sparkles
} from "lucide-react"
import { fetchJson } from "@/lib/api"
import StageCardWithTasks from "@/components/stage-card-tasks"

type Assignee = {
  id: number
  name: string
  email: string
  role: string
}

type Project = {
  id: number
  name: string
  description?: string | null
  image_url?: string | null
  image?: string 
  status: string
  start_date?: string | null
  end_date?: string | null
  priority?: number
  budget_allocated: number
  budget_spent: number
  budget_currency: string
  assignees: Assignee[]
}

type Stage = {
  id: number
  title: string
  description?: string | null
  status: "todo" | "in_progress" | "done"
  start_date?: string | null
  end_date?: string | null
  order_index: number
}

// Mock phases for the "Design" if real grouping isn't available
const PHASES = [
  { id: 1, title: "I. Предпроектная подготовка (до начала проектирования)", delay: "5 дней" },
  { id: 2, title: "II. Проектирование (архитектура, инженерия, дизайн)", delay: null },
  { id: 3, title: "III. Строительный этап", delay: null }
]

export default function ProjectDetail() {
  const router = useRouter()
  const params = useParams()
  // Ensure we get a string ID safely
  const projectId = Array.isArray(params.id) ? params.id[0] : params.id
  
  const [project, setProject] = useState<Project | null>(null)
  const [stages, setStages] = useState<Stage[]>([])
  const [loading, setLoading] = useState(true)
  const [stageModalOpen, setStageModalOpen] = useState(false)
  const [stageTitle, setStageTitle] = useState("")
  const [stageDescription, setStageDescription] = useState("")
  const [stageSubmitting, setStageSubmitting] = useState(false)
  const [stageError, setStageError] = useState<string | null>(null)

  useEffect(() => {
    async function load() {
      if (!projectId) return
      try {
        const p = await fetchJson<Project>(`/projects/${projectId}`)
        setProject(p)
        const s = await fetchJson<Stage[]>(`/projects/${projectId}/stages`)
        if (Array.isArray(s)) setStages(s)
      } catch (e) {
        console.error(e)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [projectId])

  if (loading || !project) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-white">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-black border-t-transparent" />
      </div>
    )
  }

  // Helper to distribute stages into phases roughly (for visual demo)
  const getPhaseStages = (phaseIndex: number) => {
    const stagesPerPhase = Math.ceil(stages.length / 3) || 1
    if (stages.length === 0) return []
    const start = phaseIndex * stagesPerPhase
    return stages.slice(start, start + stagesPerPhase)
  }

  const handleCreateStage = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!projectId) return
    if (!stageTitle.trim()) {
      setStageError("Введите название этапа")
      return
    }
    try {
      setStageSubmitting(true)
      setStageError(null)
      const stage = await fetchJson<Stage>(`/projects/${projectId}/stages`, {
        method: "POST",
        body: JSON.stringify({
          title: stageTitle.trim(),
          description: stageDescription.trim() || undefined,
        }),
      })
      setStages([...stages, stage])
      setStageTitle("")
      setStageDescription("")
      setStageModalOpen(false)
    } catch (err) {
      console.error(err)
      setStageError("Не удалось создать этап")
    } finally {
      setStageSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen bg-white pb-20 text-black">
      {/* Top Navigation Bar */}
      <div className="sticky top-0 z-50 flex items-center justify-between border-b border-gray-100 bg-white/80 px-6 py-4 backdrop-blur-md">
        <button 
          onClick={() => router.back()}
          className="flex items-center gap-2 rounded-full bg-black px-6 py-2.5 text-sm font-bold text-white transition hover:bg-neutral-800"
        >
          <ArrowLeft className="h-4 w-4" />
          Назад
        </button>
        
        <div className="flex gap-2">
          <button className="rounded-full bg-black px-8 py-2.5 text-sm font-bold text-white shadow-lg transition hover:bg-neutral-800">
            Задача
          </button>
          <button className="rounded-full bg-gray-500 px-8 py-2.5 text-sm font-bold text-white transition hover:bg-gray-600">
            Отчеты
          </button>
        </div>
        <div className="w-24" /> {/* Spacer for balance */}
      </div>

      <div className="container mx-auto max-w-6xl space-y-8 p-6">
        
        {/* Project Header Card */}
        <div className="relative overflow-hidden rounded-[2.5rem] bg-white">
          <div className="flex flex-col gap-6 md:flex-row md:items-start">
             {/* Thumbnail */}
            <div className="relative h-40 w-40 shrink-0 overflow-hidden rounded-2xl bg-gray-200 shadow-md">
               {project.image || project.image_url ? (
                 <Image 
                   src={project.image || project.image_url || ""} 
                   alt={project.name}
                   fill
                   className="object-cover"
                 />
               ) : (
                <div className="flex h-full w-full items-center justify-center bg-gray-100 text-4xl font-bold text-gray-400">
                  {project.name[0]}
                </div>
               )}
            </div>

            <div className="flex-1 space-y-6">
              <h1 className="text-4xl font-black tracking-tight text-black">
                Проект: {project.name}
              </h1>

              {/* Finance Bar */}
              <div className="flex flex-wrap items-center justify-between gap-4 rounded-3xl bg-[#A3E635] p-2 pl-6 pr-2 shadow-sm">
                <span className="font-bold text-black/80">
                  Финансы: <span className="font-black text-black">{project.budget_spent.toLocaleString()} / {project.budget_allocated.toLocaleString()}</span>
                </span>
                <button className="rounded-full bg-black px-6 py-2 text-sm font-bold text-white transition hover:bg-neutral-800">
                   Отчеты расходов
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Info Grid */}
        <div className="grid gap-4 md:grid-cols-12">
           {/* Deadline Card */}
           <div className="col-span-12 flex flex-col justify-center gap-1 rounded-3xl bg-[#EAD8B1] p-6 text-[#4A3B18] md:col-span-4">
              <div className="flex items-center gap-3">
                 <Clock className="h-5 w-5" />
                 <span className="font-semibold">Дедлайн: {project.end_date ? new Date(project.end_date).toLocaleDateString() : 'Не указан'}</span>
              </div>
              <div className="pl-8 text-sm opacity-80">
                 {project.start_date && `Дата начала: ${new Date(project.start_date).toLocaleDateString()}`}
              </div>
           </div>

           {/* Team Card */}
           <div className="col-span-12 flex items-center justify-between rounded-3xl border border-gray-200 bg-white p-6 shadow-sm transition hover:shadow-md md:col-span-5">
              <div className="flex items-center gap-4">
                 <div className="flex -space-x-3">
                    {[1,2,3].map(i => (
                      <div key={i} className="flex h-10 w-10 items-center justify-center rounded-full border-2 border-white bg-gray-100 text-xs font-bold text-gray-600">
                        U{i}
                      </div>
                    ))}
                 </div>
                 <div className="text-sm font-medium leading-tight text-gray-600">
                    Ответственные: <br/>
                    <span className="text-black">Омар Ахмет, Зейнулла...</span>
                 </div>
              </div>
              <Maximize2 className="h-5 w-5 text-gray-400" />
           </div>

           {/* Alert Card */}
           <div className="col-span-12 flex items-center justify-between rounded-3xl bg-black p-6 text-white md:col-span-3">
              <div className="flex items-center gap-3">
                 <div className="h-3 w-3 rounded-full bg-red-500 animate-pulse" />
                 <span className="text-sm font-medium leading-tight">Причины отсрочки по всем дедлайнам</span>
              </div>
              <Maximize2 className="h-4 w-4 text-white/50" />
           </div>
        </div>

        {/* Phases & Stages */}
        <div className="space-y-10">
          {PHASES.map((phase, idx) => {
             const phaseStages = getPhaseStages(idx)
             
             return (
              <div key={phase.id} className="space-y-6">
                {/* Phase Header */}
                <div className="flex flex-wrap items-center justify-between gap-4">
                   <div className="flex items-center gap-4 rounded-full bg-black py-3 pl-6 pr-3 text-white max-sm:w-full max-sm:justify-between">
                      <span className="font-bold">{phase.title}</span>
                      {phase.delay && (
                         <span className="flex items-center gap-2 rounded-full bg-[#FDA4AF] px-3 py-1 text-xs font-bold text-[#881337]">
                            <Clock className="h-3 w-3" /> {phase.delay}
                         </span>
                      )}
                   </div>
                   
                   <button className="flex items-center gap-2 rounded-full bg-black px-5 py-3 text-sm font-bold text-white transition hover:bg-neutral-800">
                      <Plus className="h-4 w-4" />
                      Добавить задачу
                   </button>
                </div>

                {/* Horizontal Scroll / Grid of Stages */}
                <div className="flex gap-4 overflow-x-auto pb-4 items-stretch scrollbar-hide">
                   {phaseStages.length > 0 ? phaseStages.map((stage, sIdx) => (
                      <StageCardWithTasks 
                        key={stage.id} 
                        stage={stage} 
                        projectName={project.name} 
                      />
                   )) : (
                     <div className="flex w-full items-center justify-center rounded-2xl border-2 border-dashed border-gray-100 py-8">
                       <span className="text-gray-400">Нет этапов</span>
                     </div>
                   )}
                   
                   {/* Add More Arrow Card - Visual Only */}
                   {phaseStages.length > 0 && (
                      <button className="flex h-auto min-w-[60px] items-center justify-center rounded-full bg-black transition hover:scale-105">
                         <ArrowRight className="h-8 w-8 text-white" />
                      </button>
                   )}
                </div>
              </div>
             )
          })}
        </div>

        {/* Bottom Actions */}
          <div className="flex justify-center pt-8">
            <button
             onClick={() => {
              setStageError(null)
              setStageModalOpen(true)
             }}
             className="flex w-full max-w-md items-center justify-center gap-2 rounded-full bg-black py-4 font-bold text-white shadow-xl transition hover:bg-neutral-800"
            >
              <Plus className="h-5 w-5" />
              Добавить этап проекта
           </button>
        </div>
      </div>

      {/* AI Floating Button */}
      <button className="fixed bottom-8 right-8 z-50 flex items-center gap-2 rounded-full bg-[#4B5563] p-4 text-white shadow-2xl transition hover:scale-110">
         <Sparkles className="h-6 w-6" />
      </button>

      {stageModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
          <form
            onSubmit={handleCreateStage}
            className="w-full max-w-xl rounded-3xl bg-white p-6 shadow-2xl"
          >
            <div className="mb-4 flex items-center justify-between">
              <h3 className="text-lg font-bold">Новый этап</h3>
              <button
                type="button"
                onClick={() => setStageModalOpen(false)}
                className="text-sm text-gray-400 hover:text-black"
              >
                Закрыть
              </button>
            </div>
            <div className="space-y-3">
              <div className="space-y-2">
                <label className="text-sm font-semibold">Название этапа</label>
                <input
                  value={stageTitle}
                  onChange={(e) => setStageTitle(e.target.value)}
                  className="w-full rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm"
                  placeholder="Например: Предпроектная подготовка"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-semibold">Описание</label>
                <textarea
                  value={stageDescription}
                  onChange={(e) => setStageDescription(e.target.value)}
                  className="w-full rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm"
                  rows={3}
                  placeholder="Опишите цели этапа"
                />
              </div>
              {stageError && (
                <div className="rounded-xl bg-red-50 px-4 py-2 text-sm text-red-600">{stageError}</div>
              )}
              <div className="flex items-center gap-3 pt-2">
                <button
                  type="submit"
                  disabled={stageSubmitting}
                  className="rounded-full bg-black px-5 py-2 text-sm font-semibold text-white transition hover:bg-neutral-800 disabled:opacity-60"
                >
                  {stageSubmitting ? "Создание..." : "Создать этап"}
                </button>
                <button
                  type="button"
                  onClick={() => setStageModalOpen(false)}
                  className="text-sm text-gray-500 hover:text-black"
                >
                  Отмена
                </button>
              </div>
            </div>
          </form>
        </div>
      )}

    </div>
  )
}
