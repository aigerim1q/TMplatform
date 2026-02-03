# Sprint 3 — Auth Integration (Frontend + Backend)

## Что сделано
- Интегрированы формы Login/Register с backend:
  - POST /auth/login
  - POST /auth/register
- JWT (accessToken) сохраняется в localStorage
- Route guard: без токена редирект на /login
- Подключены защищённые запросы:
  - GET /users/:id (профиль)
  - GET /hierarchy (иерархия)
- На backend включён CORS для запросов с frontend

## Быстрый старт (одним блоком команд)


# 1) Склонировать репозиторий и перейти в нужную ветку
git clone https://github.com/thatgingergirl/TMplatform.git
cd TMplatform
git checkout feat/sprint3-auth-integration

# 2) Backend: зависимости + запуск
cd backend
go mod tidy
go run ./cmd/server
# Backend: http://localhost:3001

# 3) (в новом терминале) Проверка, что backend работает
curl http://localhost:3001/health

# 4) Frontend: переменные окружения
cd ../frontend
printf "NEXT_PUBLIC_API_URL=http://localhost:3001\n" > .env.local
cat .env.local

# 5) Frontend: установка зависимостей + запуск
pnpm install
pnpm dev
# Frontend: http://localhost:3000 (или другой порт, если 3000 занят)
