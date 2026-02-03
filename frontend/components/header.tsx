'use client';

import { useState } from 'react';
import { MessageCircle, Moon } from 'lucide-react';
import { useRouter, usePathname } from 'next/navigation';
import ChatModal from './chat-modal';

export default function Header() {
  const router = useRouter();
  const pathname = usePathname();
  const isLifecyclePage = pathname === '/lifecycle';
  const isDashboard = pathname === '/dashboard' || pathname === '/';
  const isHierarchy = pathname === '/hierarchy';
  const [isChatModalOpen, setIsChatModalOpen] = useState(false);

  return (
    <>
      <header className="inline-flex items-center gap-8 rounded-full border border-gray-200 bg-white px-6 py-3 shadow-sm">
        {/* Logo */}
        <div className="flex items-center gap-1">
          <div className="flex h-7 w-7 items-center justify-center rounded-full bg-amber-400 font-bold text-white text-[10px]">
            THE
          </div>
          <span className="text-sm font-semibold text-amber-500">QURYLYS</span>
        </div>

        {/* Navigation */}
        <nav className="flex items-center gap-6">
          <a href="#" className="text-sm text-gray-600 hover:text-gray-900">
            Календар
          </a>
          <button
            onClick={() => router.push('/hierarchy')}
            className={`text-sm ${isHierarchy ? 'font-semibold text-gray-900 underline underline-offset-4' : 'text-gray-600 hover:text-gray-900'}`}
          >
            Иерархия
          </button>
          <button
            onClick={() => router.push('/dashboard')}
            className={`text-sm ${isDashboard ? 'font-semibold text-gray-900 underline underline-offset-4' : 'text-gray-600 hover:text-gray-900'}`}
          >
            Дашборд
          </button>
          <button 
            onClick={() => router.push('/lifecycle')}
            className={`text-sm ${isLifecyclePage ? 'font-semibold text-gray-900 underline underline-offset-4' : 'text-gray-600 hover:text-gray-900'}`}
          >
            ЖЦП
          </button>
        </nav>

        {/* Right Icons */}
        <div className="flex items-center gap-3">
          <button 
            onClick={() => setIsChatModalOpen(true)}
            className="relative text-red-400 hover:text-red-500"
          >
            <MessageCircle size={22} fill="currentColor" />
          </button>
          <button className="flex h-8 w-8 items-center justify-center rounded-full bg-teal-600 text-white">
            <Moon size={16} />
          </button>
          <button className="h-8 w-8 overflow-hidden rounded-full ring-2 ring-gray-200">
            <img
              src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
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
