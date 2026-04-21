"use client";

import { useCallback, useEffect, useMemo, useState, type ReactNode } from "react";
import { ProjectService } from "@/services/project";
import { Phase13Service } from "@/services/phase13";
import { Phase5CodeFile } from "@/types/phase5";
import { Project } from "@/types/project";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Markdown } from "@/components/ui/markdown";
import { toast } from "sonner";
import Link from "next/link";
import {
  Archive,
  CheckCircle2,
  Download,
  FileCode2,
  FileText,
  Gift,
  Loader2,
  Rocket,
  Server,
  Sparkles,
} from "lucide-react";

const DOC_FILES = ["README.md", "docs/API_REFERENCE.md", "docs/OPERATIONS.md", "CONTRIBUTING.md", "PROJECT_SUMMARY.md"];
const INFRA_FILES = [
  "infra/Dockerfile.backend",
  "infra/Dockerfile.frontend",
  "infra/docker-compose.yml",
  ".github/workflows/ci.yml",
  ".github/workflows/cd.yml",
  "k8s/backend-deployment.yaml",
  "k8s/frontend-deployment.yaml",
  "k8s/services.yaml",
  "k8s/hpa.yaml",
  "k8s/pdb.yaml",
  "k8s/configmap.yaml",
  "k8s/secret.template.yaml",
  "k8s/ingress.yaml",
];

interface Phase13DeliveryCenterProps {
  project: Project;
}

export function Phase13DeliveryCenter({ project }: Phase13DeliveryCenterProps) {
  const [files, setFiles] = useState<Phase5CodeFile[]>([]);
  const [selectedDoc, setSelectedDoc] = useState<string>("README.md");
  const [selectedInfra, setSelectedInfra] = useState<string>("infra/docker-compose.yml");
  const [running, setRunning] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [showCompletionModal, setShowCompletionModal] = useState(false);

  const fetchFiles = useCallback(async () => {
    try {
      const data = await ProjectService.getProjectFiles(project.id);
      setFiles(data.filter((item) => item.phase_number >= 8));
    } catch (error) {
      console.error(error);
      toast.error("Falha ao carregar artefatos da entrega.");
    }
  }, [project.id]);

  useEffect(() => {
    const bootstrap = setTimeout(() => {
      void fetchFiles();
    }, 0);

    const interval = setInterval(() => void fetchFiles(), 12000);
    return () => {
      clearTimeout(bootstrap);
      clearInterval(interval);
    };
  }, [fetchFiles]);

  const docs = useMemo(() => files.filter((file) => DOC_FILES.includes(file.path) || file.path.startsWith("docs/")), [files]);
  const infra = useMemo(() => files.filter((file) => INFRA_FILES.includes(file.path) || file.path.startsWith("infra/") || file.path.startsWith("k8s/") || file.path.startsWith(".github/workflows/")), [files]);

  const docCompletion = Math.round((DOC_FILES.filter((path) => files.some((file) => file.path === path)).length / DOC_FILES.length) * 100);
  const infraCompletion = Math.round((INFRA_FILES.filter((path) => files.some((file) => file.path === path)).length / INFRA_FILES.length) * 100);

  const selectedDocFile = docs.find((file) => file.path === selectedDoc) ?? docs[0] ?? null;
  const selectedInfraFile = infra.find((file) => file.path === selectedInfra) ?? infra[0] ?? null;

  const summaryMetrics = useMemo(() => {
    const summary = files.find((file) => file.path === "PROJECT_SUMMARY.md")?.content ?? "";
    return {
      phases: summary.match(/Fases concluídas\*\*: ([^\n]+)/)?.[1] ?? "-",
      generatedFiles: summary.match(/Arquivos gerados\*\*: ([^\n]+)/)?.[1] ?? "-",
      estimatedCost: summary.match(/Custo estimado\*\*: ([^\n]+)/)?.[1] ?? "-",
      coverage: summary.match(/Cobertura de testes\*\*: ([^\n]+)/)?.[1] ?? "-",
      security: summary.match(/Score de segurança\*\*: ([^\n]+)/)?.[1] ?? "-",
    };
  }, [files]);

  useEffect(() => {
    const hasSummary = files.some((file) => file.path === "PROJECT_SUMMARY.md");
    if (project.status !== "COMPLETED" || !hasSummary) return;

    const storageKey = `project-complete-modal-${project.id}`;
    if (typeof window !== "undefined" && localStorage.getItem(storageKey) !== "shown") {
      const timer = setTimeout(() => {
        setShowCompletionModal(true);
      }, 0);
      localStorage.setItem(storageKey, "shown");
      return () => clearTimeout(timer);
    }
  }, [files, project.id, project.status]);

  const runPhase13 = async () => {
    try {
      setRunning(true);
      const report = await Phase13Service.run(project.id, {
        include_devops: true,
        backend_base_url: process.env.NEXT_PUBLIC_API_URL?.replace("/api/v1", "") ?? "http://localhost:8080",
        frontend_url: typeof window !== "undefined" ? window.location.origin : "http://localhost:3000",
      });
      toast.success(`Entrega gerada com ${report.artifacts.length} artefatos.`);
      await fetchFiles();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao executar a Fase 13.");
    } finally {
      setRunning(false);
    }
  };

  const downloadZip = async () => {
    try {
      setDownloading(true);
      const blob = await ProjectService.downloadProjectFilesZip(project.id);
      const anchor = document.createElement("a");
      anchor.href = URL.createObjectURL(blob);
      anchor.download = `project-${project.id}-delivery.zip`;
      document.body.appendChild(anchor);
      anchor.click();
      URL.revokeObjectURL(anchor.href);
      anchor.remove();
      toast.success("Download iniciado.");
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível baixar os artefatos.");
    } finally {
      setDownloading(false);
    }
  };

  return (
    <div className="space-y-6">
      <Card className="border-border/60 bg-card/70 backdrop-blur-sm">
        <CardHeader className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Gift className="h-5 w-5 text-primary" />
              Entrega Final — Documentação & DevOps
            </CardTitle>
            <CardDescription>
              Visualização premium dos artefatos da Fase 13 com completude, leitura técnica e exportação em lote.
            </CardDescription>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" onClick={downloadZip} disabled={downloading} className="gap-2">
              {downloading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Archive className="h-4 w-4" />} Download ZIP
            </Button>
            <Button onClick={runPhase13} disabled={running} className="gap-2">
              {running ? <Loader2 className="h-4 w-4 animate-spin" /> : <Rocket className="h-4 w-4" />} Gerar/Atualizar Entrega
            </Button>
          </div>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          <MetricCard title="Documentação" value={`${docCompletion}%`} icon={<FileText className="h-4 w-4 text-primary" />} />
          <MetricCard title="Infraestrutura" value={`${infraCompletion}%`} icon={<Server className="h-4 w-4 text-primary" />} />
          <MetricCard title="Artefatos" value={`${files.length}`} icon={<FileCode2 className="h-4 w-4 text-primary" />} />
          <MetricCard title="Status" value={project.status.replace("_", " ")} icon={<CheckCircle2 className="h-4 w-4 text-primary" />} />
          <div className="md:col-span-2 xl:col-span-4 space-y-3">
            <div>
              <p className="text-xs text-muted-foreground mb-2">Completude Fase 8 (Documentação)</p>
              <Progress value={docCompletion} className="h-2" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground mb-2">Completude Fase 9 (DevOps/Deploy)</p>
              <Progress value={infraCompletion} className="h-2" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue="documentation" className="w-full">
        <TabsList className="bg-card/50 border border-border p-1">
          <TabsTrigger value="documentation">Documentação</TabsTrigger>
          <TabsTrigger value="infrastructure">Infraestrutura</TabsTrigger>
        </TabsList>

        <TabsContent value="documentation" className="mt-5">
          <div className="grid gap-5 lg:grid-cols-[300px_1fr]">
            <Card className="border-border/60 bg-card/70">
              <CardHeader>
                <CardTitle className="text-base">Artefatos de Docs</CardTitle>
                <CardDescription>README, API, operações, contribuição e relatório final.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-2">
                {docs.map((file) => (
                  <button
                    key={file.id}
                    onClick={() => setSelectedDoc(file.path)}
                    className={`w-full rounded-lg border p-2 text-left text-sm transition ${selectedDocFile?.id === file.id ? "border-primary bg-primary/10" : "border-border/60 hover:border-primary/50"}`}
                  >
                    <p className="font-medium truncate">{file.path}</p>
                    <p className="text-xs text-muted-foreground">Atualizado {new Date(file.updated_at).toLocaleString()}</p>
                  </button>
                ))}
                {!docs.length && <p className="text-sm text-muted-foreground">Nenhum artefato de documentação disponível.</p>}
              </CardContent>
            </Card>

            <Card className="border-border/60 bg-card/70">
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle className="text-base">{selectedDocFile?.path ?? "Selecione um artefato"}</CardTitle>
                  <CardDescription>Renderização com leitura otimizada para revisão técnica.</CardDescription>
                </div>
                {selectedDocFile && <Badge variant="outline">{selectedDocFile.language}</Badge>}
              </CardHeader>
              <CardContent>
                <ScrollArea className="h-[560px] rounded-lg border border-border/60 bg-background/50 p-4">
                  {selectedDocFile ? <Markdown content={selectedDocFile.content} /> : <p className="text-sm text-muted-foreground">Sem conteúdo.</p>}
                </ScrollArea>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="infrastructure" className="mt-5">
          <div className="grid gap-5 lg:grid-cols-[300px_1fr]">
            <Card className="border-border/60 bg-card/70">
              <CardHeader>
                <CardTitle className="text-base">Artefatos de Infra</CardTitle>
                <CardDescription>Docker, pipelines de CI/CD e manifests Kubernetes.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-2">
                {infra.map((file) => (
                  <button
                    key={file.id}
                    onClick={() => setSelectedInfra(file.path)}
                    className={`w-full rounded-lg border p-2 text-left text-sm transition ${selectedInfraFile?.id === file.id ? "border-primary bg-primary/10" : "border-border/60 hover:border-primary/50"}`}
                  >
                    <p className="font-medium truncate">{file.path}</p>
                    <p className="text-xs text-muted-foreground">Atualizado {new Date(file.updated_at).toLocaleString()}</p>
                  </button>
                ))}
                {!infra.length && <p className="text-sm text-muted-foreground">Nenhum artefato de infraestrutura disponível.</p>}
              </CardContent>
            </Card>

            <Card className="border-border/60 bg-card/70">
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle className="text-base">{selectedInfraFile?.path ?? "Selecione um artefato"}</CardTitle>
                  <CardDescription>Leitor técnico para Dockerfiles, workflows e manifests.</CardDescription>
                </div>
                {selectedInfraFile && <Badge variant="outline">{selectedInfraFile.language}</Badge>}
              </CardHeader>
              <CardContent>
                <ScrollArea className="h-[560px] rounded-lg border border-border/60 bg-[#0b1220] p-4">
                  <pre className="text-xs leading-6 text-emerald-100 whitespace-pre-wrap">{selectedInfraFile?.content ?? "Sem conteúdo."}</pre>
                </ScrollArea>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>

      <Dialog open={showCompletionModal} onOpenChange={setShowCompletionModal}>
        <DialogContent className="border-primary/30 bg-card/95">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-2xl">
              <Sparkles className="h-6 w-6 text-yellow-500" />
              🎉 Projeto Concluído!
            </DialogTitle>
          </DialogHeader>
          <div className="grid gap-3 sm:grid-cols-2">
            <Highlight label="Fases concluídas" value={summaryMetrics.phases} />
            <Highlight label="Arquivos gerados" value={summaryMetrics.generatedFiles} />
            <Highlight label="Custo estimado" value={summaryMetrics.estimatedCost} />
            <Highlight label="Cobertura de testes" value={summaryMetrics.coverage} />
            <Highlight label="Segurança" value={summaryMetrics.security} />
          </div>
          <div className="flex gap-2">
            <Button className="flex-1 gap-2" onClick={downloadZip}>
              <Download className="h-4 w-4" /> Baixar Artefatos
            </Button>
            <Link href="/projects/new" className="flex-1">
              <Button variant="outline" className="w-full">Criar Novo Projeto</Button>
            </Link>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}

function MetricCard({ title, value, icon }: { title: string; value: string; icon: ReactNode }) {
  return (
    <div className="rounded-xl border border-border/60 bg-background/40 p-4">
      <p className="text-xs uppercase tracking-wide text-muted-foreground">{title}</p>
      <div className="mt-2 flex items-center justify-between">
        <p className="text-2xl font-semibold">{value}</p>
        {icon}
      </div>
    </div>
  );
}

function Highlight({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg border border-border/60 bg-background/50 p-3">
      <p className="text-xs text-muted-foreground">{label}</p>
      <p className="text-sm font-semibold mt-1">{value}</p>
    </div>
  );
}
