import { api } from "./api";
import { Project, ProjectListResponse, ProjectCreateRequest, ProjectStatus, FlowType } from "../types/project";
import { TaskListResponse, TaskStatus } from "../types/task";
import { Phase6AnalyzeCoverageResponse, Phase6ValidationResult } from "../types/phase6";
import { Phase5CodeContext, Phase5CodeFile, Phase5ExecutionMode, Phase5FileListResponse, Phase5Summary } from "../types/phase5";

const FLOW_TO_API: Record<FlowType, "SOFTWARE" | "LANDING_PAGE" | "MARKETING"> = {
  A: "SOFTWARE",
  B: "LANDING_PAGE",
  C: "MARKETING",
};

const FLOW_FROM_API: Record<string, FlowType> = {
  A: "A",
  B: "B",
  C: "C",
  SOFTWARE: "A",
  LANDING_PAGE: "B",
  MARKETING: "C",
};

const normalizeProject = (raw: Record<string, unknown>): Project => ({
  id: String(raw.id ?? ""),
  name: String(raw.name ?? ""),
  description: String(raw.description ?? ""),
  status: (raw.status as ProjectStatus) ?? "DRAFT",
  flow_type: FLOW_FROM_API[String(raw.flow_type ?? "A")] ?? "A",
  linked_project_id: raw.linked_project_id ? String(raw.linked_project_id) : undefined,
  current_phase: Number(raw.current_phase ?? raw.current_phase_number ?? 1),
  progress_percentage: Number(raw.progress_percentage ?? 0),
  tokens_used: Number(raw.tokens_used ?? raw.total_tokens_used ?? 0),
  dynamic_mode: Boolean(raw.dynamic_mode ?? raw.dynamic_mode_enabled ?? false),
  created_at: String(raw.created_at ?? ""),
  updated_at: String(raw.updated_at ?? ""),
});

export const ProjectService = {
  getProjects: async (page = 1, size = 10, status?: ProjectStatus, flow_type?: FlowType): Promise<ProjectListResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: size.toString(),
    });
    if (status) params.append("status", status);
    if (flow_type) params.append("flow_type", flow_type);

    const response = await api.get("/projects", { params });
    const payload = response.data;
    const items = (payload?.items ?? []).map((item: Record<string, unknown>) => normalizeProject(item));
    const total = Number(payload?.total ?? items.length);
    const currentPage = Number(payload?.page ?? page);
    const currentSize = Number(payload?.size ?? payload?.limit ?? size);
    const pages = Number(payload?.pages ?? Math.max(1, Math.ceil(total / (currentSize || 1))));

    return { items, total, page: currentPage, size: currentSize, pages };
  },

  getProjectById: async (id: string): Promise<Project> => {
    const response = await api.get(`/projects/${id}`);
    return normalizeProject(response.data);
  },

  createProject: async (data: ProjectCreateRequest): Promise<Project> => {
    const payload = {
      name: data.name,
      description: data.description,
      flow_type: FLOW_TO_API[data.flow_type],
      linked_project_id: data.linked_project_id,
      dynamic_mode_enabled: data.dynamic_mode,
    };
    const response = await api.post("/projects", payload);
    return normalizeProject(response.data);
  },

  updateProject: async (id: string, data: Partial<Project>): Promise<Project> => {
    const response = await api.put(`/projects/${id}`, data);
    return normalizeProject(response.data);
  },

  pauseProject: async (id: string): Promise<void> => {
    await api.post(`/projects/${id}/pause`);
  },

  resumeProject: async (id: string): Promise<void> => {
    await api.post(`/projects/${id}/resume`);
  },

  archiveProject: async (id: string): Promise<void> => {
    await api.post(`/projects/${id}/archive`);
  },

  getTasks: async (projectId: string, page = 1, size = 10): Promise<TaskListResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: size.toString(),
    });
    const response = await api.get(`/projects/${projectId}/tasks`, { params });
    return response.data;
  },

  updateTaskStatus: async (projectId: string, taskId: string, status: TaskStatus): Promise<void> => {
    await api.put(`/projects/${projectId}/tasks/${taskId}/status`, { status });
  },

  analyzePhase6Coverage: async (projectId: string, payload: { backend_dir: string; threshold?: number }): Promise<Phase6AnalyzeCoverageResponse> => {
    const response = await api.post(`/projects/${projectId}/phases/6/analyze-coverage`, payload);
    return response.data;
  },

  validatePhase6Tests: async (projectId: string, payload: { backend_dir: string; frontend_dir: string }): Promise<Phase6ValidationResult> => {
    const response = await api.post(`/projects/${projectId}/phases/6/validate-tests`, payload);
    return response.data;
  },

  setPhase5Mode: async (projectId: string, mode: Phase5ExecutionMode): Promise<{ mode: Phase5ExecutionMode }> => {
    const response = await api.post(`/projects/${projectId}/phases/5/mode`, { mode });
    return response.data;
  },

  executeAllPhase5Tasks: async (projectId: string): Promise<{ executed_tasks: number }> => {
    const response = await api.post(`/projects/${projectId}/phases/5/execute`);
    return response.data;
  },

  executePhase5Task: async (projectId: string, taskId: string): Promise<void> => {
    await api.post(`/projects/${projectId}/phases/5/tasks/${taskId}/execute`);
  },

  getPhase5Summary: async (projectId: string): Promise<Phase5Summary> => {
    const response = await api.get(`/projects/${projectId}/phases/5/summary`);
    return response.data;
  },

  getPhase5CodeContext: async (projectId: string): Promise<Phase5CodeContext> => {
    const response = await api.get(`/projects/${projectId}/phases/5/code-context`);
    return response.data;
  },

  getProjectFiles: async (projectId: string): Promise<Phase5CodeFile[]> => {
    const response = await api.get<Phase5FileListResponse>(`/projects/${projectId}/files`);
    return response.data.items;
  },

  getProjectFileById: async (projectId: string, fileId: string): Promise<Phase5CodeFile> => {
    const response = await api.get(`/projects/${projectId}/files/${fileId}`);
    return response.data;
  },

  downloadProjectFilesZip: async (projectId: string): Promise<Blob> => {
    const response = await api.get(`/projects/${projectId}/files/download`, { responseType: "blob" });
    return response.data;
  },
};
