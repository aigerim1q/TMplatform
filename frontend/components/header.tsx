'use client';

import Image from 'next/image';
import { Bell, Moon, Sun } from 'lucide-react';
import { usePathname, useRouter } from 'next/navigation';
import { useTheme } from 'next-themes';
import { useEffect, useState } from 'react';

export default function Header() {
  const router = useRouter();
  const pathname = usePathname();
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => setMounted(true), []);

  const nav = [
    { label: 'Календарь', href: '/calendar' },
    { label: 'Иерархия', href: '/hierarchy' },
    { label: 'Дашборд', href: '/dashboard' },
    { label: 'ЖЦП', href: '/lifecycle' },
  ];

  const isActive = (href: string) => href === '/' ? pathname === '/' : pathname.startsWith(href);

  const toggleTheme = () => setTheme(theme === 'dark' ? 'light' : 'dark');

  return (
    <div className="flex w-full justify-center bg-transparent py-6">
      <header className="flex w-[92%] max-w-6xl items-center justify-between rounded-full border border-slate-200 bg-white px-8 py-4 shadow-sm backdrop-blur-sm dark:border-slate-700 dark:bg-slate-900">
        {/* Logo */}
        <button
          onClick={() => router.push("/dashboard")}
          className="flex items-center gap-3 rounded-full px-2 py-1 text-left transition hover:opacity-80"
          aria-label="На дашборд"
        >
          <div className="flex h-11 w-11 items-center justify-center rounded-full bg-[#FEC42E] text-sm font-black text-white shadow-sm">
            THE
          </div>
          <span className="text-xl font-semibold text-slate-900 dark:text-white">QURYLS</span>
        </button>

        {/* Navigation */}
        <nav className="flex items-center gap-8 text-sm">
          {nav.map((item) => (
            <button
              key={item.href}
              onClick={() => router.push(item.href)}
              className={`transition-colors ${isActive(item.href)
                ? 'font-extrabold text-slate-900 underline decoration-2 underline-offset-[8px] dark:text-white'
                : 'font-medium text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white'}`}
            >
              {item.label}
            </button>
          ))}
        </nav>

        {/* Right Icons */}
        <div className="flex items-center gap-4">
          <button className="text-slate-500 hover:text-slate-900 transition-colors dark:text-slate-300 dark:hover:text-white" aria-label="Уведомления">
            <Bell size={20} />
          </button>

          <button
            onClick={toggleTheme}
            className="text-slate-500 hover:text-slate-900 transition-colors dark:text-slate-300 dark:hover:text-white"
            aria-label="Переключить тему"
          >
            {mounted && theme === 'dark' ? <Sun size={20} /> : <Moon size={20} />}
          </button>

          <button
            onClick={() => router.push("/profile")}
            className="relative h-10 w-10 overflow-hidden rounded-full border border-slate-200 bg-gray-100 transition hover:ring-2 hover:ring-amber-300 dark:border-slate-700"
            aria-label="Профиль"
          >
            <Image
              src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=128&h=128&q=80"
              alt="User avatar"
              fill
              unoptimized
              sizes="40px"
              className="object-cover"
            />
          </button>
        </div>
      </header>
    </div>
  );
}
