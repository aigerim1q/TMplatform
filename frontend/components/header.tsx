'use client';

import { useEffect, useState } from 'react';
import { MessageCircle, Moon, Bell, Sun } from 'lucide-react';
import { useRouter, usePathname } from 'next/navigation';
import { useTheme } from 'next-themes';
import { fetchJson } from '@/lib/api';
import ChatModal from './chat-modal';

export default function Header() {
  const router = useRouter();
  const pathname = usePathname();
  const { theme, setTheme } = useTheme();
  const isCalendarPage = pathname === '/calendar';
  const isLifecyclePage = pathname === '/lifecycle';
  const isDashboard = pathname === '/dashboard' || pathname === '/';
  const isHierarchy = pathname === '/hierarchy';
  const [avatarUrl, setAvatarUrl] = useState<string | null>(null);
  const [isChatModalOpen, setIsChatModalOpen] = useState(false);

  useEffect(() => {
    let alive = true;

    const loadMe = async () => {
      try {
        const me = await fetchJson<{ avatar_url?: string | null }>("/me");
        if (!alive) return;
        setAvatarUrl(me.avatar_url || null);
      } catch (err) {
        // Silent fallback to default avatar
      }
    };

    loadMe();

    const handleMeUpdated = (event: Event) => {
      const detail = (event as CustomEvent).detail as { avatar_url?: string | null } | undefined;
      if (detail) {
        setAvatarUrl(detail.avatar_url ?? null);
      }
    };

    window.addEventListener('me-updated', handleMeUpdated);

    return () => {
      alive = false;
      window.removeEventListener('me-updated', handleMeUpdated);
    };
  }, []);

  return (
    <>
      <header className="inline-flex items-center gap-8 rounded-full border border-gray-200 bg-white px-6 py-3 shadow-sm dark:border-slate-800 dark:bg-slate-900/80 dark:shadow-[0_8px_24px_rgba(0,0,0,0.35)] backdrop-blur">
        {/* Logo */}
        <button
          onClick={() => router.push('/dashboard')}
          className="flex items-center gap-1 focus:outline-none"
          aria-label="На дашборд"
        >
          <div className="flex h-7 w-7 items-center justify-center rounded-full bg-amber-400 font-bold text-white text-[10px] shadow dark:shadow-amber-500/30">
            THE
          </div>
          <span className="text-sm font-semibold text-amber-500 dark:text-amber-300">QURYLYS</span>
        </button>

        {/* Navigation */}
        <nav className="flex items-center gap-6">
          <button
            onClick={() => router.push('/calendar')}
            className={`text-sm ${isCalendarPage ? 'font-semibold text-gray-900 underline underline-offset-4 dark:text-slate-100' : 'text-gray-600 hover:text-gray-900 dark:text-slate-400 dark:hover:text-slate-100'}`}
          >
            Календар
          </button>
          <button
            onClick={() => router.push('/hierarchy')}
            className={`text-sm ${isHierarchy ? 'font-semibold text-gray-900 underline underline-offset-4 dark:text-slate-100' : 'text-gray-600 hover:text-gray-900 dark:text-slate-400 dark:hover:text-slate-100'}`}
          >
            Иерархия
          </button>
          <button
            onClick={() => router.push('/dashboard')}
            className={`text-sm ${isDashboard ? 'font-semibold text-gray-900 underline underline-offset-4 dark:text-slate-100' : 'text-gray-600 hover:text-gray-900 dark:text-slate-400 dark:hover:text-slate-100'}`}
          >
            Дашборд
          </button>
          <button 
            onClick={() => router.push('/lifecycle')}
            className={`text-sm ${isLifecyclePage ? 'font-semibold text-gray-900 underline underline-offset-4 dark:text-slate-100' : 'text-gray-600 hover:text-gray-900 dark:text-slate-400 dark:hover:text-slate-100'}`}
          >
            ЖЦП
          </button>
        </nav>

        {/* Right Icons */}
        <div className="flex items-center gap-3">
          <button
            className="flex h-8 w-8 items-center justify-center rounded-full bg-slate-100 text-slate-600 hover:text-slate-800 dark:bg-slate-800 dark:text-slate-200"
            aria-label="Уведомления"
          >
            <Bell size={18} />
          </button>
          <button 
            onClick={() => setIsChatModalOpen(true)}
            className="relative text-red-400 hover:text-red-500 dark:text-red-300"
          >
            <MessageCircle size={22} fill="currentColor" />
          </button>
          <button
            onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
            className="flex h-8 w-8 items-center justify-center rounded-full bg-teal-600 text-white dark:bg-slate-800 dark:text-amber-200"
            aria-label="Переключить тему"
            type="button"
          >
            {theme === 'dark' ? <Sun size={16} /> : <Moon size={16} />}
          </button>
          <button
            onClick={() => router.push('/profile')}
            className="h-8 w-8 overflow-hidden rounded-full ring-2 ring-gray-200 hover:ring-gray-300 dark:ring-slate-700 dark:hover:ring-slate-500"
            aria-label="Профиль"
            type="button"
          >
            <img
              src={avatarUrl || "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"}
              alt="User avatar"
              className="h-full w-full object-cover"
            />
          </button>
        </div>
      </header>

      {/* Chat Modal */}
      <ChatModal isOpen={isChatModalOpen} onClose={() => setIsChatModalOpen(false)} />
    </>
  );
}
