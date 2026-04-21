"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Skeleton } from "@/components/ui/skeleton";
import { Phase19Service } from "@/services/phase19";
import { AdminQualityReport } from "@/types/phase19";
import { Activity, Bot, Gauge, RefreshCcw, ShieldCheck, TestTube2, TimerReset } from "lucide-react";

const pct = (value: number) => `${value.toFixed(1)}%`;
const asCurrency = (value: number) => new Intl.NumberFormat("pt-BR", { style: "currency", currency: "USD" }).format(value);

const phaseLabel = (phase: string) =>
  phase.replace("phase_", "Fase ").replace("_", " ").replace(/\b\w/g, (char) => char.toUpperCase());

const flowLabel = (flow: string) => `Fluxo ${flow}`;

export function AdminQualityReportPanel() {
  const [report, setReport] = useState<AdminQualityReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const fetchReport = useCallback(async (silent = false) => {
    if (silent) {
      setRefreshing(true);
    } else {
      setLoading(true);
    }

    try {
      const response = await Phase19Service.getAdminQualityReport();
      setReport(response);
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível carregar o relatório de qualidade.");
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchReport();
  }, [fetchReport]);

  const judgeRows = useMemo(() => Object.entries(report?.judge_average_by_phase ?? {}).sort(([a], [b]) => a.localeCompare(b)), [report]);
  const executionRows = useMemo(() => Object.entries(report?.avg_execution_minutes_by_phase ?? {}).sort(([a], [b]) => a.localeCompare(b)), [report]);
  const costRows = useMemo(() => Object.entries(report?.average_cost_by_flow_type ?? {}).sort(([a], [b]) => a.localeCompare(b)), [report]);

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          {Array.from({ length: 4 }).map((_, idx) => (
            <Card key={idx} className="bg-card/70">
              <CardHeader>
                <Skeleton className="h-4 w-28" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-24" />
                <Skeleton className="mt-3 h-2 w-full" />
              </CardContent>
            </Card>
          ))}
        </div>
        <Skeleton className="h-80 w-full" />
      </div>
    );
  }

  if (!report) {
    return (
      <Card className="border-destructive/40 bg-destructive/5">
        <CardHeader>
          <CardTitle>Relatório indisponível</CardTitle>
          <CardDescription>Verifique conectividade com o backend e permissões de acesso ADMIN.</CardDescription>
        </CardHeader>
        <CardContent>
          <Button onClick={() => fetchReport()} variant="outline">Tentar novamente</Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-3 rounded-xl border border-primary/20 bg-gradient-to-r from-primary/10 via-transparent to-secondary/10 p-5 md:flex-row md:items-center md:justify-between">
        <div className="space-y-1">
          <h2 className="text-xl font-semibold tracking-tight">Pulse de Qualidade da Plataforma</h2>
          <p className="text-sm text-muted-foreground">
            Atualizado em {new Date(report.generated_at).toLocaleString("pt-BR")} · Base de {report.project_sample_size} projetos avaliados.
          </p>
        </div>
        <Button variant="secondary" className="w-fit gap-2" onClick={() => fetchReport(true)} disabled={refreshing}>
          <RefreshCcw className={`h-4 w-4 ${refreshing ? "animate-spin" : ""}`} />
          Atualizar agora
        </Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <Card className="border-border/80 bg-card/70">
          <CardHeader className="pb-2">
            <CardDescription>Cobertura de testes</CardDescription>
            <CardTitle className="flex items-center gap-2 text-2xl"><TestTube2 className="h-5 w-5 text-primary" /> {pct(report.test_coverage_percent)}</CardTitle>
          </CardHeader>
          <CardContent><Progress value={report.test_coverage_percent} className="h-2" /></CardContent>
        </Card>
        <Card className="border-border/80 bg-card/70">
          <CardHeader className="pb-2">
            <CardDescription>Taxa de sucesso da Tríade</CardDescription>
            <CardTitle className="flex items-center gap-2 text-2xl"><ShieldCheck className="h-5 w-5 text-emerald-400" /> {pct(report.triad_success_rate_percent)}</CardTitle>
          </CardHeader>
          <CardContent><Progress value={report.triad_success_rate_percent} className="h-2" /></CardContent>
        </Card>
        <Card className="border-border/80 bg-card/70">
          <CardHeader className="pb-2">
            <CardDescription>Uptime nos últimos 30 dias</CardDescription>
            <CardTitle className="flex items-center gap-2 text-2xl"><Gauge className="h-5 w-5 text-sky-400" /> {pct(report.platform_uptime_30d_percent)}</CardTitle>
          </CardHeader>
          <CardContent><Progress value={report.platform_uptime_30d_percent} className="h-2" /></CardContent>
        </Card>
        <Card className="border-border/80 bg-card/70">
          <CardHeader className="pb-2">
            <CardDescription>Projetos (concluídos vs abandonados)</CardDescription>
            <CardTitle className="flex items-center gap-2 text-2xl"><Activity className="h-5 w-5 text-orange-400" /> {report.projects_completed} / {report.projects_abandoned}</CardTitle>
          </CardHeader>
          <CardContent className="text-xs text-muted-foreground">Concluídos · Abandonados</CardContent>
        </Card>
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        <Card className="bg-card/70 lg:col-span-1">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base"><Bot className="h-4 w-4 text-primary" /> Score médio do LLM Judge</CardTitle>
            <CardDescription>Pontuação por fase do pipeline.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {judgeRows.length === 0 ? (
              <p className="text-sm text-muted-foreground">Nenhuma métrica de Judge encontrada ainda.</p>
            ) : judgeRows.map(([phase, score]) => (
              <div key={phase} className="space-y-1">
                <div className="flex items-center justify-between text-sm">
                  <span>{phaseLabel(phase)}</span>
                  <Badge variant="secondary">{score.toFixed(2)} / 10</Badge>
                </div>
                <Progress value={score * 10} className="h-1.5" />
              </div>
            ))}
          </CardContent>
        </Card>

        <Card className="bg-card/70 lg:col-span-1">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base"><TimerReset className="h-4 w-4 text-secondary" /> Tempo médio por fase</CardTitle>
            <CardDescription>Minutos médios para completar cada fase.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            {executionRows.length === 0 ? (
              <p className="text-sm text-muted-foreground">Sem dados de execução suficientes.</p>
            ) : executionRows.map(([phase, minutes]) => (
              <div key={phase} className="flex items-center justify-between rounded-md border border-border/60 px-3 py-2 text-sm">
                <span>{phaseLabel(phase)}</span>
                <span className="font-medium">{minutes.toFixed(1)} min</span>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card className="bg-card/70 lg:col-span-1">
          <CardHeader>
            <CardTitle className="text-base">Custo médio por fluxo</CardTitle>
            <CardDescription>Referência para forecast financeiro.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            {costRows.length === 0 ? (
              <p className="text-sm text-muted-foreground">Sem custos agregados no período analisado.</p>
            ) : costRows.map(([flow, cost]) => (
              <div key={flow} className="flex items-center justify-between rounded-md border border-border/60 px-3 py-2 text-sm">
                <span>{flowLabel(flow)}</span>
                <span className="font-semibold">{asCurrency(cost)}</span>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>

      {report.notes && report.notes.length > 0 && (
        <Card className="bg-card/60">
          <CardHeader>
            <CardTitle className="text-base">Notas técnicas do relatório</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm text-muted-foreground">
            {report.notes.map((note) => (
              <p key={note}>• {note}</p>
            ))}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
