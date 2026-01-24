"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Sparkles, Upload, FileText, ArrowLeft } from "lucide-react"
import Header from "@/components/header"
import { fetchJson, getAccessToken, getApiBaseUrl } from "@/lib/api"

export default function NewProjectPage() {
  const router = useRouter()
  const [name, setName] = useState("Новый проект")
  const [description, setDescription] = useState("Опишите цели, сроки и основные параметры")
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [imagePreview, setImagePreview] = useState<string | null>(null)
  const [startDate, setStartDate] = useState("")
  const [endDate, setEndDate] = useState("")
  const [priority, setPriority] = useState(0)
  const [budgetAllocatedStr, setBudgetAllocatedStr] = useState("0")
  const [budgetSpentStr, setBudgetSpentStr] = useState("0")
  const [budgetCurrency, setBudgetCurrency] = useState("KZT")
  const [status, setStatus] = useState("draft")
  const [fileName, setFileName] = useState<string | null>(null)
  const [submitted, setSubmitted] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const formatNumber = (val: string | number) => {
    if (!val) return ""
    return val.toString().replace(/\D/g, "").replace(/\B(?=(\d{3})+(?!\d))/g, " ")
  }

  const parseNumber = (val: string) => {
    return parseInt(val.replace(/\s/g, "") || "0", 10)
  }

  const handleBudgetChange = (
    val: string, 
    setterStr: (v: string) => void
  ) => {
    const raw = val.replace(/\D/g, "")
    const formatted = raw.replace(/\B(?=(\d{3})+(?!\d))/g, " ")
    setterStr(formatted)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    const budgetAllocated = parseNumber(budgetAllocatedStr)
    const budgetSpent = parseNumber(budgetSpentStr)
    
    try {
      setSubmitted(true)
      let project: { id: number }

      if (imageFile) {
        const formData = new FormData()
        formData.append("name", name)
        formData.append("description", description)
        formData.append("status", status)
        formData.append("start_date", startDate ? new Date(startDate).toISOString() : "")
        formData.append("end_date", endDate ? new Date(endDate).toISOString() : "")
        formData.append("priority", String(priority))
        formData.append("budget_allocated", String(budgetAllocated))
        formData.append("budget_spent", String(budgetSpent))
        formData.append("budget_currency", budgetCurrency)
        formData.append("image", imageFile)

        const token = getAccessToken()
        const res = await fetch(`${getApiBaseUrl()}/projects`, {
          method: "POST",
          headers: token ? { Authorization: `Bearer ${token}` } : undefined,
          body: formData,
        })


        if (!res.ok) {
          const text = await res.text()
          throw new Error(text || "Не удалось создать проект")
        }

        project = await res.json()
      } else {
        const payload = {
          name,
          description,
          status,
          start_date: startDate ? new Date(startDate).toISOString() : "",
          end_date: endDate ? new Date(endDate).toISOString() : "",
          priority,
          budget: {
            allocated: budgetAllocated,
            spent: budgetSpent,
            currency: budgetCurrency,
          },
        }

        project = await fetchJson<{ id: number }>("/projects", {
          method: "POST",
          body: JSON.stringify(payload),
        })
      }

      router.push(`/project/${project.id}`)
    } catch (err) {
      const message = err instanceof Error ? err.message : "Не удалось создать проект"
      setError(message)
      setSubmitted(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 text-slate-900 dark:bg-slate-950 dark:text-slate-50">
      <div className="flex w-screen justify-center border-b border-gray-200 bg-gray-50 py-4 dark:border-slate-800 dark:bg-slate-950">
        <Header />
      </div>

      <main className="mx-auto max-w-4xl px-6 py-10 space-y-8">
        <button
          onClick={() => router.back()}
          className="inline-flex items-center gap-2 text-sm text-slate-600 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
        >
          <ArrowLeft className="h-4 w-4" /> Назад
        </button>

        <div className="flex items-center gap-2 rounded-full bg-purple-100 px-4 py-2 text-sm font-semibold text-purple-700 dark:bg-purple-900/40 dark:text-purple-100">
          <Sparkles className="h-4 w-4" /> AI поможет собрать ЖЦП и бюджет
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900">
          <div className="mb-6">
            <p className="text-sm font-semibold text-green-600 dark:text-green-300">Создание проекта</p>
            <h1 className="text-2xl font-bold">Новый проект</h1>
            <p className="text-sm text-slate-600 dark:text-slate-300">Заполните поля или загрузите ТЗ — AI предложит структуру и сроки.</p>
          </div>

          <form className="space-y-6" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <label className="text-sm font-semibold">Название проекта</label>
              <input
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-lg font-semibold tracking-tight focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-semibold">Описание</label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={4}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-base focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
              />
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-sm font-semibold">Изображение (файл)</label>
                <input
                  type="file"
                  accept="image/*"
                  onChange={(e) => {
                    const file = e.target.files?.[0] || null
                    setImageFile(file)
                    if (file) {
                      const reader = new FileReader()
                      reader.onload = () => setImagePreview(reader.result as string)
                      reader.readAsDataURL(file)
                    } else {
                      setImagePreview(null)
                    }
                  }}
                  className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
                {imagePreview && (
                  <img
                    src={imagePreview}
                    alt="Preview"
                    className="mt-3 h-28 w-full rounded-xl object-cover"
                  />
                )}
              </div>
              <div className="space-y-2">
                <label className="text-sm font-semibold">Статус проекта</label>
                <select
                  value={status}
                  onChange={(e) => setStatus(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-base focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                >
                  <option value="draft">Черновик</option>
                  <option value="active">Активный</option>
                  <option value="paused">Пауза</option>
                  <option value="done">Завершён</option>
                </select>
              </div>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-sm font-semibold">Дата начала</label>
                <input
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-base focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-semibold">Дата завершения</label>
                <input
                  type="date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-base focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div className="space-y-2">
                <label className="text-sm font-semibold">Приоритет (0-10)</label>
                <input
                  type="number"
                  min="0"
                  max="10"
                  value={priority}
                  onChange={(e) => setPriority(Number(e.target.value))}
                  className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-lg font-bold focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-semibold">Бюджет (выделено)</label>
                <input
                  type="text"
                  value={budgetAllocatedStr}
                  onChange={(e) => handleBudgetChange(e.target.value, setBudgetAllocatedStr)}
                  placeholder="0"
                  className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-lg font-bold tracking-wide focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-semibold">Бюджет (потрачено)</label>
                <input
                  type="text"
                  value={budgetSpentStr}
                  onChange={(e) => handleBudgetChange(e.target.value, setBudgetSpentStr)}
                  placeholder="0"
                  className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-lg font-bold tracking-wide focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-sm font-semibold">Валюта</label>
                <input
                  value={budgetCurrency}
                  onChange={(e) => setBudgetCurrency(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-lg font-semibold focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
            </div>

            <div className="space-y-3 rounded-xl border border-dashed border-amber-300 bg-amber-50 p-4 text-sm text-amber-900 dark:border-amber-500/60 dark:bg-amber-500/10 dark:text-amber-100">
              <div className="flex items-center gap-2 font-semibold">
                <Upload className="h-4 w-4" /> Загрузить ТЗ или спецификацию (PDF/DOCX)
              </div>
              <label className="inline-flex cursor-pointer items-center gap-2 rounded-full bg-white px-3 py-2 text-xs font-semibold text-amber-900 shadow-sm transition hover:bg-amber-100 dark:bg-amber-500/20 dark:text-amber-50 dark:hover:bg-amber-500/30">
                <FileText className="h-4 w-4" /> Выбрать файл
                <input
                  type="file"
                  className="hidden"
                  onChange={(e) => setFileName(e.target.files?.[0]?.name ?? null)}
                  accept=".pdf,.doc,.docx"
                />
              </label>
              {fileName && <p className="text-xs">Вы выбрали: {fileName}</p>}
            </div>

            <div className="flex flex-wrap gap-3">
              <button
                type="submit"
                className="rounded-full bg-black px-5 py-3 text-sm font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90"
              >
                Создать проект
              </button>
            </div>

            {submitted && (
              <div className="rounded-xl border border-green-200 bg-green-50 px-4 py-3 text-sm font-semibold text-green-700 dark:border-green-700 dark:bg-green-900/30 dark:text-green-100">
                Проект создан. Перенаправляем...
              </div>
            )}

            {error && (
              <div className="rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">
                {error}
              </div>
            )}
          </form>
        </div>
      </main>
    </div>
  )
}
