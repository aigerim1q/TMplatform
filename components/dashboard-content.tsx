'use client';

import { PenTool, AlertTriangle, Users, ArrowUp, Lock } from 'lucide-react';
import { useState } from 'react';

export default function DashboardContent() {
  const [query, setQuery] = useState('');

  const suggestions = [
    {
      icon: PenTool,
      title: 'Спланировать проект',
      description: 'Создам структуру этапов и задач для нового объекта',
    },
    {
      icon: AlertTriangle,
      title: 'Проанализировать риски',
      description: 'Найду слабые места в текущем графике поставок',
    },
    {
      icon: Users,
      title: 'Распределить задачу',
      description: 'Оптимально начну исполнителей по компетенциям',
    },
  ];

  return (
    <main className="flex-1 flex flex-col items-center px-6 py-12">
      {/* Sparkle Icon */}
      <div className="mb-8 rounded-2xl bg-gradient-to-br from-purple-400 to-purple-600 p-4">
        <svg
          className="h-12 w-12 text-white"
          fill="currentColor"
          viewBox="0 0 24 24"
        >
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
      </div>

      {/* Main Heading */}
      <h1 className="mb-2 text-center text-4xl font-bold text-gray-900">
        Чем я могу помочь вам сегодня?
      </h1>
      <p className="mb-12 text-center text-gray-600">
        Задайте вопрос или выберите одно из предложений ниже
      </p>

      {/* Suggestion Cards */}
      <div className="mb-12 grid w-full max-w-4xl grid-cols-3 gap-4">
        {suggestions.map((suggestion, index) => {
          const Icon = suggestion.icon;
          return (
            <button
              key={index}
              className="rounded-lg border border-gray-200 bg-white p-6 text-left hover:border-gray-300 hover:shadow-sm transition-all"
            >
              <div className="mb-3 inline-block rounded-lg bg-purple-100 p-2">
                <Icon className="h-6 w-6 text-purple-600" />
              </div>
              <h3 className="mb-2 font-semibold text-gray-900">
                {suggestion.title}
              </h3>
              <p className="text-sm text-gray-600">{suggestion.description}</p>
            </button>
          );
        })}
      </div>

      {/* Input Area */}
      <div className="w-full max-w-2xl">
        <div className="relative">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Опишите задачу или задайте вопрос AI..."
            className="w-full rounded-full border-2 border-gray-300 bg-white px-6 py-3 pr-12 text-gray-900 placeholder-gray-400 focus:border-purple-500 focus:outline-none"
          />
          <button className="absolute right-2 top-1/2 -translate-y-1/2 rounded-full bg-gray-700 p-2 text-white hover:bg-gray-800">
            <ArrowUp size={20} />
          </button>
        </div>

        {/* Bottom Text */}
        <div className="mt-3 flex items-center justify-between text-xs text-gray-500">
          <div className="flex items-center gap-1">
            <Lock size={14} />
            Ваши данные защищены
          </div>
          <span>Нажмите Enter для отправки</span>
        </div>
      </div>
    </main>
  );
}
