## TM Platform – Frontend

Проект: Next.js 14 + Tailwind (shadcn/ui) для внутренних страниц (дашборд, календарь, иерархия, профиль и т.д.).

### Стек
- Next.js 14, React 18
- Tailwind CSS + shadcn/ui

### Требования
- Node.js 18+
- npm или pnpm (рекомендовано pnpm)

### Установка
```bash
cd frontend
pnpm install      # можно npm install
```

### Запуск
- Dev: `pnpm dev`
- Проверка: открыть http://localhost:3000

### Сборка и прод
- Build: `pnpm build`
- Start: `pnpm start` (использует собранный `.next`)

### Линт
- `pnpm lint`

### Структура
- `frontend/app` — страницы Next.js
- `frontend/components` — общие UI-компоненты
- `frontend/styles` — глобальные стили

### Переменные окружения
Если понадобятся, создайте `frontend/.env.local` (файл игнорируется Git). Сейчас обязательных переменных нет.
