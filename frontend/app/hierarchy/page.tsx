"use client";

import Header from "@/components/header";

type Status = "free" | "busy" | "sick";

type PersonNode = {
  title: string;
  name: string;
  status?: Status;
  note?: string;
  children?: PersonNode[];
  kind?: "dept" | "ceo" | "default";
};

const hierarchy: PersonNode = {
  title: "Генеральный директор",
  name: "Саламат Ахмедов",
  status: "busy",
  kind: "ceo",
  note: "Руководит всей компанией, утверждает проекты и формирует ключевую иерархию.",
  children: [
    {
      title: "Отдел внутреннего аудита",
      name: "Руководитель: Алмас Садыков",
      status: "free",
      kind: "dept",
      note: "Отвечает за проверки проектов, прозрачность процессов и контроль качества.",
      children: [
        {
          title: "Внутренний аудитор",
          name: "Айдын Рахимбаев",
          status: "busy",
          note: "Проводит внутренние аудиты, пишет комментарии по рискам и нарушениям.",
        },
        {
          title: "Внутренний аудитор",
          name: "Мади Ержанов",
          status: "free",
        },
      ],
    },
    {
      title: "Технический директор (главный инженер)",
      name: "Нуржан Ибраев",
      status: "busy",
      kind: "dept",
      note: "Курирует архитектурный, инженерный, конструкторский отделы и ПТО.",
      children: [
        {
          title: "Архитектурно-проектный отдел",
          name: "Руководитель: Ршыман Зейнулла",
          status: "busy",
          kind: "dept",
          children: [
            {
              title: "Главный архитектор",
              name: "Ршыман Зейнулла",
              status: "busy",
              note: "Утверждает архитектурные решения, планировки и фасады по объектам.",
            },
            {
              title: "Архитектор",
              name: "Омар Ахмед",
              status: "free",
              note: "Разрабатывает рабочие чертежи и узлы здания.",
            },
            {
              title: "Архитектор",
              name: "Диас Мухамеджанов",
              status: "busy",
            },
          ],
        },
        {
          title: "Отдел дизайнеров",
          name: "Руководитель: Алия Токжан",
          status: "busy",
          kind: "dept",
          note: "Делает дизайн экстерьера и интерьера, после чего подключается IT-отдел.",
          children: [
            {
              title: "Руководитель отдела дизайнеров",
              name: "Алия Токжан",
              status: "busy",
            },
            {
              title: "Дизайнер интерьера",
              name: "Айжан Курманова",
              status: "free",
            },
            {
              title: "Дизайнер экстерьера",
              name: "Темирлан Оспанов",
              status: "busy",
            },
          ],
        },
        {
          title: "ПТО",
          name: "Руководитель: Марат Алиев",
          status: "free",
          kind: "dept",
          children: [
            {
              title: "Начальник ПТО",
              name: "Марат Алиев",
              status: "busy",
            },
            {
              title: "Инженер ПТО",
              name: "Ильяс Койшыбаев",
              status: "free",
            },
          ],
        },
      ],
    },
    {
      title: "Директор по строительству",
      name: "Ербол Кенжебаев",
      status: "busy",
      kind: "dept",
      children: [
        {
          title: "Руководители строительных участков",
          name: "Куратор: Ернар Абдрахманов",
          status: "busy",
          kind: "dept",
          children: [
            {
              title: "Руководитель группы прорабов",
              name: "Ернар Абдрахманов",
              status: "busy",
            },
            {
              title: "Прораб",
              name: "Бекзат Жанабаев",
              status: "busy",
            },
            {
              title: "Прораб",
              name: "Расул Даулетов",
              status: "free",
            },
          ],
        },
      ],
    },
    {
      title: "IT-отдел",
      name: "Руководитель: Тимур Азимов",
      status: "busy",
      kind: "dept",
      note: "Занимается разработкой портала и внутренних систем после того, как дизайнеры завершат визуальную часть.",
      children: [
        {
          title: "Руководитель IT-отдела",
          name: "Тимур Азимов",
          status: "busy",
        },
        {
          title: "Frontend-разработчик",
          name: "Захар Ким",
          status: "free",
        },
        {
          title: "Backend-разработчик",
          name: "Мухаммед Алиев",
          status: "busy",
        },
        {
          title: "IT-поддержка",
          name: "Асель Ибрашева",
          status: "sick",
        },
      ],
    },
    {
      title: "Отдел кадров (HR)",
      name: "Руководитель: Динара Байжан",
      status: "free",
      kind: "dept",
      children: [
        {
          title: "Руководитель HR-отдела",
          name: "Динара Байжан",
          status: "busy",
        },
        {
          title: "HR-специалист",
          name: "Сания Рахматулла",
          status: "free",
          note: "В реальной системе именно HR создаёт профили, размещает в иерархии и отмечает статус (в том числе «болен»).",
        },
      ],
    },
    {
      title: "Юридический отдел",
      name: "Руководитель: Аскар Утегенов",
      status: "free",
      kind: "dept",
    },
    {
      title: "Коммерческий отдел / Отдел продаж",
      name: "Руководитель: Рустам Жаксылыков",
      status: "busy",
      kind: "dept",
    },
  ],
};

const statusMap: Record<Status, { label: string; dot: string; bg: string; text: string }> = {
  free: { label: "Свободен", dot: "#22c55e", bg: "#ecfdf3", text: "#166534" },
  busy: { label: "Занят", dot: "#f59e0b", bg: "#fef3c7", text: "#92400e" },
  sick: { label: "Болен", dot: "#ef4444", bg: "#fee2e2", text: "#b91c1c" },
};

function initials(name: string) {
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((n) => n[0]?.toUpperCase() ?? "?")
    .join("") || "?";
}

function NodeCard({ node }: { node: PersonNode }) {
  const status = node.status ? statusMap[node.status] : undefined;
  const typeClass = node.kind === "ceo" ? "ceo" : node.kind === "dept" ? "dept" : "";
  const statusClass = node.status ? `status-${node.status}` : "";

  return (
    <div className={`card ${typeClass} ${statusClass}`}>
      <div className="card-top">
        <div className="avatar-wrapper">
          <div className="avatar">{initials(node.name)}</div>
          <div className="person-main">
            <div className="title">{node.title}</div>
            <div className="name">{node.name}</div>
          </div>
        </div>
        {status && (
          <div className="status" style={{ color: status.text }}>
            <span className="status-dot" style={{ background: status.dot }} />
            <span>{status.label}</span>
          </div>
        )}
      </div>
      {node.note && <div className="note">{node.note}</div>}
    </div>
  );
}

function renderTree(node: PersonNode) {
  return (
    <li key={`${node.title}-${node.name}`}>
      <NodeCard node={node} />
      {node.children && node.children.length > 0 && (
        <ul>
          {node.children.map((child) => renderTree(child))}
        </ul>
      )}
    </li>
  );
}

export default function HierarchyPage() {
  return (
    <div className="min-h-screen bg-white text-slate-900">
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4">
        <Header />
      </div>

      <main className="mx-auto max-w-6xl px-4 py-8">
        <h1 className="text-center text-2xl font-bold">Иерархия строительной компании</h1>
        <p className="subtitle">Организационная структура без фотографий</p>

        <div className="org-chart">
          <ul>
            {renderTree(hierarchy)}
          </ul>
        </div>
      </main>

      <style jsx>{`
        .subtitle {
          text-align: center;
          font-size: 13px;
          color: #4b5563;
          margin-top: 4px;
          margin-bottom: 24px;
        }
        .org-chart {
          display: flex;
          justify-content: center;
          overflow-x: auto;
          padding-bottom: 40px;
        }
        .org-chart ul {
          padding-top: 20px;
          position: relative;
          padding-left: 0;
          display: flex;
          justify-content: center;
        }
        .org-chart ul ul { padding-top: 30px; }
        .org-chart li {
          list-style-type: none;
          text-align: center;
          position: relative;
          padding: 18px 8px 0 8px;
        }
        .org-chart li::before,
        .org-chart li::after {
          content: "";
          position: absolute;
          top: 0;
          right: 50%;
          border-top: 1px solid #e5e7eb;
          width: 50%;
          height: 18px;
        }
        .org-chart li::after {
          right: auto;
          left: 50%;
          border-left: 1px solid #e5e7eb;
        }
        .org-chart li:only-child::before,
        .org-chart li:only-child::after { display: none; }
        .org-chart li:only-child { padding-top: 0; }
        .org-chart li:first-child::before,
        .org-chart li:last-child::after { border: 0 none; }
        .org-chart li:last-child::before {
          border-right: 1px solid #e5e7eb;
          border-radius: 0 5px 0 0;
        }
        .org-chart li:first-child::after { border-radius: 5px 0 0 0; }
        .org-chart ul ul::before {
          content: "";
          position: absolute;
          top: 0;
          left: 50%;
          border-left: 1px solid #e5e7eb;
          width: 0;
          height: 26px;
        }
        .card {
          position: relative;
          background: #ffffff;
          border-radius: 14px;
          padding: 12px 14px;
          border: 1px solid #e5e7eb;
          min-width: 220px;
          max-width: 270px;
          margin: 0 auto;
          box-shadow: 0 10px 30px rgba(15,23,42,0.08);
          text-align: left;
          overflow: hidden;
        }
        .card::before {
          content: "";
          position: absolute;
          left: 0;
          top: 0;
          bottom: 0;
          width: 6px;
          background: #cbd5e1;
        }
        .card .card-top {
          display: flex;
          align-items: flex-start;
          justify-content: space-between;
          gap: 12px;
        }
        .card.status-free {
          background: #f7fef9;
          border-color: #bbf7d0;
        }
        .card.status-free::before { background: #22c55e; }
        .card.status-busy {
          background: #fffaf0;
          border-color: #fde68a;
        }
        .card.status-busy::before { background: #f59e0b; }
        .card.status-sick {
          background: #fff5f5;
          border-color: #fecdd3;
        }
        .card.status-sick::before { background: #ef4444; }
        .card.ceo {
          background: linear-gradient(135deg, #0b132b, #111827);
          color: #e5e7eb;
          border: 1px solid #1f2937;
          box-shadow: 0 18px 38px rgba(15,23,42,0.28);
        }
        .card.ceo::before { background: #2563eb; }
        .card.dept { background: #fdfcf5; }
        .avatar-wrapper {
          display: flex;
          align-items: center;
          gap: 10px;
          margin-bottom: 4px;
        }
        .avatar {
          width: 42px;
          height: 42px;
          border-radius: 12px;
          border: 1px solid #e5e7eb;
          background: #f3f4f6;
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 700;
          color: #9ca3af;
          flex-shrink: 0;
        }
        .person-main { display: flex; flex-direction: column; gap: 2px; }
        .title { font-size: 13px; font-weight: 700; color: #0f172a; }
        .card.ceo .title { color: #e5e7eb; }
        .name { font-size: 13px; color: #374151; }
        .card.ceo .name { color: #c7d2fe; }
        .status {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 4px 8px;
          border-radius: 999px;
          font-size: 11px;
          background: rgba(255,255,255,0.72);
          border: 1px solid #e5e7eb;
          box-shadow: 0 4px 10px rgba(15,23,42,0.06);
        }
        .card.ceo .status {
          background: rgba(255,255,255,0.1);
          border: 1px solid rgba(255,255,255,0.12);
          color: #f8fafc;
        }
        .status-dot {
          width: 7px;
          height: 7px;
          border-radius: 999px;
        }
        .note { font-size: 11px; color: #6b7280; margin-top: 6px; line-height: 1.45; }
        .card.ceo .note { color: #cbd5e1; }
        @media (max-width: 768px) {
          .card { min-width: 200px; }
        }
      `}</style>
    </div>
  );
}
