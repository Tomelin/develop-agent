"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import type { ComponentType } from "react";
import {
  Activity,
  AlertCircle,
  CheckCircle2,
  Clock3,
  Database,
  RefreshCw,
  Server,
  Workflow,
  Wrench,
  XCircle,
} from "lucide-react";
import { statusService } from "@/services/statusService";
import {
  ComponentHealth,
  PlatformComponentStatus,
  PlatformStatusLevel,
  PlatformStatusResponse,
  StatusIncident,
} from "@/types/status";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Progress } from "@/components/ui/progress";

const POLLING_INTERVAL_MS = 30_000;

type ComponentConfig = {
  label: string;
  icon: ComponentType<{ className?: string }>;
};

const COMPONENT_CATALOG: Record<string, ComponentConfig> = {
  api_backend: { label: "API Backend", icon: Server },
  workers: { label: "Workers de IA", icon: Workflow },
  mongodb: { label: "MongoDB", icon: Database },
  redis: { label: "Redis", icon: Database },
  rabbitmq: { label: "RabbitMQ", icon: Activity },
  openai: { label: "OpenAI", icon: Wrench },
  anthropic: { label: "Anthropic", icon: Wrench },
  google: { label: "Google", icon: Wrench },
  ollama: { label: "Ollama", icon: Wrench },
};

const EXPECTED_COMPONENTS_ORDER = [
  "api_backend",
  "workers",
  "mongodb",
  "redis",
  "rabbitmq",
  "openai",
  "anthropic",
  "google",
  "ollama",
];

const statusTone: Record<PlatformStatusLevel, { label: string; className: string; icon: ComponentType<{ className?: string }> }> = {
  OPERATIONAL: {
    label: "Operacional",
    className: "bg-emerald-500/15 text-emerald-400 border-emerald-500/40",
    icon: CheckCircle2,
  },
  DEGRADED: {
    label: "Degradado",
    className: "bg-amber-500/15 text-amber-400 border-amber-500/40",
    icon: AlertCircle,
  },
  OUTAGE: {
    label: "Instável / Indisponível",
    className: "bg-red-500/15 text-red-400 border-red-500/40",
    icon: XCircle,
  },
};

const componentTone: Record<ComponentHealth, { label: string; chip: string; progress: string }> = {
  healthy: {
    label: "Saudável",
    chip: "bg-emerald-500/15 text-emerald-300 border-emerald-500/40",
    progress: "[&>div]:bg-emerald-500",
  },
  degraded: {
    label: "Degradado",
    chip: "bg-amber-500/15 text-amber-300 border-amber-500/40",
    progress: "[&>div]:bg-amber-500",
  },
  unhealthy: {
    label: "Indisponível",
    chip: "bg-red-500/15 text-red-300 border-red-500/40",
    progress: "[&>div]:bg-red-500",
  },
  unknown: {
    label: "Desconhecido",
    chip: "bg-slate-500/20 text-slate-300 border-slate-500/40",
    progress: "[&>div]:bg-slate-500",
  },
};

const getComponentStatus = (
  key: string,
  data?: PlatformStatusResponse | null,
): PlatformComponentStatus => {
  if (!data) return { status: "unknown" };

  if (data.components?.[key]) {
    return data.components[key];
  }

  if (data.providers?.[key]) {
    const providerStatus = data.providers[key].toLowerCase();
    if (providerStatus === "healthy" || providerStatus === "degraded" || providerStatus === "unhealthy") {
      return { status: providerStatus };
    }
  }

  return { status: "unknown" };
};

const formatCheckedAt = (checkedAt?: string): string => {
  if (!checkedAt) return "—";
  const date = new Date(checkedAt);
  return date.toLocaleString("pt-BR", {
    dateStyle: "short",
    timeStyle: "medium",
  });
};

const formatUptime = (uptime?: number): string => {
  if (typeof uptime !== "number") return "—";
  return `${uptime.toFixed(2)}%`;
};

const formatIncidentDate = (date: string) =>
  new Date(date).toLocaleString("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  });

export default function StatusPage() {
  const [statusData, setStatusData] = useState<PlatformStatusResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchStatus = useCallback(async (silent?: boolean) => {
    try {
      setError(null);
      if (silent) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }

      const response = await statusService.getPlatformStatus();
      setStatusData(response);
    } catch {
      setError("Não foi possível carregar o status da plataforma agora. Tente novamente em instantes.");
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(() => fetchStatus(true), POLLING_INTERVAL_MS);
    return () => clearInterval(interval);
  }, [fetchStatus]);

  const overview = useMemo(() => {
    if (!statusData) {
      return {
        healthy: 0,
        degraded: 0,
        unhealthy: 0,
      };
    }

    return EXPECTED_COMPONENTS_ORDER.reduce(
      (acc, key) => {
        const componentStatus = getComponentStatus(key, statusData).status;
        if (componentStatus === "healthy") acc.healthy += 1;
        if (componentStatus === "degraded") acc.degraded += 1;
        if (componentStatus === "unhealthy") acc.unhealthy += 1;
        return acc;
      },
      { healthy: 0, degraded: 0, unhealthy: 0 },
    );
  }, [statusData]);

  const incidents: StatusIncident[] = statusData?.incidents ?? [];
  const currentStatus = statusData?.status ?? "DEGRADED";
  const currentTone = statusTone[currentStatus] ?? statusTone.DEGRADED;
  const StatusIcon = currentTone.icon;

  return (
    <main className="mx-auto w-full max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
      <section className="rounded-3xl border border-border/70 bg-gradient-to-br from-card via-card/95 to-background p-6 shadow-xl shadow-black/20 sm:p-8">
        <div className="flex flex-col gap-5 md:flex-row md:items-start md:justify-between">
          <div className="space-y-3">
            <Badge className={`rounded-full border px-3 py-1 text-xs tracking-wide ${currentTone.className}`}>
              <StatusIcon className="mr-2 h-3.5 w-3.5" />
              Status Geral: {currentTone.label}
            </Badge>
            <h1 className="text-3xl font-semibold tracking-tight sm:text-4xl">Status da Plataforma</h1>
            <p className="max-w-2xl text-sm text-muted-foreground sm:text-base">
              Transparência operacional em tempo real. A página é atualizada a cada 30 segundos para refletir a saúde
              da API, infraestrutura e providers de IA.
            </p>
          </div>
          <div className="flex flex-col items-start gap-3 sm:items-end">
            <div className="flex items-center gap-2 text-xs text-muted-foreground">
              <Clock3 className="h-4 w-4" />
              Última atualização: {formatCheckedAt(statusData?.checked_at)}
            </div>
            <Button
              variant="outline"
              onClick={() => fetchStatus(true)}
              className="gap-2 border-border/70 bg-background/50"
              disabled={refreshing}
            >
              <RefreshCw className={`h-4 w-4 ${refreshing ? "animate-spin" : ""}`} />
              Atualizar agora
            </Button>
          </div>
        </div>
      </section>

      <section className="mt-6 grid gap-4 sm:grid-cols-3">
        <Card className="border-border/70 bg-card/80">
          <CardHeader className="pb-2">
            <CardDescription>Serviços saudáveis</CardDescription>
            <CardTitle className="text-3xl text-emerald-400">{overview.healthy}</CardTitle>
          </CardHeader>
        </Card>
        <Card className="border-border/70 bg-card/80">
          <CardHeader className="pb-2">
            <CardDescription>Serviços degradados</CardDescription>
            <CardTitle className="text-3xl text-amber-400">{overview.degraded}</CardTitle>
          </CardHeader>
        </Card>
        <Card className="border-border/70 bg-card/80">
          <CardHeader className="pb-2">
            <CardDescription>Serviços indisponíveis</CardDescription>
            <CardTitle className="text-3xl text-red-400">{overview.unhealthy}</CardTitle>
          </CardHeader>
        </Card>
      </section>

      {error && (
        <Card className="mt-6 border-red-500/40 bg-red-500/10">
          <CardContent className="flex items-center gap-2 p-4 text-sm text-red-200">
            <AlertCircle className="h-4 w-4" />
            {error}
          </CardContent>
        </Card>
      )}

      <section className="mt-6 grid gap-4 lg:grid-cols-2 xl:grid-cols-3">
        {EXPECTED_COMPONENTS_ORDER.map((key) => {
          const item = COMPONENT_CATALOG[key] ?? { label: key, icon: Server };
          const Icon = item.icon;
          const component = getComponentStatus(key, statusData);
          const tone = componentTone[component.status] ?? componentTone.unknown;
          const uptimeValue = component.uptime_30d_percent ?? (component.status === "healthy" ? 100 : 0);

          return (
            <Card key={key} className="border-border/70 bg-card/80 backdrop-blur">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between gap-3">
                  <div className="flex items-center gap-3">
                    <span className="rounded-xl border border-border/70 bg-background/40 p-2">
                      <Icon className="h-5 w-5 text-primary" />
                    </span>
                    <div>
                      <CardTitle className="text-base">{item.label}</CardTitle>
                      <CardDescription className="text-xs">Monitoramento contínuo</CardDescription>
                    </div>
                  </div>
                  <Badge className={`border ${tone.chip}`}>{tone.label}</Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <div className="flex items-center justify-between text-muted-foreground">
                  <span>Latência</span>
                  <span className="font-medium text-foreground">{component.latency_ms ? `${component.latency_ms}ms` : "—"}</span>
                </div>
                <div className="space-y-2">
                  <div className="flex items-center justify-between text-muted-foreground">
                    <span>Uptime (30d)</span>
                    <span className="font-medium text-foreground">{formatUptime(component.uptime_30d_percent)}</span>
                  </div>
                  <Progress value={uptimeValue} className={`h-1.5 bg-muted ${tone.progress}`} />
                </div>
                {component.error && (
                  <p className="rounded-md border border-red-500/30 bg-red-500/10 px-2 py-1 text-xs text-red-200">
                    {component.error}
                  </p>
                )}
              </CardContent>
            </Card>
          );
        })}
      </section>

      <section className="mt-8">
        <Card className="border-border/70 bg-card/80">
          <CardHeader>
            <CardTitle>Incidentes — últimos 7 dias</CardTitle>
            <CardDescription>Histórico público e resolução operacional de eventos.</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <p className="text-sm text-muted-foreground">Carregando incidentes...</p>
            ) : incidents.length === 0 ? (
              <p className="text-sm text-muted-foreground">Nenhum incidente registrado no período.</p>
            ) : (
              <div className="space-y-4">
                {incidents.map((incident, index) => (
                  <article key={`${incident.title}-${incident.occurred_at}-${index}`}>
                    <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
                      <h3 className="font-medium">{incident.title}</h3>
                      <span className="text-xs text-muted-foreground">{formatIncidentDate(incident.occurred_at)}</span>
                    </div>
                    <p className="text-sm text-muted-foreground">Duração: {incident.duration}</p>
                    <p className="mt-1 text-sm">{incident.resolution}</p>
                    {index < incidents.length - 1 && <Separator className="mt-4" />}
                  </article>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </section>
    </main>
  );
}
