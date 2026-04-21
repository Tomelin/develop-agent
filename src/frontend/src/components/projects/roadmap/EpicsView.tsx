import { useState, useEffect } from 'react';
import { ProjectService } from '@/services/project';
import { RoadmapTask } from '@/types/task';

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Network, CheckCircle2, Circle, AlertCircle, ArrowRight } from 'lucide-react';
import { Progress } from '@/components/ui/progress';

interface EpicsViewProps {
  projectId: string;
}

export function EpicsView({ projectId }: EpicsViewProps) {
  const [tasks, setTasks] = useState<RoadmapTask[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await ProjectService.getTasks(projectId, 1, 500);
        setTasks(response.items as RoadmapTask[] || []);
      } catch (error) {
        console.error("Failed to fetch tasks", error);
      } finally {
        setLoading(false);
      }
    };
    fetchTasks();
  }, [projectId]);

  if (loading) {
    return <div className="animate-pulse h-64 bg-card/30 rounded-xl" />;
  }

  // Agrupamento manual caso o backend não entregue as phases/epics pre-aninhadas
  const groupedByPhase: Record<string, Record<string, RoadmapTask[]>> = {};

  tasks.forEach(task => {
      const phaseId = task.phase_id || "Fase Desconhecida";
      const epicId = task.epic_id || "Épico Geral";

      if (!groupedByPhase[phaseId]) groupedByPhase[phaseId] = {};
      if (!groupedByPhase[phaseId][epicId]) groupedByPhase[phaseId][epicId] = [];

      groupedByPhase[phaseId][epicId].push(task);
  });

  const getTaskStatusIcon = (status: string, isBlocked: boolean) => {
      if (isBlocked) return <AlertCircle className="h-4 w-4 text-destructive" />;
      if (status === 'DONE') return <CheckCircle2 className="h-4 w-4 text-primary" />;
      if (status === 'IN_PROGRESS') return <Circle className="h-4 w-4 text-blue-500 fill-blue-500/20" />;
      return <Circle className="h-4 w-4 text-muted-foreground" />;
  };

  return (
    <div className="space-y-6">
        {Object.entries(groupedByPhase).map(([phaseId, epics]) => (
            <div key={phaseId} className="bg-card/30 border border-border rounded-xl p-4">
                <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                    <span className="text-primary text-2xl font-black">#</span>
                    Fase {phaseId.substring(0,8)}
                </h2>

                <Accordion className="w-full space-y-4">
                    {Object.entries(epics).map(([epicId, epicTasks]) => {
                        const doneTasks = epicTasks.filter(t => t.status === 'DONE').length;
                        const progress = epicTasks.length > 0 ? (doneTasks / epicTasks.length) * 100 : 0;
                        const epicStatus = progress === 0 ? 'PENDING' : (progress === 100 ? 'COMPLETED' : 'IN_PROGRESS');

                        return (
                            <AccordionItem key={epicId} value={epicId} className="border border-border bg-card/50 rounded-lg px-4 data-[state=open]:bg-card/80">
                                <AccordionTrigger className="hover:no-underline">
                                    <div className="flex flex-col sm:flex-row sm:items-center justify-between w-full pr-4 gap-4">
                                        <div className="flex items-center gap-3">
                                            <Network className="h-5 w-5 text-secondary" />
                                            <div className="text-left">
                                                <h3 className="font-semibold text-base">{epicId}</h3>
                                                <p className="text-xs text-muted-foreground">{epicTasks.length} tasks</p>
                                            </div>
                                        </div>

                                        <div className="flex items-center gap-4 w-full sm:w-64">
                                            <Badge variant={epicStatus === 'COMPLETED' ? 'default' : 'secondary'} className="text-[10px]">
                                                {epicStatus}
                                            </Badge>
                                            <div className="flex-1">
                                                <div className="flex justify-between text-xs mb-1">
                                                    <span>Progresso</span>
                                                    <span>{Math.round(progress)}%</span>
                                                </div>
                                                <Progress value={progress} className="h-2" />
                                            </div>
                                        </div>
                                    </div>
                                </AccordionTrigger>
                                <AccordionContent className="pt-4 pb-6">
                                    <div className="space-y-3">
                                        {epicTasks.map(task => {
                                            // Check if blocked by dependencies
                                            const hasPendingDeps = task.dependencies?.some(depId => {
                                                const depTask = tasks.find(t => t.id === depId);
                                                return depTask && depTask.status !== 'DONE';
                                            });

                                            return (
                                                <Card key={task.id} className={`bg-background/50 border-border/50 ${hasPendingDeps ? 'border-destructive/30 bg-destructive/5' : ''}`}>
                                                    <CardContent className="p-4 flex items-center justify-between gap-4">
                                                        <div className="flex items-center gap-3 flex-1">
                                                            {getTaskStatusIcon(task.status, !!hasPendingDeps)}
                                                            <div>
                                                                <div className="font-medium text-sm flex items-center gap-2">
                                                                    {task.title}
                                                                    {hasPendingDeps && (
                                                                        <Badge variant="destructive" className="text-[10px] h-4 px-1">BLOQUEADA</Badge>
                                                                    )}
                                                                </div>
                                                                <div className="text-xs text-muted-foreground mt-1">
                                                                    {task.type} • {task.complexity} • {task.estimated_hours}h
                                                                </div>
                                                            </div>
                                                        </div>

                                                        {task.dependencies && task.dependencies.length > 0 && (
                                                            <div className="hidden sm:flex items-center gap-2 text-xs text-muted-foreground bg-muted/50 px-3 py-1.5 rounded-md">
                                                                <ArrowRight className="h-3 w-3" />
                                                                Depende de {task.dependencies.length} tasks
                                                            </div>
                                                        )}
                                                    </CardContent>
                                                </Card>
                                            );
                                        })}
                                    </div>
                                </AccordionContent>
                            </AccordionItem>
                        );
                    })}
                </Accordion>
            </div>
        ))}
        {Object.keys(groupedByPhase).length === 0 && (
            <div className="text-center py-12 text-muted-foreground">
                Nenhum épico encontrado para este projeto.
            </div>
        )}
    </div>
  );
}
