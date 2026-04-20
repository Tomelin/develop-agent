import React from 'react';
import { Task, TaskType, TaskComplexity } from '@/types/task';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Clock, GripVertical } from 'lucide-react';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';

interface TaskCardProps {
  task: Task;
  isDragging?: boolean;
}

export function TaskCard({ task, isDragging }: TaskCardProps) {
  const getTypeColor = (type: TaskType) => {
    switch (type) {
      case 'FRONTEND': return 'bg-blue-500/10 text-blue-500 border-blue-500/20';
      case 'BACKEND': return 'bg-green-500/10 text-green-500 border-green-500/20';
      case 'INFRA': return 'bg-purple-500/10 text-purple-500 border-purple-500/20';
      case 'TEST': return 'bg-orange-500/10 text-orange-500 border-orange-500/20';
      case 'DOC': return 'bg-gray-500/10 text-gray-500 border-gray-500/20';
    }
  };

  const getComplexityIndicator = (complexity: TaskComplexity) => {
    switch (complexity) {
      case 'LOW': return <div className="flex gap-0.5"><div className="w-1.5 h-3 bg-green-500 rounded-sm"></div><div className="w-1.5 h-3 bg-muted rounded-sm"></div><div className="w-1.5 h-3 bg-muted rounded-sm"></div></div>;
      case 'MEDIUM': return <div className="flex gap-0.5"><div className="w-1.5 h-3 bg-yellow-500 rounded-sm"></div><div className="w-1.5 h-3 bg-yellow-500 rounded-sm"></div><div className="w-1.5 h-3 bg-muted rounded-sm"></div></div>;
      case 'HIGH': return <div className="flex gap-0.5"><div className="w-1.5 h-3 bg-orange-500 rounded-sm"></div><div className="w-1.5 h-3 bg-orange-500 rounded-sm"></div><div className="w-1.5 h-3 bg-orange-500 rounded-sm"></div></div>;
      case 'CRITICAL': return <div className="flex gap-0.5"><div className="w-1.5 h-3 bg-destructive rounded-sm"></div><div className="w-1.5 h-3 bg-destructive rounded-sm"></div><div className="w-1.5 h-3 bg-destructive rounded-sm animate-pulse"></div></div>;
    }
  };

  return (
    <Card className={`group relative bg-card border-border hover:border-primary/50 transition-colors ${isDragging ? 'shadow-xl scale-105 border-primary z-50' : 'hover:shadow-md'}`}>
      <div className="absolute left-2 top-1/2 -translate-y-1/2 text-muted/50 group-hover:text-muted-foreground transition-colors cursor-grab active:cursor-grabbing">
        <GripVertical className="h-4 w-4" />
      </div>
      <CardContent className="p-4 pl-8">
        <div className="flex justify-between items-start mb-2">
          <Badge variant="outline" className={`text-[10px] uppercase ${getTypeColor(task.type)}`}>
            {task.type}
          </Badge>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger>
                {getComplexityIndicator(task.complexity)}
              </TooltipTrigger>
              <TooltipContent>
                <p>Complexidade: {task.complexity}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>

        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger >
              <h4 className="text-sm font-semibold mb-2 line-clamp-2 text-left cursor-help">{task.title}</h4>
            </TooltipTrigger>
            <TooltipContent className="max-w-[250px]">
              <p className="text-sm">{task.description}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>

        <div className="flex items-center justify-between mt-4">
          <div className="flex items-center text-xs text-muted-foreground bg-muted/50 px-2 py-1 rounded-md">
            <Clock className="h-3 w-3 mr-1" />
            {task.estimated_hours}h
          </div>

          {task.assigned_agent_id && (
            <div className="h-6 w-6 rounded-full bg-primary/20 border border-primary/30 flex items-center justify-center" title="Atribuído a um Agente">
               <span className="text-[10px] text-primary">IA</span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}