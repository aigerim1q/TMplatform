'use client';

import { useRouter, useParams } from 'next/navigation';
import { Clock, Users, AlertCircle, Plus, ChevronRight } from 'lucide-react';
import Header from '@/components/header';

export default function ProjectOverviewPage() {
  const router = useRouter();
  const params = useParams();

  const projectImages: { [key: string]: string } = {
    'shyraq': '/images/building-1.jpg',
    'ansau': '/images/building-2.jpg',
    'dariya': '/images/building-3.jpg',
  };

  const projectId = params.id as string;
  const projectImage = projectImages[projectId] || '/images/building-1.jpg';
  const projectName = projectId.charAt(0).toUpperCase() + projectId.slice(1);

  const projectData = {
    title: `Проект: ${projectName}`,
    image: projectImage,
    deadline: '12.12.2025 23:59 (25 дней)',
    responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    hasWarning: true,
    warningText: 'Причины отсрочки по всем дедлайнам',
  };

  const phases = [
    {
      id: 1,
      title: 'I. Предпроектная подготовка',
      tasks: [
        {
          project: `Проект: ${projectName}`,
          title: '1. Инициирование проекта',
          description: 'Определение целей и ограничений проекта. Анализ потребностей рынка...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          project: `Проект: ${projectName}`,
          title: '2. Финансово-экономический...',
          description: 'Прогноз стоимости строительства. Расчет рентабельности проекта (ROI)...',
          status: 'Задержка: 5...',
          statusColor: 'red',
        },
      ],
    },
    {
      id: 2,
      title: 'II. Проектирование (архитектура, инженерия, дизайн)',
      tasks: [
        {
          project: `Проект: ${projectName}`,
          title: '1. Архитектурная концепция',
          description: 'Конструктивный проект здания. Общие планы этажей. Количество квартир и их типы...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          project: `Проект: ${projectName}`,
          title: '2. Инженерные разделы',
          description: 'Конструктивные (фундамент, колонны, электроснабжение...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          project: `Проект: ${projectName}`,
          title: '3. Дизайн интерьера и фас...',
          description: 'Внешний вид здания с соседства. Интерьеры квартир и этажей проекта...',
          status: 'Выполнено',
          statusColor: 'green',
        },
      ],
    },
    {
      id: 3,
      title: 'III. Строительный этап',
      tasks: [
        {
          project: `Проект: ${projectName}`,
          title: '1. Подготовка площадки',
          description: 'Освещение участка. Установка бытовок и складов. Организация временного...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          project: `Проект: ${projectName}`,
          title: '2. Фундаментные работы',
          description: 'Геодезическая разбивка. Рытьё котлована. Подготовка основания...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          project: `Проект: ${projectName}`,
          title: '3. Возведение колонн на 1 э...',
          description: 'Перед началом работ нужно провести подготовку...',
          status: '9 часов',
          statusColor: 'orange',
        },
      ],
    },
  ];

  return (
    <div className="min-h-screen bg-white">
      {/* Header - centered */}
      <div className="flex justify-center pt-6">
        <Header />
      </div>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-6 py-8">
        {/* Navigation Buttons - All in one line */}
        <div className="flex items-center gap-3 mb-8">
          <button
            onClick={() => router.push('/dashboard')}
            className="px-4 py-2 rounded-full bg-gray-100 text-gray-900 hover:bg-gray-200 transition-colors text-sm font-medium"
          >
            ← Назад
          </button>
          <button className="bg-black text-white px-6 py-2 rounded-full text-sm font-semibold">
            Задача
          </button>
          <button 
            onClick={() => router.push(`/project/${params.id}/reports`)}
            className="bg-gray-400 text-white px-6 py-2 rounded-full text-sm font-semibold hover:bg-gray-500 transition-colors"
          >
            Отчеты
          </button>
        </div>

        {/* Project Image */}
        <div className="mb-8">
          <img
            src={projectImage || "/placeholder.svg"}
            alt="Project"
            className="w-80 h-48 rounded-3xl object-cover shadow-md"
          />
        </div>

        {/* Project Title */}
        <h1 className="text-2xl font-bold text-gray-900 mb-6">{projectData.title}</h1>

        {/* Project Info Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          {/* Deadline Card */}
          <div className="bg-yellow-100 rounded-3xl p-4 flex items-start gap-3">
            <Clock className="w-6 h-6 text-yellow-600 flex-shrink-0 mt-1" />
            <div>
              <p className="text-xs text-gray-600 mb-1">Дедлайн:</p>
              <p className="text-sm font-semibold text-gray-900">{projectData.deadline}</p>
            </div>
          </div>

          {/* Responsibility Card */}
          <div className="bg-white border-2 border-gray-300 rounded-3xl p-4 flex items-start gap-3">
            <Users className="w-6 h-6 text-gray-600 flex-shrink-0 mt-1" />
            <div>
              <p className="text-xs text-gray-600 mb-1">Ответственные:</p>
              <p className="text-sm font-semibold text-gray-900">{projectData.responsible}</p>
            </div>
          </div>

          {/* Warning Card */}
          {projectData.hasWarning && (
            <div className="bg-red-100 rounded-3xl p-4 flex items-start gap-3">
              <AlertCircle className="w-6 h-6 text-red-600 flex-shrink-0 mt-1" />
              <div>
                <p className="text-xs text-gray-600 mb-1">Внимание</p>
                <p className="text-sm font-semibold text-gray-900">{projectData.warningText}</p>
              </div>
            </div>
          )}
        </div>

        {/* Phases */}
        {phases.map((phase) => (
          <div key={phase.id} className="mb-12">
            {/* Phase Header */}
            <h2 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
              <div className="bg-black text-white rounded-full px-3 py-1 text-xs font-bold">
                {phase.id}
              </div>
              {phase.title}
            </h2>

            {/* Phase Tasks */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
              {phase.tasks.map((task, taskIdx) => (
                <div key={taskIdx} className="bg-white border-2 border-gray-200 rounded-2xl p-4 hover:shadow-md transition-shadow">
                  {/* Project name and status */}
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-xs text-gray-600">{task.project}</span>
                    <div
                      className={`px-3 py-1 rounded-full text-xs font-semibold ${
                        task.statusColor === 'green'
                          ? 'bg-green-100 text-green-700'
                          : task.statusColor === 'red'
                          ? 'bg-red-100 text-red-700'
                          : 'bg-orange-100 text-orange-700'
                      }`}
                    >
                      {task.status === 'Выполнено' && 'Выполнено'}
                      {task.status.includes('Задержка') && task.status}
                      {task.status.includes('часов') && task.status}
                    </div>
                  </div>

                  {/* Task title */}
                  <h3 className="font-bold text-gray-900 mb-2 text-sm">{task.title}</h3>

                  {/* Task description */}
                  <p className="text-xs text-gray-600 line-clamp-2 mb-3">{task.description}</p>
                </div>
              ))}
            </div>
          </div>
        ))}

        {/* Budget Section */}
        <div className="bg-white rounded-2xl p-6 mb-8 border border-gray-200">
          <h3 className="text-lg font-bold text-gray-900 mb-4">Бюджет проекта</h3>
          <div className="flex items-end justify-between mb-6">
            <div>
              <p className="text-xs text-gray-600 mb-2">ОСВОЕНО СРЕДСТВ</p>
              <p className="text-3xl font-bold text-gray-900">
                1,500,900,000 <span className="text-lg text-gray-600">₸</span>
              </p>
            </div>
            <p className="text-sm text-gray-600">2,400,800,000 ₸</p>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div className="bg-yellow-600 h-2 rounded-full" style={{ width: '62.5%' }}></div>
          </div>
        </div>

        {/* Add Phase Button */}
        <div className="flex justify-center">
          <button className="bg-yellow-600 hover:bg-yellow-700 text-white px-6 py-3 rounded-full font-semibold transition-colors flex items-center gap-2">
            <Plus className="w-5 h-5" />
            Добавить этап проекта
          </button>
        </div>
      </main>
    </div>
  );
}
