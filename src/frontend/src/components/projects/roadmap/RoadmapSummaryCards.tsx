import { RoadmapSummary } from "@/types/roadmap";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Clock, Layers, GitBranch, LayoutGrid } from "lucide-react";

interface RoadmapSummaryCardsProps {
  summary: RoadmapSummary;
}

export function RoadmapSummaryCards({ summary }: RoadmapSummaryCardsProps) {
  const totalTasks = Object.values(summary.tasks_by_type || {}).reduce((acc, val) => acc + val, 0);
  const totalHours = Object.values(summary.estimated_hours_by_type || {}).reduce((acc, val) => acc + val, 0);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <Card className="bg-card/50 border-border shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">Progresso e Escopo</CardTitle>
          <Layers className="h-4 w-4 text-primary" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{summary.total_phases} Fases</div>
          <p className="text-xs text-muted-foreground mt-1">
            {summary.total_epics} Épicos • {totalTasks} Tasks
          </p>
        </CardContent>
      </Card>

      <Card className="bg-card/50 border-border shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">Esforço Total</CardTitle>
          <Clock className="h-4 w-4 text-primary" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{totalHours}h estimadas</div>
          <p className="text-xs text-muted-foreground mt-1">
            Distribuídas por {Object.keys(summary.estimated_hours_by_type || {}).length} tipos de tasks
          </p>
        </CardContent>
      </Card>

      <Card className="bg-card/50 border-border shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">Path Crítico</CardTitle>
          <GitBranch className="h-4 w-4 text-destructive" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-destructive">{summary.critical_path_hours}h</div>
          <p className="text-xs text-muted-foreground mt-1">
            Soma de horas das tasks CRITICAL
          </p>
        </CardContent>
      </Card>

      <Card className="bg-card/50 border-border shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">Distribuição Principal</CardTitle>
          <LayoutGrid className="h-4 w-4 text-secondary" />
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-1 text-sm">
             <div className="flex justify-between items-center">
                 <span className="text-muted-foreground">Backend:</span>
                 <span className="font-medium">{summary.tasks_by_type?.BACKEND || 0} tasks</span>
             </div>
             <div className="flex justify-between items-center">
                 <span className="text-muted-foreground">Frontend:</span>
                 <span className="font-medium">{summary.tasks_by_type?.FRONTEND || 0} tasks</span>
             </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
