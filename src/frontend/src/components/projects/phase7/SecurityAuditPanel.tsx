"use client";

import { useMemo, useState } from "react";
import { Phase7Service } from "@/services/phase7";
import { ProjectService } from "@/services/project";
import { SecurityAuditReport, SecurityFinding, SecurityFindingStatus, SecuritySeverity } from "@/types/phase7";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AlertTriangle, CheckCircle2, ChevronDown, ChevronUp, Download, Loader2, ShieldAlert, ShieldCheck } from "lucide-react";
import { toast } from "sonner";

interface SecurityAuditPanelProps {
  projectId: string;
}

const SEVERITIES: Array<"ALL" | SecuritySeverity> = ["ALL", "CRITICAL", "HIGH", "MEDIUM", "LOW"];
const STATUSES: Array<"ALL" | SecurityFindingStatus> = ["ALL", "OPEN", "FIXED", "ACCEPTED"];

export function SecurityAuditPanel({ projectId }: SecurityAuditPanelProps) {
  const [backendDir, setBackendDir] = useState("src/backend");
  const [frontendDir, setFrontendDir] = useState("src/frontend");
  const [projectRootDir, setProjectRootDir] = useState(".");
  const [highRetryCount, setHighRetryCount] = useState("2");

  const [loadingAudit, setLoadingAudit] = useState(false);
  const [downloadingReport, setDownloadingReport] = useState(false);
  const [report, setReport] = useState<SecurityAuditReport | null>(null);
  const [severityFilter, setSeverityFilter] = useState<"ALL" | SecuritySeverity>("ALL");
  const [statusFilter, setStatusFilter] = useState<"ALL" | SecurityFindingStatus>("ALL");
  const [query, setQuery] = useState("");
  const [expandedFinding, setExpandedFinding] = useState<string | null>(null);

  const runAudit = async () => {
    try {
      setLoadingAudit(true);
      const payload = {
        backend_dir: backendDir,
        frontend_dir: frontendDir,
        project_root_dir: projectRootDir,
        high_retry_count: Number(highRetryCount) || 2,
      };
      const data = await Phase7Service.runAudit(projectId, payload);
      setReport(data);
      toast.success("Auditoria de segurança concluída com sucesso.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao executar auditoria de segurança.");
    } finally {
      setLoadingAudit(false);
    }
  };

  const downloadSecurityReport = async () => {
    try {
      setDownloadingReport(true);
      const files = await ProjectService.getProjectFiles(projectId);
      const auditFile = files.find((file) => file.path.endsWith("SECURITY_AUDIT.md"));

      if (!auditFile) {
        toast.error("Arquivo SECURITY_AUDIT.md ainda não disponível.");
        return;
      }

      const blob = new Blob([auditFile.content], { type: "text/markdown;charset=utf-8" });
      const url = URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = url;
      anchor.download = "SECURITY_AUDIT.md";
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();
      URL.revokeObjectURL(url);
      toast.success("Relatório SECURITY_AUDIT.md baixado.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao baixar relatório de segurança.");
    } finally {
      setDownloadingReport(false);
    }
  };

  const scoreColor = useMemo(() => {
    if (!report) return "#64748b";
    if (report.summary.score >= 80) return "#10b981";
    if (report.summary.score >= 60) return "#f59e0b";
    return "#ef4444";
  }, [report]);

  const filteredFindings = useMemo(() => {
    if (!report) return [];
    const normalizedQuery = query.trim().toLowerCase();

    return report.findings.filter((finding) => {
      const matchesSeverity = severityFilter === "ALL" ? true : finding.severity === severityFilter;
      const matchesStatus = statusFilter === "ALL" ? true : finding.status === statusFilter;
      const matchesQuery =
        !normalizedQuery ||
        finding.title.toLowerCase().includes(normalizedQuery) ||
        finding.category.toLowerCase().includes(normalizedQuery) ||
        finding.id.toLowerCase().includes(normalizedQuery) ||
        finding.description.toLowerCase().includes(normalizedQuery) ||
        finding.file?.toLowerCase().includes(normalizedQuery);

      return matchesSeverity && matchesStatus && matchesQuery;
    });
  }, [report, severityFilter, statusFilter, query]);

  const severityCardData = useMemo(() => {
    if (!report) {
      return [
        { label: "CRITICAL", count: 0, className: "border-red-500/30 bg-red-500/5 text-red-500" },
        { label: "HIGH", count: 0, className: "border-orange-500/30 bg-orange-500/5 text-orange-500" },
        { label: "MEDIUM", count: 0, className: "border-yellow-500/30 bg-yellow-500/5 text-yellow-500" },
        { label: "LOW", count: 0, className: "border-blue-500/30 bg-blue-500/5 text-blue-500" },
      ];
    }

    return [
      { label: "CRITICAL", count: report.summary.critical_count, className: "border-red-500/30 bg-red-500/5 text-red-500" },
      { label: "HIGH", count: report.summary.high_count, className: "border-orange-500/30 bg-orange-500/5 text-orange-500" },
      { label: "MEDIUM", count: report.summary.medium_count, className: "border-yellow-500/30 bg-yellow-500/5 text-yellow-500" },
      { label: "LOW", count: report.summary.low_count, className: "border-blue-500/30 bg-blue-500/5 text-blue-500" },
    ];
  }, [report]);

  return (
    <div className="space-y-6">
      <Card className="border-border bg-card/60 backdrop-blur-sm">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-lg">
            <ShieldAlert className="h-5 w-5 text-primary" />
            Centro de Auditoria — Fase 12 (Segurança)
          </CardTitle>
          <CardDescription>
            Execute a auditoria OWASP completa com dados reais do backend e visualize findings com filtros avançados.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <div className="space-y-2">
              <Label htmlFor="phase7BackendDir">Diretório backend</Label>
              <Input id="phase7BackendDir" value={backendDir} onChange={(e) => setBackendDir(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="phase7FrontendDir">Diretório frontend</Label>
              <Input id="phase7FrontendDir" value={frontendDir} onChange={(e) => setFrontendDir(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="phase7RootDir">Diretório raiz do projeto</Label>
              <Input id="phase7RootDir" value={projectRootDir} onChange={(e) => setProjectRootDir(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="phase7RetryCount">Tentativas para HIGH (0-2)</Label>
              <Input id="phase7RetryCount" type="number" min={0} max={2} value={highRetryCount} onChange={(e) => setHighRetryCount(e.target.value)} />
            </div>
          </div>

          <div className="flex flex-col gap-3 sm:flex-row">
            <Button onClick={runAudit} disabled={loadingAudit || downloadingReport} className="gap-2">
              {loadingAudit ? <Loader2 className="h-4 w-4 animate-spin" /> : <ShieldCheck className="h-4 w-4" />}
              Executar auditoria de segurança
            </Button>
            <Button variant="secondary" onClick={downloadSecurityReport} disabled={loadingAudit || downloadingReport} className="gap-2">
              {downloadingReport ? <Loader2 className="h-4 w-4 animate-spin" /> : <Download className="h-4 w-4" />}
              Baixar SECURITY_AUDIT.md
            </Button>
          </div>
        </CardContent>
      </Card>

      {report && (
        <>
          <div className="grid gap-4 lg:grid-cols-3">
            <Card className="border-border bg-card/50 lg:col-span-1">
              <CardHeader>
                <CardTitle className="text-base">Security Score</CardTitle>
                <CardDescription>Status geral da auditoria: {report.status}</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-4">
                  <div
                    className="relative grid h-28 w-28 place-items-center rounded-full"
                    style={{ backgroundImage: `conic-gradient(${scoreColor} ${Math.min(report.summary.score, 100) * 3.6}deg, rgba(120,120,120,0.18) 0deg)` }}
                  >
                    <div className="grid h-20 w-20 place-items-center rounded-full bg-background text-lg font-semibold">{report.summary.score}</div>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Pontuação (0-100)</p>
                    <p className="text-3xl font-bold">{report.summary.score}</p>
                    <p className="text-xs text-muted-foreground">
                      {report.summary.score >= 80 ? "Excelente" : report.summary.score >= 60 ? "Atenção" : "Crítico"}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="border-border bg-card/50 lg:col-span-2">
              <CardHeader>
                <CardTitle className="text-base">Resumo por severidade</CardTitle>
                <CardDescription>{report.summary.total_findings} findings identificados na auditoria.</CardDescription>
              </CardHeader>
              <CardContent className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
                {severityCardData.map((item) => (
                  <div key={item.label} className={`rounded-xl border p-4 ${item.className}`}>
                    <p className="text-xs font-semibold tracking-wide">{item.label}</p>
                    <p className="mt-2 text-3xl font-bold">{item.count}</p>
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>

          <Card className="border-border bg-card/50">
            <CardHeader>
              <CardTitle className="text-base">Tabela de findings</CardTitle>
              <CardDescription>
                Filtre por severidade/status, expanda para ver detalhes técnicos e acompanhe as correções do Refinador.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-3 md:grid-cols-[1fr_auto_auto]">
                <Input
                  placeholder="Buscar por ID, título, categoria, arquivo ou descrição..."
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                />
                <div className="flex gap-2 overflow-x-auto pb-1">
                  {SEVERITIES.map((severity) => (
                    <Button
                      key={severity}
                      type="button"
                      size="sm"
                      variant={severityFilter === severity ? "default" : "outline"}
                      onClick={() => setSeverityFilter(severity)}
                    >
                      {severity}
                    </Button>
                  ))}
                </div>
                <div className="flex gap-2 overflow-x-auto pb-1">
                  {STATUSES.map((status) => (
                    <Button
                      key={status}
                      type="button"
                      size="sm"
                      variant={statusFilter === status ? "default" : "outline"}
                      onClick={() => setStatusFilter(status)}
                    >
                      {status}
                    </Button>
                  ))}
                </div>
              </div>

              <ScrollArea className="h-[520px] rounded-xl border border-border/60 bg-background/30">
                <div className="divide-y divide-border/50">
                  {filteredFindings.map((finding) => {
                    const isExpanded = expandedFinding === finding.id;
                    return (
                      <div key={finding.id} className="p-4">
                        <button
                          type="button"
                          className="flex w-full items-center justify-between gap-4 text-left"
                          onClick={() => setExpandedFinding(isExpanded ? null : finding.id)}
                        >
                          <div className="min-w-0 space-y-2">
                            <div className="flex flex-wrap items-center gap-2">
                              <Badge variant="outline">{finding.id}</Badge>
                              <SeverityBadge severity={finding.severity} />
                              <Badge variant="secondary">CVSS {finding.cvss.toFixed(1)}</Badge>
                              {finding.status === "FIXED" && (
                                <Badge className="bg-emerald-600 text-white hover:bg-emerald-600">CORRIGIDO</Badge>
                              )}
                            </div>
                            <p className="font-medium leading-tight">{finding.title}</p>
                            <p className="text-sm text-muted-foreground">
                              {finding.category} • {finding.detected_by} • {finding.file ? `${finding.file}:${finding.line ?? "-"}` : "Sem arquivo vinculado"}
                            </p>
                          </div>
                          <div className="shrink-0">{isExpanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}</div>
                        </button>

                        {isExpanded && <ExpandedFindingDetails finding={finding} />}
                      </div>
                    );
                  })}

                  {!filteredFindings.length && (
                    <div className="p-8 text-center text-muted-foreground">
                      Nenhum finding encontrado com os filtros atuais.
                    </div>
                  )}
                </div>
              </ScrollArea>
            </CardContent>
          </Card>

          <Card className="border-border bg-card/50">
            <CardHeader>
              <CardTitle className="text-base">Auto Rejeição (Fase 7)</CardTitle>
              <CardDescription>Rastreie se o gatilho foi ativado para CVSS crítico, secrets expostos ou falhas estruturais.</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-3">
              <StatusTile
                label="Gatilho"
                value={report.auto_rejection.triggered ? "ATIVADO" : "NÃO ATIVADO"}
                tone={report.auto_rejection.triggered ? "destructive" : "success"}
              />
              <StatusTile
                label="Retorno para Fase 5"
                value={report.auto_rejection.returned_phase5 ? "SIM" : "NÃO"}
                tone={report.auto_rejection.returned_phase5 ? "destructive" : "neutral"}
              />
              <StatusTile
                label="Tentativas para HIGH"
                value={`${report.auto_rejection.retry_count}`}
                tone="neutral"
              />
              {report.auto_rejection.reason && (
                <div className="md:col-span-3 rounded-xl border border-destructive/40 bg-destructive/5 p-4 text-sm">
                  <p className="font-semibold text-destructive">Motivo</p>
                  <p className="mt-1 text-muted-foreground">{report.auto_rejection.reason}</p>
                  {report.auto_rejection.findings?.length ? (
                    <ul className="mt-3 list-disc space-y-1 pl-5 text-muted-foreground">
                      {report.auto_rejection.findings.map((item) => (
                        <li key={item}>{item}</li>
                      ))}
                    </ul>
                  ) : null}
                </div>
              )}
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}

function SeverityBadge({ severity }: { severity: SecuritySeverity }) {
  const style =
    severity === "CRITICAL"
      ? "border-red-500/40 bg-red-500/10 text-red-500"
      : severity === "HIGH"
        ? "border-orange-500/40 bg-orange-500/10 text-orange-500"
        : severity === "MEDIUM"
          ? "border-yellow-500/40 bg-yellow-500/10 text-yellow-500"
          : "border-blue-500/40 bg-blue-500/10 text-blue-500";

  return (
    <Badge variant="outline" className={style}>
      {severity}
    </Badge>
  );
}

function StatusTile({ label, value, tone }: { label: string; value: string; tone: "success" | "destructive" | "neutral" }) {
  const toneClass =
    tone === "success"
      ? "border-emerald-600/40 bg-emerald-600/10"
      : tone === "destructive"
        ? "border-destructive/40 bg-destructive/10"
        : "border-border/60 bg-background/40";

  return (
    <div className={`rounded-xl border p-4 ${toneClass}`}>
      <p className="text-xs text-muted-foreground">{label}</p>
      <p className="mt-2 text-lg font-semibold">{value}</p>
    </div>
  );
}

function ExpandedFindingDetails({ finding }: { finding: SecurityFinding }) {
  return (
    <div className="mt-4 space-y-3 rounded-xl border border-border/60 bg-background/30 p-4 text-sm">
      <div className="grid gap-3 md:grid-cols-2">
        <Detail label="Status" value={finding.status} />
        <Detail label="CVE" value={finding.cve || "N/A"} />
        <Detail label="Arquivo" value={finding.file || "N/A"} />
        <Detail label="Linha" value={finding.line?.toString() || "N/A"} />
      </div>

      <Block title="Descrição técnica" content={finding.description} />
      <Block title="Prova de conceito (POC)" content={finding.poc || "Não informado"} />
      <Block title="Remediação" content={finding.remediation || "Não informado"} />

      <div className="flex items-center gap-2 text-xs text-muted-foreground">
        {finding.status === "FIXED" ? (
          <CheckCircle2 className="h-3.5 w-3.5 text-emerald-600" />
        ) : (
          <AlertTriangle className="h-3.5 w-3.5 text-yellow-500" />
        )}
        Detectado por: {finding.detected_by}
      </div>
    </div>
  );
}

function Detail({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg border border-border/60 bg-background/40 p-3">
      <p className="text-xs text-muted-foreground">{label}</p>
      <p className="mt-1 font-medium">{value}</p>
    </div>
  );
}

function Block({ title, content }: { title: string; content: string }) {
  return (
    <div className="space-y-1">
      <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">{title}</p>
      <p className="whitespace-pre-wrap leading-relaxed text-muted-foreground">{content}</p>
    </div>
  );
}
