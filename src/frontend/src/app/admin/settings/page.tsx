"use client";

import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { useAuth } from "@/contexts/AuthContext";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ShieldAlert } from "lucide-react";
import { AdminSettingsPanel } from "@/components/admin/AdminSettingsPanel";
import { FeatureFlagsPanel } from "@/components/admin/FeatureFlagsPanel";

export default function AdminSettingsPage() {
  const { user } = useAuth();

  return (
    <PrivateRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Admin · Configurações</h1>
          <p className="text-muted-foreground">Parâmetros globais da plataforma e controle de rollout por feature flags.</p>
        </div>

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
