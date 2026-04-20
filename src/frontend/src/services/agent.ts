import { api } from "./api";
import { AgentListResponse, Agent } from "../types/agent";

export const AgentService = {
  getAgents: async (page = 1, size = 10): Promise<AgentListResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      size: size.toString(),
    });
    const response = await api.get("/agents", { params });
    return response.data;
  },

  getAgentById: async (id: string): Promise<Agent> => {
    const response = await api.get(`/agents/${id}`);
    return response.data;
  }
};
