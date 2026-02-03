"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { ArrowRight, Plus, Sparkles } from "lucide-react";
import Header from "@/components/header";

const urgentTasks = [
  {
    title: "Нужно привезти 10 плиток",
    description:
      "На объекте Shyraq срочно требуется новая плитка, 5x3 метра, 2 штуки для наличников вокруг дверей на 17 этаже",
    project: "Shyraq",
    due: "18 часов",
  },
];

const myTasks = [
  {
    title: "Возведение колонн на 13 этаже",
    description: "Интерьеры подъездов и этажей. Перед началом работ нужно провести подготовку...",
    project: "Shyraq",
    overdue: "-9 часов",
  },
  {
    title: "Нужно сделать новую планировку",
    description: "На объекте Ansau срочно требуется новая планировка 5-го этажа, чтобы совпадала с новым...",
    project: "Ansau",
    due: "2 дня",
  },
  {
    title: "Нужно сделать новую планировку",
    description: "На объекте Dariya срочно требуется новая планировка 5-го этажа, чтобы совпадала с новым...",
    project: "Dariya",
    due: "20 дней",
  },
];

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
];

const badge = (text: string, color: "red" | "green" | "yellow") => {
  const map = {
    red: "bg-rose-100 text-rose-700 dark:bg-rose-500/20 dark:text-rose-100",
    green: "bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-100",
    yellow: "bg-amber-100 text-amber-700 dark:bg-amber-500/20 dark:text-amber-100",
  } as const;
  return <span className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold ${map[color]}`}>{text}</span>;
};

export default function DashboardScreen() {
  const router = useRouter();

  const sectionCard = (title: string, children: React.ReactNode, accent: "red" | "green" | "yellow") => {
    const accentColors = {
      red: "bg-rose-500",
      green: "bg-emerald-400",
      yellow: "bg-amber-400",
    } as const;
    return (
      <section className="space-y-3">
        <div className="flex items-center gap-3">
          <span className={`h-4 w-4 rounded-full ${accentColors[accent]}`} aria-hidden />
          <div className="flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white dark:bg-white dark:text-black">
            {title}
          </div>
        </div>
        {children}
      </section>
    );
  };

  const taskCard = (task: typeof myTasks[number]) => (
    <div className="flex h-full flex-col justify-between rounded-2xl border border-white/10 bg-white/70 p-4 shadow-sm backdrop-blur dark:border-[#2d2648] dark:bg-[#161126]">
      <div className="flex items-center gap-2 text-xs text-slate-700 dark:text-slate-200">
        <span>Проект: {task.project}</span>
        {task.overdue
          ? badge(task.overdue, "red")
          : badge(task.due ?? "", "green")}
      </div>
      <div className="mt-3 space-y-2">
        <h3 className="text-lg font-semibold leading-snug text-slate-900 dark:text-white">{task.title}</h3>
        <p className="text-sm text-slate-700 dark:text-slate-200">{task.description}</p>
      </div>
      <button
        className="mt-4 inline-flex items-center justify-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90"
        onClick={() => router.push(`/project/${task.project.toLowerCase()}/preview`)}
      >
        Открыть
      </button>
    </div>
  );

  const projectCard = (project: typeof projects[number]) => (
    <div className="relative overflow-hidden rounded-2xl border border-white/10 bg-white/70 shadow-sm backdrop-blur dark:border-[#2d2648] dark:bg-[#161126]">
      <div className="relative h-44 w-full">
        <Image
          src={project.image}
          alt={project.name}
          fill
          unoptimized
          sizes="(max-width: 768px) 100vw, 33vw"
          className="object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent" />
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
  );

  return (
    <div className="min-h-screen bg-white text-slate-900 transition-colors dark:bg-[#0f0a1f] dark:text-white">
      <div className="flex w-screen justify-center border-b border-transparent bg-transparent px-4 py-4">
        <Header />
      </div>

      <main className="mx-auto flex max-w-6xl flex-col gap-8 px-4 pb-14">
        <div className="rounded-3xl border border-white/10 bg-gradient-to-br from-[#120d24] via-[#0f0a1f] to-[#090716] p-6 shadow-lg dark:border-[#241c3d] dark:from-[#120d24] dark:via-[#0f0a1f] dark:to-[#090716]">
          {sectionCard(`Срочные задачи: ${urgentTasks.length}`, (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {urgentTasks.map((task) => (
                <div key={task.title} className="rounded-2xl border border-emerald-200/70 bg-white/80 p-4 shadow-sm backdrop-blur dark:border-emerald-500/20 dark:bg-[#161126]">
                  <div className="flex items-center gap-2 text-xs text-slate-700 dark:text-slate-200">
                    <span>Проект: {task.project}</span>
                    {badge(task.due, "green")}
                  </div>
                  <h3 className="mt-3 text-lg font-semibold leading-tight text-slate-900 dark:text-white">{task.title}</h3>
                  <p className="mt-2 text-sm text-slate-700 dark:text-slate-200">{task.description}</p>
                </div>
              ))}
            </div>
          ), "green")}

          <div className="mt-6 flex flex-col gap-8">
            {sectionCard(`Мои задачи: ${myTasks.length}`, (
              <div className="grid gap-4 md:grid-cols-3">
                {myTasks.map((task) => (
                  <div key={task.title} className="relative">
                    {taskCard(task)}
                    <button
                      className="absolute right-[-16px] top-1/2 hidden h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full bg-black text-white shadow-lg transition hover:bg-black/80 lg:flex dark:bg-white dark:text-black"
                      aria-label="Навигация"
                    >
                      <ArrowRight className="h-5 w-5" />
                    </button>
                  </div>
                ))}
              </div>
            ), "red")}

            {sectionCard(`Проекты: ${projects.length}`, (
              <div className="grid gap-6 md:grid-cols-3">
                {projects.map((project) => (
                  <div key={project.name} className="relative">
                    {projectCard(project)}
                    <button
                      className="absolute right-[-16px] top-1/2 hidden h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full bg-black text-white shadow-lg transition hover:bg-black/80 lg:flex dark:bg-white dark:text-black"
                      aria-label="Навигация"
                    >
                      <ArrowRight className="h-5 w-5" />
                    </button>
                  </div>
                ))}
              </div>
            ), "yellow")}

            {sectionCard(`Задачи подчинённых`, (
              <div className="grid gap-4 md:grid-cols-3">
                {myTasks.map((task, idx) => (
                  <div key={`${task.title}-sub-${idx}`} className="relative">
                    {taskCard(task)}
                    <button
                      className="absolute right-[-16px] top-1/2 hidden h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full bg-black text-white shadow-lg transition hover:bg-black/80 lg:flex dark:bg-white dark:text-black"
                      aria-label="Навигация"
                    >
                      <ArrowRight className="h-5 w-5" />
                    </button>
                  </div>
                ))}
              </div>
            ), "red")}

            {sectionCard(`Проекты подчинённых: ${projects.length}`, (
              <div className="grid gap-6 md:grid-cols-3">
                {projects.map((project, idx) => (
                  <div key={`${project.name}-sub-${idx}`} className="relative">
                    {projectCard(project)}
                  </div>
                ))}
              </div>
            ), "yellow")}
          </div>

          <div className="mt-8 flex flex-wrap items-center gap-3">
            <button
              onClick={() => router.push("/project/new")}
              className="inline-flex items-center gap-2 rounded-full bg-black px-4 py-2 text-sm font-semibold text-white transition hover:bg-black/80 dark:bg-white dark:text-black dark:hover:bg-white/90"
            >
              <Plus className="h-4 w-4" /> Добавить задачу
            </button>
            <button
              onClick={() => router.push("/chat")}
              className="inline-flex items-center gap-2 rounded-full bg-amber-300 px-4 py-2 text-sm font-semibold text-slate-900 transition hover:bg-amber-200"
            >
              <Sparkles className="h-4 w-4" /> AI: предложить план
            </button>
          </div>
        </div>
      </main>
    </div>
  );
}
