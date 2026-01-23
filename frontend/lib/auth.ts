// frontend/lib/auth.ts
import { apiFetch, setToken, clearToken } from "./api";

type LoginResponse = {
  token: string;
};

type RegisterResponse = {
  token: string;
};

export async function login(params: {
  email: string;
  password: string;
  remember: boolean;
}) {
  const res = await apiFetch<LoginResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email: params.email, password: params.password }),
  });

  setToken(res.token, params.remember);
  return res;
}

export async function register(params: {
  email: string;
  password: string;
  org_id: number;
  role: string;
  remember: boolean;
}) {
  const res = await apiFetch<RegisterResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify({
      email: params.email,
      password: params.password,
      org_id: params.org_id,
      role: params.role,
    }),
  });

  setToken(res.token, params.remember);
  return res;
}

export function logoutLocal() {
  clearToken();
}
