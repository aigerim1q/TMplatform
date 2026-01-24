"use client";

import Image from "next/image";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { ArrowLeft, ArrowRight, Plus, Sparkles } from "lucide-react";
import Header from "@/components/header";
import { fetchJson } from "@/lib/api";

type Assignee = {
  id: number;
  name: string;
  email: string;
  role: string;
};

type Project = {
  id: number;
  name: string;
  description?: string | null;
  image_url?: string | null;
  start_date?: string | null;
  end_date?: string | null;
  budget_allocated: number;
  budget_spent: number;
  budget_currency: string;
  progress_percent: number;
  assignees: Assignee[];
};

type Stage = {
  id: number;
  title: string;
  description?: string | null;
  status: "todo" | "in_progress" | "done";
  end_date?: string | null;
};

export default function ProjectPreview() {
  const router = useRouter();
  const params = useParams();
  const projectId = Array.isArray(params.id) ? params.id[0] : params.id;

  const [project, setProject] = useState<Project | null>(null);
  const [stages, setStages] = useState<Stage[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const assigneesText = useMemo(() => {
    if (!project?.assignees?.length) return "Ответственные: не назначены";
    return `Ответственные: ${project.assignees.map((a) => a.name || a.email).join(", ")}`;
  }, [project?.assignees]);

  const budgetText = useMemo(() => {
    if (!project) return "—";
    return `${project.budget_spent.toLocaleString()} / ${project.budget_allocated.toLocaleString()} ${project.budget_currency || "₸"}`;
  }, [project]);

  const dueText = useMemo(() => {
    if (!project?.end_date) return "Без дедлайна";
    const end = new Date(project.end_date);
    if (Number.isNaN(end.getTime())) return "Без дедлайна";
    const diff = end.getTime() - Date.now();
    const days = Math.max(0, Math.ceil(diff / (1000 * 60 * 60 * 24)));
    return `${days} дней`;
  }, [project?.end_date]);

  useEffect(() => {
    let active = true;
    const load = async () => {
      if (!projectId) return;
      try {
        setLoading(true);
        const [projectData, stageData] = await Promise.all([
          fetchJson<Project>(`/projects/${projectId}`),
          fetchJson<{ items: Stage[] }>(`/projects/${projectId}/stages?sort=urgency&order=asc`),
        ]);
        if (!active) return;
        setProject(projectData);
        setStages(stageData.items || []);
        setError(null);
      } catch (err) {
        if (!active) return;
        setError("Не удалось загрузить проект");
      } finally {
        if (!active) return;
        setLoading(false);
      }
    };
    load();
    return () => {
      active = false;
    };
  }, [projectId]);

  const badge = (text: string, className: string) => (
    <span className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold ${className}`}>
      {text}
    </span>
  );

  const statusBadge = (stage: Stage) => {
    if (stage.status === "done") return badge("Выполнено", "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-100");
    if (stage.status === "in_progress")
      return badge("В работе", "bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-100");
    return badge("К выполнению", "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200");
  };

  return (
    <div className="min-h-screen bg-white text-slate-900 dark:bg-slate-950 dark:text-slate-50">
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4 dark:border-slate-800 dark:bg-slate-950">
        <Header />
      </div>

      <main className="mx-auto max-w-6xl px-6 py-8 space-y-8">
        {error && (
          <div className="rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700">
            {error}
          </div>
        )}
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
              src={project?.image_url || "https://images.unsplash.com/photo-1486325212027-8081e485255e?auto=format&fit=crop&w=900&q=80"}
              alt={project?.name || "Проект"}
              fill
              unoptimized
              sizes="(max-width: 640px) 100vw, 208px"
              className="object-cover"
            />
          </div>
          <div className="flex-1 space-y-3">
            <h1 className="text-2xl font-bold">Проект: {project?.name || "..."}</h1>
            <div className="flex flex-wrap items-center gap-3 text-sm">
              {badge("Финансы: " + budgetText, "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-100")}
              {badge("Отчеты расходов", "bg-black text-white dark:bg-white dark:text-black")}
            </div>
            <div className="flex flex-wrap items-center gap-3 text-sm text-slate-700 dark:text-slate-300">
              <div className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 dark:bg-slate-800">
                ⏱ Дедлайн: {project?.end_date ? new Date(project.end_date).toLocaleDateString("ru-RU") : "—"} ({dueText})
              </div>
              <div className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 dark:bg-slate-800">
                ▶ Дата начала: {project?.start_date ? new Date(project.start_date).toLocaleDateString("ru-RU") : "—"}
              </div>
              <div className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 dark:bg-slate-800">
                {assigneesText}
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
        <div className="space-y-4">
          <div className="rounded-3xl bg-white p-5 shadow-sm ring-1 ring-slate-200 dark:bg-slate-900 dark:ring-slate-800">
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
              <div className="flex items-center gap-3 text-sm font-semibold">
                <span className="inline-flex h-8 items-center rounded-full bg-slate-900 px-4 text-white dark:bg-white dark:text-slate-900">
                  Этапы проекта
                </span>
                {badge(`Прогресс: ${project?.progress_percent ?? 0}%`, "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-100")}
              </div>
              <div className="flex items-center gap-2">
                <button className="inline-flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90">
                  <Plus className="h-4 w-4" /> Добавить этап
                </button>
                <button className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-black text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90" aria-label="Навигация">
                  <ArrowRight className="h-5 w-5" />
                </button>
              </div>
            </div>

            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {loading && <div className="text-sm text-slate-500">Загрузка этапов...</div>}
              {!loading && stages.length === 0 && (
                <div className="text-sm text-slate-500">Этапы пока не созданы.</div>
              )}
              {stages.map((stage) => (
                <div
                  key={stage.id}
                  className="flex h-full flex-col justify-between rounded-2xl border border-slate-200 bg-white p-4 shadow-sm transition hover:-translate-y-0.5 dark:border-slate-800 dark:bg-slate-900"
                >
                  <div className="flex items-center justify-between text-xs text-slate-600 dark:text-slate-300">
                    <span>Проект: {project?.name || "—"}</span>
                    {statusBadge(stage)}
                  </div>
                  <div className="mt-3 space-y-2">
                    <h3 className="text-sm font-semibold text-slate-900 dark:text-white">{stage.title}</h3>
                    {stage.description && <p className="text-xs text-slate-600 dark:text-slate-300">{stage.description}</p>}
                  </div>
                  <div className="mt-3 text-xs text-slate-500">
                    {stage.end_date ? `Дедлайн: ${new Date(stage.end_date).toLocaleDateString("ru-RU")}` : "Без дедлайна"}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
