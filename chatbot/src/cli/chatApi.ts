/**
 * Client for the Go chatbot HTTP API (POST /api/chat).
 * Used when CHATBOT_API_URL is set to connect the frontend (main) to the backend.
 */

export interface ChatSession {
  active_project_id: string;
  active_project_title: string;
}

export interface ChatRequest {
  message: string;
  active_project_id?: string;
  active_project_title?: string;
}

export interface ChatResponse {
  reply: string;
  exit: boolean;
  active_project_id: string;
  active_project_title: string;
}

export function getChatbotApiUrl(): string | undefined {
  const url = process.env.CHATBOT_API_URL || process.env.CHATBOT_HTTP_URL;
  return url ? url.replace(/\/$/, '') : undefined;
}

export async function sendToChatbot(
  baseUrl: string,
  message: string,
  session: ChatSession
): Promise<{ reply: string; exit: boolean; session: ChatSession }> {
  const endpoint = `${baseUrl}/api/chat`;
  const body: ChatRequest = {
    message,
    active_project_id: session.active_project_id || '',
    active_project_title: session.active_project_title || '',
  };
  const res = await fetch(endpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Chatbot API error ${res.status}: ${text || res.statusText}`);
  }
  const data = (await res.json()) as ChatResponse;
  return {
    reply: data.reply ?? '',
    exit: data.exit ?? false,
    session: {
      active_project_id: data.active_project_id ?? '',
      active_project_title: data.active_project_title ?? '',
    },
  };
}
