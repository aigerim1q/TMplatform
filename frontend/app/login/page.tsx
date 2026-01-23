"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { cn } from "@/lib/utils";
import { login } from "@/lib/auth";

function SocialButton({
  children,
  "aria-label": ariaLabel,
}: {
  children: React.ReactNode;
  "aria-label": string;
}) {
  return (
    <button
      type="button"
      aria-label={ariaLabel}
      className="grid h-[38px] w-[38px] place-items-center rounded-full bg-[#FCFCFF] shadow-sm ring-1 ring-black/5 transition hover:scale-[1.02]"
    >
      {children}
    </button>
  );
}

function IconPlaceholder({ label }: { label: string }) {
  return (
    <span className="grid h-7 w-7 place-items-center rounded-md bg-black/5 text-[10px] font-medium text-black/60">
      {label}
    </span>
  );
}

export default function LoginPage() {
  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [remember, setRemember] = useState(false);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isValid = useMemo(() => {
    if (!email.trim() || !password.trim()) return false;
    return true;
  }, [email, password]);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      setLoading(true);

      await login({
        email: email.trim(),
        password,
        remember,
      });

      router.push("/dashboard");
    } catch (err: any) {
      // если в login() ты кидаешь Error(message), он сюда попадёт
      setError(err?.message || "Ошибка входа");
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-white text-black">
      <div className="mx-auto flex items-center justify-center px-8 py-16">
        <div className="flex w-full max-w-[1080px] flex-col items-center justify-center gap-10 lg:flex-row lg:gap-14">
          {/* Left */}
          <section className="flex-1">
            <div className="max-w-[558px]">
              <div className="mb-10 flex items-center gap-4">
                <div className="grid h-[60px] w-[60px] place-items-center rounded-full bg-[#F3C24C] text-sm font-bold">
                  THE
                </div>
                <div className="text-2xl font-bold text-[#8B6B4E]">Quryls</div>
              </div>

              <h1 className="text-[37px] font-bold leading-tight">
                Управление строительными проектами
              </h1>
              <p className="mt-4 text-lg leading-relaxed text-[#505872]">
                Комплексная платформа для эффективного управления строительными
                процессами, контроля качества и координации команд.
              </p>

              <div className="mt-10 space-y-6">
                <div className="flex items-start gap-4">
                  <div className="grid h-[41px] w-[41px] place-items-center rounded-xl bg-black/5">
                    <IconPlaceholder label="A" />
                  </div>
                  <div>
                    <div className="text-[16.5px] font-medium">
                      Аналитика проектов
                    </div>
                    <div className="mt-1 text-sm text-black/60">
                      Отслеживайте процесс в реальном времени
                    </div>
                  </div>
                </div>

                <div className="flex items-start gap-4">
                  <div className="grid h-[41px] w-[41px] place-items-center rounded-xl bg-black/5">
                    <IconPlaceholder label="T" />
                  </div>
                  <div>
                    <div className="text-[16.5px] font-medium">
                      Командная работа
                    </div>
                    <div className="mt-1 text-sm text-black/60">
                      Координация всех участников проекта
                    </div>
                  </div>
                </div>

                <div className="flex items-start gap-4">
                  <div className="grid h-[41px] w-[41px] place-items-center rounded-xl bg-black/5">
                    <IconPlaceholder label="S" />
                  </div>
                  <div>
                    <div className="text-[16.5px] font-medium">
                      Безопасность данных
                    </div>
                    <div className="mt-1 text-sm text-black/60">
                      Надежная защита корпоративной информации
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </section>

          {/* Right */}
          <section className="w-full max-w-[450px]">
            <Card className="rounded-[15px] border-0 bg-[rgba(255,255,255,0.5)] shadow-[0px_4px_90px_0px_rgba(240,230,218,1)]">
              <CardHeader className="pb-0">
                <div className="mx-auto mt-[35px] w-[254px] text-center">
                  <div className="text-[24.5px] font-semibold">
                    Вход в систему
                  </div>
                  <div className="mt-2 text-base text-[#747D8A]">
                    Войдите в свою учетную запись
                  </div>
                </div>
              </CardHeader>

              <CardContent className="px-[39px] pb-8 pt-[37px]">
                <form onSubmit={onSubmit} className="space-y-6">
                  <div className="space-y-2">
                    <Label className="text-[14.5px] font-medium text-[#4C5464]">
                      Электронная почта
                    </Label>
                    <Input
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="your@email.com"
                      autoComplete="email"
                      className={cn(
                        "h-[45px] rounded-[10px] border-[#FFF2E6] bg-[rgba(255,252,249,0.6)]"
                      )}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-[14.5px] font-medium text-[#4C5464]">
                      Пароль
                    </Label>
                    <Input
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      placeholder="Введите пароль"
                      type="password"
                      autoComplete="current-password"
                      className={cn(
                        "h-[45px] rounded-[10px] border-[#FFF2E6] bg-[rgba(255,252,249,0.6)]"
                      )}
                    />
                  </div>

                  <div className="flex items-center justify-between text-[12.5px] font-medium">
                    <label className="flex items-center gap-2 text-[#545C6B]">
                      <Checkbox
                        checked={remember}
                        onCheckedChange={(v) => setRemember(v === true)}
                        className="h-4 w-4 rounded-[2px] border-[#969697] bg-[#F5F7FF]"
                      />
                      Запомнить меня
                    </label>

                    <Link
                      href="/forgot-password"
                      className="text-[#BA8D51] hover:underline"
                    >
                      Забыли пароль?
                    </Link>
                  </div>

                  {error && (
                    <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
                      {error}
                    </div>
                  )}

                  <Button
                    type="submit"
                    disabled={!isValid || loading}
                    className="h-[49px] w-full rounded-[11px] bg-[#C19A6B] text-white hover:bg-[#C19A6B]/90 disabled:opacity-60"
                  >
                    {loading ? "Входим..." : "Войти"}
                  </Button>

                  <div className="pt-2 text-center text-xs text-[#58616F]">
                    Войти через корпоративные сервисы
                  </div>

                  <div className="flex justify-center gap-6 pt-1">
                    <SocialButton aria-label="Corporate service 1">
                      <IconPlaceholder label="S1" />
                    </SocialButton>
                    <SocialButton aria-label="Corporate service 2">
                      <IconPlaceholder label="S2" />
                    </SocialButton>
                    <SocialButton aria-label="Corporate service 3">
                      <IconPlaceholder label="S3" />
                    </SocialButton>
                  </div>

                  <div className="pt-2">
                    <div className="mx-auto h-px w-[372px] bg-[#C7C7C7]" />
                  </div>

                  <div className="pt-4 text-center text-sm">
                    <span className="text-[#58616F]">Нет учетной записи?</span>{" "}
                    <Link
                      href="/register"
                      className="font-medium text-[#BC915A] hover:underline"
                    >
                      Зарегистрироваться
                    </Link>
                  </div>
                </form>
              </CardContent>
            </Card>
          </section>
        </div>
      </div>
    </main>
  );
}
