import { useState, useEffect, useRef } from 'react';
import { ProjectService } from '@/services/project';
import { RoadmapTask } from '@/types/task';
import { Download } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface GanttTimelineViewProps {
  projectId: string;
}

export function GanttTimelineView({ projectId }: GanttTimelineViewProps) {
  const [tasks, setTasks] = useState<RoadmapTask[]>([]);
  const [loading, setLoading] = useState(true);
  const ganttRef = useRef<HTMLDivElement>(null);

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

  const handleExportPNG = () => {
    // This is a placeholder for actual PNG export logic.
    // In a real scenario, we'd use html2canvas or similar to snapshot ganttRef.current
    alert("Exportação de PNG será implementada integrando html2canvas.");
  };

  if (loading) {
    return <div className="animate-pulse h-64 bg-card/30 rounded-xl" />;
  }

  // Basic Gantt Logic:
  // 1. Group tasks by Epic to allow visual "parallelism"
  // 2. Order by ID or creation to simulate timeline progression
  // 3. Color by complexity

  const getComplexityColor = (complexity: string) => {
    switch (complexity) {
      case 'CRITICAL': return 'bg-destructive/80 border-destructive';
      case 'HIGH': return 'bg-orange-500/80 border-orange-600';
      case 'MEDIUM': return 'bg-yellow-500/80 border-yellow-600';
      case 'LOW': return 'bg-green-500/80 border-green-600';
      default: return 'bg-primary/80 border-primary';
    }
  };

  const groupedByEpic = tasks.reduce((acc, task) => {
    const epic = task.epic_id || 'Geral';
    if (!acc[epic]) acc[epic] = [];
    acc[epic].push(task);
    return acc;
  }, {} as Record<string, RoadmapTask[]>);

  // Calculate generic block width based on hours (1h = 10px roughly, min 40px, max 300px)
  const getWidth = (hours: number) => {
      const w = hours * 15;
      return Math.max(80, Math.min(w, 300));
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center bg-card/50 p-4 rounded-xl border border-border">
          <div className="text-sm text-muted-foreground flex gap-4">
              <span className="flex items-center gap-1"><div className="w-3 h-3 rounded-full bg-destructive"></div> Crítica</span>
              <span className="flex items-center gap-1"><div className="w-3 h-3 rounded-full bg-orange-500"></div> Alta</span>
              <span className="flex items-center gap-1"><div className="w-3 h-3 rounded-full bg-yellow-500"></div> Média</span>
              <span className="flex items-center gap-1"><div className="w-3 h-3 rounded-full bg-green-500"></div> Baixa</span>
          </div>
          <Button variant="outline" size="sm" onClick={handleExportPNG} className="gap-2">
            <Download className="h-4 w-4" /> Exportar PNG
          </Button>
      </div>

      <div
        ref={ganttRef}
        className="bg-card/30 border border-border rounded-xl p-6 overflow-x-auto"
      >
        <div className="min-w-[800px] space-y-8">
            {Object.entries(groupedByEpic).map(([epicId, epicTasks]) => (
                <div key={epicId} className="relative">
                    <h3 className="text-sm font-bold text-muted-foreground mb-4 sticky left-0 uppercase tracking-wider">{epicId}</h3>

                    {/* Simulated Parallelism: Render tasks side by side if possible, or staggered */}
                    <div className="flex gap-4 flex-wrap relative pb-4">
                        {/* A line connecting them loosely */}
                        <div className="absolute top-1/2 left-0 w-full h-[1px] bg-border/50 -z-10"></div>

                        {epicTasks.map((task) => (
                            <div
                                key={task.id}
                                className={`h-12 rounded-md border text-xs text-foreground p-2 flex flex-col justify-center overflow-hidden whitespace-nowrap text-ellipsis cursor-help shadow-sm transition-transform hover:scale-105 ${getComplexityColor(task.complexity)}`}
                                style={{ width: `${getWidth(task.estimated_hours)}px` }}
                                title={`${task.title}\nHoras: ${task.estimated_hours}h\nTipo: ${task.type}`}
                            >
                                <span className="font-semibold truncate text-white mix-blend-difference drop-shadow-md">{task.title}</span>
                                <span className="text-[10px] text-white/80 mix-blend-difference drop-shadow-md">{task.estimated_hours}h</span>
                            </div>
                        ))}
                    </div>
                </div>
            ))}
            {Object.keys(groupedByEpic).length === 0 && (
                <div className="text-center py-12 text-muted-foreground">
                    Nenhuma task para gerar a timeline.
                </div>
            )}
        </div>
      </div>
    </div>
  );
}
