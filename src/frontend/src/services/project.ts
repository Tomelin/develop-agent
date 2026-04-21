import { api } from "./api";
import { Project, ProjectListResponse, ProjectCreateRequest, ProjectStatus, FlowType } from "../types/project";
import { TaskListResponse, TaskStatus } from "../types/task";
import { Phase6AnalyzeCoverageResponse, Phase6ValidationResult } from "../types/phase6";
import { Phase5CodeContext, Phase5CodeFile, Phase5ExecutionMode, Phase5FileListResponse, Phase5Summary } from "../types/phase5";

export const ProjectService = {
  getProjects: async (page = 1, size = 10, status?: ProjectStatus, flow_type?: FlowType): Promise<ProjectListResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: size.toString(),
    });
    if (status) params.append("status", status);
    if (flow_type) params.append("flow_type", flow_type);

    const response = await api.get("/projects", { params });
    return response.data;
  },

  getProjectById: async (id: string): Promise<Project> => {
    const response = await api.get(`/projects/${id}`);
    return response.data;
  },

  createProject: async (data: ProjectCreateRequest): Promise<Project> => {
    const response = await api.post("/projects", data);
    return response.data;
  },

  updateProject: async (id: string, data: Partial<Project>): Promise<Project> => {
    const response = await api.put(`/projects/${id}`, data);
    return response.data;
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
