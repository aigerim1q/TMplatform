'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Send, FileText } from 'lucide-react';
import Header from '@/components/header';
import Sidebar from '@/components/sidebar';

export default function ChatPage() {
  const router = useRouter();
  const [messages, setMessages] = useState([
    {
      type: 'ai',
      text: 'Привет! Я твой AI помощник. Ты можешь спросить меня о любом проекте, и я дам тебе информацию. Попробуй написать например "Расскажи о проекте Shyraq"',
    },
  ]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const projects = {
    shyraq: {
      name: 'Shyraq',
      description: 'Искусственный интеллект проанализировал ваш документ и сформировал структуру жилищного цикла проекта.',
      fileName: 'Техническое_задание_Shyraq.pdf',
      fileSize: '2.4 MB',
    },
    ansau: {
      name: 'Ansau',
      description: 'AI обработал документацию проекта и подготовил анализ всех этапов работы.',
      fileName: 'Документация_Ansau.pdf',
      fileSize: '1.8 MB',
    },
    dariya: {
      name: 'Dariya',
      description: 'Проект успешно загружен и проанализирован AI системой.',
      fileName: 'Спецификация_Dariya.pdf',
      fileSize: '3.2 MB',
    },
  };

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (!input.trim()) return;

    // Add user message
    const userMessage = input.toLowerCase();
    setMessages((prev) => [...prev, { type: 'user', text: input }]);
    setInput('');
    setIsLoading(true);

    // Simulate AI response delay
    setTimeout(() => {
      let aiResponse = null;

      // Check if user asked about a specific project
      if (userMessage.includes('shyraq')) {
        aiResponse = {
          type: 'ai',
          text: 'Вот информация о проекте Shyraq:',
          project: 'shyraq',
        };
      } else if (userMessage.includes('ansau')) {
        aiResponse = {
          type: 'ai',
          text: 'Вот информация о проекте Ansau:',
          project: 'ansau',
        };
      } else if (userMessage.includes('dariya')) {
        aiResponse = {
          type: 'ai',
          text: 'Вот информация о проекте Dariya:',
          project: 'dariya',
        };
      } else {
        aiResponse = {
          type: 'ai',
          text: 'Я помогу тебе! Напиши название проекта, например "Shyraq", "Ansau" или "Dariya".',
        };
      }

      setMessages((prev) => [...prev, aiResponse]);
      setIsLoading(false);
    }, 500);
  };

  const handleProjectClick = (projectId) => {
    router.push(`/project/${projectId}/preview`);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="flex w-screen justify-center border-b border-gray-200 bg-gray-50 py-4">
        <Header />
      </div>

      {/* Main Layout with Sidebar */}
      <div className="flex">
        <Sidebar />
        
        {/* Chat Content */}
        <main className="flex-1 flex flex-col">
          {/* Chat Messages */}
          <div className="flex-1 overflow-y-auto p-6">
            <div className="max-w-3xl mx-auto space-y-4">
              {messages.map((message, idx) => (
                <div key={idx} className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}>
                  <div
                    className={`max-w-xl ${
                      message.type === 'user'
                        ? 'bg-purple-500 text-white rounded-3xl px-6 py-3'
                        : 'bg-white text-gray-900 rounded-3xl px-6 py-4 border border-gray-200 space-y-3'
                    }`}
                  >
                    <p className="text-sm">{message.text}</p>

                    {/* Project Document Card */}
                    {message.project && projects[message.project] && (
                      <button
                        onClick={() => handleProjectClick(message.project)}
                        className="w-full mt-3 border-2 border-purple-400 rounded-2xl p-4 bg-white hover:bg-purple-50 transition-colors text-left group"
                      >
                        <div className="flex items-start gap-3">
                          <div className="flex-shrink-0 mt-1">
                            <FileText className="w-8 h-8 text-purple-500" />
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="flex items-start justify-between">
                              <div>
                                <h4 className="font-bold text-gray-900 text-sm">
                                  Проект: {projects[message.project].name}
                                </h4>
                                <p className="text-xs text-gray-600 mt-1">
                                  {projects[message.project].description}
                                </p>
                              </div>
                              <svg
                                className="w-4 h-4 text-purple-500 flex-shrink-0 ml-2 group-hover:translate-x-1 transition-transform"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                              >
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                              </svg>
                            </div>
                            <div className="mt-2 flex items-center gap-2 text-xs text-gray-500">
                              <FileText className="w-3 h-3" />
                              <span>{projects[message.project].fileName}</span>
                              <span>{projects[message.project].fileSize}</span>
                            </div>
                          </div>
                        </div>
                      </button>
                    )}
                  </div>
                </div>
              ))}

              {isLoading && (
                <div className="flex justify-start">
                  <div className="bg-white border border-gray-200 rounded-3xl px-6 py-3">
                    <div className="flex gap-2">
                      <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                      <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce delay-100"></div>
                      <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce delay-200"></div>
                    </div>
                  </div>
                </div>
              )}

              <div ref={messagesEndRef} />
            </div>
          </div>

          {/* Input Area */}
          <div className="border-t border-gray-200 bg-white p-6">
            <form onSubmit={handleSendMessage} className="max-w-3xl mx-auto">
              <div className="relative">
                <input
                  type="text"
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  placeholder="Опишите задачу или задайте вопрос AI..."
                  className="w-full rounded-full border-2 border-red-300 px-6 py-3 pr-12 focus:outline-none focus:border-purple-500 transition-colors text-sm"
                />
                <button
                  type="submit"
                  disabled={!input.trim()}
                  className="absolute right-2 top-1/2 -translate-y-1/2 bg-gray-700 hover:bg-gray-800 disabled:bg-gray-300 text-white rounded-full p-2 transition-colors flex-shrink-0"
                >
                  <Send className="w-5 h-5" />
                </button>
              </div>
              <div className="mt-2 flex items-center justify-between text-xs text-gray-500 px-2">
                <div className="flex items-center gap-1">
                  🔒 Ваши данные защищены
                </div>
                <span>Нажмите Enter для отправки</span>
              </div>
            </form>
          </div>
        </main>
      </div>
    </div>
  );
}
