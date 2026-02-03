"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Image from "next/image"
import Header from "@/components/header"
import { Camera, Save, User } from "lucide-react"

export default function ProfilePage() {
  const router = useRouter()
  const [name, setName] = useState("Омар Ахмет")
  const [title, setTitle] = useState("Руководитель проектов")
  const [email, setEmail] = useState("omar@quryls.kz")
  const [avatar, setAvatar] = useState(
    "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
  )

  const handleAvatar = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const url = URL.createObjectURL(file)
      setAvatar(url)
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
              <div className="relative h-24 w-24 overflow-hidden rounded-full border border-slate-200 dark:border-slate-700">
                <Image src={avatar} alt="Avatar" fill unoptimized sizes="96px" className="object-cover" />
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

            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-xs font-semibold uppercase text-slate-500">Имя</label>
                <input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
              <div className="space-y-2">
                <label className="text-xs font-semibold uppercase text-slate-500">Должность</label>
                <input
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
              <div className="space-y-2 md:col-span-2">
                <label className="text-xs font-semibold uppercase text-slate-500">Email</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm focus:border-amber-400 focus:outline-none dark:border-slate-700 dark:bg-slate-800"
                />
              </div>
            </div>

            <div className="mt-6 flex flex-wrap gap-3">
              <button className="inline-flex items-center gap-2 rounded-full bg-black px-5 py-3 text-sm font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90">
                <Save className="h-4 w-4" /> Сохранить изменения
              </button>
              <button
                onClick={() => router.back?.()}
                className="inline-flex items-center gap-2 rounded-full border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 dark:border-slate-700 dark:text-slate-100 dark:hover:bg-slate-800"
              >
                Отменить
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
