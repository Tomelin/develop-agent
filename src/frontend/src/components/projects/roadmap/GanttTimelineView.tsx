import { useEffect, useMemo, useRef, useState } from 'react';
import { Download } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { ProjectService } from '@/services/project';
import { RoadmapTask } from '@/types/task';
import { toast } from 'sonner';

interface GanttTimelineViewProps {
  projectId: string;
}

type TimelineBlock = {
  task: RoadmapTask;
  lane: number;
  start: number;
  duration: number;
  end: number;
};

const PX_PER_HOUR = 22;
const MIN_BLOCK_WIDTH = 88;
const LANE_HEIGHT = 56;

const getComplexityStyle = (complexity: string) => {
  switch (complexity) {
    case 'CRITICAL':
      return {
        wrapper: 'border-red-500/60 bg-red-500/25',
        dot: 'bg-red-500',
        text: 'Crítica',
      };
    case 'HIGH':
      return {
        wrapper: 'border-orange-500/60 bg-orange-500/25',
        dot: 'bg-orange-500',
        text: 'Alta',
      };
    case 'MEDIUM':
      return {
        wrapper: 'border-yellow-500/60 bg-yellow-500/25',
        dot: 'bg-yellow-500',
        text: 'Média',
      };
    default:
      return {
        wrapper: 'border-emerald-500/60 bg-emerald-500/25',
        dot: 'bg-emerald-500',
        text: 'Baixa',
      };
  }
};

function topologicalSort(tasks: RoadmapTask[]) {
  const taskById = new Map(tasks.map((task) => [task.id, task]));
  const indegree = new Map<string, number>(tasks.map((task) => [task.id, 0]));
  const graph = new Map<string, string[]>();

  tasks.forEach((task) => {
    graph.set(task.id, []);
  });

  tasks.forEach((task) => {
    task.dependencies?.forEach((dependencyId) => {
      if (!taskById.has(dependencyId)) return;
      graph.get(dependencyId)?.push(task.id);
      indegree.set(task.id, (indegree.get(task.id) || 0) + 1);
    });
  });

  const queue: string[] = [];
  indegree.forEach((value, key) => {
    if (value === 0) queue.push(key);
  });

  const sorted: string[] = [];
  while (queue.length > 0) {
    const current = queue.shift();
    if (!current) continue;

    sorted.push(current);
    graph.get(current)?.forEach((next) => {
      const nextInDegree = (indegree.get(next) || 0) - 1;
      indegree.set(next, nextInDegree);
      if (nextInDegree === 0) queue.push(next);
    });
  }

  if (sorted.length !== tasks.length) {
    const remaining = tasks
      .map((task) => task.id)
      .filter((id) => !sorted.includes(id))
      .sort((a, b) => a.localeCompare(b));

    return [...sorted, ...remaining];
  }

  return sorted;
}

function computeTimeline(tasks: RoadmapTask[]): TimelineBlock[] {
  const orderedIds = topologicalSort(tasks);
  const taskById = new Map(tasks.map((task) => [task.id, task]));

  const laneAvailability: number[] = [];
  const endByTask = new Map<string, number>();
  const blocks: TimelineBlock[] = [];

  for (const taskId of orderedIds) {
    const task = taskById.get(taskId);
    if (!task) continue;

    const dependencyEnd = Math.max(
      0,
      ...(task.dependencies || []).map((dependencyId) => endByTask.get(dependencyId) || 0)
    );

    let lane = 0;
    while ((laneAvailability[lane] ?? 0) > dependencyEnd) {
      lane += 1;
    }

    const duration = Math.max(1, task.estimated_hours || 1);
    const start = dependencyEnd;
    const end = start + duration;

    laneAvailability[lane] = end;
    endByTask.set(task.id, end);

    blocks.push({ task, lane, start, duration, end });
  }

  return blocks;
}

function downloadCanvasAsPng(canvas: HTMLCanvasElement, fileName: string) {
  canvas.toBlob((blob) => {
    if (!blob) {
      toast.error('Não foi possível gerar o PNG da timeline.');
      return;
    }

    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = fileName;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }, 'image/png');
}

export function GanttTimelineView({ projectId }: GanttTimelineViewProps) {
  const [tasks, setTasks] = useState<RoadmapTask[]>([]);
  const [loading, setLoading] = useState(true);
  const timelineRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await ProjectService.getTasks(projectId, 1, 500);
        setTasks((response.items as RoadmapTask[]) || []);
      } catch (error) {
        console.error('Failed to fetch tasks', error);
        toast.error('Não foi possível carregar os dados da timeline.');
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, [projectId]);

  const timeline = useMemo(() => computeTimeline(tasks), [tasks]);
  const maxEnd = useMemo(() => Math.max(8, ...timeline.map((block) => block.end), 0), [timeline]);
  const lanesCount = useMemo(() => Math.max(1, ...timeline.map((block) => block.lane + 1), 1), [timeline]);

  const handleExportPng = () => {
    if (timeline.length === 0) {
      toast.error('Não há tasks para exportar na timeline.');
      return;
    }

    const canvas = document.createElement('canvas');
    const width = Math.max(1200, maxEnd * PX_PER_HOUR + 360);
    const height = lanesCount * LANE_HEIGHT + 220;
    canvas.width = width;
    canvas.height = height;

    const context = canvas.getContext('2d');
    if (!context) {
      toast.error('Falha ao criar contexto de desenho para exportação.');
      return;
    }

    context.fillStyle = '#09090b';
    context.fillRect(0, 0, width, height);

    context.fillStyle = '#fafafa';
    context.font = '700 24px Inter, system-ui';
    context.fillText('Roadmap Timeline (Gantt)', 36, 42);

    context.fillStyle = '#a1a1aa';
    context.font = '500 14px Inter, system-ui';
    context.fillText('Escala relativa baseada em estimated_hours e dependências entre tasks', 36, 66);

    const axisY = 104;
    context.strokeStyle = '#27272a';
    context.lineWidth = 1;

    for (let step = 0; step <= maxEnd; step += 2) {
      const x = 260 + step * PX_PER_HOUR;
      context.beginPath();
      context.moveTo(x, axisY);
      context.lineTo(x, height - 32);
      context.stroke();

      context.fillStyle = '#71717a';
      context.font = '11px Inter, system-ui';
      context.fillText(`T+${step}h`, x - 14, axisY - 12);
    }

    timeline.forEach((block) => {
      const x = 260 + block.start * PX_PER_HOUR;
      const y = axisY + block.lane * LANE_HEIGHT + 6;
      const blockWidth = Math.max(MIN_BLOCK_WIDTH, block.duration * PX_PER_HOUR);
      const blockHeight = 34;

      const fill =
        block.task.complexity === 'CRITICAL'
          ? '#ef4444'
          : block.task.complexity === 'HIGH'
            ? '#f97316'
            : block.task.complexity === 'MEDIUM'
              ? '#eab308'
              : '#22c55e';

      context.fillStyle = fill;
      context.globalAlpha = 0.85;
      context.fillRect(x, y, blockWidth, blockHeight);
      context.globalAlpha = 1;

      context.strokeStyle = '#111827';
      context.strokeRect(x, y, blockWidth, blockHeight);

      context.fillStyle = '#f8fafc';
      context.font = '600 11px Inter, system-ui';
      context.fillText(block.task.title.slice(0, 28), x + 8, y + 14);
      context.font = '500 10px Inter, system-ui';
      context.fillText(`${block.duration}h • ${block.task.type}`, x + 8, y + 28);

      context.fillStyle = '#94a3b8';
      context.font = '10px Inter, system-ui';
      context.fillText(`L${block.lane + 1}`, 215, y + 20);
    });

    downloadCanvasAsPng(canvas, `roadmap-gantt-${projectId}.png`);
    toast.success('PNG da timeline exportado com sucesso.');
  };

  if (loading) {
    return <div className='h-64 animate-pulse rounded-xl bg-card/30' />;
  }

  return (
    <div className='space-y-4'>
      <div className='flex flex-wrap items-center justify-between gap-3 rounded-xl border border-border bg-card/50 p-4'>
        <div className='flex flex-wrap gap-4 text-sm text-muted-foreground'>
          {(['CRITICAL', 'HIGH', 'MEDIUM', 'LOW'] as const).map((complexity) => {
            const style = getComplexityStyle(complexity);
            return (
              <span key={complexity} className='inline-flex items-center gap-2'>
                <span className={`h-3 w-3 rounded-full ${style.dot}`} />
                {style.text}
              </span>
            );
          })}
        </div>

        <Button variant='outline' size='sm' onClick={handleExportPng} className='gap-2'>
          <Download className='h-4 w-4' /> Exportar PNG
        </Button>
      </div>

      <div ref={timelineRef} className='overflow-x-auto rounded-xl border border-border bg-card/30 p-6'>
        <div
          className='relative min-w-[980px]'
          style={{ height: `${Math.max(300, lanesCount * LANE_HEIGHT + 80)}px`, width: `${maxEnd * PX_PER_HOUR + 280}px` }}
        >
          {Array.from({ length: maxEnd + 1 }).map((_, index) => {
            const isMajor = index % 2 === 0;
            return (
              <div
                key={index}
                className={`absolute bottom-0 top-8 ${isMajor ? 'border-border/80' : 'border-border/35'} border-l`}
                style={{ left: `${180 + index * PX_PER_HOUR}px` }}
              >
                {isMajor && (
                  <span className='absolute -top-6 -translate-x-1/2 text-[10px] text-muted-foreground'>
                    T+{index}h
                  </span>
                )}
              </div>
            );
          })}

          {timeline.map((block) => {
            const style = getComplexityStyle(block.task.complexity);
            return (
              <div
                key={block.task.id}
                className={`absolute rounded-md border px-3 py-2 shadow-sm transition-all hover:scale-[1.01] ${style.wrapper}`}
                style={{
                  left: `${180 + block.start * PX_PER_HOUR}px`,
                  top: `${42 + block.lane * LANE_HEIGHT}px`,
                  width: `${Math.max(MIN_BLOCK_WIDTH, block.duration * PX_PER_HOUR)}px`,
                  height: '42px',
                }}
                title={`${block.task.title}\n${block.duration}h • ${block.task.type}`}
              >
                <p className='truncate text-xs font-semibold text-foreground'>{block.task.title}</p>
                <p className='text-[10px] text-muted-foreground'>
                  {block.duration}h • {block.task.type} • L{block.lane + 1}
                </p>
              </div>
            );
          })}

          {timeline.length === 0 && (
            <div className='flex h-full items-center justify-center text-sm text-muted-foreground'>
              Nenhuma task disponível para montar a timeline.
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
