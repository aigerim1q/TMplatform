'use client';

import { useRouter, useParams } from 'next/navigation';
import { CheckCircle2, Edit2 } from 'lucide-react';

export default function ProjectPreview() {
  const router = useRouter();
  const params = useParams();

  const projectData = {
    shyraq: {
      name: 'Shyraq',
      image: 'https://images.unsplash.com/photo-1486325212027-8081e485255e?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&h=200',
      description: 'Искусственный интеллект проанализировал ваш документ и сформировал структуру жилищного цикла проекта. Проверьте данные ниже.',
      document: {
        name: 'Техническое_задание_Shyraq.pdf',
        size: '2.4 MB',
      },
      deadline: '15 сентября 2025',
      model: 'Смета (Fixed Price)',
      stages: [
        { number: 1, title: 'Подготовительный этап и мобилизация ресурсов', duration: '14 дней' },
        { number: 2, title: 'Разработка и утверждение проектно-сметной документации', duration: '45 дней' },
        { number: 3, title: 'Закупка материалов и оборудования', duration: '30 дней' },
        { number: 4, title: 'Строительно-монтажные работы (СМР)', duration: '180 дней' },
        { number: 5, title: 'Пусконаладочные работы и тестирование систем', duration: '21 день' },
        { number: 6, title: 'Ввод в эксплуатацию и передача заказчику', duration: '10 дней' },
      ],
      team: [
        { name: 'Омар Ахмет', role: 'РУКОВОДИТЕЛЬ ПРОЕКТА', avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop' },
        { name: 'Расул Даулетов', role: 'ТЕХНИЧЕСКИЙ ДИРЕКТОР', avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&h=150&fit=crop' },
        { name: 'Айдын Рахимбаев', role: 'ГЛАВНЫЙ ИНЖЕНЕР', avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=150&h=150&fit=crop' },
      ],
    },
    ansau: {
      name: 'Ansau',
      image: 'https://images.unsplash.com/photo-1486406146926-c627a92ad1ab?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&h=200',
      description: 'Проект успешно загружен и обработан AI системой.',
      document: { name: 'Документация_Ansau.pdf', size: '1.8 MB' },
      deadline: '20 октября 2025',
      model: 'Смета (Fixed Price)',
      stages: [
        { number: 1, title: 'Подготовка площадки', duration: '10 дней' },
        { number: 2, title: 'Проектирование', duration: '30 дней' },
        { number: 3, title: 'Согласование', duration: '15 дней' },
        { number: 4, title: 'Строительство', duration: '120 дней' },
        { number: 5, title: 'Отделочные работы', duration: '45 дней' },
        { number: 6, title: 'Сдача объекта', duration: '5 дней' },
      ],
      team: [
        { name: 'Айтурган Сатпаева', role: 'РУКОВОДИТЕЛЬ ПРОЕКТА', avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop' },
        { name: 'Мария Иванова', role: 'АРХИТЕКТОР', avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&h=150&fit=crop' },
        { name: 'Максим Петров', role: 'ИНЖЕНЕР', avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=150&h=150&fit=crop' },
      ],
    },
    dariya: {
      name: 'Dariya',
      image: 'https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&h=200',
      description: 'Проект успешно загружен и проанализирован AI системой.',
      document: { name: 'Спецификация_Dariya.pdf', size: '3.2 MB' },
      deadline: '25 ноября 2025',
      model: 'Смета (Fixed Price)',
      stages: [
        { number: 1, title: 'Инициирование', duration: '7 дней' },
        { number: 2, title: 'Планирование', duration: '20 дней' },
        { number: 3, title: 'Проектирование', duration: '40 дней' },
        { number: 4, title: 'Строительство', duration: '100 дней' },
        { number: 5, title: 'Тестирование', duration: '15 дней' },
        { number: 6, title: 'Запуск', duration: '7 дней' },
      ],
      team: [
        { name: 'Дарья Смирнова', role: 'РУКОВОДИТЕЛЬ ПРОЕКТА', avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop' },
        { name: 'Артем Морозов', role: 'АНАЛИТИК', avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&h=150&fit=crop' },
        { name: 'Елена Волкова', role: 'КООРДИНАТОР', avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=150&h=150&fit=crop' },
      ],
    },
  };

  const project = projectData[params.id] || projectData.shyraq;

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
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
              <a href="#" className="text-sm text-gray-600 hover:text-gray-900">Календар</a>
              <a href="#" className="text-sm text-gray-600 hover:text-gray-900">Иерархия</a>
              <button onClick={() => router.push('/')} className="text-sm text-gray-600 hover:text-gray-900 font-bold">Дашборд</button>
              <button onClick={() => router.push('/lifecycle')} className="text-sm text-gray-600 hover:text-gray-900">ЖЦП</button>
            </nav>
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
                <img src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" alt="Avatar" className="h-full w-full object-cover" />
              </button>
            </div>
          </div>
        </header>
      </div>

      {/* Main Content */}
      <main className="max-w-5xl mx-auto px-6 py-8">
        {/* Back Button */}
        <button onClick={() => router.back()} className="flex items-center gap-2 mb-6 px-4 py-2 rounded-full bg-yellow-100 text-gray-900 hover:bg-yellow-200 transition-colors text-sm font-semibold">
          ← Назад
        </button>

        {/* Project Header */}
        <div className="flex gap-6 mb-8">
          <img src={project.image || "/placeholder.svg"} alt={project.name} className="w-40 h-40 rounded-2xl object-cover" />
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-2 h-2 rounded-full bg-purple-500"></div>
              <span className="text-xs font-semibold text-purple-600">AI Обработка завершена</span>
            </div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Проект: {project.name}</h1>
            <p className="text-gray-600 text-sm">{project.description}</p>
          </div>
        </div>

        {/* Content Grid */}
        <div className="grid grid-cols-3 gap-6">
          {/* Left Column */}
          <div className="col-span-2 space-y-6">
            {/* Document Card */}
            <div className="border-2 border-purple-300 rounded-3xl p-6 bg-purple-50">
              <div className="flex items-center gap-3">
                <CheckCircle2 className="w-6 h-6 text-purple-600 flex-shrink-0" />
                <div className="flex-1">
                  <h3 className="font-bold text-gray-900 mb-1">AI успешно обработал ваш документ</h3>
                  <div className="flex items-center gap-2 text-sm text-gray-600 mb-3">
                    <span className="text-red-500">📄</span>
                    <span>{project.document.name}</span>
                    <span className="text-gray-500">{project.document.size}</span>
                  </div>
                </div>
                <button className="text-gray-400 hover:text-gray-600 text-sm">⟲ Заменить файл</button>
              </div>
            </div>

            {/* Parameters */}
            <div className="border border-gray-200 rounded-3xl p-6">
              <div className="flex items-center justify-between mb-6">
                <h3 className="font-bold text-gray-900 flex items-center gap-2">
                  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M9 2a1 1 0 000 2h2V2H9zm0 1h2v13H9V3zm4-1a1 1 0 000 2h2V2h-2zm0 1h2v13h-2V3zm4-1a1 1 0 100 2h2V2h-2zm0 1h2v13h-2V3z" /></svg>
                  Параметры из документа
                </h3>
                <span className="text-xs font-semibold text-gray-500 uppercase">Автозаполнение</span>
              </div>

              <div className="grid grid-cols-2 gap-6 mb-6">
                <div>
                  <span className="text-xs text-gray-500 uppercase">Срок завершения (дедлайн)</span>
                  <p className="text-lg font-bold text-gray-900 mt-2">{project.deadline}</p>
                  <p className="text-xs text-gray-500 mt-1">Примерно за 360 дней</p>
                </div>
                <div>
                  <span className="text-xs text-gray-500 uppercase">Финансовая модель</span>
                  <p className="text-lg font-bold text-gray-900 mt-2">{project.model}</p>
                  <p className="text-xs text-gray-500 mt-1">Определено по 12 проектов</p>
                </div>
              </div>

              <p className="text-xs text-gray-500 mb-4">ЭТАПЫ ПРОЕКТА (ЖЦП)</p>
              <p className="text-right text-xs text-gray-500 mb-4">Надаю этапов: 6</p>

              {/* Stages List */}
              <div className="space-y-3">
                {project.stages.map((stage) => (
                  <div key={stage.number} className="flex items-center gap-4 p-3 border border-gray-200 rounded-lg hover:bg-gray-50 group">
                    <span className="text-sm font-semibold text-gray-600 min-w-6">{stage.number}</span>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900">{stage.title}</p>
                      <p className="text-xs text-gray-500">Срок: {stage.duration}</p>
                    </div>
                    <button className="text-gray-400 hover:text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity">
                      <Edit2 className="w-4 h-4" />
                    </button>
                  </div>
                ))}
              </div>

              {/* Add Stage Button */}
              <button className="w-full mt-4 text-center text-gray-400 hover:text-gray-600 text-sm py-2 border border-dashed border-gray-300 rounded-lg transition-colors">
                + Добавить новый этап
              </button>
            </div>

            {/* Action Buttons */}
            <div className="flex gap-4 mt-8">
              <button className="flex-1 bg-yellow-200 hover:bg-yellow-300 text-gray-900 font-semibold py-3 rounded-2xl transition-colors">
                Подтвердить и продолжить ✓
              </button>
              <button 
                onClick={() => router.push(`/project/${params.id}/edit`)}
                className="flex-1 border border-gray-300 hover:bg-gray-50 text-gray-900 font-semibold py-3 rounded-2xl transition-colors"
              >
                Редактировать
              </button>
            </div>
          </div>

          {/* Right Column */}
          <div className="space-y-6">
            {/* Team Recommendations */}
            <div className="bg-white border border-gray-200 rounded-2xl p-6">
              <h3 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
                Рекомендуемые ответственные
                <span className="text-blue-500">✓</span>
              </h3>
              <div className="space-y-3">
                {project.team.map((member, idx) => (
                  <div key={idx} className="flex items-center gap-3">
                    <img src={member.avatar || "/placeholder.svg"} alt={member.name} className="w-10 h-10 rounded-full object-cover" />
                    <div className="min-w-0">
                      <p className="text-sm font-semibold text-gray-900 truncate">{member.name}</p>
                      <p className="text-xs text-gray-500 uppercase">{member.role}</p>
                    </div>
                  </div>
                ))}
              </div>
              <p className="text-xs text-gray-600 mt-4">AI подобрал сотрудников на основе опыта в аналогичных проектах и текущей загрузки.</p>
            </div>

            {/* Control Accuracy */}
            <div className="bg-white border border-gray-200 rounded-2xl p-6">
              <h3 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
                💖 Контроль точности
              </h3>
              <p className="text-sm text-gray-600 mb-3">Точность распознавания</p>
              <div className="w-full bg-gray-300 rounded-full h-2 mb-4">
                <div className="bg-yellow-400 h-2 rounded-full" style={{ width: '98%' }}></div>
              </div>
              <p className="text-xs text-gray-500 text-center">98%</p>
              <p className="text-xs text-gray-600 mt-4">Все этапы были согласованы с внутренним регламентом качества.</p>
              <button className="w-full mt-4 bg-yellow-200 hover:bg-yellow-300 text-gray-900 font-semibold py-2 rounded-full text-sm transition-colors">
                ⬇ Скачать структуру ЖЦП
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
