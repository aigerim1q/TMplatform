"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Sparkles, Upload, FileText, ArrowLeft } from "lucide-react"
import Header from "@/components/header"

export default function NewProjectPage() {
  const router = useRouter()
  const [name, setName] = useState("Новый проект")
  const [description, setDescription] = useState("Опишите цели, сроки и основные параметры")
  const [fileName, setFileName] = useState<string | null>(null)
  const [submitted, setSubmitted] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitted(true)
    setTimeout(() => {
      router.push("/lifecycle")
    }, 800)
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
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-semibold">Описание</label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={4}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
              />
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
              <button
                type="button"
                onClick={() => router.push("/lifecycle")}
                className="rounded-full border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800"
              >
                Выбрать шаблон ЖЦП
              </button>
            </div>

            {submitted && (
              <div className="rounded-xl border border-green-200 bg-green-50 px-4 py-3 text-sm font-semibold text-green-700 dark:border-green-700 dark:bg-green-900/30 dark:text-green-100">
                Проект отправлен в обработку AI. Перенаправляем на настройку ЖЦП...
              </div>
            )}
          </form>
        </div>
      </main>
    </div>
  )
}
