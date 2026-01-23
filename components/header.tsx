'use client';

import { Bell, Moon } from 'lucide-react';
import { useRouter, usePathname } from 'next/navigation';

export default function Header() {
  const router = useRouter();
  const pathname = usePathname();
  const isLifecyclePage = pathname === '/lifecycle';

  return (
    <header className="mx-auto mb-12 w-fit rounded-full border border-gray-200 bg-white px-8 py-3 shadow-sm">
      <div className="flex items-center justify-between gap-12">
        {/* Logo */}
        <div className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-yellow-400 font-bold text-white text-xs">
            THE
          </div>
          <span className="text-sm font-semibold text-gray-800">QURYLS</span>
        </div>

        {/* Navigation */}
        <nav className="flex gap-8">
          <a href="#" className="text-sm text-gray-600 hover:text-gray-900">
            Календар
          </a>
          <a href="#" className="text-sm text-gray-600 hover:text-gray-900">
            Иерархия
          </a>
          <button
            onClick={() => router.push('/')}
            className={`text-sm font-bold ${!isLifecyclePage ? 'border-b-2 border-black text-black' : 'text-gray-600 hover:text-gray-900'}`}
          >
            Дашборд
          </button>
          <button 
            onClick={() => router.push('/lifecycle')}
            className={`text-sm font-bold ${isLifecyclePage ? 'border-b-2 border-black text-black underline' : 'text-gray-600 hover:text-gray-900'}`}
          >
            ЖЦП
          </button>
        </nav>

        {/* Right Icons */}
        <div className="flex items-center gap-4">
          <button className="text-gray-600 hover:text-gray-900">
            <Bell size={18} />
          </button>
          <button className="text-gray-600 hover:text-gray-900">
            <Moon size={18} />
          </button>
          <button className="h-8 w-8 overflow-hidden rounded-full bg-gray-300">
            <img
              src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
              alt="User avatar"
              className="h-full w-full object-cover"
            />
          </button>
        </div>
      </div>
    </header>
  );
}
