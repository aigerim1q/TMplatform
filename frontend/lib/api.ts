// Unified API helpers (supports legacy tokens and new access tokens)
const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001"
const TOKEN_KEY = "accessToken"
const LEGACY_TOKEN_KEY = "token"

export function getApiUrl() {
  return API_URL.replace(/\/$/, "")
}

// Alias to keep both naming conventions working
export const getApiBaseUrl = () => getApiUrl()

export function getAccessToken() {
  if (typeof window === "undefined") return null
  return (
    localStorage.getItem(TOKEN_KEY) ||
    sessionStorage.getItem(TOKEN_KEY) ||
    localStorage.getItem(LEGACY_TOKEN_KEY) ||
    sessionStorage.getItem(LEGACY_TOKEN_KEY)
  )
}

// Legacy alias
export function getToken() {
  return getAccessToken()
}

export function setToken(token: string, remember: boolean) {
  if (typeof window === "undefined") return

  if (remember) {
    localStorage.setItem(TOKEN_KEY, token)
    localStorage.setItem(LEGACY_TOKEN_KEY, token)
    sessionStorage.removeItem(TOKEN_KEY)
    sessionStorage.removeItem(LEGACY_TOKEN_KEY)
  } else {
    sessionStorage.setItem(TOKEN_KEY, token)
    sessionStorage.setItem(LEGACY_TOKEN_KEY, token)
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(LEGACY_TOKEN_KEY)
  }
}

export function clearToken() {
  if (typeof window === "undefined") return
  localStorage.removeItem(TOKEN_KEY)
  sessionStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(LEGACY_TOKEN_KEY)
  sessionStorage.removeItem(LEGACY_TOKEN_KEY)
}

export async function apiFetch<T>(
  path: string,
  opts: RequestInit & { auth?: boolean } = {},
): Promise<T> {
  const url = `${getApiUrl()}${path.startsWith("/") ? "" : "/"}${path}`

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(opts.headers as Record<string, string> | undefined),
  }

  if (opts.auth) {
    const token = getAccessToken()
    if (token) headers.Authorization = `Bearer ${token}`
  }

  const res = await fetch(url, {
    ...opts,
    headers,
  })

  const text = await res.text()
  const data = text ? safeJson(text) : null

  if (!res.ok) {
    const msg =
      (data && (data.error || data.message)) ||
      `Request failed: ${res.status} ${res.statusText}`
    throw new Error(msg)
  }

  return data as T
}

export async function fetchJson<T>(path: string, options: RequestInit = {}): Promise<T> {
  const base = getApiBaseUrl()
  const token = getAccessToken()

  const headers = new Headers(options.headers || {})
  const isFormData = typeof FormData !== "undefined" && options.body instanceof FormData

  if (!headers.has("Content-Type") && options.body && !isFormData) {
    headers.set("Content-Type", "application/json")
  }
  if (token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const res = await fetch(`${base}${path}`, {
    ...options,
    headers,
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Request failed with status ${res.status}`)
  }

  return (await res.json()) as T
}

function safeJson(text: string) {
  try {
    return JSON.parse(text)
  } catch {
    return null
  }
}