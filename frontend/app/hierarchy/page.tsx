'use client';

import Header from '@/components/header';
import { cn } from '@/lib/utils';

type StatusType = 'free' | 'busy' | 'sick';

interface Person {
  title: string;
  name: string;
  status: StatusType;
  note?: string;
  children?: Person[];
  isDept?: boolean;
  isCeo?: boolean;
}

const orgData: Person = {
  title: 'Генеральный директор',
  name: 'Саламат Ахмедов',
  status: 'busy',
  isCeo: true,
  note: 'Руководит всей компанией, утверждает проекты и формирует ключевую иерархию.',
  children: [
    {
      title: 'Отдел внутреннего аудита',
      name: 'Руководитель: Алмас Садыков',
      status: 'free',
      isDept: true,
      note: 'Отвечает за проверки проектов, прозрачность процессов и контроль качества.',
      children: [
        {
          title: 'Внутренний аудитор',
          name: 'Айдын Рахимбаев',
          status: 'busy',
          note: 'Проводит внутренние аудиты, пишет комментарии по рискам и нарушениям.',
        },
        {
          title: 'Внутренний аудитор',
          name: 'Мади Ержанов',
          status: 'free',
        },
      ],
    },
    {
      title: 'Технический директор (главный инженер)',
      name: 'Нуржан Ибраев',
      status: 'busy',
      isDept: true,
      note: 'Курирует архитектурный, инженерный, конструкторский отделы и ПТО.',
      children: [
        {
          title: 'Архитектурно-проектный отдел',
          name: 'Руководитель: Ршыман Зейнулла',
          status: 'busy',
          isDept: true,
          children: [
            {
              title: 'Главный архитектор',
              name: 'Ршыман Зейнулла',
              status: 'busy',
              note: 'Утверждает архитектурные решения, планировки и фасады по объектам.',
            },
            {
              title: 'Архитектор',
              name: 'Омар Ахмед',
              status: 'free',
              note: 'Разрабатывает рабочие чертежи и узлы здания.',
            },
            {
              title: 'Архитектор',
              name: 'Диас Мухамеджанов',
              status: 'busy',
            },
          ],
        },
        {
          title: 'Отдел дизайнеров',
          name: 'Руководитель: Алия Токжан',
          status: 'busy',
          isDept: true,
          note: 'Делает дизайн экстерьера и интерьера, после чего подключается IT-отдел.',
          children: [
            {
              title: 'Руководитель отдела дизайнеров',
              name: 'Алия Токжан',
              status: 'busy',
            },
            {
              title: 'Дизайнер интерьера',
              name: 'Айжан Курманова',
              status: 'free',
            },
            {
              title: 'Дизайнер экстерьера',
              name: 'Темирлан Оспанов',
              status: 'busy',
            },
          ],
        },
        {
          title: 'ПТО',
          name: 'Руководитель: Марат Алиев',
          status: 'free',
          isDept: true,
          children: [
            {
              title: 'Начальник ПТО',
              name: 'Марат Алиев',
              status: 'busy',
            },
            {
              title: 'Инженер ПТО',
              name: 'Ильяс Койшыбаев',
              status: 'free',
            },
          ],
        },
      ],
    },
    {
      title: 'Директор по строительству',
      name: 'Ербол Кенжебаев',
      status: 'busy',
      isDept: true,
      children: [
        {
          title: 'Руководители строительных участков',
          name: 'Куратор: Ернар Абдрахманов',
          status: 'busy',
          isDept: true,
          children: [
            {
              title: 'Руководитель группы прорабов',
              name: 'Ернар Абдрахманов',
              status: 'busy',
            },
            {
              title: 'Прораб',
              name: 'Бекзат Жанабаев',
              status: 'busy',
            },
            {
              title: 'Прораб',
              name: 'Расул Даулетов',
              status: 'free',
            },
          ],
        },
      ],
    },
    {
      title: 'IT-отдел',
      name: 'Руководитель: Тимур Азимов',
      status: 'busy',
      isDept: true,
      note: 'Занимается разработкой портала и внутренних систем после того, как дизайнеры завершат визуальную часть.',
      children: [
        {
          title: 'Руководитель IT-отдела',
          name: 'Тимур Азимов',
          status: 'busy',
        },
        {
          title: 'Frontend-разработчик',
          name: 'Захар Ким',
          status: 'free',
        },
        {
          title: 'Backend-разработчик',
          name: 'Мухаммед Алиев',
          status: 'busy',
        },
        {
          title: 'IT-поддержка',
          name: 'Асель Ибрашева',
          status: 'sick',
        },
      ],
    },
    {
      title: 'Отдел кадров (HR)',
      name: 'Руководитель: Динара Байжан',
      status: 'free',
      isDept: true,
      children: [
        {
          title: 'Руководитель HR-отдела',
          name: 'Динара Байжан',
          status: 'busy',
        },
        {
          title: 'HR-специалист',
          name: 'Сания Рахматулла',
          status: 'free',
          note: 'В реальной системе именно HR создаёт профили, размещает в иерархии и отмечает статус (в том числе «болен»).',
        },
      ],
    },
    {
      title: 'Юридический отдел',
      name: 'Руководитель: Аскар Утегенов',
      status: 'free',
      isDept: true,
    },
    {
      title: 'Коммерческий отдел / Отдел продаж',
      name: 'Руководитель: Рустам Жаксылыков',
      status: 'busy',
      isDept: true,
    },
  ],
};

function StatusBadge({ status }: { status: StatusType }) {
  const statusConfig = {
    free: { bg: 'bg-emerald-50 dark:bg-emerald-500/15', text: 'text-emerald-700 dark:text-emerald-200', dot: 'bg-emerald-500', label: 'Свободен' },
    busy: { bg: 'bg-amber-50 dark:bg-amber-500/15', text: 'text-amber-700 dark:text-amber-200', dot: 'bg-amber-500', label: 'Занят' },
    sick: { bg: 'bg-red-50 dark:bg-red-500/15', text: 'text-red-700 dark:text-red-200', dot: 'bg-red-500', label: 'Болен' },
  };

  const config = statusConfig[status];

  return (
    <span className={cn('inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs', config.bg, config.text)}>
      <span className={cn('w-1.5 h-1.5 rounded-full', config.dot)} />
      {config.label}
    </span>
  );
}

function PersonCard({ person }: { person: Person }) {
  return (
    <div
      className={cn(
        'bg-white rounded-xl p-3 border min-w-[210px] max-w-[260px] mx-auto shadow-sm text-left dark:bg-slate-900 dark:border-slate-800 dark:shadow-[0_8px_24px_rgba(0,0,0,0.35)]',
        person.isDept && 'bg-gray-50 dark:bg-slate-800/70'
      )}
    >
      <div className="flex items-center gap-2.5 mb-1.5">
        <div className="w-10 h-10 rounded-full border-2 border-gray-200 bg-gray-100 flex items-center justify-center text-gray-400 font-bold flex-shrink-0 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300">
          ?
        </div>

        <div className="flex flex-col gap-0.5">
          <div className="text-sm font-semibold text-gray-900 leading-tight dark:text-slate-100">
            {person.title}
          </div>
          <div className="text-sm text-gray-600 dark:text-slate-300">{person.name}</div>
          <StatusBadge status={person.status} />
        </div>
      </div>

      {person.note && (
        <p className="text-xs text-gray-500 mt-1.5 dark:text-slate-400">{person.note}</p>
      )}

      {/* ✅ КНОПКА ВЫБРАТЬ */}
      <button
        className="mt-3 w-full rounded-lg bg-[#cdbb9a] py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50 dark:bg-amber-400 dark:text-slate-900"
        disabled={person.status === 'busy'}
      >
        Выбрать
      </button>
    </div>
  );
}


function OrgNode({ person }: { person: Person }) {
  const hasChildren = person.children && person.children.length > 0;

  return (
    <li className="relative text-center pt-5 px-2">
      <PersonCard person={person} />
      {hasChildren && (
        <ul className="flex justify-center pt-7 relative before:content-[''] before:absolute before:top-0 before:left-1/2 before:border-l before:border-gray-300 before:h-6">
          {person.children!.map((child, index) => (
            <OrgNode key={index} person={child} />
          ))}
        </ul>
      )}
    </li>
  );
}

export default function HierarchyPage() {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 dark:bg-slate-950 dark:text-slate-100">
      {/* Header - centered */}
      <div className="flex justify-center pt-6">
        <Header />
      </div>

      {/* Content */}
      <div className="px-6 py-8">
        <h1 className="text-xl font-semibold text-center mb-2">Иерархия строительной компании</h1>
        <p className="text-sm text-gray-500 text-center mb-8 dark:text-slate-300">Организационная структура с отделами и сотрудниками</p>

        {/* Org Chart */}
        <div className="overflow-x-auto pb-10">
          <div className="flex justify-center min-w-max">
            <ul className="flex justify-center org-chart">
              <OrgNode person={orgData} />
            </ul>
          </div>
        </div>
      </div>

      <style jsx>{`
        .org-chart li::before,
        .org-chart li::after {
          content: "";
          position: absolute;
          top: 0;
          right: 50%;
          border-top: 1px solid #d1d5db;
          width: 50%;
          height: 20px;
        }
        .org-chart li::after {
          right: auto;
          left: 50%;
          border-left: 1px solid #d1d5db;
        }
        .org-chart li:only-child::before,
        .org-chart li:only-child::after {
          display: none;
        }
        .org-chart li:only-child {
          padding-top: 0;
        }
        .org-chart li:first-child::before,
        .org-chart li:last-child::after {
          border: 0 none;
        }
        .org-chart li:last-child::before {
          border-right: 1px solid #d1d5db;
          border-radius: 0 5px 0 0;
        }
        .org-chart li:first-child::after {
          border-radius: 5px 0 0 0;
        }
        :global(.dark) .org-chart li::before,
        :global(.dark) .org-chart li::after {
          border-color: #475569;
        }
        :global(.dark) .org-chart li:last-child::before {
          border-color: #475569;
        }
      `}</style>
    </div>
  );
}
