"use client";

import { useAuth } from "@/contexts/AuthContext";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Sparkles, Activity, Plus, Search, FolderKanban, Clock, Bot } from "lucide-react";
import { AgentStatusPanel } from "@/components/dashboard/AgentStatusPanel";
import { useEffect, useState } from "react";
import { Project, ProjectStatus, FlowType } from "@/types/project";
import { ProjectService } from "@/services/project";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Progress } from "@/components/ui/progress";

export default function DashboardPage() {
  const { user } = useAuth();
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState<ProjectStatus | "ALL">("ALL");
  const [flowFilter, setFlowFilter] = useState<FlowType | "ALL">("ALL");

  const [stats, setStats] = useState({
    total: 0,
    inProgress: 0,
    completed: 0,
    tokens: 0,
  });

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        // In a real scenario we might need pagination, but let's fetch first page with large size for dashboard overview
        const response = await ProjectService.getProjects(1, 50, statusFilter === "ALL" ? undefined : statusFilter, flowFilter === "ALL" ? undefined : flowFilter);
        setProjects(response.items || []);

        // Calculate stats based on all items (if we had a separate stats endpoint we'd use that)
        // For now, we estimate from fetched projects if no filters applied
        if (statusFilter === "ALL" && flowFilter === "ALL") {
          const inProgress = (response.items || []).filter(p => p.status === "IN_PROGRESS").length;
          const completed = (response.items || []).filter(p => p.status === "COMPLETED").length;
          const totalTokens = (response.items || []).reduce((acc, p) => acc + (p.tokens_used || 0), 0);

          setStats({
            total: response.total || response.items?.length || 0,
            inProgress,
            completed,
            tokens: totalTokens,
          });
        }
      } catch (error) {
        console.error("Failed to fetch projects:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchProjects();
    const interval = setInterval(fetchProjects, 30000); // 30s polling
    return () => clearInterval(interval);
  }, [statusFilter, flowFilter]);

  const filteredProjects = projects.filter(p =>
    p.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    p.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getStatusColor = (status: ProjectStatus) => {
    switch (status) {
      case "DRAFT": return "bg-gray-500/20 text-gray-400";
      case "IN_PROGRESS": return "bg-primary/20 text-primary";
      case "PAUSED": return "bg-yellow-500/20 text-yellow-500";
      case "COMPLETED": return "bg-green-500/20 text-green-500";
      case "ARCHIVED": return "bg-destructive/20 text-destructive";
      default: return "bg-gray-500/20 text-gray-400";
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground mt-1">
            Bem-vindo de volta, <span className="text-primary font-medium">{user?.name}</span>!
          </p>
        </div>
        <Button  className="gap-2 bg-primary text-primary-foreground hover:bg-primary/90">
          <Link href="/projects/new">
            <Plus className="h-4 w-4" /> Novo Projeto
          </Link>
        </Button>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card className="bg-card/50 backdrop-blur-sm border-border hover:border-primary/50 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Projetos</CardTitle>
            <FolderKanban className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
          </CardContent>
        </Card>
        <Card className="bg-card/50 backdrop-blur-sm border-border hover:border-primary/50 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Em Andamento</CardTitle>
            <Activity className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.inProgress}</div>
          </CardContent>
        </Card>
        <Card className="bg-card/50 backdrop-blur-sm border-border hover:border-primary/50 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Concluídos</CardTitle>
            <Sparkles className="h-4 w-4 text-secondary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.completed}</div>
          </CardContent>
        </Card>
        <Card className="bg-card/50 backdrop-blur-sm border-border hover:border-primary/50 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Tokens Utilizados (Mês)</CardTitle>
            <Bot className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.tokens.toLocaleString()}</div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-3 lg:grid-cols-4">
        {/* Main Projects Section */}
        <div className="md:col-span-2 lg:col-span-3 space-y-4">
          <div className="flex flex-col sm:flex-row gap-4 items-center justify-between bg-card/50 p-4 rounded-xl border border-border">
            <div className="relative w-full sm:w-96">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Buscar projetos..."
                className="pl-9 bg-background"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
            <div className="flex w-full sm:w-auto gap-2">
              <Select value={statusFilter} onValueChange={(v: string | null) => v && setStatusFilter(v as any)}>
                <SelectTrigger className="w-[140px] bg-background">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ALL">Todos Status</SelectItem>
                  <SelectItem value="DRAFT">Rascunho</SelectItem>
                  <SelectItem value="IN_PROGRESS">Em Andamento</SelectItem>
                  <SelectItem value="PAUSED">Pausado</SelectItem>
                  <SelectItem value="COMPLETED">Concluído</SelectItem>
                </SelectContent>
              </Select>
              <Select value={flowFilter} onValueChange={(v: string | null) => v && setFlowFilter(v as any)}>
                <SelectTrigger className="w-[140px] bg-background">
                  <SelectValue placeholder="Fluxo" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ALL">Todos Fluxos</SelectItem>
                  <SelectItem value="A">Software (A)</SelectItem>
                  <SelectItem value="B">Landing Page (B)</SelectItem>
                  <SelectItem value="C">Marketing (C)</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {loading ? (
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {[1, 2, 3].map(i => (
                <Card key={i} className="animate-pulse bg-card/30 h-48"></Card>
              ))}
            </div>
          ) : filteredProjects.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center border border-dashed border-border rounded-xl bg-card/30">
              <div className="bg-primary/10 p-4 rounded-full mb-4">
                <Sparkles className="h-8 w-8 text-primary" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Nenhum projeto encontrado</h3>
              <p className="text-muted-foreground max-w-md mb-6">
                {searchTerm || statusFilter !== "ALL" || flowFilter !== "ALL"
                  ? "Tente ajustar seus filtros ou termos de busca para encontrar o que procura."
                  : "Você ainda não possui projetos. Crie seu primeiro projeto para começar a utilizar a plataforma Agency AI."}
              </p>
              {!(searchTerm || statusFilter !== "ALL" || flowFilter !== "ALL") && (
                <Button  size="lg">
                  <Link href="/projects/new">Criar Meu Primeiro Projeto</Link>
                </Button>
              )}
            </div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {filteredProjects.map((project) => (
                <Link key={project.id} href={`/projects/${project.id}`} className="block group">
                  <Card className="h-full bg-card/50 hover:bg-card transition-all duration-200 border-border hover:border-primary/50 flex flex-col">
                    <CardHeader className="pb-3">
                      <div className="flex justify-between items-start mb-2">
                        <Badge variant="outline" className={getStatusColor(project.status)}>
                          {project.status.replace("_", " ")}
                        </Badge>
                        <Badge variant="secondary" className="font-mono">
                          Fluxo {project.flow_type}
                        </Badge>
                      </div>
                      <CardTitle className="text-lg group-hover:text-primary transition-colors line-clamp-1">{project.name}</CardTitle>
                    </CardHeader>
                    <CardContent className="flex-1 pb-3">
                      <div className="space-y-4">
                        <div className="space-y-1">
                          <div className="flex justify-between text-xs text-muted-foreground mb-1">
                            <span>Progresso (Fase {project.current_phase})</span>
                            <span>{project.progress_percentage}%</span>
                          </div>
                          <Progress value={project.progress_percentage} className="h-1.5" />
                        </div>
                      </div>
                    </CardContent>
                    <CardFooter className="pt-0 pb-4 text-xs text-muted-foreground flex items-center border-t border-border/50 mt-auto px-6 pt-4">
                      <Clock className="h-3 w-3 mr-1" />
                      Atualizado em {new Date(project.updated_at).toLocaleDateString()}
                    </CardFooter>
                  </Card>
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Right Sidebar - Agent Status Panel */}
        <div className="md:col-span-1 h-[calc(100vh-12rem)] sticky top-24">
          <AgentStatusPanel />
        </div>
      </div>
    </div>
  );
}