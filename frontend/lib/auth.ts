import { apiFetch, setToken, clearToken } from "./api"

type LoginResponse = {
  accessToken?: string
  token?: string
  name?: string
  avatar_url?: string | null
  email?: string
  org_id?: number
  role?: string
}

type RegisterResponse = {
  accessToken?: string
  token?: string
  id?: number
  email?: string
  org_id?: number
  role?: string
  name?: string
}

export async function login(params: {
  email: string
  password: string
  remember: boolean
}) {
  const res = await apiFetch<LoginResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email: params.email, password: params.password }),
  })

  const token = res.accessToken || res.token
  if (!token) {
    throw new Error("Не удалось получить токен доступа")
  }

  setToken(token, params.remember)
  return res
}

export async function register(params: {
  email: string
  password: string
  org_id: number
  role: string
  name?: string
  remember: boolean
}) {
  const res = await apiFetch<RegisterResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify({
      email: params.email,
      password: params.password,
      org_id: params.org_id,
      role: params.role,
      name: params.name,
    }),
  })

  const token = res.accessToken || res.token
  if (!token) {
    throw new Error("Не удалось получить токен доступа")
  }

  setToken(token, params.remember)
  return res
}

export function logoutLocal() {
  clearToken()
}
