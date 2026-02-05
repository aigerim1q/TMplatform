'use client';

import { useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { ChevronRight, Clock, Users, AlertCircle, Send, Flag, Share2, Copy, X, Search, UserPlus, Calendar, ChevronDown, Paperclip, Star } from 'lucide-react';
import Header from '@/components/header';

export default function TaskDetail() {
  const router = useRouter();
  const params = useParams();
  const [isDelegateModalOpen, setIsDelegateModalOpen] = useState(false);
  const [selectedPerson, setSelectedPerson] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [delegateDueDate, setDelegateDueDate] = useState('');
  const [delegatePriority, setDelegatePriority] = useState('Высокий');
  const [delegateComment, setDelegateComment] = useState('');
  const [activeTab, setActiveTab] = useState<'comments' | 'history'>('comments');

  const historyItems = [
    {
      id: 1,
      type: 'expense',
      user: 'Вы',
      action: 'добавили новый расход',
      detail: 'Арматура A500C',
      time: '2 минуты назад',
      avatar: 'you',
    },
    {
      id: 2,
      type: 'deadline',
      user: 'Евгений С.',
      action: 'обновил дедлайн задачи',
      oldDate: '24.11.2025',
      newDate: '26.01.2026',
      time: '15 минут назад',
      avatar: 'evgeniy',
    },
    {
      id: 3,
      type: 'delegate',
      user: '',
      action: 'Автоматически делегирована задача',
      detail: '"Подготовка опалубки" перешла к следующему этапу согласно графику',
      time: '1 час назад',
      avatar: 'system',
    },
    {
      id: 4,
      type: 'status',
      user: 'Серик Р.',
      action: 'изменил статус',
      oldStatus: 'Ожидание',
      newStatus: 'В работе',
      time: '3 часа назад',
      avatar: 'serik',
    },
    {
      id: 5,
      type: 'file',
      user: 'Омар А.',
      action: 'прикрепил файл',
      fileName: 'smetka_v3',
      time: 'Вчера, 18:30',
      avatar: 'omar',
    },
    {
      id: 6,
      type: 'created',
      user: '',
      action: 'Задача создано в системе',
      time: '22.12.2025',
      avatar: 'system',
    },
  ];

  const teamMembers = [
    { id: '1', name: 'Алия К.', role: 'Инженер ПТО', avatar: 'aliya', status: 'available', statusText: 'Свободна: 4ч сегодня', recommended: true, match: 98 },
    { id: '2', name: 'Данияр С.', role: 'Снабжение', avatar: 'daniyar', status: 'busy', statusText: 'Занят до 16:00', recommended: false, match: 0 },
    { id: '3', name: 'Ержан Б.', role: 'Прораб', avatar: 'erzhan', status: 'offline', statusText: 'Оффлайн', recommended: false, match: 0 },
  ];

  const filteredMembers = teamMembers.filter(member =>
    member.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    member.role.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const taskData = {
    title: 'Возведение колонн на 1 этаже несущих конструкции',
    deadline: '26.11.2025 23:59 (-9 часов)',
    startDate: 'Дата начала: 15.11.2025 12:00',
    responsible: ['Омар Ахмет', 'Зейнулла Ршыман', 'Серик Рах...'],
    issue: 'Причина просрочки: Айды Рахимбаев болел 5 дней',
    preparation: [
      'Проверка геодезической разбивки осей и отметок',
      'Подготовка и выравнивание опубликованных для колонн',
      'Проверка качества арматурных каркасов (диаметр, шаг, фиксация)',
      'Затем переидем к задачам:',
    ],
    stages: [
      {
        title: 'Пересмотреть и доработать чертежи',
        description: 'Пересмотреть расположение всех существующих квартир\nПереработать конфигурацию несущих стен',
        status: 'Выполнено',
      },
      {
        title: 'С новой конфигурацией добавить...',
        description: 'Добавить больше квартир за счет оптимизации пространства и корректировки нежиных зон...',
        status: 'Выполнено',
      },
      {
        title: 'Провести проверку новой планировки',
        description: 'Проверить соответствие новой планировки требованиям пожарной безопасности',
        status: 'Ошибка',
        days: '-9 часов',
      },
    ],
    comments: [
      {
        author: 'Зейнулла Ршыман',
        time: '1 ч.з.',
        text: 'А почему она может быть крива?\nМы же проверили уровни на прошлой неделе.',
      },
      {
        author: 'Айды Р',
        role: 'изменил статус на',
        status: 'в работе',
        time: '',
      },
      {
        author: 'Там была проблема с виниграми. Используемые ножи заменены сегодня.',
        text: 'Там была проблема с виниграми. Используемые ножи заменены сегодня.',
        time: 'Прошедшая запись',
      },
    ],
  };

  return (
    <div className="min-h-screen bg-white text-gray-900 dark:bg-slate-950 dark:text-slate-100">
      {/* Header - centered */}
      <div className="flex justify-center pt-6">
        <Header />
      </div>

      <div className="max-w-7xl mx-auto px-6 py-8">
        {/* Top Navigation - All Buttons in One Line */}
        <div className="flex items-center gap-3 mb-8">
          <button
            onClick={() => router.push('/')}
            className="flex items-center gap-2 px-4 py-2 rounded-full bg-gray-100 text-gray-900 hover:bg-gray-200 transition-colors dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
          >
            ← Назад
          </button>
          
          <button className="bg-black text-white px-6 py-2 rounded-full text-sm font-semibold dark:bg-amber-500 dark:text-slate-950">
            Задача
          </button>
          <button 
            onClick={() => router.push(`/project/${params.id}/reports`)}
            className="bg-gray-400 text-white px-6 py-2 rounded-full text-sm font-semibold hover:bg-gray-500 transition-colors dark:bg-slate-700 dark:hover:bg-slate-600"
          >
            Отчеты
          </button>
        </div>

        {/* Title */}
        <h1 className="text-3xl font-bold text-gray-900 mb-8 dark:text-white">{taskData.title}</h1>

        {/* Info Cards */}
        <div className="grid grid-cols-3 gap-4 mb-8">
          {/* Deadline Card */}
          <div className="bg-red-100 rounded-2xl p-4 dark:bg-red-500/15">
            <div className="flex items-start gap-3">
              <Clock className="w-6 h-6 text-red-600 flex-shrink-0 mt-1 dark:text-red-300" />
              <div>
                <p className="text-sm text-gray-700 font-semibold dark:text-slate-100">{taskData.deadline}</p>
                <p className="text-sm text-gray-600 dark:text-slate-300">{taskData.startDate}</p>
              </div>
            </div>
          </div>

          {/* Responsible Card */}
          <div className="bg-white border-2 border-gray-300 rounded-2xl p-4 dark:bg-slate-900 dark:border-slate-800">
            <div className="flex items-start gap-3">
              <Users className="w-6 h-6 text-gray-700 flex-shrink-0 mt-1 dark:text-slate-200" />
              <div>
                <p className="text-sm text-gray-900 font-semibold dark:text-slate-100">Ответственные: {taskData.responsible.join(', ')}</p>
                <button className="text-blue-600 text-sm mt-1 hover:underline dark:text-blue-300">↗</button>
              </div>
            </div>
          </div>

          {/* Status Card */}
          <div className="bg-black rounded-2xl p-4 text-white dark:bg-slate-800 dark:text-slate-100">
            <div className="flex items-start gap-3">
              <AlertCircle className="w-6 h-6 flex-shrink-0 mt-1" />
              <div>
                <p className="text-sm font-semibold">{taskData.issue}</p>
                <button className="text-gray-300 text-sm mt-1 hover:underline dark:text-slate-300">↗</button>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content */}
        <div className="grid grid-cols-3 gap-8">
          {/* Left Column */}
          <div className="col-span-2">
            {/* Preparation Section */}
            <div className="mb-8">
              <h2 className="text-lg font-bold text-gray-900 mb-4 dark:text-white">Подготовка</h2>
              <ul className="space-y-2">
                {taskData.preparation.map((item, idx) => (
                  <li key={idx} className={`text-sm ${idx === 1 ? 'bg-yellow-200 px-3 py-2 rounded dark:bg-amber-500/20' : ''} text-gray-700 dark:text-slate-200`}>
                    • {item}
                  </li>
                ))}
              </ul>
            </div>

            {/* Stages Section */}
            <div>
              <h2 className="text-lg font-bold text-gray-900 mb-4 dark:text-white">Этапы выполнения</h2>
              <div className="space-y-4">
                {taskData.stages.map((stage, idx) => (
                  <div key={idx} className="border border-gray-200 rounded-xl p-4 dark:border-slate-800 dark:bg-slate-900/70">
                    <div className="flex items-start gap-3 mb-3">
                      <div className={`w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 ${
                        stage.status === 'Выполнено' ? 'bg-green-100' : 'bg-orange-100'
                      }`}>
                        <div className={`w-4 h-4 rounded-full ${
                          stage.status === 'Выполнено' ? 'bg-green-500' : 'bg-orange-500'
                        }`} />
                      </div>
                      <div className="flex-1">
                        <h3 className="font-semibold text-gray-900 dark:text-white">{stage.title}</h3>
                        <p className="text-sm text-gray-600 mt-2 dark:text-slate-300">{stage.description}</p>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className={`text-xs font-semibold px-3 py-1 rounded-full ${
                          stage.status === 'Выполнено' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                        }`}>
                          {stage.status}
                        </span>
                        {stage.days && <span className="text-xs text-red-600 font-semibold dark:text-red-300">{stage.days}</span>}
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* Add Stage Button */}
              <button className="w-full mt-6 py-2 px-4 border-2 border-yellow-600 text-yellow-600 rounded-full font-semibold hover:bg-yellow-50 transition-colors dark:border-amber-400 dark:text-amber-300 dark:hover:bg-amber-400/20">
                + Добавить этап проекта
              </button>
            </div>
          </div>

          {/* Right Column */}
          <div className="col-span-1">
            {/* Action Buttons */}
            <div className="bg-yellow-600 rounded-2xl p-4 text-white font-semibold text-center cursor-pointer hover:bg-yellow-700 transition-colors mb-4 dark:bg-amber-500 dark:text-slate-950">
              ✓ Завершить задачу
            </div>

            <div className="grid grid-cols-2 gap-3 mb-6">
              <button className="py-3 px-4 border-2 border-gray-300 rounded-lg text-gray-900 font-semibold hover:bg-gray-50 transition-colors dark:border-slate-700 dark:text-slate-100 dark:hover:bg-slate-800">
                Отложить
              </button>
              <button 
                onClick={() => setIsDelegateModalOpen(true)}
                className="py-3 px-4 border-2 border-gray-300 rounded-lg text-gray-900 font-semibold hover:bg-gray-50 transition-colors dark:border-slate-700 dark:text-slate-100 dark:hover:bg-slate-800"
              >
                Делегировать
              </button>
            </div>

            {/* Comments Section */}
            <div className="bg-gray-50 rounded-3xl p-6 border border-gray-200 dark:bg-slate-900 dark:border-slate-800">
              <div className="flex gap-6 mb-6">
                <button 
                  onClick={() => setActiveTab('comments')}
                  className={`font-semibold text-base pb-1 transition-all ${
                    activeTab === 'comments' 
                      ? 'text-gray-900 border-b-2 border-gray-900 dark:text-white dark:border-white' 
                      : 'text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200'
                  }`}
                >
                  Комментарии
                </button>
                <button 
                  onClick={() => setActiveTab('history')}
                  className={`font-semibold text-base pb-1 transition-all ${
                    activeTab === 'history' 
                      ? 'text-gray-900 border-b-2 border-gray-900 dark:text-white dark:border-white' 
                      : 'text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200'
                  }`}
                >
                  История
                </button>
              </div>

              {activeTab === 'comments' ? (
                <>
                  <div className="space-y-4 mb-6">
                    {/* Date Separator */}
                    <div className="flex items-center gap-3 py-2">
                      <div className="flex-1 h-px bg-gray-300 dark:bg-slate-700"></div>
                      <span className="text-gray-400 text-sm font-medium dark:text-slate-400">Сегодня</span>
                      <div className="flex-1 h-px bg-gray-300 dark:bg-slate-700"></div>
                    </div>

                    {/* First Comment */}
                    <div>
                      <div className="flex items-start gap-3 mb-2">
                        <img 
                          src="https://api.dicebear.com/7.x/avataaars/svg?seed=user1" 
                          alt="avatar" 
                          className="w-10 h-10 rounded-full"
                        />
                        <div>
                          <div className="flex items-center gap-2">
                            <p className="font-semibold text-gray-900 dark:text-slate-100">Зейнулла Ршыман</p>
                            <p className="text-gray-500 text-xs dark:text-slate-400">14:32</p>
                          </div>
                          <p className="text-gray-700 text-sm mt-1 dark:text-slate-200">А почему она может быть крива?<br />Мы же проверили уровни на<br />прошлой неделе.</p>
                        </div>
                      </div>
                    </div>

                    {/* Status Change Separator */}
                    <div className="flex items-center gap-3 py-3 mt-4">
                      <div className="flex-1 h-px bg-gray-300 dark:bg-slate-700"></div>
                      <span className="text-gray-400 text-sm font-medium dark:text-slate-400">Статус изменен</span>
                      <div className="flex-1 h-px bg-gray-300 dark:bg-slate-700"></div>
                    </div>

                    {/* Status Update */}
                    <div className="flex items-center gap-2 mb-4">
                      <img 
                        src="https://api.dicebear.com/7.x/avataaars/svg?seed=user2" 
                        alt="avatar" 
                        className="w-6 h-6 rounded-full"
                      />
                      <p className="text-gray-600 text-sm dark:text-slate-300">Айдын Р. изменил статус на <span className="text-yellow-600 font-semibold dark:text-amber-300">В работе</span></p>
                    </div>

                    {/* Your Comment */}
                    <div className="mt-4">
                      <div className="flex items-end justify-end gap-2 mb-1">
                        <p className="text-gray-500 text-xs dark:text-slate-400">14:45</p>
                        <span className="text-gray-600 text-xs font-semibold dark:text-slate-300">Вы</span>
                        <img 
                          src="https://i.pinimg.com/736x/23/00/64/230064c1d688b4553c244292c8bb3220.jpg" 
                          alt="avatar" 
                          className="w-8 h-8 rounded-full"
                        />
                      </div>
                      <div className="bg-yellow-500 text-white rounded-2xl p-4 max-w-xs ml-auto dark:bg-amber-500 dark:text-slate-950">
                        <p className="text-sm">Там была проблема с фиксаторами. Исправили, но нужен доп. контроль сегодня.</p>
                      </div>
                      <p className="text-gray-500 text-xs text-right mt-2 dark:text-slate-400">Просмотрено</p>
                    </div>
                  </div>

                  {/* Comment Input */}
                  <div className="mt-8">
                    <div className="flex gap-3 items-end mb-3">
                      <input
                        type="text"
                        placeholder="Напишать комментарий..."
                        className="flex-1 bg-white border border-gray-200 rounded-full px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-yellow-500 placeholder-gray-400 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100 dark:placeholder-slate-500"
                      />
                      <button className="bg-yellow-500 text-white w-10 h-10 rounded-full hover:bg-yellow-600 transition-colors flex items-center justify-center dark:bg-amber-500 dark:hover:bg-amber-400">
                        <Send className="w-5 h-5" />
                      </button>
                    </div>
                    <div className="flex gap-4 px-4">
                      <button className="text-gray-600 text-sm hover:text-gray-900 transition-colors flex items-center gap-1 dark:text-slate-300 dark:hover:text-white">
                        📎 Файл
                      </button>
                      <button className="text-gray-600 text-sm hover:text-gray-900 transition-colors flex items-center gap-1 dark:text-slate-300 dark:hover:text-white">
                        📷 Фото
                      </button>
                      <button className="text-gray-600 text-sm hover:text-gray-900 transition-colors flex items-center gap-1 dark:text-slate-300 dark:hover:text-white">
                        @ Упомянуть
                      </button>
                    </div>
                  </div>
                </>
              ) : (
                /* History Tab Content */
                <div className="space-y-1">
                  {historyItems.map((item, idx) => (
                    <div key={item.id} className="flex items-start gap-3 py-3">
                      {/* Timeline line */}
                      <div className="flex flex-col items-center">
                        {item.type === 'delegate' ? (
                          <div className="w-8 h-8 rounded-full bg-amber-100 flex items-center justify-center">
                            <Star className="w-4 h-4 text-amber-600" />
                          </div>
                        ) : item.type === 'created' ? (
                          <div className="w-8 h-8 rounded-full bg-red-100 flex items-center justify-center">
                            <span className="text-red-600 font-bold text-sm">Q</span>
                          </div>
                        ) : (
                          <img 
                            src={`https://picsum.photos/seed/${item.avatar}/32/32`}
                            alt="avatar" 
                            className="w-8 h-8 rounded-full object-cover"
                          />
                        )}
                        {idx < historyItems.length - 1 && (
                          <div className="w-px h-full min-h-8 bg-gray-200 mt-2" />
                        )}
                      </div>

                      {/* Content */}
                      <div className="flex-1 min-w-0">
                        {item.type === 'expense' && (
                          <>
                            <p className="text-sm text-gray-900">
                              <span className="font-semibold">{item.user}</span> {item.action}{' '}
                              <span className="text-amber-700">{item.detail}</span>
                            </p>
                            <p className="text-xs text-gray-400 mt-1">{item.time}</p>
                          </>
                        )}

                        {item.type === 'deadline' && (
                          <>
                            <p className="text-sm text-gray-900">
                              <span className="font-semibold">{item.user}</span> {item.action}
                            </p>
                            <p className="text-sm mt-1">
                              <span className="text-gray-500">{item.oldDate}</span>
                              <span className="text-gray-400 mx-2">→</span>
                              <span className="text-amber-600">{item.newDate}</span>
                            </p>
                            <p className="text-xs text-gray-400 mt-1">{item.time}</p>
                          </>
                        )}

                        {item.type === 'delegate' && (
                          <>
                            <p className="text-sm font-semibold text-gray-900">{item.action}</p>
                            <div className="bg-gray-100 rounded-lg p-2 mt-1">
                              <p className="text-xs text-gray-600 italic">{item.detail}</p>
                            </div>
                            <p className="text-xs text-gray-400 mt-1">{item.time}</p>
                          </>
                        )}

                        {item.type === 'status' && (
                          <>
                            <p className="text-sm text-gray-900">
                              <span className="font-semibold">{item.user}</span> {item.action}
                            </p>
                            <div className="flex items-center gap-2 mt-1">
                              <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">{item.oldStatus}</span>
                              <span className="text-gray-400">→</span>
                              <span className="text-xs text-amber-800 bg-amber-100 px-2 py-1 rounded">{item.newStatus}</span>
                            </div>
                            <p className="text-xs text-gray-400 mt-1">{item.time}</p>
                          </>
                        )}

                        {item.type === 'file' && (
                          <>
                            <p className="text-sm text-gray-900">
                              <span className="font-semibold">{item.user}</span> {item.action}{' '}
                              <span className="text-blue-600">{item.fileName}</span>
                            </p>
                            <p className="text-xs text-gray-400 mt-1">{item.time}</p>
                          </>
                        )}

                        {item.type === 'created' && (
                          <>
                            <p className="text-sm font-semibold text-gray-900">{item.action}</p>
                            <p className="text-xs text-gray-400 mt-1">{item.time}</p>
                          </>
                        )}
                      </div>
                    </div>
                  ))}

                  {/* Show all history button */}
                  <button className="w-full mt-4 py-3 text-gray-500 text-sm font-medium hover:text-gray-700 transition-colors flex items-center justify-center gap-2">
                    <Clock className="w-4 h-4" />
                    Показать всю историю
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Delegate Modal */}
        {isDelegateModalOpen && (
          <div 
            className="fixed inset-0 z-50 flex items-center justify-center"
            onClick={() => setIsDelegateModalOpen(false)}
          >
            {/* Dark overlay */}
            <div className="absolute inset-0 bg-black/40" />
            
            {/* Modal content */}
            <div 
              className="relative bg-white rounded-2xl w-full max-w-3xl mx-4 shadow-2xl"
              onClick={(e) => e.stopPropagation()}
            >
              {/* Modal header */}
              <div className="flex items-center justify-between p-6 pb-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center">
                    <UserPlus className="w-5 h-5 text-gray-600" />
                  </div>
                  <h2 className="text-xl font-bold text-gray-900">Делегировать задачу</h2>
                </div>
                <button
                  onClick={() => setIsDelegateModalOpen(false)}
                  className="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 transition-colors"
                >
                  <X className="w-5 h-5 text-gray-500" />
                </button>
              </div>

              {/* Modal body - Two columns */}
              <div className="flex gap-8 px-6 pb-6">
                {/* Left column - Selected task */}
                <div className="w-64 flex-shrink-0">
                  <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">Выбранная задача</p>
                  
                  <div className="bg-white border border-gray-200 rounded-xl p-4 mb-4">
                    <div className="flex items-center gap-2 mb-2">
                      <div className="w-2 h-2 rounded-full bg-yellow-500" />
                      <span className="text-sm text-gray-600">Проект: Shyraq</span>
                    </div>
                    <p className="font-semibold text-gray-900 text-sm leading-tight mb-3">Нужно привести 10 плиток</p>
                    <div className="flex items-center gap-2 text-gray-500 text-xs">
                      <Calendar className="w-3.5 h-3.5" />
                      <span>Создано: Сегодня</span>
                    </div>
                  </div>

                  {/* Due date */}
                  <div className="mb-4">
                    <p className="text-sm text-gray-600 mb-2">Срок выполнения</p>
                    <div className="relative">
                      <input
                        type="text"
                        placeholder="mm/dd/yyyy"
                        value={delegateDueDate}
                        onChange={(e) => setDelegateDueDate(e.target.value)}
                        className="w-full px-3 py-2.5 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-yellow-500 placeholder-gray-400"
                      />
                      <Calendar className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                    </div>
                  </div>

                  {/* Priority */}
                  <div className="mb-4">
                    <p className="text-sm text-gray-600 mb-2">Приоритет</p>
                    <div className="relative">
                      <select
                        value={delegatePriority}
                        onChange={(e) => setDelegatePriority(e.target.value)}
                        className="w-full px-3 py-2.5 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-yellow-500 appearance-none cursor-pointer"
                      >
                        <option value="Высокий">Высокий</option>
                        <option value="Средний">Средний</option>
                        <option value="Низкий">Низкий</option>
                      </select>
                      <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
                    </div>
                  </div>

                  {/* Comment */}
                  <div>
                    <p className="text-sm text-gray-600 mb-2">Комментарий к задаче</p>
                    <textarea
                      placeholder="Добавьте инструкции..."
                      value={delegateComment}
                      onChange={(e) => setDelegateComment(e.target.value)}
                      className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-yellow-500 placeholder-gray-400 resize-none h-20"
                    />
                  </div>
                </div>

                {/* Right column - Select executor */}
                <div className="flex-1">
                  <p className="text-base font-semibold text-gray-900 mb-3">Выберите исполнителя</p>
                  
                  {/* Search input */}
                  <div className="relative mb-4">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <input
                      type="text"
                      placeholder="Поиск сотрудника..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="w-full pl-10 pr-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-yellow-500 placeholder-gray-400"
                    />
                  </div>

                  {/* Team members list */}
                  <div className="space-y-2 max-h-64 overflow-y-auto">
                    {filteredMembers.map((member) => (
                      <div
                        key={member.id}
                        onClick={() => setSelectedPerson(member.id)}
                        className={`flex items-center gap-3 p-3 rounded-xl cursor-pointer transition-colors border ${
                          selectedPerson === member.id
                            ? 'bg-green-50 border-green-200'
                            : 'bg-white border-gray-100 hover:bg-gray-50'
                        }`}
                      >
                        {/* Radio button */}
                        <div className={`w-5 h-5 rounded-full border-2 flex items-center justify-center flex-shrink-0 ${
                          selectedPerson === member.id ? 'border-green-500' : 'border-gray-300'
                        }`}>
                          {selectedPerson === member.id && (
                            <div className="w-2.5 h-2.5 rounded-full bg-green-500" />
                          )}
                        </div>

                        {/* Avatar with status */}
                        <div className="relative">
                          <img
                            src={`https://picsum.photos/seed/${member.avatar}/40/40`}
                            alt={member.name}
                            className="w-10 h-10 rounded-full object-cover"
                          />
                          {member.status === 'available' && (
                            <div className="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-500 rounded-full border-2 border-white" />
                          )}
                        </div>

                        {/* Name and status */}
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2">
                            <p className="font-semibold text-gray-900 text-sm">{member.name}</p>
                            {member.recommended && (
                              <span className="px-2 py-0.5 bg-green-100 text-green-700 text-xs font-medium rounded-full">
                                Рекомендовано
                              </span>
                            )}
                          </div>
                          <p className="text-xs text-gray-500">{member.statusText}</p>
                        </div>

                        {/* Match percentage and role */}
                        <div className="text-right flex-shrink-0">
                          {member.match > 0 && (
                            <p className="text-sm font-semibold text-green-600">{member.match}% совпадение</p>
                          )}
                          <p className="text-xs text-gray-400">{member.role}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Modal footer */}
              <div className="flex items-center justify-end gap-4 px-6 py-4 border-t border-gray-100">
                <button
                  onClick={() => setIsDelegateModalOpen(false)}
                  className="px-6 py-2.5 text-gray-700 font-semibold hover:bg-gray-50 rounded-lg transition-colors"
                >
                  Отмена
                </button>
                <button
                  onClick={() => {
                    if (selectedPerson) {
                      setIsDelegateModalOpen(false);
                      setSelectedPerson(null);
                      setSearchQuery('');
                      setDelegateComment('');
                      setDelegateDueDate('');
                    }
                  }}
                  className="flex items-center gap-2 px-6 py-2.5 bg-amber-600 text-white font-semibold rounded-full hover:bg-amber-700 transition-colors"
                >
                  Делегировать
                  <Send className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
