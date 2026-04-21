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
import { AlertTriangle, CheckCircle2, FlaskConical, Loader2, ShieldCheck, TestTube2 } from "lucide-react";
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
  const [validation, setValidation] = useState<Phase6ValidationResult | null>(null);

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
