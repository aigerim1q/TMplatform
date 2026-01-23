"use client"

import { usePathname, useRouter } from "next/navigation"
import { useState } from "react"
import { Maximize2, MessageCircle, Send, Sparkles, X } from "lucide-react"

export default function AIChatWidget() {
  const pathname = usePathname()
  const router = useRouter()
  const [open, setOpen] = useState(false)

  if (pathname.startsWith("/chat")) return null

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end gap-3">
      {!open && (
        <button
          onClick={() => setOpen(true)}
          className="flex h-14 w-14 items-center justify-center rounded-full bg-slate-800 text-white shadow-lg ring-2 ring-white/10 transition hover:-translate-y-1 hover:bg-slate-900"
          aria-label="Открыть AI чат"
        >
          <MessageCircle className="h-6 w-6" />
        </button>
      )}

      {open && (
        <div className="w-[360px] overflow-hidden rounded-3xl bg-slate-800 text-white shadow-2xl ring-1 ring-black/10">
          {/* Header */}
          <div className="flex items-center justify-between rounded-t-3xl bg-slate-700 px-4 py-3">
            <div className="flex items-center gap-2">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-amber-300 text-xs font-black text-slate-900 shadow-sm">
                THE
              </div>
              <span className="text-lg font-bold tracking-tight">QURYLYS AI</span>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => router.push("/chat")}
                className="rounded-full p-2 text-white/80 transition hover:bg-white/10 hover:text-white"
                aria-label="Открыть полноэкранно"
              >
                <Maximize2 className="h-5 w-5" />
              </button>
              <button
                onClick={() => setOpen(false)}
                className="rounded-full p-2 text-white/80 transition hover:bg-white/10 hover:text-white"
                aria-label="Закрыть"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
          </div>

          {/* Body */}
          <div className="space-y-3 bg-slate-800 px-4 pb-4 pt-3">
            <div className="rounded-2xl bg-slate-700/70 px-4 py-3 text-sm leading-relaxed text-white">
              Создай пожалуйста проект 5ти этажного ЖК на ровно-квадратной земле 50 соток, со всеми этапами, задачами и подзадачами, прикрепляя ответственных
            </div>

            <div className="inline-flex items-center gap-2 rounded-full bg-white text-slate-900 px-3 py-1 text-xs font-semibold shadow-sm">
              Время ответа: ~4 минуты
            </div>

            <div className="space-y-2 text-sm leading-relaxed text-white">
              <p>Добрый день</p>
              <p>
                Создал проект 5ти этажного дома на ровно-квадратной земле 5000 кв. м. с 4 этапами. Отметьте дату начала работ по этапу №1 и согласуйте создание задачи во вкладке «Уведомления».
              </p>
              <p>
                После этого проект со всеми задачами запустится, и всем ответственным по задачам первого этапа придёт уведомление. По окончанию работ первого этапа второй этап инициируется автоматически.
              </p>
            </div>
          </div>

          {/* Input */}
          <div className="flex items-center gap-2 border-t border-white/10 bg-slate-800 px-4 py-3">
            <button
              className="flex h-10 w-10 items-center justify-center rounded-full bg-amber-300 text-slate-900 shadow-sm transition hover:bg-amber-200"
              aria-label="Новый вопрос"
            >
              <Sparkles className="h-5 w-5" />
            </button>
            <div className="flex flex-1 items-center gap-2 rounded-full bg-white px-3 py-2 text-slate-900 shadow-inner">
              <input
                className="w-full border-none bg-transparent text-sm outline-none placeholder:text-slate-400"
                placeholder="Напишите свой вопрос"
              />
              <button
                className="rounded-full p-1 text-slate-500 transition hover:text-slate-800"
                aria-label="Отправить"
              >
                <Send className="h-5 w-5" />
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
