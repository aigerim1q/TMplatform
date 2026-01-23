"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import Link from "next/link";

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    firstName: "",
    lastName: "",
    email: "",
    password: "",
    confirmPassword: "",
    agreeToTerms: false,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Handle form submission
    console.log("Form submitted:", formData);
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-white p-8">
      <div className="flex w-full max-w-[1221px] flex-col items-center justify-center gap-[73px] lg:flex-row">
        {/* Left Section - Logo and Features */}
        <div className="w-full max-w-[558px] space-y-8">
          {/* Logo */}
          <div className="mb-8">
            <div className="text-2xl font-bold text-[#8B6B4E]">Qurylys</div>
          </div>

          {/* Title and Description */}
          <div className="space-y-4">
            <h1 className="text-[37px] font-bold leading-tight text-black">
              Управление строительными проектами
            </h1>
            <p className="text-lg text-[#505872] leading-relaxed">
              Комплексная платформа для эффективного управления строительными процессами, контроля качества и координации команд.
            </p>
          </div>

          {/* Features */}
          <div className="space-y-6 mt-8">
            <div className="flex items-start gap-4">
              <div className="w-10 h-10 flex-shrink-0 bg-gray-100 rounded"></div>
              <div>
                <h3 className="text-base font-medium text-black mb-1">
                  Аналитика проектов
                </h3>
                <p className="text-sm text-black">
                  Отслеживайте процесс в реальном времени
                </p>
              </div>
            </div>

            <div className="flex items-start gap-4">
              <div className="w-10 h-10 flex-shrink-0 bg-gray-100 rounded"></div>
              <div>
                <h3 className="text-base font-medium text-black mb-1">
                  Командная работа
                </h3>
                <p className="text-sm text-black">
                  Координация всех участников проекта
                </p>
              </div>
            </div>

            <div className="flex items-start gap-4">
              <div className="w-10 h-10 flex-shrink-0 bg-gray-100 rounded"></div>
              <div>
                <h3 className="text-base font-medium text-black mb-1">
                  Безопасность данных
                </h3>
                <p className="text-sm text-black">
                  Надежная защита корпоративной информации
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Right Section - Registration Card */}
        <Card className="w-full max-w-[517px] rounded-xl border-0 bg-[rgba(250,250,255,0.5)] shadow-[0px_4px_90px_0px_rgba(255,235,211,1)]">
          <CardHeader className="space-y-4 pb-6">
            <div className="text-center space-y-2">
              <Link 
                href="/login" 
                className="text-sm font-semibold text-black hover:underline"
              >
                Уже есть аккаунт? Войти
              </Link>
              <CardTitle className="text-[26px] font-semibold text-black">
                Создать аккаунт
              </CardTitle>
              <CardDescription className="text-base text-[#51789E]">
                Начните управлять проектами эффективно
              </CardDescription>
            </div>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Name Fields */}
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="firstName">Имя</Label>
                  <Input
                    id="firstName"
                    placeholder="Ахмет"
                    value={formData.firstName}
                    onChange={(e) =>
                      setFormData({ ...formData, firstName: e.target.value })
                    }
                    className="h-12"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="lastName">Фамилия</Label>
                  <Input
                    id="lastName"
                    placeholder="Омар"
                    value={formData.lastName}
                    onChange={(e) =>
                      setFormData({ ...formData, lastName: e.target.value })
                    }
                    className="h-12"
                  />
                </div>
              </div>

              {/* Email */}
              <div className="space-y-2">
                <Label htmlFor="email">Электронная почта</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="your@email.com"
                  value={formData.email}
                  onChange={(e) =>
                    setFormData({ ...formData, email: e.target.value })
                  }
                  className="h-12"
                />
              </div>

              {/* Password */}
              <div className="space-y-2">
                <Label htmlFor="password">Пароль</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  value={formData.password}
                  onChange={(e) =>
                    setFormData({ ...formData, password: e.target.value })
                  }
                  className="h-12"
                />
              </div>

              {/* Confirm Password */}
              <div className="space-y-2">
                <Label htmlFor="confirmPassword">Подтвердите пароль</Label>
                <Input
                  id="confirmPassword"
                  type="password"
                  placeholder="••••••••"
                  value={formData.confirmPassword}
                  onChange={(e) =>
                    setFormData({ ...formData, confirmPassword: e.target.value })
                  }
                  className="h-12"
                />
              </div>

              {/* Terms Checkbox */}
              <div className="flex items-start gap-3">
                <Checkbox
                  id="terms"
                  checked={formData.agreeToTerms}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, agreeToTerms: checked === true })
                  }
                  className="mt-1"
                />
                <Label
                  htmlFor="terms"
                  className="text-sm text-[#005088] leading-relaxed cursor-pointer"
                >
                  Я соглашаюсь с Условиями использования и Политикой конфиденциальность.
                </Label>
              </div>

              {/* Submit Button */}
              <Button
                type="submit"
                className="w-full h-12 bg-[#C19A6B] hover:bg-[#C19A6B]/90 text-white font-semibold rounded-lg"
              >
                Зарегистрироваться
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
