'use client';

import { useRef, useState } from 'react';
import { Plus, Lightbulb, Upload } from 'lucide-react';
import { getAccessToken } from '@/lib/api';

const DEFAULT_TASK_ID = Number(process.env.NEXT_PUBLIC_CHAT_TASK_ID ?? 1) || 1;

export default function Sidebar() {
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadMsg, setUploadMsg] = useState<string | null>(null);
  const [uploadError, setUploadError] = useState<string | null>(null);

  const handleSelect = () => {
    setUploadMsg(null);
    setUploadError(null);
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    setUploadMsg(null);
    setUploadError(null);

    try {
      const token = getAccessToken();
      const formData = new FormData();
      formData.append('file', file);

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/tasks/${DEFAULT_TASK_ID}/files`, {
        method: 'POST',
        headers: token ? { Authorization: `Bearer ${token}` } : undefined,
        body: formData,
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || 'Не удалось загрузить файл');
      }

      setUploadMsg(`Файл ${file.name} загружен в задачу #${DEFAULT_TASK_ID}`);
    } catch (err: any) {
      setUploadError(err?.message || 'Ошибка загрузки');
    } finally {
      setUploading(false);
      if (fileInputRef.current) fileInputRef.current.value = '';
    }
  };

  return (
    <aside className="w-64 border-r border-gray-200 bg-white p-6 dark:border-slate-800 dark:bg-slate-900">
      {/* Context Section */}
      <div>
        <div className="flex items-center gap-2 mb-1">
          <h3 className="text-xs font-semibold uppercase tracking-wider text-gray-900 dark:text-slate-100">
            КОНТЕКСТ
          </h3>
          <span className="text-xs text-gray-500 dark:text-slate-400">0</span>
        </div>

        {/* Empty State */}
        <div className="mt-6 rounded-lg bg-gray-50 p-6 text-center dark:bg-slate-800">
          <div className="mb-3 flex justify-center">
            <div className="rounded bg-gray-200 p-2 dark:bg-slate-700">
              <svg
                className="h-5 w-5 text-gray-400 dark:text-slate-300"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              </svg>
            </div>
          </div>
          <p className="text-xs text-gray-600 dark:text-slate-300">
            Добавьте проекты или файлы, чтобы AI лучше понимал задачу
          </p>
        </div>

        {/* Add Context Button */}
        <button
          onClick={handleSelect}
          disabled={uploading}
          className="mt-4 flex w-full items-center justify-center gap-2 rounded-lg border border-gray-300 bg-white py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-70 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
        >
          {uploading ? <Upload size={16} className="animate-pulse" /> : <Plus size={16} />} 
          {uploading ? 'Загрузка...' : 'Добавить контекст'}
        </button>
        <input
          ref={fileInputRef}
          type="file"
          accept=".pdf,.doc,.docx,.ppt,.pptx"
          className="hidden"
          onChange={handleFileChange}
        />
        {uploadMsg && <p className="mt-2 text-xs text-green-600 dark:text-green-400">{uploadMsg}</p>}
        {uploadError && <p className="mt-2 text-xs text-red-600 dark:text-red-400">{uploadError}</p>}
      </div>

      {/* Tip Section */}
      <div className="mt-8 border-t border-gray-200 pt-6 dark:border-slate-800">
        <div className="flex gap-3">
          <div className="mt-0.5 flex-shrink-0">
            <Lightbulb size={18} className="text-blue-500 dark:text-blue-300" />
          </div>
          <div>
            <h4 className="text-xs font-semibold text-gray-900 dark:text-slate-100">Подсказка</h4>
            <p className="mt-1 text-xs text-gray-600 dark:text-slate-300">
              Вы можете выбрать конкретный проект в качестве контекста, чтобы ответы были точнее.
            </p>
          </div>
        </div>
      </div>
    </aside>
  );
}
