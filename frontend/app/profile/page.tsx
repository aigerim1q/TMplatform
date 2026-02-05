"use client"

import { useEffect, useMemo, useState } from "react"
import { useRouter } from "next/navigation"
import Image from "next/image"
import Header from "@/components/header"
import { Camera, Save, User, Loader2 } from "lucide-react"
import { logoutLocal } from "@/lib/auth"
import { fetchJson } from "@/lib/api"

export default function ProfilePage() {
  const router = useRouter()
  const [name, setName] = useState("Пользователь")
  const [title, setTitle] = useState("Роль не указана")
  const [email, setEmail] = useState("")
  const [avatar, setAvatar] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)

  const handleAvatar = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = () => {
        if (typeof reader.result === "string") {
          setAvatar(reader.result)
        }
      }
      reader.readAsDataURL(file)
    }
  }

  const initials = useMemo(() => {
    const source = (name || email || "").trim()
    if (!source) return "?"
    const parts = source.split(/\s+/).filter(Boolean)
    if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
    return source.slice(0, 2).toUpperCase()
  }, [name, email])

  useEffect(() => {
    let alive = true
    ;(async () => {
      try {
        setLoading(true)
        const me = await fetchJson<{ id: number; email: string; name?: string; role?: string; avatar_url?: string | null }>("/me")
        if (!alive) return
        setEmail(me.email || "")
        setName(me.name?.trim() || me.email?.split("@")[0] || "Пользователь")
        setTitle(me.role ? me.role.toUpperCase() : "Роль не указана")
        setAvatar(me.avatar_url || null)
        setError(null)
      } catch (err: any) {
        if (!alive) return
        setError(err?.message || "Не удалось загрузить профиль")
      } finally {
        if (alive) setLoading(false)
      }
    })()
    return () => {
      alive = false
    }
  }, [])

  const handleSave = async () => {
    try {
      setSaving(true)
      setSaved(false)
      await fetchJson("/me", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name,
          avatar_url: avatar,
        }),
      })
      window.dispatchEvent(
        new CustomEvent("me-updated", { detail: { avatar_url: avatar, name } })
      )
      setSaved(true)
    } catch (err: any) {
      setError(err?.message || "Не удалось сохранить профиль")
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="min-h-screen bg-white text-slate-900 dark:bg-slate-950 dark:text-slate-50">
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4 dark:border-slate-800 dark:bg-slate-950">
        <Header />
      </div>

      <main className="mx-auto flex max-w-5xl flex-col gap-8 px-6 py-10">
        <div className="space-y-2">
          <p className="text-sm font-semibold text-amber-500">Профиль</p>
          <h1 className="text-3xl font-bold">Настройки аккаунта</h1>
          <p className="text-sm text-slate-600 dark:text-slate-300">Обновите имя, должность, фото и контактные данные.</p>
        </div>

        <div className="grid gap-6 lg:grid-cols-[320px,1fr]">
          {/* Avatar card */}
          <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900">
            <div className="flex items-center gap-4">
              <div className="relative h-24 w-24 overflow-hidden rounded-full border border-slate-200 dark:border-slate-700 bg-slate-100 dark:bg-slate-800 flex items-center justify-center">
                {avatar ? (
                  <Image src={avatar} alt="Avatar" fill unoptimized sizes="96px" className="object-cover" />
                ) : (
                  <span className="text-xl font-semibold text-slate-700 dark:text-slate-100">{initials}</span>
                )}
              </div>
              <div>
                <p className="text-sm text-slate-500 dark:text-slate-300">Фото профиля</p>
                <p className="text-base font-semibold">{name}</p>
              </div>
            </div>
            <label className="mt-4 inline-flex cursor-pointer items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90">
              <Camera className="h-4 w-4" /> Обновить фото
              <input type="file" className="hidden" accept="image/*" onChange={handleAvatar} />
            </label>
          </div>

          {/* Form card */}
          <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900">
            <div className="mb-4 inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-200">
              <User className="h-4 w-4" /> Основная информация
            </div>

            {error && (
              <div className="mb-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-900/30 dark:text-red-100">
                {error}
              </div>
            )}

            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-xs font-semibold uppercase text-slate-500">Имя</label>
                <input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                  disabled={loading}
                />
              </div>
              <div className="space-y-2">
                <label className="text-xs font-semibold uppercase text-slate-500">Должность</label>
                <input
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                  disabled={loading}
                />
              </div>
              <div className="space-y-2 md:col-span-2">
                <label className="text-xs font-semibold uppercase text-slate-500">Email</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                  disabled={loading}
                />
              </div>
            </div>

            <div className="mt-6 flex flex-wrap gap-3">
              <button
                onClick={handleSave}
                disabled={loading || saving}
                className="inline-flex items-center gap-2 rounded-full bg-black px-5 py-3 text-sm font-semibold text-white transition hover:bg-black/80 disabled:opacity-60 dark:bg-white dark:text-black dark:hover:bg-white/90"
              >
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />} Сохранить изменения
              </button>
              <button
                onClick={() => router.back?.()}
                className="inline-flex items-center gap-2 rounded-full border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 dark:border-slate-700 dark:text-slate-100 dark:hover:bg-slate-800"
              >
                Отменить
              </button>
              <button
                onClick={() => {
                  logoutLocal()
                  router.push("/login")
                }}
                className="inline-flex items-center gap-2 rounded-full border border-red-200 px-5 py-3 text-sm font-semibold text-red-600 transition hover:bg-red-50 dark:border-red-900/50 dark:text-red-200 dark:hover:bg-red-900/30"
              >
                Выйти из аккаунта
              </button>
            </div>
            {saved && (
              <div className="mt-3 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700 dark:border-emerald-900/50 dark:bg-emerald-900/30 dark:text-emerald-100">
                Профиль сохранён
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
