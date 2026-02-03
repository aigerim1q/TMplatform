"use client"

import Image from "next/image"
import { useRouter } from "next/navigation"
import { ArrowLeft, ArrowRight, Plus, Sparkles } from "lucide-react"
import Header from "@/components/header"

const urgentTasks = [
  {
    title: "Нужно привезти 10 плиток",
    description: "На объекте Shyraq срочно требуется новая плитка, 5x3 метра, 2 штуки для наличников вокруг л...",
    project: "Shyraq",
    due: "18 часов",
  },
]

const myTasks = [
  {
    title: "Возведение колонн на 13 этаже",
    description: "Интерьер подъездов и этажей Перед началом работ нужно провести подготовку...",
    project: "Shyraq",
    overdue: "-9 часов",
  },
  {
    title: "Нужно сделать новую планировку",
    description: "На объекте Ansau срочно требуется новая планировка 5-го этажа, чтобы совпадала с новым ...",
    project: "Ansau",
    due: "2 дня",
  },
  {
    title: "Нужно сделать новую планировку",
    description: "На объекте Dariya срочно требуется новая планировка 5-го этажа, чтобы совпадала с новым ...",
    project: "Dariya",
    due: "20 дней",
  },
]

const projects = [
  {
    name: "Shyraq",
    due: "25 дней",
    image:
      "https://images.unsplash.com/photo-1448630360428-65456885c650?auto=format&fit=crop&w=900&q=80",
  },
  {
    name: "Ansau",
    due: "55 дней",
    image:
      "https://images.unsplash.com/photo-1505693416388-ac5ce068fe85?auto=format&fit=crop&w=900&q=80",
  },
  {
    name: "Dariya",
    due: "55 дней",
    image:
      "https://images.unsplash.com/photo-1505691938895-1758d7feb511?auto=format&fit=crop&w=900&q=80",
  },
]

export default function Home() {
  const router = useRouter()

  return (
    <div className="min-h-screen bg-white text-gray-900 dark:bg-slate-950 dark:text-white">
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4 dark:border-slate-800 dark:bg-slate-950">
        <Header />
      </div>

      <main className="mx-auto flex max-w-6xl flex-col gap-10 px-6 py-8">
        {/* Срочные задачи */}
        <section className="space-y-3">
          <div className="flex items-center gap-3">
            <span className="h-4 w-4 rounded-full bg-emerald-400" aria-hidden />
            <div className="flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white">
              Срочные задачи: {urgentTasks.length}
            </div>
          </div>

          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {urgentTasks.map((task) => (
              <div key={task.title} className="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900">
                <div className="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-300">
                  <span>Проект: {task.project}</span>
                  <span className="inline-flex items-center gap-1 rounded-full bg-emerald-50 px-3 py-1 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-100">
                    <Sparkles className="h-3 w-3" /> {task.due}
                  </span>
                </div>
                <h3 className="mt-3 text-lg font-semibold leading-tight">{task.title}</h3>
                <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{task.description}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Мои задачи */}
        <section className="space-y-4">
          <div className="flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-3 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white">
              <span className="h-3 w-3 rounded-full bg-red-500" aria-hidden /> Мои задачи: {myTasks.length}
            </div>
            <button
              onClick={() => router.push("/project/new")}
              className="inline-flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80"
            >
              <Plus className="h-4 w-4" /> Добавить задачу
            </button>
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            {myTasks.map((task) => (
              <div key={task.title} className="flex h-full flex-col justify-between rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900">
                <div className="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                  <span>Проект: {task.project}</span>
                  <span
                    className={`inline-flex items-center gap-1 rounded-full px-3 py-1 text-xs font-semibold ${
                      task.overdue
                        ? "bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-100"
                        : "bg-emerald-50 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-100"
                    }`}
                  >
                    {task.overdue || task.due}
                  </span>
                </div>
                <div className="mt-3 space-y-2">
                  <h3 className="text-lg font-semibold leading-snug">{task.title}</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300">{task.description}</p>
                </div>
                <button
                  className="mt-4 inline-flex items-center justify-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80"
                  onClick={() => router.push(`/project/${task.project.toLowerCase()}/preview`)}
                >
                  Открыть
                </button>
              </div>
            ))}
          </div>
        </section>

        {/* Проекты */}
        <section className="space-y-4">
          <div className="flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-3 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white">
              <span className="h-3 w-3 rounded-full bg-yellow-400" aria-hidden /> Проекты: {projects.length}
            </div>
            <button
              onClick={() => router.push("/project/new")}
              className="inline-flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80"
            >
              <Plus className="h-4 w-4" /> Добавить проект
            </button>
          </div>

          <div className="grid gap-6 md:grid-cols-3">
            {projects.map((project) => (
              <div key={project.name} className="relative overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900">
                <div className="relative h-44 w-full">
                  <Image
                    src={project.image}
                    alt={project.name}
                    fill
                    unoptimized
                    sizes="(max-width: 768px) 100vw, 33vw"
                    className="object-cover"
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/50 via-black/10 to-transparent" />
                  <div className="absolute left-4 bottom-4 space-y-2 text-white">
                    <div className="text-sm font-semibold">Проект: {project.name}</div>
                    <div className="inline-flex items-center gap-1 rounded-full bg-white/80 px-3 py-1 text-xs font-semibold text-emerald-700">
                      ⏱ {project.due}
                    </div>
                  </div>
                  <button
                    className="absolute right-3 top-3 inline-flex h-10 w-10 items-center justify-center rounded-full bg-black/70 text-white transition hover:bg-black"
                    onClick={() => router.push(`/project/${project.name.toLowerCase()}/preview`)}
                    aria-label="Открыть проект"
                  >
                    <ArrowRight className="h-5 w-5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Задачи подчинённых */}
        <section className="space-y-4">
          <div className="flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-3 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white">
              <span className="h-3 w-3 rounded-full bg-red-500" aria-hidden /> Задачи подчинённых
            </div>
            <button
              onClick={() => router.push("/project/new")}
              className="inline-flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80"
            >
              <Plus className="h-4 w-4" /> Добавить задачу
            </button>
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            {myTasks.map((task) => (
              <div key={`${task.title}-sub`} className="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900">
                <div className="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                  <span>Проект: {task.project}</span>
                  <span
                    className={`inline-flex items-center gap-1 rounded-full px-3 py-1 text-xs font-semibold ${
                      task.overdue
                        ? "bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-100"
                        : "bg-emerald-50 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-100"
                    }`}
                  >
                    {task.overdue || task.due}
                  </span>
                </div>
                <h3 className="mt-3 text-lg font-semibold leading-tight">{task.title}</h3>
                <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{task.description}</p>
                <div className="mt-4 flex items-center justify-between">
                  <button
                    onClick={() => router.push(`/project/${task.project.toLowerCase()}/preview`)}
                    className="inline-flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80"
                  >
                    Открыть
                  </button>
                  <button
                    className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-gray-200 bg-white text-gray-900 transition hover:bg-gray-50 dark:border-slate-700 dark:bg-slate-900 dark:text-white"
                    aria-label="Следующая"
                  >
                    <ArrowRight className="h-5 w-5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Проекты подчинённых */}
        <section className="space-y-4 pb-10">
          <div className="flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-3 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white">
              <span className="h-3 w-3 rounded-full bg-yellow-400" aria-hidden /> Проекты подчинённых
            </div>
            <button
              onClick={() => router.push("/project/new")}
              className="inline-flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80"
            >
              <Plus className="h-4 w-4" /> Добавить проект
            </button>
          </div>

          <div className="grid gap-6 md:grid-cols-3">
            {projects.map((project) => (
              <div key={`${project.name}-sub`} className="relative overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900">
                <div className="relative h-44 w-full">
                  <Image
                    src={project.image}
                    alt={project.name}
                    fill
                    unoptimized
                    sizes="(max-width: 768px) 100vw, 33vw"
                    className="object-cover"
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/50 via-black/10 to-transparent" />
                  <div className="absolute left-4 bottom-4 space-y-2 text-white">
                    <div className="text-sm font-semibold">Проект: {project.name}</div>
                    <div className="inline-flex items-center gap-1 rounded-full bg-white/80 px-3 py-1 text-xs font-semibold text-emerald-700">
                      ⏱ {project.due}
                    </div>
                  </div>
                  <div className="absolute right-3 bottom-3 flex gap-2">
                    <button
                      className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-black/70 text-white transition hover:bg-black"
                      onClick={() => router.push(`/project/${project.name.toLowerCase()}/preview`)}
                      aria-label="Открыть проект"
                    >
                      <ArrowRight className="h-5 w-5" />
                    </button>
                    <button
                      className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-black/50 text-white transition hover:bg-black"
                      aria-label="Навигация"
                    >
                      <ArrowLeft className="h-5 w-5" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>
      </main>
    </div>
  )
}
