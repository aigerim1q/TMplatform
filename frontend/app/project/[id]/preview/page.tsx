"use client";

import Image from "next/image";
import { useParams, useRouter } from "next/navigation";
import { ArrowLeft, ArrowRight, Plus, Sparkles } from "lucide-react";
import Header from "@/components/header";

type ProjectKey = "shyraq" | "ansau" | "dariya";

type Stage = {
  title: string;
  subtitle?: string;
  status: "done" | "progress" | "delay";
  delay?: string;
};

type Section = {
  title: string;
  statusLabel: string;
  statusColor: string;
  tasks: Stage[];
};

export default function ProjectPreview() {
  const router = useRouter();
  const params = useParams();

  const paramId = Array.isArray(params.id) ? params.id[0] : params.id;
  const projectKey: ProjectKey =
    paramId === "ansau" ? "ansau" : paramId === "dariya" ? "dariya" : "shyraq";

  const projectMeta: Record<ProjectKey, { name: string; image: string; finance: string }> = {
    shyraq: {
      name: "Shyraq",
      image:
        "https://images.unsplash.com/photo-1486325212027-8081e485255e?auto=format&fit=crop&w=900&q=80",
      finance: "1,500,900,000/2,400,800,000",
    },
    ansau: {
      name: "Ansau",
      image:
        "https://images.unsplash.com/photo-1505691938895-1758d7feb511?auto=format&fit=crop&w=900&q=80",
      finance: "980,200,000/1,800,000,000",
    },
    dariya: {
      name: "Dariya",
      image:
        "https://images.unsplash.com/photo-1505691938895-1758d7feb511?auto=format&fit=crop&w=900&q=80",
      finance: "640,000,000/1,200,000,000",
    },
  };

  const sections: Section[] = [
    {
      title: "I. Предпроектная подготовка (до начала проектирования)",
      statusLabel: "Задержка: 5 дней",
      statusColor: "bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-100",
      tasks: [
        {
          title: "1. Инициирование проекта",
          subtitle: "Определение ориентированной площади и этажности, Анализ потребностей рынка...",
          status: "done",
        },
        {
          title: "2. Финансово-экономический анализ",
          subtitle: "Прогноз стоимости строительства, Расчёт рентабельности проекта (ROI)...",
          status: "delay",
          delay: "Задержка: 5 дней",
        },
      ],
    },
    {
      title: "II. Проектирование (архитектура, инженерия, дизайн)",
      statusLabel: "Выполнено",
      statusColor: "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-100",
      tasks: [
        {
          title: "1. Архитектурная концепция",
          subtitle: "Концептуальный проект здания, Общее планировочное решение...",
          status: "done",
        },
        {
          title: "2. Инженерные разделы",
          subtitle: "Конструктив (фундамент, колонны, плиты), Электроснабжение...",
          status: "done",
        },
        {
          title: "3. Дизайн интерьера и фасадов",
          subtitle: "Интерьеры подъездов и этажей, Интерьеры квартир...",
          status: "done",
        },
      ],
    },
    {
      title: "III. Строительный этап",
      statusLabel: "25 дней",
      statusColor: "bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-100",
      tasks: [
        {
          title: "1. Подготовка площадки",
          subtitle: "Ограживание участка, Установка бытовок и складов...",
          status: "done",
        },
        {
          title: "2. Фундаментные работы",
          subtitle: "Геодезическая разбивка, Рытьё котлована, Подготовка основания...",
          status: "done",
        },
        {
          title: "3. Возведение колонн на 13 этаже",
          subtitle: "Интерьеры подъездов и этажей. Перед началом работ нужно провести подготовку...",
          status: "delay",
          delay: "-9 часов",
        },
      ],
    },
  ];

  const project = projectMeta[projectKey];

  const badge = (text: string, className: string) => (
    <span className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold ${className}`}>
      {text}
    </span>
  );

  const statusBadge = (task: Stage) => {
    if (task.status === "done") return badge("Выполнено", "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-100");
    if (task.status === "delay")
      return badge(task.delay || "Задержка", "bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-100");
    return badge("В работе", "bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-100");
  };

  return (
    <div className="min-h-screen bg-white text-slate-900 dark:bg-slate-950 dark:text-slate-50">
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4 dark:border-slate-800 dark:bg-slate-950">
        <Header />
      </div>

      <main className="mx-auto max-w-6xl px-6 py-8 space-y-8">
        {/* Top bar */}
        <div className="flex flex-wrap items-center gap-3">
          <button
            onClick={() => router.back()}
            className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-4 py-2 text-sm font-semibold text-slate-900 transition hover:bg-slate-200 dark:bg-slate-800 dark:text-white dark:hover:bg-slate-700"
          >
            <ArrowLeft className="h-4 w-4" /> Назад
          </button>
          <div className="flex flex-wrap gap-2">
            {badge("Задача", "bg-black text-white dark:bg-white dark:text-black")}
            {badge("Отчеты", "bg-slate-800 text-white dark:bg-slate-200 dark:text-slate-900")}
          </div>
        </div>

        {/* Hero */}
        <div className="flex flex-col gap-4 rounded-3xl bg-white p-4 shadow-sm ring-1 ring-slate-200 dark:bg-slate-900 dark:ring-slate-800 sm:flex-row">
          <div className="relative h-36 w-full overflow-hidden rounded-2xl sm:w-52">
            <Image
              src={project.image}
              alt={project.name}
              fill
              unoptimized
              sizes="(max-width: 640px) 100vw, 208px"
              className="object-cover"
            />
          </div>
          <div className="flex-1 space-y-3">
            <h1 className="text-2xl font-bold">Проект: {project.name}</h1>
            <div className="flex flex-wrap items-center gap-3 text-sm">
              {badge("Финансы: " + project.finance, "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-100")}
              {badge("Отчеты расходов", "bg-black text-white dark:bg-white dark:text-black")}
            </div>
            <div className="flex flex-wrap items-center gap-3 text-sm text-slate-700 dark:text-slate-300">
              <div className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 dark:bg-slate-800">
                ⏱ Дедлайн: 12.12.2025 23:59 (25 дней)
              </div>
              <div className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 dark:bg-slate-800">
                ▶ Дата начала: 9.06.2025 12:00
              </div>
              <div className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 dark:bg-slate-800">
                Ответственные: Омар Ахмет, Зейнула Рышым, Серик Рахым...
              </div>
              <button
                onClick={() => router.push("/chat")}
                className="inline-flex items-center gap-2 rounded-full bg-black px-3 py-1 text-xs font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90"
              >
                <Sparkles className="h-4 w-4" /> AI: Причины отсрочки →
              </button>
            </div>
          </div>
        </div>

        {/* Sections */}
        <div className="space-y-8">
          {sections.map((section, idx) => (
            <div key={section.title} className="rounded-3xl bg-white p-5 shadow-sm ring-1 ring-slate-200 dark:bg-slate-900 dark:ring-slate-800">
              <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
                <div className="flex items-center gap-3 text-sm font-semibold">
                  <span className="inline-flex h-8 items-center rounded-full bg-slate-900 px-4 text-white dark:bg-white dark:text-slate-900">
                    {section.title}
                  </span>
                  {badge(section.statusLabel, section.statusColor)}
                </div>
                <div className="flex items-center gap-2">
                  <button className="inline-flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90">
                    <Plus className="h-4 w-4" /> Добавить задачу
                  </button>
                  <button className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-black text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90" aria-label="Навигация">
                    <ArrowRight className="h-5 w-5" />
                  </button>
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {section.tasks.map((task, tIdx) => (
                  <div
                    key={task.title + tIdx}
                    className="flex h-full flex-col justify-between rounded-2xl border border-slate-200 bg-white p-4 shadow-sm transition hover:-translate-y-0.5 dark:border-slate-800 dark:bg-slate-900"
                  >
                    <div className="flex items-center justify-between text-xs text-slate-600 dark:text-slate-300">
                      <span>Проект: {project.name}</span>
                      {statusBadge(task)}
                    </div>
                    <div className="mt-3 space-y-2">
                      <h3 className="text-sm font-semibold text-slate-900 dark:text-white">{task.title}</h3>
                      {task.subtitle && <p className="text-xs text-slate-600 dark:text-slate-300">{task.subtitle}</p>}
                    </div>
                  </div>
                ))}
              </div>

              {idx === sections.length - 1 && (
                <div className="mt-5 flex justify-center">
                  <button className="inline-flex items-center gap-2 rounded-full bg-black px-5 py-3 text-sm font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90">
                    <Plus className="h-4 w-4" /> Добавить этап проекта
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
