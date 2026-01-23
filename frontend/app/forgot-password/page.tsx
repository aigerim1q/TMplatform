"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const canSubmit = useMemo(() => email.trim().length > 0, [email]);

  function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    const trimmed = email.trim();
    // Simple email validation
    const isEmailLike = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(trimmed);
    if (!isEmailLike) {
      setError("Введите корректный email адрес.");
      return;
    }
    // Mock "not found" case - redirect to error page (replace with API call)
    if (trimmed.toLowerCase().endsWith("@example.com")) {
      router.push("/forgot-password/error-email");
      return;
    }
  
    setError(null);
    console.log({ email: trimmed });
    router.push("/forgot-password/sent");
  }

  return (
    <main className="min-h-screen bg-white text-black">
      <div className="relative mx-auto flex min-h-screen items-center justify-center px-6 py-16">
        <Card className="w-full max-w-[480px] rounded-[14px] border-0 bg-[rgba(255,255,255,0.5)] shadow-[0_0_30px_rgba(0,0,0,0.06)]">
          <CardContent className="px-[41px] py-10 text-center">
            <h1 className="text-center text-[27px] font-bold">
              Восстановить пароль
            </h1>
            <p className="mx-auto mt-4 max-w-[393px] text-center text-base text-[#455D7B]">
              Введите ваш email адрес для сброса пароля.
            </p>

            <form onSubmit={onSubmit} className="mt-8 space-y-5">
              <div className="space-y-2 text-left">
                <Label className="text-base font-normal">Email</Label>
                <Input
                  value={email}
                  onChange={(e) => {
                    setEmail(e.target.value);
                    if (error) setError(null);
                  }}
                  placeholder="omarakhmet@gmail.com"
                  className={cn(
                    "h-[50px] rounded-[10px] bg-white px-4",
                    error ? "border-[#EF7C7D] bg-[#FDF7F7]" : "border-[#D7E1EB]"
                  )}
                />
                {error ? (
                  <p className="max-w-[364px] text-sm text-[#EB3223]">{error}</p>
                ) : null}
              </div>

              <Button
                type="submit"
                disabled={!canSubmit}
                className="h-[50px] w-full rounded-[10px] bg-[#D1C7B3] text-white hover:bg-[#D1C7B3]/90 disabled:opacity-60"
              >
                Отправить запрос →
              </Button>

              <div className="pt-2 text-center text-sm text-[#969696]">
                <Link href="/login" className="inline-flex items-center gap-2">
                  <span>→</span>
                  <span>Вернуться ко Входу</span>
                </Link>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </main>
  );
}

