"use client";

import Header from "@/components/header";
import DocumentsContent from "@/components/documents-content";

export default function DocumentsPage() {
  return (
    <div className="min-h-screen bg-background">
      <div className="flex justify-center pt-6">
        <Header />
      </div>
      <main className="mx-auto max-w-6xl px-6 py-8">
        <DocumentsContent />
      </main>
    </div>
  );
}
