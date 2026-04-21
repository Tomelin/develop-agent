"use client";

import { useEffect, useState, use } from "react";
import { useRouter } from "next/navigation";
import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { ProjectService } from "@/services/project";
import { Project, PhaseStatus } from "@/types/project";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ArrowLeft, Play, Pause, Archive, Settings2, Sparkles, CheckCircle2, Clock, XCircle, LayoutTemplate, MonitorPlay, Megaphone, Activity, FileCode2, MessagesSquare, CheckSquare, Coins } from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import { Progress } from "@/components/ui/progress";

import { Phase6ExecutionPanel } from "@/components/projects/Phase6ExecutionPanel";
import { Phase8Workspace } from "@/components/projects/phase8/Phase8Workspace";
import { Phase5DevelopmentCenter } from "@/components/projects/phase5/Phase5DevelopmentCenter";
import { SecurityAuditPanel } from "@/components/projects/phase7/SecurityAuditPanel";
import { Phase13DeliveryCenter } from "@/components/projects/phase13/Phase13DeliveryCenter";
import { Phase14LandingCenter } from "@/components/projects/phase14/Phase14LandingCenter";
import { Phase15MarketingCenter } from "@/components/projects/phase15/Phase15MarketingCenter";
import { DynamicModeControl } from "@/components/projects/phase17/DynamicModeControl";
import { TriadCompositionPanel } from "@/components/projects/phase17/TriadCompositionPanel";
import { DiversityInsightsCard } from "@/components/projects/phase17/DiversityInsightsCard";
import { AgentPhaseMatrix } from "@/components/projects/phase17/AgentPhaseMatrix";
import { ProjectTeamPanel } from "@/components/phase20/ProjectTeamPanel";
import { ProjectIntegrationsPanel } from "@/components/phase20/ProjectIntegrationsPanel";

export default function ProjectDetailsPage({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = use(params);
  const router = useRouter();
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProject = async () => {
      try {
        const data = await ProjectService.getProjectById(resolvedParams.id);
        setProject(data);
      } catch (error) {
        console.error(error);
        toast.error("Erro ao carregar projeto");
        router.push("/dashboard");
      } finally {
        setLoading(false);
      }
    };

    fetchProject();

    // Setup SSE for real-time phase updates
    // Using polling as fallback since standard EventSource might have auth header issues
    const interval = setInterval(fetchProject, 10000); // 10s polling

    return () => {
      clearInterval(interval);
    };
  }, [resolvedParams.id, router]);

  const handleAction = async (action: 'pause' | 'resume' | 'archive') => {
    if (!project) return;
    try {
      if (action === 'pause') await ProjectService.pauseProject(project.id);
      if (action === 'resume') await ProjectService.resumeProject(project.id);
      if (action === 'archive') await ProjectService.archiveProject(project.id);

      toast.success(`Ação executada com sucesso!`);
      // Simulating a state update after action. In a real scenario we could call fetchProject here.
      // Doing a simple reload for now since fetchProject is encapsulated in useEffect.
      window.location.reload();
    } catch (err) {
      console.error(err);
      toast.error("Erro ao executar ação.");
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "DRAFT": return "bg-gray-500/20 text-gray-400";
      case "IN_PROGRESS": return "bg-primary/20 text-primary";
      case "PAUSED": return "bg-yellow-500/20 text-yellow-500";
      case "COMPLETED": return "bg-green-500/20 text-green-500";
      case "ARCHIVED": return "bg-destructive/20 text-destructive";
      default: return "bg-gray-500/20 text-gray-400";
    }
  };

  const getPhaseStatusIcon = (status: PhaseStatus, isCurrent: boolean) => {
    if (status === "COMPLETED") return <CheckCircle2 className="h-5 w-5 text-green-500" />;
    if (status === "REJECTED") return <XCircle className="h-5 w-5 text-destructive" />;
    if (status === "REVIEW") return <MessagesSquare className="h-5 w-5 text-yellow-500" />;
    if (isCurrent || status === "IN_PROGRESS") return <Activity className="h-5 w-5 text-primary animate-pulse" />;
    return <Clock className="h-5 w-5 text-muted-foreground" />;
  };

  const flowIcon = project?.flow_type === "A" ? MonitorPlay : project?.flow_type === "B" ? LayoutTemplate : Megaphone;
  const FlowIconComponent = flowIcon;

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  if (!project) return null;

  // Mocking phase definitions based on flow type for the timeline UI
  const totalPhases = project.flow_type === "A" ? 8 : project.flow_type === "B" ? 5 : 6;
  const phases = Array.from({ length: totalPhases }).map((_, i) => ({
    id: i + 1,
    name: `Fase ${i + 1}`,
    status: (i + 1 < project.current_phase ? "COMPLETED" : i + 1 === project.current_phase ? "IN_PROGRESS" : "PENDING") as PhaseStatus,
  }));

  return (
    <PrivateRoute>
      <div className="container py-8 max-w-6xl mx-auto space-y-6 animate-in fade-in duration-500">
        <Button variant="ghost"  className="-ml-4 text-muted-foreground">
          <Link href="/dashboard"><ArrowLeft className="mr-2 h-4 w-4" /> Voltar ao Dashboard</Link>
        </Button>

        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-card/50 p-6 rounded-2xl border border-border backdrop-blur-sm">
          <div className="flex items-start gap-4">
            <div className={`p-4 rounded-xl ${project.flow_type === 'A' ? 'bg-primary/10' : project.flow_type === 'B' ? 'bg-secondary/10' : 'bg-chart-3/10'}`}>
              <FlowIconComponent className={`h-8 w-8 ${project.flow_type === 'A' ? 'text-primary' : project.flow_type === 'B' ? 'text-secondary' : 'text-chart-3'}`} />
            </div>
            <div>
              <div className="flex items-center gap-3 mb-1">
                <h1 className="text-2xl font-bold">{project.name}</h1>
                <Badge variant="outline" className={getStatusColor(project.status)}>
                  {project.status.replace("_", " ")}
                </Badge>
              </div>
              <div className="flex items-center gap-4 text-sm text-muted-foreground">
                <span className="flex items-center gap-1">
                  <Settings2 className="h-4 w-4" />
                  Fluxo {project.flow_type}
                </span>
                {project.dynamic_mode && (
                  <span className="flex items-center gap-1 text-primary">
                    <Sparkles className="h-4 w-4" />
                    Modo Dinâmico
                  </span>
                )}
              </div>
            </div>
          </div>

          <div className="flex items-center gap-2">
            {project.flow_type === "A" && (
              <Button variant="secondary" size="sm" >
                <Link href={`/projects/${project.id}/interview`}>
                  <MessagesSquare className="mr-2 h-4 w-4" /> Entrevista
                </Link>
              </Button>
            )}
            {project.status === "IN_PROGRESS" && (
              <Button variant="outline" size="sm" onClick={() => handleAction('pause')}>
                <Pause className="mr-2 h-4 w-4" /> Pausar
              </Button>
            )}
            {project.status === "PAUSED" && (
              <Button variant="outline" size="sm" onClick={() => handleAction('resume')}>
                <Play className="mr-2 h-4 w-4" /> Retomar
              </Button>
            )}
            {(project.status !== "ARCHIVED" && project.status !== "COMPLETED") && (
              <Button variant="destructive" size="sm" onClick={() => handleAction('archive')}>
                <Archive className="mr-2 h-4 w-4" /> Arquivar
              </Button>
            )}
          </div>
        </div>

        <Tabs defaultValue="timeline" className="w-full">
          <TabsList className="bg-card/50 border border-border p-1">
            <TabsTrigger value="timeline">Timeline do Projeto</TabsTrigger>
            <TabsTrigger value="kanban">Roadmap</TabsTrigger>
            <TabsTrigger value="phase5">Fase 05</TabsTrigger>
            <TabsTrigger value="phase6">Fase 06</TabsTrigger>
            <TabsTrigger value="phase8">Fase 08</TabsTrigger>
            {project.flow_type === "B" && <TabsTrigger value="phase14">Fluxo B</TabsTrigger>}
            {project.flow_type === "C" && <TabsTrigger value="phase15">Fluxo C</TabsTrigger>}
            <TabsTrigger value="phase12">Segurança</TabsTrigger>
            <TabsTrigger value="delivery">Entrega</TabsTrigger>
            <TabsTrigger value="triad">Tríade IA</TabsTrigger>
            <TabsTrigger value="agent-config">Config. de Agentes</TabsTrigger>
            <TabsTrigger value="team">Equipe</TabsTrigger>
            <TabsTrigger value="integrations">Integrações</TabsTrigger>
            <TabsTrigger value="settings">Configurações</TabsTrigger>
          </TabsList>

          <TabsContent value="timeline" className="mt-6 space-y-6">
            <div className="grid md:grid-cols-3 gap-6">

              {/* Timeline Column */}
              <div className="md:col-span-2 space-y-6">
                <Card className="bg-card/50 border-border">
                  <CardHeader>
                    <CardTitle className="text-lg">Progresso das Fases</CardTitle>
                    <div className="flex items-center gap-4 mt-2">
                      <Progress value={project.progress_percentage} className="h-2 w-full flex-1" />
                      <span className="text-sm font-medium">{project.progress_percentage}%</span>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4 relative before:absolute before:inset-0 before:ml-5 before:-translate-x-px md:before:mx-auto md:before:translate-x-0 before:h-full before:w-0.5 before:bg-gradient-to-b before:from-transparent before:via-border before:to-transparent">
                      {phases.map((phase) => (
                        <div key={phase.id} className="relative flex items-center justify-between md:justify-normal md:odd:flex-row-reverse group is-active">
                          <div className={`flex items-center justify-center w-10 h-10 rounded-full border-4 border-background shrink-0 md:order-1 md:group-odd:-translate-x-1/2 md:group-even:translate-x-1/2 shadow
                            ${phase.id === project.current_phase ? "bg-primary" : phase.id < project.current_phase ? "bg-green-500" : "bg-card"}
                          `}>
                            {getPhaseStatusIcon(phase.status, phase.id === project.current_phase)}
                          </div>
                          <div className={`w-[calc(100%-4rem)] md:w-[calc(50%-2.5rem)] p-4 rounded-xl border ${phase.id === project.current_phase ? "border-primary/50 bg-primary/5" : "border-border bg-card/50"}`}>
                            <div className="flex items-center justify-between mb-1">
                              <h3 className={`font-bold ${phase.id === project.current_phase ? "text-primary" : "text-foreground"}`}>{phase.name}</h3>
                              <Badge variant="outline" className={`text-[10px] uppercase ${phase.id === project.current_phase ? "bg-primary/20 text-primary border-primary/20" : ""}`}>
                                {phase.status === "IN_PROGRESS" ? "Em andamento" : phase.status === "COMPLETED" ? "Concluída" : "Pendente"}
                              </Badge>
                            </div>
                            <p className="text-sm text-muted-foreground">
                              {phase.id === 1 && "Planejamento inicial e definição de escopo."}
                              {phase.id === 2 && "Arquitetura e design de sistema."}
                              {phase.id === 3 && "Configuração de infraestrutura e repositórios."}
                              {phase.id > 3 && "Desenvolvimento e iterações."}
                            </p>

                            {phase.id === project.current_phase && (
                              <div className="mt-4 pt-4 border-t border-primary/20">
                                <Button size="sm" className="w-full gap-2">
                                  <Sparkles className="h-4 w-4" /> Acompanhar Agentes
                                </Button>
                              </div>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Sidebar Column */}
              <div className="space-y-6">
                <Card className="bg-card/50 border-border">
                  <CardHeader>
                    <CardTitle className="text-md flex items-center gap-2">
                      <Coins className="h-4 w-4 text-yellow-500" />
                      Custos do Projeto
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="flex justify-between items-center pb-4 border-b border-border/50">
                        <span className="text-sm text-muted-foreground">Tokens Utilizados</span>
                        <span className="font-mono font-medium">{project.tokens_used.toLocaleString()}</span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-muted-foreground">Custo Estimado</span>
                        <span className="font-mono font-medium text-destructive">~${((project.tokens_used / 1000) * 0.01).toFixed(2)}</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card className="bg-card/50 border-border">
                  <CardHeader>
                    <CardTitle className="text-md flex items-center gap-2">
                      <FileCode2 className="h-4 w-4 text-primary" />
                      Artefatos (Fase Atual)
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      <div className="p-3 rounded-lg border border-border/50 bg-background flex items-center justify-between group hover:border-primary/50 cursor-pointer transition-colors">
                        <div className="flex items-center gap-3">
                          <CheckSquare className="h-4 w-4 text-muted-foreground group-hover:text-primary" />
                          <span className="text-sm">Requirements.md</span>
                        </div>
                      </div>
                      <div className="p-3 rounded-lg border border-border/50 bg-background flex items-center justify-between group hover:border-primary/50 cursor-pointer transition-colors">
                        <div className="flex items-center gap-3">
                          <CheckSquare className="h-4 w-4 text-muted-foreground group-hover:text-primary" />
                          <span className="text-sm">Architecture_ADR_01.md</span>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </TabsContent>

          <TabsContent value="kanban" className="mt-6">

            <div className="flex flex-col items-center justify-center p-12 text-center bg-card/30 rounded-xl border border-border mt-4">
              <LayoutTemplate className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-xl font-semibold mb-2">Roadmap & KanBan</h3>
              <p className="text-muted-foreground mb-6 max-w-md">
                Gerencie o progresso do desenvolvimento, visualize épicos, dependências e a timeline do projeto.
              </p>
              <Link href={`/projects/${project.id}/roadmap`}>
                <Button className="gap-2">
                  <LayoutTemplate className="h-4 w-4" />
                  Abrir Dashboard do Roadmap
                </Button>
              </Link>
            </div>

          </TabsContent>


          <TabsContent value="phase5" className="mt-6">
            <Phase5DevelopmentCenter projectId={project.id} />
          </TabsContent>

          <TabsContent value="phase6" className="mt-6">
            <Phase6ExecutionPanel projectId={project.id} />
          </TabsContent>



          <TabsContent value="phase8" className="mt-6">
            <Phase8Workspace projectId={project.id} phaseNumber={2} />
          </TabsContent>

          {project.flow_type === "B" && (
            <TabsContent value="phase14" className="mt-6">
              <Phase14LandingCenter project={project} />
            </TabsContent>
          )}

          {project.flow_type === "C" && (
            <TabsContent value="phase15" className="mt-6">
              <Phase15MarketingCenter project={project} />
            </TabsContent>
          )}

          <TabsContent value="phase12" className="mt-6">
            <SecurityAuditPanel projectId={project.id} />
          </TabsContent>


          <TabsContent value="delivery" className="mt-6">
            <Phase13DeliveryCenter project={project} />
          </TabsContent>


          <TabsContent value="triad" className="mt-6 space-y-6">
            <DynamicModeControl projectId={project.id} initialEnabled={project.dynamic_mode} />
            <TriadCompositionPanel projectId={project.id} />
            <DiversityInsightsCard projectId={project.id} />
          </TabsContent>

          <TabsContent value="agent-config" className="mt-6">
            <AgentPhaseMatrix projectId={project.id} />
          </TabsContent>

          <TabsContent value="team" className="mt-6">
            <ProjectTeamPanel projectId={project.id} />
          </TabsContent>

          <TabsContent value="integrations" className="mt-6">
            <ProjectIntegrationsPanel projectId={project.id} />
          </TabsContent>

          <TabsContent value="settings" className="mt-6">
            <Card className="bg-card/50 border-border">
              <CardHeader>
                <CardTitle>Configurações do Projeto</CardTitle>
                <CardDescription>Detalhes do projeto fornecidos durante a criação.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <h4 className="text-sm font-medium text-muted-foreground">Descrição</h4>
                  <p className="mt-1">{project.description}</p>
                </div>
                <div>
                  <h4 className="text-sm font-medium text-muted-foreground">Criado em</h4>
                  <p className="mt-1">{new Date(project.created_at).toLocaleString()}</p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </PrivateRoute>
  );
}
