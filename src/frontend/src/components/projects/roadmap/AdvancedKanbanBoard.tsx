import { useState, useEffect } from 'react';
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
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { TaskStatus, TaskComplexity, TaskType, RoadmapTask } from '@/types/task';
import { ProjectService } from '@/services/project';
import { SortableTaskCard } from '@/components/projects/SortableTaskCard';
import { TaskCard } from '@/components/projects/TaskCard';
import { Badge } from '@/components/ui/badge';
import { LayoutGrid, Search } from 'lucide-react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { toast } from 'sonner';

interface AdvancedKanbanBoardProps {
  projectId: string;
}

const COLUMNS: { id: TaskStatus; title: string }[] = [
  { id: 'TODO', title: 'To Do' },
  { id: 'IN_PROGRESS', title: 'Em Progresso' },
  { id: 'REVIEW', title: 'Revisão' },
  { id: 'BLOCKED', title: 'Bloqueado' },
  { id: 'DONE', title: 'Concluído' },
];

const TASK_TYPE_VALUES: Array<TaskType | 'ALL'> = ['ALL', 'FRONTEND', 'BACKEND', 'INFRA', 'TEST', 'DOC'];
const TASK_COMPLEXITY_VALUES: Array<TaskComplexity | 'ALL'> = ['ALL', 'LOW', 'MEDIUM', 'HIGH', 'CRITICAL'];

const isTaskTypeFilter = (value: string): value is TaskType | 'ALL' =>
  TASK_TYPE_VALUES.includes(value as TaskType | 'ALL');

const isTaskComplexityFilter = (value: string): value is TaskComplexity | 'ALL' =>
  TASK_COMPLEXITY_VALUES.includes(value as TaskComplexity | 'ALL');

export function AdvancedKanbanBoard({ projectId }: AdvancedKanbanBoardProps) {
  const [tasks, setTasks] = useState<RoadmapTask[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTask, setActiveTask] = useState<RoadmapTask | null>(null);

  // Filters
  const [search, setSearch] = useState('');
  const [typeFilter, setTypeFilter] = useState<TaskType | 'ALL'>('ALL');
  const [complexityFilter, setComplexityFilter] = useState<TaskComplexity | 'ALL'>('ALL');
  const [epicFilter, setEpicFilter] = useState<string>('ALL');

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
        const response = await ProjectService.getTasks(projectId, 1, 500);
        setTasks(response.items as RoadmapTask[] || []);
      } catch {
        console.error("Failed to fetch tasks");
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, [projectId]);

  const uniqueEpics = Array.from(new Set(tasks.map(t => t.epic_id).filter(Boolean))) as string[];

  const filteredTasks = tasks.filter(t => {
    const matchesSearch = t.title.toLowerCase().includes(search.toLowerCase()) ||
                          t.description.toLowerCase().includes(search.toLowerCase());
    const matchesType = typeFilter === 'ALL' || t.type === typeFilter;
    const matchesComplexity = complexityFilter === 'ALL' || t.complexity === complexityFilter;
    const matchesEpic = epicFilter === 'ALL' || t.epic_id === epicFilter;

    return matchesSearch && matchesType && matchesComplexity && matchesEpic;
  });

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

    const activeTask = tasks.find((t) => t.id === activeId);
    if (!activeTask) return;

    const activeStatus = activeTask.status;
    let overStatus = overId as TaskStatus;

    if (!COLUMNS.find(c => c.id === overId)) {
        const overTask = tasks.find((t) => t.id === overId);
        if (overTask) overStatus = overTask.status;
    }

    if (activeStatus !== overStatus) {
      // Validations (Frontend pre-check before backend call)
      if (overStatus === 'DONE' && activeTask.dependencies && activeTask.dependencies.length > 0) {
          const unresolvedDeps = activeTask.dependencies.filter(depId => {
              const dep = tasks.find(t => t.id === depId);
              return dep && dep.status !== 'DONE';
          });
          if (unresolvedDeps.length > 0) {
              toast.error("Não é possível mover para DONE. Existem dependências pendentes.");
              return;
          }
      }

      setTasks((prev) =>
        prev.map(t => t.id === activeId ? { ...t, status: overStatus } : t)
      );

      try {
        await ProjectService.updateTaskStatus(projectId, activeId.toString(), overStatus);
        toast.success(`Status atualizado para ${COLUMNS.find(c => c.id === overStatus)?.title}`);
      } catch {
        toast.error("Erro ao atualizar status na API. Revertendo...");
        setTasks((prev) =>
          prev.map(t => t.id === activeId ? { ...t, status: activeStatus } : t)
        );
      }
    }
  };

  if (loading) {
    return <div className="animate-pulse h-96 bg-card/30 rounded-xl" />;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-4 p-4 bg-card/50 rounded-xl border border-border">
        <div className="flex-1 min-w-[200px] relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
                placeholder="Buscar tasks..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9 bg-background"
            />
        </div>

        <div className="flex items-center gap-2">
            <LayoutGrid className="h-4 w-4 text-muted-foreground hidden sm:block" />

            <Select value={typeFilter} onValueChange={(v: string | null) => v && isTaskTypeFilter(v) && setTypeFilter(v)}>
                <SelectTrigger className="w-[140px] bg-background">
                <SelectValue placeholder="Tipo" />
                </SelectTrigger>
                <SelectContent>
                <SelectItem value="ALL">Todos Tipos</SelectItem>
                <SelectItem value="FRONTEND">Frontend</SelectItem>
                <SelectItem value="BACKEND">Backend</SelectItem>
                <SelectItem value="INFRA">Infra</SelectItem>
                <SelectItem value="TEST">Testes</SelectItem>
                <SelectItem value="DOC">Docs</SelectItem>
                </SelectContent>
            </Select>

            <Select value={complexityFilter} onValueChange={(v: string | null) => v && isTaskComplexityFilter(v) && setComplexityFilter(v)}>
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

            <Select value={epicFilter} onValueChange={(v: string | null) => v && setEpicFilter(v)}>
                <SelectTrigger className="w-[150px] bg-background">
                <SelectValue placeholder="Épico" />
                </SelectTrigger>
                <SelectContent>
                <SelectItem value="ALL">Todos Épicos</SelectItem>
                {uniqueEpics.map(epicId => (
                    <SelectItem key={epicId} value={epicId}>Épico {epicId.substring(0,6)}</SelectItem>
                ))}
                </SelectContent>
            </Select>
        </div>
      </div>

      <div className="flex gap-6 overflow-x-auto pb-4 min-h-[600px]">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          {COLUMNS.map((col) => {
            const columnTasks = getTasksByStatus(col.id);
            const totalHours = columnTasks.reduce((acc, task) => acc + (task.estimated_hours || 0), 0);

            return (
              <div key={col.id} className="flex flex-col min-w-[320px] max-w-[320px] bg-card/30 rounded-xl border border-border">
                <div className="p-4 border-b border-border/50 flex flex-col gap-2 bg-card/50 rounded-t-xl">
                  <div className="flex justify-between items-center">
                    <h3 className="font-semibold">{col.title}</h3>
                    <Badge variant="secondary" className="bg-background">
                        {columnTasks.length}
                    </Badge>
                  </div>
                  <div className="text-xs text-muted-foreground">
                      {totalHours} horas estimadas
                  </div>
                </div>

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
