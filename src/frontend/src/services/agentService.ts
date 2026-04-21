import { api } from "./api";
import { Agent, AgentListResponse } from "../types/agent";

export const agentService = {
  async getAgents(params?: {
    page?: number;
    size?: number;
    enabled?: boolean;
    skill?: string;
    provider?: string;
    search?: string;
  }): Promise<AgentListResponse> {
    const response = await api.get("/agents", { params });
    const payload = response.data;

    if (Array.isArray(payload)) {
      return {
        items: payload,
        total: payload.length,
        page: params?.page ?? 1,
        size: params?.size ?? payload.length,
        pages: 1,
      };
    }

    return {
      items: payload?.items ?? [],
      total: payload?.total ?? payload?.items?.length ?? 0,
      page: payload?.page ?? params?.page ?? 1,
      size: payload?.size ?? params?.size ?? payload?.items?.length ?? 0,
      pages: payload?.pages ?? 1,
    };
  },

  async getAgentById(id: string): Promise<Agent> {
    const response = await api.get(`/agents/${id}`);
    return response.data;
  },

  async createAgent(agentData: Partial<Agent>): Promise<Agent> {
    const response = await api.post("/agents", agentData);
    return response.data;
  },

  async updateAgent(id: string, agentData: Partial<Agent>): Promise<Agent> {
    const response = await api.put(`/agents/${id}`, agentData);
    return response.data;
  },

  async deleteAgent(id: string): Promise<void> {
    await api.delete(`/agents/${id}`);
  },

  async testConnection(id: string): Promise<{ success: boolean; message: string }> {
    const response = await api.post(`/agents/${id}/test`);
    return response.data;
  },

  async testConfiguration(agentData: Partial<Agent>): Promise<{ success: boolean; response: string }> {
    const response = await api.post('/agents/test-config', agentData);
    return response.data;
  }
};
