'use client';

import { Clock, Plus } from 'lucide-react';
import { useRouter } from 'next/navigation';

export default function Dashboard() {
  const router = useRouter();

  const sections = [
    {
      title: 'Срочные задачи',
      count: 1,
      badgeColor: 'bg-green-100',
      dotColor: 'bg-green-500',
      type: 'urgent',
      items: [
        {
          project: 'Shyraq',
          title: 'Нужно привести 10 плиток',
          description: 'На обьекте Shyraq срочно требуется новая плитка, 5х3 метра, 2 штуки для наличников вокруг л....',
          time: '18 часов',
          timeColor: 'text-green-500',
          borderColor: 'border-green-500',
        },
      ],
    },
    {
      title: 'Мои задачи',
      count: 10,
      badgeColor: 'bg-red-100',
      dotColor: 'bg-red-500',
      type: 'tasks',
      showAdd: true,
      items: [
        {
          project: 'Shyraq',
          title: 'Возведение колонн на 1 э...',
          description: 'Интерьеры подъездов и этажей\nПеред началом работ нужно провести подготовку...',
          time: '-9 часов',
          timeColor: 'text-red-500',
          borderColor: 'border-red-500',
        },
        {
          project: 'Ansau',
          title: 'Нужно сделать новую пла...',
          description: 'На обьекте Ansau срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ...',
          time: '2 дня',
          timeColor: 'text-green-500',
          borderColor: 'border-green-500',
        },
        {
          project: 'Dariya',
          title: 'Нужно сделать новую пла...',
          description: 'На обьекте Dariya срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ...',
          time: '20 дней',
          timeColor: 'text-green-500',
          borderColor: 'border-green-500',
        },
      ],
    },
    {
      title: 'Проекты',
      count: 3,
      badgeColor: 'bg-yellow-100',
      dotColor: 'bg-yellow-500',
      type: 'projects',
      showAdd: true,
      items: [
        { project: 'Shyraq', time: '25 дней', image: true },
        { project: 'Ansau', time: '55 дней', image: true },
        { project: 'Dariya', time: '55 дней', image: true },
      ],
    },
    {
      title: 'Задачи подчинённых',
      count: null,
      badgeColor: 'bg-red-100',
      dotColor: 'bg-red-500',
      type: 'subordinate-tasks',
      showAdd: true,
      items: [
        {
          project: 'Shyraq',
          title: 'Нужно сделать новую пла...',
          description: 'На обьекте Shyraq срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ...',
          time: '-9 часов',
          timeColor: 'text-red-500',
          borderColor: 'border-red-500',
        },
        {
          project: 'Ansau',
          title: 'Нужно сделать новую пла...',
          description: 'На обьекте Ansau срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ...',
          time: '2 дня',
          timeColor: 'text-green-500',
          borderColor: 'border-green-500',
        },
        {
          project: 'Dariya',
          title: 'Нужно сделать новую пла...',
          description: 'На обьекте Dariya срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ...',
          time: '20 дней',
          timeColor: 'text-green-500',
          borderColor: 'border-green-500',
        },
      ],
    },
    {
      title: 'Проекты подчинённых',
      count: 3,
      badgeColor: 'bg-yellow-100',
      dotColor: 'bg-yellow-500',
      type: 'subordinate-projects',
      showAdd: true,
      items: [
        { project: 'Shyraq', time: '25 дней', image: true },
        { project: 'Ansau', time: '55 дней', image: true },
        { project: 'Dariya', time: '55 дней', image: true },
      ],
    },
  ];

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4">
        <header className="mx-auto w-fit rounded-full border border-gray-200 bg-white px-8 py-3 shadow-sm">
          <div className="flex items-center justify-between gap-12">
            {/* Logo */}
            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-yellow-400 font-bold text-white text-xs">
                THE
              </div>
              <span className="text-sm font-semibold text-gray-800">QURYLS</span>
            </div>

            {/* Navigation */}
            <nav className="flex gap-8">
              <a href="#" className="text-sm text-gray-600 hover:text-gray-900">
                Календар
              </a>
              <a href="#" className="text-sm text-gray-600 hover:text-gray-900">
                Иерархия
              </a>
              <a href="#" className="border-b-2 border-black text-sm text-black font-bold">
                Дашборд
              </a>
              <button
                onClick={() => router.push('/lifecycle')}
                className="text-sm text-gray-600 hover:text-gray-900 transition-colors"
              >
                ЖЦП
              </button>
            </nav>

            {/* Right Icons */}
            <div className="flex items-center gap-4">
              <button className="text-gray-600 hover:text-gray-900">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
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
      <main className="max-w-6xl mx-auto px-6 py-8">
        {sections.map((section, idx) => (
          <div key={idx} className="mb-12">
            {/* Section Header */}
            <div className="flex items-center justify-between mb-6">
              <div className={`${section.badgeColor} rounded-full px-4 py-2 inline-flex items-center gap-2`}>
                <div className={`w-3 h-3 rounded-full ${section.dotColor}`}></div>
                <span className="text-sm font-semibold text-gray-900">
                  {section.title}{section.count !== null ? `: ${section.count}` : ''}
                </span>
              </div>
              {section.showAdd && (
                <button className={`${section.badgeColor} rounded-full px-6 py-2 inline-flex items-center gap-2 hover:opacity-80 transition-opacity`}>
                  <Plus className="w-4 h-4" />
                  <span className="text-sm font-semibold text-gray-900">
                    Добавить {section.type.includes('project') ? 'проект' : 'задачу'}
                  </span>
                </button>
              )}
            </div>

            {/* Items Container */}
            <div className="flex gap-6 overflow-x-auto pb-4">
              {section.items.map((item, itemIdx) => (
                <div key={itemIdx} className="flex-shrink-0">
                  {section.type.includes('project') ? (
                    // Project Card with Image
                    <div className="w-80 h-48 rounded-3xl overflow-hidden relative bg-gray-900 group cursor-pointer">
                      <img
                        src={`https://images.unsplash.com/photo-${itemIdx === 0 ? '1486325212027' : itemIdx === 1 ? '1486406146926' : '1486312338219'}?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80`}
                        alt={item.project}
                        className="w-full h-full object-cover group-hover:opacity-75 transition-opacity"
                      />
                      <div className="absolute inset-0 bg-black/40 flex flex-col justify-between p-6">
                        <div className="flex items-center justify-between">
                          <span className="text-white font-bold text-lg">Проект: {item.project}</span>
                        </div>
                        <div className="flex items-center justify-between">
                          <div className="bg-black/60 rounded-full px-3 py-1 flex items-center gap-2">
                            <Clock className="w-4 h-4 text-white" />
                            <span className="text-white text-sm font-semibold">{item.time}</span>
                          </div>
                          <div className="w-12 h-12 rounded-full bg-yellow-200 flex items-center justify-center text-gray-900 font-bold">
                            →
                          </div>
                        </div>
                      </div>
                    </div>
                  ) : (
                    // Task Card - Clickable
                    <button
                      onClick={() => router.push(`/project/${item.project.toLowerCase()}`)}
                      className="w-80 bg-white border-2 border-gray-200 rounded-3xl p-6 hover:shadow-md transition-shadow text-left"
                    >
                      <div className="flex items-start justify-between mb-3">
                        <span className="text-sm text-gray-600">Проект: {item.project}</span>
                        <div className={`bg-black rounded-full px-3 py-1 flex items-center gap-2 ${item.borderColor} border-2`}>
                          <Clock className={`w-4 h-4 ${item.timeColor}`} />
                          <span className="text-white text-xs font-semibold">{item.time}</span>
                        </div>
                      </div>
                      <h3 className="text-lg font-bold text-gray-900 mb-2">{item.title}</h3>
                      <p className="text-sm text-gray-600 mb-4 line-clamp-2">{item.description}</p>
                      <div className="flex items-center justify-between">
                        <span className="text-xs text-gray-500"></span>
                        <div className="w-10 h-10 rounded-full bg-yellow-200 flex items-center justify-center">
                          →
                        </div>
                      </div>
                    </button>
                  )}
                </div>
              ))}
            </div>
          </div>
        ))}
      </main>
    </div>
  );
}
