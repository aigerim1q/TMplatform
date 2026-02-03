'use client';

import React from "react"
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Clock, Plus, X, Users, UserPlus, ChevronDown, ChevronRight } from 'lucide-react';

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
      className="rounded-2xl bg-white p-5 shadow-sm cursor-pointer hover:shadow-md transition-shadow"
    >
      <div className="flex items-center justify-between mb-3">
        <span className="text-sm text-gray-500">Проект: {project}</span>
        <div className="flex items-center gap-1.5 rounded-full bg-amber-50 px-3 py-1">
          <Clock size={12} className="text-gray-600" />
          <span className="text-sm text-gray-700">{time}</span>
          <span className={`h-2.5 w-2.5 rounded-full ${statusColors[timeStatus]}`} />
        </div>
      </div>
      <h3 className="font-bold text-base text-gray-900 mb-2 line-clamp-1">{title}</h3>
      <p className="text-sm text-gray-600 line-clamp-3">{description}</p>
    </div>
  );
}

interface ProjectCardProps {
  id: string;
  name: string;
  days: string;
  image: string;
  onClick?: () => void;
}

function ProjectCard({ id, name, days, image, onClick }: ProjectCardProps) {
  return (
    <div 
      onClick={onClick}
      className="relative rounded-2xl overflow-hidden aspect-[16/10] cursor-pointer hover:shadow-lg transition-shadow"
    >
      <img
        src={image || "/placeholder.svg"}
        alt={name}
        className="w-full h-full object-cover"
      />
      <div className="absolute inset-0 bg-gradient-to-b from-black/50 via-transparent to-black/30" />
      <div className="absolute top-4 left-4 right-4 flex items-center justify-between">
        <span className="text-white font-semibold text-base drop-shadow-md">Проект: {name}</span>
        <div className="flex items-center gap-1.5 rounded-full bg-white/95 px-3 py-1">
          <Clock size={12} className="text-gray-600" />
          <span className="text-sm text-gray-700">{days}</span>
          <span className="h-2.5 w-2.5 rounded-full bg-green-500" />
        </div>
      </div>
      <button className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center justify-center w-8 h-8 bg-white/90 rounded-full shadow-md hover:bg-white">
        <ChevronRight size={18} className="text-gray-600" />
      </button>
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
    <div className="inline-flex items-center gap-2 rounded-full bg-amber-100/80 px-4 py-2">
      <span className={`h-3 w-3 rounded-full ${dotColors[color]}`} />
      <span className="font-medium text-gray-800 text-sm">{title}: {count}</span>
    </div>
  );
}

function AddButton({ text, onClick }: { text: string; onClick?: () => void }) {
  return (
    <button 
      onClick={onClick}
      className="flex items-center gap-2 text-gray-700 hover:text-gray-900 font-medium text-sm"
    >
      <Plus size={16} />
      {text}
    </button>
  );
}

interface AddTaskModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (task: { project: string; title: string; description: string; deadline: string }) => void;
}

function AddTaskModal({ isOpen, onClose, onSubmit }: AddTaskModalProps) {
  const [project, setProject] = useState('Shyraq');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [deadline, setDeadline] = useState('');

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ project, title, description, deadline });
    setTitle('');
    setDescription('');
    setDeadline('');
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl p-6 w-full max-w-md shadow-xl">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900">Добавить задачу</h2>
          <button 
            onClick={onClose}
            className="p-1 rounded-full hover:bg-gray-100"
          >
            <X size={20} className="text-gray-500" />
          </button>
        </div>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Проект
            </label>
            <select
              value={project}
              onChange={(e) => setProject(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500"
            >
              <option value="Shyraq">Shyraq</option>
              <option value="Ansau">Ansau</option>
              <option value="Dariya">Dariya</option>
            </select>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Название задачи
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Введите название задачи"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500"
              required
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Описание
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Введите описание задачи"
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none"
              required
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Срок выполнения
            </label>
            <input
              type="date"
              value={deadline}
              onChange={(e) => setDeadline(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500"
              required
            />
          </div>
          
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-full text-gray-700 font-medium hover:bg-gray-50 bg-transparent"
            >
              Отмена
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-2 bg-amber-500 text-white rounded-full font-medium hover:bg-amber-600"
            >
              Создать
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
      <div className="bg-amber-50 rounded-2xl p-6 w-full max-w-md shadow-xl">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2">
            <Users size={24} className="text-gray-700" />
            <h2 className="text-xl font-bold text-gray-900">Ответственнные</h2>
          </div>
          <div className="flex items-center gap-3">
            <button className="flex items-center gap-1 text-gray-600 hover:text-gray-900">
              <UserPlus size={18} />
              <span className="text-sm">Делегировать</span>
            </button>
            <button 
              onClick={onClose}
              className="p-1 rounded-full hover:bg-amber-100"
            >
              <X size={20} className="text-gray-500" />
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
                  <span className="font-semibold text-gray-900">{person.name}</span>
                  {person.isYou && (
                    <span className="text-xs bg-gray-300 text-gray-700 px-2 py-0.5 rounded-full">
                      это вы
                    </span>
                  )}
                </div>
                <span className="text-sm text-amber-700">{person.role}</span>
              </div>
            </div>
          ))}
        </div>

        {/* Quick Delegation */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">
              Быстрое делегирование
            </span>
            <span className="text-xs text-gray-500">
              Выбрать ответственного
            </span>
          </div>
          <div className="relative">
            <select
              value={selectedResponsible}
              onChange={(e) => setSelectedResponsible(e.target.value)}
              className="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl appearance-none focus:ring-2 focus:ring-amber-500 focus:border-amber-500 text-gray-700"
            >
              <option value="">Выбрать ответст...</option>
              <option value="omar">Омар Ахмет</option>
              <option value="zeinulla">Зейнулла Рышман</option>
              <option value="aidyn">Айдын Рахимбаев</option>
            </select>
            <ChevronDown size={20} className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none" />
          </div>
        </div>

        {/* Add Button */}
        <button className="w-full flex items-center justify-center gap-2 bg-amber-200 hover:bg-amber-300 text-amber-900 font-semibold py-3 rounded-xl transition-colors">
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
  const [myTasks, setMyTasks] = useState([
    {
      id: 'task-1',
      project: 'Shyraq',
      time: '-9 часов',
      timeStatus: 'danger' as const,
      title: 'Возведение колонн на 1 э...',
      description: 'Интерьеры подъездов и этажей Перед началом работ нужно провести подготовку...',
      responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    },
    {
      id: 'task-2',
      project: 'Ansau',
      time: '2 дня',
      timeStatus: 'success' as const,
      title: 'Нужно сделать новую пла...',
      description: 'На объекте Ansau срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ....',
      responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    },
    {
      id: 'task-3',
      project: 'Dariya',
      time: '20 дней',
      timeStatus: 'success' as const,
      title: 'Нужно сделать новую пла...',
      description: 'На объекте Dariya срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ....',
      responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    },
  ]);

  const urgentTasks = [
    {
      id: 'urgent-1',
      project: 'Shyraq',
      time: '18 часов',
      timeStatus: 'success' as const,
      title: 'Нужно привезти 10 плиток',
      description: 'На объекте Shyraq срочно требуется новая плитка, 5x3 метра, 2 штуки для наличников вокруг л....',
      responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    },
  ];

  const projects = [
    {
      id: 'shyraq',
      name: 'Shyraq',
      days: '25 дней',
      image: '/images/building-1.jpg',
    },
    {
      id: 'ansau',
      name: 'Ansau',
      days: '55 дней',
      image: '/images/building-2.jpg',
    },
    {
      id: 'dariya',
      name: 'Dariya',
      days: '55 дней',
      image: '/images/building-3.jpg',
    },
  ];

  const subordinateTasks = [
    {
      id: 'sub-task-1',
      project: 'Shyraq',
      time: '-9 часов',
      timeStatus: 'danger' as const,
      title: 'Нужно сделать новую пла...',
      description: 'На объекте Shyraq срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ....',
      responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    },
    {
      id: 'sub-task-2',
      project: 'Ansau',
      time: '2 дня',
      timeStatus: 'success' as const,
      title: 'Нужно сделать новую пла...',
      description: 'На объекте Ansau срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ....',
      responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    },
    {
      id: 'sub-task-3',
      project: 'Dariya',
      time: '20 дней',
      timeStatus: 'success' as const,
      title: 'Нужно сделать новую пла...',
      description: 'На объекте Dariya срочно требуется новая планировку 5-го этажа, чтобы совпадала с новым ....',
      responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    },
  ];

  const handleAddTask = (task: { project: string; title: string; description: string; deadline: string }) => {
    const newTask = {
      id: `task-${Date.now()}`,
      project: task.project,
      time: '7 дней',
      timeStatus: 'success' as const,
      title: task.title.length > 25 ? task.title.slice(0, 25) + '...' : task.title,
      description: task.description.length > 80 ? task.description.slice(0, 80) + '...' : task.description,
      responsible: 'Омар Ахмет, Зейнулла Рышман, Серик Рах...',
    };
    setMyTasks([...myTasks, newTask]);
  };

  const handleTaskClick = (taskId: string) => {
    router.push(`/project/${taskId}`);
  };

  const handleProjectClick = (projectId: string) => {
    router.push(`/project-overview/${projectId}`);
  };

  const handleResponsibleClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsResponsibleModalOpen(true);
  };

  return (
    <main className="w-full max-w-5xl mx-auto px-6 py-8">
      {/* Add Task Modal */}
      <AddTaskModal 
        isOpen={isAddTaskModalOpen} 
        onClose={() => setIsAddTaskModalOpen(false)}
        onSubmit={handleAddTask}
      />

      {/* Responsible Modal */}
      <ResponsibleModal
        isOpen={isResponsibleModalOpen}
        onClose={() => setIsResponsibleModalOpen(false)}
      />

      {/* Urgent Tasks Section */}
      <div className="mb-8">
        <SectionHeader color="green" title="Срочные задачи" count={1} />
        <div className="mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {urgentTasks.map((task) => (
            <TaskCard 
              key={task.id} 
              {...task} 
              onClick={() => handleTaskClick(task.id)}
              onResponsibleClick={handleResponsibleClick}
            />
          ))}
        </div>
      </div>

      {/* My Tasks Section */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <SectionHeader color="red" title="Мои задачи" count={10} />
          <AddButton text="Добавить задачу" onClick={() => setIsAddTaskModalOpen(true)} />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {myTasks.map((task) => (
            <TaskCard 
              key={task.id} 
              {...task} 
              onClick={() => handleTaskClick(task.id)}
              onResponsibleClick={handleResponsibleClick}
            />
          ))}
        </div>
      </div>

      {/* Projects Section */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <SectionHeader color="yellow" title="Проекты" count={3} />
          <AddButton text="Добавить проект" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {projects.map((project) => (
            <ProjectCard 
              key={project.id} 
              {...project} 
              onClick={() => handleProjectClick(project.id)}
            />
          ))}
        </div>
      </div>

      {/* Subordinate Tasks Section */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <SectionHeader color="red" title="Задачи подчинённых" count={subordinateTasks.length} />
          <AddButton text="Добавить проект" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {subordinateTasks.map((task) => (
            <TaskCard 
              key={task.id} 
              {...task} 
              onClick={() => handleTaskClick(task.id)}
              onResponsibleClick={handleResponsibleClick}
            />
          ))}
        </div>
      </div>

      {/* Subordinate Projects Section */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <SectionHeader color="yellow" title="Проекты подчиненных" count={3} />
          <AddButton text="Добавить проект" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {projects.map((project) => (
            <ProjectCard 
              key={project.id} 
              {...project} 
              onClick={() => handleProjectClick(project.id)}
            />
          ))}
        </div>
      </div>
    </main>
  );
}
