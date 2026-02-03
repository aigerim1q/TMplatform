"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { FileText, Send } from "lucide-react";
import Header from "@/components/header";
import Sidebar from "@/components/sidebar";

type ProjectKey = 'shyraq' | 'ansau' | 'dariya';

type ChatMessage = {
  type: 'ai' | 'user';
  text: string;
  project?: ProjectKey;
};

export default function ChatPage() {
  const router = useRouter();
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      type: "ai",
      text: "Привет! Я твой AI помощник. Ты можешь спросить меня о любом проекте, и я дам тебе информацию. Попробуй написать например \"Расскажи о проекте Shyraq\"",
    },
  ]);
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement | null>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const projects: Record<ProjectKey, {
    name: string;
    description: string;
    fileName: string;
    fileSize: string;
  }> = {
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

  const handleSendMessage = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!input.trim()) return;

    const userMessage = input.toLowerCase();
    setMessages((prev) => [...prev, { type: "user", text: input }]);
    setInput("");
    setIsLoading(true);

    setTimeout(() => {
      let aiResponse: ChatMessage;

      if (userMessage.includes("shyraq")) {
        aiResponse = {
          type: "ai",
          text: "Вот информация о проекте Shyraq:",
          project: "shyraq",
        };
      } else if (userMessage.includes("ansau")) {
        aiResponse = {
          type: "ai",
          text: "Вот информация о проекте Ansau:",
          project: "ansau",
        };
      } else if (userMessage.includes("dariya")) {
        aiResponse = {
          type: "ai",
          text: "Вот информация о проекте Dariya:",
          project: "dariya",
        };
      } else {
        aiResponse = {
          type: "ai",
          text: "Я помогу тебе! Напиши название проекта, например \"Shyraq\", \"Ansau\" или \"Dariya\".",
        };
      }

      setMessages((prev) => [...prev, aiResponse]);
      setIsLoading(false);
    }, 500);
  };

  const handleProjectClick = (projectId: string) => {
    router.push(`/project/${projectId}/preview`);
  };

  return (
    <div className="min-h-screen bg-gray-50 text-slate-900 dark:bg-slate-950 dark:text-slate-50">
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4 dark:border-slate-800 dark:bg-slate-950">
        <Header />
      </div>

      <div className="flex">
        <Sidebar />

        <main className="flex-1">
          <div className="mx-auto flex min-h-[calc(100vh-140px)] max-w-5xl flex-col px-6 py-8 lg:px-10">
            <div className="mb-6 flex justify-center">
              <div className="relative max-w-3xl rounded-3xl bg-white px-6 py-4 text-sm leading-relaxed shadow-sm ring-1 ring-slate-200 dark:bg-slate-900 dark:ring-slate-800">
                Привет! Я твой AI помощник. Ты можешь спросить меня о любом проекте, и я дам тебе информацию. Попробуй написать, например, «Расскажи о проекте Shyraq».
                <div className="mt-3 inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                  ⏱ Время ответа: ~4 минуты
                </div>
              </div>
            </div>

            <div className="flex-1 space-y-4 overflow-y-auto rounded-3xl bg-white/90 p-6 shadow-sm ring-1 ring-slate-200 dark:bg-slate-900/80 dark:ring-slate-800">
              {messages.map((message, idx) => (
                <div
                  key={idx}
                  className={`flex ${message.type === "user" ? "justify-end" : "justify-start"}`}
                >
                  <div
                    className={`max-w-2xl rounded-3xl px-5 py-4 text-sm leading-relaxed shadow-sm transition ${
                      message.type === "user"
                        ? "bg-amber-400 text-slate-900 dark:bg-amber-500 dark:text-slate-950"
                        : "bg-slate-100 text-slate-900 ring-1 ring-slate-200 dark:bg-slate-800 dark:text-slate-50 dark:ring-slate-700"
                    }`}
                  >
                    <p>{message.text}</p>

                    {message.project && projects[message.project] && (
                      <button
                        onClick={() => message.project && handleProjectClick(message.project)}
                        className="mt-4 w-full rounded-2xl border border-amber-300 bg-white p-4 text-left transition hover:-translate-y-0.5 hover:shadow-sm dark:border-amber-500/60 dark:bg-slate-900"
                      >
                        <div className="flex items-start gap-3">
                          <div className="mt-1 flex h-10 w-10 items-center justify-center rounded-full bg-amber-100 text-amber-800 dark:bg-amber-500/20 dark:text-amber-100">
                            <FileText className="h-5 w-5" />
                          </div>
                          <div className="min-w-0 flex-1 space-y-1 text-sm">
                            <div className="flex items-start justify-between gap-3">
                              <div>
                                <h4 className="font-semibold text-slate-900 dark:text-white">Проект: {projects[message.project].name}</h4>
                                <p className="text-xs text-slate-600 dark:text-slate-300">{projects[message.project].description}</p>
                              </div>
                              <span className="text-xs font-semibold text-amber-600 dark:text-amber-200">Подробнее →</span>
                            </div>
                            <div className="flex flex-wrap items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                              <FileText className="h-3 w-3" />
                              <span>{projects[message.project].fileName}</span>
                              <span>• {projects[message.project].fileSize}</span>
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
                  <div className="rounded-3xl bg-slate-100 px-6 py-3 ring-1 ring-slate-200 dark:bg-slate-800 dark:ring-slate-700">
                    <div className="flex gap-2">
                      <div className="h-2 w-2 animate-bounce rounded-full bg-slate-400" />
                      <div className="h-2 w-2 animate-bounce rounded-full bg-slate-400" />
                      <div className="h-2 w-2 animate-bounce rounded-full bg-slate-400" />
                    </div>
                  </div>
                </div>
              )}

              <div ref={messagesEndRef} />
            </div>

            <div className="mt-6">
              <form onSubmit={handleSendMessage} className="relative mx-auto max-w-3xl">
                <div className="flex items-center gap-3 rounded-full border border-slate-300 bg-white px-4 py-2 shadow-inner focus-within:border-amber-400 dark:border-slate-700 dark:bg-slate-900">
                  <input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    placeholder="Опишите задачу или задайте вопрос AI..."
                    className="w-full bg-transparent px-2 py-2 text-sm text-slate-900 outline-none placeholder:text-slate-400 dark:text-slate-50 dark:placeholder:text-slate-500"
                  />
                  <button
                    type="submit"
                    disabled={!input.trim()}
                    className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-amber-400 text-white transition hover:bg-amber-500 disabled:cursor-not-allowed disabled:bg-slate-300 dark:bg-amber-500 dark:hover:bg-amber-400"
                    aria-label="Отправить"
                  >
                    <Send className="h-4 w-4" />
                  </button>
                </div>
                <div className="mt-2 flex items-center justify-between px-2 text-xs text-slate-500 dark:text-slate-400">
                  <span>🔒 Ваши данные защищены</span>
                  <span>Нажмите Enter для отправки</span>
                </div>
              </form>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
