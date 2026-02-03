"use client";

import { useRef, useState } from "react";
import { Lightbulb, Paperclip, Plus } from "lucide-react";

export default function Sidebar() {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [files, setFiles] = useState<File[]>([]);

  const handlePick = () => inputRef.current?.click();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFiles(Array.from(e.target.files));
    }
  };

  return (
    <aside className="w-64 border-r border-gray-200 bg-white p-6 dark:border-slate-800 dark:bg-slate-900">
      <div>
        <div className="mb-1 flex items-center gap-2">
          <h3 className="text-xs font-semibold uppercase tracking-wider text-gray-900 dark:text-slate-100">КОНТЕКСТ</h3>
          <span className="text-xs text-gray-500 dark:text-slate-400">{files.length}</span>
        </div>

        <div className="mt-6 rounded-lg bg-gray-50 p-6 text-center dark:bg-slate-800">
          <div className="mb-3 flex justify-center">
            <div className="rounded bg-gray-200 p-2 dark:bg-slate-700">
              <Paperclip className="h-5 w-5 text-gray-500 dark:text-slate-200" />
            </div>
          </div>
          <p className="text-xs text-gray-600 dark:text-slate-300">
            Добавьте проекты или файлы, чтобы AI лучше понимал задачу
          </p>
          {files.length > 0 && (
            <div className="mt-4 space-y-1 text-left text-xs text-gray-700 dark:text-slate-200">
              {files.map((file) => (
                <div key={file.name} className="truncate rounded bg-white px-2 py-1 ring-1 ring-gray-200 dark:bg-slate-900 dark:ring-slate-700">
                  {file.name}
                </div>
              ))}
            </div>
          )}
        </div>

        <button
          onClick={handlePick}
          className="mt-4 flex w-full items-center justify-center gap-2 rounded-lg border border-gray-300 bg-white py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
        >
          <Plus size={16} />
          Добавить контекст
        </button>
        <input
          ref={inputRef}
          type="file"
          multiple
          onChange={handleChange}
          accept=".pdf,.doc,.docx,.txt,.xlsx,.xls"
          className="hidden"
        />
      </div>

      <div className="mt-8 border-t border-gray-200 pt-6 dark:border-slate-800">
        <div className="flex gap-3">
          <div className="mt-0.5 flex-shrink-0">
            <Lightbulb size={18} className="text-blue-500" />
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
