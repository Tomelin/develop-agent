"use client";

import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { useAuth } from "@/contexts/AuthContext";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { BarChart3, ShieldAlert } from "lucide-react";
import { AdminSettingsPanel } from "@/components/admin/AdminSettingsPanel";
import { FeatureFlagsPanel } from "@/components/admin/FeatureFlagsPanel";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function AdminSettingsPage() {
  const { user } = useAuth();

  return (
    <PrivateRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Admin · Configurações</h1>
          <p className="text-muted-foreground">Parâmetros globais da plataforma e controle de rollout por feature flags.</p>
        </div>

        {user?.role === "ADMIN" && (
          <Card className="border-primary/30 bg-primary/5">
            <CardContent className="flex flex-col gap-3 p-5 md:flex-row md:items-center md:justify-between">
              <p className="text-sm text-muted-foreground">
                Acompanhe os indicadores operacionais da PHASE-19 em um dashboard dedicado de qualidade.
              </p>
              <Link href="/dashboard/admin/quality-report">
                <Button variant="secondary" className="w-fit gap-2">
                  <BarChart3 className="h-4 w-4" /> Ver relatório de qualidade
                </Button>
              </Link>
            </CardContent>
          </Card>
        )}

        {user?.role !== "ADMIN" ? (
          <Card className="border-destructive/40 bg-destructive/5">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-destructive"><ShieldAlert className="h-5 w-5" /> Acesso restrito</CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-muted-foreground">Apenas usuários com role ADMIN podem acessar esta área.</CardContent>
          </Card>
        ) : (
          <div className="space-y-6">
            <AdminSettingsPanel />
            <FeatureFlagsPanel />
          </div>
        )}
      </div>
    </PrivateRoute>
  );
}
