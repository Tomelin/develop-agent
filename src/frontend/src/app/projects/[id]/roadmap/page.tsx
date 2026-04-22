"use client";

import { useEffect, useState, use } from "react";
import { useRouter } from "next/navigation";
import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { ProjectService } from "@/services/project";
import { RoadmapService } from "@/services/roadmap";
import { Project } from "@/types/project";
import { RoadmapSummary } from "@/types/roadmap";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ArrowLeft, LayoutTemplate, Network, Clock, Download } from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import { AdvancedKanbanBoard } from "@/components/projects/roadmap/AdvancedKanbanBoard";
import { EpicsView } from "@/components/projects/roadmap/EpicsView";
import { GanttTimelineView } from "@/components/projects/roadmap/GanttTimelineView";
import { RoadmapSummaryCards } from "@/components/projects/roadmap/RoadmapSummaryCards";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export default function RoadmapPage({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = use(params);
  const router = useRouter();
  const [project, setProject] = useState<Project | null>(null);
  const [summary, setSummary] = useState<RoadmapSummary | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [projectData, summaryData] = await Promise.all([
          ProjectService.getProjectById(resolvedParams.id),
          RoadmapService.getRoadmapSummary(resolvedParams.id).catch(e => {
            console.warn("Could not fetch roadmap summary, returning null", e);
            return null;
          }),
        ]);
        setProject(projectData);
        if (summaryData) {
            setSummary(summaryData);
        }
      } catch (error) {
        console.error(error);
        toast.error("Erro ao carregar projeto");
        router.push("/dashboard");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [resolvedParams.id, router]);

  const handleExport = (format: string) => {
    const url = RoadmapService.getExportUrl(resolvedParams.id, format);
    // Assuming the backend endpoint returns the file as an attachment
    window.open(url, "_blank");
  };

  if (loading || !project) {
    return (
      <PrivateRoute>
        <div className="flex h-screen items-center justify-center bg-background text-foreground">
          <div className="animate-pulse flex flex-col items-center gap-4">
            <div className="h-8 w-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
            <p className="text-muted-foreground">Carregando roadmap...</p>
          </div>
        </div>
      </PrivateRoute>
    );
  }

  return (
    <PrivateRoute>
      <div className="min-h-screen bg-background text-foreground p-6">
        <div className="max-w-7xl mx-auto space-y-8">

          {/* Header */}
          <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-card/50 p-6 rounded-xl border border-border backdrop-blur-sm">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Link href={`/projects/${project.id}`}>
                    <Button variant="ghost" size="icon" className="h-8 w-8">
                        <ArrowLeft className="h-4 w-4" />
                    </Button>
                </Link>
                <h1 className="text-3xl font-bold tracking-tight text-primary">Roadmap: {project.name}</h1>
              </div>
              <p className="text-muted-foreground ml-10">
                Acompanhamento e planejamento estratégico de tarefas e épicos.
              </p>
            </div>

            <div className="flex items-center gap-3">
              <DropdownMenu>
                <DropdownMenuTrigger className={cn(buttonVariants({ variant: "outline" }), "gap-2 border-border/50 hover:bg-card")}>
                        <Download className="h-4 w-4 text-primary" />
                        Exportar Roadmap
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-48 bg-card border-border">
                  <DropdownMenuItem onClick={() => handleExport("json")} className="cursor-pointer hover:bg-muted">JSON</DropdownMenuItem>
                  <DropdownMenuItem onClick={() => handleExport("csv")} className="cursor-pointer hover:bg-muted">CSV</DropdownMenuItem>
                  <DropdownMenuItem onClick={() => handleExport("markdown")} className="cursor-pointer hover:bg-muted">Markdown</DropdownMenuItem>
                  <DropdownMenuItem onClick={() => handleExport("jira")} className="cursor-pointer hover:bg-muted">CSV Jira</DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>

          {/* Metrics Summary */}
          {summary && <RoadmapSummaryCards summary={summary} />}

          {/* Tabs View */}
          <Tabs defaultValue="kanban" className="w-full">
            <TabsList className="grid w-full grid-cols-3 bg-muted/50 p-1 rounded-xl h-14">
              <TabsTrigger value="kanban" className="rounded-lg data-[state=active]:bg-primary data-[state=active]:text-primary-foreground gap-2">
                <LayoutTemplate className="h-4 w-4" /> KanBan
              </TabsTrigger>
              <TabsTrigger value="epics" className="rounded-lg data-[state=active]:bg-primary data-[state=active]:text-primary-foreground gap-2">
                <Network className="h-4 w-4" /> Épicos & Dependências
              </TabsTrigger>
              <TabsTrigger value="timeline" className="rounded-lg data-[state=active]:bg-primary data-[state=active]:text-primary-foreground gap-2">
                <Clock className="h-4 w-4" /> Timeline (Gantt)
              </TabsTrigger>
            </TabsList>

            <TabsContent value="kanban" className="mt-6">
              <AdvancedKanbanBoard projectId={project.id} />
            </TabsContent>

            <TabsContent value="epics" className="mt-6">
              <EpicsView projectId={project.id} />
            </TabsContent>

            <TabsContent value="timeline" className="mt-6">
              <GanttTimelineView projectId={project.id} />
            </TabsContent>
          </Tabs>

        </div>
      </div>
    </PrivateRoute>
  );
}
