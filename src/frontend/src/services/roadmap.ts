import { api } from './api';
import { RoadmapData, RoadmapSummary } from '../types/roadmap';

export const RoadmapService = {
  getRoadmapSummary: async (projectId: string): Promise<RoadmapSummary> => {
    const response = await api.get(`/projects/${projectId}/roadmap/summary`);
    return response.data;
  },

  getRoadmapData: async (projectId: string): Promise<RoadmapData> => {
    const response = await api.get(`/projects/${projectId}/roadmap`);
    return response.data;
  },

  getExportUrl: (projectId: string, format: string): string => {
    const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
    return `${baseUrl}/projects/${projectId}/roadmap/export?format=${format}`;
  },
};
