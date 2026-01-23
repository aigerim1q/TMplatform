'use client';

import { ChevronRight, CheckCircle2, Zap, Pencil } from 'lucide-react';
import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function MainContent() {
  const router = useRouter();
  const [selectedCard, setSelectedCard] = useState(null);

  const cards = [
    {
      id: 'construction',
      title: 'Полный цикл ЖК',
      description: 'Стандартный шаблон для возведения жилых комплексов. Включает этапы от котлована до сдачи в эксплуатацию и снятелстройства.',
      tag: 'Строительство',
      tagColor: 'bg-purple-100 text-purple-600',
      items: ['24 этапа работ', 'Авто-генерация сметы', 'Контроль порядников'],
      used: 'Использовано 120 раз.',
      icon: '🔨',
    },
    {
      id: 'renovation',
      title: 'Капитальный ремонт',
      description: 'Оптимизирован для ремонтных работ в существующих зданиях. Фокус на демонтаже, отделке и инженерных сетях.',
      tag: 'Реновация',
      tagColor: 'bg-blue-100 text-blue-600',
      items: ['15 этапов работ', 'Учет материалов', ''],
      used: 'Использовано 85 раз.',
      icon: '🔧',
    },
    {
      id: 'architecture',
      title: 'Архитектурный проект',
      description: 'Фокус на создании чертежей, получении разрешений и согласовании документации. Идеально для пред-строительного этапа.',
      tag: 'Проектирование',
      tagColor: 'bg-amber-100 text-amber-600',
      items: ['Согласование с гос. органами', 'BIM интеграция', ''],
      used: 'Использовано 40 раз.',
      icon: '📋',
    },
  ];

  return (
    <main className="w-full flex flex-col items-center px-6 py-12">
      {/* AI Assistant Badge */}
      <div className="mb-8 flex items-center gap-2 rounded-full bg-purple-100 px-4 py-2">
        <svg className="h-5 w-5 text-purple-600" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
        <span className="text-sm font-semibold text-purple-600">AI ASSISTANT</span>
      </div>

      {/* Main Heading */}
      <h1 className="mb-3 text-center text-4xl font-bold text-gray-900">
        Выберите жизненный цикл проекта
      </h1>
      <p className="mb-8 max-w-2xl text-center text-gray-600">
        Наш ИИ поможет вам настроить структуру проекта. Выберите шаблон, который лучше всего соответствует вашим текущим задачам.
      </p>

      {/* Recommended Badge */}
      <div className="mb-8 flex items-center gap-2 rounded-full bg-gray-900 px-4 py-2">
        <div className="h-2 w-2 rounded-full bg-green-500"></div>
        <span className="text-sm font-medium text-white">Рекомендуемые шаблоны</span>
      </div>

      {/* Cards Grid */}
      <div className="mb-12 grid w-full max-w-4xl grid-cols-3 gap-6">
        {cards.map((card) => (
          <button
            key={card.id}
            onClick={() => setSelectedCard(card.id)}
            className="relative flex flex-col overflow-hidden rounded-lg border-2 bg-white p-6 text-left transition-all duration-200"
            style={{
              borderColor: selectedCard === card.id ? '#D4AF37' : '#E5E7EB',
              boxShadow: selectedCard === card.id ? '0 0 0 3px rgba(212, 175, 55, 0.1)' : 'none',
            }}
          >
            {/* Checkmark Badge */}
            {selectedCard === card.id && (
              <div className="absolute top-3 left-3">
                <CheckCircle2 className="h-6 w-6 text-yellow-500" />
              </div>
            )}

            {/* Tag */}
            <div className={`mb-4 inline-block rounded-full px-3 py-1 text-xs font-semibold w-fit ${card.tagColor}`}>
              {card.tag}
            </div>

            {/* Title */}
            <h3 className="mb-2 text-lg font-bold text-gray-900">{card.title}</h3>

            {/* Description */}
            <p className="mb-4 flex-1 text-sm text-gray-600">{card.description}</p>

            {/* Items List */}
            <div className="mb-4 space-y-2">
              {card.items.map((item, idx) => (
                item && (
                  <div key={idx} className="flex items-center gap-2 text-sm text-gray-600">
                    <div className="h-2 w-2 rounded-full bg-green-500"></div>
                    {item}
                  </div>
                )
              ))}
            </div>

            {/* Footer */}
            <div className="flex items-center justify-between border-t border-gray-200 pt-4">
              <span className="text-xs text-gray-500">{card.used}</span>
              <ChevronRight className="h-5 w-5 text-gray-400" />
            </div>
          </button>
        ))}
      </div>

      {/* Bottom Section */}
      <div className="w-full max-w-4xl">
        {/* Action Buttons */}
        <div className="flex items-center justify-center gap-4 mb-12">
          <button
            onClick={() => selectedCard && router.push('/chat')}
            style={{
              opacity: selectedCard ? 1 : 0.6,
            }}
            className="rounded-full bg-amber-100 px-8 py-3 text-sm font-semibold text-amber-900 transition-opacity duration-200 hover:bg-amber-200 disabled:cursor-not-allowed"
            disabled={!selectedCard}
          >
            Продолжить →
          </button>
          <div className="flex items-center gap-2 rounded-full border border-gray-300 bg-white px-6 py-3">
            <span className="text-sm text-gray-600">Выбран:</span>
            <span className="text-sm font-semibold text-gray-900">{selectedCard ? cards.find(c => c.id === selectedCard)?.title : '?'}</span>
          </div>
        </div>

        {/* Specific Tasks Section */}
        <div className="mb-8">
          <div className="mb-6 flex items-center gap-2 rounded-full bg-green-100 px-4 py-2 w-fit">
            <Zap className="h-4 w-4 text-green-600" />
            <span className="text-sm font-semibold text-green-600">Специфические задачи</span>
          </div>

          {/* Task Items */}
          <div className="space-y-4">
            {[
              {
                icon: Pencil,
                title: 'Ландшафтный дизайн',
                desc: 'Планировка территории, озеленение, дорожки',
                complexity: 'СЛОЖНОСТЬ',
                dots: 3,
              },
              {
                icon: Zap,
                title: 'Электросети',
                desc: 'Монтаж проводки, щитков, подключение к сети',
                complexity: 'СЛОЖНОСТЬ',
                dots: 3,
              },
              {
                icon: Pencil,
                title: 'Свой шаблон с нуля',
                desc: 'Опишите задачу, AI он составит план',
                hasAI: true,
              },
            ].map((task, idx) => {
              const TaskIcon = task.icon;
              return (
                <div key={idx} className="flex items-center justify-between rounded-lg border border-gray-200 bg-white p-4 hover:bg-gray-50 transition-colors">
                  <div className="flex items-center gap-4">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-gray-300 bg-gray-50">
                      <TaskIcon className="h-5 w-5 text-gray-600" />
                    </div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{task.title}</h4>
                      <p className="text-sm text-gray-600">{task.desc}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    {task.hasAI ? (
                      <button className="rounded-full bg-purple-600 px-4 py-1 text-sm font-semibold text-white hover:bg-purple-700">
                        AI Генератор
                      </button>
                    ) : (
                      <div className="text-right">
                        <p className="text-xs text-gray-500">{task.complexity}</p>
                        <div className="flex gap-1">
                          {[...Array(task.dots)].map((_, i) => (
                            <div
                              key={i}
                              className="h-2 w-2 rounded-full"
                              style={{
                                backgroundColor: i < 2 ? '#FFA500' : (i < 3 ? '#FFB84D' : '#FFD700'),
                              }}
                            ></div>
                          ))}
                        </div>
                      </div>
                    )}
                    <ChevronRight className="h-5 w-5 text-gray-400" />
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </main>
  );
}
