import { api } from "./api";
import {
  AdminPlatformSettings,
  CostPreviewResponse,
  DiversityMetrics,
  DynamicModeConfig,
  DynamicModePreview,
  FeatureFlag,
  PhaseAgentMatrixResponse,
  TriadSelection,
  TriadSelectionLog,
} from "@/types/phase17";

export const Phase17Service = {
  getTriadSelections: async (projectId: string): Promise<TriadSelection[]> => {
    const response = await api.get(`/projects/${projectId}/triad-selections`);
    return response.data.items ?? response.data;
  },

  getSelectionLogs: async (projectId: string, phaseNumber?: number): Promise<TriadSelectionLog[]> => {
    const response = await api.get(`/projects/${projectId}/selection-logs`, {
      params: phaseNumber ? { phase: phaseNumber } : undefined,
    });
    return response.data.items ?? response.data;
  },

  updateDynamicMode: async (projectId: string, payload: DynamicModeConfig): Promise<{ enabled: boolean }> => {
    const response = await api.put(`/projects/${projectId}/dynamic-mode`, payload);
    return response.data;
  },

  previewDynamicSelection: async (projectId: string): Promise<DynamicModePreview> => {
    const response = await api.get(`/projects/${projectId}/dynamic-mode/preview`);
    return response.data;
  },

  getDiversityMetrics: async (projectId: string): Promise<DiversityMetrics> => {
    const response = await api.get(`/projects/${projectId}/diversity-metrics`);
    return response.data;
  },

  getAgentMatrix: async (projectId: string): Promise<PhaseAgentMatrixResponse> => {
    const response = await api.get(`/projects/${projectId}/agent-config/matrix`);
    return response.data;
  },

  updateAgentMatrix: async (projectId: string, payload: PhaseAgentMatrixResponse): Promise<PhaseAgentMatrixResponse> => {
    const response = await api.put(`/projects/${projectId}/agent-config/matrix`, payload);
    return response.data;
  },

  previewConfigurationCost: async (projectId: string, payload: PhaseAgentMatrixResponse): Promise<CostPreviewResponse> => {
    const response = await api.post(`/projects/${projectId}/agent-config/cost-preview`, payload);
    return response.data;
  },

  getAdminSettings: async (): Promise<AdminPlatformSettings> => {
    const response = await api.get("/admin/settings");
    return response.data;
  },

  saveAdminSettings: async (payload: AdminPlatformSettings): Promise<AdminPlatformSettings> => {
    const response = await api.put("/admin/settings", payload);
    return response.data;
  },

  getFeatureFlags: async (): Promise<FeatureFlag[]> => {
    const response = await api.get("/admin/feature-flags");
    return response.data.items ?? response.data;
  },

  updateFeatureFlags: async (flags: FeatureFlag[]): Promise<FeatureFlag[]> => {
    const response = await api.put("/admin/feature-flags", { flags });
    return response.data.items ?? response.data;
  },

  getFeatureFlagsPublic: async (): Promise<FeatureFlag[]> => {
    const response = await api.get("/feature-flags");
    return response.data.items ?? response.data;
  },
};
