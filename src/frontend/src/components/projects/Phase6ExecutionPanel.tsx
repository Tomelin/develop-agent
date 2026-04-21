"use client";

import { useMemo, useState } from "react";
import { ProjectService } from "@/services/project";
import { Phase6AnalyzeCoverageResponse, Phase6ValidationResult } from "@/types/phase6";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AlertTriangle, CheckCircle2, ChevronDown, ChevronUp, FlaskConical, Loader2, ShieldCheck, TestTube2 } from "lucide-react";
import { toast } from "sonner";

interface Phase6ExecutionPanelProps {
  projectId: string;
}

export function Phase6ExecutionPanel({ projectId }: Phase6ExecutionPanelProps) {
  const [backendDir, setBackendDir] = useState("src/backend");
  const [frontendDir, setFrontendDir] = useState("src/frontend");
  const [threshold, setThreshold] = useState("80");

  const [isCoverageLoading, setIsCoverageLoading] = useState(false);
  const [isValidationLoading, setIsValidationLoading] = useState(false);

  const [coverage, setCoverage] = useState<Phase6AnalyzeCoverageResponse | null>(null);
  const [previousCoverage, setPreviousCoverage] = useState<Phase6AnalyzeCoverageResponse | null>(null);
  const [validation, setValidation] = useState<Phase6ValidationResult | null>(null);
  const [expandedFailure, setExpandedFailure] = useState<string | null>(null);

  const thresholdValue = useMemo(() => {
    const parsed = Number(threshold);
    if (Number.isNaN(parsed) || parsed <= 0) return 80;
    return parsed;
  }, [threshold]);

  const runCoverage = async () => {
    try {
      setIsCoverageLoading(true);
      const data = await ProjectService.analyzePhase6Coverage(projectId, {
        backend_dir: backendDir,
        threshold: thresholdValue,
      });
      setPreviousCoverage(coverage);
      setCoverage(data);
      toast.success("Análise de cobertura concluída.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao executar análise de cobertura.");
    } finally {
      setIsCoverageLoading(false);
    }
  };

  const runValidation = async () => {
    try {
      setIsValidationLoading(true);
      const data = await ProjectService.validatePhase6Tests(projectId, {
        backend_dir: backendDir,
        frontend_dir: frontendDir,
      });
      setValidation(data);
      toast.success("Validação de testes concluída.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao validar testes da fase 06.");
    } finally {
      setIsValidationLoading(false);
    }
  };

  const coverageGaugeColor = useMemo(() => {
    if (!coverage) return "#64748b";
    if (coverage.report.total_percent >= 80) return "#10b981";
    if (coverage.report.total_percent >= 60) return "#f59e0b";
    return "#ef4444";
  }, [coverage]);

  const coverageDelta = useMemo(() => {
    if (!coverage || !previousCoverage) return null;
    return coverage.report.total_percent - previousCoverage.report.total_percent;
  }, [coverage, previousCoverage]);

  const failedChecks = useMemo(() => {
    if (!validation?.details) return [];
    return validation.details
      .split("\n")
      .map((line) => line.trim())
      .filter((line) => line.startsWith("FAIL") || line.toLowerCase().includes("error"))
      .slice(0, 8)
      .map((line, index) => ({ id: `${index}-${line}`, summary: line }));
  }, [validation]);

  return (
    <div className="space-y-6">
      <Card className="border-border bg-card/60 backdrop-blur-sm">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-lg">
            <FlaskConical className="h-5 w-5 text-primary" />
            Centro de Execução — Fase 06
          </CardTitle>
          <CardDescription>
            Integração real com o backend para executar análise de cobertura e validação de testes sem mock.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="grid gap-4 md:grid-cols-3">
            <div className="space-y-2">
              <Label htmlFor="backendDir">Diretório backend</Label>
              <Input id="backendDir" value={backendDir} onChange={(e) => setBackendDir(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="frontendDir">Diretório frontend</Label>
              <Input id="frontendDir" value={frontendDir} onChange={(e) => setFrontendDir(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="threshold">Threshold de cobertura (%)</Label>
              <Input id="threshold" type="number" min={1} max={100} value={threshold} onChange={(e) => setThreshold(e.target.value)} />
            </div>
          </div>

          <div className="flex flex-col sm:flex-row gap-3">
            <Button onClick={runCoverage} disabled={isCoverageLoading || isValidationLoading} className="gap-2">
              {isCoverageLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <ShieldCheck className="h-4 w-4" />}
              Executar cobertura
            </Button>
            <Button variant="secondary" onClick={runValidation} disabled={isCoverageLoading || isValidationLoading} className="gap-2">
              {isValidationLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <TestTube2 className="h-4 w-4" />}
              Validar testes
            </Button>
          </div>
        </CardContent>
      </Card>

      {coverage && (
        <Card className="border-border bg-card/50">
          <CardHeader>
            <div className="flex flex-wrap items-center justify-between gap-3">
              <CardTitle className="text-base">Relatório de cobertura</CardTitle>
              <Badge variant="outline" className={coverage.below_threshold ? "text-destructive border-destructive/40" : "text-green-600 border-green-600/40"}>
                {coverage.below_threshold ? "Abaixo do threshold" : "Dentro do threshold"}
              </Badge>
            </div>
            <CardDescription>
              Total de cobertura: <strong>{coverage.report.total_percent.toFixed(2)}%</strong> (meta {coverage.report.threshold_percent.toFixed(2)}%)
            </CardDescription>
            <Progress value={Math.min(coverage.report.total_percent, 100)} className="h-2" />
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-2">
            <div className="space-y-3 rounded-xl border border-border/60 p-4">
              <p className="text-sm font-medium">Pacotes com menor cobertura</p>
              {coverage.report.packages.slice(0, 6).map((pkg) => (
                <div key={pkg.package} className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground line-clamp-1">{pkg.package}</span>
                  <span className="font-mono">{pkg.percent.toFixed(2)}%</span>
                </div>
              ))}
            </div>
            <div className="space-y-3 rounded-xl border border-border/60 p-4">
              <p className="text-sm font-medium">Funções críticas (menor cobertura)</p>
              {coverage.report.functions
                .slice()
                .sort((a, b) => a.percent - b.percent)
                .slice(0, 6)
                .map((fn) => (
                  <div key={fn.source} className="flex items-center justify-between gap-3 text-sm">
                    <span className="text-muted-foreground line-clamp-1">{fn.name}</span>
                    <span className="font-mono">{fn.percent.toFixed(2)}%</span>
                  </div>
                ))}
            </div>
          </CardContent>
        </Card>
      )}

      {validation && (
        <Card className="border-border bg-card/50">
          <CardHeader>
            <CardTitle className="text-base">Validação de testes</CardTitle>
            <CardDescription>Resultado real da execução de `go test`, `go test -race` e testes frontend.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-3 md:grid-cols-3">
              <StatusChip label="Go test" ok={validation.go_test_passed} />
              <StatusChip label="Go race" ok={validation.go_race_passed} />
              <StatusChip label="Frontend" ok={validation.frontend_test_passed} />
            </div>

            <div className="rounded-xl border border-border/60 p-4">
              <p className="text-sm text-muted-foreground">Classificação de falha</p>
              <p className="mt-1 font-medium">{validation.failure_kind}</p>
            </div>

            {validation.details && (
              <div className="rounded-xl border border-border/60 p-4">
                <p className="mb-2 text-sm font-medium">Detalhes da execução</p>
                <ScrollArea className="h-52 rounded-md border border-border/60 bg-background p-3">
                  <pre className="text-xs whitespace-pre-wrap">{validation.details}</pre>
                </ScrollArea>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {(coverage || validation) && (
        <Card className="border-border bg-card/50">
          <CardHeader>
            <CardTitle className="text-base">Dashboard de Testes (Fase 11)</CardTitle>
            <CardDescription>
              Visão consolidada de cobertura, qualidade dos testes e status de CI para tomada de decisão rápida.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {coverage && (
              <div className="grid gap-4 lg:grid-cols-3">
                <div className="rounded-2xl border border-border/60 bg-background/60 p-5 flex items-center gap-4">
                  <div
                    className="relative grid h-24 w-24 place-items-center rounded-full"
                    style={{
                      backgroundImage: `conic-gradient(${coverageGaugeColor} ${Math.min(coverage.report.total_percent, 100) * 3.6}deg, rgba(120,120,120,0.18) 0deg)`,
                    }}
                  >
                    <div className="grid h-16 w-16 place-items-center rounded-full bg-background text-xs font-semibold">
                      {coverage.report.total_percent.toFixed(1)}%
                    </div>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Cobertura total</p>
                    <p className="text-xl font-semibold">{coverage.report.total_percent.toFixed(2)}%</p>
                    <p className="text-xs text-muted-foreground">Meta: {coverage.report.threshold_percent.toFixed(1)}%</p>
                  </div>
                </div>

                <div className="rounded-2xl border border-border/60 bg-background/60 p-5 space-y-2">
                  <p className="text-sm text-muted-foreground">Comparação entre ciclos</p>
                  {coverageDelta === null ? (
                    <p className="text-sm">Execute a cobertura novamente para comparar evolução.</p>
                  ) : (
                    <div className="flex items-center gap-2">
                      {coverageDelta >= 0 ? (
                        <ChevronUp className="h-4 w-4 text-emerald-600" />
                      ) : (
                        <ChevronDown className="h-4 w-4 text-destructive" />
                      )}
                      <p className={`text-lg font-semibold ${coverageDelta >= 0 ? "text-emerald-600" : "text-destructive"}`}>
                        {coverageDelta >= 0 ? "+" : ""}
                        {coverageDelta.toFixed(2)} p.p.
                      </p>
                    </div>
                  )}
                  <p className="text-xs text-muted-foreground">Antes → Depois (última execução na sessão atual)</p>
                </div>

                <div className="rounded-2xl border border-border/60 bg-background/60 p-5 space-y-3">
                  <p className="text-sm text-muted-foreground">Status do CI</p>
                  <div className="flex flex-wrap gap-2">
                    <CIBadge label="Go test" ok={validation?.go_test_passed ?? false} />
                    <CIBadge label="Go race" ok={validation?.go_race_passed ?? false} />
                    <CIBadge label="Frontend" ok={validation?.frontend_test_passed ?? false} />
                  </div>
                </div>
              </div>
            )}

            {coverage && (
              <div className="rounded-2xl border border-border/60 bg-background/60 p-4">
                <p className="mb-4 text-sm font-medium">Cobertura por package/módulo</p>
                <div className="space-y-3">
                  {coverage.report.packages.slice().sort((a, b) => b.percent - a.percent).map((pkg) => (
                    <div key={pkg.package} className="grid gap-2 md:grid-cols-[1fr_auto] md:items-center">
                      <div className="space-y-1">
                        <div className="flex items-center justify-between text-sm">
                          <span className="truncate text-muted-foreground">{pkg.package}</span>
                          <span className="font-mono">{pkg.percent.toFixed(2)}%</span>
                        </div>
                        <Progress value={Math.min(pkg.percent, 100)} className="h-2" />
                      </div>
                      <Badge
                        variant="outline"
                        className={
                          pkg.percent >= 80
                            ? "border-emerald-600/40 text-emerald-600"
                            : pkg.percent >= 60
                              ? "border-yellow-500/40 text-yellow-600"
                              : "border-destructive/40 text-destructive"
                        }
                      >
                        {pkg.percent >= 80 ? "Excelente" : pkg.percent >= 60 ? "Atenção" : "Crítico"}
                      </Badge>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {validation && (
              <div className="rounded-2xl border border-border/60 bg-background/60 p-4">
                <p className="mb-2 text-sm font-medium">Testes com falha</p>
                {failedChecks.length === 0 ? (
                  <p className="text-sm text-muted-foreground">Nenhuma falha detectada na última execução.</p>
                ) : (
                  <div className="space-y-2">
                    {failedChecks.map((item) => {
                      const isExpanded = expandedFailure === item.id;
                      return (
                        <div key={item.id} className="rounded-lg border border-border/60">
                          <button
                            type="button"
                            className="w-full px-3 py-2 text-left text-sm flex items-center justify-between hover:bg-muted/30 transition-colors"
                            onClick={() => setExpandedFailure(isExpanded ? null : item.id)}
                          >
                            <span className="truncate">{item.summary}</span>
                            {isExpanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                          </button>
                          {isExpanded && (
                            <pre className="px-3 pb-3 text-xs whitespace-pre-wrap text-muted-foreground">{validation.details}</pre>
                          )}
                        </div>
                      );
                    })}
                  </div>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}

function StatusChip({ label, ok }: { label: string; ok: boolean }) {
  return (
    <div className="rounded-xl border border-border/60 bg-background p-4 flex items-center justify-between">
      <span className="text-sm">{label}</span>
      <span className="inline-flex items-center gap-1 text-sm">
        {ok ? <CheckCircle2 className="h-4 w-4 text-green-600" /> : <AlertTriangle className="h-4 w-4 text-destructive" />}
        {ok ? "OK" : "Falhou"}
      </span>
    </div>
  );
}

function CIBadge({ label, ok }: { label: string; ok: boolean }) {
  return (
    <Badge
      variant="outline"
      className={ok ? "border-emerald-600/40 text-emerald-600" : "border-destructive/40 text-destructive"}
    >
      {label}: {ok ? "passando" : "falhando"}
    </Badge>
  );
}
