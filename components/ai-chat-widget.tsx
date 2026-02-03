'use client';

import { usePathname, useRouter } from 'next/navigation';
import { useState, useRef, useEffect } from 'react';
import { Maximize2, MessageCircle, Send, X } from 'lucide-react';
import { ScrollArea } from '@/components/ui/scroll-area';

const CHATBOT_API_URL = typeof process !== 'undefined' ? process.env.NEXT_PUBLIC_CHATBOT_API_URL || '' : '';

type ChatMessage = { type: 'ai' | 'user'; text: string };

const INITIAL_MESSAGE: ChatMessage = {
  type: 'ai',
  text: CHATBOT_API_URL
    ? 'Chatbot connected. Try /project list or "create project My App".'
    : 'Привет! Я твой AI помощник. Спроси о любом проекте или напиши "Расскажи о проекте Shyraq".',
};

export default function AIChatWidget() {
  const pathname = usePathname();
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([INITIAL_MESSAGE]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [chatSession, setChatSession] = useState({ active_project_id: '', active_project_title: '' });
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    scrollRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim()) return;

    const userText = input.trim();
    setMessages((prev) => [...prev, { type: 'user', text: userText }]);
    setInput('');
    setIsLoading(true);

    if (CHATBOT_API_URL) {
      try {
        const res = await fetch(`${CHATBOT_API_URL.replace(/\/$/, '')}/api/chat`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            message: userText,
            active_project_id: chatSession.active_project_id || '',
            active_project_title: chatSession.active_project_title || '',
          }),
        });
        if (!res.ok) throw new Error(await res.text() || res.statusText);
        const data = (await res.json()) as {
          reply: string;
          exit?: boolean;
          active_project_id?: string;
          active_project_title?: string;
        };
        setChatSession({
          active_project_id: data.active_project_id || '',
          active_project_title: data.active_project_title || '',
        });
        if (data.reply) setMessages((prev) => [...prev, { type: 'ai', text: data.reply }]);
        if (data.exit) setMessages((prev) => [...prev, { type: 'ai', text: 'Goodbye!' }]);
      } catch (err) {
        setMessages((prev) => [
          ...prev,
          { type: 'ai', text: 'Error: ' + (err instanceof Error ? err.message : String(err)) },
        ]);
      }
      setIsLoading(false);
      return;
    }

    // Demo fallback when no chatbot API
    const userMessage = userText.toLowerCase();
    setTimeout(() => {
      let reply = 'Я помогу тебе! Напиши название проекта, например "Shyraq", "Ansau" или "Dariya".';
      if (userMessage.includes('shyraq')) reply = 'Вот информация о проекте Shyraq.';
      else if (userMessage.includes('ansau')) reply = 'Вот информация о проекте Ansau.';
      else if (userMessage.includes('dariya')) reply = 'Вот информация о проекте Dariya.';
      setMessages((prev) => [...prev, { type: 'ai', text: reply }]);
      setIsLoading(false);
    }, 500);
  };

  if (pathname.startsWith('/chat')) return null;

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end gap-3">
      {!open && (
        <button
          onClick={() => setOpen(true)}
          className="flex h-14 w-14 items-center justify-center rounded-full bg-slate-800 text-white shadow-lg ring-2 ring-white/10 transition hover:-translate-y-1 hover:bg-slate-900"
          aria-label="Открыть AI чат"
        >
          <MessageCircle className="h-6 w-6" />
        </button>
      )}

      {open && (
        <div className="flex h-[420px] w-[360px] flex-col overflow-hidden rounded-3xl bg-slate-800 text-white shadow-2xl ring-1 ring-black/10">
          {/* Header */}
          <div className="flex items-center justify-between rounded-t-3xl bg-slate-700 px-4 py-3">
            <div className="flex items-center gap-2">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-amber-300 text-xs font-black text-slate-900 shadow-sm">
                THE
              </div>
              <span className="text-lg font-bold tracking-tight">QURYLYS AI</span>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => router.push('/chat')}
                className="rounded-full p-2 text-white/80 transition hover:bg-white/10 hover:text-white"
                aria-label="Открыть полноэкранно"
              >
                <Maximize2 className="h-5 w-5" />
              </button>
              <button
                onClick={() => setOpen(false)}
                className="rounded-full p-2 text-white/80 transition hover:bg-white/10 hover:text-white"
                aria-label="Закрыть"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
          </div>

          {/* Messages */}
          <ScrollArea className="flex-1 px-4 py-3">
            <div className="space-y-3 pb-2">
              {messages.map((msg, idx) => (
                <div
                  key={idx}
                  className={`flex ${msg.type === 'user' ? 'justify-end' : 'justify-start'}`}
                >
                  <div
                    className={`max-w-[85%] rounded-2xl px-4 py-2.5 text-sm ${
                      msg.type === 'user'
                        ? 'bg-amber-400 text-slate-900'
                        : 'bg-slate-700/70 text-white'
                    }`}
                  >
                    {msg.text}
                  </div>
                </div>
              ))}
              {isLoading && (
                <div className="flex justify-start">
                  <div className="flex gap-1 rounded-2xl bg-slate-700/70 px-4 py-2.5">
                    <span className="h-2 w-2 animate-bounce rounded-full bg-white/80" />
                    <span className="h-2 w-2 animate-bounce rounded-full bg-white/80 [animation-delay:0.1s]" />
                    <span className="h-2 w-2 animate-bounce rounded-full bg-white/80 [animation-delay:0.2s]" />
                  </div>
                </div>
              )}
              <div ref={scrollRef} />
            </div>
          </ScrollArea>

          {/* Input */}
          <form
            onSubmit={handleSend}
            className="flex items-center gap-2 border-t border-white/10 bg-slate-800 px-4 py-3"
          >
            <input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Напишите вопрос..."
              className="flex-1 rounded-full bg-slate-700/70 px-4 py-2.5 text-sm text-white placeholder:text-slate-400 outline-none ring-1 ring-white/10 focus:ring-amber-400/50"
            />
            <button
              type="submit"
              disabled={!input.trim() || isLoading}
              className="flex h-10 w-10 items-center justify-center rounded-full bg-amber-400 text-slate-900 shadow-sm transition hover:bg-amber-300 disabled:opacity-50 disabled:hover:bg-amber-400"
              aria-label="Отправить"
            >
              <Send className="h-5 w-5" />
            </button>
          </form>
        </div>
      )}
    </div>
  );
}
