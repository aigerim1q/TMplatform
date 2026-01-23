"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

export default function ForgotPasswordErrorEmailPage() {
  const router = useRouter();

  const handleTryAgain = () => {
    router.push("/forgot-password");
  };

  return (
    <main className="min-h-screen bg-white text-black">
      <div className="relative mx-auto flex min-h-screen items-center justify-center px-6 py-16">
        <Card className="w-full max-w-[480px] rounded-[14px] border-0 bg-[rgba(255,255,255,0.5)] shadow-[0_0_30px_rgba(0,0,0,0.06)]">
          <CardContent className="px-[41px] py-10 text-center">
            <div className="mx-auto grid h-[50px] w-[50px] place-items-center rounded-full bg-[#FEE9E9]">
              <svg 
                className="h-6 w-6 text-[#EB3223]" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24" 
                xmlns="http://www.w3.org/2000/svg"
              >
                <path 
                  strokeLinecap="round" 
                  strokeLinejoin="round" 
                  strokeWidth={2} 
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>

            <h1 className="mt-6 text-center text-[27px] font-bold text-[#EB3223]">
              Email не найден
            </h1>
            
            <p className="mx-auto mt-4 max-w-[393px] text-center text-base text-[#455D7B]">
              Пользователь с таким Email не найден. Пожалуйста, проверьте введенный Email или зарегистрируйтесь.
            </p>

            <div className="mt-8 space-y-4">
              <Button
                onClick={handleTryAgain}
                className="h-[50px] w-full rounded-[10px] bg-[#D1C7B3] text-white hover:bg-[#D1C7B3]/90"
              >
                Попробовать снова
              </Button>

              <Button
                asChild
                variant="outline"
                className="h-[50px] w-full rounded-[10px] border-[#D7E1EB] text-[#455D7B] hover:bg-[#F8F9FA]"
              >
                <Link href="/register">Зарегистрироваться</Link>
              </Button>
            </div>

            <div className="mt-6 pt-6 text-center text-sm text-[#969696] border-t border-[#D7E1EB]">
              <Link href="/login" className="inline-flex items-center gap-2 hover:underline">
                <span>→</span>
                <span>Вернуться ко Входу</span>
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
