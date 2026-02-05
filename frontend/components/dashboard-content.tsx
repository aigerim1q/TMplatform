'use client';

import React, { useCallback, useEffect, useState } from "react"
import { useRouter } from 'next/navigation';
import { Clock, Plus, X, Users, UserPlus, ChevronDown, ChevronRight, Upload, Trash } from 'lucide-react';
import { fetchJson } from '@/lib/api';

type ProjectApi = {
  id: number;
  name: string;
  status: string;
  image_url?: string | null;
  stages_count?: number;
  done_stages_count?: number;
  progress_percent?: number;
  end_date?: string | null;
};

type TaskApi = {
  id: number;
  stage_id: number;
  stage_title: string;
  project_id: number;
  project_name: string;
  project_image_url?: string | null;
  title: string;
  description?: string | null;
  status: string;
  priority: number;
  due_date?: string | null;
  assignee_id?: number | null;
  assignee_name?: string | null;
  author_name?: string | null;
};

interface TaskCardProps {
  id: string;
  project: string;
  time: string;
  timeStatus: 'danger' | 'warning' | 'success';
  title: string;
  description: string;
  responsible?: string;
  onClick?: () => void;
  onResponsibleClick?: (e: React.MouseEvent) => void;
}

function TaskCard({ id, project, time, timeStatus, title, description, responsible, onClick, onResponsibleClick }: TaskCardProps) {
  const statusColors = {
    danger: 'bg-red-500',
    warning: 'bg-amber-500',
    success: 'bg-green-500',
  };

  return (
    <div 
      onClick={onClick}
      className="rounded-2xl bg-white p-5 shadow-sm cursor-pointer hover:shadow-md transition-shadow dark:bg-slate-900/70 dark:border dark:border-slate-800 dark:hover:border-slate-700 dark:shadow-[0_8px_24px_rgba(0,0,0,0.35)]"
    >
      <div className="flex items-center justify-between mb-3">
        <span className="text-sm text-gray-500 dark:text-slate-300">Проект: {project}</span>
        <div className="flex items-center gap-1.5 rounded-full bg-amber-50 px-3 py-1 dark:bg-amber-400/10">
          <Clock size={12} className="text-gray-600 dark:text-amber-200" />
          <span className="text-sm text-gray-700 dark:text-slate-100">{time}</span>
          <span className={`h-2.5 w-2.5 rounded-full ${statusColors[timeStatus]}`} />
        </div>
      </div>
      <h3 className="font-bold text-base text-gray-900 mb-2 line-clamp-1 dark:text-white">{title}</h3>
      <p className="text-sm text-gray-600 line-clamp-3 dark:text-slate-200">{description}</p>
    </div>
  );
}

function getDueMeta(dueDate?: string | null) {
  if (!dueDate) return { label: 'Без срока', status: 'success' as const };
  const due = new Date(dueDate);
  if (Number.isNaN(due.getTime())) return { label: 'Без срока', status: 'success' as const };
  const now = new Date();
  const diffMs = due.getTime() - now.getTime();
  const diffDays = Math.round(diffMs / (1000 * 60 * 60 * 24));
  if (diffMs < 0) return { label: `Просрочено ${Math.abs(diffDays)} дн`, status: 'danger' as const };
  if (diffDays <= 2) return { label: `До срока ${diffDays} дн`, status: 'warning' as const };
  return { label: `До срока ${diffDays} дн`, status: 'success' as const };
}

interface ProjectCardProps {
  id: string;
  name: string;
  days: string;
  image: string;
  onClick?: () => void;
  onDelete?: () => void;
}

function ProjectCard({ id, name, days, image, onClick, onDelete }: ProjectCardProps) {
  return (
    <div 
      onClick={onClick}
      className="relative rounded-2xl overflow-hidden aspect-[16/10] cursor-pointer hover:shadow-lg transition-shadow"
    >
      <img
        src={image || "/placeholder.svg"}
        alt={name}
        className="w-full h-full object-cover"
        onError={(e) => {
          // Gracefully fall back if backend image path is missing
          e.currentTarget.onerror = null;
          e.currentTarget.src = '/placeholder.svg';
        }}
      />
      <div className="absolute inset-0 bg-gradient-to-b from-black/55 via-black/10 to-black/35" />

      <div className="absolute top-3 left-3 right-3 flex items-start gap-2">
        <div className="flex-1 min-w-0">
          <div className="inline-flex items-center gap-2 rounded-full bg-white/90 px-3 py-1 text-xs font-semibold text-gray-800 shadow-sm">
            <span className="truncate">Проект: {name}</span>
          </div>
        </div>
        {onDelete && (
          <button
            onClick={(e) => {
              e.stopPropagation();
              onDelete();
            }}
            className="flex items-center justify-center w-9 h-9 bg-white/95 rounded-full shadow-md hover:bg-white transition"
            aria-label="Удалить проект"
          >
            <Trash size={16} className="text-red-500" />
          </button>
        )}
      </div>

      <div className="absolute top-14 left-3 right-3 flex items-center justify-between gap-2">
        <div className="flex items-center gap-2 rounded-full bg-white/92 px-3 py-1 text-xs text-gray-800 shadow-sm">
          <Clock size={12} className="text-gray-600" />
          <span className="font-medium whitespace-nowrap">{days}</span>
          <span className="h-2.5 w-2.5 rounded-full bg-green-500" />
        </div>
        <button className="flex items-center justify-center w-9 h-9 bg-white/95 rounded-full shadow-md hover:bg-white transition">
          <ChevronRight size={18} className="text-gray-700" />
        </button>
      </div>
    </div>
  );
}

interface SectionHeaderProps {
  color: 'green' | 'red' | 'yellow';
  title: string;
  count: number;
}

function SectionHeader({ color, title, count }: SectionHeaderProps) {
  const dotColors = {
    green: 'bg-green-500',
    red: 'bg-red-500',
    yellow: 'bg-amber-400',
  };

  return (
    <div className="inline-flex items-center gap-2 rounded-full bg-amber-100/80 px-4 py-2 dark:bg-slate-900/60 dark:border dark:border-slate-800">
      <span className={`h-3 w-3 rounded-full ${dotColors[color]}`} />
      <span className="font-medium text-gray-800 text-sm dark:text-slate-100">{title}: {count}</span>
    </div>
  );
}

function AddButton({ text, onClick }: { text: string; onClick?: () => void }) {
  return (
    <button 
      onClick={onClick}
      className="flex items-center gap-2 text-gray-700 hover:text-gray-900 font-medium text-sm dark:text-slate-200 dark:hover:text-white"
    >
      <Plus size={16} />
      {text}
    </button>
  );
}

interface AddProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreated: () => Promise<void>;
  projects: ProjectApi[];
}

function AddProjectModal({ isOpen, onClose, onCreated, projects }: AddProjectModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [endDate, setEndDate] = useState('');
  const [imageDataUrl, setImageDataUrl] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = () => {
      if (typeof reader.result === 'string') {
        setImageDataUrl(reader.result);
      }
    };
    reader.readAsDataURL(file);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setSaving(true);
      setError(null);
      await fetchJson('/projects', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: name.trim(),
          description: description.trim() ? description : undefined,
          status: 'draft',
          end_date: endDate ? new Date(endDate).toISOString() : undefined,
          image_url: imageDataUrl || undefined,
        }),
      });
      await onCreated();
      onClose();
      setName('');
      setDescription('');
      setEndDate('');
      setImageDataUrl(null);
    } catch (err: any) {
      setError(err?.message || 'Не удалось создать проект');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="w-full max-w-lg rounded-2xl bg-white p-6 shadow-xl dark:bg-slate-900 dark:border dark:border-slate-800">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-gray-900 dark:text-white">Создать проект</h2>
          <button onClick={onClose} className="p-1 rounded-full hover:bg-gray-100 dark:hover:bg-slate-800">
            <X size={20} className="text-gray-500 dark:text-slate-300" />
          </button>
        </div>

        {error && (
          <div className="mb-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-900/30 dark:text-red-100">
            {error}
          </div>
        )}
        {!projects.length && (
          <div className="mb-4 rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-700 dark:border-amber-900/50 dark:bg-amber-900/30 dark:text-amber-100">
            Сначала создайте проект, чтобы добавить задачу.
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-sm font-medium text-gray-700 dark:text-slate-200">Название</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:ring-2 focus:ring-amber-500 focus:border-amber-500 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
              placeholder="Например, Shyraq"
            />
          </div>
          <div>
            <label className="text-sm font-medium text-gray-700 dark:text-slate-200">Описание</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:ring-2 focus:ring-amber-500 focus:border-amber-500 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
              placeholder="Коротко о проекте"
            />
          </div>
          <div className="grid grid-cols-2 gap-3 items-end">
            <div>
              <label className="text-sm font-medium text-gray-700 dark:text-slate-200">Дедлайн</label>
              <input
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:ring-2 focus:ring-amber-500 focus:border-amber-500 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
              />
            </div>
            <div>
              <label className="text-sm font-medium text-gray-700 dark:text-slate-200">Картинка</label>
              <label className="mt-1 inline-flex items-center gap-2 w-full justify-center rounded-xl border border-dashed border-gray-300 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 cursor-pointer dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800">
                <Upload size={16} /> Загрузить
                <input type="file" accept="image/*" className="hidden" onChange={handleFile} />
              </label>
              {imageDataUrl && <p className="mt-1 text-xs text-gray-500 dark:text-slate-400">Изображение выбрано</p>}
            </div>
          </div>

          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 rounded-full border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-slate-700 dark:text-slate-100 dark:hover:bg-slate-800"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={saving}
              className="flex-1 rounded-full bg-amber-500 px-4 py-2 text-sm font-semibold text-white hover:bg-amber-600 disabled:opacity-60"
            >
              {saving ? 'Создаём…' : 'Создать проект'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

interface AddTaskModalProps {
  isOpen: boolean;
  onClose: () => void;
  projects: ProjectApi[];
  onCreated: () => Promise<void>;
}

type StageOption = { id: number; title: string };

function AddTaskModal({ isOpen, onClose, projects, onCreated }: AddTaskModalProps) {
  const [projectId, setProjectId] = useState<number | null>(null);
  const [stages, setStages] = useState<StageOption[]>([]);
  const [stageId, setStageId] = useState<number | null>(null);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [deadline, setDeadline] = useState('');
  const [loadingStages, setLoadingStages] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen && projects.length && !projectId) {
      setProjectId(projects[0].id);
    }
    if (!isOpen) {
      setError(null);
    }
  }, [isOpen, projects, projectId]);

  useEffect(() => {
    const loadStages = async () => {
      if (!projectId || !isOpen) {
        setStages([]);
        setStageId(null);
        return;
      }
      try {
        setLoadingStages(true);
        const res = await fetchJson<{ items: StageOption[] }>(`/projects/${projectId}/stages`);
        setStages(res.items || []);
        setStageId((res.items || [])[0]?.id ?? null);
      } catch (err: any) {
        setError(err?.message || 'Не удалось загрузить стадии');
      } finally {
        setLoadingStages(false);
      }
    };
    loadStages();
  }, [projectId, isOpen]);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!projectId) {
      setError('Выберите проект');
      return;
    }
    try {
      setSaving(true);
      setError(null);
      await fetchJson('/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          project_id: projectId,
          stage_id: stageId || undefined,
          title: title.trim(),
          description: description.trim() || undefined,
          due_date: deadline || undefined,
        }),
      });
      await onCreated();
      setTitle('');
      setDescription('');
      setDeadline('');
      onClose();
    } catch (err: any) {
      setError(err?.message || 'Не удалось создать задачу');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl p-6 w-full max-w-md shadow-xl dark:bg-slate-900 dark:border dark:border-slate-800">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-white">Добавить задачу</h2>
          <button 
            onClick={onClose}
            className="p-1 rounded-full hover:bg-gray-100 dark:hover:bg-slate-800"
          >
            <X size={20} className="text-gray-500 dark:text-slate-300" />
          </button>
        </div>

        {error && (
          <div className="mb-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-900/30 dark:text-red-100">
            {error}
          </div>
        )}
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-1 gap-3">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1 dark:text-slate-200">
                Проект
              </label>
              <select
                value={projectId ?? ''}
                onChange={(e) => setProjectId(Number(e.target.value) || null)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
                required
                disabled={!projects.length}
              >
                {projects.map((p) => (
                  <option key={p.id} value={p.id}>{p.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1 dark:text-slate-200">
                Стадия (опционально)
              </label>
              <select
                value={stageId ?? ''}
                onChange={(e) => setStageId(e.target.value ? Number(e.target.value) : null)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
                disabled={loadingStages}
              >
                <option value="">Без стадии</option>
                {stages.map((s) => (
                  <option key={s.id} value={s.id}>{s.title}</option>
                ))}
              </select>
            </div>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1 dark:text-slate-200">
              Название задачи
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Введите название задачи"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
              required
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1 dark:text-slate-200">
              Описание
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Введите описание задачи"
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1 dark:text-slate-200">
              Срок выполнения
            </label>
            <input
              type="date"
              value={deadline}
              onChange={(e) => setDeadline(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
            />
          </div>
          
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-full text-gray-700 font-medium hover:bg-gray-50 bg-transparent dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={saving || !projects.length}
              className="flex-1 px-4 py-2 bg-amber-500 text-white rounded-full font-medium hover:bg-amber-600 disabled:opacity-60"
            >
              {saving ? 'Создаём…' : 'Создать'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

interface ResponsibleModalProps {
  isOpen: boolean;
  onClose: () => void;
}

function ResponsibleModal({ isOpen, onClose }: ResponsibleModalProps) {
  const [selectedResponsible, setSelectedResponsible] = useState('');

  if (!isOpen) return null;

  const responsiblePeople = [
    {
      id: 1,
      name: 'Омар Ахмет',
      role: 'Архитектор',
      avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face',
      isYou: true,
    },
    {
      id: 2,
      name: 'Зейнулла Рышман',
      role: 'Архитектор',
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop&crop=face',
      isYou: false,
    },
    {
      id: 3,
      name: 'Айдын Рахимбаев',
      role: 'Аудитор',
      avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=100&h=100&fit=crop&crop=face',
      isYou: false,
    },
  ];

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-amber-50 rounded-2xl p-6 w-full max-w-md shadow-xl dark:bg-slate-900 dark:border dark:border-slate-800">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2">
            <Users size={24} className="text-gray-700 dark:text-slate-200" />
            <h2 className="text-xl font-bold text-gray-900 dark:text-white">Ответственнные</h2>
          </div>
          <div className="flex items-center gap-3">
            <button className="flex items-center gap-1 text-gray-600 hover:text-gray-900 dark:text-slate-200 dark:hover:text-white">
              <UserPlus size={18} />
              <span className="text-sm">Делегировать</span>
            </button>
            <button 
              onClick={onClose}
              className="p-1 rounded-full hover:bg-amber-100 dark:hover:bg-slate-800"
            >
              <X size={20} className="text-gray-500 dark:text-slate-300" />
            </button>
          </div>
        </div>

        {/* People List */}
        <div className="space-y-4 mb-6">
          {responsiblePeople.map((person) => (
            <div key={person.id} className="flex items-center gap-3">
              <img
                src={person.avatar || "/placeholder.svg"}
                alt={person.name}
                className="w-14 h-14 rounded-full object-cover"
              />
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-gray-900 dark:text-slate-100">{person.name}</span>
                  {person.isYou && (
                    <span className="text-xs bg-gray-300 text-gray-700 px-2 py-0.5 rounded-full dark:bg-slate-800 dark:text-slate-200">
                      это вы
                    </span>
                  )}
                </div>
                <span className="text-sm text-amber-700 dark:text-amber-200">{person.role}</span>
              </div>
            </div>
          ))}
        </div>

        {/* Quick Delegation */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs font-medium text-gray-500 uppercase tracking-wide dark:text-slate-400">
              Быстрое делегирование
            </span>
            <span className="text-xs text-gray-500 dark:text-slate-400">
              Выбрать ответственного
            </span>
          </div>
          <div className="relative">
            <select
              value={selectedResponsible}
              onChange={(e) => setSelectedResponsible(e.target.value)}
              className="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl appearance-none focus:ring-2 focus:ring-amber-500 focus:border-amber-500 text-gray-700 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-100"
            >
              <option value="">Выбрать ответст...</option>
              <option value="omar">Омар Ахмет</option>
              <option value="zeinulla">Зейнулла Рышман</option>
              <option value="aidyn">Айдын Рахимбаев</option>
            </select>
            <ChevronDown size={20} className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none dark:text-slate-500" />
          </div>
        </div>

        {/* Add Button */}
        <button className="w-full flex items-center justify-center gap-2 bg-amber-200 hover:bg-amber-300 text-amber-900 font-semibold py-3 rounded-xl transition-colors dark:bg-amber-300/20 dark:text-amber-200 dark:hover:bg-amber-300/30">
          <UserPlus size={20} />
          Добавить ответственных
        </button>
      </div>
    </div>
  );
}

export default function DashboardContent() {
  const router = useRouter();
  const [isAddTaskModalOpen, setIsAddTaskModalOpen] = useState(false);
  const [isResponsibleModalOpen, setIsResponsibleModalOpen] = useState(false);
  const [myTasks, setMyTasks] = useState<TaskApi[]>([]);
  const [urgentTasks, setUrgentTasks] = useState<TaskApi[]>([]);
  const [subordinateTasks, setSubordinateTasks] = useState<TaskApi[]>([]);
  const [loadingTasks, setLoadingTasks] = useState(true);
  const [tasksError, setTasksError] = useState<string | null>(null);

  const [projects, setProjects] = useState<ProjectApi[]>([]);
  const [loadingProjects, setLoadingProjects] = useState(true);
  const [projectsError, setProjectsError] = useState<string | null>(null);
  const [isProjectModalOpen, setIsProjectModalOpen] = useState(false);
  const [deletingProjectId, setDeletingProjectId] = useState<number | null>(null);

  useEffect(() => {
    let alive = true;
    (async () => {
      try {
        setLoadingProjects(true);
        const res = await fetchJson<{ items: ProjectApi[] }>(`/projects?mine=1`);
        if (!alive) return;
        setProjects(res?.items || []);
        setProjectsError(null);
      } catch (err: any) {
        if (!alive) return;
        setProjectsError(err?.message || 'Не удалось загрузить проекты');
      } finally {
        if (alive) setLoadingProjects(false);
      }
    })();
    return () => {
      alive = false;
    };
  }, []);

  const refetchProjects = async () => {
    const res = await fetchJson<{ items: ProjectApi[] }>(`/projects?mine=1`);
    setProjects(res?.items || []);
  };

  const handleDeleteProject = async (projectId: number) => {
    try {
      setDeletingProjectId(projectId);
      await fetchJson(`/projects/${projectId}`, { method: 'DELETE' });
      setProjectsError(null);
      await refetchProjects();
    } catch (err: any) {
      setProjectsError(err?.message || 'Не удалось удалить проект');
    } finally {
      setDeletingProjectId(null);
    }
  };

  const loadTasks = useCallback(async () => {
    try {
      setLoadingTasks(true);
      const [mine, urgent, sub] = await Promise.all([
        fetchJson<TaskApi[]>(`/tasks?mine=1`),
        fetchJson<TaskApi[]>(`/tasks?scope=urgent&mine=1`),
        fetchJson<TaskApi[]>(`/tasks?scope=subordinate`),
      ]);
      setMyTasks(Array.isArray(mine) ? mine : []);
      setUrgentTasks(Array.isArray(urgent) ? urgent : []);
      setSubordinateTasks(Array.isArray(sub) ? sub : []);
      setTasksError(null);
    } catch (err: any) {
      setTasksError(err?.message || 'Не удалось загрузить задачи');
    } finally {
      setLoadingTasks(false);
    }
  }, []);

  useEffect(() => {
    loadTasks();
  }, [loadTasks]);

  const handleTaskClick = (task: TaskApi) => {
    router.push(`/project-overview/${task.project_id}`);
  };

  const handleProjectClick = (projectId: string) => {
    router.push(`/project-overview/${projectId}`);
  };

  const renderTaskCard = (task: TaskApi) => {
    const dueMeta = getDueMeta(task.due_date);
    return (
      <TaskCard
        key={task.id}
        id={String(task.id)}
        project={task.project_name}
        time={dueMeta.label}
        timeStatus={dueMeta.status}
        title={task.title}
        description={task.description || ''}
        responsible={task.assignee_name || task.author_name || 'Не назначен'}
        onClick={() => handleTaskClick(task)}
        onResponsibleClick={handleResponsibleClick}
      />
    );
  };

  const handleResponsibleClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsResponsibleModalOpen(true);
  };

  return (
    <main className="w-full max-w-5xl mx-auto px-6 py-8 text-gray-900 dark:text-slate-100">
      {/* Add Task Modal */}
      <AddTaskModal 
        isOpen={isAddTaskModalOpen} 
        onClose={() => setIsAddTaskModalOpen(false)}
        projects={projects}
        onCreated={async () => {
          await loadTasks();
        }}
      />

      {/* Add Project Modal */}
      <AddProjectModal
        isOpen={isProjectModalOpen}
        onClose={() => setIsProjectModalOpen(false)}
        onCreated={async () => {
          await refetchProjects();
          setIsProjectModalOpen(false);
        }}
        projects={projects}
      />

      {/* Responsible Modal */}
      <ResponsibleModal
        isOpen={isResponsibleModalOpen}
        onClose={() => setIsResponsibleModalOpen(false)}
      />

      {/* Urgent Tasks Section */}
      <div className="mb-8">
        <SectionHeader color="green" title="Срочные задачи" count={urgentTasks.length} />
        {tasksError && (
          <div className="mt-3 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-900/30 dark:text-red-100">
            {tasksError}
          </div>
        )}
        <div className="mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {loadingTasks && Array.from({ length: 3 }).map((_, idx) => (
            <div key={idx} className="h-32 rounded-2xl bg-black/5 animate-pulse dark:bg-white/5" />
          ))}
          {!loadingTasks && urgentTasks.map(renderTaskCard)}
        </div>
      </div>

      {/* My Tasks Section */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <SectionHeader color="red" title="Мои задачи" count={myTasks.length} />
          <AddButton text="Добавить задачу" onClick={() => setIsAddTaskModalOpen(true)} />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {loadingTasks && Array.from({ length: 3 }).map((_, idx) => (
            <div key={idx} className="h-32 rounded-2xl bg-black/5 animate-pulse dark:bg-white/5" />
          ))}
          {!loadingTasks && myTasks.map(renderTaskCard)}
        </div>
      </div>

      {/* Projects Section */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <SectionHeader color="yellow" title="Проекты" count={projects.length} />
          <AddButton text="Добавить проект" onClick={() => setIsProjectModalOpen(true)} />
        </div>
        {projectsError && (
          <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-900/30 dark:text-red-100">
            {projectsError}
          </div>
        )}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {loadingProjects && Array.from({ length: 3 }).map((_, idx) => (
            <div key={idx} className="h-40 rounded-2xl bg-black/5 animate-pulse dark:bg-white/5" />
          ))}
          {!loadingProjects && projects.map((project) => (
            <ProjectCard 
              key={project.id} 
              id={String(project.id)}
              name={project.name}
              days={project.end_date ? `Дедлайн: ${new Date(project.end_date).toLocaleDateString()}` : `${project.done_stages_count ?? 0}/${project.stages_count ?? 0} стадий`}
              image={project.image_url || '/placeholder.svg'}
              onClick={() => handleProjectClick(String(project.id))}
              onDelete={() => {
                if (deletingProjectId === project.id) return;
                handleDeleteProject(project.id);
              }}
            />
          ))}
        </div>
      </div>

      {/* Subordinate Tasks Section */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <SectionHeader color="red" title="Задачи подчинённых" count={subordinateTasks.length} />
          <AddButton text="Добавить задачу" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {loadingTasks && Array.from({ length: 3 }).map((_, idx) => (
            <div key={idx} className="h-32 rounded-2xl bg-black/5 animate-pulse dark:bg-white/5" />
          ))}
          {!loadingTasks && subordinateTasks.map(renderTaskCard)}
        </div>
      </div>

      {/* Subordinate Projects Section */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <SectionHeader color="yellow" title="Проекты подчиненных" count={projects.length} />
          <AddButton text="Добавить проект" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {projects.map((project) => (
            <ProjectCard 
              key={project.id} 
              id={String(project.id)}
              name={project.name}
              days={`${project.done_stages_count ?? 0}/${project.stages_count ?? 0} стадий`}
              image={project.image_url || '/placeholder.svg'}
              onClick={() => handleProjectClick(String(project.id))}
            />
          ))}
        </div>
      </div>
    </main>
  );
}
