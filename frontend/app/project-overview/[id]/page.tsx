'use client';

import { useEffect, useMemo, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { Clock, Users, AlertCircle, Plus, ChevronRight, Loader2, X } from 'lucide-react';
import Header from '@/components/header';
import { fetchJson } from '@/lib/api';

type ProjectApi = {
  id: number;
  name: string;
  description?: string | null;
  end_date?: string | null;
  image_url?: string | null;
  assignees?: Array<{ id: number; email: string; role: string; name?: string }>;
};

type StageApi = {
  id: number;
  project_id: number;
  title: string;
  description?: string | null;
  status: string;
};

type TaskApi = {
  id: number;
  title: string;
  description?: string | null;
  status: string;
  priority: number;
  assignee_name?: string | null;
  author_name?: string | null;
  due_date?: string | null;
};

type UserOption = { id: number; label: string; role: string };

function formatDate(d?: string | null) {
  if (!d) return 'Без дедлайна';
  const dt = new Date(d);
  if (Number.isNaN(dt.getTime())) return 'Без дедлайна';
  return dt.toLocaleDateString();
}

export default function ProjectOverviewPage() {
  const router = useRouter();
  const params = useParams();
  const projectId = params.id as string;

  const [project, setProject] = useState<ProjectApi | null>(null);
  const [stages, setStages] = useState<StageApi[]>([]);
  const [tasksByStage, setTasksByStage] = useState<Record<number, TaskApi[]>>({});
  const [assignees, setAssignees] = useState<UserOption[]>([]);
  const [allUsers, setAllUsers] = useState<UserOption[]>([]);
  const [loading, setLoading] = useState(true);
  const [savingAssignees, setSavingAssignees] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [assigneeModalOpen, setAssigneeModalOpen] = useState(false);
  const [selectedAssignees, setSelectedAssignees] = useState<number[]>([]);
  const [stageModalOpen, setStageModalOpen] = useState(false);
  const [newStageTitle, setNewStageTitle] = useState('');
  const [savingStage, setSavingStage] = useState(false);
  const [taskModalOpen, setTaskModalOpen] = useState(false);
  const [savingTask, setSavingTask] = useState(false);
  const [taskForm, setTaskForm] = useState({
    title: '',
    description: '',
    dueDate: '',
    stageId: 0,
    assigneeId: 0,
  });

  useEffect(() => {
    const load = async () => {
      try {
        setLoading(true);
        setError(null);
        const [p, sRaw, a, users] = await Promise.all([
          fetchJson<ProjectApi>(`/projects/${projectId}`),
          fetchJson<StageApi[] | { items: StageApi[] }>(`/projects/${projectId}/stages`),
          fetchJson<{ items: UserOption[] }>(`/projects/${projectId}/assignees`).catch(() => ({ items: [] })),
          fetchJson<{ items: any[] }>(`/users`).catch(() => ({ items: [] })),
        ]);

        const stageItems = Array.isArray(sRaw) ? sRaw : sRaw?.items || [];

        setProject(p);
        setStages(stageItems);

        const assigneeOpts = (a?.items || []).map((u) => ({
          id: (u as any).id,
          label: (u as any).name || (u as any).email,
          role: (u as any).role || '',
        }));
        setAssignees(assigneeOpts);
        setSelectedAssignees(assigneeOpts.map((u) => u.id));

        setAllUsers((users.items || []).map((u: any) => ({ id: u.id, label: u.name || u.email, role: u.role })));

        if (stageItems && stageItems.length) {
          const tasksEntries = await Promise.all(
            stageItems.map(async (st) => {
              const items = await fetchJson<TaskApi[]>(`/stages/${st.id}/tasks`).catch(() => []);
              return [st.id, Array.isArray(items) ? items : []] as const;
            })
          );
          setTasksByStage(Object.fromEntries(tasksEntries));
        } else {
          setTasksByStage({});
        }
      } catch (e: any) {
        setError(e?.message || 'Не удалось загрузить данные');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [projectId]);

  const handleSaveAssignees = async () => {
    try {
      setSavingAssignees(true);
      await fetchJson(`/projects/${projectId}/assignees`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_ids: selectedAssignees }),
      });
      setAssignees(allUsers.filter((u) => selectedAssignees.includes(u.id)));
      setAssigneeModalOpen(false);
    } catch (e: any) {
      setError(e?.message || 'Не удалось сохранить назначенных');
    } finally {
      setSavingAssignees(false);
    }
  };

  const tasksTotal = useMemo(() => Object.values(tasksByStage).reduce((acc, arr) => acc + arr.length, 0), [tasksByStage]);

  const projectImage = project?.image_url || '/placeholder.svg';

  const refreshStageTasks = async (stageId: number) => {
    const items = await fetchJson<TaskApi[]>(`/stages/${stageId}/tasks`).catch(() => []);
    setTasksByStage((prev) => ({ ...prev, [stageId]: Array.isArray(items) ? items : [] }));
  };

  const handleCreateStage = async () => {
    const title = newStageTitle.trim();
    if (!title) return;
    try {
      setSavingStage(true);
      const created = await fetchJson<StageApi>(`/projects/${projectId}/stages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title }),
      });
      setStages((prev) => [...prev, created]);
      setTasksByStage((prev) => ({ ...prev, [created.id]: [] }));
      setNewStageTitle('');
      setStageModalOpen(false);
    } catch (e: any) {
      setError(e?.message || 'Не удалось создать этап');
    } finally {
      setSavingStage(false);
    }
  };

  const openTaskModal = (stageId?: number) => {
    const defaultStage = stageId || stages[0]?.id || 0;
    setTaskForm({ title: '', description: '', dueDate: '', stageId: defaultStage, assigneeId: 0 });
    setTaskModalOpen(true);
  };

  const handleTaskSave = async () => {
    if (!taskForm.title.trim() || !taskForm.stageId) return;
    try {
      setSavingTask(true);
      const payload: any = {
        project_id: Number(projectId),
        stage_id: taskForm.stageId,
        title: taskForm.title.trim(),
        description: taskForm.description.trim() || undefined,
        priority: 0,
      };
      if (taskForm.dueDate) payload.due_date = taskForm.dueDate;
      if (taskForm.assigneeId) payload.assignee_id = taskForm.assigneeId;

      await fetchJson(`/tasks`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      await refreshStageTasks(taskForm.stageId);
      setTaskModalOpen(false);
    } catch (e: any) {
      setError(e?.message || 'Не удалось создать задачу');
    } finally {
      setSavingTask(false);
    }
  };

  return (
    <div className="min-h-screen bg-white text-gray-900 dark:bg-slate-950 dark:text-slate-50">
      <div className="flex justify-center pt-6">
        <Header />
      </div>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <div className="flex items-center gap-3 mb-8">
          <button
            onClick={() => router.push('/dashboard')}
            className="px-4 py-2 rounded-full bg-gray-100 text-gray-900 hover:bg-gray-200 transition-colors text-sm font-medium dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
          >
            ← Назад
          </button>
          <button className="bg-black text-white px-6 py-2 rounded-full text-sm font-semibold dark:bg-amber-500">
            Задача
          </button>
          <button
            onClick={() => router.push(`/project/${params.id}/reports`)}
            className="bg-gray-400 text-white px-6 py-2 rounded-full text-sm font-semibold hover:bg-gray-500 transition-colors dark:bg-slate-700 dark:hover:bg-slate-600"
          >
            Отчеты
          </button>
        </div>

        {error && (
          <div className="mb-6 rounded-2xl bg-red-50 border border-red-200 text-red-700 px-4 py-3 text-sm dark:bg-red-900/30 dark:border-red-800 dark:text-red-100">{error}</div>
        )}

        <div className="mb-8 flex flex-col md:flex-row gap-6">
          <img
            src={projectImage}
            alt="Project"
            className="w-80 h-48 rounded-3xl object-cover shadow-md"
            onError={(e) => {
              e.currentTarget.onerror = null;
              e.currentTarget.src = '/placeholder.svg';
            }}
          />
          <div className="flex-1 space-y-3">
            <h1 className="text-2xl font-bold text-gray-900 dark:text-slate-50">{project ? project.name : 'Загрузка проекта...'}</h1>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              <div className="bg-yellow-100 rounded-3xl p-4 flex items-start gap-3 dark:bg-amber-500/20">
                <Clock className="w-6 h-6 text-yellow-600 flex-shrink-0 mt-1" />
                <div>
                  <p className="text-xs text-gray-600 mb-1 dark:text-amber-100/80">Дедлайн:</p>
                  <p className="text-sm font-semibold text-gray-900 dark:text-slate-50">{formatDate(project?.end_date)}</p>
                </div>
              </div>

              <div className="bg-white border-2 border-gray-300 rounded-3xl p-4 flex items-start gap-3 dark:bg-slate-900 dark:border-slate-700">
                <Users className="w-6 h-6 text-gray-600 flex-shrink-0 mt-1" />
                <div className="flex-1">
                  <div className="flex items-center justify-between mb-1">
                    <p className="text-xs text-gray-600 dark:text-slate-300">Ответственные</p>
                    <button
                      className="text-xs font-semibold text-amber-600 hover:text-amber-700"
                      onClick={() => setAssigneeModalOpen(true)}
                    >
                      Назначить
                    </button>
                  </div>
                  <p className="text-sm font-semibold text-gray-900 line-clamp-2 dark:text-slate-50">
                    {assignees.length ? assignees.map((a) => a.label).join(', ') : 'Не назначены'}
                  </p>
                </div>
              </div>

              <div className="bg-red-100 rounded-3xl p-4 flex items-start gap-3 dark:bg-red-500/20">
                <AlertCircle className="w-6 h-6 text-red-600 flex-shrink-0 mt-1" />
                <div>
                  <p className="text-xs text-gray-600 mb-1 dark:text-red-100/80">Задач</p>
                  <p className="text-sm font-semibold text-gray-900 dark:text-slate-50">{tasksTotal}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {loading && (
          <div className="flex items-center gap-2 text-gray-600 mb-8"><Loader2 className="w-4 h-4 animate-spin" /> Загружаем...</div>
        )}

        {stages.map((stage, idx) => (
          <div key={stage.id} className="mb-12">
            <h2 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2 dark:text-slate-50">
              <div className="bg-black text-white rounded-full px-3 py-1 text-xs font-bold">{idx + 1}</div>
              <span className="flex-1">{stage.title}</span>
              <button
                onClick={() => openTaskModal(stage.id)}
                className="text-xs font-semibold text-amber-600 hover:text-amber-700"
              >
                Добавить задачу
              </button>
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-4">
              {(tasksByStage[stage.id] || []).map((task) => (
                <div key={task.id} className="bg-white border-2 border-gray-200 rounded-2xl p-4 hover:shadow-md transition-shadow dark:bg-slate-900 dark:border-slate-700">
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-xs text-gray-600 dark:text-slate-300">Проект: {project?.name}</span>
                    <div className="px-3 py-1 rounded-full text-xs font-semibold bg-gray-100 text-gray-800 dark:bg-slate-800 dark:text-slate-100">
                      {task.status}
                    </div>
                  </div>
                  <h3 className="font-bold text-gray-900 mb-2 text-sm dark:text-slate-50">{task.title}</h3>
                  <p className="text-xs text-gray-600 line-clamp-2 mb-3 dark:text-slate-300">{task.description || 'Без описания'}</p>
                  <div className="text-xs text-gray-500 flex items-center justify-between dark:text-slate-400">
                    <span>Исполнитель: {task.assignee_name || 'Не назначен'}</span>
                    <ChevronRight className="w-4 h-4 text-gray-400 dark:text-slate-500" />
                  </div>
                </div>
              ))}
              {!tasksByStage[stage.id]?.length && (
                <div className="text-sm text-gray-500 dark:text-slate-400">Нет задач на этом этапе</div>
              )}
            </div>
          </div>
        ))}

        <div className="flex justify-center">
          <button
            onClick={() => setStageModalOpen(true)}
            className="bg-yellow-600 hover:bg-yellow-700 text-white px-6 py-3 rounded-full font-semibold transition-colors flex items-center gap-2 dark:bg-amber-500 dark:hover:bg-amber-600"
          >
            <Plus className="w-5 h-5" />
            Добавить этап проекта
          </button>
        </div>
      </main>

      {assigneeModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-lg shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold text-gray-900">Назначить участников</h3>
              <button onClick={() => setAssigneeModalOpen(false)} className="p-1 rounded-full hover:bg-gray-100">
                <X className="w-5 h-5 text-gray-500" />
              </button>
            </div>
            <div className="max-h-72 overflow-auto space-y-2 mb-4">
              {allUsers.map((u) => {
                const checked = selectedAssignees.includes(u.id);
                return (
                  <label key={u.id} className="flex items-center justify-between border border-gray-200 rounded-xl px-3 py-2 hover:bg-gray-50">
                    <div>
                      <p className="text-sm font-semibold text-gray-900">{u.label}</p>
                      <p className="text-xs text-gray-500">{u.role}</p>
                    </div>
                    <input
                      type="checkbox"
                      checked={checked}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setSelectedAssignees((prev) => [...prev, u.id]);
                        } else {
                          setSelectedAssignees((prev) => prev.filter((id) => id !== u.id));
                        }
                      }}
                      className="w-4 h-4"
                    />
                  </label>
                );
              })}
              {!allUsers.length && <p className="text-sm text-gray-500">Нет пользователей</p>}
            </div>
            <div className="flex gap-3">
              <button
                onClick={() => setAssigneeModalOpen(false)}
                className="flex-1 rounded-full border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                Отмена
              </button>
              <button
                onClick={handleSaveAssignees}
                disabled={savingAssignees}
                className="flex-1 rounded-full bg-amber-500 px-4 py-2 text-sm font-semibold text-white hover:bg-amber-600 disabled:opacity-60 flex items-center justify-center gap-2"
              >
                {savingAssignees && <Loader2 className="w-4 h-4 animate-spin" />} Сохранить
              </button>
            </div>
          </div>
        </div>
      )}

      {stageModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md shadow-xl dark:bg-slate-900 dark:border dark:border-slate-700">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold text-gray-900 dark:text-slate-50">Новый этап</h3>
              <button onClick={() => setStageModalOpen(false)} className="p-1 rounded-full hover:bg-gray-100">
                <X className="w-5 h-5 text-gray-500 dark:text-slate-400" />
              </button>
            </div>
            <div className="space-y-3 mb-4">
              <label className="block text-sm font-medium text-gray-700 dark:text-slate-200">Название</label>
              <input
                value={newStageTitle}
                onChange={(e) => setNewStageTitle(e.target.value)}
                className="w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500 dark:bg-slate-900 dark:border-slate-700 dark:text-slate-100"
                placeholder="Например, Предпроектная подготовка"
              />
            </div>
            <div className="flex gap-3">
              <button
                onClick={() => setStageModalOpen(false)}
                className="flex-1 rounded-full border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-slate-700 dark:text-slate-100 dark:hover:bg-slate-800"
              >
                Отмена
              </button>
              <button
                onClick={handleCreateStage}
                disabled={savingStage || !newStageTitle.trim()}
                className="flex-1 rounded-full bg-amber-500 px-4 py-2 text-sm font-semibold text-white hover:bg-amber-600 disabled:opacity-60 flex items-center justify-center gap-2 dark:bg-amber-500 dark:hover:bg-amber-600"
              >
                {savingStage && <Loader2 className="w-4 h-4 animate-spin" />} Создать
              </button>
            </div>
          </div>
        </div>
      )}

      {taskModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-lg shadow-xl dark:bg-slate-900 dark:border dark:border-slate-700">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold text-gray-900 dark:text-slate-50">Новая задача</h3>
              <button onClick={() => setTaskModalOpen(false)} className="p-1 rounded-full hover:bg-gray-100">
                <X className="w-5 h-5 text-gray-500 dark:text-slate-400" />
              </button>
            </div>
            <div className="space-y-3 mb-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-slate-200">Название</label>
                <input
                  value={taskForm.title}
                  onChange={(e) => setTaskForm({ ...taskForm, title: e.target.value })}
                  className="w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500 dark:bg-slate-900 dark:border-slate-700 dark:text-slate-100"
                  placeholder="Что нужно сделать"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-slate-200">Описание</label>
                <textarea
                  value={taskForm.description}
                  onChange={(e) => setTaskForm({ ...taskForm, description: e.target.value })}
                  className="w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500 dark:bg-slate-900 dark:border-slate-700 dark:text-slate-100"
                  rows={3}
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-slate-200">Этап</label>
                  <select
                    value={taskForm.stageId}
                    onChange={(e) => setTaskForm({ ...taskForm, stageId: Number(e.target.value) })}
                    className="w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500 dark:bg-slate-900 dark:border-slate-700 dark:text-slate-100"
                  >
                    {stages.map((s) => (
                      <option key={s.id} value={s.id}>
                        {s.title}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-slate-200">Дедлайн</label>
                  <input
                    type="date"
                    value={taskForm.dueDate}
                    onChange={(e) => setTaskForm({ ...taskForm, dueDate: e.target.value })}
                    className="w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500 dark:bg-slate-900 dark:border-slate-700 dark:text-slate-100"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-slate-200">Исполнитель</label>
                <select
                  value={taskForm.assigneeId}
                  onChange={(e) => setTaskForm({ ...taskForm, assigneeId: Number(e.target.value) })}
                  className="w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500 dark:bg-slate-900 dark:border-slate-700 dark:text-slate-100"
                >
                  <option value={0}>Не назначать</option>
                  {allUsers.map((u) => (
                    <option key={u.id} value={u.id}>
                      {u.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div className="flex gap-3">
              <button
                onClick={() => setTaskModalOpen(false)}
                className="flex-1 rounded-full border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                Отмена
              </button>
              <button
                onClick={handleTaskSave}
                disabled={savingTask || !taskForm.title.trim() || !taskForm.stageId}
                className="flex-1 rounded-full bg-amber-500 px-4 py-2 text-sm font-semibold text-white hover:bg-amber-600 disabled:opacity-60 flex items-center justify-center gap-2"
              >
                {savingTask && <Loader2 className="w-4 h-4 animate-spin" />} Создать
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
