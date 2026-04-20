"use client";

import { useAuth } from "@/contexts/AuthContext";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Sparkles, Users, Activity } from "lucide-react";

export default function DashboardPage() {
  const { user } = useAuth();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Bem-vindo de volta, <span className="text-primary font-medium">{user?.name}</span>!
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card className="bg-card/50 backdrop-blur-sm border-primary/20 hover:border-primary/50 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Projetos Ativos</CardTitle>
            <Sparkles className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">0</div>
            <p className="text-xs text-muted-foreground">Novos projetos esta semana</p>
          </CardContent>
        </Card>

        <Card className="bg-card/50 backdrop-blur-sm border-secondary/20 hover:border-secondary/50 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Equipe</CardTitle>
            <Users className="h-4 w-4 text-secondary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1</div>
            <p className="text-xs text-muted-foreground">Usuários registrados no sistema</p>
          </CardContent>
        </Card>

        <Card className="bg-card/50 backdrop-blur-sm border-muted-foreground/20 hover:border-muted-foreground/50 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Atividade da IA</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">0</div>
            <p className="text-xs text-muted-foreground">Ações automatizadas concluídas</p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4 bg-card/50 backdrop-blur-sm">
          <CardHeader>
            <CardTitle>Visão Geral</CardTitle>
            <CardDescription>Resumo de atividades dos últimos 7 dias</CardDescription>
          </CardHeader>
          <CardContent className="pl-2 h-[300px] flex items-center justify-center text-muted-foreground">
            Gráficos serão implementados nas próximas fases.
          </CardContent>
        </Card>

        <Card className="col-span-3 bg-card/50 backdrop-blur-sm">
          <CardHeader>
            <CardTitle>Atividades Recentes</CardTitle>
            <CardDescription>Histórico de operações no sistema</CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-center h-[300px] text-muted-foreground">
            Nenhuma atividade registrada ainda.
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
