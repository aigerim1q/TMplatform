# TMplatform

Full platform (frontend, backend, auth, notifications) **with Chatbot** — like [aigerim1q/TMplatform](https://github.com/aigerim1q/TMplatform) plus the knowledge-planning chatbot.

## Project structure

- **`app/`** – Next.js app (dashboard, chat, hierarchy, documents, project)
- **`backend/`** – Go API (auth, users, hierarchy, notifications) — `http://localhost:3001`
- **`frontend/`** – Next.js frontend (alternate)
- **`components/`**, **`lib/`**, **`hooks/`** – Shared UI and utilities
- **`chatbot/`** – Knowledge-planning chatbot (RAG, projects/tasks, HTTP API + `/chat` UI)
- **`zhcp-parser-go/`** – Parser component
- **`parsing AI/`** – Parsing AI component

---

## Sprint 3 — Auth integration (frontend + backend)

- Login/Register with backend: `POST /auth/login`, `POST /auth/register`
- JWT (accessToken) in localStorage; route guard → `/login` when no token
- Protected: `GET /users/:id`, `GET /hierarchy`
- CORS enabled for frontend
- Notifications: `GET /notifications`, `PUT /notifications/:id/read` (JWT)

---

## Chatbot

The **chatbot** in `chatbot/` provides:

- Projects & tasks (create project, list, show, assign, plan)
- HTTP API: `POST /api/chat`, `GET /api/health`, in-browser UI at `GET /chat`
- Optional CLI frontend that uses the API when `CHATBOT_API_URL` is set

### Run the chatbot backend

```bash
cd "$(git rev-parse --show-toplevel)" && ./chatbot/TMplatform/run_chatbot_http.sh
```

Then open **http://localhost:8080/chat** for the in-browser chat, or use the main app’s Chat page (see below).

### Run the main app with chatbot

1. **Backend (platform API)**  
   ```bash
   cd "$(git rev-parse --show-toplevel)/backend" && go mod tidy && go run ./cmd/server
   ```  
   → http://localhost:3001

2. **Chatbot backend (optional, for real AI chat)**  
   ```bash
   cd "$(git rev-parse --show-toplevel)" && ./chatbot/TMplatform/run_chatbot_http.sh
   ```  
   → http://localhost:8080

3. **Frontend**  
   ```bash
   cd "$(git rev-parse --show-toplevel)" && printf "NEXT_PUBLIC_API_URL=http://localhost:3001\nNEXT_PUBLIC_CHATBOT_API_URL=http://localhost:8080\n" > .env.local && pnpm install && pnpm dev
   ```  
   → http://localhost:3000

The **Chat** page in the app uses the chatbot when `NEXT_PUBLIC_CHATBOT_API_URL` is set; otherwise it uses the built-in demo.

---

## Quick start (platform only)

```bash
git clone https://github.com/aigerim1q/TMplatform.git
cd TMplatform
```

**Backend** (terminal 1):
```bash
cd "$(git rev-parse --show-toplevel)/backend" && go mod tidy && go run ./cmd/server
```
→ http://localhost:3001

**Frontend** (terminal 2, replace path if your clone is elsewhere):
```bash
cd ~/TMplatform && printf "NEXT_PUBLIC_API_URL=http://localhost:3001\n" > .env.local && pnpm install && pnpm dev
```
→ http://localhost:3000

For **Chatbot** as well:
```bash
cd ~/TMplatform && ./chatbot/TMplatform/run_chatbot_http.sh
```
Then add `NEXT_PUBLIC_CHATBOT_API_URL=http://localhost:8080` to `.env.local`.
