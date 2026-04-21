"use client";

import { ReactNode, useEffect, useMemo, useState } from "react";
import { ProjectService } from "@/services/project";
import { Task } from "@/types/task";
import { Phase5CodeContext, Phase5CodeFile, Phase5ExecutionMode, Phase5Summary, Phase5TaskExecutionEvent } from "@/types/phase5";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import {
  Bot,
  CircleAlert,
  Code2,
  Download,
  FileCode2,
  Folder,
  FolderOpen,
  Loader2,
  PlayCircle,
  RefreshCcw,
  Search,
  TerminalSquare,
  Timer,
} from "lucide-react";

interface Phase5DevelopmentCenterProps {
  projectId: string;
}

interface FileNode {
  name: string;
  fullPath: string;
  file?: Phase5CodeFile;
  children?: FileNode[];
}

export function Phase5DevelopmentCenter({ projectId }: Phase5DevelopmentCenterProps) {
  const [summary, setSummary] = useState<Phase5Summary | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [files, setFiles] = useState<Phase5CodeFile[]>([]);
  const [context, setContext] = useState<Phase5CodeContext | null>(null);

  const [selectedMode, setSelectedMode] = useState<Phase5ExecutionMode>("MANUAL");
  const [selectedTaskId, setSelectedTaskId] = useState<string | null>(null);
  const [selectedFileId, setSelectedFileId] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [expandedFolders, setExpandedFolders] = useState<Record<string, boolean>>({});
  const [terminalEvents, setTerminalEvents] = useState<Phase5TaskExecutionEvent[]>([]);

  const [isLoading, setIsLoading] = useState(true);
  const [isUpdatingMode, setIsUpdatingMode] = useState(false);
  const [isExecutingAll, setIsExecutingAll] = useState(false);
  const [isExecutingTask, setIsExecutingTask] = useState(false);
  const [isDownloading, setIsDownloading] = useState(false);

  const selectedFile = useMemo(
    () => files.find((file) => file.id === selectedFileId) ?? null,
    [files, selectedFileId],
  );

  const filteredFiles = useMemo(() => {
    if (!searchTerm.trim()) return files;
    const normalized = searchTerm.trim().toLowerCase();
    return files.filter((file) => file.path.toLowerCase().includes(normalized));
  }, [files, searchTerm]);

  const fileTree = useMemo(() => buildTree(filteredFiles), [filteredFiles]);

  const appendTerminalEvent = (event: Omit<Phase5TaskExecutionEvent, "timestamp">) => {
    setTerminalEvents((current) => [
      {
        timestamp: new Date().toISOString(),
        ...event,
      },
      ...current,
    ].slice(0, 60));
  };

  const loadPhaseData = async () => {
    try {
      const [summaryData, taskData, fileData, contextData] = await Promise.all([
        ProjectService.getPhase5Summary(projectId),
        ProjectService.getTasks(projectId, 1, 200),
        ProjectService.getProjectFiles(projectId),
        ProjectService.getPhase5CodeContext(projectId),
      ]);

      setSummary(summaryData);
      setTasks(taskData.items);
      setFiles(fileData);
      setContext(contextData);
      setSelectedMode(summaryData.execution_mode);

      if (!selectedTaskId && taskData.items.length > 0) {
        const pending = taskData.items.find((task) => task.status !== "DONE");
        setSelectedTaskId(pending?.id ?? taskData.items[0].id);
      }

      if (!selectedFileId && fileData.length > 0) {
        setSelectedFileId(fileData[0].id);
      }
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível carregar os dados da Fase 05.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadPhaseData();
    const interval = setInterval(loadPhaseData, 10000);
    return () => clearInterval(interval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [projectId]);

  const updateMode = async () => {
    try {
      setIsUpdatingMode(true);
      await ProjectService.setPhase5Mode(projectId, selectedMode);
      appendTerminalEvent({
        level: "SUCCESS",
        message: `Modo alterado para ${selectedMode === "AUTOMATIC" ? "Automático" : "Manual"}.`,
      });
      toast.success("Modo de execução atualizado.");
      await loadPhaseData();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao atualizar modo da Fase 05.");
    } finally {
      setIsUpdatingMode(false);
    }
  };

  const runAutomaticExecution = async () => {
    try {
      setIsExecutingAll(true);
      appendTerminalEvent({ level: "INFO", message: "Iniciando execução automática da Fase 05..." });
      const response = await ProjectService.executeAllPhase5Tasks(projectId);
      appendTerminalEvent({
        level: "SUCCESS",
        message: `${response.executed_tasks} task(s) executadas automaticamente.`,
      });
      toast.success("Execução automática concluída.");
      await loadPhaseData();
    } catch (error) {
      console.error(error);
      appendTerminalEvent({ level: "ERROR", message: "Falha na execução automática das tasks." });
      toast.error("Falha ao executar tasks automaticamente.");
    } finally {
      setIsExecutingAll(false);
    }
  };

  const executeSelectedTask = async () => {
    if (!selectedTaskId) return;
    const task = tasks.find((item) => item.id === selectedTaskId);
    if (!task) return;

    try {
      setIsExecutingTask(true);
      appendTerminalEvent({
        level: "INFO",
        taskId: task.id,
        taskTitle: task.title,
        message: `Executando task ${task.title}...`,
      });
      await ProjectService.executePhase5Task(projectId, task.id);
      appendTerminalEvent({
        level: "SUCCESS",
        taskId: task.id,
        taskTitle: task.title,
        message: `Task ${task.title} concluída e versionada.`,
      });
      toast.success("Task executada com sucesso.");
      await loadPhaseData();
    } catch (error) {
      console.error(error);
      appendTerminalEvent({
        level: "ERROR",
        taskId: task.id,
        taskTitle: task.title,
        message: `Task ${task.title} bloqueada ou falhou na validação.`,
      });
      toast.error("Falha ao executar task selecionada.");
    } finally {
      setIsExecutingTask(false);
    }
  };

  const downloadZip = async () => {
    try {
      setIsDownloading(true);
      const blob = await ProjectService.downloadProjectFilesZip(projectId);
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = `project-${projectId}-phase5-files.zip`;
      document.body.appendChild(link);
      link.click();
      URL.revokeObjectURL(link.href);
      link.remove();
      toast.success("Download iniciado.");
    } catch (error) {
      console.error(error);
      toast.error("Falha no download do ZIP.");
    } finally {
      setIsDownloading(false);
    }
  };

  const toggleFolder = (path: string) => {
    setExpandedFolders((current) => ({ ...current, [path]: !current[path] }));
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[30vh]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Card className="border-border bg-card/60 backdrop-blur-sm">
        <CardHeader className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Code2 className="h-5 w-5 text-primary" />
              Fase 05 — Centro de Desenvolvimento
            </CardTitle>
            <CardDescription>
              Execução manual/automática das tasks, observabilidade de progresso e repositório virtual de código.
            </CardDescription>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" className="gap-2" onClick={loadPhaseData}>
              <RefreshCcw className="h-4 w-4" /> Atualizar
            </Button>
            <Button onClick={runAutomaticExecution} disabled={isExecutingAll || isExecutingTask} className="gap-2">
              {isExecutingAll ? <Loader2 className="h-4 w-4 animate-spin" /> : <PlayCircle className="h-4 w-4" />}
              Executar pendentes
            </Button>
            <Button variant="secondary" onClick={downloadZip} disabled={isDownloading} className="gap-2">
              {isDownloading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Download className="h-4 w-4" />}
              Download ZIP
            </Button>
          </div>
        </CardHeader>

        {summary && (
          <CardContent className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            <MetricCard title="Progresso" value={`${summary.completion_percent.toFixed(0)}%`} subtitle={`${summary.done_tasks}/${summary.total_tasks} tasks`} />
            <MetricCard title="Arquivos gerados" value={`${summary.backend_files + summary.frontend_files}`} subtitle={`${summary.backend_files} back • ${summary.frontend_files} front`} />
            <MetricCard title="Linhas geradas" value={summary.generated_lines_of_code.toLocaleString()} subtitle="LOC total da fase" />
            <MetricCard title="Auto rejeições" value={`${summary.auto_rejections}`} subtitle={`Modo: ${summary.execution_mode}`} />
            <div className="sm:col-span-2 xl:col-span-4 pt-1">
              <Progress value={Math.min(summary.completion_percent, 100)} className="h-2" />
            </div>
          </CardContent>
        )}
      </Card>

      <div className="grid gap-6 xl:grid-cols-[1.25fr_0.95fr]">
        <Card className="border-border bg-card/60">
          <CardHeader>
            <CardTitle className="text-base">Modo de execução da fase</CardTitle>
            <CardDescription>
              Manual para controle task-a-task ou automático para throughput máximo.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-5">
            <RadioGroup
              value={selectedMode}
              onValueChange={(value) => setSelectedMode(value as Phase5ExecutionMode)}
              className="grid gap-3 md:grid-cols-2"
            >
              <label className="rounded-xl border border-border/60 p-4 cursor-pointer hover:border-primary/50 transition-colors">
                <div className="flex items-center gap-3">
                  <RadioGroupItem value="MANUAL" id="manual" />
                  <div>
                    <Label htmlFor="manual" className="text-sm font-semibold">Manual</Label>
                    <p className="text-xs text-muted-foreground mt-1">Aprovação por task, ideal para tarefas críticas.</p>
                  </div>
                </div>
              </label>

              <label className="rounded-xl border border-border/60 p-4 cursor-pointer hover:border-primary/50 transition-colors">
                <div className="flex items-center gap-3">
                  <RadioGroupItem value="AUTOMATIC" id="automatic" />
                  <div>
                    <Label htmlFor="automatic" className="text-sm font-semibold">Automático</Label>
                    <p className="text-xs text-muted-foreground mt-1">Executa todas as pendências em lote sequencial.</p>
                  </div>
                </div>
              </label>
            </RadioGroup>

            <div className="flex flex-wrap items-center gap-3">
              <Button onClick={updateMode} disabled={isUpdatingMode} className="gap-2">
                {isUpdatingMode ? <Loader2 className="h-4 w-4 animate-spin" /> : <Bot className="h-4 w-4" />}
                Salvar modo
              </Button>
              <Badge variant="outline" className="text-xs uppercase">modo atual: {summary?.execution_mode ?? "MANUAL"}</Badge>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border bg-card/60">
          <CardHeader>
            <CardTitle className="text-base">Code Context</CardTitle>
            <CardDescription>Manifesto de contexto acumulado utilizado nas próximas tasks.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm">
            <div className="grid grid-cols-2 gap-3">
              <div className="rounded-lg border border-border/60 p-3 bg-background/50">
                <p className="text-xs text-muted-foreground">Arquivos no contexto</p>
                <p className="text-xl font-semibold">{context?.files.length ?? 0}</p>
              </div>
              <div className="rounded-lg border border-border/60 p-3 bg-background/50">
                <p className="text-xs text-muted-foreground">Tokens aproximados</p>
                <p className="text-xl font-semibold">{context?.approx_tokens ?? 0}</p>
              </div>
            </div>

            <div className="rounded-lg border border-border/60 p-3 bg-background/50">
              <p className="text-xs text-muted-foreground mb-2">Dependências</p>
              <div className="flex flex-wrap gap-2">
                {(context?.dependencies.length ? context.dependencies : ["Nenhuma dependência detectada"]).map((dep) => (
                  <Badge variant="secondary" key={dep} className="font-normal">{dep}</Badge>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
        <Card className="border-border bg-card/60">
          <CardHeader>
            <CardTitle className="text-base">Execução Task-by-Task</CardTitle>
            <CardDescription>Selecione uma task e execute com validação e persistência no repositório virtual.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <ScrollArea className="h-64 rounded-lg border border-border/60 p-2">
              <div className="space-y-2">
                {tasks.map((task) => (
                  <button
                    key={task.id}
                    type="button"
                    onClick={() => setSelectedTaskId(task.id)}
                    className={`w-full rounded-lg border p-3 text-left transition-colors ${
                      selectedTaskId === task.id ? "border-primary bg-primary/10" : "border-border/60 hover:border-primary/40"
                    }`}
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div>
                        <p className="text-sm font-medium line-clamp-1">{task.title}</p>
                        <p className="text-xs text-muted-foreground mt-1 line-clamp-2">{task.description}</p>
                      </div>
                      <Badge variant="outline" className={statusClass(task.status)}>{task.status}</Badge>
                    </div>
                  </button>
                ))}
              </div>
            </ScrollArea>

            <div className="flex gap-3">
              <Button onClick={executeSelectedTask} disabled={!selectedTaskId || isExecutingTask || isExecutingAll} className="gap-2 w-full">
                {isExecutingTask ? <Loader2 className="h-4 w-4 animate-spin" /> : <PlayCircle className="h-4 w-4" />}
                Executar task selecionada
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border bg-card/60">
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <TerminalSquare className="h-4 w-4 text-primary" />
              Terminal de desenvolvimento
            </CardTitle>
            <CardDescription>Eventos de execução em tempo real da sessão atual.</CardDescription>
          </CardHeader>
          <CardContent>
            <ScrollArea className="h-80 rounded-lg border border-border/60 bg-zinc-950 text-zinc-100 p-3">
              <div className="space-y-3 text-xs font-mono">
                {terminalEvents.length === 0 && (
                  <p className="text-zinc-400">Nenhum evento ainda. Execute uma task para iniciar o stream de eventos operacionais.</p>
                )}

                {terminalEvents.map((event, index) => (
                  <div key={`${event.timestamp}-${index}`} className="space-y-1">
                    <div className="flex items-center gap-2 text-zinc-400">
                      <Timer className="h-3.5 w-3.5" />
                      <span>{new Date(event.timestamp).toLocaleTimeString()}</span>
                      {event.taskTitle && <span>• {event.taskTitle}</span>}
                    </div>
                    <p className={event.level === "ERROR" ? "text-red-300" : event.level === "SUCCESS" ? "text-emerald-300" : "text-zinc-100"}>
                      {event.message}
                    </p>
                  </div>
                ))}
              </div>
            </ScrollArea>
            <p className="mt-3 text-xs text-muted-foreground flex items-center gap-2">
              <CircleAlert className="h-3.5 w-3.5" />
              Quando o backend disponibilizar stream SSE dedicado para Fase 05, este painel poderá consumir o feed contínuo sem polling.
            </p>
          </CardContent>
        </Card>
      </div>

      <Card className="border-border bg-card/60">
        <CardHeader className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div>
            <CardTitle className="text-base flex items-center gap-2">
              <FileCode2 className="h-4 w-4 text-primary" />
              File Explorer
            </CardTitle>
            <CardDescription>Navegue pelos arquivos gerados e inspecione conteúdo com busca por caminho.</CardDescription>
          </div>

          <div className="relative w-full md:max-w-sm">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} placeholder="Buscar arquivo por nome..." className="pl-9" />
          </div>
        </CardHeader>
        <CardContent className="grid gap-4 lg:grid-cols-[320px_1fr]">
          <ScrollArea className="h-[460px] rounded-xl border border-border/60 p-2">
            <div className="space-y-1">{fileTree.map((node) => renderNode(node, expandedFolders, toggleFolder, setSelectedFileId, selectedFileId))}</div>
          </ScrollArea>

          <div className="rounded-xl border border-border/60 bg-background/40 p-4 min-h-[460px]">
            {selectedFile ? (
              <div className="space-y-4">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <div>
                    <p className="text-sm font-semibold break-all">{selectedFile.path}</p>
                    <p className="text-xs text-muted-foreground">Task {selectedFile.task_id} • versão {new Date(selectedFile.version).toLocaleString()}</p>
                  </div>
                  <Badge variant="outline">{selectedFile.language}</Badge>
                </div>
                <ScrollArea className="h-[360px] rounded-md border border-border/60 bg-zinc-950 p-3">
                  <pre className="text-xs text-zinc-100 whitespace-pre-wrap font-mono">{selectedFile.content}</pre>
                </ScrollArea>
              </div>
            ) : (
              <div className="h-full flex items-center justify-center text-sm text-muted-foreground">
                Selecione um arquivo para visualizar o conteúdo.
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function MetricCard({ title, value, subtitle }: { title: string; value: string; subtitle: string }) {
  return (
    <div className="rounded-xl border border-border/60 bg-background/40 p-4">
      <p className="text-xs text-muted-foreground">{title}</p>
      <p className="text-2xl font-semibold mt-2">{value}</p>
      <p className="text-xs text-muted-foreground mt-1">{subtitle}</p>
    </div>
  );
}

function statusClass(status: Task["status"]) {
  switch (status) {
    case "DONE":
      return "border-emerald-500/40 text-emerald-600";
    case "IN_PROGRESS":
      return "border-blue-500/40 text-blue-600";
    case "BLOCKED":
      return "border-destructive/40 text-destructive";
    case "REVIEW":
      return "border-yellow-500/40 text-yellow-600";
    default:
      return "border-border text-muted-foreground";
  }
}

function buildTree(files: Phase5CodeFile[]): FileNode[] {
  const root: FileNode[] = [];

  for (const file of files) {
    const parts = file.path.split("/").filter(Boolean);
    let currentLevel = root;
    let accumulatedPath = "";

    parts.forEach((part, index) => {
      accumulatedPath = accumulatedPath ? `${accumulatedPath}/${part}` : part;
      const isFile = index === parts.length - 1;

      let node = currentLevel.find((entry) => entry.name === part);
      if (!node) {
        node = {
          name: part,
          fullPath: accumulatedPath,
          children: isFile ? undefined : [],
          file: isFile ? file : undefined,
        };
        currentLevel.push(node);
      }

      if (!isFile) {
        currentLevel = node.children ?? [];
        node.children = currentLevel;
      }
    });
  }

  const sortNodes = (nodes: FileNode[]) => {
    nodes.sort((a, b) => {
      if (a.children && !b.children) return -1;
      if (!a.children && b.children) return 1;
      return a.name.localeCompare(b.name);
    });
    nodes.forEach((node) => node.children && sortNodes(node.children));
  };

  sortNodes(root);
  return root;
}

function renderNode(
  node: FileNode,
  expandedFolders: Record<string, boolean>,
  toggleFolder: (path: string) => void,
  onFileSelect: (fileId: string) => void,
  selectedFileId: string | null,
  depth = 0,
): ReactNode {
  const isFolder = Boolean(node.children);
  const isExpanded = expandedFolders[node.fullPath] ?? depth < 2;

  if (isFolder) {
    return (
      <div key={node.fullPath}>
        <button
          type="button"
          onClick={() => toggleFolder(node.fullPath)}
          className="w-full flex items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-primary/10"
          style={{ paddingLeft: `${8 + depth * 14}px` }}
        >
          {isExpanded ? <FolderOpen className="h-4 w-4 text-primary" /> : <Folder className="h-4 w-4 text-muted-foreground" />}
          <span className="truncate">{node.name}</span>
        </button>

        {isExpanded && node.children?.map((child) => renderNode(child, expandedFolders, toggleFolder, onFileSelect, selectedFileId, depth + 1))}
      </div>
    );
  }

  return (
    <button
      key={node.fullPath}
      type="button"
      onClick={() => node.file?.id && onFileSelect(node.file.id)}
      className={`w-full flex items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-primary/10 ${selectedFileId === node.file?.id ? "bg-primary/10 text-primary" : ""
        }`}
      style={{ paddingLeft: `${8 + depth * 14}px` }}
    >
      <FileCode2 className="h-3.5 w-3.5" />
      <span className="truncate">{node.name}</span>
    </button>
  );
}
