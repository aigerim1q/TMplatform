"use client"

import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { 
  Calendar, 
  Users, 
  LayoutDashboard, 
  FileText, 
  Bell, 
  User,
  Clock,
  AlertCircle,
  CheckCircle2,
  Plus,
  ArrowRight
} from "lucide-react"

export default function DashboardPage() {
  return (
    <div className="min-h-screen bg-[#F7F6F2]">
      <Header />

      <main className="mx-auto max-w-[1200px] px-6 py-8 space-y-10">
        {/* Срочные задачи */}
        <div className="space-y-4">
          <SectionBadge color="green" text="Срочные задачи: 1" />
          <div className="flex gap-4">
            <TaskCard
              project="Shyraq"
              time="18 часов"
              title="Нужно привести 10 плиток"
              desc="На объекте Shyraq срочно требуется новая плитка. 5×3 метра."
              status="green"
              icon={<AlertCircle className="h-4 w-4" />}
            />
          </div>
        </div>

        {/* Мои задачи */}
        <div className="space-y-4">
          <SectionHeader
            left="Мои задачи: 10"
            right="+ Добавить задачу"
            color="red"
          />
          <div className="flex gap-4 overflow-x-auto pb-2">
            <TaskCard
              project="Shyraq"
              time="-9 часов"
              title="Возведение колонн на 1 этаже"
              desc="Перед началом работ нужно провести подготовку."
              status="red"
              icon={<AlertCircle className="h-4 w-4" />}
            />
            <TaskCard
              project="Ansau"
              time="2 дня"
              title="Нужно сделать новую планировку"
              desc="На объекте Ansau срочно требуется новая планировка."
              status="green"
              icon={<CheckCircle2 className="h-4 w-4" />}
            />
            <TaskCard
              project="Dariya"
              time="20 дней"
              title="Нужно сделать новую планировку"
              desc="На объекте Dariya срочно требуется новая планировка."
              status="green"
              icon={<CheckCircle2 className="h-4 w-4" />}
            />
          </div>
        </div>

        {/* Проекты */}
        <div className="space-y-4">
          <SectionHeader
            left="Проекты: 3"
            right="+ Добавить проект"
            color="yellow"
          />
          <div className="flex gap-4 overflow-x-auto pb-2">
            <ProjectCard 
              title="Shyraq" 
              days="25 дней" 
              progress={65}
            />
            <ProjectCard 
              title="Ansau" 
              days="55 дней" 
              progress={45}
            />
            <ProjectCard 
              title="Dariya" 
              days="55 дней" 
              progress={30}
            />
          </div>
        </div>
      </main>
    </div>
  )
}

/* ================= HEADER ================= */

function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-white/80 backdrop-blur-sm">
      <div className="mx-auto flex max-w-[1200px] items-center justify-between px-6 py-4">
        {/* Logo */}
        <Link href="/" className="text-xl font-bold text-[#8B6B4E] hover:text-[#8B6B4E]/80 transition-colors">
          Qurylys
        </Link>

        {/* Nav */}
        <nav className="hidden md:flex items-center gap-8 text-sm">
          <Link 
            href="/calendar" 
            className="flex items-center gap-2 text-[#505872] hover:text-[#8B6B4E] transition-colors"
          >
            <Calendar className="h-4 w-4" />
            <span>Календарь</span>
          </Link>
          <Link 
            href="/hierarchy" 
            className="flex items-center gap-2 text-[#505872] hover:text-[#8B6B4E] transition-colors"
          >
            <Users className="h-4 w-4" />
            <span>Иерархия</span>
          </Link>
          <Link 
            href="/" 
            className="flex items-center gap-2 font-semibold text-[#8B6B4E] border-b-2 border-[#8B6B4E] pb-1"
          >
            <LayoutDashboard className="h-4 w-4" />
            <span>Дашборд</span>
          </Link>
          <Link 
            href="/lifecycle" 
            className="flex items-center gap-2 text-[#505872] hover:text-[#8B6B4E] transition-colors"
          >
            <FileText className="h-4 w-4" />
            <span>ЖЦП</span>
          </Link>
        </nav>

        {/* Right */}
        <div className="flex items-center gap-4">
          <Button 
            variant="ghost" 
            size="icon"
            className="relative text-[#505872] hover:text-[#8B6B4E]"
          >
            <Bell className="h-5 w-5" />
            <span className="absolute top-1 right-1 h-2 w-2 rounded-full bg-red-500" />
          </Button>
          <Avatar className="h-9 w-9 cursor-pointer hover:ring-2 ring-[#8B6B4E] transition-all">
            <AvatarFallback className="bg-[#C19A6B] text-white font-semibold">
              <User className="h-5 w-5" />
            </AvatarFallback>
          </Avatar>
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
  const colorMap = {
    green: "bg-green-50 text-green-700 border-green-200",
    red: "bg-red-50 text-red-700 border-red-200",
    yellow: "bg-yellow-50 text-yellow-700 border-yellow-200",
  }

  return (
    <Badge 
      variant="outline" 
      className={`${colorMap[color]} border px-4 py-1.5 text-sm font-semibold rounded-full`}
    >
      {text}
    </Badge>
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
      <Button 
        className="rounded-full bg-[#C19A6B] text-white hover:bg-[#C19A6B]/90 font-medium shadow-sm"
      >
        <Plus className="h-4 w-4 mr-2" />
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
  icon,
}: {
  project: string
  time: string
  title: string
  desc: string
  status: "red" | "green" | "yellow"
  icon: React.ReactNode
}) {
  const statusColorMap = {
    red: "bg-red-500 hover:bg-red-600 text-white",
    green: "bg-green-500 hover:bg-green-600 text-white",
    yellow: "bg-yellow-400 hover:bg-yellow-500 text-white",
  }

  const isOverdue = time.startsWith("-")

  return (
    <Card className="w-[320px] shrink-0 hover:shadow-lg transition-shadow duration-200 border border-gray-200">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="p-1.5 rounded-lg bg-[#F7F6F2]">
              {icon}
            </div>
            <span className="text-xs font-medium text-[#505872]">Проект: {project}</span>
          </div>
          <Badge className={`${statusColorMap[status]} border-0 text-xs font-medium`}>
            <Clock className="h-3 w-3 mr-1" />
            {time}
          </Badge>
        </div>
      </CardHeader>
      <CardContent className="pt-0 space-y-2">
        <CardTitle className="text-base font-semibold text-black leading-tight">
          {title}
        </CardTitle>
        <CardDescription className="text-sm text-[#505872] leading-relaxed">
          {desc}
        </CardDescription>
        {isOverdue && (
          <div className="pt-2">
            <Badge variant="destructive" className="text-xs">
              Просрочено
            </Badge>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

function ProjectCard({ 
  title, 
  days,
  progress 
}: { 
  title: string
  days: string
  progress: number
}) {
  return (
    <Card className="relative h-[200px] w-[280px] shrink-0 overflow-hidden border-0 bg-gradient-to-br from-gray-100 to-gray-200 hover:shadow-xl transition-shadow duration-200 cursor-pointer group">
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-5">
        <div className="absolute inset-0" style={{
          backgroundImage: `repeating-linear-gradient(45deg, transparent, transparent 10px, rgba(0,0,0,0.1) 10px, rgba(0,0,0,0.1) 20px)`
        }}></div>
      </div>

      {/* Content */}
      <div className="relative h-full p-4 flex flex-col justify-between">
        <div className="flex items-start justify-between">
          <Badge className="bg-black/80 text-white border-0 px-3 py-1.5 text-xs font-semibold backdrop-blur-sm">
            Проект: {title}
          </Badge>
          <Badge className="bg-[#C19A6B] text-white border-0 px-2.5 py-1.5 text-xs font-semibold hover:bg-[#C19A6B]/90">
            {days}
          </Badge>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs text-gray-700">
            <span className="font-medium">Прогресс</span>
            <span className="font-semibold">{progress}%</span>
          </div>
          <div className="h-2 bg-white/50 rounded-full overflow-hidden">
            <div 
              className="h-full bg-[#C19A6B] rounded-full transition-all duration-300"
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>

        {/* Hover Arrow */}
        <div className="absolute bottom-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
          <div className="p-2 rounded-full bg-white/90 backdrop-blur-sm">
            <ArrowRight className="h-4 w-4 text-[#8B6B4E]" />
          </div>
        </div>
      </div>
    </Card>
  )
}
