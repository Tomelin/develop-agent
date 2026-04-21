import { api } from "./api";
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
    const response = await api.get(`/projects/${projectId}/collaborators`);
    return response.data.items ?? response.data;
  },

  addProjectCollaborator: async (projectId: string, email: string, role: CollaboratorRole): Promise<void> => {
    await api.post(`/projects/${projectId}/collaborators`, { email, role });
  },

  updateProjectCollaboratorRole: async (projectId: string, userId: string, role: CollaboratorRole): Promise<void> => {
    await api.put(`/projects/${projectId}/collaborators/${userId}/role`, { role });
  },

  removeProjectCollaborator: async (projectId: string, userId: string): Promise<void> => {
    await api.delete(`/projects/${projectId}/collaborators/${userId}`);
  },

  listPublicTemplates: async (params?: PublicTemplateQuery): Promise<PublicTemplateListResponse> => {
    const response = await api.get("/marketplace/templates", { params });
    return {
      items: response.data.items ?? response.data,
      total: response.data.total ?? (response.data.items?.length || 0),
      page: response.data.page ?? 1,
      size: response.data.size ?? 20,
    };
  },

  publishTemplate: async (payload: Pick<PublicTemplate, "title" | "description" | "category" | "content" | "group" | "tags">): Promise<void> => {
    await api.post("/marketplace/templates", payload);
  },

  activateTemplate: async (templateId: string): Promise<void> => {
    await api.post(`/marketplace/templates/${templateId}/use`);
  },

  starTemplate: async (templateId: string): Promise<void> => {
    await api.post(`/marketplace/templates/${templateId}/star`);
  },

  getIntegrationsStatus: async (projectId: string): Promise<IntegrationConnectionState[]> => {
    const response = await api.get(`/projects/${projectId}/integrations`);
    return response.data.items ?? response.data;
  },

  getGithubAuthUrl: async (): Promise<{ auth_url: string }> => {
    const response = await api.get("/integrations/github/auth");
    return response.data;
  },

  configureJiraIntegration: async (payload: { base_url: string; email: string; api_token: string; project_key: string }): Promise<void> => {
    await api.post("/integrations/jira", payload);
  },

  syncProjectToJira: async (projectId: string): Promise<void> => {
    await api.post(`/projects/${projectId}/integrations/jira/sync`);
  },

  configureSlackWebhook: async (payload: { webhook_url: string; channel: string }): Promise<void> => {
    await api.post("/integrations/slack/webhook", payload);
  },

  getPricingPlans: async (): Promise<PricingPlan[]> => {
    const response = await api.get("/pricing/plans");
    return response.data.items ?? response.data;
  },

  createStripeCheckout: async (planCode: string): Promise<{ checkout_url: string }> => {
    const response = await api.post("/pricing/checkout", { plan_code: planCode });
    return response.data;
  },

  getPublicRoadmap: async (): Promise<PublicRoadmapResponse> => {
    const response = await api.get("/roadmap/public");
    return response.data;
  },

  voteRoadmapFeature: async (featureId: string): Promise<void> => {
    await api.post(`/roadmap/features/${featureId}/vote`);
  },

  suggestRoadmapFeature: async (payload: { title: string; description: string }): Promise<void> => {
    await api.post("/roadmap/features/suggestions", payload);
  },

  updateRoadmapFeatureStatus: async (featureId: string, status: string): Promise<void> => {
    await api.put(`/admin/roadmap/features/${featureId}/status`, { status });
  },
};
