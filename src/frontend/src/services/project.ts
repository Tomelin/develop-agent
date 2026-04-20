import { api } from "./api";
import { Project, ProjectListResponse, ProjectCreateRequest, ProjectStatus, FlowType } from "../types/project";
import { Task, TaskListResponse, TaskStatus } from "../types/task";

export const ProjectService = {
  getProjects: async (page = 1, size = 10, status?: ProjectStatus, flow_type?: FlowType): Promise<ProjectListResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      size: size.toString(),
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
      size: size.toString(),
    });
    const response = await api.get(`/projects/${projectId}/tasks`, { params });
    return response.data;
  },

  updateTaskStatus: async (projectId: string, taskId: string, status: TaskStatus): Promise<Task> => {
    const response = await api.put(`/projects/${projectId}/tasks/${taskId}/status`, { status });
    return response.data;
  }
};
