"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import Header from "@/components/header";
import { AlertTriangle, CalendarDays, CheckCircle2, ChevronLeft, ChevronRight, Clock, Flag, MapPin } from "lucide-react";
import { fetchJson } from "@/lib/api";

type CalendarEvent = {
  id: string;
  projectId: number;
  projectName: string;
  title: string;
  date: string; // ISO date (YYYY-MM-DD)
  status: "deadline" | "milestone" | "risk";
  location?: string;
};

type Project = {
  id: number;
  name: string;
  end_date?: string | null;
  start_date?: string | null;
};

const weekdays = ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"];

const statusStyles: Record<CalendarEvent["status"], string> = {
  deadline: "bg-rose-100 text-rose-700 ring-rose-200 dark:bg-rose-900/30 dark:text-rose-100 dark:ring-rose-800/80",
  milestone: "bg-emerald-100 text-emerald-700 ring-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-100 dark:ring-emerald-800/80",
  risk: "bg-amber-100 text-amber-800 ring-amber-200 dark:bg-amber-900/30 dark:text-amber-100 dark:ring-amber-800/80"
};

function startOfMonth(date: Date) {
  return new Date(date.getFullYear(), date.getMonth(), 1);
}

function buildMonthDays(base: Date) {
  const start = startOfMonth(base);
  const end = new Date(base.getFullYear(), base.getMonth() + 1, 0);
  const offset = (start.getDay() + 6) % 7; // shift Sunday=0 to Monday=0
  const days: Array<Date | null> = [];
  for (let i = 0; i < offset; i += 1) days.push(null);
  for (let d = 1; d <= end.getDate(); d += 1) {
    days.push(new Date(base.getFullYear(), base.getMonth(), d));
  }
  return days;
}

function dateKey(date: Date) {
  return date.toISOString().slice(0, 10);
}

function sameMonth(a: Date, b: Date) {
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth();
}

function formatHuman(dateStr: string) {
  return new Date(`${dateStr}T00:00:00`).toLocaleDateString("ru-RU", {
    day: "2-digit",
    month: "long",
    year: "numeric"
  });
}

export default function CalendarPage() {
  const router = useRouter();
  const [currentMonth, setCurrentMonth] = useState(() => startOfMonth(new Date()));
  const todayKey = dateKey(new Date());
  const [calendarEvents, setCalendarEvents] = useState<CalendarEvent[]>([]);

  const days = useMemo(() => buildMonthDays(currentMonth), [currentMonth]);

  useEffect(() => {
    async function loadProjectDeadlines() {
      try {
        const res = await fetchJson<any>(`/projects`);
        const items: Project[] = Array.isArray(res) ? res : res?.items || [];
        const mapped = items
          .filter((p) => p.end_date)
          .map((p) => ({
            id: `project-${p.id}`,
            projectId: p.id,
            projectName: p.name,
            title: "Дедлайн проекта",
            date: (p.end_date || "").slice(0, 10),
            status: "deadline" as const,
          }));
        setCalendarEvents(mapped);
      } catch (err) {
        console.error("Failed to load project deadlines", err);
        setCalendarEvents([]);
      }
    }
    loadProjectDeadlines();
  }, []);

  const eventsByDay = useMemo(() => {
    return calendarEvents.reduce<Record<string, CalendarEvent[]>>((acc, event) => {
      const key = event.date;
      acc[key] = acc[key] ? [...acc[key], event] : [event];
      return acc;
    }, {});
  }, [calendarEvents]);

  const eventsThisMonth = useMemo(
    () => calendarEvents.filter((evt) => sameMonth(new Date(`${evt.date}T00:00:00`), currentMonth)),
    [calendarEvents, currentMonth]
  );

  const upcoming = useMemo(() => {
    return [...calendarEvents].sort((a, b) => a.date.localeCompare(b.date));
  }, [calendarEvents]);

  const monthLabel = currentMonth.toLocaleDateString("ru-RU", { month: "long", year: "numeric" });

  const totalDeadlines = calendarEvents.filter((evt) => evt.status === "deadline").length;

  return (
    <div className="min-h-screen bg-white text-slate-900 dark:bg-slate-950 dark:text-slate-50">
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4 dark:border-slate-800 dark:bg-slate-950">
        <Header />
      </div>

      <main className="mx-auto flex max-w-6xl flex-col gap-8 px-6 py-8">
        <section className="flex flex-col gap-4 rounded-3xl border border-gray-200 bg-gradient-to-r from-amber-50 via-white to-white p-6 shadow-sm dark:border-slate-800 dark:from-slate-900 dark:via-slate-950 dark:to-slate-950">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <p className="text-sm uppercase tracking-[0.2em] text-amber-600 dark:text-amber-300">Полный календарь</p>
              <h1 className="mt-1 text-3xl font-bold leading-tight">Проекты по датам и дедлайнам</h1>
              <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">Контроль ключевых событий, рисков и сроков по всем активным объектам.</p>
            </div>
            <div className="flex items-center gap-3 rounded-2xl bg-white px-4 py-3 text-sm font-semibold text-slate-800 shadow-sm ring-1 ring-gray-200 dark:bg-slate-900 dark:text-slate-100 dark:ring-slate-800">
              <CalendarDays className="h-5 w-5 text-amber-500" />
              {monthLabel}
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div className="flex items-center justify-between rounded-2xl bg-white px-4 py-3 shadow-sm ring-1 ring-gray-200 dark:bg-slate-900 dark:ring-slate-800">
              <div>
                <p className="text-xs uppercase tracking-wide text-slate-500">События всего</p>
                <p className="text-2xl font-bold">{calendarEvents.length}</p>
              </div>
              <Flag className="h-6 w-6 text-amber-500" />
            </div>
            <div className="flex items-center justify-between rounded-2xl bg-white px-4 py-3 shadow-sm ring-1 ring-gray-200 dark:bg-slate-900 dark:ring-slate-800">
              <div>
                <p className="text-xs uppercase tracking-wide text-slate-500">В этом месяце</p>
                <p className="text-2xl font-bold">{eventsThisMonth.length}</p>
              </div>
              <Clock className="h-6 w-6 text-emerald-500" />
            </div>
            <div className="flex items-center justify-between rounded-2xl bg-white px-4 py-3 shadow-sm ring-1 ring-gray-200 dark:bg-slate-900 dark:ring-slate-800">
              <div>
                <p className="text-xs uppercase tracking-wide text-slate-500">Дедлайны</p>
                <p className="text-2xl font-bold">{totalDeadlines}</p>
              </div>
              <AlertTriangle className="h-6 w-6 text-rose-500" />
            </div>
            <div className="flex items-center justify-between rounded-2xl bg-white px-4 py-3 shadow-sm ring-1 ring-gray-200 dark:bg-slate-900 dark:ring-slate-800">
              <div>
                <p className="text-xs uppercase tracking-wide text-slate-500">События дня</p>
                <p className="text-2xl font-bold">{eventsByDay[todayKey]?.length ?? 0}</p>
              </div>
              <CheckCircle2 className="h-6 w-6 text-emerald-500" />
            </div>
          </div>
        </section>

        <section className="grid gap-6 lg:grid-cols-[2fr_1fr]">
          <div className="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900">
            <div className="mb-4 flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-500 dark:text-slate-400">Месяц</p>
                <h2 className="text-2xl font-bold capitalize">{monthLabel}</h2>
              </div>
              <div className="flex items-center gap-2">
                <button
                  className="flex h-10 w-10 items-center justify-center rounded-full border border-gray-200 bg-white text-slate-700 transition hover:bg-gray-50 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100"
                  onClick={() => setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1, 1))}
                  aria-label="Предыдущий месяц"
                >
                  <ChevronLeft className="h-5 w-5" />
                </button>
                <button
                  className="flex h-10 w-10 items-center justify-center rounded-full border border-gray-200 bg-white text-slate-700 transition hover:bg-gray-50 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100"
                  onClick={() => setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 1))}
                  aria-label="Следующий месяц"
                >
                  <ChevronRight className="h-5 w-5" />
                </button>
              </div>
            </div>

            <div className="grid grid-cols-7 gap-2 text-center text-xs font-semibold text-slate-500 uppercase tracking-wide">
              {weekdays.map((day) => (
                <div key={day} className="rounded-lg bg-slate-50 py-2 dark:bg-slate-800/70">
                  {day}
                </div>
              ))}
            </div>

            <div className="mt-3 grid grid-cols-7 gap-2">
              {days.map((day, idx) => {
                const key = day ? dateKey(day) : `empty-${idx}`;
                const eventsForDay = day ? eventsByDay[key] || [] : [];
                const isToday = day ? key === todayKey : false;

                return (
                  <div
                    key={key}
                    className={`min-h-[90px] rounded-2xl border p-3 text-sm transition hover:border-amber-400/80 hover:shadow-sm dark:hover:border-amber-400/60 ${
                      day
                        ? "border-gray-200 bg-white dark:border-slate-800 dark:bg-slate-900"
                        : "border-dashed border-gray-200 bg-slate-50 dark:border-slate-800 dark:bg-slate-900/40"
                    } ${isToday ? "ring-2 ring-amber-400/80" : ""}`}
                  >
                    <div className="flex items-center justify-between">
                      <span className={`text-base font-semibold ${day ? "text-slate-900 dark:text-slate-50" : "text-slate-300"}`}>
                        {day ? day.getDate() : ""}
                      </span>
                      {isToday && (
                        <span className="rounded-full bg-amber-100 px-2 py-0.5 text-[11px] font-semibold text-amber-700 dark:bg-amber-900/40 dark:text-amber-100">
                          Сегодня
                        </span>
                      )}
                    </div>
                    <div className="mt-2 space-y-1">
                      {eventsForDay.slice(0, 2).map((evt) => (
                        <div
                          key={evt.id}
                          onClick={() => router.push(`/project/${evt.projectId}`)}
                          className={`truncate rounded-xl px-2 py-1 text-[11px] font-semibold ring-1 transition hover:opacity-80 ${statusStyles[evt.status]}`}
                          role="button"
                          aria-label={`Открыть проект ${evt.projectName}`}
                          title={`${evt.title} · ${evt.projectName}`}
                        >
                          {evt.projectName}: {evt.title}
                        </div>
                      ))}
                      {eventsForDay.length > 2 && (
                        <div className="text-[11px] font-medium text-slate-500 dark:text-slate-400">+{eventsForDay.length - 2} ещё</div>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          <div className="flex flex-col gap-4 rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold">Ближайшие события</h3>
              <CalendarDays className="h-5 w-5 text-amber-500" />
            </div>

            <div className="space-y-3">
              {upcoming.map((evt) => (
                <div
                  key={evt.id}
                  onClick={() => router.push(`/project/${evt.projectId}`)}
                  className="rounded-2xl border border-gray-200 bg-white px-4 py-3 shadow-sm transition hover:border-amber-300 hover:shadow-md dark:border-slate-800 dark:bg-slate-950"
                  role="button"
                  aria-label={`Открыть проект ${evt.projectName}`}
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-semibold text-slate-900 dark:text-slate-50">{evt.projectName}</p>
                      <p className="text-xs text-slate-500 dark:text-slate-400">{evt.title}</p>
                    </div>
                    <span className={`rounded-full px-3 py-1 text-[11px] font-semibold ring-1 ${statusStyles[evt.status]}`}>
                      {evt.status === "deadline" ? "Дедлайн" : evt.status === "milestone" ? "Этап" : "Риск"}
                    </span>
                  </div>
                  <div className="mt-2 flex flex-wrap items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                    <span className="inline-flex items-center gap-1 rounded-full bg-slate-100 px-2 py-1 text-slate-700 dark:bg-slate-800 dark:text-slate-100">
                      <Clock className="h-3 w-3" /> {formatHuman(evt.date)}
                    </span>
                    {evt.location && (
                      <span className="inline-flex items-center gap-1 rounded-full bg-slate-100 px-2 py-1 text-slate-700 dark:bg-slate-800 dark:text-slate-100">
                        <MapPin className="h-3 w-3" /> {evt.location}
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
