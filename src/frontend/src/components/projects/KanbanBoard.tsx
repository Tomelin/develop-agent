"use client";

import React, { useState, useEffect } from 'react';
import {
  DndContext,
  DragOverlay,
  closestCorners,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragStartEvent,
  DragEndEvent,
} from '@dnd-kit/core';
import {
  SortableContext,
  arrayMove,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { Task, TaskStatus, TaskComplexity, TaskType } from '@/types/task';
import { ProjectService } from '@/services/project';
import { SortableTaskCard } from './SortableTaskCard';
import { TaskCard } from './TaskCard';
import { Badge } from '@/components/ui/badge';
import { Download, Filter, LayoutGrid } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { toast } from 'sonner';

interface KanbanBoardProps {
  projectId: string;
}

const COLUMNS: { id: TaskStatus; title: string }[] = [
  { id: 'TODO', title: 'To Do' },
  { id: 'IN_PROGRESS', title: 'Em Progresso' },
  { id: 'BLOCKED', title: 'Bloqueado' },
  { id: 'DONE', title: 'Concluído' },
];

export function KanbanBoard({ projectId }: KanbanBoardProps) {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTask, setActiveTask] = useState<Task | null>(null);

  // Filters
  const [typeFilter, setTypeFilter] = useState<TaskType | 'ALL'>('ALL');
  const [complexityFilter, setComplexityFilter] = useState<TaskComplexity | 'ALL'>('ALL');

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 5,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await ProjectService.getTasks(projectId, 1, 100);
        setTasks(response.items || []);
      } catch (error) {
        console.error("Failed to fetch tasks", error);
        // We will seed mock data for visual showcase if API is not fully returning tasks yet
        setTasks([
          {
            id: "t1", project_id: projectId, phase_id: "p1", title: "Setup React Router", description: "Configurar rotas base no frontend",
            type: "FRONTEND", complexity: "LOW", estimated_hours: 4, status: "DONE", created_at: new Date().toISOString(), updated_at: new Date().toISOString()
          },
          {
            id: "t2", project_id: projectId, phase_id: "p1", title: "API Authentication", description: "Implementar JWT no backend Gin",
            type: "BACKEND", complexity: "HIGH", estimated_hours: 8, status: "IN_PROGRESS", created_at: new Date().toISOString(), updated_at: new Date().toISOString()
          },
          {
            id: "t3", project_id: projectId, phase_id: "p1", title: "Database Schema", description: "Modelagem MongoDB collections",
            type: "INFRA", complexity: "MEDIUM", estimated_hours: 6, status: "TODO", created_at: new Date().toISOString(), updated_at: new Date().toISOString()
          },
          {
            id: "t4", project_id: projectId, phase_id: "p1", title: "Docker Compose", description: "Configurar ambiente local",
            type: "INFRA", complexity: "LOW", estimated_hours: 2, status: "BLOCKED", created_at: new Date().toISOString(), updated_at: new Date().toISOString()
          }
        ]);
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, [projectId]);

  const filteredTasks = tasks.filter(t =>
    (typeFilter === 'ALL' || t.type === typeFilter) &&
    (complexityFilter === 'ALL' || t.complexity === complexityFilter)
  );

  const getTasksByStatus = (status: TaskStatus) => {
    return filteredTasks.filter((t) => t.status === status);
  };

  const handleDragStart = (event: DragStartEvent) => {
    const { active } = event;
    const task = tasks.find((t) => t.id === active.id);
    if (task) setActiveTask(task);
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveTask(null);

    if (!over) return;

    const activeId = active.id;
    const overId = over.id;

    // Find the task and the target column
    const activeTask = tasks.find((t) => t.id === activeId);
    if (!activeTask) return;

    const activeStatus = activeTask.status;
    let overStatus = overId as TaskStatus;

    // If dropped over a task, get that task's status instead
    if (!COLUMNS.find(c => c.id === overId)) {
        const overTask = tasks.find((t) => t.id === overId);
        if (overTask) overStatus = overTask.status;
    }

    if (activeStatus !== overStatus) {
      // Optimistic UI update
      setTasks((prev) =>
        prev.map(t => t.id === activeId ? { ...t, status: overStatus } : t)
      );

      try {
        await ProjectService.updateTaskStatus(projectId, activeId.toString(), overStatus);
        toast.success(`Task movida para ${COLUMNS.find(c => c.id === overStatus)?.title}`);
      } catch (error) {
        toast.error("Erro ao atualizar status da task");
        // Revert on error
        setTasks((prev) =>
          prev.map(t => t.id === activeId ? { ...t, status: activeStatus } : t)
        );
      }
    }
  };

  const exportJSON = () => {
    const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(tasks, null, 2));
    const downloadAnchorNode = document.createElement('a');
    downloadAnchorNode.setAttribute("href",     dataStr);
    downloadAnchorNode.setAttribute("download", `project-${projectId}-tasks.json`);
    document.body.appendChild(downloadAnchorNode); // required for firefox
    downloadAnchorNode.click();
    downloadAnchorNode.remove();
  };

  if (loading) {
    return <div className="animate-pulse h-96 bg-card/30 rounded-xl" />;
  }

  return (
    <div className="space-y-6">
      {/* Toolbar */}
      <div className="flex flex-wrap items-center justify-between gap-4 p-4 bg-card/50 rounded-xl border border-border">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <LayoutGrid className="h-4 w-4 text-muted-foreground" />
            <span className="font-medium">Filtros:</span>
          </div>

          <Select value={typeFilter} onValueChange={(v: string | null) => v && setTypeFilter(v as any)}>
            <SelectTrigger className="w-[140px] bg-background">
              <SelectValue placeholder="Tipo" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ALL">Todos os Tipos</SelectItem>
              <SelectItem value="FRONTEND">Frontend</SelectItem>
              <SelectItem value="BACKEND">Backend</SelectItem>
              <SelectItem value="INFRA">Infra</SelectItem>
              <SelectItem value="TEST">Testes</SelectItem>
              <SelectItem value="DOC">Docs</SelectItem>
            </SelectContent>
          </Select>

          <Select value={complexityFilter} onValueChange={(v: string | null) => v && setComplexityFilter(v as any)}>
            <SelectTrigger className="w-[150px] bg-background">
              <SelectValue placeholder="Complexidade" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ALL">Qualquer Compl.</SelectItem>
              <SelectItem value="LOW">Baixa</SelectItem>
              <SelectItem value="MEDIUM">Média</SelectItem>
              <SelectItem value="HIGH">Alta</SelectItem>
              <SelectItem value="CRITICAL">Crítica</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <Button variant="outline" size="sm" onClick={exportJSON} className="gap-2">
          <Download className="h-4 w-4" /> Exportar JSON
        </Button>
      </div>

      {/* Board */}
      <div className="flex gap-6 overflow-x-auto pb-4 min-h-[600px]">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          {COLUMNS.map((col) => {
            const columnTasks = getTasksByStatus(col.id);
            return (
              <div key={col.id} className="flex flex-col min-w-[300px] max-w-[300px] bg-card/30 rounded-xl border border-border">
                <div className="p-4 border-b border-border/50 flex justify-between items-center bg-card/50 rounded-t-xl">
                  <h3 className="font-semibold">{col.title}</h3>
                  <Badge variant="secondary" className="bg-background">
                    {columnTasks.length}
                  </Badge>
                </div>

                {/* Dropzone Column */}
                <div className="flex-1 p-3">
                  <SortableContext
                    id={col.id}
                    items={columnTasks.map((t) => t.id)}
                    strategy={verticalListSortingStrategy}
                  >
                    <div className="space-y-3 min-h-[150px]">
                      {columnTasks.map((task) => (
                        <SortableTaskCard key={task.id} task={task} />
                      ))}
                    </div>
                  </SortableContext>
                </div>
              </div>
            );
          })}

          <DragOverlay>
            {activeTask ? <TaskCard task={activeTask} isDragging /> : null}
          </DragOverlay>
        </DndContext>
      </div>
    </div>
  );
}