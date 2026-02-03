"use client";

import Link from "next/link";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

export default function ForgotPasswordSentPage() {
  return (
    <main className="min-h-screen bg-white text-black">
      <div className="mx-auto flex min-h-screen items-center justify-center px-6 py-16">
        <Card className="w-full max-w-[483px] rounded-[17px] border-0 bg-[rgba(255,255,255,0.5)] shadow-[0_0_30px_rgba(0,0,0,0.06)]">
          <CardContent className="px-10 py-12 text-center">
            <div className="mx-auto grid h-[46px] w-[46px] place-items-center rounded-2xl bg-black/5">
              {/* icon placeholder */}
              <span className="text-xs font-semibold text-black/60">Mail</span>
            </div>

            <h1 className="mt-6 text-[27px] font-bold">
              Письмо для сброса отправлено
            </h1>

            <p className="mx-auto mt-6 max-w-[385px] text-[18px] text-[#78889D]">
              Инструкции по сбросу пароля отправлены на ваш Email. Пожалуйста,
              проверьте папку Входящие и Спам.
            </p>

            <div className="mt-10">
              <Button
                asChild
                className="h-[49px] w-full rounded-[7px] bg-[#C19A6B] text-white hover:bg-[#C19A6B]/90"
              >
                <Link href="/login">Вернуться ко Входу →</Link>
              </Button>
            </div>

            <div className="mt-12 h-px w-full bg-[#D9D9D9]" />

            <div className="mt-8 space-y-3">
              <button type="button" className="text-[15px] text-[#78889D] hover:underline">
                Не получили письмо? Отправить снова
              </button>
              <div className="text-xs text-[#78889D]">
                Нужна помощь?{" "}
                <a href="#" className="hover:underline">
                  Связаться с поддержкой
                </a>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </main>
  );
}

