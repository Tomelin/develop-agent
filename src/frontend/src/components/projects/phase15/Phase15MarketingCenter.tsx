"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { Project } from "@/types/project";
import { Phase5CodeFile } from "@/types/phase5";
import { MarketingChannel, MarketingManualBrief, Phase15DeliveryReport } from "@/types/phase15";
import { ProjectService } from "@/services/project";
import { Phase15Service } from "@/services/phase15";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Progress } from "@/components/ui/progress";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Markdown } from "@/components/ui/markdown";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { toast } from "sonner";
import { CalendarRange, Download, Loader2, Megaphone, Rocket, Webhook } from "lucide-react";

interface Phase15MarketingCenterProps {
  project: Project;
}

interface CalendarEntry {
  date: string;
  channel: MarketingChannel;
  assetType: string;
  title: string;
  bestTimeUtc: string;
  cta: string;
}

const CHANNELS: { value: MarketingChannel; label: string; color: string }[] = [
  { value: "linkedin", label: "LinkedIn", color: "bg-sky-100 text-sky-700 border-sky-300" },
  { value: "instagram", label: "Instagram", color: "bg-rose-100 text-rose-700 border-rose-300" },
  { value: "google-ads", label: "Google Ads", color: "bg-amber-100 text-amber-700 border-amber-300" },
];

const DEFAULT_BRIEF: MarketingManualBrief = {
  product_name: "",
  problem_solved: "",
  target_audience: "",
  main_benefits: ["", "", ""],
  communication_tone: "consultivo",
  primary_cta: "Agendar demonstração",
  secondary_cta: "Saiba mais",
};

function parseCsvCalendar(content: string): CalendarEntry[] {
  const lines = content.split(/\r?\n/).filter(Boolean);
  if (lines.length <= 1) return [];

  return lines.slice(1).map((line) => {
    const [date, channel, assetType, title, bestTimeUtc, cta] = line.split(",");
    return {
      date,
      channel: channel as MarketingChannel,
      assetType,
      title,
      bestTimeUtc,
      cta,
    };
  });
}

function buildIcs(entries: CalendarEntry[]): string {
  const lines = ["BEGIN:VCALENDAR", "VERSION:2.0", "PRODID:-//Develop Agent Frontend//Marketing Calendar//PT-BR"];
  entries.forEach((entry, index) => {
    const dt = entry.date.replaceAll("-", "");
    lines.push("BEGIN:VEVENT");
    lines.push(`UID:${index}-${entry.channel}@develop-agent`);
    lines.push(`DTSTAMP:${new Date().toISOString().replace(/[-:]/g, "").replace(/\.\d{3}/, "")}`);
    lines.push(`DTSTART;VALUE=DATE:${dt}`);
    lines.push(`SUMMARY:${entry.title}`);
    lines.push(`DESCRIPTION:Canal ${entry.channel} | CTA ${entry.cta} | Horário ${entry.bestTimeUtc} UTC`);
    lines.push("END:VEVENT");
  });
  lines.push("END:VCALENDAR");
  return lines.join("\n");
}

export function Phase15MarketingCenter({ project }: Phase15MarketingCenterProps) {
  const [files, setFiles] = useState<Phase5CodeFile[]>([]);
  const [running, setRunning] = useState(false);
  const [streamProgress, setStreamProgress] = useState(0);
  const [useLinkedProject, setUseLinkedProject] = useState(Boolean(project.linked_project_id));
  const [channels, setChannels] = useState<MarketingChannel[]>(["linkedin", "instagram", "google-ads"]);
  const [budgetUsd, setBudgetUsd] = useState(3000);
  const [manualBrief, setManualBrief] = useState<MarketingManualBrief>(DEFAULT_BRIEF);
  const [deliveryReport, setDeliveryReport] = useState<Phase15DeliveryReport | null>(null);
  const [selectedEntry, setSelectedEntry] = useState<CalendarEntry | null>(null);
  const [webhookUrl, setWebhookUrl] = useState("");

  const fetchArtifacts = useCallback(async () => {
    try {
      const allFiles = await ProjectService.getProjectFiles(project.id);
      setFiles(allFiles.filter((file) => file.path.startsWith("artifacts/marketing/") || file.path.startsWith("docs/prompts/marketing/")));
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível carregar os artefatos do Fluxo C.");
    }
  }, [project.id]);

  useEffect(() => {
    const bootstrap = setTimeout(() => {
      void fetchArtifacts();
    }, 0);

    const timer = setInterval(() => void fetchArtifacts(), 12000);
    return () => {
      clearTimeout(bootstrap);
      clearInterval(timer);
    };
  }, [fetchArtifacts]);

  useEffect(() => {
    if (!running) return;
    const interval = setInterval(() => {
      setStreamProgress((value) => Math.min(95, value + 6));
    }, 500);

    return () => clearInterval(interval);
  }, [running]);

  const strategyFile = files.find((file) => file.path === "artifacts/marketing/strategy/MARKETING_STRATEGY.md");
  const forecastFile = files.find((file) => file.path === "artifacts/marketing/PERFORMANCE_FORECAST.md");
  const calendarCsv = files.find((file) => file.path === "artifacts/marketing/strategy/CALENDAR.csv")?.content ?? "";

  const entries = useMemo(() => parseCsvCalendar(calendarCsv), [calendarCsv]);
  const filteredEntries = useMemo(() => entries.filter((entry) => channels.includes(entry.channel)), [entries, channels]);

  const monthSeeds = useMemo(() => {
    const now = new Date();
    return [0, 1, 2].map((offset) => new Date(now.getFullYear(), now.getMonth() + offset, 1));
  }, []);

  const validationErrors = useMemo(() => {
    if (useLinkedProject) return [] as string[];
    const errors: string[] = [];
    if (!manualBrief.product_name.trim()) errors.push("Nome do produto é obrigatório.");
    if (!manualBrief.problem_solved.trim()) errors.push("Problema resolvido é obrigatório.");
    if (!manualBrief.target_audience.trim()) errors.push("Público-alvo é obrigatório.");
    if (manualBrief.main_benefits.map((item) => item.trim()).filter(Boolean).length < 3) errors.push("Informe ao menos 3 benefícios principais.");
    return errors;
  }, [manualBrief, useLinkedProject]);

  const handleChannel = (channel: MarketingChannel, checked: boolean) => {
    setChannels((prev) => {
      if (checked) return Array.from(new Set([...prev, channel]));
      const next = prev.filter((item) => item !== channel);
      return next.length ? next : prev;
    });
  };

  const runPhase15 = async () => {
    if (channels.length === 0) {
      toast.error("Selecione ao menos um canal.");
      return;
    }
    if (validationErrors.length > 0) {
      toast.error("Revise o brief antes de iniciar.");
      return;
    }

    try {
      setRunning(true);
      setStreamProgress(0);
      const report = await Phase15Service.run(project.id, {
        use_linked_project: useLinkedProject,
        channels,
        monthly_budget_usd: budgetUsd,
        manual_brief: {
          ...manualBrief,
          main_benefits: manualBrief.main_benefits.map((item) => item.trim()).filter(Boolean),
        },
      });

      setDeliveryReport(report);
      setStreamProgress(100);
      toast.success(`Fluxo C executado: ${report.total_pieces} peças geradas.`);
      await fetchArtifacts();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao executar a phase 15.");
    } finally {
      setRunning(false);
    }
  };

  const downloadPack = async () => {
    try {
      const { blob, pieces, filename } = await Phase15Service.downloadPack(project.id, channels);
      const anchor = document.createElement("a");
      anchor.href = URL.createObjectURL(blob);
      anchor.download = filename ?? `marketing-pack-${project.id}.zip`;
      document.body.appendChild(anchor);
      anchor.click();
      URL.revokeObjectURL(anchor.href);
      anchor.remove();
      toast.success(`Download iniciado (${pieces} peças de conteúdo).`);
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível exportar o pack de marketing.");
    }
  };

  const exportCalendarCsv = () => {
    if (!filteredEntries.length) {
      toast.error("Nenhum item no calendário para exportação.");
      return;
    }

    const csv = ["date,channel,asset_type,title,best_time_utc,cta", ...filteredEntries.map((entry) => `${entry.date},${entry.channel},${entry.assetType},${entry.title},${entry.bestTimeUtc},${entry.cta}`)].join("\n");
    const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });
    const anchor = document.createElement("a");
    anchor.href = URL.createObjectURL(blob);
    anchor.download = `marketing-calendar-${project.id}.csv`;
    anchor.click();
    URL.revokeObjectURL(anchor.href);
  };

  const exportCalendarIcs = () => {
    if (!filteredEntries.length) {
      toast.error("Nenhum item no calendário para exportação.");
      return;
    }

    const blob = new Blob([buildIcs(filteredEntries)], { type: "text/calendar;charset=utf-8" });
    const anchor = document.createElement("a");
    anchor.href = URL.createObjectURL(blob);
    anchor.download = `marketing-calendar-${project.id}.ics`;
    anchor.click();
    URL.revokeObjectURL(anchor.href);
  };

  const configureWebhook = async () => {
    if (!webhookUrl.trim()) {
      toast.error("Informe a URL do webhook.");
      return;
    }

    try {
      const result = await Phase15Service.configureWebhook(project.id, webhookUrl.trim());
      toast.success(`Webhook validado (${result.last_test.response_status}).`);
      await fetchArtifacts();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao validar webhook.");
    }
  };

  return (
    <div className="space-y-6">
      <Card className="border-border/60 bg-card/70 backdrop-blur-sm">
        <CardHeader className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Megaphone className="h-5 w-5 text-primary" />
              Fluxo C — Marketing Multi-Canal
            </CardTitle>
            <CardDescription>
              Operação completa da phase 15 com calendário interativo, exportáveis (CSV/ICS), download de pack por canal e configuração de webhooks.
            </CardDescription>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" onClick={downloadPack} className="gap-2">
              <Download className="h-4 w-4" /> Download pack
            </Button>
            <Button onClick={runPhase15} disabled={running} className="gap-2">
              {running ? <Loader2 className="h-4 w-4 animate-spin" /> : <Rocket className="h-4 w-4" />} Executar Phase 15
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {running && (
            <div className="space-y-2 rounded-xl border border-primary/30 bg-primary/5 p-3">
              <p className="text-sm font-medium">Gerando estratégia e conteúdos por canal...</p>
              <Progress value={streamProgress} className="h-2" />
            </div>
          )}

          <div className="grid gap-4 lg:grid-cols-3">
            <Card className="border-border/60 bg-background/40 lg:col-span-2">
              <CardHeader>
                <CardTitle className="text-base">Configuração da execução</CardTitle>
                <CardDescription>Sem mocks: todos os dados exibidos são lidos dos artefatos reais gerados pelo backend.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between rounded-lg border border-border/60 p-3">
                  <Label htmlFor="use-linked">Herdar contexto do Fluxo A vinculado</Label>
                  <Switch id="use-linked" checked={useLinkedProject} onCheckedChange={setUseLinkedProject} />
                </div>

                <div className="space-y-2">
                  <Label>Orçamento mensal (USD)</Label>
                  <Input type="number" min={500} step={100} value={budgetUsd} onChange={(event) => setBudgetUsd(Number(event.target.value || 0))} />
                </div>

                <div className="space-y-2">
                  <Label>Canais ativos</Label>
                  <div className="grid gap-2 sm:grid-cols-3">
                    {CHANNELS.map((channel) => (
                      <label key={channel.value} className="flex items-center gap-2 rounded-lg border border-border/60 bg-card px-3 py-2 text-sm">
                        <Checkbox checked={channels.includes(channel.value)} onCheckedChange={(checked) => handleChannel(channel.value, Boolean(checked))} />
                        {channel.label}
                      </label>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="border-border/60 bg-background/40">
              <CardHeader>
                <CardTitle className="text-base">Resumo da geração</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <p>Canal(is): <span className="font-medium">{channels.length}</span></p>
                <p>Peças totais: <span className="font-medium">{deliveryReport?.total_pieces ?? 0}</span></p>
                <p>Fonte do brief: <span className="font-medium">{deliveryReport?.brief_source ?? "-"}</span></p>
                <p>Artefatos: <span className="font-medium">{deliveryReport?.artifact_paths.length ?? files.length}</span></p>
              </CardContent>
            </Card>
          </div>

          {!useLinkedProject && (
            <div className="grid gap-3 lg:grid-cols-2">
              <Input placeholder="Nome do produto" value={manualBrief.product_name} onChange={(event) => setManualBrief((prev) => ({ ...prev, product_name: event.target.value }))} />
              <Input placeholder="Tagline (opcional)" value={manualBrief.tagline ?? ""} onChange={(event) => setManualBrief((prev) => ({ ...prev, tagline: event.target.value }))} />
              <Textarea placeholder="Problema resolvido" value={manualBrief.problem_solved} onChange={(event) => setManualBrief((prev) => ({ ...prev, problem_solved: event.target.value }))} className="lg:col-span-2" />
              <Textarea placeholder="Público-alvo" value={manualBrief.target_audience} onChange={(event) => setManualBrief((prev) => ({ ...prev, target_audience: event.target.value }))} className="lg:col-span-2" />
              <Input placeholder="Benefício #1" value={manualBrief.main_benefits[0] ?? ""} onChange={(event) => setManualBrief((prev) => ({ ...prev, main_benefits: [event.target.value, prev.main_benefits[1] ?? "", prev.main_benefits[2] ?? ""] }))} />
              <Input placeholder="Benefício #2" value={manualBrief.main_benefits[1] ?? ""} onChange={(event) => setManualBrief((prev) => ({ ...prev, main_benefits: [prev.main_benefits[0] ?? "", event.target.value, prev.main_benefits[2] ?? ""] }))} />
              <Input placeholder="Benefício #3" value={manualBrief.main_benefits[2] ?? ""} onChange={(event) => setManualBrief((prev) => ({ ...prev, main_benefits: [prev.main_benefits[0] ?? "", prev.main_benefits[1] ?? "", event.target.value] }))} />
              <Input placeholder="CTA primária" value={manualBrief.primary_cta ?? ""} onChange={(event) => setManualBrief((prev) => ({ ...prev, primary_cta: event.target.value }))} />
            </div>
          )}
        </CardContent>
      </Card>

      <Tabs defaultValue="calendar" className="w-full">
        <TabsList className="border border-border bg-card/50 p-1">
          <TabsTrigger value="calendar">Calendário</TabsTrigger>
          <TabsTrigger value="strategy">Estratégia</TabsTrigger>
          <TabsTrigger value="forecast">Performance</TabsTrigger>
          <TabsTrigger value="integrations">Integrações</TabsTrigger>
        </TabsList>

        <TabsContent value="calendar" className="mt-5 space-y-4">
          <Card className="border-border/60 bg-card/70">
            <CardHeader className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <CardTitle className="text-base flex items-center gap-2"><CalendarRange className="h-4 w-4 text-primary" /> Calendário editorial (3 meses)</CardTitle>
                <CardDescription>Filtro por canal com detalhes por item e exportação para CSV/ICS.</CardDescription>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={exportCalendarCsv}>Exportar CSV</Button>
                <Button variant="outline" size="sm" onClick={exportCalendarIcs}>Exportar ICS</Button>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              {monthSeeds.map((monthDate) => {
                const monthLabel = monthDate.toLocaleDateString("pt-BR", { month: "long", year: "numeric" });
                const daysInMonth = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0).getDate();
                return (
                  <div key={monthLabel} className="space-y-2">
                    <h3 className="text-sm font-semibold capitalize">{monthLabel}</h3>
                    <div className="grid gap-2 md:grid-cols-4 xl:grid-cols-7">
                      {Array.from({ length: daysInMonth }, (_, idx) => {
                        const day = idx + 1;
                        const date = `${monthDate.getFullYear()}-${String(monthDate.getMonth() + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
                        const dayEntries = filteredEntries.filter((entry) => entry.date === date);
                        return (
                          <div key={date} className="min-h-[100px] rounded-xl border border-border/60 bg-background/50 p-2">
                            <p className="mb-2 text-xs font-semibold text-muted-foreground">{String(day).padStart(2, "0")}</p>
                            <div className="space-y-1">
                              {dayEntries.slice(0, 3).map((entry) => {
                                const channel = CHANNELS.find((item) => item.value === entry.channel);
                                return (
                                  <button key={`${date}-${entry.channel}-${entry.title}`} onClick={() => setSelectedEntry(entry)} className={`w-full rounded-md border px-2 py-1 text-left text-[11px] font-medium ${channel?.color ?? ""}`}>
                                    {entry.title}
                                  </button>
                                );
                              })}
                              {dayEntries.length > 3 && <p className="text-[10px] text-muted-foreground">+{dayEntries.length - 3} itens</p>}
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                );
              })}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="strategy" className="mt-5">
          <Card className="border-border/60 bg-card/70">
            <CardHeader>
              <CardTitle className="text-base">MARKETING_STRATEGY.md</CardTitle>
              <CardDescription>Documento estratégico gerado pela Tríade (TASK-15-003).</CardDescription>
            </CardHeader>
            <CardContent>
              <ScrollArea className="h-[560px] rounded-lg border border-border/60 bg-background/60 p-4">
                {strategyFile ? <Markdown content={strategyFile.content} /> : <p className="text-sm text-muted-foreground">Execute a phase 15 para gerar a estratégia.</p>}
              </ScrollArea>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="forecast" className="mt-5">
          <div className="grid gap-4 lg:grid-cols-3">
            {(deliveryReport?.channel_summaries ?? []).map((summary) => (
              <Card key={summary.channel} className="border-border/60 bg-card/70">
                <CardHeader>
                  <CardTitle className="text-base capitalize">{summary.channel}</CardTitle>
                  <CardDescription>{summary.pieces} peças | Budget US$ {summary.budget_usd.toFixed(0)}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-2 text-sm">
                  <p>CTR esperado: <strong>{summary.expected_ctr}</strong></p>
                  <p>Conversão esperada: <strong>{summary.expected_conversion}</strong></p>
                </CardContent>
              </Card>
            ))}
          </div>
          <Card className="mt-4 border-border/60 bg-card/70">
            <CardHeader>
              <CardTitle className="text-base">PERFORMANCE_FORECAST.md</CardTitle>
            </CardHeader>
            <CardContent>
              <ScrollArea className="h-[320px] rounded-lg border border-border/60 bg-background/60 p-4">
                {forecastFile ? <Markdown content={forecastFile.content} /> : <p className="text-sm text-muted-foreground">Nenhum forecast disponível.</p>}
              </ScrollArea>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="integrations" className="mt-5">
          <Card className="border-border/60 bg-card/70">
            <CardHeader>
              <CardTitle className="text-base flex items-center gap-2"><Webhook className="h-4 w-4 text-primary" /> Webhook de marketing</CardTitle>
              <CardDescription>Valida URL (200 OK) e persiste configuração real em artifacts/marketing/webhooks/config.json.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex flex-col gap-2 md:flex-row">
                <Input placeholder="https://example.com/webhook" value={webhookUrl} onChange={(event) => setWebhookUrl(event.target.value)} />
                <Button onClick={configureWebhook}>Validar e salvar</Button>
              </div>
              <div className="flex flex-wrap gap-2">
                {channels.map((channel) => (
                  <Badge key={channel} variant="outline">{CHANNELS.find((item) => item.value === channel)?.label}</Badge>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <Sheet open={Boolean(selectedEntry)} onOpenChange={(open) => !open && setSelectedEntry(null)}>
        <SheetContent className="w-full sm:max-w-xl">
          <SheetHeader>
            <SheetTitle>{selectedEntry?.title}</SheetTitle>
            <SheetDescription>Detalhes do post selecionado no calendário.</SheetDescription>
          </SheetHeader>
          {selectedEntry && (
            <div className="mt-6 space-y-3 text-sm">
              <p><strong>Canal:</strong> {selectedEntry.channel}</p>
              <p><strong>Data:</strong> {new Date(`${selectedEntry.date}T00:00:00Z`).toLocaleDateString("pt-BR")}</p>
              <p><strong>Melhor horário:</strong> {selectedEntry.bestTimeUtc} UTC</p>
              <p><strong>Tipo:</strong> {selectedEntry.assetType}</p>
              <p><strong>CTA:</strong> {selectedEntry.cta}</p>
              <div className="rounded-lg border border-border/60 bg-background/50 p-3">
                <p className="mb-1 font-medium">Copy</p>
                <p className="text-muted-foreground">{selectedEntry.title}</p>
              </div>
              <div className="rounded-lg border border-border/60 bg-background/50 p-3">
                <p className="mb-1 font-medium">Hashtags</p>
                <p className="text-muted-foreground">Não disponível neste item de calendário. Use o pack por canal para copiar hashtags completas.</p>
              </div>
            </div>
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}
