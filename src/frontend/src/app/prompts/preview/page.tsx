"use client";

import Link from "next/link";
import { ArrowLeft } from "lucide-react";

import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { PromptCompositionPreview } from "@/components/prompts/PromptCompositionPreview";

export default function PromptPreviewPage() {
  return (
    <PrivateRoute>
      <div className="min-h-screen bg-background">
        <header className="border-b bg-background/90 backdrop-blur sticky top-0 z-30">
          <div className="container mx-auto flex h-16 items-center justify-between px-4">
            <Link href="/prompts" className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-4 w-4" /> Voltar para Prompts
            </Link>
          </div>
        </header>
        <main className="container mx-auto max-w-6xl p-4 md:p-8">
          <PromptCompositionPreview />
        </main>
      </div>
    </PrivateRoute>
  );
}
