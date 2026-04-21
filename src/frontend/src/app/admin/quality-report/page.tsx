"use client";

import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { useAuth } from "@/contexts/AuthContext";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ShieldAlert } from "lucide-react";
import { AdminQualityReportPanel } from "@/components/admin/AdminQualityReportPanel";

export default function AdminQualityReportPage() {
  const { user } = useAuth();

  return (
    <PrivateRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Admin · Relatório de Qualidade</h1>
          <p className="text-muted-foreground">
            Monitoramento contínuo da saúde da plataforma: cobertura, desempenho das fases, custos e confiabilidade.
          </p>
        </div>

        {user?.role !== "ADMIN" ? (
          <Card className="border-destructive/40 bg-destructive/5">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-destructive"><ShieldAlert className="h-5 w-5" /> Acesso restrito</CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-muted-foreground">Apenas usuários com role ADMIN podem acessar esta área.</CardContent>
          </Card>
        ) : (
          <AdminQualityReportPanel />
        )}
      </div>
    </PrivateRoute>
  );
}

