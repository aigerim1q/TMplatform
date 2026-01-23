'use client';

import { useRouter, useParams } from 'next/navigation';
import { Plus, Clock, AlertCircle, ArrowLeft, CheckCircle2, Edit3, Download } from 'lucide-react';

export default function ProjectDetail() {
  const router = useRouter();
  const params = useParams();

  const projectPhases = [
    {
      phase: 'I. Предпроектная подготовка',
      tasks: [
        {
          number: 1,
          title: 'Инициирование проекта',
          description: 'Определение ориентировочной площади и этажности. Анализ потребностей рынка...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          number: 2,
          title: 'Финансово-экономический анализ',
          description: 'Прогноз стоимости строительства проекта (РСМ) и ренгабельность проекта',
          status: 'Задержка: 5...',
          statusColor: 'red',
        },
      ],
    },
    {
      phase: 'II. Проектирование (архитектура, инженерия, дизайн)',
      tasks: [
        {
          number: 1,
          title: 'Архитектурная концепция',
          description: 'Концептуальный проект зданий. Общие планы этажей. Количество квартир и их типы...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          number: 2,
          title: 'Инженерные разделы',
          description: 'Конструкты (фундамент, колонны, плиты) Электроснабжение...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          number: 3,
          title: 'Дизайн интерьера и фас...',
          description: 'Интерьеры подъезда и этажей. Интерьеры квартир (всех бизнес-класса)...',
          status: 'Выполнено',
          statusColor: 'green',
          daysLeft: '25 дней',
        },
      ],
    },
    {
      phase: 'III. Строительный этап',
      tasks: [
        {
          number: 1,
          title: 'Подготовка площадки',
          description: 'Ограждение участка. Установка бытовок и складов. Организация временного...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          number: 2,
          title: 'Фундаментные работы',
          description: 'Геодезическая разбивка. Раытье котлована. Подготовка основания...',
          status: 'Выполнено',
          statusColor: 'green',
        },
        {
          number: 3,
          title: 'Возведение колонн на 1 э...',
          description: 'Интерьеры подъезда и этажей. Перед началом работ нужно провести подготовку...',
          status: '-9 часов',
          statusColor: 'red',
        },
      ],
    },
  ];

  const budgetData = {
    allocated: 1500900000,
    total: 2400800000,
  };

  const projects = {
    shyraq: {
      name: 'Shyraq',
      image: '/shyraq.png',
      aiStatus: 'Обработано',
      title: 'Проект Shyraq',
      description: 'Детальное описание проекта Shyraq...',
      document: {
        name: 'document.pdf',
        size: '5MB',
      },
      deadline: '2024-01-01',
      model: 'Модель A',
      budget: '2,400,800,000',
      stages: [
        { title: 'Планирование', duration: '30 дней' },
        { title: 'Проектирование', duration: '60 дней' },
        { title: 'Строительство', duration: '180 дней' },
        { title: 'Монтаж', duration: '90 дней' },
        { title: 'Проверка', duration: '30 дней' },
        { title: 'Ввод в эксплуатацию и передача заказчику', duration: '5 дней' },
      ],
      team: [
        { name: 'Омар Аммет', role: 'Руководитель ПРОЕКТА' },
        { name: 'Расул Даулетов', role: 'ТЕХНИЧЕСКИЙ ДИРЕКТОР' },
        { name: 'Айдын Рахимбаев', role: 'Главный инженер' },
      ],
      completion: 98,
    },
  };

  const project = projects[params.id] || projects.shyraq;

  return (
    <div className="min-h-screen bg-white">
      {/* Header - centered */}
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4">
        <header className="mx-auto w-fit rounded-full border border-gray-200 bg-white px-8 py-3 shadow-sm">
          <div className="flex items-center justify-between gap-12">
            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-yellow-400 font-bold text-white text-xs">
                THE
              </div>
              <span className="text-sm font-semibold text-gray-800">QURYLS</span>
            </div>

            <nav className="flex gap-8">
              <a href="#" className="text-sm text-gray-600 hover:text-gray-900">
                Календарь
              </a>
              <a href="#" className="text-sm text-gray-600 hover:text-gray-900">
                Иерархия
              </a>
              <button
                onClick={() => router.push('/')}
                className="text-sm font-bold border-b-2 border-black text-black"
              >
                Дашборд
              </button>
              <button
                onClick={() => router.push('/lifecycle')}
                className="text-sm text-gray-600 hover:text-gray-900"
              >
                ЖЦП
              </button>
            </nav>

            <div className="flex items-center gap-4">
              <button className="text-gray-600 hover:text-gray-900">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                </svg>
              </button>
              <button className="text-gray-600 hover:text-gray-900">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                </svg>
              </button>
              <button className="h-8 w-8 overflow-hidden rounded-full bg-gray-300">
                <img
                  src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                  alt="Avatar"
                  className="h-full w-full object-cover"
                />
              </button>
            </div>
          </div>
        </header>
      </div>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-8">
        {/* Top Navigation with Tabs */}
        <div className="flex items-center justify-between mb-8">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-2 px-4 py-2 rounded-full bg-gray-100 text-gray-900 hover:bg-gray-200 transition-colors"
          >
            ← Назад
          </button>
          
          <div className="flex gap-3">
            <button className="bg-black text-white px-6 py-2 rounded-full text-sm font-semibold">
              Задача
            </button>
            <button className="bg-gray-400 text-white px-6 py-2 rounded-full text-sm font-semibold hover:bg-gray-500">
              Отчеты
            </button>
          </div>
        </div>

        {/* Project Info */}
        <div className="mb-8">
          <div className="flex gap-6 mb-6">
            <img
              src="https://images.unsplash.com/photo-1486325212027-8081e485255e?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&h=200"
              alt="Project"
              className="w-48 h-40 rounded-3xl object-cover"
            />
            <div className="flex-1">
              <h1 className="text-3xl font-bold text-gray-900 mb-4">Проект: Shyraq</h1>
              
              <div className="grid grid-cols-3 gap-4 mb-6">
                {/* Deadline Card */}
                <div className="bg-yellow-100 rounded-2xl p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Clock className="w-5 h-5 text-gray-900" />
                    <span className="text-sm font-semibold text-gray-900">Дедлайн: 12.12.2025 23:59 (25 дней)</span>
                  </div>
                  <p className="text-xs text-gray-700">Дата начала: 9.06.2025 12:00</p>
                </div>

                {/* Responsible Card */}
                <div className="bg-gray-100 rounded-2xl p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-sm font-semibold text-gray-900">Ответственные: Омар Ахмет, Зейнулла Рыщман, Серик Рақ...</span>
                  </div>
                  <button className="text-purple-600 text-xs font-semibold hover:underline">
                    Редактировать →
                  </button>
                </div>

                {/* Priority Card */}
                <div className="bg-black text-white rounded-2xl p-4 flex items-center justify-between">
                  <span className="text-sm font-semibold">Приоритетные отсрочки по всем дедлайнам</span>
                  <AlertCircle className="w-5 h-5" />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Project Phases */}
        <div className="space-y-8">
          {projectPhases.map((phaseGroup, phaseIdx) => (
            <div key={phaseIdx}>
              {/* Phase Title */}
              <div className="bg-black text-white rounded-full px-6 py-3 inline-block mb-4 font-semibold text-sm">
                {phaseGroup.phase}
              </div>

              {/* Phase Tasks Grid */}
              <div className="grid grid-cols-2 gap-6 mb-6">
                {phaseGroup.tasks.map((task, taskIdx) => (
                  <div key={taskIdx} className="bg-white border border-gray-200 rounded-2xl p-6 hover:shadow-md transition-shadow">
                    <div className="flex items-start justify-between mb-3">
                      <div>
                        <p className="text-gray-600 text-sm">Проект: Shyraq</p>
                        {task.statusColor === 'green' ? (
                          <span className="inline-block mt-2 bg-green-100 text-green-700 px-3 py-1 rounded-full text-xs font-semibold">
                            ✓ Выполнено
                          </span>
                        ) : (
                          <span className="inline-block mt-2 bg-red-100 text-red-700 px-3 py-1 rounded-full text-xs font-semibold">
                            ● {task.status}
                          </span>
                        )}
                      </div>
                      {task.daysLeft && (
                        <span className="text-xs text-gray-600 font-semibold">{task.daysLeft}</span>
                      )}
                    </div>
                    <h3 className="font-bold text-gray-900 mb-2">{task.number}. {task.title}</h3>
                    <p className="text-sm text-gray-600 line-clamp-2">{task.description}</p>
                  </div>
                ))}
              </div>

              {/* Add Task Button */}
              {phaseIdx === 1 && (
                <div className="text-right mb-8">
                  <button className="inline-flex items-center gap-2 text-purple-600 hover:text-purple-700 font-semibold">
                    <Plus className="w-5 h-5" /> Добавить задачу
                  </button>
                </div>
              )}

              {/* Budget Section for middle phase */}
              {phaseIdx === 1 && (
                <div className="bg-gray-100 rounded-2xl p-6 mb-8">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="font-semibold text-gray-900 flex items-center gap-2">
                      📊 Бюджет проекта
                    </h3>
                    <span className="text-sm text-gray-600">Отчет расходов →</span>
                  </div>
                  <div className="flex items-baseline gap-4">
                    <div>
                      <p className="text-gray-600 text-sm">Выделено</p>
                      <p className="text-2xl font-bold text-gray-900">{budgetData.allocated.toLocaleString()}</p>
                    </div>
                    <span className="text-gray-400">/</span>
                    <div>
                      <p className="text-gray-600 text-sm">Всего</p>
                      <p className="text-xl font-bold text-gray-900">{budgetData.total.toLocaleString()} ₸</p>
                    </div>
                  </div>
                  <div className="w-full bg-gray-300 rounded-full h-2 mt-4">
                    <div className="bg-yellow-400 h-2 rounded-full" style={{ width: `${(budgetData.allocated / budgetData.total) * 100}%` }}></div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Add Stage Button */}
        <div className="text-center mt-12">
          <button className="bg-yellow-200 hover:bg-yellow-300 text-gray-900 px-8 py-3 rounded-full font-semibold transition-colors">
            + Добавить этап проекта
          </button>
        </div>

        {/* Floating AI Icon */}
        <div className="fixed bottom-8 right-8 w-16 h-16 bg-yellow-200 rounded-full flex items-center justify-center cursor-pointer hover:bg-yellow-300 transition-colors shadow-lg">
          <span className="text-2xl">🤖</span>
        </div>
      </main>
    </div>
  );
}
