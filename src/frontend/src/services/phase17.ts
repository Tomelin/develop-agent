import { api } from "./api";
import { isAxiosError } from "axios";
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
    try {
      const response = await api.get(`/projects/${projectId}/triad-selections`);
      return response.data.items ?? response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return [];
      throw error;
    }
  },

  getSelectionLogs: async (projectId: string, phaseNumber?: number): Promise<TriadSelectionLog[]> => {
    try {
      const response = await api.get(`/projects/${projectId}/selection-logs`, {
        params: phaseNumber ? { phase: phaseNumber } : undefined,
      });
      return response.data.items ?? response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return [];
      throw error;
    }
  },

  updateDynamicMode: async (projectId: string, payload: DynamicModeConfig): Promise<{ enabled: boolean }> => {
    try {
      const response = await api.put(`/projects/${projectId}/dynamic-mode`, payload);
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return { enabled: payload.enabled };
      throw error;
    }
  },

  previewDynamicSelection: async (projectId: string): Promise<DynamicModePreview> => {
    try {
      const response = await api.get(`/projects/${projectId}/dynamic-mode/preview`);
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return {
          eligible_agents: 0,
          triad: {
            producer: { id: "", name: "", provider: "OPENAI", model: "" },
            reviewer: { id: "", name: "", provider: "OPENAI", model: "" },
            refiner: { id: "", name: "", provider: "OPENAI", model: "" },
          },
          notes: [],
        };
      }
      throw error;
    }
  },

  getDiversityMetrics: async (projectId: string): Promise<DiversityMetrics> => {
    try {
      const response = await api.get(`/projects/${projectId}/diversity-metrics`);
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return {
          project_id: projectId,
          diversity_score: 0,
          providers: [],
          models: [],
          full_diversity_triads: 0,
          repeated_provider_triads: 0,
          role_distribution: [],
        };
      }
      throw error;
    }
  },

  getAgentMatrix: async (projectId: string): Promise<PhaseAgentMatrixResponse> => {
    try {
      const response = await api.get(`/projects/${projectId}/agent-config/matrix`);
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return { rows: [] };
      }
      throw error;
    }
  },

  updateAgentMatrix: async (projectId: string, payload: PhaseAgentMatrixResponse): Promise<PhaseAgentMatrixResponse> => {
    try {
      const response = await api.put(`/projects/${projectId}/agent-config/matrix`, payload);
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return payload;
      throw error;
    }
  },

  previewConfigurationCost: async (projectId: string, payload: PhaseAgentMatrixResponse): Promise<CostPreviewResponse> => {
    try {
      const response = await api.post(`/projects/${projectId}/agent-config/cost-preview`, payload);
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return { monthly_estimated_usd: 0, note: "Cost preview indisponível no backend atual." };
      }
      throw error;
    }
  },

  getAdminSettings: async (): Promise<AdminPlatformSettings> => {
    try {
      const response = await api.get("/admin/settings");
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return {
          workers: { max_concurrency: 1, agent_timeout_seconds: 60, triad_timeout_seconds: 120 },
          models: { default_model: "", spec_generation_model: "" },
          limits: { max_projects_per_user: 1, max_parallel_phases_per_user: 1, max_spec_tokens: 4000 },
          retry: { max_attempts: 1, backoff_seconds: 1 },
        };
      }
      throw error;
    }
  },

  saveAdminSettings: async (payload: AdminPlatformSettings): Promise<AdminPlatformSettings> => {
    try {
      const response = await api.put("/admin/settings", payload);
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return payload;
      throw error;
    }
  },

  getFeatureFlags: async (): Promise<FeatureFlag[]> => {
    try {
      const response = await api.get("/admin/feature-flags");
      return response.data.items ?? response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return [];
      throw error;
    }
  },

  updateFeatureFlags: async (flags: FeatureFlag[]): Promise<FeatureFlag[]> => {
    try {
      const response = await api.put("/admin/feature-flags", { flags });
      return response.data.items ?? response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return flags;
      throw error;
    }
  },

  getFeatureFlagsPublic: async (): Promise<FeatureFlag[]> => {
    try {
      const response = await api.get("/feature-flags");
      return response.data.items ?? response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return [];
      throw error;
    }
  },
};
