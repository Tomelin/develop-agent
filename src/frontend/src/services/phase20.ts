import { api } from "./api";
import { isAxiosError } from "axios";
import {
  CollaboratorRole,
  IntegrationConnectionState,
  OrganizationMember,
  OrganizationRole,
  PricingPlan,
  ProjectCollaborator,
  PublicRoadmapResponse,
  PublicTemplate,
  PublicTemplateListResponse,
} from "@/types/phase20";

interface PublicTemplateQuery {
  search?: string;
  category?: string;
  group?: string;
  page?: number;
  size?: number;
}

export const Phase20Service = {
  listOrganizationMembers: async (): Promise<OrganizationMember[]> => {
    const response = await api.get("/org/members");
    return response.data.items ?? response.data;
  },

  inviteMember: async (email: string, role: OrganizationRole): Promise<void> => {
    await api.post("/org/invite", { email, role });
  },

  updateOrganizationMemberRole: async (userId: string, role: OrganizationRole): Promise<void> => {
    await api.put(`/org/members/${userId}/role`, { role });
  },

  removeOrganizationMember: async (userId: string): Promise<void> => {
    await api.delete(`/org/members/${userId}`);
  },

  listProjectCollaborators: async (projectId: string): Promise<ProjectCollaborator[]> => {
    try {
      const response = await api.get(`/projects/${projectId}/collaborators`);
      return response.data.items ?? response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return [];
      throw error;
    }
  },

  addProjectCollaborator: async (projectId: string, email: string, role: CollaboratorRole): Promise<void> => {
    try {
      await api.post(`/projects/${projectId}/collaborators`, { email, role });
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  updateProjectCollaboratorRole: async (projectId: string, userId: string, role: CollaboratorRole): Promise<void> => {
    try {
      await api.put(`/projects/${projectId}/collaborators/${userId}/role`, { role });
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  removeProjectCollaborator: async (projectId: string, userId: string): Promise<void> => {
    try {
      await api.delete(`/projects/${projectId}/collaborators/${userId}`);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  listPublicTemplates: async (params?: PublicTemplateQuery): Promise<PublicTemplateListResponse> => {
    try {
      const response = await api.get("/marketplace/templates", { params });
      return {
        items: response.data.items ?? response.data,
        total: response.data.total ?? (response.data.items?.length || 0),
        page: response.data.page ?? 1,
        size: response.data.size ?? 20,
      };
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return { items: [], total: 0, page: params?.page ?? 1, size: params?.size ?? 20 };
      }
      throw error;
    }
  },

  publishTemplate: async (payload: Pick<PublicTemplate, "title" | "description" | "category" | "content" | "group" | "tags">): Promise<void> => {
    try {
      await api.post("/marketplace/templates", payload);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  activateTemplate: async (templateId: string): Promise<void> => {
    try {
      await api.post(`/marketplace/templates/${templateId}/use`);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  starTemplate: async (templateId: string): Promise<void> => {
    try {
      await api.post(`/marketplace/templates/${templateId}/star`);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  getIntegrationsStatus: async (projectId: string): Promise<IntegrationConnectionState[]> => {
    try {
      const response = await api.get(`/projects/${projectId}/integrations`);
      return response.data.items ?? response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return [];
      throw error;
    }
  },

  getGithubAuthUrl: async (): Promise<{ auth_url: string }> => {
    try {
      const response = await api.get("/integrations/github/auth");
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return { auth_url: "" };
      throw error;
    }
  },

  configureJiraIntegration: async (payload: { base_url: string; email: string; api_token: string; project_key: string }): Promise<void> => {
    try {
      await api.post("/integrations/jira", payload);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  syncProjectToJira: async (projectId: string): Promise<void> => {
    try {
      await api.post(`/projects/${projectId}/integrations/jira/sync`);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  configureSlackWebhook: async (payload: { webhook_url: string; channel: string }): Promise<void> => {
    try {
      await api.post("/integrations/slack/webhook", payload);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  getPricingPlans: async (): Promise<PricingPlan[]> => {
    try {
      const response = await api.get("/pricing/plans");
      return response.data.items ?? response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return [];
      throw error;
    }
  },

  createStripeCheckout: async (planCode: string): Promise<{ checkout_url: string }> => {
    try {
      const response = await api.post("/pricing/checkout", { plan_code: planCode });
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return { checkout_url: "" };
      throw error;
    }
  },

  getPublicRoadmap: async (): Promise<PublicRoadmapResponse> => {
    try {
      const response = await api.get("/roadmap/public");
      return response.data;
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return { vision: "", features: [], changelog: [] };
      throw error;
    }
  },

  voteRoadmapFeature: async (featureId: string): Promise<void> => {
    try {
      await api.post(`/roadmap/features/${featureId}/vote`);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  suggestRoadmapFeature: async (payload: { title: string; description: string }): Promise<void> => {
    try {
      await api.post("/roadmap/features/suggestions", payload);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },

  updateRoadmapFeatureStatus: async (featureId: string, status: string): Promise<void> => {
    try {
      await api.put(`/admin/roadmap/features/${featureId}/status`, { status });
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) return;
      throw error;
    }
  },
};
