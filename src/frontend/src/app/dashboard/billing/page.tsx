"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { BillingService, BillingQueryParams } from "@/services/billing";
import { BillingGroupedItem, BillingPricingTable, BillingRecord, BillingSummary } from "@/types/billing";
import { ProjectService } from "@/services/project";
import { Project } from "@/types/project";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Progress } from "@/components/ui/progress";
import {
  AlertCircle,
  BarChart3,
  CalendarRange,
  Download,
  DollarSign,
  Gauge,
  Loader2,
  PieChart,
  TrendingUp,
} from "lucide-react";

const PERIODS = [
  { label: "7 dias", value: "7d" },
  { label: "30 dias", value: "30d" },
  { label: "90 dias", value: "90d" },
  { label: "12 meses", value: "365d" },
  { label: "Personalizado", value: "custom" },
] as const;

type PeriodValue = (typeof PERIODS)[number]["value"];

const formatCurrency = (value: number) =>
  new Intl.NumberFormat("pt-BR", { style: "currency", currency: "USD", minimumFractionDigits: 2 }).format(value || 0);

const formatNumber = (value: number) => new Intl.NumberFormat("pt-BR").format(value || 0);

const toDateInput = (date: Date) => date.toISOString().split("T")[0];

const buildDateRange = (period: PeriodValue, customFrom: string, customTo: string): BillingQueryParams => {
  if (period === "custom" && customFrom && customTo) {
    return {
      from: new Date(`${customFrom}T00:00:00Z`).toISOString(),
      to: new Date(`${customTo}T23:59:59Z`).toISOString(),
    };
  }

  const now = new Date();
  const from = new Date(now);

  switch (period) {
    case "7d":
      from.setDate(now.getDate() - 7);
      break;
    case "30d":
      from.setDate(now.getDate() - 30);
      break;
    case "90d":
      from.setDate(now.getDate() - 90);
      break;
    case "365d":
      from.setDate(now.getDate() - 365);
      break;
    default:
      from.setDate(now.getDate() - 30);
  }

  return {
    from: from.toISOString(),
    to: now.toISOString(),
  };
};

const maxBy = (values: number[]) => Math.max(...values, 0.00001);

const DEFAULT_CUSTOM_FROM = toDateInput(new Date(Date.now() - 30 * 24 * 3600 * 1000));
const DEFAULT_CUSTOM_TO = toDateInput(new Date());

const BarStack = ({ data }: { data: BillingGroupedItem[] }) => {
  const max = maxBy(data.map((item) => item.cost_usd));
  return (
    <div className="space-y-3">
      {data.slice(0, 8).map((item) => (
        <div key={item.key} className="space-y-1">
          <div className="flex items-center justify-between text-xs">
            <span className="text-muted-foreground truncate max-w-[70%]">{item.key}</span>
            <span className="font-medium">{formatCurrency(item.cost_usd)}</span>
          </div>
          <div className="h-2 w-full rounded-full bg-muted">
            <div
              className="h-full rounded-full bg-primary"
              style={{ width: `${Math.max((item.cost_usd / max) * 100, 3)}%` }}
            />
          </div>
        </div>
      ))}
    </div>
  );
};

const SERIES_COLORS = ["#3b82f6", "#8b5cf6", "#f59e0b", "#10b981", "#ef4444", "#06b6d4"]

const DonutLegend = ({ data }: { data: BillingGroupedItem[] }) => {
  const total = data.reduce((acc, item) => acc + item.cost_usd, 0) || 1;
  return (
    <div className="grid gap-2">
      {data.slice(0, 6).map((item, index) => {
        const percentage = (item.cost_usd / total) * 100;
        return (
          <div key={item.key} className="flex items-center justify-between text-xs rounded-lg border p-2">
            <div className="flex items-center gap-2">
              <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: SERIES_COLORS[index % SERIES_COLORS.length] }} />
              <span className="font-medium">{item.key}</span>
            </div>
            <span className="text-muted-foreground">{percentage.toFixed(1)}%</span>
          </div>
        );
      })}
    </div>
  );
};

export default function BillingDashboardPage() {
  const [period, setPeriod] = useState<PeriodValue>("30d");
  const [customFrom, setCustomFrom] = useState(DEFAULT_CUSTOM_FROM);
  const [customTo, setCustomTo] = useState(DEFAULT_CUSTOM_TO);
  const [selectedProject, setSelectedProject] = useState<string>("ALL");
  const [exportFormat, setExportFormat] = useState<"csv" | "json">("csv");
  const [budgetUsd, setBudgetUsd] = useState<number>(1000);

  const [loading, setLoading] = useState(true);
  const [summary, setSummary] = useState<BillingSummary | null>(null);
  const [byModel, setByModel] = useState<BillingGroupedItem[]>([]);
  const [byPhase, setByPhase] = useState<BillingGroupedItem[]>([]);
  const [topProjects, setTopProjects] = useState<BillingGroupedItem[]>([]);
  const [records, setRecords] = useState<BillingRecord[]>([]);
  const [pricing, setPricing] = useState<BillingPricingTable | null>(null);
  const [projects, setProjects] = useState<Project[]>([]);

  const range = useMemo(() => buildDateRange(period, customFrom, customTo), [period, customFrom, customTo]);

  const expensiveModel = useMemo(() => byModel[0]?.key || "-", [byModel]);
  const avgPerProject = useMemo(() => {
    if (!summary?.by_project?.length) return 0;
    return summary.total_cost_usd / summary.by_project.length;
  }, [summary]);

  const selectedProjectCost = useMemo(() => {
    if (selectedProject === "ALL") return summary?.total_cost_usd || 0;
    return topProjects.find((item) => item.key === selectedProject)?.cost_usd || 0;
  }, [selectedProject, summary, topProjects]);

  const budgetProgress = useMemo(() => {
    if (!budgetUsd) return 0;
    return Math.min((selectedProjectCost / budgetUsd) * 100, 100);
  }, [selectedProjectCost, budgetUsd]);

  const loadDashboard = useCallback(async () => {
    setLoading(true);
    try {
      const [summaryData, byModelData, byPhaseData, topProjectsData, pricingData, recordsData, projectsData] = await Promise.all([
        BillingService.getSummary(range),
        BillingService.getByModel(range),
        BillingService.getByPhase(range),
        BillingService.getTopProjects(range),
        BillingService.getPricing(),
        BillingService.getRecords({ ...range, limit: 200 }),
        ProjectService.getProjects(1, 100),
      ]);

      setSummary(summaryData);
      setByModel(byModelData.sort((a, b) => b.cost_usd - a.cost_usd));
      setByPhase(byPhaseData.sort((a, b) => b.cost_usd - a.cost_usd));
      setTopProjects(topProjectsData.sort((a, b) => b.cost_usd - a.cost_usd));
      setPricing(pricingData);
      setRecords(recordsData);
      setProjects(projectsData.items || []);
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível carregar o painel de billing.");
    } finally {
      setLoading(false);
    }
  }, [range]);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      void loadDashboard();
    }, 0);

    return () => window.clearTimeout(timer);
  }, [loadDashboard]);

  const handleExport = async () => {
    try {
      const blob = await BillingService.exportBilling(exportFormat, {
        ...range,
        project_id: selectedProject === "ALL" ? undefined : selectedProject,
      });
      const url = URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = url;
      anchor.download = `billing-${new Date().toISOString().split("T")[0]}.${exportFormat}`;
      anchor.click();
      URL.revokeObjectURL(url);
      toast.success(`Relatório ${exportFormat.toUpperCase()} exportado com sucesso.`);
    } catch (error) {
      console.error(error);
      toast.error("Falha ao exportar relatório.");
    }
  };

  const handleBudgetUpdate = async () => {
    if (selectedProject === "ALL") {
      toast.info("Selecione um projeto específico para configurar budget.");
      return;
    }

    try {
      await BillingService.updateProjectBudget(selectedProject, budgetUsd);
      toast.success("Budget atualizado com sucesso.");
    } catch (error) {
      console.error(error);
      toast.warning("Endpoint de budget indisponível no backend atual.");
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Billing & Auditoria de Custos</h1>
          <p className="text-muted-foreground mt-1">Transparência completa de consumo LLM por projeto, fase, agente e modelo.</p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Select value={period} onValueChange={(value) => setPeriod(value as PeriodValue)}>
            <SelectTrigger className="w-[170px]">
              <SelectValue placeholder="Período" />
            </SelectTrigger>
            <SelectContent>
              {PERIODS.map((item) => (
                <SelectItem key={item.value} value={item.value}>
                  {item.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          {period === "custom" && (
            <>
              <Input type="date" value={customFrom} onChange={(event) => setCustomFrom(event.target.value)} className="w-[150px]" />
              <Input type="date" value={customTo} onChange={(event) => setCustomTo(event.target.value)} className="w-[150px]" />
            </>
          )}

          <Button variant="outline" onClick={loadDashboard} className="gap-2">
            <CalendarRange className="h-4 w-4" /> Atualizar
          </Button>
        </div>
      </div>

      {loading ? (
        <Card>
          <CardContent className="py-20 flex items-center justify-center text-muted-foreground gap-2">
            <Loader2 className="h-5 w-5 animate-spin" /> Carregando métricas de billing...
          </CardContent>
        </Card>
      ) : (
        <>
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            <Card>
              <CardHeader className="pb-3">
                <CardDescription>Custo total no período</CardDescription>
                <CardTitle className="text-2xl">{formatCurrency(summary?.total_cost_usd || 0)}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <DollarSign className="h-3.5 w-3.5" /> {formatNumber(summary?.total_tokens || 0)} tokens processados
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-3">
                <CardDescription>Custo médio por projeto</CardDescription>
                <CardTitle className="text-2xl">{formatCurrency(avgPerProject)}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-xs text-muted-foreground">{summary?.by_project?.length || 0} projetos com consumo registrado</div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-3">
                <CardDescription>Modelo de maior custo</CardDescription>
                <CardTitle className="text-2xl truncate">{expensiveModel}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-xs text-muted-foreground">Baseado no consolidado por modelo do período</div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-3">
                <CardDescription>Auto-rejections</CardDescription>
                <CardTitle className="text-2xl">{records.filter((record) => record.is_auto_rejection).length}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-xs text-muted-foreground">Execuções acionadas por rejeição automática</div>
              </CardContent>
            </Card>
          </div>

          <div className="grid gap-4 xl:grid-cols-3">
            <Card className="xl:col-span-2">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="h-5 w-5 text-primary" /> Custo por fase do pipeline
                </CardTitle>
                <CardDescription>Identifique as fases com maior concentração de custos LLM.</CardDescription>
              </CardHeader>
              <CardContent>
                <BarStack data={byPhase} />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="h-5 w-5 text-primary" /> Distribuição por modelo
                </CardTitle>
                <CardDescription>Visão de representatividade financeira por modelo/provider.</CardDescription>
              </CardHeader>
              <CardContent>
                <DonutLegend data={byModel} />
              </CardContent>
            </Card>
          </div>

          <Tabs defaultValue="records">
            <TabsList className="grid w-full grid-cols-4">
              <TabsTrigger value="records">Registros</TabsTrigger>
              <TabsTrigger value="models">Análise de Modelos</TabsTrigger>
              <TabsTrigger value="budget">Budget</TabsTrigger>
              <TabsTrigger value="pricing">Tabela de Preços</TabsTrigger>
            </TabsList>

            <TabsContent value="records" className="space-y-4">
              <Card>
                <CardHeader className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                  <div>
                    <CardTitle>Tabela detalhada de billing</CardTitle>
                    <CardDescription>Auditoria detalhada por execução da Tríade.</CardDescription>
                  </div>
                  <div className="flex items-center gap-2">
                    <Select value={selectedProject} onValueChange={(value: string | null) => value && setSelectedProject(value)}>
                      <SelectTrigger className="w-[220px]">
                        <SelectValue placeholder="Filtrar projeto" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="ALL">Todos os projetos</SelectItem>
                        {projects.map((project) => (
                          <SelectItem key={project.id} value={project.id}>
                            {project.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>

                    <Select value={exportFormat} onValueChange={(value) => setExportFormat(value as "csv" | "json")}>
                      <SelectTrigger className="w-[110px]">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="csv">CSV</SelectItem>
                        <SelectItem value="json">JSON</SelectItem>
                      </SelectContent>
                    </Select>

                    <Button onClick={handleExport} className="gap-2">
                      <Download className="h-4 w-4" /> Exportar
                    </Button>
                  </div>
                </CardHeader>
                <CardContent className="overflow-x-auto">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Data</TableHead>
                        <TableHead>Fase</TableHead>
                        <TableHead>Agente</TableHead>
                        <TableHead>Modelo</TableHead>
                        <TableHead>Tokens</TableHead>
                        <TableHead>Custo</TableHead>
                        <TableHead>Tipo</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {records
                        .filter((item) => selectedProject === "ALL" || item.project_id === selectedProject)
                        .slice(0, 120)
                        .map((record) => (
                          <TableRow key={record.id}>
                            <TableCell>{new Date(record.timestamp).toLocaleString("pt-BR")}</TableCell>
                            <TableCell>
                              <div className="text-sm">Fase {record.phase_number}</div>
                              <div className="text-xs text-muted-foreground">{record.phase_name}</div>
                            </TableCell>
                            <TableCell>
                              <div className="text-sm">{record.agent_name || record.agent_id}</div>
                              <div className="text-xs text-muted-foreground">{record.triad_role}</div>
                            </TableCell>
                            <TableCell>
                              <div className="text-sm">{record.model}</div>
                              <div className="text-xs text-muted-foreground">{record.provider}</div>
                            </TableCell>
                            <TableCell className="text-sm">{formatNumber(record.total_tokens)}</TableCell>
                            <TableCell className="font-medium">{formatCurrency(record.estimated_cost_usd)}</TableCell>
                            <TableCell>
                              {record.is_auto_rejection ? (
                                <Badge variant="destructive" className="gap-1">
                                  <AlertCircle className="h-3 w-3" /> Auto-Reject
                                </Badge>
                              ) : (
                                <Badge variant="secondary">Normal</Badge>
                              )}
                            </TableCell>
                          </TableRow>
                        ))}
                    </TableBody>
                  </Table>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="models">
              <div className="grid gap-4 xl:grid-cols-2">
                <Card>
                  <CardHeader>
                    <CardTitle>Comparativo de custo por modelo</CardTitle>
                    <CardDescription>Ranking de eficiência financeira por volume e custo total.</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>Modelo</TableHead>
                          <TableHead>Execuções</TableHead>
                          <TableHead>Tokens</TableHead>
                          <TableHead>Custo Total</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {byModel.map((item) => (
                          <TableRow key={item.key}>
                            <TableCell className="font-medium">{item.key}</TableCell>
                            <TableCell>{item.executions}</TableCell>
                            <TableCell>{formatNumber(item.tokens)}</TableCell>
                            <TableCell>{formatCurrency(item.cost_usd)}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <TrendingUp className="h-5 w-5 text-primary" /> Recomendação orientada a histórico
                    </CardTitle>
                    <CardDescription>Modelo sugerido com melhor relação entre custo e volume de execuções.</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {byModel.length > 0 ? (
                      <>
                        <div className="rounded-xl border bg-muted/30 p-4">
                          <div className="text-sm text-muted-foreground">Modelo recomendado</div>
                          <div className="text-xl font-semibold mt-1">{byModel[byModel.length - 1].key}</div>
                          <div className="text-xs text-muted-foreground mt-2">
                            Menor custo agregado no período filtrado mantendo histórico ativo de execuções.
                          </div>
                        </div>
                        <BarStack data={byModel} />
                      </>
                    ) : (
                      <div className="text-sm text-muted-foreground">Sem dados suficientes para gerar recomendação.</div>
                    )}
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            <TabsContent value="budget">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Gauge className="h-5 w-5 text-primary" /> Alertas de budget por projeto
                  </CardTitle>
                  <CardDescription>Configure um limite de gasto e monitore o consumo em tempo real.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-5">
                  <div className="grid gap-4 lg:grid-cols-[1fr_240px_140px]">
                    <Select value={selectedProject} onValueChange={(value: string | null) => value && setSelectedProject(value)}>
                      <SelectTrigger>
                        <SelectValue placeholder="Escolha um projeto" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="ALL">Selecione um projeto</SelectItem>
                        {projects.map((project) => (
                          <SelectItem key={project.id} value={project.id}>
                            {project.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <Input
                      type="number"
                      min={0}
                      step={50}
                      value={budgetUsd}
                      onChange={(event) => setBudgetUsd(Number(event.target.value || 0))}
                      placeholder="Budget em USD"
                    />
                    <Button onClick={handleBudgetUpdate}>Salvar Budget</Button>
                  </div>

                  <div className="rounded-xl border p-4 space-y-3">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">Consumo atual</span>
                      <span className="font-semibold">
                        {formatCurrency(selectedProjectCost)} / {formatCurrency(budgetUsd)}
                      </span>
                    </div>
                    <Progress value={budgetProgress} className="h-2" />
                    <div className="flex items-center justify-between text-xs text-muted-foreground">
                      <span>Alerta em 80%</span>
                      <span>Pausa automática em 100%</span>
                    </div>
                    {budgetProgress >= 80 && (
                      <div className="rounded-lg border border-yellow-500/30 bg-yellow-500/10 p-3 text-sm text-yellow-700 dark:text-yellow-300">
                        Atenção: projeto próximo do limite. Reavalie execução antes de continuar.
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="pricing">
              <Card>
                <CardHeader>
                  <CardTitle>Tabela de preços por provider/modelo</CardTitle>
                  <CardDescription>
                    Snapshot da tabela vigente no backend. Última atualização: {pricing?.last_updated || "n/d"}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Provider</TableHead>
                        <TableHead>Modelo</TableHead>
                        <TableHead>Prompt / 1M</TableHead>
                        <TableHead>Completion / 1M</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {pricing?.models?.map((item) => (
                        <TableRow key={`${item.provider}-${item.model}`}>
                          <TableCell>{item.provider}</TableCell>
                          <TableCell className="font-medium">{item.model}</TableCell>
                          <TableCell>{formatCurrency(item.prompt_price_per_million_tokens)}</TableCell>
                          <TableCell>{formatCurrency(item.completion_price_per_million_tokens)}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </>
      )}
    </div>
  );
}
