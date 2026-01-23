import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/components/theme-provider";
import { cn } from "@/lib/utils";
import AIChatWidget from "@/components/ai-chat-widget";

const fontSans = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
  display: "swap"
});

export const metadata: Metadata = {
  title: "Next + Shadcn UI",
  description: "Starter kit with Tailwind and Shadcn components"
};

export default function RootLayout({
  children
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={cn(
          "min-h-screen bg-white text-slate-900 transition-colors dark:bg-slate-950 dark:text-slate-50",
          fontSans.variable,
          "font-sans"
        )}
      >
        <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
          {children}
          <AIChatWidget />
        </ThemeProvider>
      </body>
    </html>
  );
}

