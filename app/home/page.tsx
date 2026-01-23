import Image from "next/image"
import { Button } from "@/components/ui/button"

export default function DashboardPage() {
  return (
    <div className="min-h-screen bg-[#F7F6F2] px-6 py-6">
      <Header />

      <main className="mx-auto mt-8 max-w-[1200px] space-y-10">
        {/* Срочные задачи */}
        <SectionBadge color="green" text="Срочные задачи: 1" />
        <div className="flex gap-4">
          <TaskCard
            project="Shyraq"
            time="18 часов"
            title="Нужно привести 10 плиток"
            desc="На объекте Shyraq срочно требуется новая плитка. 5×3 метра."
            status="green"
          />
        </div>

        {/* Мои задачи */}
        <SectionHeader
          left="Мои задачи: 10"
          right="+ Добавить задачу"
          color="red"
        />
        <div className="flex gap-4 overflow-x-auto">
          <TaskCard
            project="Shyraq"
            time="-9 часов"
            title="Возведение колонн на 1 этаже"
            desc="Перед началом работ нужно провести подготовку."
            status="red"
          />
          <TaskCard
            project="Ansau"
            time="2 дня"
            title="Нужно сделать новую планировку"
            desc="На объекте Ansau срочно требуется новая планировка."
            status="green"
          />
          <TaskCard
            project="Dariya"
            time="20 дней"
            title="Нужно сделать новую планировку"
            desc="На объекте Dariya срочно требуется новая планировка."
            status="green"
          />
        </div>

        {/* Проекты */}
        <SectionHeader
          left="Проекты: 3"
          right="+ Добавить проект"
          color="yellow"
        />
        <div className="flex gap-4 overflow-x-auto">
          <ProjectCard title="Shyraq" days="25 дней" />
          <ProjectCard title="Ansau" days="55 дней" />
          <ProjectCard title="Dariya" days="55 дней" />
        </div>
      </main>
    </div>
  )
}

/* ================= HEADER ================= */

function Header() {
  return (
    <header className="flex justify-center">
      <div className="flex w-full max-w-[1200px] items-center justify-between rounded-full bg-white px-8 py-4 shadow-sm">
        {/* Logo */}
        <div className="text-xs font-semibold tracking-widest text-gray-700">
          THE QURYLYS
        </div>

        {/* Nav */}
        <nav className="flex items-center gap-8 text-sm">
          <span className="text-gray-400">Календарь</span>
          <span className="text-gray-400">Иерархия</span>
          <span className="font-medium text-black underline">
            Дашборд
          </span>
          <span className="text-gray-400">ЖЦП</span>
        </nav>

        {/* Right */}
        <div className="flex items-center gap-3">
          <div className="h-6 w-6 rounded-full bg-gray-200" />
          <div className="h-8 w-8 rounded-full bg-gray-300" />
        </div>
      </div>
    </header>
  )
}

/* ================= UI BLOCKS ================= */

function SectionBadge({
  text,
  color,
}: {
  text: string
  color: "green" | "red" | "yellow"
}) {
  const map = {
    green: "bg-green-100 text-green-700",
    red: "bg-red-100 text-red-700",
    yellow: "bg-yellow-100 text-yellow-700",
  }

  return (
    <span className={`inline-block rounded-full px-4 py-1 text-sm ${map[color]}`}>
      {text}
    </span>
  )
}

function SectionHeader({
  left,
  right,
  color,
}: {
  left: string
  right: string
  color: "green" | "red" | "yellow"
}) {
  return (
    <div className="flex items-center justify-between">
      <SectionBadge text={left} color={color} />
      <Button className="rounded-full bg-[#E6D8AE] text-black hover:bg-[#dccb95]">
        {right}
      </Button>
    </div>
  )
}

/* ================= CARDS ================= */

function TaskCard({
  project,
  time,
  title,
  desc,
  status,
}: any) {
  const statusColor =
    status === "red"
      ? "bg-red-500"
      : status === "green"
      ? "bg-green-500"
      : "bg-yellow-400"

  return (
    <div className="w-[300px] shrink-0 rounded-xl bg-white p-4 shadow-sm">
      <div className="mb-2 flex items-center justify-between text-xs text-gray-500">
        <span>Проект: {project}</span>
        <span className={`rounded-full px-2 py-1 text-white ${statusColor}`}>
          {time}
        </span>
      </div>

      <h3 className="mb-1 text-sm font-semibold">{title}</h3>
      <p className="text-xs text-gray-500">{desc}</p>
    </div>
  )
}

function ProjectCard({ title, days }: any) {
  return (
    <div className="relative h-[160px] w-[260px] shrink-0 overflow-hidden rounded-xl bg-gray-300">
      <div className="absolute left-3 top-3 rounded-full bg-black px-3 py-1 text-xs text-white">
        Проект: {title}
      </div>
      <div className="absolute right-3 top-3 rounded-full bg-green-500 px-2 py-1 text-xs text-white">
        {days}
      </div>
    </div>
  )
}
