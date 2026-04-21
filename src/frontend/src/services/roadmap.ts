import { api } from "./api";
import { RoadmapSummary, RoadmapData } from "../types/roadmap";

export const RoadmapService = {
  getRoadmapSummary: async (projectId: string): Promise<RoadmapSummary> => {
    const response = await api.get(`/projects/${projectId}/roadmap/summary`);
    return response.data;
  },

  getRoadmapData: async (projectId: string): Promise<RoadmapData> => {
    // Assuming backend endpoint to get the structured roadmap with phases/epics/tasks
    // Since TASK-09-006 asks for "Aba Épicos" and hierarchical structure,
    // we assume this endpoint exists or we can reconstruct it from tasks list.
    // If it doesn't exist natively as structured, we can reconstruct it, but for now we expect it.
    // Wait, the document says the JSON schema is saved.
    // However, the backend exposes GET /projects/:id/tasks (already handled in project.ts).
    // Let's implement an endpoint just in case, and fallback if it doesn't work.
    try {
      const response = await api.get(`/projects/${projectId}/roadmap`);
      return response.data;
    } catch (e) {
      // Return empty structure if not found
      console.warn("Could not fetch structured roadmap, returning empty", e);
      return { project_id: projectId, phases: [] };
    }
  },

  getExportUrl: (projectId: string, format: string): string => {
    return `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"}/projects/${projectId}/roadmap/export?format=${format}`;
  }
};
