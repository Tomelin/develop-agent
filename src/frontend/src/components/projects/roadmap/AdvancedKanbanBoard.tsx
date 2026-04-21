import { useEffect, useMemo, useState, type ReactNode } from 'react';
import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  KeyboardSensor,
  PointerSensor,
  closestCorners,
  useDroppable,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import { SortableContext, sortableKeyboardCoordinates, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { LayoutGrid, Search, SlidersHorizontal } from 'lucide-react';

import { SortableTaskCard } from '@/components/projects/SortableTaskCard';
import { TaskCard } from '@/components/projects/TaskCard';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ProjectService } from '@/services/project';
import { TaskComplexity, TaskStatus, TaskType, RoadmapTask } from '@/types/task';
import { toast } from 'sonner';

interface AdvancedKanbanBoardProps {
  projectId: string;
}

const COLUMNS: { id: TaskStatus; title: string; accent: string }[] = [
  { id: 'TODO', title: 'To Do', accent: 'bg-slate-500/15 text-slate-200' },
  { id: 'IN_PROGRESS', title: 'Em Progresso', accent: 'bg-blue-500/15 text-blue-300' },
  { id: 'REVIEW', title: 'Revisão', accent: 'bg-violet-500/15 text-violet-300' },
  { id: 'DONE', title: 'Concluído', accent: 'bg-emerald-500/15 text-emerald-300' },
  { id: 'BLOCKED', title: 'Bloqueado', accent: 'bg-red-500/15 text-red-300' },
];

const TASK_TYPE_VALUES: Array<TaskType | 'ALL'> = ['ALL', 'FRONTEND', 'BACKEND', 'INFRA', 'TEST', 'DOC'];
const TASK_COMPLEXITY_VALUES: Array<TaskComplexity | 'ALL'> = ['ALL', 'LOW', 'MEDIUM', 'HIGH', 'CRITICAL'];

const isTaskTypeFilter = (value: string): value is TaskType | 'ALL' =>
  TASK_TYPE_VALUES.includes(value as TaskType | 'ALL');

const isTaskComplexityFilter = (value: string): value is TaskComplexity | 'ALL' =>
  TASK_COMPLEXITY_VALUES.includes(value as TaskComplexity | 'ALL');

function DroppableColumn({
  id,
  children,
  isOver,
}: {
  id: TaskStatus;
  children: ReactNode;
  isOver: boolean;
}) {
  return (
    <div
      data-column-id={id}
      className={`flex-1 p-3 min-h-[180px] rounded-b-xl transition-colors ${
        isOver ? 'bg-primary/5' : 'bg-transparent'
      }`}
    >
      {children}
    </div>
  );
}

function ColumnDropZone({
  column,
  tasks,
}: {
  column: { id: TaskStatus; title: string; accent: string };
  tasks: RoadmapTask[];
}) {
  const { isOver, setNodeRef } = useDroppable({ id: column.id });

  return (
    <div
      ref={setNodeRef}
      className='flex flex-col min-w-[340px] max-w-[340px] rounded-xl border border-border bg-card/40 backdrop-blur'
    >
      <div className='rounded-t-xl border-b border-border/50 bg-card/70 p-4'>
        <div className='mb-2 flex items-center justify-between'>
          <h3 className='font-semibold tracking-tight'>{column.title}</h3>
          <Badge className={column.accent}>{tasks.length}</Badge>
        </div>
        <p className='text-xs text-muted-foreground'>
          {tasks.reduce((sum, task) => sum + (task.estimated_hours || 0), 0)} horas estimadas
        </p>
      </div>

      <DroppableColumn id={column.id} isOver={isOver}>
        <SortableContext items={tasks.map((task) => task.id)} strategy={verticalListSortingStrategy}>
          <div className='space-y-3'>
            {tasks.map((task) => (
              <SortableTaskCard key={task.id} task={task} />
            ))}
          </div>
        </SortableContext>
      </DroppableColumn>
    </div>
  );
}

export function AdvancedKanbanBoard({ projectId }: AdvancedKanbanBoardProps) {
  const [tasks, setTasks] = useState<RoadmapTask[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTask, setActiveTask] = useState<RoadmapTask | null>(null);

  const [search, setSearch] = useState('');
  const [typeFilter, setTypeFilter] = useState<TaskType | 'ALL'>('ALL');
  const [complexityFilter, setComplexityFilter] = useState<TaskComplexity | 'ALL'>('ALL');
  const [epicFilter, setEpicFilter] = useState<string>('ALL');
  const [phaseFilter, setPhaseFilter] = useState<string>('ALL');

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  );

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await ProjectService.getTasks(projectId, 1, 500);
        setTasks((response.items as RoadmapTask[]) || []);
      } catch (error) {
        console.error('Failed to fetch tasks', error);
        toast.error('Não foi possível carregar as tasks do roadmap.');
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, [projectId]);

  const uniqueEpics = useMemo(
    () => Array.from(new Set(tasks.map((task) => task.epic_id).filter(Boolean))) as string[],
    [tasks]
  );

  const uniquePhases = useMemo(
    () => Array.from(new Set(tasks.map((task) => task.phase_id).filter(Boolean))),
    [tasks]
  );

  const filteredTasks = useMemo(
    () =>
      tasks.filter((task) => {
        const searchValue = search.trim().toLowerCase();
        const matchesSearch =
          searchValue.length === 0 ||
          task.title.toLowerCase().includes(searchValue) ||
          task.description.toLowerCase().includes(searchValue);
        const matchesType = typeFilter === 'ALL' || task.type === typeFilter;
        const matchesComplexity = complexityFilter === 'ALL' || task.complexity === complexityFilter;
        const matchesEpic = epicFilter === 'ALL' || task.epic_id === epicFilter;
        const matchesPhase = phaseFilter === 'ALL' || task.phase_id === phaseFilter;

        return matchesSearch && matchesType && matchesComplexity && matchesEpic && matchesPhase;
      }),
    [tasks, search, typeFilter, complexityFilter, epicFilter, phaseFilter]
  );

  const tasksByColumn = useMemo(
    () =>
      COLUMNS.reduce<Record<TaskStatus, RoadmapTask[]>>(
        (acc, column) => ({ ...acc, [column.id]: filteredTasks.filter((task) => task.status === column.id) }),
        { TODO: [], IN_PROGRESS: [], REVIEW: [], DONE: [], BLOCKED: [] }
      ),
    [filteredTasks]
  );

  const handleDragStart = (event: DragStartEvent) => {
    const task = tasks.find((item) => item.id === event.active.id);
    if (task) setActiveTask(task);
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveTask(null);

    if (!over) return;

    const activeId = String(active.id);
    const activeTaskRecord = tasks.find((task) => task.id === activeId);
    if (!activeTaskRecord) return;

    const sourceStatus = activeTaskRecord.status;

    let targetStatus = over.id as TaskStatus;
    if (!COLUMNS.some((column) => column.id === over.id)) {
      const overTask = tasks.find((task) => task.id === over.id);
      if (!overTask) return;
      targetStatus = overTask.status;
    }

    if (sourceStatus === targetStatus) return;

    if (targetStatus === 'DONE' && activeTaskRecord.dependencies?.length) {
      const unresolvedDependencies = activeTaskRecord.dependencies.filter((depId) => {
        const dependency = tasks.find((task) => task.id === depId);
        return dependency && dependency.status !== 'DONE';
      });

      if (unresolvedDependencies.length > 0) {
        toast.error('Não é possível concluir a task com dependências pendentes.');
        return;
      }
    }

    setTasks((previousTasks) =>
      previousTasks.map((task) =>
        task.id === activeId
          ? {
              ...task,
              status: targetStatus,
            }
          : task
      )
    );

    try {
      await ProjectService.updateTaskStatus(projectId, activeId, targetStatus);
      toast.success(`Status atualizado para ${COLUMNS.find((column) => column.id === targetStatus)?.title}.`);
    } catch (error) {
      console.error('Failed to update task status', error);
      toast.error('Falha ao persistir status da task. Revertendo alteração.');
      setTasks((previousTasks) =>
        previousTasks.map((task) => (task.id === activeId ? { ...task, status: sourceStatus } : task))
      );
    }
  };

  if (loading) {
    return <div className='h-96 animate-pulse rounded-xl bg-card/30' />;
  }

  return (
    <div className='space-y-5'>
      <div className='rounded-xl border border-border bg-card/50 p-4'>
        <div className='mb-4 flex flex-wrap items-center gap-3'>
          <div className='relative min-w-[240px] flex-1'>
            <Search className='absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground' />
            <Input
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              placeholder='Buscar por título ou descrição...'
              className='bg-background pl-9'
            />
          </div>
          <div className='inline-flex items-center gap-2 rounded-lg border border-border/80 bg-background px-3 py-2 text-xs text-muted-foreground'>
            <SlidersHorizontal className='h-3.5 w-3.5' />
            {filteredTasks.length} de {tasks.length} tasks visíveis
          </div>
        </div>

        <div className='flex flex-wrap items-center gap-2'>
          <LayoutGrid className='hidden h-4 w-4 text-muted-foreground sm:block' />

          <Select value={typeFilter} onValueChange={(value: string | null) => value && isTaskTypeFilter(value) && setTypeFilter(value)}>
            <SelectTrigger className='w-[145px] bg-background'>
              <SelectValue placeholder='Tipo' />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='ALL'>Todos tipos</SelectItem>
              <SelectItem value='FRONTEND'>Frontend</SelectItem>
              <SelectItem value='BACKEND'>Backend</SelectItem>
              <SelectItem value='INFRA'>Infra</SelectItem>
              <SelectItem value='TEST'>Testes</SelectItem>
              <SelectItem value='DOC'>Docs</SelectItem>
            </SelectContent>
          </Select>

          <Select value={complexityFilter} onValueChange={(value: string | null) => value && isTaskComplexityFilter(value) && setComplexityFilter(value)}>
            <SelectTrigger className='w-[160px] bg-background'>
              <SelectValue placeholder='Complexidade' />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='ALL'>Qualquer complexidade</SelectItem>
              <SelectItem value='LOW'>Baixa</SelectItem>
              <SelectItem value='MEDIUM'>Média</SelectItem>
              <SelectItem value='HIGH'>Alta</SelectItem>
              <SelectItem value='CRITICAL'>Crítica</SelectItem>
            </SelectContent>
          </Select>

          <Select value={phaseFilter} onValueChange={(value: string | null) => value && setPhaseFilter(value)}>
            <SelectTrigger className='w-[170px] bg-background'>
              <SelectValue placeholder='Fase' />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='ALL'>Todas fases</SelectItem>
              {uniquePhases.map((phaseId) => (
                <SelectItem key={phaseId} value={phaseId}>
                  {phaseId}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={epicFilter} onValueChange={(value: string | null) => value && setEpicFilter(value)}>
            <SelectTrigger className='w-[170px] bg-background'>
              <SelectValue placeholder='Épico' />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='ALL'>Todos épicos</SelectItem>
              {uniqueEpics.map((epicId) => (
                <SelectItem key={epicId} value={epicId}>
                  {epicId}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <DndContext
        sensors={sensors}
        collisionDetection={closestCorners}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        <div className='flex min-h-[620px] gap-6 overflow-x-auto pb-4'>
          {COLUMNS.map((column) => (
            <ColumnDropZone key={column.id} column={column} tasks={tasksByColumn[column.id]} />
          ))}
        </div>

        <DragOverlay>{activeTask ? <TaskCard task={activeTask} isDragging /> : null}</DragOverlay>
      </DndContext>
    </div>
  );
}
